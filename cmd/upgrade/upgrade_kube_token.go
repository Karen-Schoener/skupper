package main

import (
	"fmt"
	"os"
	"path/filepath"

	"strings"

	"github.com/skupperproject/skupper/pkg/apis/skupper/v1alpha1"

	"github.com/skupperproject/skupper/api/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/skupperproject/skupper/internal/cmd/skupper/common/utils"
)

func upgradeTokens(siteInfo *SiteInfo, sitesInfo *SitesInfo, outputPath string) error {
	// Scan the network status for the site the matches input param siteInfo.
	var siteStatus SiteStatus
	found := false
	for _, siteStatus = range siteInfo.NetworkStatus.SiteStatus {
		if siteInfo.SiteConfig.Reference.UID == siteStatus.Site.Identity {
			// look for matching site in the network status info
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("Failed to upgrade tokens for site %s, site not found in network status", siteInfo.SiteConfig.Spec.SkupperName)
	}

	// For every link that originates from this site, create CRs: AcessGrant, AccessToken.
	for _, routerStatus := range siteStatus.RouterStatus {
		for _, link := range routerStatus.Links {
			targetUid := LookupTargetSiteUidForLink(siteInfo.NetworkStatus, &link)
			if len(targetUid) <= 0 {
				return fmt.Errorf("Failed to lookup target site for link %v", link)
			}
			if targetSiteInfo, ok := sitesInfo.UidToSiteInfo[targetUid]; ok {
				accessGrant, err := createAccessGrantCR(siteInfo.SiteConfig, targetSiteInfo.SiteConfig)
				if err != nil {
					return err
				}
				accessToken, err := createAccessTokenCR(siteInfo.SiteConfig, targetSiteInfo.SiteConfig)
				if err != nil {
					return err
				}
				saveAccessGrantCR(accessGrant, outputPath, targetSiteInfo.SiteConfig.Spec.SkupperName)
				saveAccessTokenCR(accessToken, outputPath, siteInfo.SiteConfig.Spec.SkupperName)
			} else {
				return fmt.Errorf("Failed to lookup target site info for link %v", link)
			}
		}
	}
	return nil
}

func createAccessGrantCR(sourceSiteConfig, targetSiteConfig *types.SiteConfig) (*v1alpha1.AccessGrant, error) {
	DefaultRedemptionsAllowed := 10
	DefaultExpirationWindow := "1h"
	resource := &v1alpha1.AccessGrant{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "skupper.io/v1alpha1",
			Kind:       "AccessGrant",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: generateTokenName(sourceSiteConfig.Spec.SkupperName, targetSiteConfig.Spec.SkupperName),
		},
		Spec: v1alpha1.AccessGrantSpec{
			RedemptionsAllowed: DefaultRedemptionsAllowed,
			ExpirationWindow:   DefaultExpirationWindow,
			// TODO how to populate in AccessGrant: RedemptionsAllowed, ExpirationWindow
		},
	}

	return resource, nil
}

func saveAccessTokenCR(resource *v1alpha1.AccessToken, outputPath string, siteName string) error {
	targetDir := filepath.Join(outputPath, siteName)
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("Failed to create directory: %v", err)
	}

	filepath := filepath.Join(targetDir, fmt.Sprintf("%s_access_token.yaml", resource.ObjectMeta.Name))

	data, err := utils.Encode("yaml", resource)
	if err != nil {
		return fmt.Errorf("Failed to marshal site resource to YAML: %w", err.Error())
	}

	err = os.WriteFile(filepath, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("Failed to write site resource to file: %w", err.Error())
	}
	fmt.Printf("Wrote access token CR to file: %s\n", filepath)
	return nil
}

func saveAccessGrantCR(resource *v1alpha1.AccessGrant, outputPath string, siteName string) error {
	targetDir := filepath.Join(outputPath, siteName)
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("Failed to create directory: %v", err)
	}

	filepath := filepath.Join(targetDir, fmt.Sprintf("%s_access_grant.yaml", resource.ObjectMeta.Name))

	data, err := utils.Encode("yaml", resource)
	if err != nil {
		return fmt.Errorf("Failed to marshal site resource to YAML: %w", err.Error())
	}

	err = os.WriteFile(filepath, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("Failed to write site resource to file: %w", err.Error())
	}
	fmt.Printf("Wrote access grant CR to file: %s\n", filepath)
	return nil
}

func createAccessTokenCR(sourceSiteConfig, targetSiteConfig *types.SiteConfig) (*v1alpha1.AccessToken, error) {
	resource := &v1alpha1.AccessToken{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "skupper.io/v1alpha1",
			Kind:       "AccessToken",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: generateTokenName(sourceSiteConfig.Spec.SkupperName, targetSiteConfig.Spec.SkupperName),
		},
		Spec: v1alpha1.AccessTokenSpec{
			//Url:  accessGrant.Status.Url,
			//Code: accessGrant.Status.Code,
			//Ca:   accessGrant.Status.Ca,
		},
		// TODO how to populate in AccessToken (after v2 sites are running: Url, Code, Ca.  Ok to use token name to to help correlate but really are missing target site name.  Annotate with comment in yaml file?
	}

	return resource, nil
}

func LookupTargetSiteUidForLink(networkStatus *NetworkStatus, link *Link) string {
	for _, siteStatus := range networkStatus.SiteStatus {
		for _, routerStatus := range siteStatus.RouterStatus {
			if strings.HasSuffix(routerStatus.Router.Name, link.Name) {
				return siteStatus.Site.Identity
			}
		}
	}
	return ""
}

func generateTokenName(sourceNamsepace, targetNamespace string) string {
	return fmt.Sprintf("%s_to_%s", sourceNamsepace, targetNamespace)
}
