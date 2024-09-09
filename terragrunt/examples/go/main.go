// A generated module for Go example functions
package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Excoriate/daggerverse/terragrunt/examples/go/internal/dagger"
	"github.com/sourcegraph/conc/pool"
)

// Go is a Dagger module that exemplifies the usage of the Terragrunt module.
//
// This module is used to create and manage containers.
type Go struct {
	TestDir *dagger.Directory
}

// New creates a new Tests instance.
//
// It's the initial constructor for the Tests struct.
func New() *Go {
	e := &Go{}

	e.TestDir = e.getTestDir()

	return e
}

var (
	errEnvVarsEmpty            = errors.New("env vars are empty")
	errEnvVarsDontMatch        = errors.New("expected env vars to be passed, got empty output")
	errNetRCFileNotFound       = errors.New("netrc file not found")
	errExpectedFoldersNotFound = errors.New("expected to have at least one folder, got empty output")
)

// getTestDir returns the test directory.
//
// This helper method retrieves the test directory, which is typically located
// in the same directory as the test file and named "testdata".
//
// Returns:
//   - *dagger.Directory: A Dagger Directory object pointing to the "testdata" directory.
func (m *Go) getTestDir() *dagger.Directory {
	return dag.
		CurrentModule().
		Source().
		Directory("./testdata")
}

// AllRecipes executes all tests.
//
// AllRecipes is a helper method for tests, executing the built-in recipes and
// other specific functionalities of the Terragrunt module.
func (m *Go) AllRecipes(ctx context.Context) error {
	polTests := pool.New().WithErrors().WithContext(ctx)

	// Test different ways to configure the base container.
	polTests.Go(m.BuiltInRecipes)
	// From this point onwards, we're testing the specific functionality of the Terragrunt module.

	if err := polTests.Wait(); err != nil {
		return fmt.Errorf("there are some failed tests: %w", err)
	}

	return nil
}

// BuiltInRecipes demonstrates how to run built-in recipes
//
// This method showcases the use of various built-in recipes provided by the Terragrunt
// module, including creating a container, running an arbitrary command, and creating a .netrc
// file for GitHub authentication.
//
// Parameters:
//   - ctx: The context for controlling the function's timeout and cancellation.
//
// Returns:
//   - An error if any of the internal methods fail, or nil otherwise.
func (m *Go) BuiltInRecipes(ctx context.Context) error {
	// Pass environment variables to the Terragrunt module using Terragrunt_PassedEnvVars
	if err := m.Terragrunt_PassedEnvVars(ctx); err != nil {
		return fmt.Errorf("failed to pass environment variables: %w", err)
	}

	// Create a configured container using Terragrunt_CreateContainer
	if _, err := m.Terragrunt_CreateContainer(ctx); err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Run an arbitrary command in the container using Terragrunt_RunArbitraryCommand
	if _, err := m.Terragrunt_RunArbitraryCommand(ctx); err != nil {
		return fmt.Errorf("failed to run arbitrary command: %w", err)
	}

	// Create a .netrc file for GitHub using Terragrunt_CreateNetRcFileForGithub
	if _, err := m.Terragrunt_CreateNetRcFileForGithub(ctx); err != nil {
		return fmt.Errorf("failed to create netrc file: %w", err)
	}

	return nil
}

// Terragrunt_PassedEnvVars demonstrates how to pass environment variables to the Terragrunt module.
//
// This method configures a Terragrunt module to use specific environment variables from the host.
func (m *Go) Terragrunt_PassedEnvVars(ctx context.Context) error {
	targetModule := dag.Terragrunt(dagger.TerragruntOpts{
		EnvVarsFromHost: []string{"SOMETHING=SOMETHING,SOMETHING=SOMETHING"},
	})

	out, err := targetModule.Ctr().
		WithExec([]string{"printenv"}).
		Stdout(ctx)

	if err != nil {
		return fmt.Errorf("failed when executing printenv: %w", err)
	}

	if out == "" {
		return errEnvVarsEmpty
	}

	if !strings.Contains(out, "SOMETHING") {
		return errEnvVarsDontMatch
	}

	return nil
}

