package main

import (
	"github.com/Excoriate/daggerx/pkg/fixtures"
	"github.com/Excoriate/daggerx/pkg/golangx"
	"path/filepath"
)

// WithSource Set the source directory.
func (m *Terratest) WithSource(
	// Src is the directory that contains all the source code, including the module directory.
	src *Directory,
	// workdir is the working directory.
	// +optional
	workdir string,
) *Terratest {
	ctr := m.Ctr.WithMountedDirectory(fixtures.MntPrefix, src)

	if workdir != "" {
		ctr = ctr.WithWorkdir(filepath.Join(fixtures.MntPrefix, workdir))
	} else {
		ctr = ctr.WithWorkdir(fixtures.MntPrefix)
	}

	m.Ctr = ctr

	return m
}

// WithCgoDisabled Set CGO_ENABLED environment variable to 0.
func (m *Terratest) WithCgoDisabled() *Terratest {
	gox := golangx.WithGoCgoDisabled()
	m.Ctr = m.Ctr.WithEnvVariable(gox.Name, gox.Value)
	return m
}

// WithEnvVar Set an environment variable.
func (m *Terratest) WithEnvVar(
	// The name of the environment variable (e.g., "HOST").
	name string,

	// The value of the environment variable (e.g., "localhost").
	value string,

	// Replace `${VAR}` or $VAR in the value according to the current environment
	// variables defined in the container (e.g., "/opt/bin:$PATH").
	// +optional
	expand bool,
) *Terratest {
	m.Ctr = m.Ctr.WithEnvVariable(name, value, ContainerWithEnvVariableOpts{
		Expand: expand,
	})

	return m
}
