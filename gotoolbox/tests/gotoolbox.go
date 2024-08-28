package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/gotoolbox/tests/internal/dagger"
)

// TestgotoolboxWithGoVersions tests various Go versions using gotoolbox.
//
// It iterates over a list of Go versions, setting up a toolbox environment
// for each version, and then verifies that the correct Go version is set up
// and reported by the environment.
//
// ctx: The context for the test execution, to control cancellation and deadlines.
//
// Returns an error if any Go version setup or verification fails.
func (m *Tests) TestgotoolboxWithGoVersions(ctx context.Context) error {
	goVersions := map[string]string{
		"1.22.6": "1.22.6-alpine3.19",
		"1.23.0": "1.23.0-alpine3.19",
		"1.21.6": "1.21.6-alpine3.19",
	}
	for expectedVersion, imageVersion := range goVersions {
		// Initialize the Go toolbox with the specified version.
		targetModDefault := dag.
			Gotoolbox(dagger.GotoolboxOpts{
				Version: imageVersion,
			})

		// Check if the Go version is set correctly.
		goVersionOut, goVersionErr := targetModDefault.
			Ctr().
			WithExec([]string{"go", "version"}).
			Stdout(ctx)

		if goVersionErr != nil {
			return WrapErrorf(goVersionErr, "failed to get Go version for %s", expectedVersion)
		}

		if goVersionOut == "" {
			return WrapErrorf(goVersionErr, "expected to have Go version "+
				"output, got empty output for %s", expectedVersion)
		}

		// Verify the output contains the expected Go version.
		if !strings.Contains(goVersionOut, expectedVersion) {
			return WrapErrorf(goVersionErr, "expected Go version %s, got %s", expectedVersion,
				goVersionOut)
		}
	}

	return nil
}

// TestgotoolboxWithOverridingContainer tests the installation of Go on an Ubuntu container.
//
// This function overrides the default container with Ubuntu, installs Go,
// and verifies that Go is correctly installed and functional.
//
// ctx: The context for the test execution, to control cancellation and deadlines.
//
// Returns an error if the Go installation or verification fails.
func (m *Tests) TestgotoolboxWithOverridingContainer(ctx context.Context) error {
	// Initialize the Go toolbox with the specified version.
	targetModDefault := dag.
		Gotoolbox(dagger.GotoolboxOpts{
			Ctr: dag.Container().From("ubuntu:22.04"),
		})

	// Installing Go on Ubuntu
	installedContainer := targetModDefault.
		Ctr().
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "golang-go"})

	// Check if Go is installed correctly
	goVersionOut, goVersionErr := installedContainer.
		WithExec([]string{"/usr/bin/go", "version"}).
		Stdout(ctx)

	if goVersionErr != nil {
		return WrapErrorf(goVersionErr, "failed to get Go version for Ubuntu")
	}

	if goVersionOut == "" {
		return NewError("expected to have Go version output, got empty output for Ubuntu")
	}

	// We're not checking for a specific version, just that Go is installed and working
	if !strings.Contains(goVersionOut, "go version go") {
		return WrapErrorf(goVersionErr, "unexpected Go version output: %s", goVersionOut)
	}

	return nil
}

