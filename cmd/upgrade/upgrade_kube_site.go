package main

import (
	"context"

	"fmt"
	"os"
	"path/filepath"

	"github.com/skupperproject/skupper/api/types"

	"github.com/skupperproject/skupper/pkg/apis/skupper/v2alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/skupperproject/skupper/internal/cmd/skupper/common/utils"

	"github.com/skupperproject/skupper/internal/kube/client"
	"github.com/skupperproject/skupper/pkg/site"
)

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

func upgradeSite(siteConfig *types.SiteConfig, outputPath string) error {
	resource, err := createSiteCR(siteConfig)
	if err != nil {
		return fmt.Errorf("Error creating site CR: %w", err.Error())
	}
	err = saveSiteCR(resource, outputPath, siteConfig.Spec.SkupperName)
	if err != nil {
		return fmt.Errorf("Error saving site CR: %w", err.Error())
	}
	return nil
}

func v1IsLinkAccessDefault(siteConfig *types.SiteConfig) bool {
	// TODO: if policy CRD is present (in the site), then check policy.allowIncomingLinks
	return true
}

func createSiteCR(siteConfig *types.SiteConfig) (*v2alpha1.Site, error) {
	// TODO: assume that service account should be unspecified during upgrade.
	//DefaultServiceAccountName = "skupper-controller"
	//options := map[string]string{
	//site.SiteConfigNameKey: siteConfig.Spec.SkupperName,
	//}

	resource := &v2alpha1.Site{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "skupper.io/v2alpha1",
			Kind:       "Site",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      siteConfig.Spec.SkupperName,
			Namespace: siteConfig.Spec.SkupperNamespace, // TODO populating namespace seems unnecessary.  Confirm.
		},
		Spec: v2alpha1.SiteSpec{
			//Settings:       options,
			//ServiceAccount: DefaultServiceAccountName,
		},
	}
	if siteConfig.Spec.RouterMode == string(types.TransportModeEdge) {
		resource.Spec.RouterMode = string(types.TransportModeEdge)
	}
	// TODO: confirm logic to set LinkAccess
	if v1IsLinkAccessDefault(siteConfig) {
		resource.Spec.LinkAccess = "default"
	}

	return resource, nil
}

func saveSiteCR(resource *v2alpha1.Site, outputPath string, siteName string) error {
	targetDir := filepath.Join(outputPath, siteName)
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("Failed to create directory: %v", err)
	}

	filepath := filepath.Join(targetDir, fmt.Sprintf("%s_site.yaml", siteName))

	data, err := utils.Encode("yaml", resource)
	if err != nil {
		return fmt.Errorf("Failed to marshal site resource to YAML: %w", err.Error())
	}

	err = os.WriteFile(filepath, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("Failed to write site resource to file: %w", err.Error())
	}
	fmt.Printf("Wrote site CR to file: %s\n", filepath)
	return nil
}
