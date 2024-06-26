// A generated module for Tests functions
package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

// TestAll executes all tests.
func (m *Tests) TestAll(ctx context.Context) error {
	polTests := pool.New().WithErrors().WithContext(ctx)

	polTests.Go(m.TestVersionOverride)
	polTests.Go(m.TestPassedEnvVars)
	polTests.Go(m.TestWithEnvVarAPI)
	polTests.Go(m.TestGoPrivate)
	polTests.Go(m.TestWithPlatformAPI)
	polTests.Go(m.TestCommandRunGoTestSimple)
	polTests.Go(m.TestCommandRunGoTestWithAdvancedOptions)
	polTests.Go(m.TestCommandRunGoTestSum)
	polTests.Go(m.TestCommandRunGoTestSumWithAdvancedOptions)

	if err := polTests.Wait(); err != nil {
		return fmt.Errorf("there are some failed tests: %w", err)
	}

	return nil
}

// getTestDir returns the test directory.
//
// This is a helper method for tests, in order to get the test directory which
// is located in the same directory as the test file, and normally named as "testdata".
func (m *Tests) getTestDir() *Directory {
	return dag.CurrentModule().Source().Directory("./testdata")
}

// TestVersionOverride tests if the version is overridden correctly.
func (m *Tests) TestVersionOverride(ctx context.Context) error {
	versions := []string{"1.21.0", "1.22.0", "1.22.1", "1.22.2"}
	for _, version := range versions {
		gt := dag.Gotest(GotestOpts{
			Version: version,
		})

		out, err := gt.Ctr().
			WithExec([]string{"go", "version"}).
			Stdout(ctx)

		if err != nil {
			return fmt.Errorf("failed to get go version: %w", err)
		}

		if out == "" {
			return fmt.Errorf("%w", errEmptyOutput)
		}

		if !strings.Contains(out, version) {
			return fmt.Errorf("mismatch of Go version, %w", errExpectedContentNotMatch)
		}
	}

	fmt.Println("All versions are correct")

	return nil
}

// TestPassedEnvVars tests if the environment variables are passed correctly.
func (m *Tests) TestPassedEnvVars(ctx context.Context) error {
	targetModule := dag.Gotest(GotestOpts{
		EnvVarsFromHost: "SOMETHING=SOMETHING,SOMETHING=SOMETHING",
	})

	out, err := targetModule.Ctr().
		WithExec([]string{"printenv"}).
		Stdout(ctx)

	if err != nil {
		return fmt.Errorf("failed when executing printenv: %w", err)
	}

	if out == "" {
		return fmt.Errorf("%w, env vars are empty", errEmptyOutput)
	}

	if !strings.Contains(out, "SOMETHING") {
		return fmt.Errorf("%w, expected env vars to be passed, got %s", errExpectedContentNotMatch, out)
	}

	return nil
}

// TestWithEnvVarAPI tests if the environment variables are passed correctly using the API.
func (m *Tests) TestWithEnvVarAPI(ctx context.Context) error {
	setOfEnvVars := []string{"RANDOM_VAR=RANDOM_VALUE", "ANOTHER_VAR=ANOTHER_VALUE"}

	for _, envVar := range setOfEnvVars {
		envVarSplit := strings.Split(envVar, "=")
		name := envVarSplit[0]
		value := envVarSplit[1]

		gt := dag.Gotest().
			WithEnvironmentVariable(name, value, GotestWithEnvironmentVariableOpts{
				Expand: true,
			})

		out, err := gt.Ctr().
			WithExec([]string{"printenv"}).
			Stdout(ctx)

		if err != nil {
			return fmt.Errorf("%w, failed with an error: %w", errUnderlyingDagger, err)
		}

		if out == "" {
			return fmt.Errorf("%w, env vars are empty", errEmptyOutput)
		}

		if !strings.Contains(out, name) {
			return fmt.Errorf("%w, expected env vars to be passed, got %s", errExpectedContentNotMatch, out)
		}
	}

	return nil
}

