package site

import (
	"log/slog"
	"testing"

	skupperv2alpha1 "github.com/skupperproject/skupper/pkg/apis/skupper/v2alpha1"
	"github.com/skupperproject/skupper/pkg/kube"
	"github.com/skupperproject/skupper/pkg/site"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//"gotest.tools/assert"

	"log"

	fakeclient "github.com/skupperproject/skupper/internal/kube/client/fake"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewMockSite(namespace string) *Site {

	site, err := newSiteMocks(namespace, nil, nil, "", false)

	log.Printf("TMPDBG: NewMockSite: err=%+v", err)

	//assert.Assert(t, err)

	return site
}

func NewMockController(namespace string) (*kube.Controller, error) {

	k8sObjects := []runtime.Object{}
	skupperObjects := []runtime.Object{}
	fakeSkupperError := ""

	client, err := fakeclient.NewFakeClient(namespace, k8sObjects, skupperObjects, fakeSkupperError)
	if err != nil {
		return nil, err
	}

	controller := kube.NewController(namespace, client)

	return controller, nil
}

func TestExtendedBindings_attachedConnectorUpdated(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Llongfile) // for full file path and line numbers
	initController := func(namespace string) *kube.Controller {
		controller, err := NewMockController(namespace)
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
	type args struct {
		name       string
		definition *skupperv2alpha1.AttachedConnector
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "No matching AttachedConnectorAnchor in site namespace",
			fields: fields{
				bindings: &site.Bindings{
				},
				connectors: map[string]*AttachedConnector{},
				controller: initController("test"),
				site:       NewMockSite("test"),
				logger: slog.New(slog.Default().Handler()).With(
					slog.String("component", "kube.site.attached_connector"),
				),
			},
			args: args{
				name: "backend",
				definition: &skupperv2alpha1.AttachedConnector{
					ObjectMeta: v1.ObjectMeta{
						Name:      "backend",
						Namespace: "test-anchor",
					},
				},
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
			if err := b.attachedConnectorUpdated(tt.args.name, tt.args.definition); (err != nil) != tt.wantErr {
				t.Errorf("ExtendedBindings.attachedConnectorUpdated() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
