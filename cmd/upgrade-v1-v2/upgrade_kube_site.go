package main

import (
	"context"

	"fmt"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"

	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/skupperproject/skupper/api/types"
	"github.com/skupperproject/skupper/internal/kube/client"

	//"log"
	//"github.com/davecgh/go-spew/spew"

	"github.com/skupperproject/skupper/pkg/site"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/skupperproject/skupper/pkg/apis/skupper/v1alpha1"
)

func getKubeConfig() string {
	var kubeconfig string // TODO
	return kubeconfig
}

func getUidToSiteConfig(cli *client.KubeClient) (map[string]*types.SiteConfig, error) {
	uidToSiteConfig := map[string]*types.SiteConfig{}

	kubeClient := cli.GetKubeClient()

	// for every namespace
	namespaces, err := kubeClient.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Error getting namepaces: %w", err.Error())
	}

	for _, ns := range namespaces.Items {
		nsName := ns.Name

		// read site configmap from the namespace
		cm, err := readConfigMap(context.Background(), nsName, types.SiteConfigMapName, cli)
		if err != nil {
			return nil, fmt.Errorf("TMPDBG: error reading configmap in namespace %s: %w", nsName, err.Error())
		}
		if cm == nil {
			continue
		}

		siteConfig, err := site.ReadSiteConfig(cm, nsName)
		if err != nil {
			return nil, fmt.Errorf("Error reading siteConfig: %w", err.Error())
		}
		uidToSiteConfig[siteConfig.Reference.UID] = siteConfig
	}

	return uidToSiteConfig, nil
}

func readConfigMap(ctx context.Context, namespace string, name string, cli *client.KubeClient) (*corev1.ConfigMap, error) {
	kubeClient := cli.GetKubeClient()
	cm, err := kubeClient.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return cm, err
}

func upgradeConfigMapSkupperSites(cli *client.KubeClient, outputPath string, uidToSiteConfig map[string]*types.SiteConfig) error {
	//log.Printf("TMPDBG: upgradeSites: uidToSiteConfigs=%+v\n", spew.Sdump(uidToSiteConfig))

	for _, siteConfig := range uidToSiteConfig {
		//log.Printf("TMPDBG: upgradeSites: uid=%+v\n", spew.Sdump(uid))
		//log.Printf("TMPDBG: upgradeSites: siteConfig=%+v\n", spew.Sdump(siteConfig))

		dirPath := filepath.Join(outputPath, siteConfig.Spec.SkupperNamespace)

		err := upgradeConfigMapSkupperSite(cli, siteConfig.Spec.SkupperNamespace, dirPath)
		if err != nil {
			return fmt.Errorf("Error upgrading skupper site configmap: %w", err.Error())
		}
	}

	return nil
}

func upgradeConfigMapSkupperSite(cli *client.KubeClient, namespace string, outputPath string) error {
	kubeconfig := getKubeConfig()

	cli, err := client.NewClient(namespace, "", kubeconfig)
	if err != nil {
		return fmt.Errorf("Error getting van client: %w", err.Error())
	}

	cm, err := readConfigMap(context.Background(), namespace, types.SiteConfigMapName, cli)
	if err != nil {
		return fmt.Errorf("Error reading configmap skupper-site: %w", err.Error())
	}

	siteConfig, err := site.ReadSiteConfig(cm, namespace)
	if err != nil {
		return fmt.Errorf("Error reading siteConfig: %w", err.Error())
	}

	//log.Printf("TMPDBG: upgradeConfigMapSkupperSite: spew.Sdump(siteConfig)=%+v\n", spew.Sdump(siteConfig))

	// TODO: Should probably validate similar to CmdSiteCreate.ValidateInput
	// in internal/cmd/skupper/site/kube/site_create.go

	resource, err := v1ToV2SkupperSite(siteConfig)
	if err != nil {
		return fmt.Errorf("Error in v1ToV2SkupperSite: %w", err.Error())
	}
	//log.Printf("TMPDBG: upgradeConfigMapSkupperSite: v2 CR resource=%+v\n", spew.Sdump(resource))

	filepath := filepath.Join(outputPath, fmt.Sprintf("%s-site.yaml", siteConfig.Spec.SkupperName))

	data, err := yaml.Marshal(resource)
	if err != nil {
		return fmt.Errorf("Failed to marshal site resource to YAML: %w", err.Error())
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("Failed to write site resource to file: %w", err.Error())
	}

	// TMPDBG remove me
	// siteConfig, err := getUidToSiteConfig(cli)
	// log.Printf("TMPDBG: upgradeConfigMapSkupperSite: siteConfigs=%+v\n", spew.Sdump(siteConfigs))
	// log.Printf("TMPDBG: upgradeConfigMapSkupperSite: err=%+v\n", err)

	return nil
}

func v1IsLinkAccessDefault(siteConfig *types.SiteConfig) bool {
	// TODO: if policy CRD is present (in the site), then check policy.allowIncomingLinks
	return true
}

func v1ToV2SkupperSite(siteConfig *types.SiteConfig) (*v1alpha1.Site, error) {
	//DefaultServiceAccountName = "skupper-controller"
	//options := map[string]string{
	//site.SiteConfigNameKey: siteConfig.Spec.SkupperName,
	//}

	resource := &v1alpha1.Site{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "skupper.io/v1alpha1",
			Kind:       "Site",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      siteConfig.Spec.SkupperName,
			Namespace: siteConfig.Spec.SkupperNamespace,
		},
		Spec: v1alpha1.SiteSpec{
			//Settings:       options,
			//ServiceAccount: DefaultServiceAccountName,
		},
	}
	if v1IsLinkAccessDefault(siteConfig) {
		resource.Spec.LinkAccess = "default"
	}

	return resource, nil
}
