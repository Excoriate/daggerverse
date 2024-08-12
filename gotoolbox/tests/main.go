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
	polTests.Go(m.TestWithNewNetrcFileGitHub)
	polTests.Go(m.TestWithNewNetrcFileAsSecretGitHub)
	polTests.Go(m.TestWithNewNetrcFileGitLab)
	polTests.Go(m.TestWithNewNetrcFileAsSecretGitLab)
	polTests.Go(m.TestWithSecretAsEnvVar)
	polTests.Go(m.TestWithDownloadedFile)
	polTests.Go(m.TestWithClonedGitRepo)
	polTests.Go(m.TestWithCacheBuster)
	// Test utility functions.
	polTests.Go(m.TestDownloadFile)
	polTests.Go(m.TestCloneGitRepo)

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
		return fmt.Errorf("failed to get Go version: %w", err)
	}

	// Verify that the output contains the expected Go version.
	if !strings.Contains(out, "go version go1.22.3 ") {
		return fmt.Errorf("expected Go version go1.22.3, got %s", out)
	}

	// Run a command to check the OS release information.
	out, err = targetModule.Ctr().
		WithExec([]string{"cat", "/etc/os-release"}).
		Stdout(ctx)

	// Check for errors executing the command to get the OS release information.
	if err != nil {
		return fmt.Errorf("failed to get OS release information: %w", err)
	}

	// Verify that the output contains "Alpine Linux".
	if !strings.Contains(out, "Alpine Linux") {
		return fmt.Errorf("expected Alpine Linux OS, got %s", out)
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
		return fmt.Errorf("failed to get Go version, the output was: %w", err)
	}

	// Verify that the output contains the expected Go version.
	if !strings.Contains(outSpecificVersion, "go version go1.20.14 ") {
		return fmt.Errorf("expected Go version go1.21.3, got %s", outSpecificVersion)
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
		return fmt.Errorf("failed to get Go version, the output was: %w", ubuntuErr)
	}

	// Verify that the output contains the expected Go version.
	if !strings.Contains(outUbuntu, "go version go1.22.3 ") {
		return fmt.Errorf("expected Go version go1.22.3, got %s", outUbuntu)
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
		return fmt.Errorf("failed to get ls output: %w", err)
	}

	if out == "" {
		return fmt.Errorf("%w, %s", errExpectedContentNotMatch, out)
	}

	if !strings.Contains(out, "total") {
		return fmt.Errorf("%w, %s", errExpectedContentNotMatch, out)
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
		return fmt.Errorf("failed to get env vars: %w", err)
	}

	// Check if the output is empty, which indicates no environment variables were found.
	if envVars == "" {
		return fmt.Errorf("%w, expected to have at least one env var, got empty output", errEmptyOutput)
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
		return fmt.Errorf("failed to inspect env var %s: %w", key, err)
	}

	// Check if the inspected value matches the expected result.
	if !strings.Contains(out, value) {
		return fmt.Errorf("%w, expected %s to be %s, got %s", errExpectedContentNotMatch, key, value, out)
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
		return fmt.Errorf("%w, failed to run shell command: %w", errUnderlyingDagger, err)
	}

	if out == "" {
		return fmt.Errorf("%w, expected to have at least one folder, got empty output", errEmptyOutput)
	}

	if !strings.Contains(out, "curl") {
		return fmt.Errorf("%w, expected 'curl' to be available in the container, got %s", errExpectedContentNotMatch, out)
	}

	return nil
}

