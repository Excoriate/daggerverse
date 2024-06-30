// A generated module for test the ModuleTemplate functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.
package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sourcegraph/conc/pool"
)

const emptyErrMsg = "the test output expected is empty"
const expectedContentNotMatchMsg = "an expected value does not match the actual value"
const underlyingDaggerErrMsg = "the dagger command failed or dagger returned an error"

var errEmptyOutput = errors.New(emptyErrMsg)
var errExpectedContentNotMatch = errors.New(expectedContentNotMatchMsg)
var errUnderlyingDagger = errors.New(underlyingDaggerErrMsg)

// Tests is a collection of tests.
//
// It's a struct that contains a single field, TestDir, which is a pointer to a Directory.
type Tests struct {
	TestDir *Directory
}

// New creates a new Tests instance.
//
// It's the initial constructor for the Tests struct.
func New() *Tests {
	t := &Tests{}

	t.TestDir = t.getTestDir()

	return t
}

// getTestDir returns the test directory.
//
// This is a helper method for tests, in order to get the test directory which
// is located in the same directory as the test file, and normally named as "testdata".
func (m *Tests) getTestDir() *Directory {
	return dag.CurrentModule().Source().Directory("./testdata")
}

// TestAll executes all tests.
//
// This is a helper method for tests, in order to execute all tests.
func (m *Tests) TestAll(ctx context.Context) error {
	polTests := pool.New().WithErrors().WithContext(ctx)

	polTests.Go(m.TestAPIWithContainer)
	polTests.Go(m.TestAPIPassingEnvVarsInConstructor)
	polTests.Go(m.TestAPIWithSource)
	polTests.Go(m.TestAPIPassingEnvVars)
	polTests.Go(m.TestRunShellCMD)

	if err := polTests.Wait(); err != nil {
		return fmt.Errorf("there are some failed tests: %w", err)
	}

	return nil
}

// TestAPIWithContainer tests if the container is set correctly.
//
// This API is used to override the container set in the Dagger instance.
func (m *Tests) TestAPIWithContainer(ctx context.Context) error {
	// Create a new container from the ubuntu:latest image
	newContainer := dag.Container().From("ubuntu:latest")

	// Ensure the container runs with root permissions
	newContainer = newContainer.
		WithUser("root")

	// Update package list
	updateCmd := []string{"apt-get", "update"}
	newContainer = newContainer.
		WithExec(updateCmd)

	// Install lsb-release package
	installCmd := []string{"apt-get", "install", "-y", "lsb-release"}
	newContainer = newContainer.
		WithExec(installCmd)

	targetModule := dag.ModuleTemplate()
	targetModule = targetModule.
		WithContainer(newContainer)

	// Specific Ubuntu command that only works in Ubuntu.
	cmd := []string{"lsb_release", "-a"}
	out, err := targetModule.Ctr().
		WithExec(cmd).
		Stdout(ctx)

	if err != nil {
		return fmt.Errorf("%w, failed to validate an specific ubuntu command: %w", errUnderlyingDagger, err)
	}

	if out == "" {
		return fmt.Errorf("%w, failed to validate the overridden container, got empty output", errUnderlyingDagger)
	}

	if !strings.Contains(out, "Ubuntu") {
		return fmt.Errorf("%w, failed to validate the container, got %s instead of Ubuntu", errUnderlyingDagger, out)
	}

	return nil
}

// TestTerminal returns a terminal for testing.
//
// This is a helper method for tests, in order to get a terminal for testing purposes.
func (m *Tests) TestTerminal() *Terminal {
	targetModule := dag.ModuleTemplate()

	_, _ = targetModule.
		Ctr().
		Stdout(context.Background())

	return targetModule.
		Ctr().
		Terminal()
}

