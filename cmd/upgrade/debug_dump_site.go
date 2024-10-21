package main

import (
	"fmt"
	"github.com/skupperproject/skupper/api/types"
	"github.com/skupperproject/skupper/pkg/site"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubetypes "k8s.io/apimachinery/pkg/types"
	"os"
	"path/filepath"
	"sort"
)

type SkupperSiteUnmarshaled struct {
	Data     map[string]string `yaml:"data"`
	Metadata struct {
		Name              string                 `yaml:"name"`
		Namespace         string                 `yaml:"namespace"`
		Labels            map[string]interface{} `yaml:"labels"`
		UID               string                 `yaml:"uid"`
		CreationTimestamp string                 `yaml:"creationTimestamp"`
	} `yaml:"metadata"`
}

type SiteInfo struct {
	SiteConfig    *types.SiteConfig
	NetworkStatus *NetworkStatus
}

type SitesInfo struct {
	UidToSiteInfo map[string]*SiteInfo
	SiteNameToUid map[string]string
	SiteNames     []string
}

func getSitesInfo(inputPath string) (*SitesInfo, error) {
	sitesInfo := &SitesInfo{
		UidToSiteInfo: map[string]*SiteInfo{},
		SiteNameToUid: map[string]string{},
		SiteNames:     []string{},
	}

	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return nil, fmt.Errorf("Error resolving input path: %v", err)
	}

	dirEntries, err := os.ReadDir(absInputPath)
	if err != nil {
		return nil, fmt.Errorf("Error reading directory contents: %v", err)
	}

	for _, entry := range dirEntries {
		debugDumpPath := filepath.Join(absInputPath, entry.Name())
		if !entry.IsDir() {
			continue
		}
		err := validateDebugDump(debugDumpPath)
		if err != nil {
			return nil, fmt.Errorf("Invalid debug dump directory: %s: %v", entry.Name(), err)
		}

		siteConfig, err := readSiteConfig(debugDumpPath)
		if err != nil {
			return nil, fmt.Errorf("Failed to read site config from directory: %s: %v", entry.Name(), err)
		}

		networkStatus, err := readNetworkStatus(debugDumpPath)
		if err != nil {
			return nil, fmt.Errorf("Failed to read network status from directory: %s: %v", entry.Name(), err)
		}

		sitesInfo.UidToSiteInfo[siteConfig.Reference.UID] = &SiteInfo{
			SiteConfig:    siteConfig,
			NetworkStatus: networkStatus,
		}
		sitesInfo.SiteNameToUid[siteConfig.Spec.SkupperName] = siteConfig.Reference.UID
		sitesInfo.SiteNames = append(sitesInfo.SiteNames, siteConfig.Spec.SkupperName)
	}

	sort.Strings(sitesInfo.SiteNames)

	return sitesInfo, nil
}

func readSiteConfig(path string) (*types.SiteConfig, error) {
	filename := "skupper-site.yaml"
	configMapsPath := filepath.Join(path, "configmaps")
	configmapFile := filepath.Join(configMapsPath, filename)
	data, err := ioutil.ReadFile(configmapFile)
	if err != nil {
		return nil, fmt.Errorf("Error reading file %s: %v", filename, err)
	}

	var skupperSiteUnmarshaled SkupperSiteUnmarshaled
	err = yaml.Unmarshal(data, &skupperSiteUnmarshaled)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling file %s data: %v", filename, err)
	}

	tmpSiteConfig := metav1.ObjectMeta{
		Name:      skupperSiteUnmarshaled.Metadata.Name,
		Namespace: skupperSiteUnmarshaled.Metadata.Namespace,
		UID:       kubetypes.UID(skupperSiteUnmarshaled.Metadata.UID),
	}
	tmpTypeMeta := metav1.TypeMeta{}
	defaultIngress := "loadbalancer"

	// TODO: for now, the prototype calls site.ReadSiteConfigFrom - with hacked input params: siteConfig, typeMeta.
	//       While this is ok for short term testing, should consider: create a copy of ReadSiteConfigFrom to this file.
	//
	//       One downside to this hack is: annotions, labels of the skupper-site configmap are currently ignored.
	siteConfig, err := site.ReadSiteConfigFrom(&tmpSiteConfig, &tmpTypeMeta, skupperSiteUnmarshaled.Data, defaultIngress)

	return siteConfig, nil
}

// TODO is validateDebugDump necessary?
func validateDebugDump(path string) error {
	expectedDirectories := []string{"configmaps", "deployments", "pods", "services", "skupper-info"}
	// Check expected directories exist.
	for _, dirName := range expectedDirectories {
		dirPath := filepath.Join(path, dirName)
		stat, err := os.Stat(dirPath)
		if err != nil {
			return fmt.Errorf("Directory %s not found: %v", dirName, err)
		}
		if !stat.IsDir() {
			return fmt.Errorf("Directory %s not not a directory: %v", dirName, err)
		}
	}
	return nil
}
