package protoio

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
	"os"
)

func WriteProtoJSONList(ctx context.Context, filename string, list []proto.Message) error {
	log := zerolog.Ctx(ctx)

	marshaler := protojson.MarshalOptions{
		Multiline:       true,
		Indent:          "  ",
		EmitUnpopulated: false,
		UseProtoNames:   true,
		Resolver:        protoregistry.GlobalTypes,
	}

	log.Info().Msg("Cathal")

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", filename, err)
	}
	defer f.Close()

	if _, err := f.WriteString("[\n"); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	for i, msg := range list {
		out, err := marshaler.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal proto message: %w", err)
		}
		if _, err := f.Write(out); err != nil {
			return fmt.Errorf("failed to write JSON to file: %w", err)
		}
		if i < len(list)-1 {
			if _, err := f.WriteString(",\n"); err != nil {
				return fmt.Errorf("failed to write comma: %w", err)
			}
		}
	}

	if _, err := f.WriteString("\n]\n"); err != nil {
		return fmt.Errorf("failed to finalize JSON array: %w", err)
	}

	return nil
}
