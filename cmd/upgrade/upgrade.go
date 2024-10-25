package main

import (
	"flag"
	"fmt"
	"github.com/skupperproject/skupper/pkg/version"
	"log"
	"os"
	"path/filepath"

	"github.com/skupperproject/skupper/internal/kube/client"
)

const (
	description = `
TODO fill in description
`
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile) // TMPDBG: for full file path and line numbers

	var outputPath string

	// if -version used, report and exit
	flag.Usage = func() {
		fmt.Println("Skupper upgrade")
		fmt.Printf("%s\n", description)
		fmt.Printf("Usage:\n  %s [options...]\n\n", os.Args[0])
		fmt.Printf("Flags:\n")
		flag.PrintDefaults()
	}
	flag.StringVar(&outputPath, "output", "", "Path to the output directory")
	isVersion := flag.Bool("v", false, "Report the version of the Skupper upgrade command")
	flag.Parse()
	if *isVersion {
		fmt.Println(version.Version)
		os.Exit(0)
	}

	// check if required params are provided
	if outputPath == "" {
		fmt.Printf("--output flag is required\n")
		os.Exit(1)
	}

	err := validateDirectory(outputPath)
	if err != nil {
		fmt.Printf("Invalid output path: %v\n", err)
		os.Exit(1)
	}

	err = performUpgrade(outputPath)
	if err != nil {
		fmt.Printf("Error upgrading sites: %v\n", err)
		os.Exit(1)
	}

}

func validateDirectory(directory string) error {
	path, err := filepath.Abs(directory)
	if err != nil {
		return fmt.Errorf("Failed to resolve file path %s: %s", directory, err)
	}

	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("Failed to stat %s: %s", directory, err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", directory)
	}
	return nil
}

func performUpgrade(outputPath string) error {
	var namespace string
	kubeconfig := getKubeConfig()

	cli, err := client.NewClient(namespace, "", kubeconfig)
	if err != nil {
		return fmt.Errorf("Error getting van client: %w", err.Error())
	}

	// TODO remove uidToSiteConfig var
	uidToSiteConfig, err := getUidToSiteConfig(cli)
	for _, siteConfig := range uidToSiteConfig {
		createDir(outputPath, siteConfig.Spec.SkupperNamespace)
	}

	// get sites info from debug dump directories
	sitesInfo, err := getSitesInfo(cli)
	if err != nil {
		return err
	}

	// iterate over site names in alphabetical order
	for _, siteName := range sitesInfo.SiteNames {
		uid := sitesInfo.SiteNameToUid[siteName]
		// TODO fix me siteInfo := sitesInfo.UidToSiteInfo[uid]

		// TODO remove me upgradeSite(siteInfo.SiteConfig, outputPath)
		//upgradeTokens(siteInfo, sitesInfo, outputPath)
		siteConfig := sitesInfo.UidToSiteConfig[uid]
		err := upgradeSite(siteConfig, outputPath)
		if err != nil {
			return err
		}

		err = upgradeTokens(cli, siteConfig, outputPath, sitesInfo.UidToSiteConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

func createDir(outputPath string, namespace string) error {
	dirPath := filepath.Join(outputPath, namespace)

	// check if the directory already exists
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		// create the directory if it doesn't exist
		err = os.MkdirAll(dirPath, os.ModePerm) // TODO be less permissive than MkdirAll
		if err != nil {
			return fmt.Errorf("Failed to create direcotry %s: %s\n", dirPath, err)
		}
	}
	return nil
}
