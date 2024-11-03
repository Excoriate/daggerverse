package main

import (
	"strconv"
	"strings"
)

// GoTestOptions holds the options for running Go tests.
//
// This struct contains the flags to be passed to the Go test command.
// It supports various flags to control the behavior of the tests.
type GoTestOptions struct {
	// Flags are the flags to pass to the Go test command.
	// +private
	Flags []string
}

// NewGoTestOptions creates a new GoTestOptions instance.
//
// This function initializes a new GoTestOptions struct with default values.
// It returns a pointer to the newly created struct.
func NewGoTestOptions() *GoTestOptions {
	return &GoTestOptions{
		Flags: []string{},
	}
}

// WithBenchmark adds a regular expression to select benchmarks to run.
//
// This function adds the -bench flag with the specified regexp to the test options.
// It selects benchmarks to run based on the provided regular expression.
func (o *GoTestOptions) WithBenchmark(regexp string) *GoTestOptions {
	o.Flags = append(o.Flags, "-bench", regexp)

	return o
}

// WithBenchmarkMemory reports memory allocations for benchmarks.
//
// This function adds the -benchmem flag to the test options.
// It provides detailed information about memory allocations during benchmarks.
func (o *GoTestOptions) WithBenchmarkMemory() *GoTestOptions {
	o.Flags = append(o.Flags, "-benchmem")

	return o
}

// WithBenchmarkTime sets the time to run each benchmark.
//
// This function adds the -benchtime flag with the specified duration to the test options.
// It controls the amount of time spent on benchmarking for accurate measurements.
func (o *GoTestOptions) WithBenchmarkTime(duration string) *GoTestOptions {
	o.Flags = append(o.Flags, "-benchtime", duration)

	return o
}

// WithBlockProfile enables block profiling.
//
// This function adds the -blockprofile flag with the specified file to the test options.
// It helps identify issues related to synchronization primitives.
func (o *GoTestOptions) WithBlockProfile(file string) *GoTestOptions {
	o.Flags = append(o.Flags, "-blockprofile", file)

	return o
}

// WithCoverage enables coverage analysis.
//
// This function adds the -cover flag to the test options.
// It helps identify which parts of the code are not being tested.
func (o *GoTestOptions) WithCoverage() *GoTestOptions {
	o.Flags = append(o.Flags, "-cover")

	return o
}

// WithCoverageProfile sets the file for the coverage profile.
//
// This function adds the -coverprofile flag with the specified file to the test options.
// It specifies where the coverage profile data will be written.
func (o *GoTestOptions) WithCoverageProfile(file string) *GoTestOptions {
	o.Flags = append(o.Flags, "-coverprofile", file)

	return o
}

// WithCPUProfile enables CPU profiling.
//
// This function adds the -cpuprofile flag with the specified file to the test options.
// It helps identify performance bottlenecks and areas of the code consuming excessive CPU resources.
func (o *GoTestOptions) WithCPUProfile(file string) *GoTestOptions {
	o.Flags = append(o.Flags, "-cpuprofile", file)

	return o
}

// WithTestCount sets the number of times to run each test.
//
// This function adds the -count flag with the specified number to the test options.
// It controls the number of test runs for each test.
func (o *GoTestOptions) WithTestCount(count int) *GoTestOptions {
	o.Flags = append(o.Flags, "-count", strconv.Itoa(count))

	return o
}

// WithFailFast stops running tests after the first failure.
//
// This function adds the -failfast flag to the test options.
// It stops the test execution process as soon as a failure is detected.
func (o *GoTestOptions) WithFailFast() *GoTestOptions {
	o.Flags = append(o.Flags, "-failfast")

	return o
}

// WithJSONOutput enables JSON output.
//
// This function adds the -json flag to the test options.
// It enables JSON output for your tests.
func (o *GoTestOptions) WithJSONOutput() *GoTestOptions {
	o.Flags = append(o.Flags, "-json")

	return o
}

// WithListTests lists tests, benchmarks, or examples matching the regular expression.
//
// This function adds the -list flag with the specified regexp to the test options.
// It lists matching tests based on the provided regular expression.
func (o *GoTestOptions) WithListTests(regexp string) *GoTestOptions {
	o.Flags = append(o.Flags, "-list", regexp)

	return o
}

// WithMemoryProfile enables memory profiling.
//
// This function adds the -memprofile flag with the specified file to the test options.
// It helps identify memory leaks, excessive memory allocations, and other memory-related issues.
func (o *GoTestOptions) WithMemoryProfile(file string) *GoTestOptions {
	o.Flags = append(o.Flags, "-memprofile", file)

	return o
}

// WithMutexProfile enables mutex profiling.
//
// This function adds the -mutexprofile flag with the specified file to the test options.
// It helps identify issues related to synchronization primitives.
func (o *GoTestOptions) WithMutexProfile(file string) *GoTestOptions {
	o.Flags = append(o.Flags, "-mutexprofile", file)

	return o
}

// WithParallelTests sets the number of parallel test executions.
//
// This function adds the -parallel flag with the specified number to the test options.
// It controls the number of tests that can be run in parallel.
func (o *GoTestOptions) WithParallelTests(count int) *GoTestOptions {
	o.Flags = append(o.Flags, "-parallel", strconv.Itoa(count))

	return o
}

// WithTestFilter adds a regular expression to select tests and examples to run.
//
// This function adds the -run flag with the specified regexp to the test options.
// It allows you to filter and run specific tests or examples based on their names or descriptions.
func (o *GoTestOptions) WithTestFilter(regexp string) *GoTestOptions {
	o.Flags = append(o.Flags, "-run", regexp)

	return o
}

// WithShortTest enables running smaller test suites.
//
// This function adds the -short flag to the test options.
// It enables shorter test runs.
func (o *GoTestOptions) WithShortTest() *GoTestOptions {
	o.Flags = append(o.Flags, "-short")

	return o
}

// WithTimeout sets the timeout for each test.
//
// This function adds the -timeout flag with the specified duration to the test options.
// It controls the maximum duration for which a test can run.
func (o *GoTestOptions) WithTimeout(duration string) *GoTestOptions {
	o.Flags = append(o.Flags, "-timeout", duration)

	return o
}

// WithVerboseOutput enables verbose output.
//
// This function adds the -v flag to the test options.
// It enables verbose output for more detailed information about the test execution process.
func (o *GoTestOptions) WithVerboseOutput() *GoTestOptions {
	o.Flags = append(o.Flags, "-v")

	return o
}

// ErrInvalidTestFlag is returned when an invalid test flag is encountered.
// var ErrInvalidTestFlag = errors.New("invalid test flag for command go test")

// Validate checks if the test options are valid.
//
// This function performs validation on the test options to ensure they are correctly formed.
// It returns an error if any issues are found, otherwise it returns nil.
func (o *GoTestOptions) Validate() error {
	for i := 0; i < len(o.Flags); i++ {
		flag := o.Flags[i]
		if strings.HasPrefix(flag, "-") && len(flag) > 1 {
			if i+1 < len(o.Flags) && !strings.HasPrefix(o.Flags[i+1], "-") {
				// Skip the next item as it's a value for this flag
				i++
			}
		} else {
			return Errorf("invalid test flag for command go test: %s", flag)
		}
	}

	return nil
}
