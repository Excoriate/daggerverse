// A generated module for Tests functions
package main

import (
"context"
"errors"
"fmt"

"github.com/sourcegraph/conc/pool"
)

const emptyErrMsg = "the test output expected is empty"
const expectedContentNotMatchMsg = "an expected value does not match the actual value"
const underlyingDaggerErrMsg = "the dagger command failed or dagger returned an error"

var errEmptyOutput = errors.New(emptyErrMsg)
var errExpectedContentNotMatch = errors.New(expectedContentNotMatchMsg)
var errUnderlyingDagger = errors.New(underlyingDaggerErrMsg)

// Tests is a collection of tests.
//
// It's a struct that contains a single field, TestDir, which is a pointer to a Directory.
type Tests struct {
TestDir *Directory
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
func (m *Tests) getTestDir() *Directory {
return dag.CurrentModule().Source().Directory("./testdata")
}

// TestAll executes all tests.
func (m *Tests) TestAll(ctx context.Context) error {
polTests := pool.New().WithErrors().WithContext(ctx)

// TODO: Add your tests here
//polTests.Go(m.TestVersionOverride)

if err := polTests.Wait(); err != nil {
return fmt.Errorf("there are some failed tests: %w", err)
}

return nil
}

// TestTerminal returns a terminal for testing.
//func (m *Tests) TestTerminal() *Terminal {
//	targetModule := dag.Gotest().WithCgoEnabled().
//		WithBuildCache().WithGcccompiler().
//		WithSource(m.TestDir)
//
//	_, _ = targetModule.Ctr().Stdout(context.Background())
//
//	return targetModule.Ctr().Terminal()
//}

//// TestGoPrivate tests if the GOPRIVATE environment variable is set correctly.
//func (m *Tests) TestGoPrivate(ctx context.Context) error {
//	targetModule := dag.Gotest().
//		WithPrivateGoPkg("github.com/privatehost/private-repo")
//
//	out, err := targetModule.Ctr().WithExec([]string{"printenv"}).Stdout(ctx)
//	if err != nil {
//		return fmt.Errorf("failed to get env vars: %w", err)
//	}
//
//	if !strings.Contains(out, "GOPRIVATE=github.com/privatehost/private-repo") {
//		return fmt.Errorf("%w, expected GOPRIVATE to be set, got %s", errExpectedContentNotMatch, out)
//	}
//
//	return nil
//}
