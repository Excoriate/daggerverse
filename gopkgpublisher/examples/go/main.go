// A generated module for Go example functions
package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Excoriate/daggerverse/gopkgpublisher/examples/go/internal/dagger"
	"github.com/sourcegraph/conc/pool"
)

// Go is a Dagger module that exemplifies the usage of the Gopkgpublisher module.
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
// other specific functionalities of the Gopkgpublisher module.
func (m *Go) AllRecipes(ctx context.Context) error {
	polTests := pool.New().WithErrors().WithContext(ctx)

	// Test different ways to configure the base container.
	polTests.Go(m.BuiltInRecipes)
	// From this point onwards, we're testing the specific functionality of the Gopkgpublisher module.

	if err := polTests.Wait(); err != nil {
		return fmt.Errorf("there are some failed tests: %w", err)
	}

	return nil
}

// BuiltInRecipes demonstrates how to run built-in recipes
//
// This method showcases the use of various built-in recipes provided by the ModuleTemplate
// module, including creating a container, running an arbitrary command, and creating a .netrc
// file for GitHub authentication.
//
// Parameters:
//   - ctx: The context for controlling the function's timeout and cancellation.
//
// Returns:
//   - An error if any of the internal methods fail, or nil otherwise.
func (m *Go) BuiltInRecipes(ctx context.Context) error {
	// Pass environment variables to the ModuleTemplate module using ModuleTemplate_PassedEnvVars
	if err := m.Gopkgpublisher_PassedEnvVars(ctx); err != nil {
		return fmt.Errorf("failed to pass environment variables: %w", err)
	}

	// Run an arbitrary command in the container using Gopkgpublisher_RunArbitraryCommand
	if _, err := m.Gopkgpublisher_RunArbitraryCommand(ctx); err != nil {
		return fmt.Errorf("failed to run arbitrary command: %w", err)
	}

	// Create a new container.
	if _, err := m.Gopkgpublisher_CreateContainer(ctx); err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	return nil
}

// Gopkgpublisher_PassedEnvVars demonstrates how to pass environment variables to the Gopkgpublisher module.
//
// This method configures a Gopkgpublisher module to use specific environment variables from the host.
func (m *Go) Gopkgpublisher_PassedEnvVars(ctx context.Context) error {
	targetModule := dag.Gopkgpublisher(dagger.GopkgpublisherOpts{
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

// Gopkgpublisher_OpenTerminal demonstrates how to open an interactive terminal session
// within a Gopkgpublisher module container.
//
// This function showcases the initialization and configuration of a
// Gopkgpublisher module container using various options like enabling Cgo,
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
func (m *Go) Gopkgpublisher_OpenTerminal() *dagger.Container {
	// Configure the Gopkgpublisher module container with necessary options
	targetModule := dag.Gopkgpublisher()

	// Retrieve and discard standard output
	_, _ = targetModule.Ctr().
		Stdout(context.Background())

	// Open and return the terminal session in the configured container
	return targetModule.Ctr().
		Terminal()
}

// Gopkgpublisher_RunArbitraryCommand runs an arbitrary shell command in the test container.
//
// This function demonstrates how to execute a shell command within the container
// using the Gopkgpublisher module.
//
// Parameters:
//
//	ctx - context for controlling the function lifetime.
//
// Returns:
//
//	A string containing the output of the executed command, or an error if the command fails or if the output is empty.
func (m *Go) Gopkgpublisher_RunArbitraryCommand(ctx context.Context) (string, error) {
	targetModule := dag.Gopkgpublisher().WithSource(m.TestDir)

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

// Gopkgpublisher_CreateContainer initializes and returns a configured Dagger container.
//
// This method exemplifies the setup of a container within the Gopkgpublisher module using the source directory.
//
// Parameters:
//   - ctx: The context for controlling the function's timeout and cancellation.
//
// Returns:
//   - A configured Dagger container if successful, or an error if the process fails.
//
// Steps Involved:
//  1. Configure the Gopkgpublisher module with the source directory.
//  2. Run a command inside the container to check the OS information.
func (m *Go) Gopkgpublisher_CreateContainer(ctx context.Context) (*dagger.Container, error) {
	targetModule := dag.
		Gopkgpublisher().
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
