package proto

import (
	"context"
	"fmt"
	"github.com/cathal-c/demistio/pkg/encoding"
	"github.com/cathal-c/demistio/pkg/envoy"
	"github.com/cathal-c/demistio/pkg/file"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
)

// WriteConfigDumpToFile creates a config dump and writes it to a file as JSON.
func WriteConfigDumpToFile(ctx context.Context, path string, listeners []*listenerv3.Listener, routes []*discoveryv3.Resource) error {
	cfgDump, err := envoy.BuildFullConfigDump(listeners, routes)
	if err != nil {
		return fmt.Errorf("build config dump: %w", err)
	}

	data, err := encoding.MarshalProtoMessageToJSON(cfgDump)
	if err != nil {
		return fmt.Errorf("marshal config dump: %w", err)
	}

	if err := file.Write(ctx, path, data); err != nil {
		return fmt.Errorf("write to file: %w", err)
	}

	return nil
}