// TestWithNewNetrcFileGitHub tests the creation of a new .netrc file with GitHub credentials.
//
// This function verifies that the GitHub credentials are set correctly in the .netrc file using a secret.
// It creates a new secret with the GitHub credentials and sets them in the target module's .netrc file.
// The function then reads the .netrc file from the container and checks if it contains the expected machine entry.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if the creation of the .netrc file fails, if the file's content does
//     not match the expected result,
//     or if there is an issue with executing commands in the container.
func (m *Tests) TestWithNewNetrcFileGitHub(ctx context.Context) error {
	targetModule := dag.
		Gotoolbox()

		// Create a new secret with the GitHub credentials.
	githubSecret := dag.SetSecret("github-username", "github-password")

	// Set the GitHub credentials as a secret in the target module's .netrc file.
	targetModule = targetModule.
		WithNewNetrcFileAsSecretGitHub("github-username", githubSecret)

	// Execute a command to read the .netrc file from the container.
	out, err := targetModule.
		Ctr().
		WithExec([]string{"cat", "/root/.netrc"}).
		Stdout(ctx)

	// Check for errors executing the command.
	if err != nil {
		return fmt.Errorf("%w, failed to get netrc file: %w", errUnderlyingDagger, err)
	}

	// Check if the .netrc file contains the expected machine entry.
	if !strings.Contains(out, "machine github.com") {
		return fmt.Errorf("%w, expected netrc file to be created, got %s", errExpectedContentNotMatch, out)
	}

	return nil
}

// TestWithNewNetrcFileAsSecretGitHub tests the creation of a new .netrc file with GitHub credentials.
//
// This function verifies that the GitHub credentials are set correctly in the .netrc file using a secret.
// It creates a new secret with the GitHub credentials and sets them in the target module's .netrc file.
// The function then reads the .netrc file from the container and checks if it contains the expected machine entry.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if the creation of the .netrc file fails, if the file's
//     content does not match the expected result,
//     or if there is an issue with executing commands in the container.
func (m *Tests) TestWithNewNetrcFileAsSecretGitHub(ctx context.Context) error {
	// Initialize the target module.
	targetModule := dag.Gotoolbox()

	// Create a new secret with the GitHub credentials.
	githubSecret := dag.SetSecret("github-username", "github-password")

	// Set the GitHub credentials as a secret in the target module's .netrc file.
	targetModule = targetModule.WithNewNetrcFileAsSecretGitHub("github-username", githubSecret)

	// Execute a command to read the .netrc file from the container.
	out, err := targetModule.Ctr().WithExec([]string{"cat", "/root/.netrc"}).Stdout(ctx)

	// Check for errors executing the command.
	if err != nil {
		return fmt.Errorf("%w, failed to get netrc file: %w", errUnderlyingDagger, err)
	}

	// Check if the .netrc file contains the expected machine entry.
	if !strings.Contains(out, "machine github.com") {
		return fmt.Errorf("%w, expected netrc file to be created, got %s", errExpectedContentNotMatch, out)
	}

	// Return nil if the netrc file is created and contains the expected entry.
	return nil
}

// TestWithNewNetrcFileGitLab tests the creation of a new .netrc file with GitLab credentials.
//
// This function verifies that the GitLab credentials are set correctly in the .netrc file using a secret.
// It creates a new secret with the GitLab credentials and sets them in the target module's .netrc file.
// The function then reads the .netrc file from the container and checks if it contains the expected machine entry.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if the creation of the .netrc file fails, if the file's
//     content does not match the expected result,
//     or if there is an issue with executing commands in the container.
func (m *Tests) TestWithNewNetrcFileGitLab(ctx context.Context) error {
	targetModule := dag.Gotoolbox()

	// Create a new secret with the GitLab credentials.
	gitlabSecret := dag.SetSecret("gitlab-username", "gitlab-password")

	// Set the GitLab credentials as a secret in the target module's .netrc file.
	targetModule = targetModule.WithNewNetrcFileAsSecretGitLab("gitlab-username", gitlabSecret)

	// Execute a command to read the .netrc file from the container.
	out, err := targetModule.Ctr().WithExec([]string{"cat", "/root/.netrc"}).Stdout(ctx)

	// Check for errors executing the command.
	if err != nil {
		return fmt.Errorf("%w, failed to get netrc file: %w", errUnderlyingDagger, err)
	}

	// Check if the .netrc file contains the expected machine entry.
	if !strings.Contains(out, "machine gitlab.com") {
		return fmt.Errorf("%w, expected netrc file to be created, got %s", errExpectedContentNotMatch, out)
	}

	return nil
}

