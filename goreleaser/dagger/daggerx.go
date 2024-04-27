package main

import (
	"fmt"
	"strings"
)

// addCMDsToContainer injects commands into a container.
// This function is useful when you want to inject commands into a container
// before running them.
func addCMDsToContainer(cmd []string, args []string, ctr *Container) *Container {
	return ctr.
		WithFocus().
		WithExec(mergeSlices(cmd, args))
}

func replaceEntryPointForShell(ctr *Container) *Container {
	return ctr.
		WithoutEntrypoint().
		WithEntrypoint(nil)
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

// addEnvVarsToContainer adds environment variables to a container.
func addEnvVarsToContainer(envVars map[string]string, ctr *Container) *Container {
	for key, value := range envVars {
		ctr = ctr.WithEnvVariable(key, value)
	}
	return ctr
}

// toEnvVars convert a string with a form SOMETHING=SOMETHING into a valid map
func toEnvVars(envVars []string) (map[string]string, error) {
	envVarsMap := make(map[string]string)
	for _, envVar := range envVars {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid environment variable format: %s", envVar)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		envVarsMap[key] = value
	}
	return envVarsMap, nil
}
