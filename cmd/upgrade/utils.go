package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/skupperproject/skupper/internal/kube/client"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	"fmt"
	"os"
	"path"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func getKubeConfig() string {
	var kubeconfig string // TODO
	return kubeconfig
}

func readConfigMap(ctx context.Context, namespace string, name string, cli *client.KubeClient) (*corev1.ConfigMap, error) {
	kubeClient := cli.GetKubeClient()
	cm, err := kubeClient.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return cm, err
}

// Routine marshal was lifted from pkg/nonkube/api/site_state.go.
// TODO consider creating a SiteState object and calling NewSiteState / MarshalSiteState
func marshal(outputDirectory, resourceType, resourceName string, resource interface{}) error {
	var err error
	writeDirectory := path.Join(outputDirectory, resourceType)
	err = os.MkdirAll(writeDirectory, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory %s: %w", writeDirectory, err)
	}
	fileName := path.Join(writeDirectory, fmt.Sprintf("%s.yaml", resourceName))
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", fileName, err)
	}
	yaml := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	err = yaml.Encode(resource.(runtime.Object), file)
	if err != nil {
		return fmt.Errorf("error marshalling resource %s: %w", resourceName, err)
	}
	fmt.Printf("Wrote CR to file: %s\n", fileName)
	return nil
}

// Note: Routine marshalMap was lifted from pkg/nonkube/api/site_state.go.
func marshalMap[V any](outputDirectory, resourceType string, resourceMap map[string]V) error {
	var err error
	for resourceName, resource := range resourceMap {
		if err = marshal(outputDirectory, resourceType, resourceName, resource); err != nil {
			return err
		}
	}
	return nil
}
