package site

import (
	"log/slog"
	"testing"

	skupperv2alpha1 "github.com/skupperproject/skupper/pkg/apis/skupper/v2alpha1"
	"github.com/skupperproject/skupper/pkg/kube"
	"github.com/skupperproject/skupper/pkg/site"
)

func TestExtendedBindings_attachedConnectorUpdated(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
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
