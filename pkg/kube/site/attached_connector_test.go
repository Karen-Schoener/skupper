package site

import (
	//"github.com/davecgh/go-spew/spew"
	"log"
	"log/slog"
	"testing"

	fakeclient "github.com/skupperproject/skupper/internal/kube/client/fake"
	skupperv2alpha1 "github.com/skupperproject/skupper/pkg/apis/skupper/v2alpha1"
	"github.com/skupperproject/skupper/pkg/kube"
	"github.com/skupperproject/skupper/pkg/site"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"gotest.tools/assert"

	"github.com/skupperproject/skupper/pkg/kube/certificates"
	"github.com/skupperproject/skupper/pkg/kube/securedaccess"
	corev1 "k8s.io/api/core/v1"
        //"github.com/skupperproject/skupper/internal/kube/client"


        "context"
        "github.com/skupperproject/skupper/api/types"
        //"github.com/skupperproject/skupper/internal/kube/client"
        "k8s.io/apimachinery/pkg/api/errors"
        "k8s.io/client-go/kubernetes"

)

func getTestSiteCr() *skupperv2alpha1.Site {
		return &skupperv2alpha1.Site{
			ObjectMeta: v1.ObjectMeta{
				Name:      "test",
				Namespace: "test",
				//UID:       "8a96ffdf-403b-4e4a-83a8-97d3d459adb6",
			},
			//Spec: skupperv2alpha1.SiteSpec{
				//DefaultIssuer: "skupper-spec-issuer-ca",
			//},
			//Status: skupperv2alpha1.SiteStatus{
				//DefaultIssuer: "skupper-status-issuer-ca",
			//},
		}
	}

func newSiteMocksCopy(name string, namespace string, k8sObjects []runtime.Object, skupperObjects []runtime.Object) (*Site, error) {
	siteCr := getTestSiteCr()
	skupperObjects = append(skupperObjects, siteCr)

	client, err := fakeclient.NewFakeClient(namespace, k8sObjects, skupperObjects, "")
	if err != nil {
		return nil, err
	}

	controller := kube.NewController("test", client)
	newSite := &Site{
		controller: controller,
		bindings:   NewExtendedBindings(controller, ""),
		links:      make(map[string]*site.Link),
		errors:     make(map[string]string),
		linkAccess: make(map[string]*skupperv2alpha1.RouterAccess),
		certs:      certificates.NewCertificateManager(controller),
		access:     securedaccess.NewSecuredAccessManager(client, nil, &securedaccess.Config{DefaultAccessType: "loadbalancer"}, securedaccess.ControllerContext{}),
		adaptor:    BindingAdaptor{},
		routerPods: make(map[string]*corev1.Pod),
		logger: slog.New(slog.Default().Handler()).With(
			slog.String("component", "kube.site.site"),
		),
	}

	newSite.site = siteCr
	newSite.name = siteCr.ObjectMeta.Name
	newSite.namespace = siteCr.ObjectMeta.Namespace

	return newSite, nil
}

func NewMockSite(name string, namespace string, k8sObjects []runtime.Object, skupperObjects []runtime.Object) *Site {
	site, err := newSiteMocksCopy(name, namespace, k8sObjects, skupperObjects)

	log.Printf("TMPDBG: NewMockSite: err=%+v", err)

	//assert.Assert(t, err)
	//site.initialised = true

	return site
}

func NewMockController(namespace string, k8sObjects []runtime.Object, skupperObjects []runtime.Object) (*kube.Controller, error) {
        client, err := fakeclient.NewFakeClient(namespace, k8sObjects, skupperObjects, "")
        if err != nil {
                return nil, err
        }

	controller := kube.NewController(namespace, client)

	return controller, nil
}

