package main

import (
	//"context"
        "sigs.k8s.io/yaml"

	"path/filepath"
	"io/ioutil"

	"fmt"
	"log"
	"os"

	//"sigs.k8s.io/yaml"

	//"k8s.io/apimachinery/pkg/api/errors"

	"github.com/skupperproject/skupper/api/types"
	//"github.com/skupperproject/skupper/internal/kube/client"
	//"log"
	//"github.com/davecgh/go-spew/spew"
	//"github.com/skupperproject/skupper/pkg/site"
	//corev1 "k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"github.com/skupperproject/skupper/pkg/apis/skupper/v1alpha1"
)

type SiteInfo struct {
	uidToSiteConfig map[string]*types.SiteConfig
	siteNameToUid   map[string]string
}

func getSiteInfo(inputPath string) (*SiteInfo, error) {
	siteInfo := &SiteInfo{
		uidToSiteConfig: map[string]*types.SiteConfig{},
		siteNameToUid:   map[string]string{},
	}

	log.Printf("TMPDBG: getSiteInfo: entering")

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
		log.Printf("TMPDBG: getSiteInfo: entry.Name(): %+v\n", entry.Name())
		siteConfig, err := readSiteConfig(debugDumpPath)
		if err != nil {
			return nil, fmt.Errorf("Failed to read site config from debug dump directory: %s: %v", entry.Name(), err)
		}
		log.Printf("TMPDBG: getSiteInfo: siteConfig: %+v\n", siteConfig)
	}

	return siteInfo, nil
}

type SkupperSite struct {
	name string `yaml:"data.name"`
	ingress string `yaml:"data.ingress"`
	routerMode string `yaml:"data.router-mode"`
}

type SkupperSiteNested struct {
	Data struct {
		Name string `yaml:"name"`
		Ingress string `yaml:"ingress"`
		RouterMode string `yaml:"router-mode"`
	} `yaml:"data"`
}

func readSiteConfig(path string) (*types.SiteConfig, error) {
	log.Printf("TMPDBG: readSiteConfig: 20241015_445: path=%+v", path)

	filename := "skupper-site.yaml"
	configMapsPath := filepath.Join(path, "configmaps")
	skupperSiteFile := filepath.Join(configMapsPath, filename)
	data, err := ioutil.ReadFile(skupperSiteFile)
	if err != nil {
		return nil, fmt.Errorf("Error reading file %s from debug dump: %v", filename, err)
	}

	readSiteConfigNested := false
	readSiteConfigNested = true
	if readSiteConfigNested {
		var siteNested SkupperSiteNested
		log.Printf("TMPDBG: readSiteConfig: data=%+v", string(data))
		err = yaml.Unmarshal(data, &siteNested)
		log.Printf("TMPDBG: readSiteConfig: case 10: after yaml.Unmarshal: err=%+v", err)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling file %s data from debug dump: %v", filename, err)
		}
		log.Printf("TMPDBG: readSiteConfig: case 11: after yaml.Unmarshal: siteNested=%+v", siteNested)
	}

	readSiteConfig := false
	readSiteConfig = true
	if readSiteConfig {
		var siteNotNested SkupperSite
		log.Printf("TMPDBG: readSiteConfig: data=%+v", string(data))
		err = yaml.Unmarshal(data, &siteNotNested)
		log.Printf("TMPDBG: readSiteConfig: case 30: after yaml.Unmarshal: err=%+v", err)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling file %s data from debug dump: %v", filename, err)
		}
		log.Printf("TMPDBG: readSiteConfig: case 31: after yaml.Unmarshal: siteNotNested=%+v", siteNotNested)
	}

	testing := true
	if testing {
		return nil, nil
	}

	siteConfig := &types.SiteConfig{}

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