// Terragrunt_OpenTerminal demonstrates how to open an interactive terminal session
// within a Terragrunt module container.
//
// This function showcases the initialization and configuration of a
// Terragrunt module container using various options like enabling Cgo,
// utilizing build cache, and including a GCC compiler.
//
// Parameters:
//   - None
//
// Returns:
//   - *dagger.Container: A configured Dagger container with an open terminal.
//
// Usage:
//
//	This function can be used to interactively debug or inspect the
//	container environment during test execution.
func (m *Go) Terragrunt_OpenTerminal() *dagger.Container {
	// Configure the Terragrunt module container with necessary options
	targetModule := dag.Terragrunt()

	// Retrieve and discard standard output
	_, _ = targetModule.Ctr().
		Stdout(context.Background())

	// Open and return the terminal session in the configured container
	return targetModule.Ctr().
		Terminal()
}

// Terragrunt_CreateNetRcFileForGithub creates and configures a .netrc file for GitHub authentication.
//
// This method exemplifies the creation of a .netrc file with credentials for accessing GitHub,
// and demonstrates how to pass this file as a secret to the Terragrunt module.
//
// Parameters:
//   - ctx: The context for controlling the function's timeout and cancellation.
//
// Returns:
//   - A configured Container with the .netrc file set as a secret, or an error if the process fails.
//
// Steps Involved:
//  1. Define GitHub password as a secret.
//  2. Configure the Terragrunt module to use the .netrc file with the provided credentials.
//  3. Run a command inside the container to verify the .netrc file's contents.
func (m *Go) Terragrunt_CreateNetRcFileForGithub(ctx context.Context) (*dagger.Container, error) {
	passwordAsSecret := dag.SetSecret("mysecret", "ohboywhatapassword")

	// Configure it for GitHub
	targetModule := dag.Terragrunt().
		WithNewNetrcFileAsSecretGitHub("supersecretuser", passwordAsSecret)

	// Check if the .netrc file is created correctly
	out, err := targetModule.
		Ctr().
		WithExec([]string{"cat", "/root/.netrc"}).
		Stdout(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get netrc file configured for github: %w", err)
	}

	expectedContent := "machine github.com\nlogin supersecretuser\npassword ohboywhatapassword"
	if !strings.Contains(out, expectedContent) {
		return nil, errNetRCFileNotFound
	}

	return targetModule.Ctr(), nil
}

// Terragrunt_RunArbitraryCommand runs an arbitrary shell command in the test container.
//
// This function demonstrates how to execute a shell command within the container
// using the Terragrunt module.
//
// Parameters:
//
//	ctx - context for controlling the function lifetime.
//
// Returns:
//
//	A string containing the output of the executed command, or an error if the command fails or if the output is empty.
func (m *Go) Terragrunt_RunArbitraryCommand(ctx context.Context) (string, error) {
	targetModule := dag.Terragrunt().WithSource(m.TestDir)

	// Execute the 'ls -l' command
	out, err := targetModule.
		Ctr().
		WithExec([]string{"ls", "-l"}).
		Stdout(ctx)

	if err != nil {
		return "", fmt.Errorf("failed to run shell command: %w", err)
	}

	if out == "" {
		return "", errExpectedFoldersNotFound
	}

	return out, nil
}

// Terragrunt_CreateContainer initializes and returns a configured Dagger container.
//
// This method exemplifies the setup of a container within the Terragrunt module using the source directory.
//
// Parameters:
//   - ctx: The context for controlling the function's timeout and cancellation.
//
// Returns:
//   - A configured Dagger container if successful, or an error if the process fails.
//
// Steps Involved:
//  1. Configure the Terragrunt module with the source directory.
//  2. Run a command inside the container to check the OS information.
func (m *Go) Terragrunt_CreateContainer(ctx context.Context) (*dagger.Container, error) {
	targetModule := dag.
		Terragrunt().
		BaseAlpine().
		WithUtilitiesInAlpineContainer(). // Install utilities
		WithGitInAlpineContainer().       // Install git
		WithSource(m.TestDir)

	// Get the OS or container information
	_, err := targetModule.
		Ctr().WithExec([]string{"uname"}).
		Stdout(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get OS information: %w", err)
	}

	return targetModule.Ctr(), nil
}
