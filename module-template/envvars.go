package main

import (
	"context"
	"regexp"
)

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
// envVars, err := ModuleTemplateInstance.PrintEnvVars()
//
//	if err != nil {
//	    log.Fatalf("Error retrieving environment variables: %v", err)
//	}
//
// fmt.Println(envVars).
func (m *ModuleTemplate) PrintEnvVars() (string, error) {
	out, err := m.
		Ctr.
		WithExec([]string{"printenv"}).
		Stdout(context.Background())

	if err != nil {
		return "", WrapError(err, "failed to get env vars")
	}

	return out, nil
}

// InspectEnvVar inspects the value of an environment variable by its key
// Arguments:
// - key: The environment variable key to inspect.
// Returns:
// - string: The value of the environment variable.
// - error: An error if the key is invalid or the environment variable is not found.
func (m *ModuleTemplate) InspectEnvVar(key string) (string, error) {
	// Validate if the key is empty or contains invalid characters
	if key == "" {
		return "", Errorf("environment variable key cannot be empty")
	}

	matched, err := regexp.MatchString(validKeyPattern, key)
	if err != nil {
		return "", WrapErrorf(err, "failed to inspect environment variable by key: %s", key)
	}

	if !matched {
		return "", Errorf("the key %s is invalid, does not match the pattern %s", key, validKeyPattern)
	}

	// Execute the printenv command to get the environment variable's value
	out, envVarErr := m.
		Ctr.
		WithExec([]string{"printenv", key}).
		Stdout(context.Background())

	if envVarErr != nil {
		return "", WrapErrorf(envVarErr, "failed to inspect the environment variable: %s", key)
	}

	// Check if the output is empty, which indicates the environment variable was not found
	if out == "" {
		return "", Errorf("the environment variable %s was not found", key)
	}

	return out, nil
}
