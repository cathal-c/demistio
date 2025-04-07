package main

import (
	"context"
	"flag"
	"github.com/cathal-c/demistio/pkg/model"
	"github.com/cathal-c/demistio/pkg/protoio"
	"github.com/rs/zerolog"
	networkingV1 "istio.io/api/networking/v1"
	securityV1 "istio.io/api/security/v1"
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
	"time"
)

const (
	version = "networkingV1.25.1"
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
		collections.DestinationRule,
		collections.PeerAuthentication,
		collections.Service,
		collections.ServiceEntry,
	))

	configs := []config.Config{
		{
			Meta: config.Meta{
				GroupVersionKind: gvk.PeerAuthentication,
				Name:             "default",
				Namespace:        "istio-system",
			},
			Spec: &securityV1.PeerAuthentication{
				Mtls: &securityV1.PeerAuthentication_MutualTLS{
					Mode: securityV1.PeerAuthentication_MutualTLS_STRICT,
				},
			},
		},
		{
			Meta: config.Meta{
				GroupVersionKind: gvk.DestinationRule,
				Name:             "default",
				Namespace:        "istio-system",
			},
			Spec: &networkingV1.DestinationRule{
				Host: "*.svc.cluster.local",
				TrafficPolicy: &networkingV1.TrafficPolicy{
					Tls: &networkingV1.ClientTLSSettings{
						Mode: networkingV1.ClientTLSSettings_ISTIO_MUTUAL,
					},
				},
			},
		},
		{
			Meta: config.Meta{
				GroupVersionKind: gvk.ServiceEntry,
				Namespace:        "default",
				Name:             "example-service",
			},
			Spec: &networkingV1.ServiceEntry{
				Hosts: []string{"example.com"},
			},
		},
	}

	for _, c := range configs {
		if _, err := store.Create(c); err != nil {
			log.Fatal().Err(err).Msg("Failed to create config")
		}
	}

	env := cfgModel.NewEnvironment()

	meshConfig := mesh.DefaultMeshConfig()

	env.Watcher = meshwatcher.NewTestWatcher(meshConfig)

	services := []*cfgModel.Service{
		{
			Attributes: cfgModel.ServiceAttributes{
				Name:      "svc-a",
				Namespace: "ns-a",
			},
			DefaultAddress: "10.0.0.2",
			Hostname:       "svc-a.ns-a.svc.cluster.local",
			Ports: []*cfgModel.Port{
				{
					Name: "http",
					Port: 8080,
					//Protocol: "HTTP",
				},
			},
			Resolution: cfgModel.ClientSideLB,
		},
		{
			Attributes: cfgModel.ServiceAttributes{
				Name:      "svc-b",
				Namespace: "ns-b",
			},
			DefaultAddress: "10.0.0.3",
			Hostname:       "svc-b.ns-b.svc.cluster.local",
			Ports: []*cfgModel.Port{
				{
					Name: "http",
					Port: 8080,
					//Protocol: "HTTP",
				},
			},
			Resolution: cfgModel.ClientSideLB,
		},
	}

	env.Init()

	env.ServiceDiscovery = model.NewLocalServiceDiscovery(services)
	env.ConfigStore = store

	// Initialize a PushContext with the config store.
	pushContext := cfgModel.NewPushContext()
	pushContext.Mesh = meshConfig

	if err := pushContext.InitContext(env, nil, nil); err != nil {
		log.Fatal().Err(err).Msg("Failed to init push context")
	}

	// Create a dummy Proxy. In a real scenario, this would reflect your proxy’s metadata.
	proxy := &cfgModel.Proxy{
		IPAddresses: []string{"10.0.0.1"},
		Metadata: &cfgModel.NodeMetadata{
			Namespace:    "default",
			IstioVersion: "1.22.0",
		},
		Type:            cfgModel.SidecarProxy,
		DNSDomain:       "default.svc.cluster.local",
		ConfigNamespace: "default",
		IstioVersion:    cfgModel.ParseIstioVersion("1.22.0"),
		SidecarScope:    cfgModel.DefaultSidecarScopeForNamespace(pushContext, "default"),
	}

	proxy.DiscoverIPMode()

	configGen := &core.ConfigGeneratorImpl{}

	listeners := configGen.BuildListeners(proxy, pushContext)

	pushReq := &cfgModel.PushRequest{
		Push:  pushContext,
		Full:  true,
		Start: time.Now(),
	}

	routes, _ := configGen.BuildHTTPRoutes(proxy, pushReq, core.ExtractRoutesFromListeners(listeners))

	if err := protoio.WriteProtoJSONList(ctx, *output, listeners, routes); err != nil {
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
