package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	// Istio internal packages. These paths and APIs may change with Istio versions.
	v1 "istio.io/api/networking/v1"
	"istio.io/istio/pilot/pkg/config/memory"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking/core"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/schema/collections"
	"istio.io/istio/pkg/config/schema/gvk"
)

func main() {
	inputPtr := flag.String("input", "", "Path to YAML file containing Istio configs")
	outputPtr := flag.String("output", "", "Output file for generated Envoy config (JSON)")
	flag.Parse()

	if *inputPtr == "" || *outputPtr == "" {
		log.Fatalf("Usage: go run main.go -input=input.yaml -output=envoy_config.json")
	}

	// // Read the YAML file
	// data, err := os.ReadFile(*inputPtr)
	// if err != nil {
	// 	log.Fatalf("Error reading file: %v", err)
	// }

	// // Load the Istio config objects from YAML
	// configs, err := loadConfigs(data)
	// if err != nil {
	// 	log.Fatalf("Error loading configs: %v", err)
	// }

	// Create an in-memory config store.
	// Note: model.IstioConfigTypes is a map defining supported config types.
	store := memory.Make(collections.All)

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
		log.Fatalf("Error creating Store: %s", err.Error())
	}

	// for _, cfg := range configs {
	// 	// Create the config in the store.
	// 	// In a full implementation, you would handle errors and possibly update existing entries.
	// 	if _, err := store.Create(cfg); err != nil {
	// 		log.Fatalf("Warning: error adding config: %v", err)
	// 	}
	// }

	// Initialize a PushContext with the config store.
	push := model.NewPushContext()
	env := &model.Environment{
		// PushContext: push,
		ConfigStore: store,
	}

	push.InitContext(env, nil, nil)

	// Create a dummy Proxy. In a real scenario, this would reflect your proxy’s metadata.
	proxy := &model.Proxy{
		IPAddresses: []string{"127.0.0.1"},
		// Populate additional fields as required.
	}

	configGen := &core.ConfigGeneratorImpl{}

	// Use Istio's v1alpha3 conversion logic to build Envoy listeners.
	listeners := configGen.BuildListeners(proxy, push)

	// Serialize the generated listeners to JSON.
	jsonData, err := json.MarshalIndent(listeners, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}

	// Write the JSON output.
	if err := os.WriteFile(*outputPtr, jsonData, 0644); err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Printf("Generated Envoy configuration written to %s\n", *outputPtr)
}

// func loadConfigs(data []byte) ([]config.Config, error) {
// 	var configs []config.Config
// 	docs := strings.Split(string(data), "\n---\n")

// 	for _, doc := range docs {
// 		if strings.TrimSpace(doc) == "" {
// 			continue
// 		}

// 		var raw interface{}
// 		if err := yaml.Unmarshal([]byte(doc), &raw); err != nil {
// 			return nil, fmt.Errorf("error unmarshalling YAML: %v", err)
// 		}
// 		raw = convertMapKeysToString(raw)
// 		typedRaw, ok := raw.(map[string]interface{})
// 		if !ok {
// 			return nil, fmt.Errorf("unexpected YAML structure")
// 		}

// 		kind, _ := typedRaw["kind"].(string)
// 		apiVersion, _ := typedRaw["apiVersion"].(string)

// 		// Split apiVersion into group/version
// 		parts := strings.Split(apiVersion, "/")
// 		var group, version string
// 		if len(parts) == 2 {
// 			group, version = parts[0], parts[1]
// 		} else {
// 			group = ""
// 			version = parts[0]
// 		}

// 		gvk := config.GroupVersionKind{
// 			Group:   group,
// 			Version: version,
// 			Kind:    kind,
// 		}

// 		// Get the schema (resource definition) for this GVK
// 		schema, found := collections.All.FindByGroupVersionKind(gvk)
// 		if !found {
// 			return nil, fmt.Errorf("unsupported kind: %s/%s", apiVersion, kind)
// 		}

// 		// Decode the spec into a concrete instance
// 		specJSON, err := json.Marshal(typedRaw["spec"])
// 		if err != nil {
// 			return nil, fmt.Errorf("error marshalling spec: %v", err)
// 		}

// 		spec, err := schema.NewInstance()
// 		if err != nil {
// 			return nil, fmt.Errorf("error creating spec instance: %v", err)
// 		}

// 		if err := json.Unmarshal(specJSON, spec); err != nil {
// 			return nil, fmt.Errorf("error unmarshalling spec into typed object: %v", err)
// 		}

// 		// Pull metadata
// 		meta := typedRaw["metadata"].(map[string]interface{})
// 		name := meta["name"].(string)
// 		namespace := "default"
// 		if ns, ok := meta["namespace"].(string); ok {
// 			namespace = ns
// 		}

// 		cfg := config.Config{
// 			Meta: config.Meta{
// 				GroupVersionKind: gvk,
// 				Name:             name,
// 				Namespace:        namespace,
// 			},
// 			Spec: spec,
// 		}

// 		configs = append(configs, cfg)
// 	}
// 	return configs, nil
// }

// func convertMapKeysToString(i interface{}) interface{} {
// 	switch x := i.(type) {
// 	case map[interface{}]interface{}:
// 		m2 := map[string]interface{}{}
// 		for k, v := range x {
// 			m2[fmt.Sprint(k)] = convertMapKeysToString(v)
// 		}
// 		return m2
// 	case []interface{}:
// 		for i, v := range x {
// 			x[i] = convertMapKeysToString(v)
// 		}
// 	}
// 	return i
// }
