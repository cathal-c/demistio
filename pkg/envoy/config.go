package envoy

import (
	"fmt"
	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"google.golang.org/protobuf/types/known/anypb"
)

func BuildListenerConfigDump(listeners []*listenerv3.Listener) (*adminv3.ListenersConfigDump, error) {
	listenerDump := &adminv3.ListenersConfigDump{}

	for _, l := range listeners {
		anyListener, err := anypb.New(l)
		if err != nil {
			return nil, fmt.Errorf("build anypb: %w", err)
		}

		listenerDump.DynamicListeners = append(listenerDump.DynamicListeners, &adminv3.ListenersConfigDump_DynamicListener{
			Name:         l.Name,
			ActiveState:  &adminv3.ListenersConfigDump_DynamicListenerState{Listener: anyListener},
			ClientStatus: adminv3.ClientResourceStatus_ACKED,
		})
	}

	return listenerDump, nil
}

func BuildClusterConfigDump(clusters []*discoveryv3.Resource) (*adminv3.ClustersConfigDump, error) {
	clusterDump := &adminv3.ClustersConfigDump{}

	for _, l := range clusters {
		anyCluster, err := anypb.New(l)
		if err != nil {
			return nil, fmt.Errorf("build anypb: %w", err)
		}

		clusterDump.DynamicActiveClusters = append(clusterDump.DynamicActiveClusters, &adminv3.ClustersConfigDump_DynamicCluster{
			Cluster:      anyCluster,
			ClientStatus: adminv3.ClientResourceStatus_ACKED,
		})
	}

	return clusterDump, nil
}

func BuildRoutesConfigDump(routes []*discoveryv3.Resource) (*adminv3.RoutesConfigDump, error) {
	routesDump := &adminv3.RoutesConfigDump{}

	for _, route := range routes {
		rc := &routev3.RouteConfiguration{}
		if err := route.Resource.UnmarshalTo(rc); err != nil {
			return nil, fmt.Errorf("unmarshal route resource: %w", err)
		}

		anyRoute, err := anypb.New(rc)
		if err != nil {
			return nil, fmt.Errorf("build anypb: %w", err)
		}

		routesDump.DynamicRouteConfigs = append(routesDump.DynamicRouteConfigs, &adminv3.RoutesConfigDump_DynamicRouteConfig{
			RouteConfig:  anyRoute,
			ClientStatus: adminv3.ClientResourceStatus_ACKED,
		})
	}

	return routesDump, nil
}

func BuildFullConfigDump(listeners []*listenerv3.Listener, clusters []*discoveryv3.Resource, routes []*discoveryv3.Resource) (*adminv3.ConfigDump, error) {
	listenerDump, err := BuildListenerConfigDump(listeners)
	if err != nil {
		return nil, fmt.Errorf("build listener config dump: %w", err)
	}

	clusterDump, err := BuildClusterConfigDump(clusters)
	if err != nil {
		return nil, fmt.Errorf("build listener config dump: %w", err)
	}

	routesDump, err := BuildRoutesConfigDump(routes)
	if err != nil {
		return nil, fmt.Errorf("build routes config dump: %w", err)
	}

	lDump, err := anypb.New(listenerDump)
	if err != nil {
		return nil, err
	}

	cDump, err := anypb.New(clusterDump)
	if err != nil {
		return nil, err
	}

	rDump, err := anypb.New(routesDump)
	if err != nil {
		return nil, err
	}

	return &adminv3.ConfigDump{
		Configs: []*anypb.Any{lDump, cDump, rDump},
	}, nil
}
