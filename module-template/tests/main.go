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

	"github.com/Excoriate/daggerverse/module-template/tests/internal/dagger"

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
	polTests.Go(m.TestUbuntuBase)
	polTests.Go(m.TestAlpineBase)
	polTests.Go(m.TestBusyBoxBase)
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
	// Test cloud-specific functions.
	polTests.Go(m.TestWithAWSKeys)
	polTests.Go(m.TestWithAzureCredentials)

	// From this point onwards, we're testing the specific functionality of the ModuleTemplate module.

	if err := polTests.Wait(); err != nil {
		return fmt.Errorf("there are some failed tests: %w", err)
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
func (m *Tests) TestTerminal() *dagger.Container {
	targetModule := dag.ModuleTemplate()

	_, _ = targetModule.
		Ctr().
		Stdout(context.Background())

	return targetModule.
		Ctr().
		Terminal()
}

// TestUbuntuBase tests that the target module is based on the Ubuntu 22.04 image.
//
// This function verifies that the target module is configured appropriately to use the base Ubuntu 22.04 image.
// It runs a command to get the OS version and confirms it matches "Ubuntu 22.04".
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if the Ubuntu image is not used or if the output is not as expected.
func (m *Tests) TestUbuntuBase(ctx context.Context) error {
	targetModule := dag.
		ModuleTemplate().
		BaseUbuntu(dagger.ModuleTemplateBaseUbuntuOpts{Version: "22.04"})

	out, err := targetModule.Ctr().
		WithExec([]string{"grep", "^ID=ubuntu$", "/etc/os-release"}).
		Stdout(ctx)

	if err != nil {
		return fmt.Errorf("failed to get Ubuntu image: %w", err)
	}

	if !strings.Contains(out, "ID=ubuntu") {
		return fmt.Errorf("%w, expected Ubuntu 22.04 or ID=ubuntu, got %s", errExpectedContentNotMatch, out)
	}

	return nil
}

// TestAlpineBase tests that the target module is based on the Alpine Linux v3.17.3 image.
//
// This function verifies that the target module is configured appropriately to use the base Alpine Linux v3.17.3 image.
// It runs a command to get the OS version and confirms it matches "Alpine Linux v3.17.3".
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if the Alpine image is not used or if the output is not as expected.
func (m *Tests) TestAlpineBase(ctx context.Context) error {
	targetModule := dag.ModuleTemplate().
		BaseAlpine(dagger.ModuleTemplateBaseAlpineOpts{Version: "3.17.3"})

	out, err := targetModule.Ctr().WithExec([]string{"cat", "/etc/os-release"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to get Alpine image: %w", err)
	}

	// Adjust the conditions to match the actual output.
	if !strings.Contains(out, "Alpine Linux") || !strings.Contains(out, "VERSION_ID=3.17.3") {
		return fmt.Errorf("%w, expected Alpine Linux VERSION_ID=3.17.3, got %s", errExpectedContentNotMatch, out)
	}

	return nil
}

// TestBusyBoxBase tests that the target module is based on the BusyBox v1.35.0 image.
//
// This function verifies that the target module is configured appropriately to use the base BusyBox v1.35.0 image.
// It runs a command to get the OS version and confirms it matches "BusyBox v1.35.0".
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if the BusyBox image is not used or if the output is not as expected.
func (m *Tests) TestBusyBoxBase(ctx context.Context) error {
	targetModule := dag.
		ModuleTemplate().
		BaseBusyBox(dagger.ModuleTemplateBaseBusyBoxOpts{Version: "1.35.0"})

	out, err := targetModule.Ctr().
		WithExec([]string{"busybox", "--help"}).
		Stdout(ctx)

	if err != nil {
		return fmt.Errorf("failed to get BusyBox image: %w", err)
	}

	if !strings.Contains(out, "v1.35.0") {
		return fmt.Errorf("%w, expected BusyBox v1.35.0, got %s", errExpectedContentNotMatch, out)
	}

	return nil
}

// TestPassingEnvVarsInConstructor tests if the environment variables are passed correctly in the constructor.
//
// This is a helper method for tests, in order to test if the env vars are passed correctly in the constructor.
func (m *Tests) TestPassingEnvVarsInConstructor(ctx context.Context) error {
	targetModule := dag.
		ModuleTemplate(dagger.ModuleTemplateOpts{
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
		ModuleTemplate().
		WithEnvironmentVariable("HOST", "localhost", dagger.ModuleTemplateWithEnvironmentVariableOpts{
			Expand: false,
		})

	targetModule = targetModule.
		WithEnvironmentVariable("PORT", "8080", dagger.ModuleTemplateWithEnvironmentVariableOpts{
			Expand: false,
		})

	targetModule = targetModule.
		WithEnvironmentVariable("USER", "me", dagger.ModuleTemplateWithEnvironmentVariableOpts{
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
	// Retrieve the environment variables using the ModuleTemplate's PrintEnvVars function.
	envVars, err := dag.ModuleTemplate().PrintEnvVars(ctx)

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
	targetModule := dag.ModuleTemplate()

	// Define the environment variable key and value.
	key := "SOMETHING"
	value := "SOMETHING"

	// Set the environment variable in the target module.
	targetModule = targetModule.
		WithEnvironmentVariable(key, value, dagger.ModuleTemplateWithEnvironmentVariableOpts{
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
		ModuleTemplate()

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
		ModuleTemplate()

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
	targetModule := dag.ModuleTemplate()

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
	targetModule := dag.ModuleTemplate()

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
	targetModule := dag.ModuleTemplate()

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
	targetModule := dag.ModuleTemplate().
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
	targetModule := dag.ModuleTemplate()

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
		WithDownloadedFile(fileURL, dagger.ModuleTemplateWithDownloadedFileOpts{
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
	targetModule := dag.ModuleTemplate()

	// This is a public repository, the token isn't required.
	targetModule = targetModule.
		WithClonedGitRepo("https://github.com/excoriate/daggerverse",
			dagger.ModuleTemplateWithClonedGitRepoOpts{})

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
	targetModule := dag.ModuleTemplate()

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
	targetModule := dag.ModuleTemplate()

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
	targetModule := dag.ModuleTemplate()

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

// TestWithAWSKeys tests the setting of AWS keys as environment variables within the target module's container.
//
// This method creates secrets for AWS Access Key ID and AWS Secret Access Key, sets these secrets
// as environment variables in the target module's container, and verifies if the expected environment
// variables are set by running the `printenv` command within the container.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if there is an issue creating secrets, setting environment variables,
//     executing commands in the container, or if the output does not contain the expected environment variables.
func (m *Tests) TestWithAWSKeys(ctx context.Context) error {
	targetModule := dag.ModuleTemplate()

	awsKeyID := dag.SetSecret("AWS_ACCESS_KEY_ID", "AKIAIOSFODNN7EXAMPLE")
	awsSecretAccessKey := dag.SetSecret("AWS_SECRET_ACCESS_KEY", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")

	// With required AWS keys only.
	targetModule = targetModule.
		WithAwskeys(awsKeyID, awsSecretAccessKey)

	out, err := targetModule.Ctr().
		WithExec([]string{"printenv"}).
		Stdout(ctx)

	if err != nil {
		return fmt.Errorf("failed to get env vars: %w", err)
	}

	if !strings.Contains(out, "AWS_ACCESS_KEY_ID") {
		return fmt.Errorf("%w, expected AWS_ACCESS_KEY_ID to be set, got %s", errExpectedContentNotMatch, out)
	}

	if !strings.Contains(out, "AWS_SECRET_ACCESS_KEY") {
		return fmt.Errorf("%w, expected AWS_SECRET_ACCESS_KEY to be set, got %s", errExpectedContentNotMatch, out)
	}

	// Check if the content of the env vars is correct.
	if !strings.Contains(out, "AKIAIOSFODNN7EXAMPLE") {
		return fmt.Errorf("%w, expected AKIAIOSFODNN7EXAMPLE to be set, got %s", errExpectedContentNotMatch, out)
	}

	if !strings.Contains(out, "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY") {
		return fmt.Errorf("%w, "+
			"expected wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY to be set, got %s",
			errExpectedContentNotMatch, out)
	}

	return nil
}

// TestWithAzureCredentials tests the setting of Azure credentials as
// environment variables within the target module's container.
//
// This method creates secrets for Azure Client ID, Azure Client Secret, and Azure Tenant ID,
// sets these secrets as environment variables in the target module's container, and verifies if the expected
// environment variables are set by running the `printenv` command within the container.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if there is an issue creating secrets, setting environment variables,
//     executing commands in the container, or if the output does not contain the expected environment variables.
func (m *Tests) TestWithAzureCredentials(ctx context.Context) error {
	targetModule := dag.ModuleTemplate()

	azureClientID := dag.SetSecret("AZURE_CLIENT_ID", "00000000-0000-0000-0000-000000000000")
	azureClientSecret := dag.SetSecret("AZURE_CLIENT_SECRET", "00000000-0000-0000-0000-000000000000")
	azureTenantID := dag.SetSecret("AZURE_TENANT_ID", "00000000-0000-0000-0000-000000000000")

	// With required Azure keys only.
	targetModule = targetModule.
		WithAzureCredentials(azureClientID, azureClientSecret, azureTenantID)

	out, err := targetModule.Ctr().
		WithExec([]string{"printenv"}).
		Stdout(ctx)

	if err != nil {
		return fmt.Errorf("failed to get env vars: %w", err)
	}

	if !strings.Contains(out, "AZURE_CLIENT_ID") {
		return fmt.Errorf("%w, expected AZURE_CLIENT_ID to be set, got %s", errExpectedContentNotMatch, out)
	}

	if !strings.Contains(out, "AZURE_CLIENT_SECRET") {
		return fmt.Errorf("%w, expected AZURE_CLIENT_SECRET to be set, got %s", errExpectedContentNotMatch, out)
	}

	if !strings.Contains(out, "AZURE_TENANT_ID") {
		return fmt.Errorf("%w, expected AZURE_TENANT_ID to be set, got %s", errExpectedContentNotMatch, out)
	}

	return nil
}
