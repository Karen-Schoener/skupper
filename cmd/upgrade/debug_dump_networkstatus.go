package main

import (
	"io/ioutil"
	"path/filepath"

	"fmt"
	"log"
	//"os"
	"gopkg.in/yaml.v3"
)

type NetworkStatusDataUnmarshaled struct {
	Data NetworkStatus `yaml:"data"`
}
type NetworkStatus struct {
	// TODO: add addresses
	//SiteStatus []SiteStatus `yaml:"siteStatus"`
	Site []interface{} `yaml:"siteStatus"`
	//RouterStatus []RouterStatus `yaml:"routerStatus"`
	//RouterStatus interface{} `yaml:"routerStatus"`
}

type SiteStatus struct {
	Site interface{} `yaml:"siteStatus"`
	//RouterStatus []RouterStatus `yaml:"routerStatus"`
	RouterStatus interface{} `yaml:"routerStatus"`
}

type NetworkStatusDataUnmarshaled_1 struct {
	Data map[string]interface{} `yaml:"data"`
}

type NetworkStatusDataUnmarshaled_2 struct {
	Data map[string]NetworkStatus_2 `yaml:"data"`
}
type NetworkStatus_2 struct {
	Addresses []interface{} `yaml:"addresses"`
	SiteStatus []interface{} `yaml:"siteStatus"`
	//SiteStatus []SiteStatus_2 `yaml:"siteStatus"`
}
type SiteStatus_2 struct {
	SiteStatus []map[string]interface{} `yaml:"site"`
}

func readNetworkStatus(path string) (*NetworkStatus, error) {
	networkStatus := &NetworkStatus{}

	log.Printf("TMPDBG: readNetworkStatus: path=%+v", path)

	filename := "skupper-network-status.yaml"
	configMapsPath := filepath.Join(path, "configmaps")
	configmapFile := filepath.Join(configMapsPath, filename)
	data, err := ioutil.ReadFile(configmapFile)
	if err != nil {
		return nil, fmt.Errorf("Error reading file %s: %v", filename, err)
	}

	log.Printf("TMPDBG: readNetworkStatus: data=%+v", string(data))

	var tmpdata_0 map[string]interface{}
	err = yaml.Unmarshal(data, &tmpdata_0)
	log.Printf("TMPDBG: readNetworkStatus: TMPDBG: case 00.1: err=%+v", err)
	log.Printf("TMPDBG: readNetworkStatus: TMPDBG: case 00.2: tmpdata_0=%+v", tmpdata_0)

	var tmpdata_1 NetworkStatusDataUnmarshaled_1
	err = yaml.Unmarshal(data, &tmpdata_1)
	log.Printf("TMPDBG: readNetworkStatus: TMPDBG: case 10.1: err=%+v", err)
	log.Printf("TMPDBG: readNetworkStatus: TMPDBG: case 10.2: tmpdata_1=%+v", tmpdata_1)

	var tmpdata_2 NetworkStatusDataUnmarshaled_2
	err = yaml.Unmarshal(data, &tmpdata_2)
	log.Printf("TMPDBG: readNetworkStatus: TMPDBG: case 20.1: err=%+v", err)
	log.Printf("TMPDBG: readNetworkStatus: TMPDBG: case 20.2: tmpdata_2=%+v", tmpdata_2)

	var networkStatusDataUnmarshaled NetworkStatusDataUnmarshaled
	err = yaml.Unmarshal(data, &networkStatusDataUnmarshaled)
	if err != nil {
		log.Printf("TMPDBG: readNetworkStatus: err=%+v", err)
		return nil, fmt.Errorf("Error unmarshalling file %s data: %v", filename, err)
	}
	log.Printf("TMPDBG: readNetworkStatus: networkStatusDataUnmarshaled=%+v", networkStatusDataUnmarshaled)

	return networkStatus, nil
}
