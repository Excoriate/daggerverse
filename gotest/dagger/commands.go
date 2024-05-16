package main

import "context"

// RunGoTest runs tests using the go test CLI.
// The default packages to test are "./...".
func (m *Gotest) RunGoTest(
	// The directory containing code to test.
	src *Directory,
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
	// Enable experimental Dagger nesting.
	// +optional
	nest bool,
	// enableCache is a flag to enable cache volumes.
	// +optional
	enableCache bool,
	// envVars is a list of environment variables to set in the container with the format "SOMETHING=SOMETHING,SOMETHING=SOMETHING".
	// +optional
	envVars []string,
	// printEnvVars is a flag to print the environment variables
	// +optional
	printEnvVars bool,
) (string, error) {
	evaluatedCtr, err := m.SetupGoTest(src, packages, enableVerbose, race, testFlags, insecureRootCapabilities, nest, enableCache, envVars, printEnvVars)
	if err != nil {
		return "", err
	}

	out, ctrErr := evaluatedCtr.Stdout(context.Background())
	return out, ctrErr
}

// RunGoTestSum runs tests using the gotestsum CLI.
func (m *Gotest) RunGoTestSum(
	// The directory containing code to test.
	src *Directory,
	// Packages to test.
	// +optional
	packages []string,
	// race is a flag to run tests with
	// +optional
	race bool,
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
	// envVars is a list of environment variables to set in the container with the format "SOMETHING=SOMETHING,SOMETHING=SOMETHING".
	// +optional
	envVars []string,
	// printEnvVars is a flag to print the environment variables
	// +optional
	printEnvVars bool,
) (string, error) {
	evaluatedCtr, err := m.SetupGoTestSum(src, packages, race, testFlags, goTestSumFlags, format, insecureRootCapabilities, enableNest, enableCache, enablePretty, envVars, printEnvVars)
	if err != nil {
		return "", err
	}

	out, ctrErr := evaluatedCtr.Stdout(context.Background())
	return out, ctrErr
}
