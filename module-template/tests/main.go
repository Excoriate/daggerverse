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
	"encoding/json"
	"strings"

	"github.com/Excoriate/daggerverse/module-template/tests/internal/dagger"

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

// getGolangAlpineContainer returns a container with the Go toolchain installed.
//
// This function returns a container with the Go toolchain installed, which is
// suitable for testing Go related functionality.
func getGolangAlpineContainer(version string) *dagger.Container {
	if version == "" {
		version = "1.20.4"
	}

	return dag.Container().
		From("golang:" + version + "-alpine")
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
	polTests.Go(m.TestWithUtilitiesInUbuntuContainer)
	polTests.Go(m.TestWithGitInAlpineContainer)
	polTests.Go(m.TestWithGitInUbuntuContainer)
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
	// Test Go specific functions.
	polTests.Go(m.TestGoWithGoPlatform)
	polTests.Go(m.TestGoWithGoBuild)
	polTests.Go(m.TestGoWithGoBuildCache)
	polTests.Go(m.TestGoWithGoModCache)
	polTests.Go(m.TestGoWithCgoEnabled)
	polTests.Go(m.TestGoWithCgoDisabled)
	polTests.Go(m.TestGoWithGoExec)
	polTests.Go(m.TestGoWithGoInstall)
	// Test HTTP specific functions.
	polTests.Go(m.TestHTTPCurl)
	polTests.Go(m.TestHTTPDoJSONAPICall)

	// From this point onwards, we're testing the specific functionality of the ModuleTemplate module.

	if err := polTests.Wait(); err != nil {
		return WrapError(err, "there are some failed tests")
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
		return WrapError(err, "failed to validate an specific ubuntu command")
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
		return WrapError(err, "failed to get Ubuntu image")
	}

	if !strings.Contains(out, "ID=ubuntu") {
		return WrapErrorf(err, "expected Ubuntu 22.04 or ID=ubuntu, got %s", out)
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
		return WrapError(err, "failed to get Alpine image")
	}

	// Adjust the conditions to match the actual output.
	if !strings.Contains(out, "Alpine Linux") || !strings.Contains(out, "VERSION_ID=3.17.3") {
		return WrapErrorf(err, "expected Alpine Linux VERSION_ID=3.17.3, got %s", out)
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
		return WrapError(err, "failed to get BusyBox image")
	}

	if !strings.Contains(out, "v1.35.0") {
		return WrapErrorf(err, "expected BusyBox v1.35.0, got %s", out)
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
		return WrapError(err, "failed to get env vars")
	}

	if out == "" {
		return WrapError(err, "expected to have at least one env var, got empty output")
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
		ModuleTemplate()

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
		return WrapError(err, "expected to have at least one folder, got empty output")
	}

	if !strings.Contains(out, "total") {
		return WrapErrorf(err, "expected to have at least one folder, got %s", out)
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
	// Retrieve the environment variables using the ModuleTemplate's PrintEnvVars function.
	envVars, err := dag.ModuleTemplate().PrintEnvVars(ctx)

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
		return WrapErrorf(err, "failed to inspect env var %s", key)
	}

	// Check if the inspected value matches the expected result.
	if !strings.Contains(out, value) {
		return WrapErrorf(err, "expected %s to be %s, got %s", key, value, out)
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
		return WrapError(err, "failed to run shell command")
	}

	if out == "" {
		return WrapError(err, "expected to have at least one folder, got empty output")
	}

	if !strings.Contains(out, "curl") {
		return WrapErrorf(err, "expected 'curl' to be available in the container, got %s", out)
	}

	return nil
}

