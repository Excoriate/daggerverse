package main

import (
	"context"

	"github.com/Excoriate/daggerverse/gotest/internal/dagger"
)

const (
	cmdEntrypoint = "go"
	cmdTest       = "test"
)

// getBaseCmd returns the base command for running Go tests.
//
// This function constructs and returns a slice of strings that represents
// the base command to execute Go tests, which includes the command
// entry point and the test subcommand.
func getBaseCmd() []string {
	return []string{cmdEntrypoint, cmdTest}
}

// GoTestCmd represents a command to run Go tests in a Dagger container.
//
// This struct holds the container in which the command will be executed
// and the command arguments to be passed to the Go test command.
type GoTestCmd struct {
	// BaseCmd is the base command to run.
	// +private
	BaseCmd []string
	// Packages are the packages to test.
	// +private
	Packages []string
	// EnvironmentVariables is the environment variables to set.
	// +private
	EnvironmentVariables []string
	// Secrets are the secrets to set.
	// +private
	Secrets []*dagger.Secret
}

// newGoTestCmd creates a new GoTestCmd instance for running Go tests.
//
// This method initializes a GoTestCmd with the specified packages, environment
// variables, and secrets. It constructs the command to be executed in a Dagger
// container, setting up the necessary parameters for the Go test command.
//
// Parameters:
//   - packages: A slice of strings representing the Go packages to be tested.
//   - environmentVariables: A slice of strings representing the environment
//     variables to set for the test execution.
//   - secrets: A slice of pointers to dagger.Secret representing the secrets
//     to be used during the test execution.
//
// Returns:
// - *GoTestCmd: A pointer to the newly created GoTestCmd instance.
func (m *Gotest) newGoTestCmd(
	packages []string,
	environmentVariables []string,
	secrets []*dagger.Secret,
) *GoTestCmd {
	return &GoTestCmd{
		BaseCmd:              getBaseCmd(),
		Packages:             packages,
		EnvironmentVariables: environmentVariables,
		Secrets:              secrets,
	}
}

// WithDefaultOptions sets the default options for running Go tests.
//
// This method configures the GoTestCmd instance with default package
// settings, enabling verbose output for test results and activating
// the race detector for the Go build process. It modifies the BaseCmd
// field to include the necessary flags for the Go test command.
//
// Returns:
//   - *GoTestCmd: A pointer to the updated GoTestCmd instance with
//     default options applied.
func (c *GoTestCmd) WithDefaultOptions() *GoTestCmd {
	c.Packages = []string{"./..."}

	testOpts := NewGoTestOptions()
	testOpts.WithVerboseOutput()

	buildOpts := NewGoBuildOptions()
	buildOpts.WithRace()

	c.BaseCmd = append(c.BaseCmd, buildOpts.Flags...)
	c.BaseCmd = append(c.BaseCmd, testOpts.Flags...)
	c.BaseCmd = append(c.BaseCmd, c.Packages...)

	return c
}

