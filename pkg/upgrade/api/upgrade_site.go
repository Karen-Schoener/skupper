package api

import (
	"context"

	"fmt"

	"github.com/skupperproject/skupper/api/types"

	"github.com/skupperproject/skupper/pkg/apis/skupper/v2alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
		cm, err := ReadConfigMap(context.Background(), nsName, types.SiteConfigMapName, cli)
		if err != nil {
			return nil, fmt.Errorf("Error reading configmap in namespace %s: %w", nsName, err.Error())
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

func v1IsLinkAccessDefault(siteConfig *types.SiteConfig) bool {
	// TODO: if policy CRD is present (in the site), then check policy.allowIncomingLinks
	return true
}

func createSiteCR(siteConfig *types.SiteConfig) (*v2alpha1.Site, error) {
	resource := &v2alpha1.Site{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "skupper.io/v2alpha1",
			Kind:       "Site",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: siteConfig.Spec.SkupperName,
			// TODO confirm that namespace should not be populated.
			// Namespace: siteConfig.Spec.SkupperNamespace,
		},
		Spec: v2alpha1.SiteSpec{
			// TODO: assume that service account should be unspecified during upgrade.
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
