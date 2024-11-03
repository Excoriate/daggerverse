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

// TestGoTestRunTestWithCustomOptions tests various configurations of the Go test command
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
//nolint:cyclop,funlen // It's okay to have this size, it's by design.
func (m *Tests) TestGoTestRunTestWithCustomOptions(ctx context.Context) error {
	testDir := m.getTestDir("testdata/golang")

	// Test with custom version
	dagModuleWithVersion := dag.Gotest(
		dagger.GotestOpts{
			Version: "1.23.0-alpine3.20",
			EnvVarsFromHost: []string{
				"GOLANG_VERSION=1.22.5",
			},
		},
	)

	versionCtr := dagModuleWithVersion.
		WithGoCgoEnabled().
		WithGcccompilerInstalled().
		RunTest(testDir, dagger.GotestRunTestOpts{
			EnableDefaultOptions: true,
			EnvVars:              []string{"MY_ENV_VAR=my_value"},
		})

	if exitCode, err := versionCtr.
		ExitCode(ctx); err != nil {
		return WrapError(err, "failed to get exit code for version test")
	} else if exitCode != 0 {
		return Errorf("version test failed with exit code %d", exitCode)
	}

	// Basic test with default options and environment variables
	basicCtr := dagModuleWithVersion.
		WithGoCgoEnabled().
		WithGcccompilerInstalled().
		RunTest(testDir, dagger.GotestRunTestOpts{
			EnableDefaultOptions: true,
			EnvVars:              []string{"MY_ENV_VAR=my_value"},
		})

	if exitCode, err := basicCtr.ExitCode(ctx); err != nil {
		return WrapError(err, "failed to get exit code for basic test")
	} else if exitCode != 0 {
		return Errorf("basic test failed with exit code %d", exitCode)
	}

	// Coverage and profiling test case
	coverageCtr := dagModuleWithVersion.
		WithGoCgoEnabled().
		WithGcccompilerInstalled().
		RunTest(testDir, dagger.GotestRunTestOpts{
			EnableDefaultOptions: true,
			Cover:                true,
			Coverprofile:         "coverage.out",
			Cpuprofile:           "cpu.prof",
			Verbose:              true,
		})

	if exitCode, err := coverageCtr.ExitCode(ctx); err != nil {
		return WrapError(err, "failed to get exit code for coverage test")
	} else if exitCode != 0 {
		return Errorf("coverage test failed with exit code %d", exitCode)
	}

	// Race detection and build tags test case
	raceCtr := dagModuleWithVersion.
		WithGoCgoEnabled().
		WithGcccompilerInstalled().
		RunTest(testDir, dagger.GotestRunTestOpts{
			EnableDefaultOptions: true,
			Race:                 true,
			BuildTags:            "integration",
			Ldflags:              "-X main.version=test",
		})

	if exitCode, err := raceCtr.ExitCode(ctx); err != nil {
		return WrapError(err, "failed to get exit code for race test")
	} else if exitCode != 0 {
		return Errorf("race test failed with exit code %d", exitCode)
	}

	// Benchmark test case
	benchCtr := dagModuleWithVersion.
		WithGoCgoEnabled().
		WithGcccompilerInstalled().
		RunTest(testDir, dagger.GotestRunTestOpts{
			EnableDefaultOptions: true,
			Benchmark:            ".",
			Benchmem:             true,
			Benchtime:            "1s",
			TestCount:            3,
		})

	if exitCode, err := benchCtr.ExitCode(ctx); err != nil {
		return WrapError(err, "failed to get exit code for benchmark test")
	} else if exitCode != 0 {
		return Errorf("benchmark test failed with exit code %d", exitCode)
	}

	// Test filtering and parallel execution test case
	filteredCtr := dagModuleWithVersion.
		WithGoCgoEnabled().
		WithGcccompilerInstalled().
		RunTest(testDir, dagger.GotestRunTestOpts{
			EnableDefaultOptions: true,
			Run:                  "TestSpecific",
			Parallel:             4,
			Timeout:              "30s",
			Failfast:             true,
			EnableJsonoutput:     true,
		})

	if exitCode, err := filteredCtr.ExitCode(ctx); err != nil {
		return WrapError(err, "failed to get exit code for filtered test")
	} else if exitCode != 0 {
		return Errorf("filtered test failed with exit code %d", exitCode)
	}

	// Build mode and compiler flags test case
	buildCtr := dagModuleWithVersion.
		WithGoCgoEnabled().
		WithGcccompilerInstalled().
		RunTest(testDir, dagger.GotestRunTestOpts{
			EnableDefaultOptions: true,
			BuildMode:            "pie",
			Gcflags:              "-N -l",
			Mod:                  "readonly",
			Trimpath:             true,
			Work:                 true,
		})

	if exitCode, err := buildCtr.ExitCode(ctx); err != nil {
		return WrapError(err, "failed to get exit code for build mode test")
	} else if exitCode != 0 {
		return Errorf("build mode test failed with exit code %d", exitCode)
	}

	// Short mode with specific packages test case
	packagesCtr := dagModuleWithVersion.
		WithGoCgoEnabled().
		WithGcccompilerInstalled().
		RunTest(testDir, dagger.GotestRunTestOpts{
			EnableDefaultOptions: true,
			Packages:             []string{"./..."},
			Short:                true,
			Verbose:              true,
			EnvVars:              []string{"TEST_ENV=test"},
		})

	if exitCode, err := packagesCtr.ExitCode(ctx); err != nil {
		return WrapError(err, "failed to get exit code for packages test")
	} else if exitCode != 0 {
		return Errorf("packages test failed with exit code %d", exitCode)
	}

	return nil
}

