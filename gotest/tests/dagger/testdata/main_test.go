package main

import (
	"bytes"
	"fmt"
	"testing"
)

// Helper function to capture the output of PrintSum.
func captureOutput(f func()) string {
	var buf bytes.Buffer
	old := stdOut
	stdOut = &buf
	defer func() { stdOut = old }()
	f()
	return buf.String()
}

// TestPrintSum tests the PrintSum function.
func TestPrintSum(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected string
	}{
		{"Positive numbers", 3, 5, "Sum: 8\n"},
		{"Negative numbers", -3, -5, "Sum: -8\n"},
		{"Mixed numbers", -3, 5, "Sum: 2\n"},
		{"Zero and positive", 0, 5, "Sum: 5\n"},
		{"Zero and negative", 0, -5, "Sum: -5\n"},
		{"Both zero", 0, 0, "Sum: 0\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintSum(tt.a, tt.b)
			})
			if output != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, output)
			}
		})
	}
}

// Dummy variable to capture the output.
var stdOut = new(bytes.Buffer)

func PrintSum(a, b int) {
	sum := a + b
	_, _ = fmt.Fprintf(stdOut, "Sum: %d\n", sum)
}
