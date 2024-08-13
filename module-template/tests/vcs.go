package main

import (
	"context"
	"strings"
)

// TestWithNewNetrcFileGitHub tests the creation of a new .netrc file with GitHub credentials.
//
// This function verifies that the GitHub credentials are set correctly in the .netrc file using a secret.
// It creates a new secret with the GitHub credentials and sets them in the target module's .netrc file.
// The function then reads the .netrc file from the container and checks if it contains the expected machine entry.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if the creation of the .netrc file fails, if the file's content does
//     not match the expected result,
//     or if there is an issue with executing commands in the container.
func (m *Tests) TestWithNewNetrcFileGitHub(ctx context.Context) error {
	targetModule := dag.
		ModuleTemplate()

	// Create a new secret with the GitHub credentials.
	githubSecret := dag.SetSecret("github-username", "github-password")

	// Set the GitHub credentials as a secret in the target module's .netrc file.
	targetModule = targetModule.
		WithNewNetrcFileAsSecretGitHub("github-username", githubSecret)

	// Execute a command to read the .netrc file from the container.
	out, err := targetModule.
		Ctr().
		WithExec([]string{"cat", "/root/.netrc"}).
		Stdout(ctx)

	// Check for errors executing the command.
	if err != nil {
		return WrapError(err, "failed to get netrc file")
	}

	// Check if the .netrc file contains the expected machine entry.
	if !strings.Contains(out, "machine github.com") {
		return WrapErrorf(err, "expected netrc file to be created, got %s", out)
	}

	return nil
}

// TestWithNewNetrcFileAsSecretGitHub tests the creation of a new .netrc file with GitHub credentials.
//
// This function verifies that the GitHub credentials are set correctly in the .netrc file using a secret.
// It creates a new secret with the GitHub credentials and sets them in the target module's .netrc file.
// The function then reads the .netrc file from the container and checks if it contains the expected machine entry.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if the creation of the .netrc file fails, if the file's
//     content does not match the expected result,
//     or if there is an issue with executing commands in the container.
func (m *Tests) TestWithNewNetrcFileAsSecretGitHub(ctx context.Context) error {
	// Initialize the target module.
	targetModule := dag.ModuleTemplate()

	// Create a new secret with the GitHub credentials.
	githubSecret := dag.SetSecret("github-username", "github-password")

	// Set the GitHub credentials as a secret in the target module's .netrc file.
	targetModule = targetModule.WithNewNetrcFileAsSecretGitHub("github-username", githubSecret)

	// Execute a command to read the .netrc file from the container.
	out, err := targetModule.Ctr().WithExec([]string{"cat", "/root/.netrc"}).Stdout(ctx)

	// Check for errors executing the command.
	if err != nil {
		return WrapError(err, "failed to get netrc file")
	}

	// Check if the .netrc file contains the expected machine entry.
	if !strings.Contains(out, "machine github.com") {
		return WrapErrorf(err, "expected netrc file to be created, got %s", out)
	}

	// Return nil if the netrc file is created and contains the expected entry.
	return nil
}

// TestWithNewNetrcFileGitLab tests the creation of a new .netrc file with GitLab credentials.
//
// This function verifies that the GitLab credentials are set correctly in the .netrc file using a secret.
// It creates a new secret with the GitLab credentials and sets them in the target module's .netrc file.
// The function then reads the .netrc file from the container and checks if it contains the expected machine entry.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if the creation of the .netrc file fails, if the file's
//     content does not match the expected result,
//     or if there is an issue with executing commands in the container.
func (m *Tests) TestWithNewNetrcFileGitLab(ctx context.Context) error {
	targetModule := dag.ModuleTemplate()

	// Create a new secret with the GitLab credentials.
	gitlabSecret := dag.SetSecret("gitlab-username", "gitlab-password")

	// Set the GitLab credentials as a secret in the target module's .netrc file.
	targetModule = targetModule.
		WithNewNetrcFileAsSecretGitLab("gitlab-username", gitlabSecret)

	// Execute a command to read the .netrc file from the container.
	out, err := targetModule.Ctr().
		WithExec([]string{"cat", "/root/.netrc"}).
		Stdout(ctx)

	// Check for errors executing the command.
	if err != nil {
		return WrapError(err, "failed to get netrc file")
	}

	// Check if the .netrc file contains the expected machine entry.
	if !strings.Contains(out, "machine gitlab.com") {
		return WrapErrorf(err, "expected netrc file to be created, got %s", out)
	}

	return nil
}

// TestWithNewNetrcFileAsSecretGitLab creates a new .netrc file with GitLab credentials.
//
// This method verifies that the GitLab credentials are set correctly in the .netrc file using a secret.
// It creates a new secret with the GitLab credentials and sets them in the target module's .netrc file.
// The method then reads the .netrc file from the container and checks if it contains the expected machine entry.
//
// Arguments:
// - ctx (context.Context): The context for the method execution.
//
// Returns:
//   - error: Returns an error if the creation of the .netrc file fails, if the file's
//     content does not match the expected result,
//     or if there is an issue with executing commands in the container.
func (m *Tests) TestWithNewNetrcFileAsSecretGitLab(ctx context.Context) error {
	targetModule := dag.ModuleTemplate()

	// Create a new secret with the GitLab credentials.
	gitlabSecret := dag.SetSecret("gitlab-username", "gitlab-password")

	// Set the GitLab credentials as a secret in the target module's .netrc file.
	targetModule = targetModule.WithNewNetrcFileAsSecretGitLab("gitlab-username", gitlabSecret)

	// Execute a command to read the .netrc file from the container.
	out, err := targetModule.
		Ctr().
		WithExec([]string{"cat", "/root/.netrc"}).
		Stdout(ctx)

	// Check for errors executing the command.
	if err != nil {
		return WrapError(err, "failed to get netrc file")
	}

	// Check if the .netrc file contains the expected machine entry.
	if !strings.Contains(out, "machine gitlab.com") {
		return WrapErrorf(err, "expected netrc file to be created, got %s", out)
	}

	return nil
}