// RunTest executes the Go tests for the specified source code and packages.
// It allows for customization of the test execution through various options,
// including environment variables, secrets, and build/test flags.
//
// Parameters:
//   - source: The source code to test.
//   - packages: A slice of strings representing the packages to test.
//   - envVars: A slice of strings representing the environment variables to set.
//   - secrets: A slice of pointers to dagger.Secret representing the secrets to set.
//   - enableDefaultOptions: A boolean flag to enable default options for the test command.
//   - race: A boolean flag to enable the race detector in the Go command.
//   - msan: A boolean flag to enable memory sanitizer in the Go command.
//   - asan: A boolean flag to enable address sanitizer in the Go command.
//   - buildTags: A string specifying build constraints for the Go command.
//   - ldflags: A string setting flags for the linker in the Go command.
//   - gcflags: A string setting flags for the Go compiler.
//   - asmflags: A string setting flags for the assembler in the Go command.
//   - trimpath: A boolean flag to remove all file system paths from the compiled binary.
//   - work: A boolean flag to enable the creation of a temporary work directory.
//   - buildMode: A string specifying the build mode for the Go command.
//   - compiler: A string specifying the compiler to use for building.
//   - gccgoflags: A string setting flags for the gccgo compiler.
//   - mod: A string specifying the module mode for the Go command.
//   - benchmark: A string specifying the benchmark to run.
//   - benchmem: A boolean flag to enable memory allocation statistics.
//   - benchtime: A string specifying the duration for benchmarks.
//   - blockprofile: A string specifying the file for block profiling.
//   - cover: A boolean flag to enable coverage analysis.
//   - coverprofile: A string specifying the file for coverage profile output.
//   - cpuprofile: A string specifying the file for CPU profiling.
//   - testCount: An integer specifying the number of test iterations.
//   - failfast: A boolean flag to stop the test run on the first failure.
//   - enableJSONOutput: A boolean flag to enable JSON output for test results.
//   - list: A string specifying a regex to filter tests.
//   - memprofile: A string specifying the file for memory profiling.
//   - mutexprofile: A string specifying the file for mutex profiling.
//   - parallel: An integer specifying the maximum number of tests to run in parallel.
//   - run: A string specifying a regex to select tests to run.
//   - short: A boolean flag to enable short test mode.
//   - timeout: A string specifying the maximum time to run tests.
//   - verbose: A boolean flag that sets the verbosity level in the Go command.
//
// Returns:
//   - A pointer to a dagger.Container representing the test execution container,
//     and an error if the execution fails.
//
//nolint:funlen,gocognit,cyclop,gocyclo,maintidx // It's okay to have this size, it's by design.
func (m *Gotest) RunTest(
	// source is the source code to test.
	source *dagger.Directory,
	// packages are the packages to test.
	// +optional
	packages []string,
	// envVars are the environment variables to set.
	// +optional
	envVars []string,
	// secrets are the secrets to set.
	// +optional
	secrets []*dagger.Secret,
	// enableDefaultOptions enables the default options for the test command.
	// +optional
	enableDefaultOptions bool,
	// Build Options
	// race enables the race detector in the Go command.
	// It's equivalent to the -race flag.
	// +optional
	race bool,
	// msan enables memory sanitizer in the Go command.
	// It's equivalent to the -msan flag.
	// +optional
	msan bool,
	// asan enables address sanitizer in the Go command.
	// It's equivalent to the -asan flag.
	// +optional
	asan bool,
	// buildTags specifies build constraints for the Go command.
	// It's equivalent to the -tags flag.
	// +optional
	buildTags string,
	// ldflags sets flags for the linker in the Go command.
	// It's equivalent to the -ldflags flag.
	// +optional
	ldflags string,
	// gcflags sets flags for the Go compiler.
	// It's equivalent to the -gcflags flag.
	// +optional
	gcflags string,
	// asmflags sets flags for the assembler in the Go command.
	// It's equivalent to the -asmflags flag.
	// +optional
	asmflags string,
	// trimpath removes all file system paths from the compiled binary.
	// It's equivalent to the -trimpath flag.
	// +optional
	trimpath bool,
	// work enables the creation of a temporary work directory.
	// It's equivalent to the -work flag.
	// +optional
	work bool,
	// buildMode specifies the build mode for the Go command.
	// It's equivalent to the -buildmode flag.
	// +optional
	buildMode string,
	// compiler specifies the compiler to use for building.
	// It's equivalent to the -compiler flag.
	// +optional
	compiler string,
	// gccgoflags sets flags for the gccgo compiler.
	// It's equivalent to the -gccgoflags flag.
	// +optional
	gccgoflags string,
	// mod specifies the module mode for the Go command.
	// It's equivalent to the -mod flag.
	// +optional
	mod string,
	// Test Options
	// benchmark specifies the benchmark to run.
	// It's equivalent to the -bench flag.
	// +optional
	benchmark string,
	// benchmem enables memory allocation statistics.
	// It's equivalent to the -benchmem flag.
	// +optional
	benchmem bool,
	// benchtime specifies the duration for benchmarks.
	// It's equivalent to the -benchtime flag.
	// +optional
	benchtime string,
	// blockprofile specifies the file for block profiling.
	// It's equivalent to the -blockprofile flag.
	// +optional
	blockprofile string,
	// cover enables coverage analysis.
	// It's equivalent to the -cover flag.
	// +optional
	cover bool,
	// coverprofile specifies the file for coverage profile output.
	// It's equivalent to the -coverprofile flag.
	// +optional
	coverprofile string,
	// cpuprofile specifies the file for CPU profiling.
	// It's equivalent to the -cpuprofile flag.
	// +optional
	cpuprofile string,
	// testCount specifies the number of test iterations.
	// It's equivalent to the -count flag.
	// +optional
	testCount int,
	// failfast stops the test run on the first failure.
	// It's equivalent to the -failfast flag.
	// +optional
	failfast bool,
	// enableJSONOutput enables JSON output for test results.
	// It's equivalent to the -json flag.
	// +optional
	enableJSONOutput bool,
	// list specifies a regex to filter tests.
	// It's equivalent to the -list flag.
	// +optional
	list string,
	// memprofile specifies the file for memory profiling.
	// It's equivalent to the -memprofile flag.
	// +optional
	memprofile string,
	// mutexprofile specifies the file for mutex profiling.
	// It's equivalent to the -mutexprofile flag.
	// +optional
	mutexprofile string,
	// parallel specifies the maximum number of tests to run in parallel.
	// It's equivalent to the -parallel flag.
	// +optional
	parallel int,
	// run specifies a regex to select tests to run.
	// It's equivalent to the -run flag.
	// +optional
	run string,
	// short enables short test mode.
	// It's equivalent to the -short flag.
	// +optional
	short bool,
	// timeout specifies the maximum time to run tests.
	// It's equivalent to the -timeout flag.
	// +optional
	timeout string,
	// verbose is the flag that sets the verbosity level in the Go command.
	// It's equivalent to the -v flag.
	// +optional
	verbose bool,
) (*dagger.Container, error) {
	gtCmd := m.newGoTestCmd(packages, envVars, secrets)
	m = m.WithSource(source, "")

	if enableDefaultOptions {
		gtCmd.WithDefaultOptions()

		return m.
			Ctr.
			WithExec(gtCmd.BaseCmd), nil
	}

	cmdToAppend := gtCmd.BaseCmd

	buildOpts := NewGoBuildOptions()
	testOptions := NewGoTestOptions()

	// Apply Build Options
	if race {
		buildOpts = buildOpts.WithRace()
	}

	if msan {
		buildOpts = buildOpts.WithMSan()
	}

	if asan {
		buildOpts = buildOpts.WithASan()
	}

	if buildTags != "" {
		buildOpts = buildOpts.WithTags(buildTags)
	}

	if ldflags != "" {
		buildOpts = buildOpts.WithLDFlags(ldflags)
	}

	if gcflags != "" {
		buildOpts = buildOpts.WithGCFlags(gcflags)
	}

	if asmflags != "" {
		buildOpts = buildOpts.WithAsmFlags(asmflags)
	}

	if trimpath {
		buildOpts = buildOpts.WithTrimPath()
	}

	if work {
		buildOpts = buildOpts.WithWork()
	}

	if buildMode != "" {
		buildOpts = buildOpts.WithBuildMode(buildMode)
	}

	if compiler != "" {
		buildOpts = buildOpts.WithCompiler(compiler)
	}

	if gccgoflags != "" {
		buildOpts = buildOpts.WithGCCGOFlags(gccgoflags)
	}

	if mod != "" {
		buildOpts = buildOpts.WithMod(mod)
	}

	// Apply Test Options
	if benchmark != "" {
		testOptions = testOptions.WithBenchmark(benchmark)
	}

	if benchmem {
		testOptions = testOptions.WithBenchmarkMemory()
	}

	if benchtime != "" {
		testOptions = testOptions.WithBenchmarkTime(benchtime)
	}

	if blockprofile != "" {
		testOptions = testOptions.WithBlockProfile(blockprofile)
	}

	if cover {
		testOptions = testOptions.WithCoverage()
	}

	if coverprofile != "" {
		testOptions = testOptions.WithCoverageProfile(coverprofile)
	}

	if cpuprofile != "" {
		testOptions = testOptions.WithCPUProfile(cpuprofile)
	}

	if testCount > 0 {
		testOptions = testOptions.WithTestCount(testCount)
	}

	if failfast {
		testOptions = testOptions.WithFailFast()
	}

	if enableJSONOutput {
		testOptions = testOptions.WithJSONOutput()
	}

	if list != "" {
		testOptions = testOptions.WithListTests(list)
	}

	if memprofile != "" {
		testOptions = testOptions.WithMemoryProfile(memprofile)
	}

	if mutexprofile != "" {
		testOptions = testOptions.WithMutexProfile(mutexprofile)
	}

	if parallel > 0 {
		testOptions = testOptions.WithParallelTests(parallel)
	}

	if run != "" {
		testOptions = testOptions.WithTestFilter(run)
	}

	if short {
		testOptions = testOptions.WithShortTest()
	}

	if timeout != "" {
		testOptions = testOptions.WithTimeout(timeout)
	}

	if verbose {
		testOptions = testOptions.WithVerboseOutput()
	}

	if err := testOptions.Validate(); err != nil {
		return nil, WrapErrorf(err, "invalid test options")
	}

	if err := buildOpts.Validate(); err != nil {
		return nil, WrapErrorf(err, "invalid build options")
	}

	cmdToAppend = append(cmdToAppend, buildOpts.Flags...)
	cmdToAppend = append(cmdToAppend, testOptions.Flags...)

	return m.
		Ctr.
		WithExec(cmdToAppend), nil
}

