package main

import (
	"fmt"
	"strings"
)

// addCMDsToContainer injects commands into a container.
// This function is useful when you want to inject commands into a container
// before running them.
func addCMDsToContainer(cmd, args []string, ctr *Container) *Container {
	fullCommand := append(cmd, args...)
	return ctr.
		WithFocus().
		WithExec(fullCommand)
}

func replaceEntryPointForShell(ctr *Container) *Container {
	return ctr.
		WithoutEntrypoint().
		WithEntrypoint(nil)
}

// buildArgs constructs an argument slice from a variadic string input,
// splitting by commas or spaces, and ignoring entirely empty or whitespace-only strings.
func buildArgs(args ...string) []string {
	var merged []string
	for _, arg := range args {
		if arg = strings.TrimSpace(arg); arg != "" {
			// Splits the string into substrings by commas first, then trims and splits by spaces
			for _, part := range strings.Split(arg, ",") {
				part = strings.TrimSpace(part)
				if part != "" {
					parts := strings.Fields(part) // Further split by any spaces remaining
					merged = append(merged, parts...)
				}
			}
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

type EnvVarDagger struct {
	Name   string
	Value  string
	Expand bool
}

func toEnvVarsDagger(envVarSlice []string) ([]EnvVarDagger, error) {
	var envVars []EnvVarDagger
	for _, envVar := range envVarSlice {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid environment variable format: %s", envVar)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		envVars = append(envVars, EnvVarDagger{
			Name:   key,
			Value:  value,
			Expand: true,
		})
	}

	return envVars, nil
}
