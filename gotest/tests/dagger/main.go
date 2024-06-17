package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/sourcegraph/conc/pool"
)

type Tests struct {
	TestDir *Directory
}

func New() *Tests {
	t := &Tests{}

	t.TestDir = t.getTestDir()

	return t
}

// TestAll executes all tests.
func (m *Tests) TestAll(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.TestVersionOverride)
	p.Go(m.TestPassedEnvVars)
	p.Go(m.TestWithEnvVarAPI)
	p.Go(m.TestGoPrivate)
	p.Go(m.TestWithPlatformAPI)
	p.Go(m.TestCommandRunGoTestSimple)
	p.Go(m.TestCommandRunGoTestWithAdvancedOptions)
	p.Go(m.TestCommandRunGoTestSum)
	p.Go(m.TestCommandRunGoTestSumWithAdvancedOptions)

	return p.Wait() //nolint:wrapcheck // no need to wrap the error
}

// getTestDir returns the test directory.
//
// This is a helper method for tests, in order to get the test directory which
// is located in the same directory as the test file, and normally named as "testdata".
func (m *Tests) getTestDir() *Directory {
	return dag.CurrentModule().Source().Directory("./testdata")
}

// TestVersionOverride tests if the version is overridden correctly.
func (m *Tests) TestVersionOverride(_ context.Context) error {
	versions := []string{"1.21.0", "1.22.0", "1.22.1", "1.22.2"}
	for _, version := range versions {
		gt := dag.Gotest(GotestOpts{
			Version: version,
		})

		out, err := gt.Ctr().
			WithExec([]string{"go", "version"}).
			Stdout(context.Background())

		if err != nil {
			return err
		}

		if out == "" {
			return fmt.Errorf("go version is empty")
		}

		if !strings.Contains(out, version) {
			return fmt.Errorf("expected go version %s, got %s", version, out)
		}
	}

	fmt.Println("All versions are correct")
	return nil
}

// TestPassedEnvVars tests if the environment variables are passed correctly.
func (m *Tests) TestPassedEnvVars(_ context.Context) error {
	gt := dag.Gotest(GotestOpts{
		EnvVarsFromHost: "SOMETHING=SOMETHING,SOMETHING=SOMETHING",
	})

	out, err := gt.Ctr().
		WithExec([]string{"printenv"}).
		Stdout(context.Background())

	if err != nil {
		return err
	}

	if out == "" {
		return fmt.Errorf("env vars are empty")
	}

	if !strings.Contains(out, "SOMETHING") {
		return fmt.Errorf("expected env vars to be passed, got %s", out)
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
			return err
		}

		if out == "" {
			return fmt.Errorf("env vars are empty")
		}

		if !strings.Contains(out, name) {
			return fmt.Errorf("expected env vars to be passed, got %s", out)
		}
	}

	return nil
}

// TestGoPrivate tests if the GOPRIVATE environment variable is set correctly.
func (m *Tests) TestGoPrivate(ctx context.Context) error {
	mt := dag.Gotest().
		WithPrivateGoPkg("github.com/privatehost/private-repo")

	out, err := mt.Ctr().WithExec([]string{"printenv"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to get env vars: %w", err)
	}

	if !strings.Contains(out, "GOPRIVATE=github.com/privatehost/private-repo") {
		return fmt.Errorf("expected GOPRIVATE to be set, got %s", out)
	}

	return nil
}

// TestWithPlatformAPI tests if the platform is set correctly.
func (m *Tests) TestWithPlatformAPI(ctx context.Context) error {
	mt := dag.Gotest().WithPlatform("darwin/arm64/v7")

	out, err := mt.Ctr().WithExec([]string{"printenv"}).Stdout(ctx)

	if err != nil {
		return fmt.Errorf("failed to get env vars: %w", err)
	}

	if !strings.Contains(out, "GOOS=darwin") {
		return fmt.Errorf("expected GOOS to be set, got %s", out)
	}

	return nil
}

// TestTerminal returns a terminal for testing.
func (m *Tests) TestTerminal() *Terminal {
	gt := dag.Gotest().WithCgoEnabled().
		WithBuildCache().WithGcccompiler().
		WithSource(m.TestDir)

	_, _ = gt.Ctr().Stdout(context.Background())

	return gt.Ctr().Terminal()
}

// TestCommandRunGoTestSimple tests running go test.
func (m *Tests) TestCommandRunGoTestSimple(ctx context.Context) error {
	mt := dag.Gotest().WithSource(m.TestDir)
	out, err := mt.RunGoTest(ctx, m.TestDir, GotestRunGoTestOpts{})

	if err != nil {
		return fmt.Errorf("failed to run go test: %w", err)
	}

	if out == "" {
		return fmt.Errorf("go test output is empty")
	}

	return nil
}

// TestCommandRunGoTestSimple tests running go test.
func (m *Tests) TestCommandRunGoTestBug(ctx context.Context) (string, error) {
	mt := dag.Gotest().WithSource(m.TestDir)
	return mt.RunGoTest(ctx, m.TestDir, GotestRunGoTestOpts{})
}

// TestCommandRunGoTestWithAdvancedOptions tests running go test with advanced options.
func (m *Tests) TestCommandRunGoTestWithAdvancedOptions(ctx context.Context) error {
	mt := dag.Gotest().WithSource(m.TestDir).WithCgoEnabled().WithGcccompiler()
	out, err := mt.RunGoTest(ctx, m.TestDir, GotestRunGoTestOpts{
		EnableVerbose: true,
		EnableCache:   true,
		Race:          true,
		EnvVars:       []string{"SOMETHING=SOMETHING,SOMETHING=SOMETHING"},
	})

	if err != nil {
		return fmt.Errorf("failed to run go test: %w", err)
	}

	if out == "" {
		return fmt.Errorf("go test output is empty")
	}

	return nil
}

// TestCommandRunGoTestSum tests running go test with gotestsum.
func (m *Tests) TestCommandRunGoTestSum(ctx context.Context) error {
	mt := dag.Gotest().WithSource(m.TestDir)
	out, err := mt.RunGoTestSum(ctx, m.TestDir, GotestRunGoTestSumOpts{})

	if err != nil {
		return fmt.Errorf("failed to run go test: %w", err)
	}

	if out == "" {
		return fmt.Errorf("go test output is empty")
	}

	return nil
}

// TestCommandRunGoTestSumWithAdvancedOptions tests running go test with gotestsum with advanced options.
func (m *Tests) TestCommandRunGoTestSumWithAdvancedOptions(ctx context.Context) error {
	mt := dag.Gotest().
		WithSource(m.TestDir).
		WithGcccompiler().
		WithGoCache().
		WithCgoEnabled()

	out, err := mt.RunGoTestSum(ctx, m.TestDir, GotestRunGoTestSumOpts{
		InsecureRootCapabilities: true,
		EnablePretty:             true,
		PrintEnvVars:             true,
	})

	if err != nil {
		return fmt.Errorf("failed to run go test: %w", err)
	}

	if out == "" {
		return fmt.Errorf("go test output is empty")
	}

	return nil
}
