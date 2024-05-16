package main

import (
	"errors"
	"fmt"
	"strings"
)

type DaggerEnvVars struct {
	Name   string
	Value  string
	Expand bool
}

// toEnvVarsFromStr converts a comma-separated string of key=value pairs into a map.
// It ensures all entries are valid and handles empty strings gracefully.
func toEnvVarsFromStr(envVars string) (map[string]string, error) {
	if envVars == "" {
		return nil, errors.New("input string is empty")
	}

	envVarsMap := make(map[string]string)
	parts := strings.Split(envVars, ",")
	for _, envVar := range parts {
		if envVar == "" {
			continue
		}

		pair := strings.SplitN(envVar, "=", 2)
		if len(pair) != 2 {
			return nil, fmt.Errorf("invalid environment variable format: %s", envVar)
		}

		key, value := strings.TrimSpace(pair[0]), strings.TrimSpace(pair[1])
		if key == "" {
			return nil, fmt.Errorf("empty key in environment variable: %s", envVar)
		}

		envVarsMap[key] = value
	}
	return envVarsMap, nil
}

// toEnvVarsDaggerFromMap converts a map of environment variables into a slice of DaggerEnvVars.
// It ensures all entries are valid.
func toEnvVarsDaggerFromMap(envVarsMap map[string]string) ([]DaggerEnvVars, error) {
	if len(envVarsMap) == 0 {
		return nil, errors.New("input map is empty")
	}

	var envVars []DaggerEnvVars
	for key, value := range envVarsMap {
		if key == "" {
			return nil, errors.New("found empty key in map")
		}
		envVars = append(envVars, DaggerEnvVars{
			Name:  key,
			Value: value,
		})
	}
	return envVars, nil
}

// toEnvVarsDaggerFromSlice converts a slice of key=value strings into a slice of DaggerEnvVars.
// It validates each entry and skips invalid entries.
func toEnvVarsDaggerFromSlice(envVarsSlice []string) ([]DaggerEnvVars, error) {
	if len(envVarsSlice) == 0 {
		return nil, errors.New("input slice is empty")
	}

	var envVars []DaggerEnvVars
	for _, envVar := range envVarsSlice {
		if envVar == "" {
			continue
		}

		pair := strings.SplitN(envVar, "=", 2)
		if len(pair) != 2 {
			return nil, fmt.Errorf("invalid environment variable format: %s", envVar)
		}

		key, value := strings.TrimSpace(pair[0]), strings.TrimSpace(pair[1])
		if key == "" {
			return nil, fmt.Errorf("empty key in environment variable: %s", envVar)
		}

		envVars = append(envVars, DaggerEnvVars{
			Name:  key,
			Value: value,
		})
	}
	return envVars, nil
}
