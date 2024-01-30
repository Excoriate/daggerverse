package main

import (
	"fmt"
	"os"
	"strings"
)

func getContainerImage(image, version string) string {
	if version == "" {
		version = "latest"
	}

	return fmt.Sprintf("%s:%s", image, version)
}

// NewOptional creates a new Optional of any type.
func toDaggerOptional[T any](value T) Optional[T] {
	return Optional[T]{value: value, isSet: true}
}

func scanEnvVarsFromHost() map[string]string {
	envVars := map[string]string{}

	for _, envVar := range os.Environ() {
		parts := strings.Split(envVar, "=")
		envVars[parts[0]] = parts[1]
	}

	return envVars
}

func getEnvVarsFromHostWithPrefix(prefix string) map[string]string {
	envVars := scanEnvVarsFromHost()
	envVarsWithPrefix := map[string]string{}

	for key, value := range envVars {
		if strings.HasPrefix(key, prefix) {
			envVarsWithPrefix[key] = value
		}
	}

	return envVarsWithPrefix
}

func getTFVARsFromHost() map[string]string {
	return getEnvVarsFromHostWithPrefix("TF_VAR_")
}

func getAWSVarsFromHost() map[string]string {
	return getEnvVarsFromHostWithPrefix("AWS_")
}
