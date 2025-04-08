package catalog

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	v1 "istio.io/api/security/v1"
	"istio.io/istio/pkg/config/schema/gvk"
)

func TestPeerAuthenticationBuilder(t *testing.T) {
	tests := []struct {
		name     string
		mtlsMode v1.PeerAuthentication_MutualTLS_Mode
		wantSpec *v1.PeerAuthentication
	}{
		{
			name:     "mtls strict",
			mtlsMode: v1.PeerAuthentication_MutualTLS_STRICT,
			wantSpec: &v1.PeerAuthentication{
				Mtls: &v1.PeerAuthentication_MutualTLS{
					Mode: v1.PeerAuthentication_MutualTLS_STRICT,
				},
			},
		},
		{
			name:     "mtls permissive",
			mtlsMode: v1.PeerAuthentication_MutualTLS_PERMISSIVE,
			wantSpec: &v1.PeerAuthentication{
				Mtls: &v1.PeerAuthentication_MutualTLS{
					Mode: v1.PeerAuthentication_MutualTLS_PERMISSIVE,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewPeerAuthenticationBuilder("default", "istio-system").
				WithMutualTlsMode(tt.mtlsMode).
				Build()

			if cfg.GroupVersionKind != gvk.PeerAuthentication {
				t.Errorf("unexpected GVK: got %+v", cfg.GroupVersionKind)
			}

			got, ok := cfg.Spec.(*v1.PeerAuthentication)
			if !ok {
				t.Fatalf("spec is not *PeerAuthentication, got %T", cfg.Spec)
			}

			if diff := cmp.Diff(tt.wantSpec, got, protocmp.Transform()); diff != "" {
				t.Errorf("Build() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
