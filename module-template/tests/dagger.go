package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/module-template/tests/internal/dagger"
)

// TestDaggerWithDaggerCLI tests various versions of the Dagger CLI from v0.13.6 to v0.13.7
//
// This function uses the provided context to run a series of tests that validate the Dagger CLI.
// It verifies that the CLI version reported matches the expected version for each specific version tested.
//
// ctx: The context for managing timeout and cancelation.
// Returns an error if any of the tests fail.
//
// Usage:
// err := m.TestDaggerWithDaggerCLI(ctx)
//
//	if err != nil {
//	  log.Fatalf("Test failed with error: %v", err)
//	}
func (m *Tests) TestDaggerWithDaggerCLI(ctx context.Context) error {
	versions := []string{"v0.13.6", "v0.13.7"}

	for _, version := range versions {
		if err := m.testDaggerVersion(ctx, version); err != nil {
			return err
		}
	}

	return nil
}

// testDaggerVersion verifies a specific version of the Dagger CLI
//
// version: The version of the Dagger CLI to test.
// Returns an error if the Dagger CLI version reported does not match the expected version.
func (m *Tests) testDaggerVersion(ctx context.Context, version string) error {
	targetModule := dag.ModuleTemplate()

	// Set the Dagger CLI version in Alpine
	targetModule = targetModule.WithDaggerClialpine(version)

	// Run the 'dagger version' command
	daggerVersionOut, daggerVersionErr := targetModule.
		Ctr().
		WithExec([]string{"dagger", "version"}).Stdout(ctx)

	// Check for errors
	if daggerVersionErr != nil {
		return WrapError(daggerVersionErr, "failed to get Dagger version for "+version)
	}

	if daggerVersionOut == "" {
		return Errorf("expected to have dagger version output, got empty output for %s", version)
	}

	expectedVersionContains := "dagger " + version

	if !strings.Contains(daggerVersionOut, expectedVersionContains) {
		return Errorf("expected Dagger version to contain %s, got %s", expectedVersionContains, daggerVersionOut)
	}

	// Set the Dagger CLI version in Ubuntu
	targetModuleUbuntu := dag.ModuleTemplate(dagger.ModuleTemplateOpts{
		Ctr: dag.Container().From("ubuntu:latest"),
	})

	targetModuleUbuntu = targetModuleUbuntu.WithDaggerCliubuntu(version)

	// Run the 'dagger version' command
	daggerVersionOutUbuntu, daggerVersionErrUbuntu := targetModuleUbuntu.Ctr().
		WithExec([]string{"dagger", "version"}).
		Stdout(ctx)

	// Check for errors
	if daggerVersionErrUbuntu != nil {
		return WrapError(daggerVersionErrUbuntu, "failed to get Dagger version for "+version)
	}

	if daggerVersionOutUbuntu == "" {
		return Errorf("expected to have dagger version output, got empty output for %s", version)
	}

	expectedVersionContainsUbuntu := "dagger " + version

	if !strings.Contains(daggerVersionOutUbuntu, expectedVersionContainsUbuntu) {
		return Errorf("expected Dagger version to contain %s, got %s", expectedVersionContainsUbuntu, daggerVersionOutUbuntu)
	}

	return nil
}

// TestDaggerSetupDaggerInDagger tests the setup of Dagger within another Dagger environment.
//
// This function performs a series of actions to ensure that Dagger can be correctly installed and run
// inside a container that is managed by another instance of Dagger. It checks the installation of the Dagger CLI,
// the Docker service, and the initialization of a Dagger module.
//
// ctx: The context for managing timeout and cancellation.
// Returns an error if any of the setup checks fail.
//
// Usage:
// err := m.TestDaggerSetupDaggerInDagger(ctx)
//
//	if err != nil {
//	  log.Fatalf("Test failed with error: %v", err)
//	}
//
//nolint:cyclop // The test handles multiple commands and environments, requiring a longer function.
func (m *Tests) TestDaggerSetupDaggerInDagger(ctx context.Context) error {
	// Initialize the target module.
	targetModule := dag.ModuleTemplate()

	// Define versions for Dagger and Docker.
	dagVersion := "v0.13.7"
	dockerVersion := "24.0"

	// Setup Dagger in Dagger environment.
	targetModule = targetModule.SetupDaggerInDagger(dagVersion, dagger.ModuleTemplateSetupDaggerInDaggerOpts{
		DockerVersion: dockerVersion,
	})

	// Verify that the Dagger CLI is installed and in the system PATH.
	daggerVersionOut, daggerVersionErr := targetModule.Ctr().
		WithExec([]string{"dagger", "version"}).
		Stdout(ctx)

	if daggerVersionErr != nil {
		return WrapError(daggerVersionErr, "failed to get Dagger version")
	}

	if daggerVersionOut == "" {
		return Errorf("expected to have dagger version output with this version %s, got empty output", dagVersion)
	}

	// Verify that the Dagger binary is set as the entry point.
	expectedEntryPoint := "/bin/dagger"
	daggerEntryPoint, daggerEntryPointErr := targetModule.Ctr().
		WithExec([]string{"which", "dagger"}).
		Stdout(ctx)

	if daggerEntryPointErr != nil {
		return WrapErrorf(daggerEntryPointErr, "failed to get Dagger entry point, expected %s", expectedEntryPoint)
	}

	if !strings.Contains(daggerEntryPoint, expectedEntryPoint) {
		return Errorf("expected to have Dagger entry point %s, got %s", expectedEntryPoint, daggerEntryPoint)
	}

	// Verify that the Dagger container is running by checking the output of `docker ps -a`.
	dockerPsCmd := []string{"docker", "ps", "-a"}
	dockerPsOut, dockerPsErr := targetModule.Ctr().WithExec(dockerPsCmd).Stdout(ctx)

	if dockerPsErr != nil {
		return WrapError(dockerPsErr, "failed to validate the docker ps command")
	}

	if dockerPsOut == "" {
		return Errorf("expected to have docker ps output, got empty output")
	}

	// Verify that the Docker service is running by executing a simple container run.
	dockerRunCmd := []string{"docker", "run", "--rm", "alpine", "echo", "Docker is working!"}
	dockerRunOut, dockerRunErr := targetModule.Ctr().WithExec(dockerRunCmd).Stdout(ctx)

	if dockerRunErr != nil {
		return WrapError(dockerRunErr, "failed to validate the "+
			"Docker service by running 'docker run --rm "+
			"alpine echo \"Docker is working!\"'")
	}

	if dockerRunOut == "" {
		return Errorf("expected to have docker run output, got empty output")
	}

	// Create and initialize a Dagger module inside the Docker container.
	daggerInitOut, daggerInitErr := targetModule.Ctr().
		WithExec([]string{"mkdir", "-p", "test-module"}).
		WithExec([]string{"sh", "-c", "cd test-module && dagger init --sdk go --name test-module"}).
		WithExec([]string{"sh", "-c", "cd test-module && dagger develop"}).
		WithExec([]string{"cat", "test-module/dagger.json"}).
		Stdout(ctx)

	if daggerInitErr != nil {
		return WrapError(daggerInitErr, "failed to validate the dagger init command")
	}

	if daggerInitOut == "" {
		return Errorf("expected to have dagger init output, got empty output")
	}

	return nil
}
