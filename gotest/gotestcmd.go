package main

import (
	"github.com/Excoriate/daggerverse/gotest/internal/dagger"
)

const (
	cmdEntrypoint = "go"
	cmdTest       = "test"
)

// getBaseCmd returns the base command for running Go tests.
//
// This function constructs and returns a slice of strings that represents
// the base command to execute Go tests, which includes the command
// entry point and the test subcommand.
func getBaseCmd() []string {
	return []string{cmdEntrypoint, cmdTest}
}

// GoTestCmd represents a command to run Go tests in a Dagger container.
//
// This struct holds the container in which the command will be executed
// and the command arguments to be passed to the Go test command.
type GoTestCmd struct {
	// BaseCmd is the base command to run.
	// +private
	BaseCmd []string
	// Packages are the packages to test.
	// +private
	Packages []string
	// EnvironmentVariables is the environment variables to set.
	// +private
	EnvironmentVariables []string
	// Secrets are the secrets to set.
	// +private
	Secrets []*dagger.Secret
	// DagModule is the DAG module to use.
	// +private
	DagModule *Gotest
	// BuildOpts are the build options to use.
	// +private
	BuildOpts *GoBuildOptions
	// TestOpts are the test options to use.
	// +private
	TestOpts *GoTestOptions
}

// newGoTestCmd creates a new GoTestCmd instance for running Go tests.
//
// This method initializes a GoTestCmd with the specified packages, environment
// variables, and secrets. It constructs the command to be executed in a Dagger
// container, setting up the necessary parameters for the Go test command.
//
// Parameters:
//   - packages: A slice of strings representing the Go packages to be tested.
//   - environmentVariables: A slice of strings representing the environment
//     variables to set for the test execution.
//   - secrets: A slice of pointers to dagger.Secret representing the secrets
//     to be used during the test execution.
//
// Returns:
// - *GoTestCmd: A pointer to the newly created GoTestCmd instance.
func (m *Gotest) newGoTestCmd(
	packages []string,
	environmentVariables []string,
	secrets []*dagger.Secret,
) (*GoTestCmd, error) {
	if len(packages) == 0 {
		return nil, Errorf("failed to create new GoTestCmd: no packages to test")
	}

	return &GoTestCmd{
		BaseCmd:              getBaseCmd(),
		Packages:             packages,
		EnvironmentVariables: environmentVariables,
		Secrets:              secrets,
		DagModule:            m,
		BuildOpts:            NewGoBuildOptions(),
		TestOpts:             NewGoTestOptions(),
	}, nil
}
