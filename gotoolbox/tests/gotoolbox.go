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
