package app

import (
	"context"
	"flag"
	"github.com/rs/zerolog"
	"os"
)

func ParseFlagsToConfig(ctx context.Context) *Config {
	log := zerolog.Ctx(ctx)

	// inputPtr := flag.String("input", "", "Path to YAML file containing Istio configs")
	output := flag.String("output", "", "Output file for generated Envoy config (JSON)")
	flag.Parse()

	// os.Args[0] is always the binary
	for _, arg := range os.Args[1:] {
		log.Info().Msg(arg)
	}

	return &Config{
		Output: *output,
	}
}
