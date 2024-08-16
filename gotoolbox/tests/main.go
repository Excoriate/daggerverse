// A generated module for test the Gotoolbox functions
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

	"github.com/Excoriate/daggerverse/gotoolbox/tests/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

const emptyErrMsg = "the test output expected is empty"
const expectedContentNotMatchMsg = "an expected value does not match the actual value"
const underlyingDaggerErrMsg = "the dagger command failed or dagger returned an error"

var errEmptyOutput = errors.New(emptyErrMsg)
var errExpectedContentNotMatch = errors.New(expectedContentNotMatchMsg)
var errUnderlyingDagger = errors.New(underlyingDaggerErrMsg)

const aptGetCMD = "apt-get"

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
func (m *Tests) TestAll(ctx context.Context) error {
	polTests := pool.New().WithErrors().WithContext(ctx)

	// Test different ways to configure the base container.
	polTests.Go(m.TestBaseContainer)
	polTests.Go(m.TestPassingEnvVarsInConstructor)
	// Test built-in commands
	polTests.Go(m.TestRunShellCMD)
	polTests.Go(m.TestPrintEnvVars)
	polTests.Go(m.TestInspectEnvVar)
	// Test API(s) usage scenarios. APIs -> With<something>
	polTests.Go(m.TestWithContainer)
	polTests.Go(m.TestWithSource)
	polTests.Go(m.TestWithEnvironmentVariable)
	polTests.Go(m.TestWithUtilitiesInAlpineContainer)
	polTests.Go(m.TestWithSecretAsEnvVar)
	polTests.Go(m.TestWithCacheBuster)

	// From this point onwards, we're testing the specific functionality of the Gotoolbox module.

	if err := polTests.Wait(); err != nil {
		return fmt.Errorf("there are some failed tests: %w", err)
	}

	return nil
}

// TestBaseContainer tests the base Go toolbox container.
//
// This method verifies that the base Go toolbox container is properly set up by checking the Go version
// and the underlying OS distribution. It runs commands within the container to retrieve the Go version
// and to check the OS release information. The method confirms that the Go version matches the expected
// version and that the OS is Alpine Linux.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if there is an issue executing commands in the container,
//     if the Go version does not match the expected version, or if the OS is not Alpine Linux.
func (m *Tests) TestBaseContainer(ctx context.Context) error {
	targetModule := dag.Gotoolbox()

	// Run a command inside the container that returns the Go version.
	out, err := targetModule.Ctr().
		WithExec([]string{"go", "version"}).
		Stdout(ctx)

	// Check for errors executing the command to get the Go version.
	if err != nil {
		return WrapError(err, "failed to get Go version")
	}

	// Verify that the output contains the expected Go version.
	if !strings.Contains(out, "go version go1.22.3 ") {
		return WrapErrorf(err, "expected Go version go1.22.3, got %s", out)
	}

	// Run a command to check the OS release information.
	out, err = targetModule.Ctr().
		WithExec([]string{"cat", "/etc/os-release"}).
		Stdout(ctx)

	// Check for errors executing the command to get the OS release information.
	if err != nil {
		return WrapError(err, "failed to get OS release information")
	}

	// Verify that the output contains "Alpine Linux".
	if !strings.Contains(out, "Alpine Linux") {
		return WrapErrorf(err, "expected Alpine Linux OS, got %s", out)
	}

	// Passing a specific version to the container.
	targetModuleSpecificVersion := dag.Gotoolbox(dagger.GotoolboxOpts{
		Version: "1.20.14",
	})

	// Run a command inside the container that returns the Go version.
	outSpecificVersion, err := targetModuleSpecificVersion.Ctr().
		WithExec([]string{"go", "version"}).
		Stdout(ctx)

	// Check for errors executing the command to get the Go version.
	if err != nil {
		return WrapErrorf(err, "failed to get Go version, the output was: %s", outSpecificVersion)
	}

	// Verify that the output contains the expected Go version.
	if !strings.Contains(outSpecificVersion, "go version go1.20.14 ") {
		return WrapErrorf(err, "expected Go version go1.21.3, got %s", outSpecificVersion)
	}

	// Pass a container with another different version of go, in an ubuntu image.
	ctrUbuntuWithGO := dag.Container().
		From("ubuntu:latest").
		WithExec([]string{aptGetCMD, "update"}).
		WithExec([]string{aptGetCMD, "install", "-y", "curl"}).
		WithExec([]string{"curl", "-L", "https://dl.google.com/go/go1.22.3.linux-amd64.tar.gz", "-o", "/tmp/go.tar.gz"}).
		WithExec([]string{"tar", "-C", "/usr/local", "-xzf", "/tmp/go.tar.gz"}).
		WithExec([]string{"rm", "/tmp/go.tar.gz"}).
		WithEnvVariable("GOPATH", "/usr/local/go").
		WithEnvVariable("PATH", "/usr/local/go/bin:$PATH")

	targetModuleGoInUbuntu := dag.Gotoolbox(dagger.GotoolboxOpts{
		Ctr: ctrUbuntuWithGO,
	})

	// Run a command inside the container that returns the Go version, and
	// validate it is ubuntu the OS
	// and, it has go 1.22.3 installed.
	outUbuntu, ubuntuErr := targetModuleGoInUbuntu.Ctr().
		WithExec([]string{"go", "version"}).
		Stdout(ctx)

	// Check for errors executing the command to get the Go version.
	if ubuntuErr != nil {
		return WrapErrorf(ubuntuErr, "failed to get Go version, the output was: %s", outUbuntu)
	}

	// Verify that the output contains the expected Go version.
	if !strings.Contains(outUbuntu, "go version go1.22.3 ") {
		return WrapErrorf(err, "expected Go version go1.22.3, got %s", outUbuntu)
	}

	return nil
}

