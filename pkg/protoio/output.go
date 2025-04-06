package protoio

import (
	"context"
	"fmt"
	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"istio.io/istio/istioctl/pkg/writer/envoy/configdump"
	"os"
)

func WriteProtoJSONList(ctx context.Context, filename string, list []proto.Message, listeners []*listenerv3.Listener) error {
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

	cfgDump, err := buildConfigDumpFromListeners(listeners)

	jDada, _ := protojson.Marshal(cfgDump)

	configWriter := configdump.ConfigWriter{Stdout: f}

	if err := configWriter.Prime(jDada); err != nil {
		log.Fatal().Err(err).Msgf("Failed to prime config writer")
	}

	if err := configWriter.PrintListenerDump(configdump.ListenerFilter{}, "json"); err != nil {
		log.Fatal().Err(err).Msg("failed to print listener dump")
	}

	return nil
}

func buildConfigDumpFromListeners(listeners []*listenerv3.Listener) (*adminv3.ConfigDump, error) {
	listenerDump := adminv3.ListenersConfigDump{}

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

	anyDump, err := anypb.New(&listenerDump)
	if err != nil {
		return nil, err
	}

	return &adminv3.ConfigDump{
		Configs: []*anypb.Any{anyDump},
	}, nil
}
