package catalog

import (
	"istio.io/api/networking/v1"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/schema/gvk"
)

type ServiceEntryBuilder struct {
	cfg config.Config
}

func NewServiceEntryBuilder(name, namespace string) *ServiceEntryBuilder {
	return &ServiceEntryBuilder{
		cfg: config.Config{
			Meta: config.Meta{
				GroupVersionKind: gvk.ServiceEntry,
				Name:             name,
				Namespace:        namespace,
			},
			Spec: &v1.ServiceEntry{},
		},
	}
}

func (b *ServiceEntryBuilder) WithHosts(hosts ...string) *ServiceEntryBuilder {
	se := b.cfg.Spec.(*v1.ServiceEntry)

	se.Hosts = append(se.Hosts, hosts...)

	return b
}

func (b *ServiceEntryBuilder) Build() config.Config {
	return b.cfg
}
