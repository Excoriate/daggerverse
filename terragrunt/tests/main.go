// A generated module for test the Terragrunt functions
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
	"strings"

	"github.com/Excoriate/daggerverse/terragrunt/tests/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

// Tests is a collection of tests.
//
// It's a struct that contains a single field, TestDir, which is a pointer to a Directory.
type Tests struct {
	TestDir *dagger.Directory
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
func (m *Tests) getTestDir() *dagger.Directory {
	return dag.CurrentModule().Source().Directory("./testdata")
}

// TestAll executes all tests.
//
// This is a helper method for tests, in order to execute all tests.
//
//nolint:funlen // The test handles multiple commands and environments, requiring a longer function.
func (m *Tests) TestAll(ctx context.Context) error {
	polTests := pool.
		New().
		WithErrors().
		WithContext(ctx)

	// Test different ways to configure the base container.
	polTests.Go(m.TestPassingEnvVarsInConstructor)
	// Test built-in commands
	polTests.Go(m.TestRunShellCMD)
	polTests.Go(m.TestPrintEnvVars)
	polTests.Go(m.TestInspectEnvVar)
	// Test API(s) usage scenarios. APIs -> With<something>
	polTests.Go(m.TestWithContainer)
	polTests.Go(m.TestWithSource)
	polTests.Go(m.TestWithEnvironmentVariable)
	polTests.Go(m.TestWithDownloadedFile)
	polTests.Go(m.TestWithCacheBuster)
	// Test utility functions.
	// Test cloud-specific functions.
	polTests.Go(m.TestWithAWSKeys)
	polTests.Go(m.TestWithAzureCredentials)
	// From this point onwards, we're testing the specific functionality of the Terragrunt module.

	if err := polTests.Wait(); err != nil {
		return WrapError(err, "there are some failed tests")
	}

	return nil
}

// TestRunShellCMD runs a shell command in the container.
//
// Tests if the shell command is executed correctly in the container.
func (m *Tests) TestRunShellCMD(ctx context.Context) error {
	targetModule := dag.
		Terragrunt()

	out, err := targetModule.
		RunShell(ctx, "ls -l")

	if err != nil {
		return WrapError(err, "failed to run shell command")
	}

	if out == "" {
		return WrapError(err, "expected to have at least one folder, got empty output")
	}

	if !strings.Contains(out, "total") {
		return WrapErrorf(err, "expected to have at least one folder, got %s", out)
	}

	return nil
}

// TestPrintEnvVars tests the PrintEnvVars function.
//
// This method verifies that environment variables can be printed within the context
// of the target module's execution. It runs the `printenv` command within the container
// and checks if any environment variables are present.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if there is an issue printing environment variables,
//     or if no environment variables are found in the output.
func (m *Tests) TestPrintEnvVars(ctx context.Context) error {
	// Retrieve the environment variables using the Terragrunt's PrintEnvVars function.
	envVars, err := dag.Terragrunt().PrintEnvVars(ctx)

	// Check for errors retrieving the environment variables.
	if err != nil {
		return WrapError(err, "failed to get env vars")
	}

	// Check if the output is empty, which indicates no environment variables were found.
	if envVars == "" {
		return WrapError(err, "expected to have at least one env var, got empty output")
	}

	// Return nil if environment variables were successfully found in the output.
	return nil
}

// TestInspectEnvVar tests the inspection of an environment variable set in the container.
//
// This method verifies that an environment variable is correctly set in the target module's container.
// It sets an environment variable and then inspects it to check if the value matches the expected result.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if there is an issue setting the environment variable, inspecting the variable,
//     or if the inspected value does not match the expected result.
func (m *Tests) TestInspectEnvVar(ctx context.Context) error {
	// Initialize the target module.
	targetModule := dag.Terragrunt()

	// Define the environment variable key and value.
	key := "SOMETHING"
	value := "SOMETHING"

	// Set the environment variable in the target module.
	targetModule = targetModule.
		WithEnvironmentVariable(key, value, dagger.TerragruntWithEnvironmentVariableOpts{
			Expand: true,
		})

	// Inspect the environment variable in the container.
	out, err := targetModule.InspectEnvVar(ctx, key)
	if err != nil {
		return WrapErrorf(err, "failed to inspect env var %s", key)
	}

	// Check if the inspected value matches the expected result.
	if !strings.Contains(out, value) {
		return WrapErrorf(err, "expected %s to be %s, got %s", key, value, out)
	}

	// Return nil if the environment variable was correctly set and inspected.
	return nil
}
