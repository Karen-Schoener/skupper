package main

import (
	"context"
	"fmt"
	"log"

	"github.com/skupperproject/skupper/api/types"
	"github.com/skupperproject/skupper/internal/kube/client"
	"github.com/skupperproject/skupper/pkg/apis/skupper/v2alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"strconv"
)

func readConnectionTokens(ctx context.Context, namespace string, cli *client.KubeClient) ([]corev1.Secret, error) {
	kubeClient := cli.GetKubeClient()

	secrets, err := kubeClient.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{LabelSelector: "skupper.io/type=connection-token"})
	if err != nil {
		return nil, err
	}
	return secrets.Items, nil
}

// Read connection tokens for the v1 siteConfig.
// Save resulting v2 AccessGrants, AccessTokens in nameToV2SiteState.
func createTokenCRs(cli *client.KubeClient, siteConfig *types.SiteConfig, uidToSiteConfig map[string]*types.SiteConfig, nameToV2SiteState map[string]*SiteState) error {

	secrets, err := readConnectionTokens(context.Background(), siteConfig.Spec.SkupperNamespace, cli)
	if err != nil {
		return fmt.Errorf("Error getting connection tokens for siteConfig.Spec.SkupperNamespace %v: %w", siteConfig.Spec.SkupperNamespace, err.Error())
	}

	for _, secret := range secrets {
		generatedBy, ok := secret.ObjectMeta.Annotations[types.TokenGeneratedBy]
		if !ok {
			return fmt.Errorf("Error getting link target site for link %s for site %s, annotations=%+v", secret.ObjectMeta.Name, siteConfig.Spec.SkupperName, secret.ObjectMeta.Annotations)
		}

		tokenCost, ok := secret.ObjectMeta.Annotations[types.TokenCost]
		cost := 1
		if ok {
			cost, err = strconv.Atoi(tokenCost)
			if err != nil {
				return fmt.Errorf("Error getting token cost for link %s for site %s, %w", secret.ObjectMeta.Name, siteConfig.Spec.SkupperName, err)
			}
		}
		debug := false
		if debug {
			log.Printf("TODO: use cost %v for token cost for link %s for site %s", cost, secret.ObjectMeta.Name, siteConfig.Spec.SkupperName)
		}

		// TODO should upgrade verify the original v1 site version?  should all sites be running a specific version?
		//if siteVersion, ok := secret.ObjectMeta.Annotations[types.SiteVersion]; ok {
		//    log.Printf("TMPDBG: upgradeTokens: siteConfig.Spec.SkupperNamespace: %v, siteVersion: %+v\n", siteConfig.Spec.SkupperNamespace, siteVersion)
		//}

		if targetSiteConfig, ok := uidToSiteConfig[generatedBy]; ok {
			name := generateTokenName(siteConfig.Spec.SkupperName, targetSiteConfig.Spec.SkupperName)
			accessGrant, err := createAccessGrantCR(siteConfig, targetSiteConfig, name)
			if err != nil {
				return err
			}
			accessToken, err := createAccessTokenCR(siteConfig, targetSiteConfig, name)
			if err != nil {
				return err
			}
			nameToV2SiteState[targetSiteConfig.Spec.SkupperName].Grants[accessGrant.Name] = accessGrant
			nameToV2SiteState[siteConfig.Spec.SkupperName].Tokens[accessGrant.Name] = accessToken
		}
	}
	return nil
}

func createAccessGrantCR(sourceSiteConfig, targetSiteConfig *types.SiteConfig, name string) (*v2alpha1.AccessGrant, error) {
	DefaultRedemptionsAllowed := 10
	DefaultExpirationWindow := "1h"
	resource := &v2alpha1.AccessGrant{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "skupper.io/v2alpha1",
			Kind:       "AccessGrant",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v2alpha1.AccessGrantSpec{
			RedemptionsAllowed: DefaultRedemptionsAllowed,
			ExpirationWindow:   DefaultExpirationWindow,
			// TODO how to populate in AccessGrant: RedemptionsAllowed, ExpirationWindow
		},
	}

	return resource, nil
}

func createAccessTokenCR(sourceSiteConfig, targetSiteConfig *types.SiteConfig, name string) (*v2alpha1.AccessToken, error) {
	resource := &v2alpha1.AccessToken{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "skupper.io/v2alpha1",
			Kind:       "AccessToken",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v2alpha1.AccessTokenSpec{
			//Url:  accessGrant.Status.Url,
			//Code: accessGrant.Status.Code,
			//Ca:   accessGrant.Status.Ca,
		},
		// TODO how to populate in AccessToken (after v2 sites are running: Url, Code, Ca.  Ok to use token name to to help correlate but really are missing target site name.  Annotate with comment in yaml file?
	}

	return resource, nil
}

func generateTokenName(sourceNamsepace, targetNamespace string) string {
	return fmt.Sprintf("%s-to-%s", sourceNamsepace, targetNamespace)
}
