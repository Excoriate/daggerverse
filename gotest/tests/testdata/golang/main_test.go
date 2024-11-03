package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestFibonacci(t *testing.T) {
	tests := []struct {
		n        int
		expected string
	}{
		{n: -1, expected: "Please provide a positive integer.\n"},
		{n: 0, expected: "Please provide a positive integer.\n"},
		{n: 1, expected: "0 \n"},
		{n: 2, expected: "0 1 \n"},
		{n: 5, expected: "0 1 1 2 3 \n"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("n=%d", tt.n), func(t *testing.T) {
			var buf bytes.Buffer
			fmt.Print(buf.String())

			fibSeries := fibonacci(tt.n)
			fmt.Println(fibSeries)
		})
	}
}
