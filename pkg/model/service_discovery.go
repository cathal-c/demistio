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
	panic("implement me")
}

func (l *LocalServiceDiscovery) GetProxyServiceTargets(*model.Proxy) []model.ServiceTarget {
	panic("implement me")
}

func (l *LocalServiceDiscovery) GetProxyWorkloadLabels(*model.Proxy) labels.Instance {
	panic("implement me")
}

func (l *LocalServiceDiscovery) GetIstioServiceAccounts(*model.Service) []string {
	return nil
}

func (l *LocalServiceDiscovery) NetworkGateways() []model.NetworkGateway {
	// TODO implement fromRegistry logic from kube controller if needed
	return nil
}

func (l *LocalServiceDiscovery) MCSServices() []model.MCSServiceInfo {
	return nil
}

func (l *LocalServiceDiscovery) GetProxyServiceInstances(proxy *model.Proxy) []*model.ServiceInstance {
	if len(proxy.IPAddresses) == 0 {
		return nil
	}

	switch proxy.IPAddresses[0] {
	case "10.0.0.2":
		return []*model.ServiceInstance{
			{
				Service:     l.services[0], // picard
				ServicePort: l.services[0].Ports[0],
				Endpoint: &model.IstioEndpoint{
					Addresses:       []string{"10.0.0.2"},
					EndpointPort:    uint32(l.services[0].Ports[0].Port),
					ServicePortName: l.services[0].Ports[0].Name,
				},
			},
		}
	case "10.0.0.3":
		return []*model.ServiceInstance{
			{
				Service:     l.services[1], // comms-operator
				ServicePort: l.services[1].Ports[0],
				Endpoint: &model.IstioEndpoint{
					Addresses:       []string{"10.0.0.3"},
					EndpointPort:    uint32(l.services[1].Ports[0].Port),
					ServicePortName: l.services[1].Ports[0].Name,
				},
			},
		}
	default:
		return nil
	}
}
