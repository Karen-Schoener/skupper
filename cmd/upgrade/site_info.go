package main

import (
	"context"

	"fmt"

	"github.com/skupperproject/skupper/api/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/skupperproject/skupper/internal/kube/client"
	"github.com/skupperproject/skupper/pkg/site"
	"sort"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

type SitesInfo struct {
	//UidToSiteInfo map[string]*SiteInfo
	UidToSiteConfig map[string]*types.SiteConfig
	SiteNameToUid   map[string]string
	SiteNames       []string
}

func getKubeConfig() string {
	var kubeconfig string // TODO
	return kubeconfig
}

func getSitesInfo(cli *client.KubeClient) (*SitesInfo, error) {
	sitesInfo := &SitesInfo{
		//UidToSiteInfo: map[string]*SiteInfo{},
		UidToSiteConfig: map[string]*types.SiteConfig{},
		SiteNameToUid:   map[string]string{},
		SiteNames:       []string{},
	}

	//sitesInfo.UidToSiteConfig, err := getUidToSiteConfig(cli)
	uidToSiteConfig, err := getUidToSiteConfig(cli)
	if err != nil {
		return nil, err
	}

	sitesInfo.UidToSiteConfig = uidToSiteConfig

	for _, siteConfig := range sitesInfo.UidToSiteConfig {
		sitesInfo.SiteNameToUid[siteConfig.Spec.SkupperName] = siteConfig.Reference.UID
		sitesInfo.SiteNames = append(sitesInfo.SiteNames, siteConfig.Spec.SkupperName)
	}

	sort.Strings(sitesInfo.SiteNames)

	return sitesInfo, nil
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
