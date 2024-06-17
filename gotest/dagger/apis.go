//nolint:nolintlint,revive // This is a method that is used to set the base image and version.
package main

import (
	"fmt"
	"path/filepath"

	"github.com/Excoriate/daggerx/pkg/fixtures"
	"github.com/Excoriate/daggerx/pkg/golangx"
	"github.com/containerd/containerd/platforms"
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

	platformSet := platforms.MustParse(string(platform))

	ctr := m.Ctr

	ctr = ctr.
		WithEnvVariable("GOOS", platformSet.OS).
		WithEnvVariable("GOARCH", platformSet.Architecture)

	if platformSet.Variant != "" {
		ctr = ctr.WithEnvVariable("GOARM", platformSet.Variant)
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

// WithEnvironmentVariable Set an environment variable.
func (m *Gotest) WithEnvironmentVariable(
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
func (m *Gotest) WithModuleCache() *Gotest {
	goModCache := dag.CacheVolume("godmodcache")

	m.Ctr = m.Ctr.WithMountedCache("/go/pkg/mod", goModCache).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod")

	return m
}

// WithBuildCache sets the build cache for the Go module.
// The default cache volume is "gobuildcache", and the default mount path is "/go/build-cache".
func (m *Gotest) WithBuildCache() *Gotest {
	goBuildCache := dag.CacheVolume("gobuildcache")

	m.Ctr = m.Ctr.WithMountedCache("/go/build-cache", goBuildCache).
		WithEnvVariable("GOCACHE", "/go/build-cache")

	return m
}

// WithGoCache mounts the Go cache directories.
func (m *Gotest) WithGoCache() *Gotest {
	return m.WithModuleCache().
		WithBuildCache()
}

// WithNewNetrcFileGitHub creates a new .netrc file with the GitHub credentials.
//
// The .netrc file is created in the root directory of the container.
func (m *Gotest) WithNewNetrcFileGitHub(username, password string) *Gotest {
	machineCMD := fmt.Sprintf("machine github.com\nlogin %s\npassword %s\n", username, password)

	//nolint:exhaustruct // This is a method that is used to set the base image and version.
	m.Ctr = m.Ctr.WithNewFile("/root/.netrc", ContainerWithNewFileOpts{
		Contents: machineCMD,
	})

	return m
}

// WithNewNetrcFileGitLab creates a new .netrc file with the GitLab credentials.
//
// The .netrc file is created in the root directory of the container.
func (m *Gotest) WithNewNetrcFileGitLab(username, password string) *Gotest {
	machineCMD := fmt.Sprintf("machine gitlab.com\nlogin %s\npassword %s\n", username, password)

	//nolint:exhaustruct // This is a method that is used to set the base image and version.
	m.Ctr = m.Ctr.WithNewFile("/root/.netrc", ContainerWithNewFileOpts{
		Contents: machineCMD,
	})

	return m
}

// WithPrivateGoPkg sets the GOPRIVATE environment variable.
//
//nolint:lll    // The GOPRIVATE environment variable is used to specify a comma-separated list of hosts for which Go modules should always be fetched directly from VCS repositories.
//nolint:exhaustruct // This is a method that is used to set the base image and version.
func (m *Gotest) WithPrivateGoPkg(privateHost string) *Gotest {
	//nolint:exhaustruct // This is a method that is used to set the base image and version.
	m.Ctr = m.Ctr.WithExec([]string{"go", "env", "GOPRIVATE", privateHost}, ContainerWithExecOpts{
		InsecureRootCapabilities: true,
	}).WithEnvVariable("GOPRIVATE", privateHost)

	return m
}

// WithGCCCompiler installs the GCC compiler and musl-dev package.
func (m *Gotest) WithGCCCompiler() *Gotest {
	m.Ctr = m.Ctr.WithExec([]string{"apk", "add", "--no-cache", "gcc", "musl-dev"})

	return m
}

// WithGoTestSum installs the gotestsum CLI.
func (m *Gotest) WithGoTestSum() *Gotest {
	goTestSumInstallCMD := []string{"go", "install", "gotest.tools/gotestsum@latest"}
	goTestInstallTparseCMD := []string{"go", "install", "github.com/mfridman/tparse@latest"}

	m.Ctr = m.Ctr.WithExec(goTestSumInstallCMD).WithExec(goTestInstallTparseCMD)

	return m
}