// TestGoTestRunTestCMDWithCustomOptions tests various configurations of the Go test command
// using different RunTestCmd options.
//
// Parameters:
//
//	ctx - The context for managing cancellation and deadlines.
//
// Returns:
//
//	An error if any test configuration fails; otherwise, it returns nil.
//
//nolint:cyclop // Test function with multiple test cases by design
func (m *Tests) TestGoTestRunTestCMDWithCustomOptions(ctx context.Context) error {
	testDir := m.getTestDir("testdata/golang")
	dagMod := dag.Gotest().
		WithGoCgoEnabled().
		WithGcccompilerInstalled()

	// Test 1: Test with race detection and build tags
	raceOut, raceErr := dagMod.RunTestCmd(ctx, testDir, dagger.GotestRunTestCmdOpts{
		Race:      true,
		BuildTags: "integration",
		Packages:  []string{"./..."},
		EnvVars:   []string{"TEST_ENV=race_test"},
	})
	if raceErr != nil {
		return WrapError(raceErr, "race detection test failed")
	}
	if !strings.Contains(raceOut, "PASS") {
		return NewError("expected PASS in race test output")
	}

	// Test 2: Test coverage and verbose output
	coverOut, coverErr := dagMod.RunTestCmd(ctx, testDir, dagger.GotestRunTestCmdOpts{
		Cover:        true,
		Coverprofile: "coverage.out",
		Verbose:      true,
	})

	if coverErr != nil {
		return WrapError(coverErr, "coverage test failed")
	}

	if !strings.Contains(coverOut, "coverage") {
		return NewError("expected coverage data in output")
	}

	// Test 3: Test with specific compiler flags and build options
	buildOut, buildErr := dagMod.RunTestCmd(ctx, testDir, dagger.GotestRunTestCmdOpts{
		Gcflags:  "-N -l",
		Ldflags:  "-w -s",
		Mod:      "readonly",
		Trimpath: true,
		Packages: []string{"./..."},
		Verbose:  true,
	})

	if buildErr != nil {
		return WrapError(buildErr, "build options test failed")
	}

	if !strings.Contains(buildOut, "PASS") {
		return NewError("expected PASS in build output")
	}

	// Test 4: Test with short mode and JSON output
	shortOut, shortErr := dagMod.RunTestCmd(ctx, testDir, dagger.GotestRunTestCmdOpts{
		Short:            true,
		EnableJsonoutput: true,
		Verbose:          true,
		Packages:         []string{"./..."},
	})

	if shortErr != nil {
		return WrapError(shortErr, "short mode test failed")
	}

	// JSON output should contain either "Action" or "Test" fields
	if !strings.Contains(shortOut, "Action") && !strings.Contains(shortOut, "Test") {
		return NewError("expected JSON test output")
	}

	// Test 5: Test with timeout and fail-fast options
	timeoutOut, timeoutErr := dagMod.RunTestCmd(ctx, testDir, dagger.GotestRunTestCmdOpts{
		Timeout:  "30s",
		Failfast: true,
		Verbose:  true,
		Packages: []string{"./..."},
	})

	if timeoutErr != nil {
		return WrapError(timeoutErr, "timeout test failed")
	}

	if !strings.Contains(timeoutOut, "PASS") {
		return NewError("expected PASS in timeout test output")
	}

	return nil
}
