package httpsrv

import (
	"reflect"
	"testing"
)

func TestConfig_buildConfig(t *testing.T) {
	type test struct {
		name     string
		opts     []Options
		expected *options
	}
	tests := []test{
		{
			name:     "no options",
			expected: &options{path: "."},
		},
		{
			name: "with port",
			opts: []Options{
				WithPort(80),
			},
			expected: &options{port: 80, path: "."},
		},
		{
			name: "with path",
			opts: []Options{
				WithPath("/tmp"),
			},
			expected: &options{port: 0, path: "/tmp"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildConfig(tt.opts...)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("[%s] buildConfig() = %v, expected %v", tt.name, got, tt.expected)
			}
		})
	}
}