// TestgotoolboxWithGoTest tests the GoTestSum functionality with various configurations.
//
// This function sets up different test cases for GoTestSum, verifies its installation,
// and runs tests using GoTestSum with different options.
//
// ctx: The context for the test execution, to control cancellation and deadlines.
//
// Returns an error if any of the GoTestSum setups or test executions fail.
func (m *Tests) TestgotoolboxWithGoTest(ctx context.Context) error {
	// Initialize the Go toolbox with the specified version.
	baseModule := dag.
		Gotoolbox(dagger.GotoolboxOpts{
			Version: "1.23.0-alpine3.19",
		}).WithSource(m.TestDir, dagger.GotoolboxWithSourceOpts{
		Workdir: "/gotoolbox",
	})

	// Test cases
	testCases := []struct {
		name             string
		goTestSumVersion string
		tParseVersion    string
		skipTParse       bool
	}{
		{"Default versions", "", "", false},
		{"Specific GoTestSum version", "v1.10.0", "", false},
		{"Specific versions for both", "v1.10.0", "v0.11.0", false},
		{"Skip TParse", "v1.10.0", "", true},
	}

	for _, testCase := range testCases {
		// Use WithGoTestSum with the test case parameters
		targetMod := baseModule.
			// WithGoTestSum(testCase.goTestSumVersion, testCase.tParseVersion, testCase.skipTParse)
			WithGoTestSum(dagger.GotoolboxWithGoTestSumOpts{
				GoTestSumVersion: testCase.goTestSumVersion,
				TParseVersion:    testCase.tParseVersion,
				SkipTparse:       testCase.skipTParse,
			})

		// Check gotestsum version
		gotestsumVersionOut, gotestsumVersionErr := targetMod.
			Ctr().
			WithExec([]string{"gotestsum", "--version"}).
			Stdout(ctx)

		if gotestsumVersionErr != nil {
			return WrapError(gotestsumVersionErr, testCase.name+": failed to get gotestsum version")
		}

		if gotestsumVersionOut == "" {
			return WrapError(gotestsumVersionErr, testCase.name+": expected to have gotestsum version output, got empty output")
		}

		// Check tparse version if not skipped
		if !testCase.skipTParse {
			tparseVersionOut, tparseVersionErr := targetMod.
				Ctr().
				WithExec([]string{"tparse", "--version"}).
				Stdout(ctx)

			if tparseVersionErr != nil {
				return WrapError(tparseVersionErr, testCase.name+": failed to get tparse version")
			}

			if tparseVersionOut == "" {
				return WrapError(tparseVersionErr, testCase.name+": expected to have tparse version output, got empty output")
			}
		}

		// Run tests with GoTestSum
		goTestSumOut, goTestSumErr := targetMod.
			Ctr().
			WithExec([]string{"gotestsum", "--format", "testname"}).
			Stdout(ctx)

		if goTestSumErr != nil {
			return WrapError(goTestSumErr, testCase.name+": failed to run gotestsum command")
		}

		if goTestSumOut == "" {
			return WrapError(goTestSumErr, testCase.name+": expected to have gotestsum output, got empty output")
		}
	}

	return nil
}

// TestgotoolboxCI is an end-to-end integration test for
// the Go toolbox in a CI environment.
//
// This function sets up the toolbox with a specified Go
// version and runs a series of Go command-line tools:
// go fmt, go vet, go test, and go build.
//
// ctx: The context for the test execution, to control cancellation and deadlines.
//
// Returns an error if any of the Go command executions fail.
func (m *Tests) TestgotoolboxCI(ctx context.Context) error {
	// Initialize the Go toolbox with the specified version.
	targetModDefault := dag.
		Gotoolbox(dagger.GotoolboxOpts{
			Version: "1.23.0-alpine3.19",
		}).WithSource(m.TestDir, dagger.GotoolboxWithSourceOpts{
		Workdir: "/gotoolbox",
	})

	// Run Go fmt to format Go source code
	ciModCfg := targetModDefault.
		WithGoExec([]string{"fmt", "./..."}).
		WithGoExec([]string{"vet", "./..."}).
		WithGoExec([]string{"test", "-v", "./..."}).
		WithGoExec([]string{"build", "-o", "/tmp/gotoolbox", "./..."})

	// Run Go command-line tools
	_, ciErr := ciModCfg.
		Ctr().
		Stdout(ctx)

	if ciErr != nil {
		return WrapError(ciErr, "failed to run Go command-line tools with WithGoExec API")
	}

	binExecOut, binExecErr := ciModCfg.
		Ctr().
		WithExec([]string{"ls", "/tmp"}).
		WithExec([]string{"sh", "-c", "cd /tmp && ./gotoolbox"}).
		Stdout(ctx)

	if binExecErr != nil {
		return WrapError(binExecErr, "failed to run ls and inspect the built Go toolbox")
	}

	if binExecOut == "" {
		return NewError("expected to have built Go toolbox output (/tmp/gotoolbox), got empty output")
	}

	if !strings.Contains(binExecOut, "Hello, Dagger") {
		return NewError("expected to have built Go toolbox output Hello, Dagger!, got " + binExecOut)
	}

	return nil
}

