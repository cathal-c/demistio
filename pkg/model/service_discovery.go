package model

import (
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/config/host"
	"istio.io/istio/pkg/config/labels"
)

func NewLocalServiceDiscovery(services []*model.Service) *LocalServiceDiscovery {
	return &LocalServiceDiscovery{
		services: services,
	}
}

// LocalServiceDiscovery is an in-memory ServiceDiscovery with mock services
type LocalServiceDiscovery struct {
	services []*model.Service
	//serviceInstances []*model.ServiceInstance

	model.NoopAmbientIndexes
	model.NetworkGatewaysHandler
}

// Compile-time interface check. Does nothing other than ensure LocalServiceDiscovery implements required methods.
var _ model.ServiceDiscovery = &LocalServiceDiscovery{}

func (l *LocalServiceDiscovery) Services() []*model.Service {
	return l.services
}

func (l *LocalServiceDiscovery) GetService(hostname host.Name) *model.Service {
	for _, service := range l.services {
		if service.Hostname == hostname {
			return service
		}

		return nil
	}

	return &model.Service{}
}

func (l *LocalServiceDiscovery) GetProxyServiceTargets(proxy *model.Proxy) []model.ServiceTarget {
	res := make([]model.ServiceTarget, 0)

	for _, service := range l.services {
		for _, address := range proxy.IPAddresses {
			if address == service.DefaultAddress {
				for _, port := range service.Ports {
					svcTarget := model.ServiceTarget{
						Service: service,
						Port: model.ServiceInstancePort{
							ServicePort: &model.Port{
								Name:     port.Name,
								Port:     port.Port,
								Protocol: port.Protocol,
							},
							TargetPort: uint32(port.Port),
						},
					}

					res = append(res, svcTarget)
				}
			}
		}
	}

	return res
}

func (l *LocalServiceDiscovery) GetProxyWorkloadLabels(node *model.Proxy) labels.Instance {
	if node != nil {
		return node.Labels
	}

	return labels.Instance{}
}

func (l *LocalServiceDiscovery) NetworkGateways() []model.NetworkGateway {
	return []model.NetworkGateway{}
}

func (l *LocalServiceDiscovery) MCSServices() []model.MCSServiceInfo {
	return []model.MCSServiceInfo{}
}