// TestWithUtilitiesInUbuntuContainer tests if the Alpine container with utilities is set correctly.
//
// This method verifies that the Alpine container includes specific utilities by running a command within the container.
// The test checks if the 'curl' utility is available and functioning as expected.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
// - error: Returns an error if the utility check fails or the output does not contain the expected results.
func (m *Tests) TestWithUtilitiesInUbuntuContainer(ctx context.Context) error {
	// Configure an ubuntu container.
	ubuntuCtr := dag.Container().
		From("ubuntu:latest")

	// Initialize the target module with the Alpine container and utilities.
	targetModule := dag.
		ModuleTemplate(dagger.ModuleTemplateOpts{
			Ctr: ubuntuCtr,
		})

	targetModule = targetModule.WithUtilitiesInUbuntuContainer()

	// Execute the 'curl --version' command within the container to verify 'curl' utility.
	out, err := targetModule.
		Ctr().
		WithExec([]string{"curl", "--version"}).
		Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to run shell command")
	}

	if out == "" {
		return WrapError(err, "expected to have at least one folder, got empty output")
	}

	if !strings.Contains(out, "curl") {
		return WrapErrorf(err, "expected 'curl' to be available in the container, got %s", out)
	}

	return nil
}

// TestWithGitInAlpineContainer tests the presence of the Git version control system
// within an Alpine-based container.
//
// This method configures the target module to include Git within an Alpine container,
// executes a command to check the Git version, and verifies the output.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if there is an issue configuring the container,
//     executing the command, or if the output does not contain the expected Git version information.
func (m *Tests) TestWithGitInAlpineContainer(ctx context.Context) error {
	targetModule := dag.ModuleTemplate()

	// Configure the target module to include Git in an Alpine container
	targetModule = targetModule.WithGitInAlpineContainer()

	out, err := targetModule.Ctr().
		WithExec([]string{"git", "--version"}).
		Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to run shell command")
	}

	// Check if the command output is empty
	if out == "" {
		return WrapError(err, "expected to have at least one folder, got empty output")
	}

	// Check if the Git version information is present in the output
	if !strings.Contains(out, "git version") {
		return WrapErrorf(err, "expected 'git' to be available in the container, got %s", out)
	}

	return nil
}

