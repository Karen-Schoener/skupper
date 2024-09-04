package main

import (
	"fmt"
	"log"
	"os"

	"path/filepath"

	"github.com/skupperproject/skupper/internal/kube/client"

	"github.com/spf13/cobra"
)

func validateOutputPath(directory string) error {
	stat, err := os.Stat(directory)
	if err != nil {
		return fmt.Errorf("Failed to stat %s: %s\n", directory, err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s is not a directory\n", directory)
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

func upgradeSites(outputPath string) error {
	var namespace string
	kubeconfig := getKubeConfig()

	cli, err := client.NewClient(namespace, "", kubeconfig)
	if err != nil {
		return fmt.Errorf("Error getting van client: %w", err.Error())
	}

	uidToSiteConfig, err := getUidToSiteConfig(cli)
	for _, siteConfig := range uidToSiteConfig {
		createDir(outputPath, siteConfig.Spec.SkupperNamespace)
	}

	err = upgradeConfigMapSkupperSites(cli, outputPath, uidToSiteConfig)
	if err != nil {
		log.Fatalf("Error upgrading site: %v: %v\n", namespace, err)
	}

	err = upgradeTokens(cli, outputPath, uidToSiteConfig)
	if err != nil {
		log.Fatalf("Error upgrading secrets: %v\n", err)
	}
	return nil
}

func main() {
	var outputPath string

	log.SetFlags(log.LstdFlags | log.Llongfile) // TMPDBG: for full file path and line numbers

	var rootCmd = &cobra.Command{Use: "upgrade-v1-v2"}

	var upgradeCmd = &cobra.Command{
		Use:   "sites",
		Short: "Upgrade v1 skupper site resources to v2 CRs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputPath == "" {
				fmt.Printf("otuput parameter is required\n")
				return nil
			}
			err := validateOutputPath(outputPath)
			if err != nil {
				fmt.Printf("output directory does not exist: %v\n", outputPath)
				return nil
			}

			return upgradeSites(outputPath)
		},
	}
	upgradeCmd.Flags().StringVarP(&outputPath, "output", "o", "./output", "Output directory")
	rootCmd.AddCommand(upgradeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
