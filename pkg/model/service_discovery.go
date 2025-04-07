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
	services         []*model.Service
	serviceInstances []*model.ServiceInstance

	model.NoopAmbientIndexes
	model.NetworkGatewaysHandler
}

var _ model.ServiceDiscovery = &LocalServiceDiscovery{}

func (l *LocalServiceDiscovery) Services() []*model.Service {
	return l.services
}

func (l *LocalServiceDiscovery) GetService(host.Name) *model.Service {
	return &model.Service{}
}

func (l *LocalServiceDiscovery) GetProxyServiceTargets(*model.Proxy) []model.ServiceTarget {
	return []model.ServiceTarget{}
}

func (l *LocalServiceDiscovery) GetProxyWorkloadLabels(*model.Proxy) labels.Instance {
	return labels.Instance{}
}

func (l *LocalServiceDiscovery) GetIstioServiceAccounts(*model.Service) []string {
	return []string{}
}

func (l *LocalServiceDiscovery) NetworkGateways() []model.NetworkGateway {
	return []model.NetworkGateway{}
}

func (l *LocalServiceDiscovery) MCSServices() []model.MCSServiceInfo {
	return []model.MCSServiceInfo{}
}

func (l *LocalServiceDiscovery) GetProxyServiceInstances(proxy *model.Proxy) []*model.ServiceInstance {
	if len(proxy.IPAddresses) == 0 {
		return nil
	}

	switch proxy.IPAddresses[0] {
	case "10.0.0.2": // picard
		return []*model.ServiceInstance{
			{
				Service:     l.services[0],
				ServicePort: l.services[0].Ports[0],
				Endpoint: &model.IstioEndpoint{
					Addresses:       []string{l.services[0].DefaultAddress},
					EndpointPort:    uint32(l.services[0].Ports[0].Port),
					ServicePortName: l.services[0].Ports[0].Name,
				},
			},
		}
	case "10.0.0.3": // comms-operator
		return []*model.ServiceInstance{
			{
				Service:     l.services[1],
				ServicePort: l.services[1].Ports[0],
				Endpoint: &model.IstioEndpoint{
					Addresses:       []string{l.services[0].DefaultAddress},
					EndpointPort:    uint32(l.services[1].Ports[0].Port),
					ServicePortName: l.services[1].Ports[0].Name,
				},
			},
		}
	default:
		return nil
	}
}
