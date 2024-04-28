package main

import (
	"strings"
)

// addCMDsToContainer injects commands into a container.
// This function is useful when you want to inject commands into a container
// before running them.
func addCMDsToContainer(cmd, args []string, ctr *Container) *Container {
	return ctr.
		WithFocus().
		WithExec(mergeSlices(cmd, args))
}

// buildArgs constructs an argument slice from a variadic string input,
// ignoring entirely empty or whitespace-only strings.
func buildArgs(args ...string) []string {
	var merged []string
	for _, arg := range args {
		if arg = strings.TrimSpace(arg); arg != "" {
			parts := strings.Fields(arg) // Splits the string into substrings removing any space characters, including newlines.
			merged = append(merged, parts...)
		}
	}
	return merged
}

// mergeSlices merges slices of strings.
func mergeSlices(slices ...[]string) []string {
	var merged []string
	for _, slice := range slices {
		merged = append(merged, slice...)
	}
	return merged
}
