package app

import (
	"flag"
	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel zerolog.Level
	Output   string
}

func ParseFlagsToConfig() *Config {
	// inputPtr := flag.String("input", "", "Path to YAML file containing Istio configs")
	debug := flag.Bool("debug", false, "Enable debug mode")
	output := flag.String("output", "out/envoy_config.json", "Output file for generated Envoy config (JSON)")
	flag.Parse()

	cfg := &Config{
		Output: *output,
	}

	if *debug {
		cfg.LogLevel = zerolog.DebugLevel
	} else {
		cfg.LogLevel = zerolog.InfoLevel
	}

	return cfg
}