// RunTestCmd executes the Go test command with the specified options and parameters.
//
// This function allows for a comprehensive configuration of the test command,
// including build options, test options, and environment settings. It takes a
// context for managing cancellation and timeouts, a source directory for the
// code to be tested, and various flags to customize the behavior of the test
// execution. The function returns the output of the command as a string and
// any error encountered during execution.
//
// Parameters:
//   - ctx: The context to run the command.
//     +optional
//   - source: The source code to test.
//   - packages: The packages to test.
//     +optional
//   - envVars: The environment variables to set.
//     +optional
//   - secrets: The secrets to set.
//     +optional
//   - enableDefaultOptions: Enables the default options for the test command.
//     +optional
//   - race: Enables the race detector in the Go command. It's equivalent to the -race flag.
//     +optional
//   - msan: Enables memory sanitizer in the Go command. It's equivalent to the -msan flag.
//     +optional
//   - asan: Enables address sanitizer in the Go command. It's equivalent to the -asan flag.
//     +optional
//   - buildTags: Specifies build constraints for the Go command. It's equivalent to the -tags flag.
//     +optional
//   - ldflags: Sets flags for the linker in the Go command. It's equivalent to the -ldflags flag.
//     +optional
//   - gcflags: Sets flags for the Go compiler. It's equivalent to the -gcflags flag.
//     +optional
//   - asmflags: Sets flags for the assembler in the Go command. It's equivalent to the -asmflags flag.
//     +optional
//   - trimpath: Removes all file system paths from the compiled binary. It's equivalent to the -trimpath flag.
//     +optional
//   - work: Enables the creation of a temporary work directory. It's equivalent to the -work flag.
//     +optional
//   - buildMode: Specifies the build mode for the Go command. It's equivalent to the -buildmode flag.
//     +optional
//   - compiler: Specifies the compiler to use for building. It's equivalent to the -compiler flag.
//     +optional
//   - gccgoflags: Sets flags for the gccgo compiler. It's equivalent to the -gccgoflags flag.
//     +optional
//   - mod: Specifies the module mode for the Go command. It's equivalent to the -mod flag.
//     +optional
//   - benchmark: Specifies the benchmark to run. It's equivalent to the -bench flag.
//     +optional
//   - benchmem: Enables memory allocation statistics. It's equivalent to the -benchmem flag.
//     +optional
//   - benchtime: Specifies the duration for benchmarks. It's equivalent to the -benchtime flag.
//     +optional
//   - blockprofile: Specifies the file for block profiling. It's equivalent to the -blockprofile flag.
//     +optional
//   - cover: Enables coverage analysis. It's equivalent to the -cover flag.
//     +optional
//   - coverprofile: Specifies the file for coverage profile output. It's equivalent to the -coverprofile flag.
//     +optional
//   - cpuprofile: Specifies the file for CPU profiling. It's equivalent to the -cpuprofile flag.
//     +optional
//   - testCount: Specifies the number of test iterations. It's equivalent to the -count flag.
//     +optional
//   - failfast: Stops the test run on the first failure. It's equivalent to the -failfast flag.
//     +optional
//   - enableJsonOutput: Enables JSON output for test results. It's equivalent to the -json flag.
//     +optional
//   - list: Specifies a regex to filter tests. It's equivalent to the -list flag.
//     +optional
//   - memprofile: Specifies the file for memory profiling. It's equivalent to the -memprofile flag.
//     +optional
//   - mutexprofile: Specifies the file for mutex profiling. It's equivalent to the -mutexprofile flag.
//     +optional
//   - parallel: Specifies the maximum number of tests to run in parallel. It's equivalent to the -parallel flag.
//     +optional
//   - run: Specifies a regex to select tests to run. It's equivalent to the -run flag.
//     +optional
//   - short: Enables short test mode. It's equivalent to the -short flag.
//     +optional
//   - timeout: Specifies the maximum time to run tests. It's equivalent to the -timeout flag.
//     +optional
//   - verbose: Is the flag that sets the verbosity level in the Go command. It's equivalent to the -v flag.
//     +optional
//
//nolint:funlen // It's okay to have this size, it's by design.
func (m *Gotest) RunTestCmd(
	// source is the source code to test.
	source *dagger.Directory,
	// packages are the packages to test.
	// +optional
	packages []string,
	// envVars are the environment variables to set.
	// +optional
	envVars []string,
	// secrets are the secrets to set.
	// +optional
	secrets []*dagger.Secret,
	// enableDefaultOptions enables the default options for the test command.
	// +optional
	enableDefaultOptions bool,
	// Build Options
	// race enables the race detector in the Go command.
	// It's equivalent to the -race flag.
	// +optional
	race bool,
	// msan enables memory sanitizer in the Go command.
	// It's equivalent to the -msan flag.
	// +optional
	msan bool,
	// asan enables address sanitizer in the Go command.
	// It's equivalent to the -asan flag.
	// +optional
	asan bool,
	// buildTags specifies build constraints for the Go command.
	// It's equivalent to the -tags flag.
	// +optional
	buildTags string,
	// ldflags sets flags for the linker in the Go command.
	// It's equivalent to the -ldflags flag.
	// +optional
	ldflags string,
	// gcflags sets flags for the Go compiler.
	// It's equivalent to the -gcflags flag.
	// +optional
	gcflags string,
	// asmflags sets flags for the assembler in the Go command.
	// It's equivalent to the -asmflags flag.
	// +optional
	asmflags string,
	// trimpath removes all file system paths from the compiled binary.
	// It's equivalent to the -trimpath flag.
	// +optional
	trimpath bool,
	// work enables the creation of a temporary work directory.
	// It's equivalent to the -work flag.
	// +optional
	work bool,
	// buildMode specifies the build mode for the Go command.
	// It's equivalent to the -buildmode flag.
	// +optional
	buildMode string,
	// compiler specifies the compiler to use for building.
	// It's equivalent to the -compiler flag.
	// +optional
	compiler string,
	// gccgoflags sets flags for the gccgo compiler.
	// It's equivalent to the -gccgoflags flag.
	// +optional
	gccgoflags string,
	// mod specifies the module mode for the Go command.
	// It's equivalent to the -mod flag.
	// +optional
	mod string,
	// Test Options
	// benchmark specifies the benchmark to run.
	// It's equivalent to the -bench flag.
	// +optional
	benchmark string,
	// benchmem enables memory allocation statistics.
	// It's equivalent to the -benchmem flag.
	// +optional
	benchmem bool,
	// benchtime specifies the duration for benchmarks.
	// It's equivalent to the -benchtime flag.
	// +optional
	benchtime string,
	// blockprofile specifies the file for block profiling.
	// It's equivalent to the -blockprofile flag.
	// +optional
	blockprofile string,
	// cover enables coverage analysis.
	// It's equivalent to the -cover flag.
	// +optional
	cover bool,
	// coverprofile specifies the file for coverage profile output.
	// It's equivalent to the -coverprofile flag.
	// +optional
	coverprofile string,
	// cpuprofile specifies the file for CPU profiling.
	// It's equivalent to the -cpuprofile flag.
	// +optional
	cpuprofile string,
	// testCount specifies the number of test iterations.
	// It's equivalent to the -count flag.
	// +optional
	testCount int,
	// failfast stops the test run on the first failure.
	// It's equivalent to the -failfast flag.
	// +optional
	failfast bool,
	// enableJSONOutput enables JSON output for test results.
	// It's equivalent to the -json flag.
	// +optional
	enableJSONOutput bool,
	// list specifies a regex to filter tests.
	// It's equivalent to the -list flag.
	// +optional
	list string,
	// memprofile specifies the file for memory profiling.
	// It's equivalent to the -memprofile flag.
	// +optional
	memprofile string,
	// mutexprofile specifies the file for mutex profiling.
	// It's equivalent to the -mutexprofile flag.
	// +optional
	mutexprofile string,
	// parallel specifies the maximum number of tests to run in parallel.
	// It's equivalent to the -parallel flag.
	// +optional
	parallel int,
	// run specifies a regex to select tests to run.
	// It's equivalent to the -run flag.
	// +optional
	run string,
	// short enables short test mode.
	// It's equivalent to the -short flag.
	// +optional
	short bool,
	// timeout specifies the maximum time to run tests.
	// It's equivalent to the -timeout flag.
	// +optional
	timeout string,
	// verbose is the flag that sets the verbosity level in the Go command.
	// It's equivalent to the -v flag.
	// +optional
	verbose bool,
) (string, error) {
	ctr, err := m.RunTest(
		source,
		packages,
		envVars,
		secrets,
		enableDefaultOptions,
		race,
		msan,
		asan,
		buildTags,
		ldflags,
		gcflags,
		asmflags,
		trimpath,
		work,
		buildMode,
		compiler,
		gccgoflags,
		mod,
		benchmark,
		benchmem,
		benchtime,
		blockprofile,
		cover,
		coverprofile,
		cpuprofile,
		testCount,
		failfast,
		enableJSONOutput,
		list,
		memprofile,
		mutexprofile,
		parallel,
		run,
		short,
		timeout,
		verbose,
	)

	if err != nil {
		return "", WrapErrorf(err, "failed to run Go test command")
	}

	stdout, err := ctr.Stdout(context.Background())

	if err != nil {
		return "", WrapErrorf(err, "failed to get stdout from Go test command")
	}

	return stdout, nil
}
