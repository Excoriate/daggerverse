package args

import (
	"github.com/Excoriate/daggerverse/utils/pkg/slices"
	"testing"
)

func TestParseArgsFromStrToSlice(t *testing.T) {
	tests := []struct {
		name string
		args string
		want []string
	}{
		{
			name: "valid input with multiple arguments",
			args: "--from-module=asdsada, --upgrade=false",
			want: []string{"--from-module=asdsada", "--upgrade=false"},
		},
		{
			name: "valid input with single argument",
			args: "--single-arg=true",
			want: []string{"--single-arg=true"},
		},
		{
			name: "invalid input with only commas",
			args: "   ,   ,   ",
			want: nil,
		},
		{
			name: "edge case: empty string",
			args: "",
			want: nil,
		},
		{
			name: "edge case: string with only spaces",
			args: "     ",
			want: nil,
		},
		{
			name: "mixed valid and invalid input",
			args: "--valid=true, , , --also-valid=false",
			want: []string{"--valid=true", "--also-valid=false"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseArgsFromStrToSlice(tt.args); !slices.SlicesAreEquivalent(got, tt.want) {
				t.Errorf("ParseArgsFromStrToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
