package main

import (
	"reflect"
	"testing"
)

func TestBuildArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "Simple arguments",
			args: []string{"--config", "--verbose"},
			want: []string{"--config", "--verbose"},
		},
		{
			name: "Arguments with spaces",
			args: []string{"--config value", "--verbose"},
			want: []string{"--config", "value", "--verbose"},
		},
		{
			name: "Combined arguments",
			args: []string{"--config=value", "--verbose"},
			want: []string{"--config=value", "--verbose"},
		},
		{
			name: "Single string with commas",
			args: []string{"--config, --verbose, --handle-this=also"},
			want: []string{"--config", "--verbose", "--handle-this=also"},
		},
		{
			name: "Empty input",
			args: []string{""},
			want: []string{},
		},
		{
			name: "Only spaces",
			args: []string{"   "},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildArgs(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
