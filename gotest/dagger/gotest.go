package main

import (
	"fmt"
	"github.com/Excoriate/daggerx/pkg/envvars"
	"github.com/Excoriate/daggerx/pkg/fixtures"
	"github.com/Excoriate/daggerx/pkg/golangx"
	"github.com/containerd/containerd/platforms"
	"path/filepath"
)

// WithSource Set the source directory.
func (m *Gotest) WithSource(
	// Src is the directory that contains all the source code, including the module directory.
	src *Directory,
	// workdir is the working directory.
	// +optional
	workdir string,
) *Gotest {
	m.Src = src
	ctr := m.Ctr.WithMountedDirectory(fixtures.MntPrefix, m.Src)

	if workdir != "" {
		ctr = ctr.WithWorkdir(filepath.Join(fixtures.MntPrefix, workdir))
	} else {
		ctr = ctr.WithWorkdir(fixtures.MntPrefix)
	}

	m.Ctr = ctr

	return m
}

// WithPlatform Set GOOS, GOARCH and GOARM environment variables.
func (m *Gotest) WithPlatform(
	// Target platform in "[os]/[platform]/[version]" format (e.g., "darwin/arm64/v7", "windows/amd64", "linux/arm64").
	platform Platform,
) *Gotest {
	if platform == "" {
		return m
	}

	p := platforms.MustParse(string(platform))

	ctr := m.Ctr

	ctr = ctr.
		WithEnvVariable("GOOS", p.OS).
		WithEnvVariable("GOARCH", p.Architecture)

	if p.Variant != "" {
		ctr = ctr.WithEnvVariable("GOARM", p.Variant)
	}

	m.Ctr = ctr
	return m
}

// WithCgoEnabled Set CGO_ENABLED environment variable to 1.
func (m *Gotest) WithCgoEnabled() *Gotest {
	gox := golangx.WithGoCgoEnabled()
	m.Ctr = m.Ctr.WithEnvVariable(gox.Name, gox.Value)
	return m
}

// WithCgoDisabled Set CGO_ENABLED environment variable to 0.
func (m *Gotest) WithCgoDisabled() *Gotest {
	gox := golangx.WithGoCgoDisabled()
	m.Ctr = m.Ctr.WithEnvVariable(gox.Name, gox.Value)
	return m
}

// WithEnvVariable Set an environment variable.
func (m *Gotest) WithEnvVariable(
	// The name of the environment variable (e.g., "HOST").
	name string,

	// The value of the environment variable (e.g., "localhost").
	value string,

	// Replace `${VAR}` or $VAR in the value according to the current environment
	// variables defined in the container (e.g., "/opt/bin:$PATH").
	// +optional
	expand bool,
) *Gotest {
	m.Ctr = m.Ctr.WithEnvVariable(name, value, ContainerWithEnvVariableOpts{
		Expand: expand,
	})

	return m
}

// WithModuleCache sets the module cache for the Go module.
// The default cache volume is "godmodcache", and the default mount path is "/go/pkg/mod".
func (m *Gotest) WithModuleCache(ctr *Container) *Container {
	goModCache := dag.CacheVolume("godmodcache")

	ctr = m.Ctr.WithMountedCache("/go/pkg/mod", goModCache).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod")

	return ctr
}

// WithBuildCache sets the build cache for the Go module.
// The default cache volume is "gobuildcache", and the default mount path is "/go/build-cache".
func (m *Gotest) WithBuildCache(ctr *Container) *Container {
	goBuildCache := dag.CacheVolume("gobuildcache")

	ctr = m.Ctr.WithMountedCache("/go/build-cache", goBuildCache).
		WithEnvVariable("GOCACHE", "/go/build-cache")

	return ctr
}

// WithGoCache mounts the Go cache directories.
func (m *Gotest) WithGoCache(ctr *Container) *Container {
	ctr = m.WithModuleCache(ctr)
	ctr = m.WithBuildCache(ctr)

	return ctr
}

