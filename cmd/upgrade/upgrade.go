package main

import (
	"flag"
	"fmt"
	"github.com/skupperproject/skupper/pkg/version"
	"os"
	"path/filepath"
)

const (
	description = `
TODO fill in description
`
)

func main() {
	var inputPath, outputPath string

	// if -version used, report and exit
	flag.Usage = func() {
		fmt.Println("Skupper upgrade")
		fmt.Printf("%s\n", description)
		fmt.Printf("Usage:\n  %s [options...]\n\n", os.Args[0])
		fmt.Printf("Flags:\n")
		flag.PrintDefaults()
	}
	flag.StringVar(&inputPath, "input", "", "Path to the input directory")
	flag.StringVar(&outputPath, "output", "", "Path to the output directory")
	isVersion := flag.Bool("v", false, "Report the version of the Skupper upgrade command")
	flag.Parse()
	if *isVersion {
		fmt.Println(version.Version)
		os.Exit(0)
	}

	// check if required params are provided
	if inputPath == "" || outputPath == "" {
		fmt.Printf("Both --input and --output flags are required\n")
		os.Exit(1)
	}

	err := validateDirectory(inputPath)
	if err != nil {
		fmt.Printf("Invalid input path: %v\n", err)
		os.Exit(1)
	}

	err = validateDirectory(outputPath)
	if err != nil {
		fmt.Printf("Invalid output path: %v\n", err)
		os.Exit(1)
	}

	err = performUpgrade(inputPath, outputPath)
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

func performUpgrade(inputPath, outputPath string) error {
	// get sites info from debug dump directories
	sitesInfo, err := getSitesInfo(inputPath)
	if err != nil {
		return err
	}

	for _, siteName := range sitesInfo.SiteNames {
		uid := sitesInfo.SiteNameToUid[siteName]
		siteInfo := sitesInfo.UidToSiteInfo[uid]

		upgradeSite(siteInfo.SiteConfig, outputPath)
		upgradeTokens(siteInfo, sitesInfo, outputPath)
	}

	return nil
}