// TestWithContainer tests if the container is set correctly.
//
// This API is used to override the container set in the Dagger instance.
func (m *Tests) TestWithContainer(ctx context.Context) error {
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

	targetModule := dag.Gotoolbox()
	targetModule = targetModule.
		WithContainer(newContainer)

	// Specific Ubuntu command that only works in Ubuntu.
	cmd := []string{"lsb_release", "-a"}
	out, err := targetModule.Ctr().
		WithExec(cmd).
		Stdout(ctx)

	if err != nil {
		return WrapErrorf(err, "failed to validate an specific ubuntu command: %s", out)
	}

	if out == "" {
		return WrapError(err, "failed to validate the overridden container, got empty output")
	}

	if !strings.Contains(out, "Ubuntu") {
		return WrapErrorf(err, "failed to validate the container, got %s instead of Ubuntu", out)
	}

	return nil
}

// TestTerminal returns a terminal for testing.
//
// This is a helper method for tests, in order to get a terminal for testing purposes.
func (m *Tests) TestTerminal() *dagger.Container {
	targetModule := dag.Gotoolbox()

	_, _ = targetModule.
		Ctr().
		Stdout(context.Background())

	return targetModule.
		Ctr().
		Terminal()
}

// TestPassingEnvVarsInConstructor tests if the environment variables are passed correctly in the constructor.
//
// This is a helper method for tests, in order to test if the env vars are passed correctly in the constructor.
func (m *Tests) TestPassingEnvVarsInConstructor(ctx context.Context) error {
	targetModule := dag.
		Gotoolbox(dagger.GotoolboxOpts{
			EnvVarsFromHost: []string{"HOST=localhost", "PORT=8080", "USER=me", "PASS=1234"},
		})

	out, err := targetModule.
		Ctr().
		WithExec([]string{"printenv"}).
		Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to get env vars")
	}

	if !strings.Contains(out, "HOST=localhost") {
		return WrapErrorf(err, "expected HOST to be set, got %s", out)
	}

	if !strings.Contains(out, "PORT=8080") {
		return WrapErrorf(err, "expected PORT to be set, got %s", out)
	}

	if !strings.Contains(out, "USER=me") {
		return WrapErrorf(err, "expected USER to be set, got %s", out)
	}

	return nil
}

