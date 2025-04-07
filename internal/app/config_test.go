package app

import (
	"flag"
	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog"
	"os"
	"testing"
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
				LogLevel: zerolog.InfoLevel,
				Output:   "out/envoy_config.json",
			},
		},
		{
			name: "output flag set",
			args: []string{"cmd", "--output=foo"},
			expected: &Config{
				LogLevel: zerolog.InfoLevel,
				Output:   "foo",
			},
		},
		{
			name: "debug flag set",
			args: []string{"cmd", "--debug=true"},
			expected: &Config{
				LogLevel: zerolog.DebugLevel,
				Output:   "out/envoy_config.json",
			},
		},
		{
			name: "all flags set",
			args: []string{"cmd", "--output=foo", "--debug=true"},
			expected: &Config{
				LogLevel: zerolog.DebugLevel,
				Output:   "foo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			os.Args = tt.args

			got := ParseFlagsToConfig()
			want := tt.expected

			if diff := cmp.Diff(want, got); diff != "" {
				//t.Errorf("ParseFlagsToConfig() mismatch = %+v, want %+v", result, tt.expected)
				t.Errorf("ParseFlagsToConfig() mismatch (-want +got):\n%s", diff)
			}

			//if !reflect.DeepEqual(result, tt.expected) {
			//
			//}
		})
	}
}