// SetupGoTest sets up the go test options, to either evaluate the container and run the test,
// or return the container to be evaluated later.
func (m *Gotest) SetupGoTest(
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
	// Enable experimental Dagger nesting. It sets the ExperimentalPrivilegedNesting option in Dagger.
	// +optional
	enableNest bool,
	// enableCache is a flag to enable cache volumes. If it's set, it'll
	// enable the cache volumes for the Go module and build cache.
	// +optional
	enableCache bool,
	// envVars is a list of environment variables to set in the container with the format "SOMETHING=SOMETHING,SOMETHING=SOMETHING".
	// +optional
	envVars []string,
	// printEnvVars is a flag to print the environment variables
	// +optional
	printEnvVars bool,
) (*Container, error) {
	goTest := []string{"go", "test"}
	ctr := m.WithSource(src, "").Ctr

	if enableCache {
		ctr = m.Ctr.With(m.WithGoCache)
	}

	if len(envVars) > 0 {
		envVarsDagger, err := envvars.ToDaggerEnvVarsFromSlice(envVars)
		if err != nil {
			return nil, err
		}

		for _, envVar := range envVarsDagger {
			ctr = m.WithEnvVariable(envVar.Name, envVar.Value, false).Ctr
		}
	}

	if printEnvVars {
		ctr = ctr.WithExec([]string{"printenv"}, ContainerWithExecOpts{
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

	ctr = ctr.WithExec(goTest, ContainerWithExecOpts{
		InsecureRootCapabilities:      insecureRootCapabilities,
		ExperimentalPrivilegedNesting: enableNest,
	})

	m.Ctr = ctr

	return m.Ctr, nil
}

// SetupGoTestSum sets up the go test options, to either evaluate the container and run the test,
// or return the container to be evaluated later.
func (m *Gotest) SetupGoTestSum(
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
) (*Container, error) {
	goTestSumInstallCMD := []string{"go", "install", "gotest.tools/gotestsum@latest"}
	goTestInstallTparseCMD := []string{"go", "install", "github.com/mfridman/tparse@latest"}
	goTestCMD := []string{"gotestsum", "--no-color=false"}
	ctr := m.WithSource(src, "").Ctr

	if enableCache {
		ctr = m.Ctr.With(m.WithGoCache)
	}

	if len(envVars) > 0 {
		envVarsDagger, err := envvars.ToDaggerEnvVarsFromSlice(envVars)
		if err != nil {
			return nil, err
		}

		for _, envVar := range envVarsDagger {
			ctr = m.WithEnvVariable(envVar.Name, envVar.Value, false).Ctr
		}
	}

	if printEnvVars {
		ctr = ctr.WithFocus().
			WithExec([]string{"printenv"}, ContainerWithExecOpts{
				InsecureRootCapabilities:      insecureRootCapabilities,
				ExperimentalPrivilegedNesting: enableNest,
			})
	}

	ctr = ctr.WithExec(goTestSumInstallCMD)

	if format == "" {
		format = "testname" // opinionated default
	}

	goTestCMD = append(goTestCMD, fmt.Sprintf("--format=%s", format))
	goTestCMD = append(goTestCMD, goTestSumFlags...)

	if race {
		goTestCMD = append(goTestCMD, "-race")
	}

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

	ctr = ctr.WithExec(goTestCMD, ContainerWithExecOpts{
		InsecureRootCapabilities:      insecureRootCapabilities,
		ExperimentalPrivilegedNesting: enableNest,
	})

	if enablePretty {
		tParseCMD := []string{"tparse", "-all", "-smallscreen", "-file=test-output.json"}
		ctr = ctr.WithExec(tParseCMD, ContainerWithExecOpts{
			InsecureRootCapabilities:      insecureRootCapabilities,
			ExperimentalPrivilegedNesting: enableNest,
		})
	}

	m.Ctr = ctr

	return m.Ctr, nil
}
