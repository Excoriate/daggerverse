package main

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/excoriate/daggerverse/module-example/internal/dagger"
)

const (
	errEmptyEnvVarKey          = "environment variable key cannot be empty"
	errInvalidEnvVarKeyPattern = "failed to validate key"
	errEnvVarPatternMismatch   = "the key %s is invalid, does not match the pattern %s"
	errInspectingEnvVar        = "failed to inspect env var %s"
	errEnvVarNotFound          = "environment variable %s not found"
	validKeyPattern            = `^[a-zA-Z0-9_]+$`
)

var (
	errEmptyEnvVarKeyError        = errors.New(errEmptyEnvVarKey)
	errEnvVarPatternMismatchError = errors.New(errEnvVarPatternMismatch)
	errEnvVarNotFoundError        = errors.New(errEnvVarNotFound)
)

// OpenTerminal returns a terminal
//
// It returns a terminal for the container.
// Arguments:
// - None.
// Returns:
// - *Terminal: The terminal for the container.
func (m *ModuleExample) OpenTerminal() *dagger.Container {
	return m.Ctr.Terminal()
}

// RunShell runs a shell command in the container.
//
// It runs a shell command in the container and returns the output.
// Arguments:
// - cmd: The command to run in the container.
// Returns:
// - string: The output of the command.
// - error: An error if the command fails.
func (m *ModuleExample) RunShell(cmd string) (string, error) {
	out, err := m.Ctr.WithoutEntrypoint().WithExec([]string{"sh", "-c", cmd}).Stdout(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to run shell command: %w", err)
	}

	return out, nil
}

// PrintEnvVars retrieves and prints the environment variables of the container.
//
// It executes the `printenv` command inside the container to get a list of all
// environment variables and their respective values.
//
// Arguments:
// - None.
//
// Returns:
//   - string: A string containing all environment variables in the format
//     "KEY=VALUE", separated by newlines.
//   - error: An error if the command fails, wrapped with additional context.
//
// Usage example:
// ```go
// envVars, err := ModuleExampleInstance.PrintEnvVars()
//
//	if err != nil {
//	    log.Fatalf("Error retrieving environment variables: %v", err)
//	}
//
// fmt.Println(envVars).
func (m *ModuleExample) PrintEnvVars() (string, error) {
	out, err := m.
		Ctr.
		WithExec([]string{"printenv"}).
		Stdout(context.Background())

	if err != nil {
		return "", fmt.Errorf("failed to get env vars: %w", err)
	}

	return out, nil
}

// InspectEnvVar inspects the value of an environment variable by its key
// Arguments:
// - key: The environment variable key to inspect.
// Returns:
// - string: The value of the environment variable.
// - error: An error if the key is invalid or the environment variable is not found.
func (m *ModuleExample) InspectEnvVar(key string) (string, error) {
	// Validate if the key is empty or contains invalid characters
	if key == "" {
		return "", fmt.Errorf("%w", errEmptyEnvVarKeyError)
	}

	matched, err := regexp.MatchString(validKeyPattern, key)
	if err != nil {
		return "", fmt.Errorf("%s: %w", errInvalidEnvVarKeyPattern, err)
	}

	if !matched {
		return "", fmt.Errorf("%w", errEnvVarPatternMismatchError)
	}

	// Execute the printenv command to get the environment variable's value
	out, envVarErr := m.
		Ctr.
		WithExec([]string{"printenv", key}).
		Stdout(context.Background())

	if envVarErr != nil {
		return "", fmt.Errorf("%s: %w", fmt.Sprintf(errInspectingEnvVar, key), envVarErr)
	}

	// Check if the output is empty, which indicates the environment variable was not found
	if out == "" {
		return "", fmt.Errorf("%w", errEnvVarNotFoundError)
	}

	return out, nil
}
