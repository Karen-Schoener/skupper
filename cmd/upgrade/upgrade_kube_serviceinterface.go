package main

import (
	"context"
	"fmt"

	"github.com/skupperproject/skupper/api/types"
	"github.com/skupperproject/skupper/internal/kube/client"
	"github.com/skupperproject/skupper/pkg/apis/skupper/v2alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"encoding/json"
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

func createServiceCRs(cli *client.KubeClient, siteConfig *types.SiteConfig, v2SiteState *SiteState) error {
	services, err := readSkupperServices(cli, siteConfig.Spec.SkupperNamespace)
	if err != nil {
		return err
	}

	for _, service := range services {
		v2SiteState.Listeners, err = createListenerCRs(siteConfig, service)
		if err != nil {
			return err
		}
		if len(service.Targets) > 0 {
			v2SiteState.Connectors, err = createConnectorCRs(siteConfig, service)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func createListenerCRs(siteConfig *types.SiteConfig, service *types.ServiceInterface) (map[string]*v2alpha1.Listener, error) {
	resources := map[string]*v2alpha1.Listener{}

	for i, port := range service.Ports {
		resource := &v2alpha1.Listener{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "skupper.io/v2alpha1",
				Kind:       "Listener",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: service.Address, // Is this right?
				//Namespace: siteConfig.Spec.SkupperNamespace,
			},
			Spec: v2alpha1.ListenerSpec{
				Host:       service.Address, // TODO is this correct?
				Port:       port,
				RoutingKey: service.Address,
				//TlsCredentials: cmd.tlsSecret,
				Type: service.Protocol,
			},
		}
		name := resource.ObjectMeta.Name
		if i > 0 {
			// TODO check if this logic is acceptable
			name = fmt.Sprintf("%s-%d", name, port)
		}
		resources[resource.ObjectMeta.Name] = resource
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

func createConnectorCRs(siteConfig *types.SiteConfig, service *types.ServiceInterface) (map[string]*v2alpha1.Connector, error) {
	resources := map[string]*v2alpha1.Connector{}

	for _, target := range service.Targets {
		for i, port := range target.TargetPorts {
			resource := &v2alpha1.Connector{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "skupper.io/v2alpha1",
					Kind:       "Connector",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: service.Address,
					// TODO confirm that namespace should not be populated.
					//Namespace: siteConfig.Spec.SkupperNamespace,
				},
				Spec: v2alpha1.ConnectorSpec{
					//Host:            cmd.host,
					Port:       port,
					RoutingKey: service.Address,
					//TlsCredentials:  cmd.tlsSecret,
					// TODO ok to populate protocol if tcp?  I notice v2 CLI did not appear to populate 'tcp'
					Type: service.Protocol,
					//IncludeNotReady: cmd.includeNotReady,
					Selector: target.Selector,
				},
			}
			name := resource.ObjectMeta.Name
			if i > 0 {
				// TODO check if this logic is acceptable
				// Actually will not work well for multiple targets.
				name = fmt.Sprintf("%s-%d", name, port)
			}
			resources[resource.ObjectMeta.Name] = resource
		}
	}
	return resources, nil
}