// TestGoPrivate tests if the GOPRIVATE environment variable is set correctly.
func (m *Tests) TestGoPrivate(ctx context.Context) error {
	targetModule := dag.Gotest().
		WithPrivateGoPkg("github.com/privatehost/private-repo")

	out, err := targetModule.Ctr().WithExec([]string{"printenv"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to get env vars: %w", err)
	}

	if !strings.Contains(out, "GOPRIVATE=github.com/privatehost/private-repo") {
		return fmt.Errorf("%w, expected GOPRIVATE to be set, got %s", errExpectedContentNotMatch, out)
	}

	return nil
}

// TestWithPlatformAPI tests if the platform is set correctly.
func (m *Tests) TestWithPlatformAPI(ctx context.Context) error {
	targetModule := dag.Gotest().
		WithPlatform("darwin/arm64/v7")

	out, err := targetModule.Ctr().
		WithExec([]string{"printenv"}).Stdout(ctx)

	if err != nil {
		return fmt.Errorf("%w, failed to get env vars: %w", errUnderlyingDagger, err)
	}

	if !strings.Contains(out, "GOOS=darwin") {
		return fmt.Errorf("%w, expected GOOS to be set, got %s", errExpectedContentNotMatch, out)
	}

	return nil
}

// TestTerminal returns a terminal for testing.
func (m *Tests) TestTerminal() *Terminal {
	targetModule := dag.Gotest().WithCgoEnabled().
		WithBuildCache().WithGcccompiler().
		WithSource(m.TestDir)

	_, _ = targetModule.Ctr().Stdout(context.Background())

	return targetModule.Ctr().Terminal()
}

// TestCommandRunGoTestSimple tests running go test.
func (m *Tests) TestCommandRunGoTestSimple(ctx context.Context) error {
	targetModule := dag.Gotest().WithSource(m.TestDir)
	out, err := targetModule.RunGoTest(ctx, m.TestDir, GotestRunGoTestOpts{})

	if err != nil {
		return fmt.Errorf("%w, failed to run go test: %w", errUnderlyingDagger, err)
	}

	if out == "" {
		return fmt.Errorf("%w", errEmptyOutput)
	}

	return nil
}

// TestCommandRunGoTestWithAdvancedOptions tests running go test with advanced options.
func (m *Tests) TestCommandRunGoTestWithAdvancedOptions(ctx context.Context) error {
	targetModule := dag.Gotest().WithSource(m.TestDir).WithCgoEnabled().WithGcccompiler()
	out, err := targetModule.RunGoTest(ctx, m.TestDir, GotestRunGoTestOpts{
		EnableVerbose: true,
		EnableCache:   true,
		Race:          true,
		EnvVars:       []string{"SOMETHING=SOMETHING,SOMETHING=SOMETHING"},
	})

	if err != nil {
		return fmt.Errorf("%w, failed to run go test: %w", errUnderlyingDagger, err)
	}

	if out == "" {
		return fmt.Errorf("%w", errEmptyOutput)
	}

	return nil
}

// TestCommandRunGoTestSum tests running go test with gotestsum.
func (m *Tests) TestCommandRunGoTestSum(ctx context.Context) error {
	targetModule := dag.Gotest().WithSource(m.TestDir)
	out, err := targetModule.RunGoTestSum(ctx, m.TestDir, GotestRunGoTestSumOpts{})

	if err != nil {
		return fmt.Errorf("%w, failed to run go test TestCommandRunGoTestSum: %w", errUnderlyingDagger, err)
	}

	if out == "" {
		return fmt.Errorf("%w", errEmptyOutput)
	}

	return nil
}

// TestCommandRunGoTestSumWithAdvancedOptions tests running go test with gotestsum with advanced options.
func (m *Tests) TestCommandRunGoTestSumWithAdvancedOptions(ctx context.Context) error {
	targetModule := dag.Gotest().
		WithSource(m.TestDir).
		WithGcccompiler().
		WithGoCache().
		WithCgoEnabled()

	out, err := targetModule.RunGoTestSum(ctx, m.TestDir, GotestRunGoTestSumOpts{
		InsecureRootCapabilities: true,
		EnablePretty:             true,
		PrintEnvVars:             true,
	})

	if err != nil {
		return fmt.Errorf("%w, failed to run  TestCommandRunGoTestSumWithAdvancedOptions: %w", errUnderlyingDagger, err)
	}

	if out == "" {
		return fmt.Errorf("%w", errEmptyOutput)
	}

	return nil
}
