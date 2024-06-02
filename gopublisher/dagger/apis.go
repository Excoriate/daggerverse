package main

import (
	"fmt"

	"github.com/Excoriate/daggerx/pkg/envvars"
	"github.com/Excoriate/daggerx/pkg/fixtures"
)

// WithEnvVariable sets an environment variable.
//
// The name of the environment variable (e.g., "HOST").
//
// The value of the environment variable (e.g., "localhost").
//
// Replace `${VAR}` or $VAR in the value according to the current environment
// variables defined in the container (e.g., "/opt/bin:$PATH").
// +optional
func (m *Gopublisher) WithEnvVariable(
	// The name of the environment variable (e.g., "HOST").
	name string,

	// The value of the environment variable (e.g., "localhost").
	value string,

	// Replace `${VAR}` or $VAR in the value according to the current environment
	// variables defined in the container (e.g., "/opt/bin:$PATH").
	// +optional
	expand bool,
) *Gopublisher {
	m.Ctr = m.Ctr.WithEnvVariable(name, value, ContainerWithEnvVariableOpts{
		Expand: expand,
	})

	return m
}

// WithCGODisabled sets the CGO_ENABLED environment variable to 0.
//
// It sets the CGO_ENABLED environment variable to 0.
func (m *Gopublisher) WithCGODisabled() *Gopublisher {

	m.Ctr = m.WithEnvVariable("CGO_ENABLED", "0", false).Ctr
	return m
}

// WithGit installs or setup the container with the git binary.
//
// It installs or sets up the container with the git binary.
func (m *Gopublisher) WithGit() *Gopublisher {
	m.Ctr = m.Ctr.WithExec([]string{"apk", "add", "--no-cache", "git"})

	return m
}

// WithCURL installs or setup the container with the curl binary.
//
// It installs or sets up the container with the curl binary.
func (m *Gopublisher) WithCURL() *Gopublisher {
	m.Ctr = m.Ctr.WithExec([]string{"apk", "add", "--no-cache", "curl"})

	return m
}

// WithSource sets the source directory.
//
// The source directory is the directory that contains all the source code, including the module directory.
func (m *Gopublisher) WithSource(
	// src is the directory that contains all the source code, including the module directory.
	src *Directory,
) *Gopublisher {
	m.Src = src

	m.Ctr = m.Ctr.WithWorkdir(fixtures.MntPrefix).
		WithMountedDirectory(fixtures.MntPrefix, m.Src)

	return m
}

// WithEnvVarsFromStrs sets the environment variables for the container.
//
// It sets the environment variables for the container. It's meant to be used as a terminal for the module.
// Arguments:
// - envVars: A set of strings (e.g., "KEY=value,KEY=value") to use as environment variables. They're
// used to set the environment variables for the container when it's required to pass multiple environment variables
// in a single argument. E.g.: "GITHUB=token,GO_VERSION=1.22.0,MYVAR=myvar"
func (m *Gopublisher) WithEnvVarsFromStrs(envVars []string) (*Gopublisher, error) {
	envVarsParsed, err := envvars.ToDaggerEnvVarsFromSlice(envVars)
	if err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	for _, envVar := range envVarsParsed {
		m.Ctr = m.WithEnvVariable(envVar.Name, envVar.Value, false).Ctr
	}

	return m, nil
}
