package main

import (
	"context"

	"log"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/skupperproject/skupper/api/types"
	"github.com/skupperproject/skupper/internal/kube/client"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func readConnectionTokens(ctx context.Context, namespace string, labelSelector string, cli *client.KubeClient) ([]corev1.Secret, error) {
	kubeClient := cli.GetKubeClient()

	secrets, err := kubeClient.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, err
	}
	return secrets.Items, nil
}

func upgradeTokens(cli *client.KubeClient, outputPath string) error {
	uidToSiteConfig, err := getUidToSiteConfig(cli)
	if err != nil {
		return fmt.Errorf("Error getting uid to site configs: %w", err.Error())
	}

	for _, siteConfig := range uidToSiteConfig {
		targetNamespaceToLinkCmd := map[string]string{}

		secrets, err := readConnectionTokens(context.Background(), siteConfig.Spec.SkupperNamespace, "skupper.io/type=connection-token", cli)
		if err != nil {
			return fmt.Errorf("Error getting connection tokens for siteConfig.Spec.SkupperNamespace %v: %w", siteConfig.Spec.SkupperNamespace, err.Error())
		}

		for _, secret := range secrets {

			//log.Printfif generatedBy, ok := secret.ObjectMeta.Annotations[types.TokenGeneratedBy] {
			//generatedBy, ok := secret.ObjectMeta.Annotations[types.TokenGeneratedBy]
			//log.Printf("TMPDBG: upgradeTokens: siteConfig.Spec.SkupperNamespace: %v, generatedBy: %+v\n", siteConfig.Spec.SkupperNamespace, generatedBy)
			//log.Printf("TMPDBG: upgradeTokens: siteConfig.Spec.SkupperNamespace: %v, ok: %+v\n", siteConfig.Spec.SkupperNamespace, ok)
			generatedBy, ok := secret.ObjectMeta.Annotations[types.TokenGeneratedBy]
			if !ok {
				log.Printf("TMPDBG: upgradeTokens: missing generatedBy annotation, siteConfig.Spec.SkupperNamespace: %v, generatedBy: %+v\n", siteConfig.Spec.SkupperNamespace, generatedBy)
			}
			//if ok {
				//log.Printf("TMPDBG: upgradeTokens: siteConfig.Spec.SkupperNamespace: %v, generatedBy: %+v\n", siteConfig.Spec.SkupperNamespace, generatedBy)
				// TODO return error if not found
			//}
			cost, ok := secret.ObjectMeta.Annotations[types.TokenCost]
			if ok {
				//log.Printf("TMPDBG: upgradeTokens: siteConfig.Spec.SkupperNamespace: %v, cost: %+v\n", siteConfig.Spec.SkupperNamespace, cost)
			}
			//if siteVersion, ok := secret.ObjectMeta.Annotations[types.SiteVersion]; ok {
				// TODO migtration should check original v1 site version.  should all sites be specific version?
				//log.Printf("TMPDBG: upgradeTokens: siteConfig.Spec.SkupperNamespace: %v, siteVersion: %+v\n", siteConfig.Spec.SkupperNamespace, siteVersion)
			//}

			//for _, ownerRef := range secret.ObjectMeta.OwnerReferences {
			//	log.Printf("TMPDBG: upgradeTokens: siteConfig.Spec.SkupperNamespace: %v, spew.Sdump(ownerRef): %+v\n", siteConfig.Spec.SkupperNamespace, spew.Sdump(ownerRef))
			//	log.Printf("TMPDBG: upgradeTokens: siteConfig.Spec.SkupperNamespace: %v, ownerRef.UID: %+v\n", siteConfig.Spec.SkupperNamespace, ownerRef.UID)
			//}

			if targetSiteConfig, ok := uidToSiteConfig[generatedBy]; ok {
				cmd := generateLinkCommand(targetSiteConfig.Spec.SkupperNamespace, cost, generateTlsSecretName(siteConfig.Spec.SkupperNamespace, targetSiteConfig.Spec.SkupperNamespace))
				targetNamespaceToLinkCmd[targetSiteConfig.Spec.SkupperNamespace] = cmd
//generateLinkCommand(targetSiteConfig.Spec.SkupperNamespace, cost, generateTlsSecretName(siteConfig.Spec.SkupperNamespace, targetSiteConfig.Spec.SkupperNamespace))
			}

		}
		for targetNamespace, linkCmd := range targetNamespaceToLinkCmd {
			dirPath := filepath.Join(outputPath, siteConfig.Spec.SkupperNamespace)
			err := writeLinkCmdToFile(dirPath, targetNamespace, linkCmd)
			if err != nil {
				return fmt.Errorf("Error writing link command to file for target namespace %s: %w", targetNamespace, err.Error())
			}
		}
	}

	return nil
}

func writeLinkCmdToFile(outputPath string, targetNamespace string, cmd string) error {
	filename := path.Join(outputPath, "link_to_"+targetNamespace+".sh")
	scriptContent := fmt.Sprintf("#!/bin/bash\n%s\n", cmd)
	err := os.WriteFile(filename, []byte(scriptContent), 0755)
	if err != nil {
		return fmt.Errorf("Error creating file: %s", err)
	}
	return nil
}

func generateTlsSecretName(sourceNamsepace, targetNamespace string) string {
	return fmt.Sprintf("%s-to-%s", sourceNamsepace, targetNamespace)
}

// Generate v2 link command
func generateLinkCommand(targetNamespace string, cost string, tlsSecret string) string {
	// TODO should the v1 link name be preserved?
	cmd := "skupper link generate -n " + targetNamespace
	if len(cost) > 0 {
		cmd += " --cost=" + cost
	}
	if len(tlsSecret) > 0 {
		cmd += " --tlsSecret=" + tlsSecret
	}
	return cmd
}
