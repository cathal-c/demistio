package catalog

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	v1 "istio.io/api/networking/v1"
	"istio.io/istio/pkg/config/schema/gvk"
)

func TestServiceEntryBuilder(t *testing.T) {
	tests := []struct {
		name     string
		hosts    []string
		wantSpec *v1.ServiceEntry
	}{
		{
			name:  "single host",
			hosts: []string{"example.com"},
			wantSpec: &v1.ServiceEntry{
				Hosts: []string{"example.com"},
			},
		},
		{
			name:  "multiple hosts",
			hosts: []string{"foo.com", "bar.com"},
			wantSpec: &v1.ServiceEntry{
				Hosts: []string{"foo.com", "bar.com"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewServiceEntryBuilder("test-entry", "default").
				WithHosts(tt.hosts...).
				Build()

			if cfg.GroupVersionKind != gvk.ServiceEntry {
				t.Errorf("unexpected GVK: got %+v", cfg.GroupVersionKind)
			}

			got, ok := cfg.Spec.(*v1.ServiceEntry)
			if !ok {
				t.Fatalf("spec is not *ServiceEntry, got %T", cfg.Spec)
			}

			want := tt.wantSpec

			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Build() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
