package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/terragrunt/tests/internal/dagger"
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

	targetModule := dag.Terragrunt()
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
	targetModule := dag.Terragrunt()

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
		Terragrunt(dagger.TerragruntOpts{
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
		Terragrunt().
		WithEnvironmentVariable("HOST", "localhost", dagger.TerragruntWithEnvironmentVariableOpts{
			Expand: false,
		})

	targetModule = targetModule.
		WithEnvironmentVariable("PORT", "8080", dagger.TerragruntWithEnvironmentVariableOpts{
			Expand: false,
		})

	targetModule = targetModule.
		WithEnvironmentVariable("USER", "me", dagger.TerragruntWithEnvironmentVariableOpts{
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
		Terragrunt()

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
	targetModule := dag.Terragrunt()

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
		WithDownloadedFile(fileURL, dagger.TerragruntWithDownloadedFileOpts{
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
	targetModule := dag.Terragrunt()

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