// TestAPIPassingEnvVarsInConstructor tests if the environment variables are passed correctly in the constructor.
//
// This is a helper method for tests, in order to test if the env vars are passed correctly in the constructor.
func (m *Tests) TestAPIPassingEnvVarsInConstructor(ctx context.Context) error {
	targetModule := dag.
		ModuleTemplate(ModuleTemplateOpts{
			EnvVarsFromHost: []string{"HOST=localhost", "PORT=8080", "USER=me", "PASS=1234"},
		})

	out, err := targetModule.
		Ctr().
		WithExec([]string{"printenv"}).
		Stdout(ctx)

	if err != nil {
		return fmt.Errorf("failed to get env vars: %w", err)
	}

	if !strings.Contains(out, "HOST=localhost") {
		return fmt.Errorf("%w, expected HOST to be set, got %s", errExpectedContentNotMatch, out)
	}

	if !strings.Contains(out, "PORT=8080") {
		return fmt.Errorf("%w, expected PORT to be set, got %s", errExpectedContentNotMatch, out)
	}

	if !strings.Contains(out, "USER=me") {
		return fmt.Errorf("%w, expected USER to be set, got %s", errExpectedContentNotMatch, out)
	}

	return nil
}

// TestAPIPassingEnvVars tests if the environment variables are passed correctly in the API.
//
// This is a helper method for tests, in order to test if the env vars are passed correctly in the API.
func (m *Tests) TestAPIPassingEnvVars(ctx context.Context) error {
	targetModule := dag.
		ModuleTemplate().
		WithEnvironmentVariable("HOST", "localhost", ModuleTemplateWithEnvironmentVariableOpts{
			Expand: false,
		})

	targetModule = targetModule.
		WithEnvironmentVariable("PORT", "8080", ModuleTemplateWithEnvironmentVariableOpts{
			Expand: false,
		})

	targetModule = targetModule.
		WithEnvironmentVariable("USER", "me", ModuleTemplateWithEnvironmentVariableOpts{
			Expand: false,
		})

	out, err := targetModule.
		Ctr().
		WithExec([]string{"printenv"}).
		Stdout(ctx)

	if err != nil {
		return fmt.Errorf("failed to get env vars: %w", err)
	}

	if out == "" {
		return fmt.Errorf("%w, expected to have at least one env var, got empty output", errEmptyOutput)
	}

	if !strings.Contains(out, "HOST=localhost") {
		return fmt.Errorf("%w, expected HOST to be set, got %s", errExpectedContentNotMatch, out)
	}

	if !strings.Contains(out, "PORT=8080") {
		return fmt.Errorf("%w, expected PORT to be set, got %s", errExpectedContentNotMatch, out)
	}

	if !strings.Contains(out, "USER=me") {
		return fmt.Errorf("%w, expected USER to be set, got %s", errExpectedContentNotMatch, out)
	}

	return nil
}

// TestAPIWithSource tests if the source directory is set correctly.
func (m *Tests) TestAPIWithSource(ctx context.Context) error {
	targetModule := dag.
		ModuleTemplate()

	targetModule.
		WithSource(m.TestDir)

	out, err := targetModule.
		Ctr().
		WithExec([]string{"ls", "-l"}).
		Stdout(ctx)

	if err != nil {
		return fmt.Errorf("failed to get ls output: %w", err)
	}

	if out == "" {
		return fmt.Errorf("%w, expected to have at least one folder, got empty output", errEmptyOutput)
	}

	if !strings.Contains(out, "total") {
		return fmt.Errorf("%w, expected to have at least one folder, got %s", errExpectedContentNotMatch, out)
	}

	return nil
}

// TestRunShellCMD runs a shell command in the container.
//
// Tests if the shell command is executed correctly in the container.
func (m *Tests) TestRunShellCMD(ctx context.Context) error {
	targetModule := dag.
		ModuleTemplate()

	out, err := targetModule.
		RunShell(ctx, "ls -l")

	if err != nil {
		return fmt.Errorf("%w, failed to run shell command: %w", errUnderlyingDagger, err)
	}

	if out == "" {
		return fmt.Errorf("%w, expected to have at least one folder, got empty output", errEmptyOutput)
	}

	if !strings.Contains(out, "total") {
		return fmt.Errorf("%w, expected to have at least one folder, got %s", errExpectedContentNotMatch, out)
	}

	return nil
}
