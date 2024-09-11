package main

import (
	"context"

	"github.com/Excoriate/daggerverse/registryauth/internal/dagger"
)

const (
	validKeyPattern = `^[a-zA-Z0-9_]+$`
)

// OpenTerminal returns a terminal
//
// It returns a terminal for the container.
// Arguments:
// - None.
// Returns:
// - *Terminal: The terminal for the container.
func (m *Registryauth) OpenTerminal() *dagger.Container {
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
func (m *Registryauth) RunShell(cmd string) (string, error) {
	out, err := m.Ctr.WithoutEntrypoint().WithExec([]string{"sh", "-c", cmd}).Stdout(context.Background())
	if err != nil {
		return "", WrapErrorf(err, "failed to run shell command: %s", cmd)
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
// envVars, err := RegistryauthInstance.PrintEnvVars()
//
//	if err != nil {
//	    log.Fatalf("Error retrieving environment variables: %v", err)
//	}
//
// fmt.Println(envVars).
func (m *Registryauth) PrintEnvVars() (string, error) {
	out, err := m.
		Ctr.
		WithExec([]string{"printenv"}).
		Stdout(context.Background())

	if err != nil {
		return "", WrapError(err, "failed to get env vars")
	}

	return out, nil
}
