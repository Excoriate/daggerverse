package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/gotest/tests/internal/dagger"
)

// TestGoTestReturningCtr runs the Go test command in the specified test directory.
// It retrieves the standard output of the test run and returns an error if
// there is a failure in obtaining the output.
//
// Parameters:
//
//	ctx - The context for managing cancellation and deadlines.
//
// Returns:
//
//	An error if the test output cannot be retrieved; otherwise, it returns nil.
func (m *Tests) TestGoTestReturningCtr(ctx context.Context) error {
	testDir := m.getTestDir("testdata/golang")

	dagModule := dag.Gotest(
		dagger.GotestOpts{
			Image: "golang:1.22.5",
			EnvVarsFromHost: []string{
				"GO_TEST_ENV_VAR=test",
				"ANOTHER_ENV_VAR=value",
				"YET_ANOTHER_ENV_VAR=another_value",
				"GOLANG_VERSION=1.22.5", // Ensure GOLANG_VERSION is set
			},
		},
	)

	// Running the goTest container with default options.
	goTestCtr := dagModule.
		RunTest(testDir, dagger.GotestRunTestOpts{
			EnableDefaultOptions: true,
		})

	exitCode, goTestErr := goTestCtr.ExitCode(ctx)

	if exitCode != 0 {
		return Errorf("tests failed with exit code %d", exitCode)
	}

	if goTestErr != nil {
		return WrapError(goTestErr, "failed to get tests output")
	}

	cmdOut, cmdOutErr := goTestCtr.Stdout(ctx)

	if cmdOutErr != nil {
		return WrapError(cmdOutErr, "failed to get tests output")
	}

	if cmdOut == "" {
		return NewError("failed to get the Go test output, got empty string")
	}

	// Validate environment variables
	expectedEnvVars := []string{
		"GO_TEST_ENV_VAR=test",
		"ANOTHER_ENV_VAR=value",
		"YET_ANOTHER_ENV_VAR=another_value",
		"GOLANG_VERSION=1.22.5",
	}

	for _, envVar := range expectedEnvVars {
		printEnvOut, printEnvErr := goTestCtr.
			WithExec([]string{"printenv"}).
			Stdout(ctx)

		if printEnvErr != nil {
			return WrapError(printEnvErr, "failed to get environment variables")
		}

		if !strings.Contains(printEnvOut, envVar) {
			return NewError("missing expected environment variable: " + envVar)
		}
	}

	return nil
}

// TestGoTestWithCustomOptions tests various configurations of the Go test command
// using different GotestOpts options.
//
// Parameters:
//
//	ctx - The context for managing cancellation and deadlines.
//
// Returns:
//
//	An error if any test configuration fails; otherwise, it returns nil.
//
//nolint:cyclop // It's okay to have this size, it's by design.
func (m *Tests) TestGoTestWithCustomOptions(ctx context.Context) error {
	testDir := m.getTestDir("testdata/golang")

	// Test with custom version
	dagModuleWithVersion := dag.Gotest(
		dagger.GotestOpts{
			Version: "v1.22.0",
			EnvVarsFromHost: []string{
				"GOLANG_VERSION=1.22.5",
			},
		},
	)

	versionCtr := dagModuleWithVersion.
		RunTest(testDir, dagger.GotestRunTestOpts{
			EnableDefaultOptions: true,
		})

	if exitCode, err := versionCtr.ExitCode(ctx); err != nil {
		return WrapError(err, "failed to get exit code for version test")
	} else if exitCode != 0 {
		return Errorf("version test failed with exit code %d", exitCode)
	}

	// Test with custom base container
	customBaseCtr := dag.Container().
		From("golang:1.22.5-alpine").
		WithEnvVariable("CGO_ENABLED", "0")

	dagModuleWithCtr := dag.Gotest(
		dagger.GotestOpts{
			Ctr: customBaseCtr,
			EnvVarsFromHost: []string{
				"GOLANG_VERSION=1.22.5",
			},
		},
	)

	customCtr := dagModuleWithCtr.
		RunTest(testDir, dagger.GotestRunTestOpts{
			EnableDefaultOptions: true,
		})

	if exitCode, err := customCtr.ExitCode(ctx); err != nil {
		return WrapError(err, "failed to get exit code for custom container test")
	} else if exitCode != 0 {
		return Errorf("custom container test failed with exit code %d", exitCode)
	}

	// Verify custom container environment
	envOut, err := customCtr.WithExec([]string{"go", "version"}).Stdout(ctx)
	if err != nil {
		return WrapError(err, "failed to get Go version from custom container")
	}

	if envOut == "" {
		return NewError("failed to get Go version output, got empty string")
	}

	// Test with custom image
	dagModuleWithImage := dag.Gotest(
		dagger.GotestOpts{
			Image: "golang:1.22.5-bullseye",
			EnvVarsFromHost: []string{
				"GOLANG_VERSION=1.22.5",
				"TEST_ENV=custom_image",
			},
		},
	)

	imageCtr := dagModuleWithImage.
		RunTest(testDir, dagger.GotestRunTestOpts{
			EnableDefaultOptions: true,
		})

	if exitCode, err := imageCtr.ExitCode(ctx); err != nil {
		return WrapError(err, "failed to get exit code for custom image test")
	} else if exitCode != 0 {
		return Errorf("custom image test failed with exit code %d", exitCode)
	}

	// Verify environment variables in custom image
	expectedEnvVars := []string{
		"GOLANG_VERSION=1.22.5",
		"TEST_ENV=custom_image",
	}

	envOutImage, err := imageCtr.WithExec([]string{"printenv"}).Stdout(ctx)
	if err != nil {
		return WrapError(err, "failed to get environment variables from custom image")
	}

	for _, envVar := range expectedEnvVars {
		if !strings.Contains(envOutImage, envVar) {
			return Errorf("missing expected environment variable in custom image: %s", envVar)
		}
	}

	return nil
}
