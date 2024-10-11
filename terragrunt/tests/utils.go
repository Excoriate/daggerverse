package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/terragrunt/tests/internal/dagger"
)

const (
	defaultTestDir = "testdata"
)

// getTestDir returns the test directory.
//
// This is a helper method for tests, in order to get the test directory which
// is located in the same directory as the test file, and normally named as "testdata".
//
//nolint:unparam // It's ok to have testDir as an empty string, as it will be replaced by the default value.
func (m *Tests) getTestDir(testDir string) *dagger.Directory {
	if testDir == "" {
		testDir = defaultTestDir
	}

	return dag.
		CurrentModule().
		Source().
		Directory(testDir)
}

// assertVersionOfBinaryInContainer executes the version command for a given tool and checks
// the output to ensure it contains the expected version string. If the command fails or the output does not contain
// the expected string, an error is returned.
func (m *Tests) assertVersionOfBinaryInContainer(
	ctx context.Context,
	ctr *dagger.Container,
	tool string,
	expectedVersion string,
	versionCmd string,
) error {
	if versionCmd == "" {
		versionCmd = "--version"
	}

	versionOut, versionOutErr := ctr.
		WithExec([]string{tool, versionCmd}).
		Stdout(ctx)

	if versionOutErr != nil {
		return WrapErrorf(versionOutErr, "failed to get %s version output", tool)
	}

	if versionOut == "" {
		return Errorf("%s version output is empty", tool)
	}

	if !strings.Contains(versionOut, expectedVersion) {
		return Errorf("%s version is expected to contain '%s', but it doesn't", tool, expectedVersion)
	}

	return nil
}

// assertEnvVarIsSetInContainer checks if a specific environment variable is set in the container.
// It executes the "printenv" command in the container and verifies if the output contains the specified variable.
// If the command fails, the output is empty, or the variable is not found, an error is returned.
func (m *Tests) assertEnvVarIsSetInContainer(
	ctx context.Context,
	ctr *dagger.Container,
	variable string,
) error {
	variableOut, variableOutErr := ctr.
		WithExec([]string{"printenv"}).
		Stdout(ctx)

	if variableOutErr != nil {
		return WrapErrorf(variableOutErr, "failed to get environment variable '%s' output", variable)
	}

	if variableOut == "" {
		return Errorf("environment variable '%s' is not set", variable)
	}

	if !strings.Contains(variableOut, variable) {
		return Errorf("environment variable '%s' is expected to contain '%s', but it doesn't", 
			variable, variableOut)
	}

	return nil
}

// assertEnvVarHasValueInContainer checks if a specific environment variable is set in the container
// and if it has the expected value. It first calls assertEnvVarIsSetInContainer to ensure the variable
// is set, then checks if the variable's value matches the expected value.
// If the variable is not set, or if its value doesn't match the expected value, an error is returned.
func (m *Tests) assertEnvVarHasValueInContainer(
	ctx context.Context,
	ctr *dagger.Container,
	variable string,
	expectedValue string,
) error {
	// First, check if the variable is set
	if err := m.assertEnvVarIsSetInContainer(ctx, ctr, variable); err != nil {
		return err
	}

	// If the variable is set, check its value
	variableOut, variableOutErr := ctr.
		WithExec([]string{"printenv", variable}).
		Stdout(ctx)

	if variableOutErr != nil {
		return WrapErrorf(variableOutErr, "failed to get value for environment variable '%s'", variable)
	}

	variableOut = strings.TrimSpace(variableOut)
	if variableOut != expectedValue {
		return Errorf("environment variable '%s' has value '%s', expected '%s'", 
			variable, variableOut, expectedValue)
	}

	return nil
}

// assertTheseFoldersExistsInContainer checks if the specified folders are present in the container.
// It executes the "ls -la <folder>" command in the container for each folder in the list.
// If the command fails or the output is empty for any folder, an error is returned.
func (m *Tests) assertTheseFoldersExistsInContainer(
	ctx context.Context,
	ctr *dagger.Container,
	folders []string,
) error {
	// list all, and pwd.
	ctr = ctr.
		WithExec([]string{"pwd"}).
		WithExec([]string{"ls", "-la"})

	for _, folder := range folders {
		// Use 'test -d' to check if the folder exists
		stdout, err := ctr.
			WithExec([]string{"test", "-d", folder}).
			Stdout(ctx)

		if err != nil {
			return WrapErrorf(err, "failed to check if folder '%s' exists", folder)
		}

		// If 'test -d' succeeds, it doesn't produce any output
		// So, an empty stdout means the folder exists
		if stdout != "" {
			return Errorf("unexpected output when checking folder '%s': %s", folder, stdout)
		}

		// Optionally, list the contents of the folder
		folderContents, err := ctr.
			WithExec([]string{"ls", "-la", folder}).
			Stdout(ctx)

		if err != nil {
			return WrapErrorf(err, "failed to list contents of folder '%s'", folder)
		}

		if folderContents == "" {
			return Errorf("folder '%s' exists but is empty", folder)
		}
	}

	return nil
}

// assertTheseFilesExistsInContainer checks if the specified files are present in the container.
// It uses 'test -f' to check file existence.
func (m *Tests) assertTheseFilesExistsInContainer(
	ctx context.Context,
	ctr *dagger.Container,
	files []string,
) error {
	for _, file := range files {
		// Check if file exists
		stdout, err := ctr.
			WithExec([]string{"test", "-f", file}).
			Stdout(ctx)

		if err != nil {
			return WrapErrorf(err, "failed to check if file '%s' exists", file)
		}

		if stdout != "" {
			return Errorf("unexpected output when checking file '%s': %s", file, stdout)
		}
	}

	return nil
}

// assertFileContentShouldContain checks if the specified file in the container contains the given content.
func (m *Tests) assertFileContentShouldContain(
	ctx context.Context,
	ctr *dagger.Container,
	file string,
	content string,
) error {
	if err := m.assertTheseFilesExistsInContainer(ctx, ctr, []string{file}); err != nil {
		return WrapErrorf(err, "file '%s' does not exist", file)
	}

	fileOut, err := ctr.
		WithExec([]string{"cat", file}).
		Stdout(ctx)

	if err != nil {
		return WrapErrorf(err, "failed to get content of file '%s'", file)
	}

	if fileOut == "" {
		return Errorf("file '%s' content is empty", file)
	}

	if !strings.Contains(fileOut, content) {
		return Errorf("file '%s' content is expected to contain '%s', but its current content is '%s'",
			file, content, fileOut)
	}

	return nil
}

// assertCommandIsSuccessfulAndOutputContains executes the specified command in the container
// and checks if the output contains the expected content.
// If the command fails, the output is empty, or the expected content is not found in the output, an error is returned.
//
//nolint:unused // This function is currently unused but may be used in the future.
func (m *Tests) assertCommandIsSuccessfulAndOutputContains(
	ctx context.Context,
	ctr *dagger.Container,
	command string,
	expectedOutput string,
) error {
	commandOut, commandOutErr := ctr.
		WithExec([]string{command}).
		Stdout(ctx)

	if commandOutErr != nil {
		return WrapErrorf(commandOutErr, "failed to execute command '%s'", command)
	}

	if commandOut == "" {
		return Errorf("command '%s' output is empty", command)
	}

	if !strings.Contains(commandOut, expectedOutput) {
		return Errorf("command '%s' output is expected to contain '%s', but it doesn't", command, expectedOutput)
	}

	return nil
}
