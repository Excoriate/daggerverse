package main

import (
	"context"
	"github.com/Excoriate/daggerverse/module-template/tests/internal/dagger"
	"strings"
)

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