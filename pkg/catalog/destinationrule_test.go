package catalog

import (
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"testing"

	v1 "istio.io/api/networking/v1"
	"istio.io/istio/pkg/config/schema/gvk"
)

func TestDestinationRuleBuilder(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		tlsMode  v1.ClientTLSSettings_TLSmode
		wantSpec *v1.DestinationRule
	}{
		{
			name:    "basic host",
			host:    "example.com",
			tlsMode: v1.ClientTLSSettings_DISABLE,
			wantSpec: &v1.DestinationRule{
				Host: "example.com",
				TrafficPolicy: &v1.TrafficPolicy{
					Tls: &v1.ClientTLSSettings{Mode: v1.ClientTLSSettings_DISABLE},
				},
			},
		},
		{
			name:    "istio mutual tls",
			host:    "*.svc.cluster.local",
			tlsMode: v1.ClientTLSSettings_ISTIO_MUTUAL,
			wantSpec: &v1.DestinationRule{
				Host: "*.svc.cluster.local",
				TrafficPolicy: &v1.TrafficPolicy{
					Tls: &v1.ClientTLSSettings{Mode: v1.ClientTLSSettings_ISTIO_MUTUAL},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewDestinationRuleBuilder("default", "istio-system").
				WithHost(tt.host).
				WithTlsMode(tt.tlsMode).
				Build()

			if cfg.GroupVersionKind != gvk.DestinationRule {
				t.Errorf("unexpected GVK: got %+v", cfg.GroupVersionKind)
			}

			got, ok := cfg.Spec.(*v1.DestinationRule)
			if !ok {
				t.Fatalf("spec is not *DestinationRule, got %T", cfg.Spec)
			}

			want := tt.wantSpec

			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Build() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
