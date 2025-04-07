package app

import (
	"context"
	"flag"
	"os"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
)

func TestParseFlagsToConfig(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected *Config
	}{
		{
			name: "no flags",
			args: []string{"cmd"},
			expected: &Config{
				Output: "",
			},
		},
		{
			name: "output flag set",
			args: []string{"cmd", "--output=/tmp/envoy.json"},
			expected: &Config{
				Output: "/tmp/envoy.json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag.CommandLine and os.Args
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			os.Args = tt.args

			// Set up logger context
			logger := zerolog.Nop()
			ctx := logger.WithContext(context.Background())

			result := ParseFlagsToConfig(ctx)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseFlagsToConfig() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}
