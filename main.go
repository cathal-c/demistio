package main

import (
	"context"
	"flag"
	"github.com/cathal-c/demistio/pkg/model"
	"github.com/cathal-c/demistio/pkg/protoio"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
	v1 "istio.io/api/networking/v1"
	"istio.io/istio/pilot/pkg/config/memory"
	cfgModel "istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking/core"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/mesh"
	"istio.io/istio/pkg/config/mesh/meshwatcher"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/collections"
	"istio.io/istio/pkg/config/schema/gvk"
	"os"
)

var (
	output *string
)

func main() {
	log := setupLogs()
	ctx := log.WithContext(context.Background())

	parseFlags(log)

	// Create an in-memory config store, registering required resource types
	store := memory.Make(collection.SchemasFor(
		collections.Service,
		collections.ServiceEntry,
	))

	if _, err := store.Create(config.Config{
		Meta: config.Meta{
			GroupVersionKind: gvk.ServiceEntry,
			Namespace:        "default",
			Name:             "example-service",
		},
		Spec: &v1.ServiceEntry{
			Hosts: []string{"example.com"},
		},
	}); err != nil {
		log.Fatal().Err(err).Msg("Failed to create ServiceEntry")
	}

	env := cfgModel.NewEnvironment()

	meshConfig := mesh.DefaultMeshConfig()

	env.Watcher = meshwatcher.NewTestWatcher(meshConfig)

	services := []*cfgModel.Service{
		{
			Hostname: "svc1",
			Ports: []*cfgModel.Port{
				{
					Name:     "http",
					Port:     8080,
					Protocol: "HTTP",
				},
			},
			Attributes: cfgModel.ServiceAttributes{
				Name:      "svc1",
				Namespace: "default",
			},
		},
	}

	env.Init()

	env.ServiceDiscovery = model.NewLocalServiceDiscovery(services)
	env.ConfigStore = store

	// Initialize a PushContext with the config store.
	push := cfgModel.NewPushContext()
	push.Mesh = mesh.DefaultMeshConfig()

	push.InitContext(env, nil, nil)

	// Create a dummy Proxy. In a real scenario, this would reflect your proxy’s metadata.
	proxy := &cfgModel.Proxy{
		IPAddresses: []string{"127.0.0.1"},
		Metadata: &cfgModel.NodeMetadata{
			Namespace:    "default",
			IstioVersion: "1.22.0",
		},
		Type:            cfgModel.SidecarProxy,
		DNSDomain:       "default.svc.cluster.local",
		ConfigNamespace: "default",
		IstioVersion:    cfgModel.ParseIstioVersion("1.22.0"),
		SidecarScope:    cfgModel.DefaultSidecarScopeForNamespace(push, "default"),
	}

	proxy.DiscoverIPMode()

	configGen := &core.ConfigGeneratorImpl{}

	// Use Istio's v1alpha3 conversion logic to build Envoy listeners.
	listeners := configGen.BuildListeners(proxy, push)

	if err := protoio.WriteProtoJSONList(ctx, *output, toProtoMessageList(listeners)); err != nil {
		log.Fatal().Err(err).Msg("Failed to generate Envoy config")
	}
}

func parseFlags(log zerolog.Logger) *string {
	// inputPtr := flag.String("input", "", "Path to YAML file containing Istio configs")
	output = flag.String("output", "", "Output file for generated Envoy config (JSON)")
	flag.Parse()

	// os.Args[0] is always the binary
	for _, arg := range os.Args[1:] {
		log.Info().Msg(arg)
	}
	return output
}

func setupLogs() zerolog.Logger {
	log := zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}).With().Timestamp().Logger()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	return log
}

func toProtoMessageList[T proto.Message](in []T) []proto.Message {
	out := make([]proto.Message, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}