// TestWithGitInUbuntuContainer verifies that the 'git' command is available
// in the target module's container which uses an Ubuntu base image.
//
// This method reconfigures the target module to include 'git' in an Ubuntu container,
// executes the 'git --version' command within the container to confirm its presence,
// and checks the output to ensure 'git' is correctly installed.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if there is an issue configuring the container,
//     executing the 'git' command, or if the 'git' command is not found in the output.
func (m *Tests) TestWithGitInUbuntuContainer(ctx context.Context) error {
	// Configure an ubuntu container.
	ubuntuCtr := dag.Container().
		From("ubuntu:latest")

	// Initialize the target module.
	targetModule := dag.ModuleTemplate(dagger.ModuleTemplateOpts{
		Ctr: ubuntuCtr,
	})

	// Configure the target module to include 'git' in an Ubuntu container.
	targetModule = targetModule.WithGitInUbuntuContainer()

	// Execute the 'git --version' command to verify 'git' is installed.
	out, err := targetModule.Ctr().
		WithExec([]string{"git", "--version"}).
		Stdout(ctx)

		// Check for errors executing the command.
	if err != nil {
		return WrapError(err, "failed to run shell command")
	}

	// Check if the output is empty, which indicates the command may have failed.
	if out == "" {
		return WrapError(err, "expected to have at least one folder, got empty output")
	}

	// Check if the output contains 'git version' to confirm 'git' is available.
	if !strings.Contains(out, "git version") {
		return WrapErrorf(err, "expected 'git' to be available in the container, got %s", out)
	}

	// Return nil if 'git' was successfully verified.
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
		return WrapError(err, "failed to get netrc file")
	}

	// Check if the .netrc file contains the expected machine entry.
	if !strings.Contains(out, "machine github.com") {
		return WrapErrorf(err, "expected netrc file to be created, got %s", out)
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
		return WrapError(err, "failed to get netrc file")
	}

	// Check if the .netrc file contains the expected machine entry.
	if !strings.Contains(out, "machine github.com") {
		return WrapErrorf(err, "expected netrc file to be created, got %s", out)
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
	targetModule = targetModule.
		WithNewNetrcFileAsSecretGitLab("gitlab-username", gitlabSecret)

	// Execute a command to read the .netrc file from the container.
	out, err := targetModule.Ctr().
		WithExec([]string{"cat", "/root/.netrc"}).
		Stdout(ctx)

	// Check for errors executing the command.
	if err != nil {
		return WrapError(err, "failed to get netrc file")
	}

	// Check if the .netrc file contains the expected machine entry.
	if !strings.Contains(out, "machine gitlab.com") {
		return WrapErrorf(err, "expected netrc file to be created, got %s", out)
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
	out, err := targetModule.
		Ctr().
		WithExec([]string{"cat", "/root/.netrc"}).
		Stdout(ctx)

	// Check for errors executing the command.
	if err != nil {
		return WrapError(err, "failed to get netrc file")
	}

	// Check if the .netrc file contains the expected machine entry.
	if !strings.Contains(out, "machine gitlab.com") {
		return WrapErrorf(err, "expected netrc file to be created, got %s", out)
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
		return WrapError(err, "failed to get env vars")
	}

	// Check if the expected environment variables are set.
	if !strings.Contains(out, "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE") {
		return WrapErrorf(err, "expected AWS_ACCESS_KEY_ID to be set, got %s", out)
	}

	if !strings.Contains(out, "GCP_PROJECT_ID=example-project-id") {
		return WrapErrorf(err, "expected GCP_PROJECT_ID to be set, got %s", out)
	}

	if !strings.Contains(out, "ANOTHER_SECRET=another-secret-value") {
		return WrapErrorf(err, "expected ANOTHER_SECRET to be set, got %s", out)
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
		return WrapErrorf(err, "failed to get download file from url %s", fileURL)
	}

	if out == "" {
		return WrapError(err, "expected to have at least one folder, got empty output")
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
		return WrapErrorf(err, "failed to get download file from url %s", fileURL)
	}

	if out == "" {
		return WrapError(err, "expected to have at least one folder, got empty output")
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
		return WrapError(err, "failed to get ls output")
	}

	if out == "" {
		return WrapError(err, "expected to have at least one folder, got empty output")
	}

	if !strings.Contains(out, "total") {
		return WrapErrorf(err, "expected to have at least one folder, got %s", out)
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
		return WrapError(err, "failed to get env vars")
	}

	// Check if the cache-busting environment variable is set
	if !strings.Contains(out, "CACHE_BUSTER") {
		return WrapErrorf(err, "expected CACHE_BUSTER to be set, got %s", out)
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
		return WrapError(err, "failed to get ls output")
	}

	if out == "" {
		return WrapError(err, "expected to have at least one folder, got empty output")
	}

	if !strings.Contains(out, "total") {
		return WrapErrorf(err, "expected to have at least one folder, got %s", out)
	}

	// Check if the .gitignore file is present.
	out, err = ctr.
		WithExec([]string{"cat", "/mnt/.gitignore"}).
		Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to get .gitignore file")
	}

	if out == "" {
		return WrapError(err, "could not inspect (cat) the .gitignore file")
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
		return WrapErrorf(err, "failed to get download file from url %s", fileURL)
	}

	// Check if the output is empty, which indicates the file was not found.
	if out == "" {
		return WrapError(err, "expected to find the file at /mnt/myfile.zip, but got empty output")
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
		return WrapError(err, "failed to get AWS Keys environment variables")
	}

	if !strings.Contains(out, "AWS_ACCESS_KEY_ID") {
		return WrapErrorf(err, "expected AWS_ACCESS_KEY_ID to be set, got %s", out)
	}

	if !strings.Contains(out, "AWS_SECRET_ACCESS_KEY") {
		return WrapErrorf(err, "expected AWS_SECRET_ACCESS_KEY to be set, got %s", out)
	}

	// Check if the content of the env vars is correct.
	if !strings.Contains(out, "AKIAIOSFODNN7EXAMPLE") {
		return WrapErrorf(err, "expected AKIAIOSFODNN7EXAMPLE to be set, got %s", out)
	}

	if !strings.Contains(out, "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY") {
		return WrapErrorf(err, "expected wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY to be set, got %s", out)
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
		return WrapError(err, "failed to get Azure Keys environment variables")
	}

	if !strings.Contains(out, "AZURE_CLIENT_ID") {
		return WrapErrorf(err, "expected AZURE_CLIENT_ID to be set, got %s", out)
	}

	if !strings.Contains(out, "AZURE_CLIENT_SECRET") {
		return WrapErrorf(err, "expected AZURE_CLIENT_SECRET to be set, got %s", out)
	}

	if !strings.Contains(out, "AZURE_TENANT_ID") {
		return WrapErrorf(err, "expected AZURE_TENANT_ID to be set, got %s", out)
	}

	return nil
}

