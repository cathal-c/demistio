package catalog

import (
	"istio.io/api/security/v1"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/schema/gvk"
)

var DefaultPeerAuthenticationEnableMutualTLS = config.Config{
	Meta: config.Meta{
		GroupVersionKind: gvk.PeerAuthentication,
		Name:             "default",
		Namespace:        "istio-system",
	},
	Spec: &v1.PeerAuthentication{
		Mtls: &v1.PeerAuthentication_MutualTLS{
			Mode: v1.PeerAuthentication_MutualTLS_STRICT,
		},
	},
}

type PeerAuthenticationBuilder struct {
	cfg config.Config
}

func NewPeerAuthenticationBuilder(name, namespace string) *PeerAuthenticationBuilder {
	return &PeerAuthenticationBuilder{
		cfg: config.Config{
			Meta: config.Meta{
				GroupVersionKind: gvk.PeerAuthentication,
				Name:             name,
				Namespace:        namespace,
			},
			Spec: &v1.PeerAuthentication{},
		},
	}
}

func (b *PeerAuthenticationBuilder) WithMutualTlsMode(mode v1.PeerAuthentication_MutualTLS_Mode) *PeerAuthenticationBuilder {
	se := b.cfg.Spec.(*v1.PeerAuthentication)

	se.Mtls = &v1.PeerAuthentication_MutualTLS{
		Mode: mode,
	}

	return b
}

func (b *PeerAuthenticationBuilder) Build() config.Config {
	return b.cfg
}
