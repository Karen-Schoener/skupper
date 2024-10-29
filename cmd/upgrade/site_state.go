package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/skupperproject/skupper/pkg/apis/skupper/v2alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

// Code in this file was lifted from pkg/nonkube/api/site_state.go.
// TODO: discuss whether upgrade should should call the site_state apis from pkg/nonkube/api.

type SiteState struct {
	//SiteId          string
	Site       *v2alpha1.Site
	Listeners  map[string]*v2alpha1.Listener
	Connectors map[string]*v2alpha1.Connector
	//RouterAccesses  map[string]*v2alpha1.RouterAccess
	Grants map[string]*v2alpha1.AccessGrant
	Tokens map[string]*v2alpha1.AccessToken
	//Links           map[string]*v2alpha1.Link
	//Secrets         map[string]*corev1.Secret
	//Claims          map[string]*v2alpha1.AccessToken
	//Certificates    map[string]*v2alpha1.Certificate
	//SecuredAccesses map[string]*v2alpha1.SecuredAccess
	//bundle          bool
}

func NewSiteState() *SiteState {
	return &SiteState{
		Site:       &v2alpha1.Site{},
		Listeners:  make(map[string]*v2alpha1.Listener),
		Connectors: make(map[string]*v2alpha1.Connector),
		//RouterAccesses:  map[string]*v2alpha1.RouterAccess{},
		Grants: make(map[string]*v2alpha1.AccessGrant),
		Tokens: make(map[string]*v2alpha1.AccessToken),
		//Links:           make(map[string]*v2alpha1.Link),
		//Secrets:         make(map[string]*corev1.Secret),
		//Claims:          make(map[string]*v2alpha1.AccessToken),
		//Certificates:    map[string]*v2alpha1.Certificate{},
		//SecuredAccesses: map[string]*v2alpha1.SecuredAccess{},
		//bundle:          bundle,
	}
}

func Render(siteState *SiteState, outputPath string) error {
	outputDirectory := filepath.Join(outputPath, siteState.Site.ObjectMeta.Name)
	if err := marshal(outputDirectory, "site", siteState.Site.ObjectMeta.Name, siteState.Site); err != nil {
		return err
	}

	siteOutputDirectory := filepath.Join(outputDirectory, siteState.Site.ObjectMeta.Name)

	if err := marshalMap(siteOutputDirectory, "accessgrants", siteState.Grants); err != nil {
		return err
	}
	if err := marshalMap(siteOutputDirectory, "accesstokens", siteState.Tokens); err != nil {
		return err
	}
	if err := marshalMap(siteOutputDirectory, "listeners", siteState.Listeners); err != nil {
		return err
	}
	if err := marshalMap(siteOutputDirectory, "connectors", siteState.Connectors); err != nil {
		return err
	}

	return nil
}

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