// TestGoWithGoPlatform tests the setting of different Go platforms within the target module's container.
//
// This method creates a target module with a Golang Alpine container and sets different Go platforms.
// It verifies if the Go platform is correctly set by running the `go version` command within the container
// for each defined platform.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if there is an issue setting the Go platform, executing
//     commands in the container, or if the `go version` output does not match the expected pattern.
func (m *Tests) TestGoWithGoPlatform(ctx context.Context) error {
	platforms := map[dagger.Platform]string{
		"linux/amd64":   "go version go1.20.4 linux/amd64",
		"linux/arm64":   "go version go1.20.4 linux/arm64",
		"windows/amd64": "go version go1.20.4 windows/amd64",
	}

	for platform := range platforms {
		targetModule := dag.
			ModuleTemplate(dagger.ModuleTemplateOpts{
				Ctr: getGolangAlpineContainer(""),
			}).WithGoPlatform(dagger.ModuleTemplateWithGoPlatformOpts{
			Platform: platform,
		})

		// Check if the Go platform is set correctly.
		out, err := targetModule.Ctr().
			WithExec([]string{"go", "version"}).
			Stdout(ctx)

		if out == "" {
			return WrapErrorf(err, "failed to run go version for platform %s", platform)
		}

		if err != nil {
			return WrapErrorf(err, "failed to run go version for platform %s", platform)
		}

		// Validate the GOOS and GOARCH environment variables.
		goosOut, goosOutErr := targetModule.Ctr().
			WithExec([]string{"printenv", "GOOS"}).
			Stdout(ctx)

		if goosOutErr != nil {
			return WrapErrorf(goosOutErr, "failed to get GOOS for platform %s", platform)
		}

		goarchOut, goarchOutErr := targetModule.Ctr().
			WithExec([]string{"printenv", "GOARCH"}).
			Stdout(ctx)

		if goarchOutErr != nil {
			return WrapErrorf(goarchOutErr, "failed to get GOARCH for platform %s", platform)
		}

		platformStr := string(platform)

		expectedGOOS := strings.Split(platformStr, "/")[0]
		expectedGOARCH := strings.Split(platformStr, "/")[1]

		if !strings.Contains(goosOut, expectedGOOS) {
			return WrapErrorf(err, "expected GOOS=%s, got %s for platform %s", expectedGOOS, goosOut, platform)
		}

		if !strings.Contains(goarchOut, expectedGOARCH) {
			return WrapErrorf(err, "expected GOARCH=%s, got %s for platform %s", expectedGOARCH, goarchOut, platform)
		}
	}

	return nil
}

// TestGoWithCgoEnabled tests enabling CGO in a Go Alpine container.
//
// This function performs the following steps:
// 1. Sets up the Go Alpine container.
// 2. Enables the CGO feature in the Go environment.
// 3. Verifies that the CGO_ENABLED environment variable is set to "1".
//
// Parameters:
// - ctx: The context to control the execution.
//
// Returns:
// - error: If any step fails or produces an unexpected output, an error is returned.
func (m *Tests) TestGoWithCgoEnabled(ctx context.Context) error {
	// Setting the Go Alpine container.
	targetModule := dag.ModuleTemplate(dagger.ModuleTemplateOpts{
		Ctr: getGolangAlpineContainer(""),
	})

	// Enable CGO.
	targetModule = targetModule.
		WithGoCgoEnabled()

	out, err := targetModule.Ctr().
		WithExec([]string{"go", "env", "CGO_ENABLED"}).
		Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to get CGO_ENABLED environment variable")
	}

	if out == "" {
		return WrapError(err, "expected to have at least one folder, got empty output")
	}

	if !strings.Contains(out, "1") {
		return WrapErrorf(err, "expected CGO_ENABLED to be set to 1, got %s", out)
	}

	return nil
}

