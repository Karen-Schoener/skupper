package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	//TMPDBG "os/exec"
	//TMPDBG "path"
	//TMPDBG "path/filepath"

	//TMPDBG "github.com/skupperproject/skupper/api/types"
	//TMPDBG internalbundle "github.com/skupperproject/skupper/internal/nonkube/bundle"
	//TMPDBG "github.com/skupperproject/skupper/pkg/config"
	//TMPDBG "github.com/skupperproject/skupper/pkg/nonkube/api"
	//TMPDBG "github.com/skupperproject/skupper/pkg/nonkube/bundle"
	//TMPDBG "github.com/skupperproject/skupper/pkg/nonkube/common"
	//TMPDBG "github.com/skupperproject/skupper/pkg/nonkube/compat"
	//TMPDBG "github.com/skupperproject/skupper/pkg/nonkube/systemd"
	//TMPDBG "github.com/skupperproject/skupper/pkg/utils"
	"github.com/skupperproject/skupper/pkg/version"
)

const (
	description = `
TODO fill me in

Bootstraps a nonkube Skupper site base on the provided flags.
`
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile) // TMPDBG: for full file path and line numbers

	var inputPath, outputPath string

	// if -version used, report and exit
	flag.Usage = func() {
		fmt.Println("Skupper bootstrap")
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
	//log.Printf("TMPDBG: inputPath=%+v", inputPath)
	//log.Printf("TMPDBG: oututPath=%+v", outputPath)

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

	err = upgradeSites(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error upgrading sites: %v\n", err)
		os.Exit(1)
	}

}

func validateDirectory(directory string) error {
	//log.Printf("TMPDBG: validateDirectory: directory=%+v", directory)
	path, err := filepath.Abs(directory)
	//log.Printf("TMPDBG: validateDirectory: path=%+v", path)
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

func upgradeSites(inputPath, outputPath string) error {
	log.Printf("TMPDBG: upgradeSites: entering")
	siteInfo, err := getSiteInfo(inputPath)
	if err != nil {
		log.Printf("TMPDBG: upgradeSites: after getSiteInfo: err=%+v", err)
		return err // TODO append error string before returning?
	}
	log.Printf("TMPDBG: upgradeSites: siteInfo=%+v", siteInfo)
	log.Printf("TMPDBG: upgradeSites: returning")
	return nil
}
