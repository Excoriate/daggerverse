// A generated module for test the Terragrunt functions
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
	"strings"

	"github.com/Excoriate/daggerverse/terragrunt/tests/internal/dagger"

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

	t.TestDir = t.getTestDir("")

	return t
}

// TestAll executes all tests.
//
// This is a helper method for tests, in order to execute all tests.
func (m *Tests) TestAll(ctx context.Context) error {
	polTests := pool.
		New().
		WithErrors().
		WithContext(ctx)

	// Test different ways to configure the base container.
	polTests.Go(m.TestPassingEnvVarsInConstructor)
	// Test built-in commands
	polTests.Go(m.TestRunShellCMD)
	// Test API(s) usage scenarios. APIs -> With<something>
	polTests.Go(m.TestWithContainer)
	polTests.Go(m.TestWithSource)
	polTests.Go(m.TestWithEnvironmentVariable)
	polTests.Go(m.TestWithDownloadedFile)
	polTests.Go(m.TestWithCacheBuster)
	// Test utility functions.
	// Specific Terragrunt/Container functionality tests.
	polTests.Go(m.TestContainerBaseApkoDefault)
	polTests.Go(m.TestContainerBaseWithPassedImage)
	polTests.Go(m.TestContainerBaseWithAWSClI)
	polTests.Go(m.TestContainerBaseApkoWithCustomVersions)
	polTests.Go(m.TestTerragruntContainerIsUp)
	polTests.Go(m.TestTerragruntBinariesAreInstalled)
	polTests.Go(m.TestTerragruntExecInitSimpleCommand)
	polTests.Go(m.TestTerragruntExecVersionCommand)
	polTests.Go(m.TestTerragruntExecPlanCommand)
	polTests.Go(m.TestTerragruntExecLifecycleCommands)

	if err := polTests.Wait(); err != nil {
		return WrapError(err, "there are some failed tests")
	}

	return nil
}

// TestRunShellCMD runs a shell command in the container.
//
// Tests if the shell command is executed correctly in the container.
func (m *Tests) TestRunShellCMD(ctx context.Context) error {
	targetModule := dag.
		Terragrunt()

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
