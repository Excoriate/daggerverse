package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/gopkgpublisher/tests/internal/dagger"
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

	targetModule := dag.Gopkgpublisher()
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
	targetModule := dag.Gopkgpublisher()

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
		Gopkgpublisher(dagger.GopkgpublisherOpts{
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
		Gopkgpublisher().
		WithEnvironmentVariable("HOST", "localhost", dagger.GopkgpublisherWithEnvironmentVariableOpts{
			Expand: false,
		})

	targetModule = targetModule.
		WithEnvironmentVariable("PORT", "8080", dagger.GopkgpublisherWithEnvironmentVariableOpts{
			Expand: false,
		})

	targetModule = targetModule.
		WithEnvironmentVariable("USER", "me", dagger.GopkgpublisherWithEnvironmentVariableOpts{
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
		Gopkgpublisher()

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
	targetModule := dag.Gopkgpublisher()

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
		Gopkgpublisher(dagger.GopkgpublisherOpts{
			Ctr: alpineCtr,
		})

	// Set ownership of directories to the specified user and group
	targetModule = targetModule.
		WithUserAsOwnerOfDirs("me", listOfDirsToOwn,
			dagger.GopkgpublisherWithUserAsOwnerOfDirsOpts{
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
		Gopkgpublisher(dagger.GopkgpublisherOpts{
			Ctr: alpineCtr,
		})

	// Set permissions on directories
	targetModule = targetModule.
		WithUserWithPermissionsOnDirs("0755", listOfDirsToOwn,
			dagger.GopkgpublisherWithUserWithPermissionsOnDirsOpts{
				User: "me",
			})

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
