package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skupperproject/skupper/api/types"
	"github.com/skupperproject/skupper/internal/cmd/skupper/common/utils"
	"github.com/skupperproject/skupper/internal/kube/client"
	"github.com/skupperproject/skupper/pkg/apis/skupper/v2alpha1"

	//corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//"strconv"

	//"bufio"
	//"bytes"
	//"k8s.io/apimachinery/pkg/runtime/serializer/json"
	//"k8s.io/client-go/kubernetes/scheme"
	//
	// jsonencoding "encoding/json"
	"encoding/json"
	"log"
)

func readSkupperServices(cli *client.KubeClient, namespace string) ([]*types.ServiceInterface, error) {
	cm, err := readConfigMap(context.Background(), namespace, types.ServiceInterfaceConfigMap, cli)
	if err != nil {
		return nil, err
	}
	var services []*types.ServiceInterface

	for _, value := range cm.Data {
		var service types.ServiceInterface
		err := json.Unmarshal([]byte(value), &service)
		if err != nil {
			return nil, err
		}
		services = append(services, &service)
	}

	return services, nil
}

func upgradeSkupperServices(cli *client.KubeClient, siteConfig *types.SiteConfig, outputPath string, uidToSiteConfig map[string]*types.SiteConfig) error {
	services, err := readSkupperServices(cli, siteConfig.Spec.SkupperNamespace)
	if err != nil {
		return err
	}

	for _, service := range services {
		log.Printf("TMPDBG: upgradeSkupperServices: siteConfig.Spec.SkupperNamespace=%+v, service=%+v", siteConfig.Spec.SkupperNamespace, service)
		resources, err := createListenerCR(siteConfig, service)
		if err != nil {
			return err
		}
		err = saveListenerCR(resources, outputPath, siteConfig.Spec.SkupperName)
		if err != nil {
			return err
		}
		if len(service.Targets) > 0 {
			resources, err := createConnectorCR(siteConfig, service)
			if err != nil {
				return err
			}
			err = saveConnectorCR(resources, outputPath, siteConfig.Spec.SkupperName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// TMPDBG ./internal/cmd/skupper/listener/kube/listener_create.go
// TMPDBG anything to validate here?
func createListenerCR(siteConfig *types.SiteConfig, service *types.ServiceInterface) ([]*v2alpha1.Listener, error) {
	var resources []*v2alpha1.Listener

	for _, port := range service.Ports {
		resource := &v2alpha1.Listener{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "skupper.io/v2alpha1",
				Kind:       "Listener",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      service.Address, // Is this right?
				Namespace: siteConfig.Spec.SkupperNamespace,
			},
			Spec: v2alpha1.ListenerSpec{
				//Host:           cmd.host,
				Port:       port,
				RoutingKey: service.Address,
				//TlsCredentials: cmd.tlsSecret,
				Type: service.Protocol,
			},
		}
		log.Printf("TMPDBG: createListenerCR: siteConfig.Spec.SkupperNamespace=%+v, resource=%+v", siteConfig.Spec.SkupperNamespace, resource)
		resources = append(resources, resource)
	}
	return resources, nil
}

// Lifted from pkg/service/bindings.go.  TODO is this necessary?
func getTargetPorts(service types.ServiceInterface, target types.ServiceInterfaceTarget) map[int]int {
	targetPorts := target.TargetPorts
	if len(targetPorts) == 0 {
		targetPorts = map[int]int{}
		for _, port := range service.Ports {
			targetPorts[port] = port
		}
	}
	return targetPorts
}

func saveListenerCR(resources []*v2alpha1.Listener, outputPath string, siteName string) error {
	if len(resources) <= 0 {
		return nil // TODO fix me
	}
	targetDir := filepath.Join(outputPath, siteName)
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("Failed to create directory: %v", err)
	}

	filepath := filepath.Join(targetDir, fmt.Sprintf("%s_listener.yaml", resources[0].ObjectMeta.Name))

	data, err := utils.Encode("yaml", resources[0]) // TODO Fix me
	if err != nil {
		return fmt.Errorf("Failed to marshal site resource to YAML: %w", err.Error())
	}

	err = os.WriteFile(filepath, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("Failed to write site resource to file: %w", err.Error())
	}
	fmt.Printf("Wrote site CR to file: %s\n", filepath)
	return nil
}

// "./internal/cmd/skupper/connector/kube/connector_create.go"
func createConnectorCR(siteConfig *types.SiteConfig, service *types.ServiceInterface) ([]*v2alpha1.Connector, error) {
	log.Printf("TMPDBG: createConnectorCR: entering, service=%+v", service)
	var resources []*v2alpha1.Connector

	// 2024/10/25 14:56:03 /home/kschoener/work/20241025_upgrade_v1_to_v2/skupper/cmd/upgrade/upgrade_kube_serviceinterface.go:55: TMPDBG: upgradeSkupperServices: siteConfig.Spec.SkupperNamespace=east-v1, service=&{Address:backend Protocol:tcp Ports:[8080] ExposeIngress: EventChannel:false Aggregate: Headless:<nil> Labels:map[] Annotations:map[] Targets:[{Name:backend Selector:app=backend TargetPorts:map[8080:8080] Service: Namespace:east-v1}] Origin: TlsCredentials: TlsCertAuthority: PublishNotReadyAddresses:false BridgeImage:}

	// TODO: lots of questions.  target.Namespace.  huh
	// TODO: lots of questions.  target.Port=map 8080:8080 huh

	for _, target := range service.Targets {
		log.Printf("TMPDBG: case 10: target=%+v", target)
		//portMap := kube.PortLabelStrToMap(target.Ports)
		//portMap := kube.PortLabelStrToMap("")
		//testing := true
		//if testing {
		//portMap = map[int]int{}
		//}
		// TODO check that this target port logic is correct.
		for _, port := range target.TargetPorts {
			resource := &v2alpha1.Connector{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "skupper.io/v2alpha1",
					Kind:       "Connector",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      service.Address, // Is this right?  // TODO connector names must be unique.
					Namespace: siteConfig.Spec.SkupperNamespace,
				},
				Spec: v2alpha1.ConnectorSpec{
					//Host:            cmd.host,
					Port:       port,
					RoutingKey: service.Address,
					//TlsCredentials:  cmd.tlsSecret,
					Type: service.Protocol,
					//IncludeNotReady: cmd.includeNotReady,
					Selector: target.Selector,
				},
			}
			log.Printf("TMPDBG: createConnectorCR: siteConfig.Spec.SkupperNamespace=%+v, resource=%+v", siteConfig.Spec.SkupperNamespace, resource)
			resources = append(resources, resource)
		}
	}
	return resources, nil
}

func saveConnectorCR(resources []*v2alpha1.Connector, outputPath string, siteName string) error {
	if len(resources) <= 0 {
		return nil // TODO fix me
	}
	targetDir := filepath.Join(outputPath, siteName)
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("Failed to create directory: %v", err)
	}

	filepath := filepath.Join(targetDir, fmt.Sprintf("%s_connector.yaml", resources[0].ObjectMeta.Name))

	data, err := utils.Encode("yaml", resources[0]) // TODO Fix me
	if err != nil {
		return fmt.Errorf("Failed to marshal site resource to YAML: %w", err.Error())
	}

	err = os.WriteFile(filepath, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("Failed to write site resource to file: %w", err.Error())
	}
	fmt.Printf("Wrote site CR to file: %s\n", filepath)
	return nil
}
