package main

import (
	"context"

	"github.com/Excoriate/daggerverse/gotoolbox/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/envvars"
)

// RunGo runs a Go command within a given context.
//
// cmd is the Go command to run, everything after the 'go' command.
// src is the optional source directory for the container.
//
// It returns the standard output of the executed command or an error if something goes wrong.
func (m *Gotoolbox) RunGo(
	// cmd is the Go command to run, everything after the 'go' command.
	cmd []string,
	// src is the optional source directory for the container.
	// +optional
	src *dagger.Directory,
	// testDir is the optional test directory for the container. If you are
	// passing a source code directory in the 'src' parameter, you can pass
	// the test directory to run the tests.
	// +optional
	testDir string,
	// platform is the optional platform to run the Go command.
	// +optional
	platform dagger.Platform,
	// envVariables are the optional environment variables to set for the container.
	// the envVariables are set in a form of a string array, where each element is a key-value pair,
	// separated by coma. E.g. ["KEY=value", "KEY2=value2"].
	// +optional
	envVariables []string,
	// enableCacheBuster is a flag to force the cache to bust.
	// +optional
	enableCacheBuster bool,
	// enableCgo is a flag to enable CGO.
	// +optional
	enableCgo bool,
	// enableGoModCache is a flag to enable Go mod cache.
	// +optional
	enableGoModCache bool,
	// enableGoBuildCache is a flag to enable Go build cache.
	// +optional
	enableGoBuildCache bool,
	// installPkgs are the packages to install.
	// +optional
	installPkgs []string,
	// enableGoGCCCompiler is a flag to enable GoGCCCompiler.
	// +optional
	enableGoGCCCompiler bool,
) (string, error) {
	// Check if the first element of cmd is "go" and remove it if present
	// It's not requires, since the WithGoExec command is used to run the Go command.
	if len(cmd) > 0 && cmd[0] == "go" {
		cmd = cmd[1:]
	}

	if src == nil {
		src = dag.
			CurrentModule().
			Source()
	}

	m = m.WithSource(src, testDir)

	if len(envVariables) > 0 {
		envVars, err := envvars.ToDaggerEnvVarsFromSlice(envVariables)
		if err != nil {
			return "", WrapErrorf(err, "failed to parse environment variables %s", envVariables)
		}

		for _, envVar := range envVars {
			m = m.WithEnvironmentVariable(envVar.Name, envVar.Value, false)
		}
	}

	if enableCacheBuster {
		m = m.WithCacheBuster()
	}

	if enableCgo {
		m = m.WithGoCgoEnabled()
	}

	if enableGoModCache {
		m = m.WithGoModCache("", nil, nil, "")
	}

	if enableGoBuildCache {
		m = m.WithGoBuildCache("", nil, nil, "")
	}

	if len(installPkgs) > 0 {
		m = m.WithGoInstall(installPkgs)
	}

	if enableGoGCCCompiler {
		m = m.WithGoGCCCompiler()
	}

	ctx := context.Background()

	if platform != "" {
		m = m.WithGoExec(cmd, platform)
	} else {
		m = m.WithGoExec(cmd, "")
	}

	ctrExecStdOut, ctrExecErr := m.Ctr.
		Stdout(ctx)

	if ctrExecErr != nil {
		return "", WrapError(ctrExecErr, "failed to execute 'RunGo' function")
	}

	return ctrExecStdOut, nil
}