// TestGoWithCgoDisabled tests disabling CGO in a Go Alpine container.
//
// This function performs the following steps:
// 1. Sets up the Go Alpine container.
// 2. Disables the CGO feature in the Go environment.
// 3. Verifies that the CGO_ENABLED environment variable is set to "0".
//
// Parameters:
// - ctx: The context to control the execution.
//
// Returns:
// - error: If any step fails or produces an unexpected output, an error is returned.
func (m *Tests) TestGoWithCgoDisabled(ctx context.Context) error {
	// Setting the Go Alpine container.
	targetModule := dag.ModuleTemplate(dagger.ModuleTemplateOpts{
		Ctr: getGolangAlpineContainer(""),
	})

	// Disable CGO.
	targetModule = targetModule.
		WithCgoDisabled()

	out, err := targetModule.Ctr().
		WithExec([]string{"go", "env", "CGO_ENABLED"}).
		Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to get CGO_ENABLED environment variable")
	}

	if out == "" {
		return WrapError(err, "expected to have at least one folder, got empty output")
	}

	if !strings.Contains(out, "0") {
		return WrapErrorf(err, "expected CGO_ENABLED to be set to 0, got %s", out)
	}

	return nil
}

// TestGoWithGoBuildCache verifies that the Go build cache (GOCACHE) is set correctly
// in the provided Alpine container.
//
// This function performs the following steps:
// 1. Sets up the Go Alpine container without specifying a particular Go version.
// 2. Configures the Go build cache.
// 3. Executes the `go env GOCACHE` command to retrieve the GOCACHE environment variable.
// 4. Validates that the GOCACHE environment variable is set to the expected path.
//
// Parameters:
// - ctx: The context to control the execution.
//
// Returns:
// - error: If any of the steps fail, an error is returned indicating what went wrong.
func (m *Tests) TestGoWithGoBuildCache(ctx context.Context) error {
	// Setting the Go Alpine container.
	targetModule := dag.ModuleTemplate(dagger.ModuleTemplateOpts{
		Ctr: getGolangAlpineContainer(""),
	})

	// Set the Go build cache.
	targetModule = targetModule.WithGoBuildCache()
	out, err := targetModule.Ctr().
		WithExec([]string{"go", "env", "GOCACHE"}).
		Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to get GOCACHE environment variable")
	}

	if out == "" {
		return WrapError(err, "expected to have at least one folder, got empty output")
	}

	if !strings.Contains(out, "/root/.cache/go-build") {
		return WrapErrorf(err, "expected GOCACHE to be set to /root/.cache/go-build, got %s", out)
	}

	return nil
}

// TestGoWithGoModCache verifies that the Go module cache (GOMODCACHE) is set correctly
// in the provided Alpine container.
//
// This function performs the following steps:
// 1. Sets up the Go Alpine container without specifying a particular Go version.
// 2. Configures the Go module cache.
// 3. Executes the `go env GOMODCACHE` command to retrieve the GOMODCACHE environment variable.
// 4. Validates that the GOMODCACHE environment variable is set to the expected path.
//
// Parameters:
// - ctx: The context to control the execution.
//
// Returns:
// - error: If any of the steps fail, an error is returned indicating what went wrong.
func (m *Tests) TestGoWithGoModCache(ctx context.Context) error {
	// Setting the Go Alpine container.
	targetModule := dag.ModuleTemplate(dagger.ModuleTemplateOpts{
		Ctr: getGolangAlpineContainer(""),
	})

	// Set the Go mod cache.
	targetModule = targetModule.
		WithGoModCache()

	out, err := targetModule.Ctr().
		WithExec([]string{"go", "env", "GOMODCACHE"}).
		Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to get GOMODCACHE environment variable")
	}

	if out == "" {
		return WrapError(err, "expected to have at least one folder, got empty output")
	}

	if !strings.Contains(out, "/go/pkg/mod") {
		return WrapErrorf(err, "expected GOMODCACHE to be set to /go/pkg/mod, got %s", out)
	}

	return nil
}