// TestWithNewNetrcFileAsSecretGitLab creates a new .netrc file with GitLab credentials.
//
// This method verifies that the GitLab credentials are set correctly in the .netrc file using a secret.
// It creates a new secret with the GitLab credentials and sets them in the target module's .netrc file.
// The method then reads the .netrc file from the container and checks if it contains the expected machine entry.
//
// Arguments:
// - ctx (context.Context): The context for the method execution.
//
// Returns:
//   - error: Returns an error if the creation of the .netrc file fails, if the file's
//     content does not match the expected result,
//     or if there is an issue with executing commands in the container.
func (m *Tests) TestWithNewNetrcFileAsSecretGitLab(ctx context.Context) error {
	targetModule := dag.Gotoolbox()

	// Create a new secret with the GitLab credentials.
	gitlabSecret := dag.SetSecret("gitlab-username", "gitlab-password")

	// Set the GitLab credentials as a secret in the target module's .netrc file.
	targetModule = targetModule.WithNewNetrcFileAsSecretGitLab("gitlab-username", gitlabSecret)

	// Execute a command to read the .netrc file from the container.
	out, err := targetModule.Ctr().WithExec([]string{"cat", "/root/.netrc"}).Stdout(ctx)

	// Check for errors executing the command.
	if err != nil {
		return fmt.Errorf("%w, failed to get netrc file: %w", errUnderlyingDagger, err)
	}

	// Check if the .netrc file contains the expected machine entry.
	if !strings.Contains(out, "machine gitlab.com") {
		return fmt.Errorf("%w, expected netrc file to be created, got %s", errExpectedContentNotMatch, out)
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
		return fmt.Errorf("failed to get env vars: %w", err)
	}

	// Check if the expected environment variables are set.
	if !strings.Contains(out, "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE") {
		return fmt.Errorf("%w, expected AWS_ACCESS_KEY_ID to be set, got %s", errExpectedContentNotMatch, out)
	}

	if !strings.Contains(out, "GCP_PROJECT_ID=example-project-id") {
		return fmt.Errorf("%w, expected GCP_PROJECT_ID to be set, got %s", errExpectedContentNotMatch, out)
	}

	if !strings.Contains(out, "ANOTHER_SECRET=another-secret-value") {
		return fmt.Errorf("%w, expected ANOTHER_SECRET to be set, got %s", errExpectedContentNotMatch, out)
	}

	// Return nil if all expected environment variables are set.
	return nil
}

// TestWithDownloadedFile tests the downloading of a file from a URL.
//
// This method verifies that a file can be downloaded from a URL and mounted in the container.
// It downloads a file from a URL, mounts it in the container, and checks if the file exists.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if there is an issue downloading the file, mounting it in the container,
//     or if the file is not found.
func (m *Tests) TestWithDownloadedFile(ctx context.Context) error {
	targetModule := dag.Gotoolbox()

	// Download a file from a URL and mount it in the container, without filename passed.
	fileURL := "https://framerusercontent.com/assets/cNNFYmZqESeYTV5PHp72ay0O2o.zip"
	targetModule = targetModule.
		WithDownloadedFile(fileURL)

	// Check if the file exists in the container.
	out, err := targetModule.
		Ctr().
		WithExec([]string{"ls", "/mnt/cNNFYmZqESeYTV5PHp72ay0O2o.zip"}).
		Stdout(ctx)

	if err != nil {
		return fmt.Errorf("failed to get download file from url %s: %w", fileURL, err)
	}

	if out == "" {
		return fmt.Errorf("%w, expected to have at least one folder, got empty output", errEmptyOutput)
	}

	// Downloading the file but with a name this time instead.
	fileName := "mydownloadedfile.zip"
	targetModule = targetModule.
		WithDownloadedFile(fileURL, dagger.GotoolboxWithDownloadedFileOpts{
			DestFileName: fileName,
		})

	// Check if the file exists in the container.
	out, err = targetModule.
		Ctr().
		WithExec([]string{"ls", "/mnt/mydownloadedfile.zip"}).
		Stdout(ctx)

	if err != nil {
		return fmt.Errorf("failed to get download file from url %s: %w", fileURL, err)
	}

	if out == "" {
		return fmt.Errorf("%w, expected to have at least one folder, got empty output", errEmptyOutput)
	}

	return nil
}