// TestWithEnvironmentVariable tests if the environment variables are passed correctly in the API.
//
// This is a helper method for tests, in order to test if the env vars are passed correctly in the API.
func (m *Tests) TestWithEnvironmentVariable(ctx context.Context) error {
	targetModule := dag.
		Gotoolbox().
		WithEnvironmentVariable("HOST", "localhost", dagger.GotoolboxWithEnvironmentVariableOpts{
			Expand: false,
		})

	targetModule = targetModule.
		WithEnvironmentVariable("PORT", "8080", dagger.GotoolboxWithEnvironmentVariableOpts{
			Expand: false,
		})

	targetModule = targetModule.
		WithEnvironmentVariable("USER", "me", dagger.GotoolboxWithEnvironmentVariableOpts{
			Expand: false,
		})

	out, err := targetModule.
		Ctr().
		WithExec([]string{"printenv"}).
		Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to get env vars")
	}

	if out == "" {
		return WrapErrorf(err, "expected to have at least one env var, got empty output")
	}

	if !strings.Contains(out, "HOST=localhost") {
		return WrapErrorf(err, "expected HOST to be set, got %s", out)
	}

	if !strings.Contains(out, "PORT=8080") {
		return WrapErrorf(err, "expected PORT to be set, got %s", out)
	}

	if !strings.Contains(out, "USER=me") {
		return WrapErrorf(err, "expected USER to be set, got %s", out)
	}

	return nil
}

// TestWithSource tests if the source directory is set correctly.
func (m *Tests) TestWithSource(ctx context.Context) error {
	targetModule := dag.
		Gotoolbox()

	targetModule.
		WithSource(m.TestDir)

	out, err := targetModule.
		Ctr().
		WithExec([]string{"ls", "-l"}).
		Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to get ls output")
	}

	if out == "" {
		return NewError("expected to have at least one folder, got empty output")
	}

	if !strings.Contains(out, "total") {
		return Errorf("expected to have at least one folder, got %s", out)
	}

	return nil
}

