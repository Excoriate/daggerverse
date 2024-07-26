package main

import (
	"fmt"

	"github.com/excoriate/daggerverse/gotest/internal/dagger"

	"github.com/Excoriate/daggerx/pkg/envvars"
)

// SetupGoTest sets up the go test options, to either evaluate the container and run the test,
// or return the container to be evaluated later.
func (m *Gotest) SetupGoTest(
	// The directory containing code to test.
	src *dagger.Directory,
	// Packages to test.
	// +optional
	packages []string,
	// enableVerbose is a flag to run tests with -v.
	// +optional
	enableVerbose bool,
	// race is a flag to run tests with
	// +optional
	race bool,
	// Arbitrary flags to pass along to go test.
	// +optional
	testFlags []string,
	// Whether to run tests insecurely, i.e. with special privileges.
	// +optional
	insecureRootCapabilities bool,
	// Enable experimental Dagger nesting. It sets the ExperimentalPrivilegedNesting option in Dagger.
	// +optional
	enableNest bool,
	// enableCache is a flag to enable cache volumes. If it's set, it'll
	// enable the cache volumes for the Go module and build cache.
	// +optional
	enableCache bool,
	// envVars is a list of environment variables with the format "SOMETHING=SOMETHING,SOMETHING=SOMETHING".
	// +optional
	envVars []string,
	// printEnvVars is a flag to print the environment variables
	// +optional
	printEnvVars bool,
) (*dagger.Container, error) {
	goTest := []string{"go", "test"}
	ctr := m.WithSource(src, "").Ctr

	if enableCache {
		ctr = m.WithGoCache().Ctr
	}

	if len(envVars) > 0 {
		envVarsDagger, err := envvars.ToDaggerEnvVarsFromSlice(envVars)
		if err != nil {
			return nil, fmt.Errorf("failed to parse environment variables: %w", err)
		}

		for _, envVar := range envVarsDagger {
			ctr = m.WithEnvironmentVariable(envVar.Name, envVar.Value, false).Ctr
		}
	}

	if printEnvVars {
		//nolint:exhaustruct // It's 'okaysh' for now, I'll decide later what's going to be the pattern here.
		ctr = ctr.WithExec([]string{"printenv"}, dagger.ContainerWithExecOpts{
			InsecureRootCapabilities:      insecureRootCapabilities,
			ExperimentalPrivilegedNesting: enableNest,
		})
	}

	pkgs := packages
	if len(pkgs) == 0 {
		pkgs = []string{"./..."}
	}

	if race {
		goTest = append(goTest, "-race")
	}

	if enableVerbose {
		goTest = append(goTest, "-v")
	}

	goTest = append(goTest, testFlags...)

	goTest = append(goTest, pkgs...)

	//nolint:exhaustruct // It's 'okaysh' for now, I'll decide later what's going to be the pattern here.
	ctr = ctr.WithExec(goTest, dagger.ContainerWithExecOpts{
		InsecureRootCapabilities:      insecureRootCapabilities,
		ExperimentalPrivilegedNesting: enableNest,
	})

	m.Ctr = ctr

	return m.Ctr, nil
}

// SetupGoTestSum sets up the go test options, to either evaluate the container and run the test,
// or return the container to be evaluated later.
//
//nolint:cyclop // It's 'okaysh' for now, I'll decide later what's going to be the pattern here.
func (m *Gotest) SetupGoTestSum(
	// The directory containing code to test.
	src *dagger.Directory,
	// Packages to test.
	// +optional
	packages []string,
	// Arbitrary flags to pass along to go test.
	// +optional
	testFlags []string,
	// goTestSumFlags is a flag to pass along to go test -json.
	// +optional
	goTestSumFlags []string,
	// format defines the option for the GoTestsum format to display
	// +optional
	format string,
	// Whether to run tests insecurely, i.e. with special privileges.
	// +optional
	insecureRootCapabilities bool,
	// Enable experimental Dagger nesting. It sets the ExperimentalPrivilegedNesting option in Dagger.
	// +optional
	enableNest bool,
	// enableCache is a flag to enable cache volumes. If it's set, it'll
	// enable the cache volumes for the Go module and build cache.
	// +optional
	enableCache bool,
	// enablePretty is a flag to enable pretty output.
	// +optional
	enablePretty bool,
	// envVars is a list of environment variables to set with the format "SOMETHING=SOMETHING,SOMETHING=SOMETHING".
	// +optional
	envVars []string,
	// printEnvVars is a flag to print the environment variables
	// +optional
	printEnvVars bool,
) (*dagger.Container, error) {
	goTestSumInstallCMD := []string{"go", "install", "gotest.tools/gotestsum@latest"}
	goTestInstallTparseCMD := []string{"go", "install", "github.com/mfridman/tparse@latest"}
	goTestCMD := []string{"gotestsum", "--no-color=false"}
	ctr := m.WithSource(src, "").Ctr

	if enableCache {
		ctr = m.WithGoCache().Ctr
	}

	if len(envVars) > 0 {
		envVarsDagger, err := envvars.ToDaggerEnvVarsFromSlice(envVars)
		if err != nil {
			return nil, fmt.Errorf("failed to parse environment variables: %w", err)
		}

		for _, envVar := range envVarsDagger {
			ctr = m.WithEnvironmentVariable(envVar.Name, envVar.Value, false).Ctr
		}
	}

	if printEnvVars {
		ctr = ctr.WithFocus().
			//nolint:exhaustruct // It's 'okaysh' for now, I'll decide later what's going to be the pattern here.
			WithExec([]string{"printenv"}, dagger.ContainerWithExecOpts{
				InsecureRootCapabilities:      insecureRootCapabilities,
				ExperimentalPrivilegedNesting: enableNest,
			})
	}

	ctr = ctr.WithExec(goTestSumInstallCMD)

	if format == "" {
		format = "testname" // opinionated default
	}

	goTestCMD = append(goTestCMD, "--format="+format)
	goTestCMD = append(goTestCMD, goTestSumFlags...)

	if len(packages) > 0 {
		goTestCMD = append(goTestCMD, packages...)
	}

	if len(testFlags) > 0 {
		goTestCMD = append(goTestCMD, "--")
		goTestCMD = append(goTestCMD, testFlags...)
	}

	if enablePretty {
		ctr = ctr.WithExec(goTestInstallTparseCMD)

		goTestCMD = append(goTestCMD, "--jsonfile", "test-output.json")
	}

	//nolint:exhaustruct // It's 'okaysh' for now, I'll decide later what's going to be the pattern here.
	ctr = ctr.WithExec(goTestCMD, dagger.ContainerWithExecOpts{
		InsecureRootCapabilities:      insecureRootCapabilities,
		ExperimentalPrivilegedNesting: enableNest,
	})

	if enablePretty {
		tParseCMD := []string{"tparse", "-all", "-smallscreen", "-file=test-output.json"}
		//nolint:exhaustruct // It's 'okaysh' for now, I'll decide later what's going to be the pattern here.
		ctr = ctr.WithExec(tParseCMD, dagger.ContainerWithExecOpts{
			InsecureRootCapabilities:      insecureRootCapabilities,
			ExperimentalPrivilegedNesting: enableNest,
		})
	}

	m.Ctr = ctr

	return m.Ctr, nil
}