// TestWithClonedGitRepo tests the WithClonedGitRepo function.
func (m *Tests) TestWithClonedGitRepo(ctx context.Context) error {
	targetModule := dag.Gotoolbox()

	// This is a public repository, the token isn't required.
	targetModule = targetModule.
		WithClonedGitRepo("https://github.com/excoriate/daggerverse",
			dagger.GotoolboxWithClonedGitRepoOpts{})

	out, err := targetModule.Ctr().
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
		return fmt.Errorf("failed to get env vars: %w", err)
	}

	// Check if the cache-busting environment variable is set
	if !strings.Contains(out, "CACHE_BUSTER") {
		return fmt.Errorf("%w, expected CACHE_BUSTER to be set, got %s", errExpectedContentNotMatch, out)
	}

	return nil
}

// TestCloneGitRepo tests the CloneGitRepo function.
func (m *Tests) TestCloneGitRepo(ctx context.Context) error {
	targetModule := dag.Gotoolbox()

	// This is a public repository, the token isn't required.
	clonedRepo := targetModule.
		CloneGitRepo("https://github.com/excoriate/daggerverse")

	// Mount it as a directory, and inspect if it contains .gitignore and LICENSE files.
	ctr := targetModule.Ctr().
		WithMountedDirectory("/mnt", clonedRepo)

	out, err := ctr.
		WithExec([]string{"ls", "-l", "/mnt"}).
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

	// Check if the .gitignore file is present.
	out, err = ctr.
		WithExec([]string{"cat", "/mnt/.gitignore"}).
		Stdout(ctx)

	if err != nil {
		return fmt.Errorf("failed to get .gitignore file: %w", err)
	}

	if out == "" {
		return fmt.Errorf("%w, expected to have at least one folder, got empty output", errEmptyOutput)
	}

	return nil
}

// TestDownloadFile tests the downloading of a file from a URL and mounts it in the container.
//
// This method verifies that a file can be downloaded from a URL, mounted
// in the container, and checks if the file exists.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if there is an issue downloading the file, mounting it in the container,
//     or if the file is not found in the mounted path.
func (m *Tests) TestDownloadFile(ctx context.Context) error {
	// Initialize the target module.
	targetModule := dag.Gotoolbox()

	// Define the URL of the file to be downloaded.
	fileURL := "https://framerusercontent.com/assets/cNNFYmZqESeYTV5PHp72ay0O2o.zip"

	// Download the file from the URL.
	fileDownloaded := targetModule.DownloadFile(fileURL)

	// Mount the downloaded file in the container at /mnt/myfile.zip.
	ctr := targetModule.
		Ctr().
		WithMountedFile("/mnt/myfile.zip", fileDownloaded)

	// Execute a command to check if the file exists in the container.
	out, err := ctr.
		WithExec([]string{"ls", "/mnt/myfile.zip"}).
		Stdout(ctx)

	// Check for errors executing the command.
	if err != nil {
		return fmt.Errorf("failed to get download file from url %s: %w", fileURL, err)
	}

	// Check if the output is empty, which indicates the file was not found.
	if out == "" {
		return fmt.Errorf("%w, expected to find the file at /mnt/myfile.zip, but got empty output", errEmptyOutput)
	}

	// Return nil if the file was successfully found.
	return nil
}