// TestGoWithGoInstall tests the installation of various Go packages
// in the provided Alpine container.
//
// This function performs the following steps:
// 1. Sets up the Go Alpine container with the expected Go version.
// 2. Installs a list of specified Go packages.
// 3. Verifies the installation by checking if the installed packages are in the PATH.
// 4. Ensures the Go module cache is correctly set.
//
// Parameters:
// - ctx: The context to control the execution.
//
// Returns:
// - error: If any of the installation or verification steps fail, an error is returned.
//
// longer function.
func (m *Tests) TestGoWithGoInstall(ctx context.Context) error {
	// Setting the Go Alpine container.
	targetModule := dag.ModuleTemplate(dagger.ModuleTemplateOpts{
		Ctr: getGolangAlpineContainer("1.22.3"),
	})

	// List of packages to install
	packages := []string{
		"gotest.tools/gotestsum@latest",
		"github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1",
		"github.com/go-delve/delve/cmd/dlv@latest",
	}

	// Prepare the installation commands for each package
	targetModule = targetModule.WithGoInstall(packages)
	// Sync, to execute the installation commands
	_, err := targetModule.Ctr().Stdout(ctx)
	if err != nil {
		return WrapError(err, "failed to get the standard output of the commands to install Go packages")
	}

	// Verify installations
	for _, pkg := range packages {
		pkgNameSplit := strings.Split(pkg, "/")
		pkgName := pkgNameSplit[len(pkgNameSplit)-1]
		pkgName = strings.Split(pkgName, "@")[0]

		out, err := targetModule.Ctr().WithExec([]string{"which", pkgName}).Stdout(ctx)
		if err != nil {
			return WrapErrorf(err, "failed to verify installation of %s", pkg)
		}

		if out == "" {
			return WrapErrorf(err, "expected to find %s in PATH, got empty output", pkgName)
		}
	}

	// Verify Go module cache
	out, goCacheErr := targetModule.Ctr().WithExec([]string{"go", "env", "GOMODCACHE"}).Stdout(ctx)
	if goCacheErr != nil {
		return WrapError(goCacheErr, "failed to get GOMODCACHE environment variable")
	}

	if !strings.Contains(out, "/go/pkg/mod") {
		return Errorf("expected GOMODCACHE to be set to /go/pkg/mod, got %s", strings.TrimSpace(out))
	}

	return nil
}

// TestGoWithGoExec tests the execution of various Go commands and ensures they produce the expected
// results in the provided Alpine container.
//
// This function performs the following steps:
// 1. Sets up the Go Alpine container with the expected Go version.
// 2. Executes the `go version` command and verifies the output against the expected Go version.
// 3. Runs additional Go commands (e.g., `go env GOPATH`) and checks their
// output against expected values.
// 4. Verifies specific Go environment variables (e.g., `GOPROXY`) to ensure
// they are set correctly.
//
// Parameters:
// - ctx: The context to control the execution.
//
// Returns:
// - error: If any of the commands fail or produce unexpected output, an error is returned.
//
//nolint:cyclop // The test handles multiple commands and environments, requiring a longer function.
func (m *Tests) TestGoWithGoExec(ctx context.Context) error {
	// Setting the Go Alpine container.
	expectedGoVersion := "1.22.3"
	targetModule := dag.ModuleTemplate(dagger.ModuleTemplateOpts{
		Ctr: getGolangAlpineContainer(expectedGoVersion),
	})

	// Execute `go version` command and capture the output.
	targetModule = targetModule.
		WithGoExec([]string{"version"})

	out, err := targetModule.Ctr().Stdout(ctx)
	if err != nil {
		return WrapErrorf(err, "failed to get the standard output of the "+
			"version command with expected Go version %s", expectedGoVersion)
	}

	if out == "" {
		return Errorf("expected to have Go version %s, got empty output", expectedGoVersion)
	}

	// Verify the Go version.
	expectedVersionString := "go" + expectedGoVersion
	if !strings.Contains(out, expectedVersionString) {
		return Errorf("expected Go version %s, but got %s", expectedVersionString, out)
	}

	// Additional Go commands to verify.
	commands := map[string]string{
		"go env GOPATH":  "/go",
		"go env GOROOT":  "/usr/local/go",
		"go env GOPROXY": "https://proxy.golang.org,direct",
	}

	for cmd, expectedOutput := range commands {
		out, execErr := targetModule.Ctr().
			WithExec(strings.Split(cmd, " ")).
			Stdout(ctx)
		if execErr != nil {
			return WrapErrorf(execErr, "failed to execute command: %s", cmd)
		}

		if out == "" {
			return Errorf("expected output for command %s, got empty output", cmd)
		}

		if !strings.Contains(out, expectedOutput) {
			return Errorf("expected output for command '%s' to contain %s, but got %s", cmd, expectedOutput, out)
		}
	}

	// Verify 'go env' command to ensure environment variables are set correctly.
	envVars := map[string]string{
		"GOPATH":  "/go",
		"GOROOT":  "/usr/local/go",
		"GOPROXY": "https://proxy.golang.org,direct",
	}

	for envVar, expectedVal := range envVars {
		out, envErr := targetModule.Ctr().
			WithExec([]string{"go", "env", envVar}).
			Stdout(ctx)
		if envErr != nil {
			return WrapErrorf(envErr, "failed to get %s environment variable", envVar)
		}

		if out == "" {
			return Errorf("expected %s environment variable to be set, got empty output", envVar)
		}

		if !strings.Contains(out, expectedVal) {
			return Errorf("expected %s to be set to %s, but got %s", envVar, expectedVal, out)
		}
	}

	return nil
}

