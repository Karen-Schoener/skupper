package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skupperproject/skupper/api/types"

	"github.com/skupperproject/skupper/pkg/apis/skupper/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/skupperproject/skupper/internal/cmd/skupper/common/utils"
)

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

func createSiteCR(siteConfig *types.SiteConfig) (*v1alpha1.Site, error) {
	// TODO: assume that service account should be unspecified during upgrade.
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
	if siteConfig.Spec.RouterMode == string(types.TransportModeEdge) {
		resource.Spec.RouterMode = string(types.TransportModeEdge)
	}
	// TODO: confirm logic to set LinkAccess
	if v1IsLinkAccessDefault(siteConfig) {
		resource.Spec.LinkAccess = "default"
	}

	return resource, nil
}

func saveSiteCR(resource *v1alpha1.Site, outputPath string, siteName string) error {
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
