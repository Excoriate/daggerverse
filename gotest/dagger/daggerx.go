package main

import (
	"fmt"
	"strings"
)

type DaggerEnvVars struct {
	Name   string
	Value  string
	Expand bool
}

// toEnvVars convert a string with a form SOMETHING=SOMETHING into a valid map
func toEnvVarsFromStr(envVars string) (map[string]string, error) {
	envVarsMap := make(map[string]string)
	parts := strings.Split(envVars, ",")
	for _, envVar := range parts {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid environment variable format: %s", envVar)
		}
		envVarsMap[parts[0]] = parts[1]
	}
	return envVarsMap, nil
}

// addEnvVarsToContainer adds environment variables to a container.
func toEnvVarsDaggerFromMap(envVarsMap map[string]string) []DaggerEnvVars {
	var envVars []DaggerEnvVars
	for key, value := range envVarsMap {
		envVars = append(envVars, DaggerEnvVars{
			Name:  key,
			Value: value,
		})
	}
	return envVars
}
