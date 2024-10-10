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

// utilValidateVersion executes the version command for a given tool and checks
// the output to ensure it contains
// the expected version string. If the command fails or the output does not contain
// the expected string, an error is returned.
func (m *Tests) utilValidateVersion(
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

// utilValidateIfEnvVarIsSetInContainer checks if a specific environment variable is set in the container.
// It executes the "printenv" command in the container and verifies if the output contains the specified variable.
// If the command fails, the output is empty, or the variable is not found, an error is returned.
func (m *Tests) utilValidateIfEnvVarIsSetInContainer(
	ctx context.Context,
	ctr *dagger.Container,
	variable string,
) error {
	variableOut, variableOutErr := ctr.
		WithExec([]string{"printenv"}).
		Stdout(ctx)

	if variableOutErr != nil {
		return WrapErrorf(variableOutErr, "failed to get environment variables output")
	}

	if variableOut == "" {
		return Errorf("environment variables output is empty")
	}

	if !strings.Contains(variableOut, variable) {
		return Errorf("environment variable '%s' is not set", variable)
	}

	return nil
}

// utilValidateEnvVarValueInContainer checks if a specific environment variable in the container has the expected value.
// It executes the "printenv <variable>" command in the container and compares the output to the expected value.
// If the command fails, the output is empty, or the value does not match the expected value, an error is returned.
//
//nolint:unused // This function is currently unused but may be used in the future.
func (m *Tests) utilValidateEnvVarValueInContainer(
	ctx context.Context,
	ctr *dagger.Container,
	variable string,
	expectedValue string,
) error {
	variableOut, variableOutErr := ctr.
		WithExec([]string{"printenv", variable}).
		Stdout(ctx)

	if variableOutErr != nil {
		return WrapErrorf(variableOutErr, "failed to get value for environment variable '%s'", variable)
	}

	if variableOut == "" {
		return Errorf("environment variable '%s' value is empty", variable)
	}

	if strings.TrimSpace(variableOut) != expectedValue {
		return Errorf("environment variable '%s' value is expected to be '%s', but it is '%s'",
			variable, expectedValue, strings.TrimSpace(variableOut))
	}

	return nil
}

// utilTheseFoldersExistsInContainer checks if the specified folders are present in the container.
// It executes the "ls -la <folder>" command in the container for each folder in the list.
// If the command fails or the output is empty for any folder, an error is returned.
func (m *Tests) utilTheseFoldersExistsInContainer(
	ctx context.Context,
	ctr *dagger.Container,
	folders []string,
) error {
	// list all, and pwd.
	ctr = ctr.
		WithExec([]string{"pwd"}).
		WithExec([]string{"ls", "-la"})

	for _, folder := range folders {
		folderOut, folderOutErr := ctr.
			WithExec([]string{"ls", "-la", folder}).
			Stdout(ctx)

		if folderOutErr != nil {
			return WrapErrorf(folderOutErr, "failed to get folder '%s' output", folder)
		}

		if folderOut == "" {
			return Errorf("folder '%s' output is empty", folder)
		}
	}

	return nil
}

// utilTheseFilesExistsInContainer checks if the specified files are present in the container.
// It executes the "ls -la <file>" command in the container for each file in the list.
// If the command fails or the output is empty for any file, an error is returned.
// If withCat is true, it also executes the "cat <file>" command to get the file content.
func (m *Tests) utilTheseFilesExistsInContainer(
	ctx context.Context,
	ctr *dagger.Container,
	files []string,
	withCat bool,
) error {
	for _, file := range files {
		fileOut, fileOutErr := ctr.
			WithExec([]string{"ls", "-la", file}).
			Stdout(ctx)

		if fileOutErr != nil {
			return WrapErrorf(fileOutErr, "failed to get file '%s' output", file)
		}

		if fileOut == "" {
			return Errorf("file '%s' output is empty", file)
		}

		if withCat {
			fileOut, fileOutErr = ctr.
				WithExec([]string{"cat", file}).
				Stdout(ctx)

			if fileOut == "" {
				return Errorf("file '%s' content is empty", file)
			}

			if fileOutErr != nil {
				return WrapErrorf(fileOutErr, "failed to get file '%s' content", file)
			}
		}
	}

	return nil
}

// utilFileShouldContainContent checks if the specified file in the container contains the given content.
// It executes the "cat <file>" command in the container to get the file content.
// If the command fails, the output is empty, or the content is not found in the file, an error is returned.
func (m *Tests) utilFileShouldContainContent(
	ctx context.Context,
	ctr *dagger.Container,
	file string,
	content string,
) error {
	fileOut, fileOutErr := ctr.
		WithExec([]string{"cat", file}).
		Stdout(ctx)

	if fileOutErr != nil {
		return WrapErrorf(fileOutErr, "failed to get file '%s' content", file)
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

// utilCommandIsSuccessfulAndOutputContains executes the specified command in the container
// and checks if the output contains the expected content.
// If the command fails, the output is empty, or the expected content is not found in the output, an error is returned.
//
//nolint:unused // This function is currently unused but may be used in the future.
func (m *Tests) utilCommandIsSuccessfulAndOutputContains(
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
