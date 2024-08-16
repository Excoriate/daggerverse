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
