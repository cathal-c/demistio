package app

import (
	"context"
	"fmt"
	"github.com/cathal-c/demistio/pkg/catalog"
	"github.com/cathal-c/demistio/pkg/model"
	"github.com/cathal-c/demistio/pkg/protoio"
	"github.com/rs/zerolog"
	v1 "istio.io/api/networking/v1"
	v2 "istio.io/api/security/v1"
	"istio.io/istio/pilot/pkg/config/memory"
	cfgModel "istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking/core"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/mesh"
	"istio.io/istio/pkg/config/mesh/meshwatcher"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/collections"
	"time"
)

func Generate(ctx context.Context, cfg *Config) error {
	//log := zerolog.Ctx(ctx)

	// Create an in-memory config store, registering required resource types
	store := memory.Make(collection.SchemasFor(
		collections.DestinationRule,
		collections.PeerAuthentication,
		collections.Service,
		collections.ServiceEntry,
	))

	configs := []config.Config{
		catalog.NewPeerAuthenticationBuilder(catalog.DefaultResourceName, catalog.DefaultIstioNamespace).WithMutualTlsMode(v2.PeerAuthentication_MutualTLS_STRICT).Build(),
		catalog.NewDestinationRuleBuilder(catalog.DefaultResourceName, catalog.DefaultIstioNamespace).WithHost("*.svc.cluster.local").WithTlsMode(v1.ClientTLSSettings_ISTIO_MUTUAL).Build(),
		catalog.NewServiceEntryBuilder("google", catalog.DefaultIstioNamespace).WithHosts("google.com").Build(),
	}

	if err := addConfigsToConfigStore(ctx, store, configs); err != nil {
		return fmt.Errorf("adding configs to config store: %w", err)
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
				},
			},
			Resolution: cfgModel.ClientSideLB,
		},
	}

	env.Init()

	env.ServiceDiscovery = model.NewLocalServiceDiscovery(services)
	env.ConfigStore = store

	pushContext := cfgModel.NewPushContext()
	pushContext.Mesh = meshConfig

	if err := pushContext.InitContext(env, nil, nil); err != nil {
		return fmt.Errorf("initializing push context: %v", err)
	}

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

	if err := protoio.WriteProtoJSONList(ctx, cfg.Output, listeners, routes); err != nil {
		return fmt.Errorf("writing proto json list: %v", err)
	}

	return nil
}

func addConfigsToConfigStore(ctx context.Context, store cfgModel.ConfigStore, configs []config.Config) error {
	log := zerolog.Ctx(ctx)

	for _, c := range configs {
		if _, err := store.Create(c); err != nil {
			return fmt.Errorf("%s: %w", c.Name, err)
		}
	}

	log.Debug().Msgf("added %d configs to store", len(configs))

	return nil
}
