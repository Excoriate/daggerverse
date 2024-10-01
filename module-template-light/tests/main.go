// A generated module for test the ModuleTemplateLight functions
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

	"github.com/Excoriate/daggerverse/module-template-light/tests/internal/dagger"

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

// TestAll executes all tests.
//
// This is a helper method for tests, in order to execute all tests.
func (m *Tests) TestAll(ctx context.Context) error {
	maxGoroutines := 10
	polTests := pool.
		New().
		WithMaxGoroutines(maxGoroutines).
		WithErrors().
		WithFirstError().
		WithContext(ctx)

	// Test different ways to configure the base container.
	polTests.Go(m.TestContainerWithUbuntuBase)
	polTests.Go(m.TestContainerWithAlpineBase)
	polTests.Go(m.TestContainerWithBusyBoxBase)
	polTests.Go(m.TestContainerWithWolfiBase)
	polTests.Go(m.TestPassingEnvVarsInConstructor)
	polTests.Go(m.TestContainerWithApkoBaseAlpine)
	polTests.Go(m.TestContainerWithApkoBaseWolfi)
	// Test built-in commands
	polTests.Go(m.TestRunShellCMD)
	// Test API(s) usage scenarios. APIs -> With<something>
	polTests.Go(m.TestWithContainer)
	polTests.Go(m.TestWithSource)
	polTests.Go(m.TestWithEnvironmentVariable)
	polTests.Go(m.TestWithSecretAsEnvVar)
	polTests.Go(m.TestWithDownloadedFile)
	polTests.Go(m.TestWithClonedGitRepo)
	polTests.Go(m.TestWithCacheBuster)
	// Test installation APIs
	polTests.Go(m.TestWithUtilitiesInAlpineContainer)
	polTests.Go(m.TestWithUtilitiesInUbuntuContainer)
	polTests.Go(m.TestWithGitInAlpineContainer)
	polTests.Go(m.TestWithGitInUbuntuContainer)
	// Test utility functions.
	polTests.Go(m.TestDownloadFile)
	polTests.Go(m.TestCloneGitRepo)
	// Test Go specific functions.
	polTests.Go(m.TestGoWithGoPlatform)
	polTests.Go(m.TestGoWithGoBuildCache)
	polTests.Go(m.TestGoWithGoModCache)
	polTests.Go(m.TestGoWithCgoEnabled)
	polTests.Go(m.TestGoWithCgoDisabled)

	// From this point onwards, we're testing the specific functionality of the ModuleTemplateLight module.

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
		ModuleTemplateLight()

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

// TestCloneGitRepo tests the CloneGitRepo function.
func (m *Tests) TestCloneGitRepo(ctx context.Context) error {
	targetModule := dag.ModuleTemplateLight()

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
	targetModule := dag.ModuleTemplateLight()

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