// TestGoWithGoBuild tests the Go build process using the provided Alpine container.
//
// This function performs the following steps:
//
// 1. Sets up the Go Alpine container with the expected Go version.
// 2. Configures the build process with specific options including the source
// directory, target platform, package to build, and output binary name.
// 3. Executes the build process and checks for errors.
// 4. Verifies the presence of the output binary in the container's directory.
// 5. Runs the binary and verifies the output against the expected string.
//
// Parameters:
// - ctx: The context to control the execution.
func (m *Tests) TestGoWithGoBuild(ctx context.Context) error {
	// Setting the Go Alpine container.
	expectedGoVersion := "1.22.3"
	targetModule := dag.ModuleTemplate(dagger.ModuleTemplateOpts{
		Ctr: getGolangAlpineContainer(expectedGoVersion),
	})

	// Configure the build process with specific options.
	targetModule = targetModule.
		WithSource(m.TestDir).
		WithGoBuild(dagger.ModuleTemplateWithGoBuildOpts{
			Platform: "linux/amd64",
			Pkg:      "golang/main.go",
			Verbose:  true,
			// Output the binary called dagger to the current directory.
			Output: "./dagger",
		})

	// Execute the build process and capture the output.
	_, err := targetModule.Ctr().Stdout(ctx)

	if err != nil {
		return WrapErrorf(err, "failed to get the standard output of "+
			"the version command with expected Go version %s", expectedGoVersion)
	}

	// Verify the presence of files in the directory.
	lsOut, lsErr := targetModule.Ctr().WithExec([]string{"ls"}).Stdout(ctx)
	if lsErr != nil {
		return WrapErrorf(lsErr, "failed to get the standard output of the "+
			"version command with expected Go version %s", expectedGoVersion)
	}

	if lsOut == "" {
		return Errorf("expected to have files listed, got empty output")
	}

	// Verify the presence of the output binary named 'dagger'.
	daggerOut, daggerErr := targetModule.Ctr().WithExec([]string{"ls", "dagger"}).Stdout(ctx)
	if daggerErr != nil {
		return WrapErrorf(daggerErr, "failed to get the standard output of the version "+
			"command with expected Go version %s", expectedGoVersion)
	}

	if daggerOut == "" {
		return Errorf("expected to have a file called 'dagger', got empty output")
	}

	// Run the binary and verify the output.
	expectedOutput := "Hello, Dagger!"
	out, binaryErr := targetModule.
		Ctr().
		WithExec([]string{"./dagger"}).Stdout(ctx)

	if binaryErr != nil {
		return WrapErrorf(binaryErr, "failed to get the standard output of the version "+
			"command with expected Go version %s", expectedGoVersion)
	}

	if out == "" {
		return Errorf("expected to have output, got empty output")
	}

	if !strings.Contains(out, expectedOutput) {
		return Errorf("expected output to contain %s, got %s", expectedOutput, out)
	}

	return nil
}