// TestRunShellCMD runs a shell command in the container.
//
// Tests if the shell command is executed correctly in the container.
func (m *Tests) TestRunShellCMD(ctx context.Context) error {
	targetModule := dag.
		Gotoolbox()

	out, err := targetModule.
		RunShell(ctx, "ls -l")

	if err != nil {
		return WrapError(err, "failed to run shell command")
	}

	if out == "" {
		return Errorf("expected to have at least one folder, got empty output")
	}

	if !strings.Contains(out, "total") {
		return Errorf("expected to have at least one folder, got %s", out)
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
	// Retrieve the environment variables using the Gotoolbox's PrintEnvVars function.
	envVars, err := dag.Gotoolbox().PrintEnvVars(ctx)

	// Check for errors retrieving the environment variables.
	if err != nil {
		return WrapError(err, "failed to get env vars")
	}

	// Check if the output is empty, which indicates no environment variables were found.
	if envVars == "" {
		return Errorf("expected to have at least one env var, got empty output")
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
	targetModule := dag.Gotoolbox()

	// Define the environment variable key and value.
	key := "SOMETHING"
	value := "SOMETHING"

	// Set the environment variable in the target module.
	targetModule = targetModule.
		WithEnvironmentVariable(key, value, dagger.GotoolboxWithEnvironmentVariableOpts{
			Expand: true,
		})

	// Inspect the environment variable in the container.
	out, err := targetModule.InspectEnvVar(ctx, key)
	if err != nil {
		return WrapErrorf(err, "failed to inspect env var %s", key)
	}

	// Check if the inspected value matches the expected result.
	if !strings.Contains(out, value) {
		return Errorf("expected %s to be %s, got %s", key, value, out)
	}

	// Return nil if the environment variable was correctly set and inspected.
	return nil
}

// TestWithUtilitiesInAlpineContainer tests if the Alpine container with utilities is set correctly.
//
// This method verifies that the Alpine container includes specific utilities by running a command within the container.
// The test checks if the 'curl' utility is available and functioning as expected.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
// - error: Returns an error if the utility check fails or the output does not contain the expected results.
func (m *Tests) TestWithUtilitiesInAlpineContainer(ctx context.Context) error {
	// Initialize the target module with the Alpine container and utilities.
	targetModule := dag.
		Gotoolbox()

	targetModule = targetModule.WithUtilitiesInAlpineContainer()

	// Execute the 'curl --version' command within the container to verify 'curl' utility.
	out, err := targetModule.
		Ctr().
		WithExec([]string{"curl", "--version"}).
		Stdout(ctx)

	if err != nil {
		return Errorf("failed to run shell command: %s", out)
	}

	if out == "" {
		return Errorf("expected to have at least one folder, got empty output")
	}

	if !strings.Contains(out, "curl") {
		return Errorf("expected 'curl' to be available in the container, got %s", out)
	}

	return nil
}

// TestWithSecretAsEnvVar tests the setting of environment variables using secrets.
//
// This method verifies that environment variables are correctly set in the container using secrets.
// It creates secrets for AWS, GCP, and another example, then sets these secrets as environment variables
// in the target module's container. The method runs the `printenv` command within the container and checks
// if the output contains the expected environment variables.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if there is an issue creating secrets, setting environment variables,
//     executing commands in the container, or if the output does not contain the expected environment variables.
func (m *Tests) TestWithSecretAsEnvVar(ctx context.Context) error {
	// Create secrets for AWS, GCP, and another example.
	secretAWS := dag.SetSecret("AWS_ACCESS_KEY_ID", "AKIAIOSFODNN7EXAMPLE")
	secretGCP := dag.SetSecret("GCP_PROJECT_ID", "example-project-id")
	secretAnother := dag.SetSecret("ANOTHER_SECRET", "another-secret-value")

	// Initialize the target module and set secrets as environment variables.
	targetModule := dag.Gotoolbox().
		WithSecretAsEnvVar("AWS_ACCESS_KEY_ID", secretAWS).
		WithSecretAsEnvVar("GCP_PROJECT_ID", secretGCP).
		WithSecretAsEnvVar("ANOTHER_SECRET", secretAnother)

	// Run the 'printenv' command within the container to check environment variables.
	out, err := targetModule.Ctr().
		WithExec([]string{"printenv"}).
		Stdout(ctx)

	// Check for errors executing the command.
	if err != nil {
		return Errorf("failed to get env vars: %s", out)
	}

	// Check if the expected environment variables are set.
	if !strings.Contains(out, "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE") {
		return Errorf("expected AWS_ACCESS_KEY_ID to be set, got %s", out)
	}

	if !strings.Contains(out, "GCP_PROJECT_ID=example-project-id") {
		return Errorf("expected GCP_PROJECT_ID to be set, got %s", out)
	}

	if !strings.Contains(out, "ANOTHER_SECRET=another-secret-value") {
		return Errorf("expected ANOTHER_SECRET to be set, got %s", out)
	}

	// Return nil if all expected environment variables are set.
	return nil
}

// TestWithCacheBuster tests the setting of a cache-busting environment variable
// within the target module's container.
//
// This method sets a cache-busting environment variable (`CACHE_BUSTER`) in
// the target module's container and verifies if it is correctly set by running
// the `printenv` command within the container.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if there is an issue setting the environment variable,
//     executing commands in the container, or if the `CACHE_BUSTER` environment variable
//     is not found in the output.
func (m *Tests) TestWithCacheBuster(ctx context.Context) error {
	targetModule := dag.Gotoolbox()

	// Set the cache-busting environment variable
	targetModule = targetModule.WithCacheBuster()

	out, err := targetModule.Ctr().
		WithExec([]string{"printenv"}).
		Stdout(ctx)

	if err != nil {
		return Errorf("failed to get env vars: %s", out)
	}

	// Check if the cache-busting environment variable is set
	if !strings.Contains(out, "CACHE_BUSTER") {
		return Errorf("expected CACHE_BUSTER to be set, got %s", out)
	}

	return nil
}