// TestgotoolboxWithGoReleaserAndGolangCILint tests the installation and setup
// of GoReleaser and GoLangCILint using gotoolbox.
//
// This function sets up the Go toolbox with a specified Go version, installs
// GoReleaser and GoLangCILint, and verifies their installation.
//
// ctx: The context for the test execution, to control cancellation and deadlines.
//
// Returns an error if the installation or verification of GoReleaser or GoLangCILint fails.
//
//nolint:cyclop // The test handles multiple commands and environments, requiring a longer function.
func (m *Tests) TestgotoolboxWithGoReleaserAndGolangCILint(ctx context.Context) error {
	// Initialize the Go toolbox with the specified version.
	targetModDefault := dag.
		Gotoolbox(dagger.GotoolboxOpts{
			Version: "1.23.0-alpine3.19",
		}).WithSource(m.TestDir, dagger.GotoolboxWithSourceOpts{
		Workdir: "/gotoolbox",
	}).WithGoReleaser().
		WithGoLint("v1.60.1")

		// Execute the container to install the tools.
	_, ctrErr := targetModDefault.
		Ctr().
		Stdout(ctx)

	if ctrErr != nil {
		return WrapError(ctrErr, "failed to install GoReleaser and GoLangCILint")
	}

	// Check golangci-lint binary installed in path.
	golangciPathOut, golangciPathErr := targetModDefault.
		Ctr().
		WithExec([]string{"which", "golangci-lint"}).
		Stdout(ctx)

	if golangciPathErr != nil {
		return WrapError(golangciPathErr, "failed to get golangci-lint path")
	}

	if golangciPathOut == "" {
		return Errorf("expected to have golangci-lint in path /go/bin/golangci-lint, got empty output")
	}

	if !strings.Contains(golangciPathOut, "/go/bin/golangci-lint") {
		return Errorf("expected to have golangci-lint "+
			"in path /go/bin/golangci-lint, got %s", golangciPathOut)
	}

	// Run golangci-lint --version
	golangciVersionOut, golangciVersionErr := targetModDefault.
		Ctr().
		WithExec([]string{"golangci-lint", "--version"}).
		Stdout(ctx)

	if golangciVersionErr != nil {
		return WrapError(golangciVersionErr, "failed to get golangci-lint version")
	}

	if golangciVersionOut == "" {
		return Errorf("expected to have golangci-lint version output, got empty output")
	}

	// Check goreleaser binary installed in path.
	goreleaserPathOut, goreleaserPathErr := targetModDefault.
		Ctr().
		WithExec([]string{"which", "goreleaser"}).
		Stdout(ctx)

	if goreleaserPathErr != nil {
		return WrapError(goreleaserPathErr, "failed to get goreleaser path")
	}

	if goreleaserPathOut == "" {
		return Errorf("expected to have goreleaser in path /go/bin/goreleaser, got empty output")
	}

	if !strings.Contains(goreleaserPathOut, "/go/bin/goreleaser") {
		return Errorf("expected to have goreleaser "+
			"in path /go/bin/goreleaser, got %s", goreleaserPathOut)
	}

	// Run GoReleaser --version
	goreleaserVersionOut, goreleaserVersionErr := targetModDefault.
		Ctr().
		WithExec([]string{"goreleaser", "--version"}).
		Stdout(ctx)

	if goreleaserVersionErr != nil {
		return WrapError(goreleaserVersionErr, "failed to get goreleaser version")
	}

	if goreleaserVersionOut == "" {
		return Errorf("expected to have goreleaser version output, got empty output")
	}

	// Run GoReleaser
	_, goreleaserErr := targetModDefault.
		Ctr().
		WithEnvVariable("GITHUB_TOKEN", "dummy-token").
		WithExec([]string{"git", "init"}).
		WithExec([]string{"git", "config", "--global", "user.name", "Test User"}).
		WithExec([]string{"git", "config", "--global", "user.email", "testuser@example.com"}).
		WithExec([]string{"rm", "-rf", "dist"}).
		WithExec([]string{"sh", "-c", "echo './dist' >>.gitignore"}).
		WithExec([]string{"git", "add", "."}).
		WithExec([]string{"git", "commit", "-m", "Initial commit"}).
		WithExec([]string{"goreleaser", "build", "--snapshot", "--skip=validate", "--clean"}).
		Stdout(ctx)

	if goreleaserErr != nil {
		return WrapError(goreleaserErr, "failed to run goreleaser")
	}

	// Run golangci-lint
	_, golangciErr := targetModDefault.
		Ctr().
		WithExec([]string{"golangci-lint", "run", "--config=.golangci.yml", "--verbose"}).
		Stdout(ctx)

	if golangciErr != nil {
		return WrapError(golangciErr, "failed to run golangci-lint")
	}

	return nil
}

