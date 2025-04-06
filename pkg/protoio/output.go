package protoio

import (
	"context"
	"fmt"
	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"os"
)

func WriteProtoJSONList(ctx context.Context, filename string, listeners []*listenerv3.Listener, routes []*discoveryv3.Resource) error {
	log := zerolog.Ctx(ctx)

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", filename, err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close file")
		}
	}(f)

	cfgDump, err := buildConfigDump(listeners, routes)

	if true {
		data, err := protojson.MarshalOptions{
			Multiline:       true,
			Indent:          "  ",
			EmitUnpopulated: false,
			UseProtoNames:   true,
		}.Marshal(cfgDump)
		if err != nil {
			return fmt.Errorf("failed to marshal config dump: %w", err)
		}

		if _, err := f.Write(data); err != nil {
			return fmt.Errorf("failed to write config dump: %w", err)
		}
	}

	return nil
}

func buildConfigDump(listeners []*listenerv3.Listener, routes []*discoveryv3.Resource) (*adminv3.ConfigDump, error) {
	listenerDump := &adminv3.ListenersConfigDump{}
	routesDump := &adminv3.RoutesConfigDump{}

	for _, l := range listeners {
		anyListener, err := anypb.New(l)
		if err != nil {
			return nil, err
		}
		listenerDump.DynamicListeners = append(listenerDump.DynamicListeners, &adminv3.ListenersConfigDump_DynamicListener{
			Name:         l.Name,
			ActiveState:  &adminv3.ListenersConfigDump_DynamicListenerState{Listener: anyListener},
			ClientStatus: adminv3.ClientResourceStatus_ACKED,
		})
	}

	for _, route := range routes {
		rc := &routev3.RouteConfiguration{}
		if err := route.Resource.UnmarshalTo(rc); err != nil {
			return nil, fmt.Errorf("unmarshal route resource: %w", err)
		}
		anyRoute, err := anypb.New(rc)
		if err != nil {
			return nil, err
		}
		routesDump.DynamicRouteConfigs = append(routesDump.DynamicRouteConfigs, &adminv3.RoutesConfigDump_DynamicRouteConfig{
			RouteConfig:  anyRoute,
			ClientStatus: adminv3.ClientResourceStatus_ACKED,
		})
	}

	lDump, err := anypb.New(listenerDump)
	if err != nil {
		return nil, err
	}
	rDump, err := anypb.New(routesDump)
	if err != nil {
		return nil, err
	}

	return &adminv3.ConfigDump{
		Configs: []*anypb.Any{lDump, rDump},
	}, nil
}
