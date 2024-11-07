package site

import (
	"testing"

	skupperv2alpha1 "github.com/skupperproject/skupper/pkg/apis/skupper/v2alpha1"
	"github.com/skupperproject/skupper/pkg/qdr"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MockBindingContext struct {
	selectors map[string]TargetSelection
}

func NewMockBindingContext(selectors map[string]TargetSelection) *MockBindingContext {
	return &MockBindingContext{
		selectors: selectors,
	}
}

func (m *MockBindingContext) Select(connector *skupperv2alpha1.Connector) TargetSelection {
	if selector, ok := m.selectors[connector.Name]; ok {
		return selector
	}
	return nil
}

func (m *MockBindingContext) Expose(ports *ExposedPortSet) {
}

func (m *MockBindingContext) Unexpose(host string) {
}

func TestBindingAdaptor_ConnectorUpdated(t *testing.T) {
	type fields struct {
		context   BindingContext
		mapping   *qdr.PortMapping
		exposed   ExposedPorts
		selectors map[string]TargetSelection
	}
	type args struct {
		connector *skupperv2alpha1.Connector
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "connector selector not populated, no matching pods",
			args: args{
				connector: &skupperv2alpha1.Connector{
					ObjectMeta: v1.ObjectMeta{
						Name:      "backend",
						Namespace: "test",
						UID:       "8a96ffdf-403b-4e4a-83a8-97d3d459adb6",
					},
				},
			},
			want: true,
		},
		{
			name: "connector selector populated, no matching pods",
			fields: fields{
				context:   NewMockBindingContext(nil),
				selectors: map[string]TargetSelection{},
			},
			args: args{
				connector: &skupperv2alpha1.Connector{
					ObjectMeta: v1.ObjectMeta{
						Name:      "backend",
						Namespace: "test",
						UID:       "8a96ffdf-403b-4e4a-83a8-97d3d459adb6",
					},
					Spec: skupperv2alpha1.ConnectorSpec{
						Selector: "app=backend",
					},
				},
			},
			want: false,
		},
		{
			name: "connector selector populated, pods match selector",
			fields: fields{
				context: NewMockBindingContext(map[string]TargetSelection{
					"backend": &TargetSelectionImpl{
						selector: "app=backend",
					},
				}),
				selectors: map[string]TargetSelection{
					"backend": &TargetSelectionImpl{
						selector: "app=backend",
					},
				},
			},
			args: args{
				connector: &skupperv2alpha1.Connector{
					ObjectMeta: v1.ObjectMeta{
						Name:      "backend",
						Namespace: "test",
						UID:       "8a96ffdf-403b-4e4a-83a8-97d3d459adb6",
					},
					Spec: skupperv2alpha1.ConnectorSpec{
						Selector: "app=backend",
					},
				},
			},
			want: true,
		},
		{
			name: "connector selector changed to empty",
			fields: fields{
				context: NewMockBindingContext(map[string]TargetSelection{
					"backend": &TargetSelectionImpl{
						selector: "app=backend",
					},
				}),
				selectors: map[string]TargetSelection{
					"backend": &TargetSelectionImpl{
						selector: "app=backend",
						watcher: &PodWatcher{
							// create channel so thta it can be closed without panic.
							stopCh: make(chan struct{}),
						},
					},
				},
			},
			args: args{
				connector: &skupperv2alpha1.Connector{
					ObjectMeta: v1.ObjectMeta{
						Name:      "backend",
						Namespace: "test",
						UID:       "8a96ffdf-403b-4e4a-83a8-97d3d459adb6",
					},
					Spec: skupperv2alpha1.ConnectorSpec{},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &BindingAdaptor{
				context:   tt.fields.context,
				mapping:   tt.fields.mapping,
				exposed:   tt.fields.exposed,
				selectors: tt.fields.selectors,
			}
			if got := a.ConnectorUpdated(tt.args.connector); got != tt.want {
				t.Errorf("BindingAdaptor.ConnectorUpdated() = %v, want %v", got, tt.want)
			}
		})
	}
}