// TestHTTPCurl tests an HTTP GET request using the curl command within an Alpine container.
//
// This function performs the following steps:
// 1. Sets up the Alpine container with necessary utilities to perform the curl operation.
// 2. Executes the curl command against the specified target URL and captures the output.
// 3. Verifies that the curl command produced non-empty output.
// 4. Checks for errors during the curl command execution.
//
// Parameters:
// - ctx: The context to control the execution.
//
// Returns:
// - error: If the curl command fails or produces an empty output, an error is returned.
func (m *Tests) TestHTTPCurl(ctx context.Context) error {
	targetURL := "https://fakestoreapiserver.reactbd.com/smart"

	// Set up the Alpine container with required utilities for HTTP operations.
	targetModule := dag.ModuleTemplate()
	targetModule = targetModule.
		WithUtilitiesInAlpineContainer().
		WithHttpcurl(targetURL)

	// Execute the curl command and capture the output.
	out, err := targetModule.Ctr().Stdout(ctx)

	// Check if the output is empty, indicating a potential issue.
	if out == "" {
		return Errorf("failed to inspect the curl output of the URL %s. Got empty output", targetURL)
	}

	// Check for any error during the curl command execution.
	if err != nil {
		return WrapErrorf(err, "failed to curl this URL %s", targetURL)
	}

	return nil
}

// ProductJSONApiTest represents the structure of the product
// information returned by the API.
//
// Fields:
// - ID: Unique identifier for each product.
// - Title: Name/title of the product.
// - IsNew: Boolean indicating if the product is new.
// - OldPrice: The previous price of the product, represented as a string.
// - Price: The current price of the product.
// - Description: A brief description of the product.
// - Category: The category to which the product belongs.
// - Image: URL to the product's image.
// - Rating: Rating of the product out of 5.
type ProductJSONApiTest struct {
	ID          int     `json:"_id"` //nolint:tagliatelle
	Title       string  `json:"title"`
	IsNew       bool    `json:"isNew"`
	OldPrice    string  `json:"oldPrice"`
	Price       float32 `json:"price"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
	Rating      int     `json:"rating"`
}

// TestHTTPDoJSONAPICall tests an HTTP GET request to fetch product information from a JSON API.
//
// This function performs the following steps:
// 1. Sends an HTTP GET request to the specified URL to fetch product information in JSON format.
// 2. Checks if the response is non-nil.
// 3. Reads the contents of the JSON response file.
// 4. Verifies that the content is not empty and that it does not contain an error message.
// 5. Unmarshals the JSON response into a slice of ProductJSONApiTest structs.
// 6. Ensures that the unmarshalling was successful and the response contains at least one product.
//
// Parameters:
// - ctx: The context to control the execution.
//
// Returns:
// - error: If any of the steps fail, an error is returned.
func (m *Tests) TestHTTPDoJSONAPICall(ctx context.Context) error {
	targetURL := "https://fakestoreapiserver.reactbd.com/products"

	targetModule := dag.ModuleTemplate()
	jsonFile := targetModule.DoJsonapicall("GET", targetURL)

	if jsonFile == nil {
		return Errorf("failed to get the JSON response from the URL %s", targetURL)
	}

	content, err := jsonFile.Contents(ctx)
	if err != nil {
		return WrapErrorf(err, "failed to get the contents of the file /response.json")
	}

	if content == "" {
		return Errorf("failed to get the contents of the file /response.json")
	}

	// Check if the response is an error message
	if strings.Contains(content, "Bad request") || strings.Contains(content, "BAD_REQUEST") {
		return Errorf("API returned an error: %s", content)
	}

	// Unmarshal the JSON content into a slice of ProductJSONApiTest structs
	var products []ProductJSONApiTest
	err = json.Unmarshal([]byte(content), &products)

	if err != nil {
		return WrapErrorf(err, "failed to unmarshal the JSON response. Raw content: %s", content)
	}

	// Ensure that the unmarshalling was successful and the response contains at least one product
	if len(products) == 0 {
		return Errorf("failed to unmarshal the JSON response or the response was empty")
	}

	return nil
}
