package slices

import "testing"

func TestSlicesAreEquivalent(t *testing.T) {
	tests := []struct {
		name   string
		sliceA []string
		sliceB []string
		want   bool
	}{
		{
			name:   "both slices nil",
			sliceA: nil,
			sliceB: nil,
			want:   true,
		},
		{
			name:   "first nil, second empty",
			sliceA: nil,
			sliceB: []string{},
			want:   true,
		},
		{
			name:   "first empty, second nil",
			sliceA: []string{},
			sliceB: nil,
			want:   true,
		},
		{
			name:   "both slices empty",
			sliceA: []string{},
			sliceB: []string{},
			want:   true,
		},
		{
			name:   "identical single-element slices",
			sliceA: []string{"element"},
			sliceB: []string{"element"},
			want:   true,
		},
		{
			name:   "identical multi-element slices",
			sliceA: []string{"hello", "world"},
			sliceB: []string{"hello", "world"},
			want:   true,
		},
		{
			name:   "non-identical single-element slices",
			sliceA: []string{"hello"},
			sliceB: []string{"world"},
			want:   false,
		},
		{
			name:   "non-identical multi-element slices",
			sliceA: []string{"hello", "there"},
			sliceB: []string{"hello", "world"},
			want:   false,
		},
		{
			name:   "different lengths, one empty",
			sliceA: []string{"hello"},
			sliceB: []string{},
			want:   false,
		},
		{
			name:   "different lengths, non-empty",
			sliceA: []string{"hello", "world"},
			sliceB: []string{"hello"},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SlicesAreEquivalent(tt.sliceA, tt.sliceB); got != tt.want {
				t.Errorf("SlicesAreEquivalent() = %v, want %v", got, tt.want)
			}
		})
	}
}
