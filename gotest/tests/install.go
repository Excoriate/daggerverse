package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/gotest/tests/internal/dagger"
)

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
		Gotest()

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
		Gotest(dagger.GotestOpts{
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
	targetModule := dag.Gotest()

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
	targetModule := dag.Gotest(dagger.GotestOpts{
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
