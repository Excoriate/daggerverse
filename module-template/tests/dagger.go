package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/module-template/tests/internal/dagger"
)

// TestDaggerWithDaggerCLI tests various versions of the Dagger CLI from v0.12.0 to v0.12.4
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
	versions := []string{"v0.12.1", "v0.12.2", "v0.12.3", "v0.12.4"}

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
	daggerVersionOut, daggerVersionErr := targetModule.Ctr().WithExec([]string{"dagger", "version"}).Stdout(ctx)

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
