package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skupperproject/skupper/api/types"
	"github.com/skupperproject/skupper/internal/cmd/skupper/common/utils"
	"github.com/skupperproject/skupper/internal/kube/client"
	"github.com/skupperproject/skupper/pkg/apis/skupper/v2alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"strconv"

	"bufio"
	"bytes"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func readConnectionTokens(ctx context.Context, namespace string, cli *client.KubeClient) ([]corev1.Secret, error) {
	kubeClient := cli.GetKubeClient()

	secrets, err := kubeClient.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{LabelSelector: "skupper.io/type=connection-token"})
	if err != nil {
		return nil, err
	}
	return secrets.Items, nil
}

func upgradeTokens(cli *client.KubeClient, siteConfig *types.SiteConfig, outputPath string, uidToSiteConfig map[string]*types.SiteConfig) error {
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
			saveAccessGrantCR(accessGrant, outputPath, targetSiteConfig.Spec.SkupperName)
			saveAccessTokenCR(accessToken, cost, outputPath, siteConfig.Spec.SkupperName)

			//linkToken, err := createLinkTokenCR(name, cost)
			//saveLinkTokenCR(linkToken, name, outputPath, siteConfig.Spec.SkupperName)
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

// TODO write cost to file in yaml comment?
func saveAccessTokenCR(resource *v2alpha1.AccessToken, cost int, outputPath string, siteName string) error {
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

func saveAccessGrantCR(resource *v2alpha1.AccessGrant, outputPath string, siteName string) error {
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

// borrowed from: ./pkg/nonkube/api/token.go
type Token struct {
	Links []*v2alpha1.Link
}

// TODO: when to create multiple links?  how to detect HA?
func createLinkTokenCR(linkName string, cost int) (*Token, error) {
	// adjusting name to match the standard used by pkg/site/link.go
	//clientSecret.Name = fmt.Sprintf("link-%s", linkName)

	resource := &Token{
		Links: []*v2alpha1.Link{
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "skupper.io/v2alpha1",
					Kind:       "Link",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: linkName,
				},
				Spec: v2alpha1.LinkSpec{
					//TlsCredentials: clientSecret.Name,
					Cost: cost,
				},
			},
		},
	}
	return resource, nil
}

// lifted Marshal from file: ./pkg/nonkube/api/token.go
func (t *Token) Marshal() ([]byte, error) {
	s := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, json.SerializerOptions{Yaml: true})
	buffer := new(bytes.Buffer)
	writer := bufio.NewWriter(buffer)
	var err error
	//_, _ = writer.Write([]byte("---\n"))
	//err := s.Encode(t.Secret, writer)
	//if err != nil {
	//	return nil, err
	//}
	for _, l := range t.Links {
		_, _ = writer.Write([]byte("---\n"))
		err = s.Encode(l, writer)
		if err != nil {
			return nil, err
		}
		if err = writer.Flush(); err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

func saveLinkTokenCR(resource *Token, name string, outputPath string, siteName string) error {
	targetDir := filepath.Join(outputPath, siteName)
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("Failed to create directory: %v", err)
	}

	filepath := filepath.Join(targetDir, fmt.Sprintf("%s_link_token.yaml", name))

	// open file in write mode
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open file %s: %w", filepath, err)
	}
	defer file.Close()

	data, err := resource.Marshal()
	_, err = file.Write(data)

	fmt.Printf("Wrote link token CR to file: %s\n", filepath)
	return nil
}