func TestExtendedBindings_attachedConnectorUpdated(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Llongfile) // for full file path and line numbers
	initAttachedConnectorCr := func() *skupperv2alpha1.AttachedConnector {
		return &skupperv2alpha1.AttachedConnector{
			ObjectMeta: v1.ObjectMeta{
				Name:      "backend",
				Namespace: "site-with-backend-pods",
				UID:       "00000000-0000-0000-0000-000000000002",
			},
			Spec: skupperv2alpha1.AttachedConnectorSpec{
				Port:          8080,
				Selector:      "app=backend",
				SiteNamespace: "test",
			},
		}
	}
	initAttachedConnectorCrSpecChanged := func() *skupperv2alpha1.AttachedConnector {
		a := initAttachedConnectorCr()
		a.Spec.Port = 8081
		return a
	}
	initAttachedConnectorAnchorCr := func() *skupperv2alpha1.AttachedConnectorAnchor {
		return &skupperv2alpha1.AttachedConnectorAnchor{
			ObjectMeta: v1.ObjectMeta{
				Name:      "backend",
				Namespace: "test",
			},
			Spec: skupperv2alpha1.AttachedConnectorAnchorSpec{
				RoutingKey:         "backend",
				ConnectorNamespace: "site-with-backend-pods",
			},
		}
	}
        initController := func(namespace string, k8sObjects []runtime.Object, skupperObjects []runtime.Object) *kube.Controller {
                controller, err := NewMockController(namespace, k8sObjects, skupperObjects)
		assert.Assert(t, err)
		if err != nil {
			t.Fatalf("Failed to initialzie controller: %v", err)
		}
		return controller
	}
	type fields struct {
		bindings   *site.Bindings
		connectors map[string]*AttachedConnector
		controller *kube.Controller
		site       *Site
		logger     *slog.Logger
	}
	type testInitSteps struct {
		setWatcher  bool
		recoverSite bool
	}
	type args struct {
		name       string
		definition *skupperv2alpha1.AttachedConnector
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		testInitSteps testInitSteps
	}{
		{
			name: "With matching AttachedConnectorAnchor in site namespace",
			testInitSteps: testInitSteps{
				setWatcher:  true,
				recoverSite: true,
			},
			fields: fields{
				bindings: &site.Bindings{},
				connectors: map[string]*AttachedConnector{
					// Attached connector with anchor
					"backend": &AttachedConnector{
						anchor: initAttachedConnectorAnchorCr(),
						definitions: map[string]*skupperv2alpha1.AttachedConnector{
							//"site-with-backend-pods": initAttachedConnectorCr(),
							"site-with-backend-pods": initAttachedConnectorCrSpecChanged(),
						},
						// parent is set in setParent api...
						// parent: &ExtendedBindings{ // TODO: circular weirdness
						//		logger: slog.New(slog.Default().Handler()).With(
						//			slog.String("component", "kube.site.attached_connector"),
						//			),
						//},
					},
				},
				controller: initController("test", 
					[]runtime.Object{},
					[]runtime.Object{
						initAttachedConnectorCr(),
						initAttachedConnectorAnchorCr(),
					},
				),
				site: NewMockSite("test", "test",
					[]runtime.Object{},
					[]runtime.Object{
						initAttachedConnectorCr(),
						initAttachedConnectorAnchorCr(),
					},
				),
				logger: slog.New(slog.Default().Handler()).With(
					slog.String("component", "kube.site.attached_connector"),
				),
			},
			args: args{
				name:       "backend",
				definition: initAttachedConnectorCr(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		log.Printf("TMPDBG: binding_test: tt.name=%+v", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			b := &ExtendedBindings{
				bindings:   tt.fields.bindings,
				connectors: tt.fields.connectors,
				controller: tt.fields.controller,
				site:       tt.fields.site,
				logger:     tt.fields.logger,
			}
			updateParent(b) // TMPDBG TODO BETTER NAME
			if tt.testInitSteps.recoverSite {
				err := b.site.Recover(getTestSiteCr())
				assert.Assert(t, err)
			}
			if tt.testInitSteps.setWatcher {
				setWatcher(b, tt.args.definition.ObjectMeta.Namespace)
			}
			if err := b.attachedConnectorUpdated(tt.args.name, tt.args.definition); (err != nil) != tt.wantErr {
				t.Errorf("ExtendedBindings.attachedConnectorUpdated() error = %v, wantErr %v", err, tt.wantErr)
			}

			namespace := "test" 
                	cm, err := readConfigMap(context.Background(), namespace, types.TransportConfigMapName, tt.fields.site.controller.GetKubeClient())
			log.Printf("TMPDBG: after readConfigMap: err=%+v", err)
			log.Printf("TMPDBG: after readConfigMap: cm=%+v", cm)


		})
	}
}

func readConfigMap(ctx context.Context, namespace string, name string, kubeClient kubernetes.Interface) (*corev1.ConfigMap, error) {
//func readConfigMap(ctx context.Context, namespace string, name string, client *client.KubeClient) (*corev1.ConfigMap, error) {
//func readConfigMap(ctx context.Context, namespace string, name string, kubeClient kubernetes.Interface) (*corev1.ConfigMap, error) {
        //kubeClient := client.GetKubeClient()
        cm, err := kubeClient.CoreV1().ConfigMaps(namespace).Get(ctx, name, v1.GetOptions{})
        if errors.IsNotFound(err) {
                return nil, nil
        } else if err != nil {
                return nil, err
        }
        return cm, err
}

func updateParent(b *ExtendedBindings) {
	for _, c := range b.connectors {
		c.parent = b
	}
}

func setWatcher(b *ExtendedBindings, namespace string) {
	for _, c := range b.connectors {
		//log.Printf("TMPDBG: setWatcher: case 1234: spew.Sdump(c)=%+v", spew.Sdump(c))
		//log.Printf("TMPDBG: setWatcher: c.ObjectMeta.Name=%+v", c.ObjectMeta.Name)
		//log.Printf("TMPDBG: setWatcher: c.ObjectMeta.Namespace=%+v", c.ObjectMeta.Namespace)
		c.watcher = c.parent.site.WatchPods(c, namespace)
	}
}
