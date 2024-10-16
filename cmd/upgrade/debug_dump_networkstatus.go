package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type NetworkStatusDataUnmarshaled struct {
	Data struct {
		NetworkStatus string `yaml:"NetworkStatus"` // This is the stringified JSON
	} `yaml:"data"`
}

type NetworkStatus struct {
	Addresses  interface{}  `json:"addresses"`  // Change yaml to json since this is a JSON string
	SiteStatus []SiteStatus `json:"siteStatus"` // List of sites
}

type SiteStatus struct {
	Site         Site           `json:"site"`
	RouterStatus []RouterStatus `json:"routerStatus"`
}

type Site struct {
	RecType     string `json:"recType"`
	Identity    string `json:"identity"`
	StartTime   int64  `json:"startTime"`
	EndTime     int64  `json:"endTime"`
	Source      string `json:"source"`
	Platform    string `json:"platform"`
	Name        string `json:"name"`
	Namespace   string `json:"nameSpace"`
	SiteVersion string `json:"siteVersion"`
	Policy      string `json:"policy"`
}

type RouterStatus struct {
	Router     Router      `json:"router"`
	Links      []Link      `json:"links,omitempty"`
	Listeners  interface{} `json:"listeners"`
	Connectors interface{} `json:"connectors"`
}

type Router struct {
	RecType      string `json:"recType"`
	Identity     string `json:"identity"`
	Parent       string `json:"parent"`
	StartTime    int64  `json:"startTime"`
	EndTime      int64  `json:"endTime"`
	Source       string `json:"source"`
	Name         string `json:"name"`
	Namespace    string `json:"namespace,omitempty"`
	Mode         string `json:"mode"`
	ImageName    string `json:"imageName,omitempty"`
	ImageVersion string `json:"imageVersion"`
	Hostname     string `json:"hostname"`
	BuildVersion string `json:"buildVersion"`
}

type Link struct {
	RecType   string `json:"recType"`
	Identity  string `json:"identity"`
	Parent    string `json:"parent"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
	Source    string `json:"source"`
	Mode      string `json:"mode"`
	Name      string `json:"name"`
	Direction string `json:"direction"`
	Cost      int    `json:"cost"` // TODO test with cost populated
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

	// Unmarshal YAML into NetworkStatusDataUnmarshaled to extract the JSON string
	var networkStatusDataUnmarshaled NetworkStatusDataUnmarshaled
	err = yaml.Unmarshal(data, &networkStatusDataUnmarshaled)
	if err != nil {
		log.Printf("Error unmarshalling YAML: %+v", err)
		return nil, fmt.Errorf("Error unmarshalling file %s data: %v", filename, err)
	}

	log.Printf("TMPDBG: NetworkStatus JSON: %s", networkStatusDataUnmarshaled.Data.NetworkStatus)

	// Unmarshal the JSON part of the NetworkStatus
	err = json.Unmarshal([]byte(networkStatusDataUnmarshaled.Data.NetworkStatus), networkStatus)
	if err != nil {
		log.Printf("Error unmarshalling NetworkStatus JSON: %+v", err)
		return nil, fmt.Errorf("Error unmarshalling NetworkStatus JSON: %v", err)
	}

	log.Printf("TMPDBG: Unmarshalled network status: %+v", networkStatus)

	DumpNetworkStatus(networkStatus)

	return networkStatus, nil
}

func DumpNetworkStatus(networkStatus *NetworkStatus) {
	for i, siteStatus := range networkStatus.SiteStatus {
		log.Printf("TMPDBG: case 50: i=%+v, siteStatus=%+v", i, siteStatus)
		for j, routerStatus := range siteStatus.RouterStatus {
			log.Printf("TMPDBG: case 51: i=%+v, j=%+v, routerStatus=%+v", i, j, routerStatus)
			log.Printf("TMPDBG: case 52: i=%+v, j=%+v, len(routerStatus.Links)=%+v", i, j, len(routerStatus.Links))
			for k, link := range routerStatus.Links {
				log.Printf("TMPDBG: case 53: i=%+v, j=%+v, k=%+v, link=%+v", i, j, k, link)
			}
		}
	}
}