// TestgotoolboxRunGo tests the functionality of the RunGo method in the Gotoolbox module.
// It performs two main checks:
// 1. Verifies that the Go version can be retrieved correctly.
// 2. Runs a specific Go test (TestFibonacci) and checks its output.
//
// This function uses the Dagger SDK to create and manipulate containers for testing.
//
// Parameters:
//   - ctx: The context for the test execution.
//
// Returns:
//   - error: An error if any part of the test fails, or nil if all checks pass.
func (m *Tests) TestgotoolboxRunGo(ctx context.Context) error {
	// Initialize the Go toolbox with the specified version.
	targetModDefault := dag.
		Gotoolbox(dagger.GotoolboxOpts{
			Version: "1.23.0-alpine3.19",
		}).WithSource(m.TestDir, dagger.GotoolboxWithSourceOpts{
		Workdir: "gotoolbox",
	})

	cmdVersion := []string{"version"}
	outVersion, versionErr := targetModDefault.RunGo(ctx,
		cmdVersion,
		dagger.GotoolboxRunGoOpts{})

	if versionErr != nil {
		return WrapError(versionErr, "failed to run Go version")
	}

	if !strings.Contains(outVersion, "go version go") {
		return WrapError(versionErr, "failed to get Go version")
	}

	cmdTest := []string{"test", "-v", "-run", "TestFibonacci"}

	targetTestMod := dag.Gotoolbox()
	outTest, testErr := targetTestMod.RunGo(ctx,
		cmdTest,
		dagger.GotoolboxRunGoOpts{
			Src:     m.TestDir,
			TestDir: "gotoolbox",
		})

	if testErr != nil {
		return WrapError(testErr, "failed to run Go test")
	}

	if outTest == "" {
		return Errorf("expected to have Go test output, got empty output")
	}

	if !strings.Contains(outTest, "PASS") {
		return Errorf("expected to have Go test output PASS, got %s", outTest)
	}

	// Test RunGo with environment variables and platform
	envVars := []string{"FOO=bar", "BAZ=qux"}
	outEnvVars, envVarsErr := targetModDefault.RunGo(ctx,
		cmdTest,
		dagger.GotoolboxRunGoOpts{
			Src:          m.TestDir,
			TestDir:      "gotoolbox",
			EnvVariables: envVars,
		},
	)

	if envVarsErr != nil {
		return WrapError(envVarsErr, "failed to run Go with environment variables")
	}

	if outEnvVars == "" {
		return Errorf("expected to have Go environment variables output, got empty output")
	}

	return nil
}
