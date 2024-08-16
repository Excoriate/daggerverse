package main

import (
	"context"
	"strings"
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
		Gotoolbox()

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
	targetModule := dag.Gotoolbox()

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
