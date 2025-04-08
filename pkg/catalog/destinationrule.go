package catalog

import (
	"istio.io/api/networking/v1"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/schema/gvk"
)

type DestinationRuleBuilder struct {
	cfg config.Config
}

func NewDestinationRuleBuilder(name, namespace string) *DestinationRuleBuilder {
	return &DestinationRuleBuilder{
		cfg: config.Config{
			Meta: config.Meta{
				GroupVersionKind: gvk.DestinationRule,
				Name:             name,
				Namespace:        namespace,
			},
			Spec: &v1.DestinationRule{},
		},
	}
}

func (b *DestinationRuleBuilder) WithHost(host string) *DestinationRuleBuilder {
	se := b.cfg.Spec.(*v1.DestinationRule)

	se.Host = host

	return b
}

func (b *DestinationRuleBuilder) WithTlsMode(mode v1.ClientTLSSettings_TLSmode) *DestinationRuleBuilder {
	se := b.cfg.Spec.(*v1.DestinationRule)

	se.TrafficPolicy = &v1.TrafficPolicy{
		Tls: &v1.ClientTLSSettings{
			Mode: mode,
		},
	}

	return b
}

func (b *DestinationRuleBuilder) Build() config.Config {
	return b.cfg
}
