package main

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/Excoriate/daggerverse/aws-tag-inspector/tests/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/fixtures"
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

	targetModule := dag.AwsTagInspector()
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
	targetModule := dag.AwsTagInspector()

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
		AwsTagInspector(dagger.AwsTagInspectorOpts{
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
		AwsTagInspector().
		WithEnvironmentVariable("HOST", "localhost", dagger.AwsTagInspectorWithEnvironmentVariableOpts{
			Expand: false,
		})

	targetModule = targetModule.
		WithEnvironmentVariable("PORT", "8080", dagger.AwsTagInspectorWithEnvironmentVariableOpts{
			Expand: false,
		})

	targetModule = targetModule.
		WithEnvironmentVariable("USER", "me", dagger.AwsTagInspectorWithEnvironmentVariableOpts{
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
		AwsTagInspector()

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
	targetModule := dag.AwsTagInspector().
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
	targetModule := dag.AwsTagInspector()

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
		WithDownloadedFile(fileURL, dagger.AwsTagInspectorWithDownloadedFileOpts{
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

// TestWithClonedGitRepoHTTPS tests the WithClonedGitRepoHTTPS function.
func (m *Tests) TestWithClonedGitRepoHTTPS(ctx context.Context) error {
	targetModule := dag.AwsTagInspector()

	// This is a public repository, the token isn't required.
	targetModule = targetModule.
		WithClonedGitRepoHTTPS("https://github.com/excoriate/daggerverse",
			dagger.AwsTagInspectorWithClonedGitRepoHTTPSOpts{})

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
	targetModule := dag.AwsTagInspector()

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

// TestWithConfigFile tests the setting of a configuration file within the target module's container.
//
// This method mounts a configuration file (`common/test-file.yml`) into the target module's container
// and verifies if it is correctly mounted by running the `cat` command within the container. It also
// checks if the environment variable `TEST_CONFIG_PATH` is set to the correct path of the configuration file.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if there is an issue mounting the configuration file, executing commands
//     in the container, or if the configuration file or environment variable is not found in the output.
func (m *Tests) TestWithConfigFile(ctx context.Context) error {
	configTestFilePath := "common/test-file.yml"
	configTestFilePathInCtr := filepath.Join(fixtures.MntPrefix, configTestFilePath)

	file := m.TestDir.
		File(configTestFilePath)

	targetModule := dag.
		AwsTagInspector().
		WithConfigFile(file, dagger.AwsTagInspectorWithConfigFileOpts{
			CfgPathInCtr: configTestFilePathInCtr,
			SetEnvVar:    "TEST_CONFIG_PATH",
		})

	// Inspecting the file mounted into the container.
	catOut, catOutErr := targetModule.
		Ctr().
		WithExec([]string{"cat", configTestFilePathInCtr}).
		Stdout(ctx)

	if catOutErr != nil {
		return WrapErrorf(catOutErr, "failed to cat the file %s", configTestFilePathInCtr)
	}

	if catOut == "" {
		return WrapError(catOutErr, "expected to have at least one folder, got empty output")
	}

	if !strings.Contains(catOut, "users") {
		return WrapErrorf(catOutErr, "expected to have at least one folder, got %s", catOut)
	}

	// Inspect if the environment variable was set, and it was the expected value.
	envVarOut, envVarOutErr := targetModule.
		Ctr().
		WithExec([]string{"printenv"}).
		Stdout(ctx)

	if envVarOutErr != nil {
		return WrapError(envVarOutErr, "failed to get env vars")
	}

	if !strings.Contains(envVarOut, "TEST_CONFIG_PATH") {
		return WrapErrorf(envVarOutErr, "expected TEST_CONFIG_PATH to be set, got %s", envVarOut)
	}

	if !strings.Contains(envVarOut, configTestFilePathInCtr) {
		return WrapErrorf(envVarOutErr, "expected TEST_CONFIG_PATH to be set, got %s", envVarOut)
	}

	return nil
}

// TestWithUserAsOwnerOfDirs tests if the specified user and group own the created directories.
func (m *Tests) TestWithUserAsOwnerOfDirs(ctx context.Context) error {
	// Create a new container from the alpine:latest image
	alpineCtr := dag.
		Container().
		From("alpine:latest")

	// Install shadow package, create a group and a user
	alpineCtr = alpineCtr.
		WithExec([]string{"apk", "add", "--no-cache", "shadow"}).
		WithExec([]string{"groupadd", "mygroup"}).
		WithExec([]string{"useradd", "-G", "mygroup", "me"})

	// List of directories to be created and owned by the user
	listOfDirsToOwn := []string{
		"/mnt/test-dir",
		"/mnt/test-dir-1",
		"/mnt/test-dir-2",
		"/mnt/test-dir-3",
		"/mnt/test-dir-4",
		"/mnt/test-dir-5",
	}

	// Create the directories
	for _, dir := range listOfDirsToOwn {
		alpineCtr = alpineCtr.
			WithExec([]string{"mkdir", "-p", dir})
	}

	// Create a new module with the container
	targetModule := dag.
		AwsTagInspector(dagger.AwsTagInspectorOpts{
			Ctr: alpineCtr,
		})

	// Set ownership of directories to the specified user and group
	targetModule = targetModule.
		WithUserAsOwnerOfDirs("me", listOfDirsToOwn,
			dagger.AwsTagInspectorWithUserAsOwnerOfDirsOpts{
				Group:           "mygroup",
				ConfigureAsRoot: true,
			})

	// Verify ownership of the directories
	out, err := targetModule.Ctr().
		WithExec(append([]string{"stat", "-c", "%U %G %n"}, listOfDirsToOwn...)).
		Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to get stat output")
	}

	// Check if the ownership is as expected
	for _, dir := range listOfDirsToOwn {
		expected := "me mygroup " + dir
		if !strings.Contains(out, expected) {
			return WrapErrorf(err,
				"expected user owner to be me and group owner to be mygroup for %s, got %s", dir, out)
		}
	}

	return nil
}

// TestWithUserWithPermissionsOnDirs tests if the specified permissions are set on the created directories.
func (m *Tests) TestWithUserWithPermissionsOnDirs(ctx context.Context) error {
	// Create a new container from the alpine:latest image
	alpineCtr := dag.
		Container().
		From("alpine:latest")

	// Install shadow package, create a group and a user
	alpineCtr = alpineCtr.
		WithExec([]string{"apk", "add", "--no-cache", "shadow"}).
		WithExec([]string{"groupadd", "mygroup"}).
		WithExec([]string{"useradd", "-G", "mygroup", "me"})

	// List of directories to be created and owned by the user
	listOfDirsToOwn := []string{
		"/mnt/test-dir",
		"/mnt/test-dir-1",
		"/mnt/test-dir-2",
		"/mnt/test-dir-3",
		"/mnt/test-dir-4",
		"/mnt/test-dir-5",
	}

	// Create the directories
	for _, dir := range listOfDirsToOwn {
		alpineCtr = alpineCtr.
			WithExec([]string{"mkdir", "-p", dir})
	}

	// Create a new module with the container
	targetModule := dag.
		AwsTagInspector(dagger.AwsTagInspectorOpts{
			Ctr: alpineCtr,
		})

	// Set permissions on directories
	targetModule = targetModule.
		WithUserWithPermissionsOnDirs("testuser", "0755", listOfDirsToOwn)

	// Verify permissions
	out, err := targetModule.Ctr().
		WithExec(append([]string{"stat", "-c", "%a %n"}, listOfDirsToOwn...)).
		Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to get stat output")
	}

	for _, dir := range listOfDirsToOwn {
		expected := "755 " + dir
		if !strings.Contains(out, expected) {
			return WrapErrorf(err, "expected permissions to be 755 for %s, got %s", dir, out)
		}
	}

	return nil
}
