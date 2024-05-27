package main

import (
	"path/filepath"

	"github.com/Excoriate/daggerx/pkg/fixtures"
)

// WithSource sets the source directory if it's passed, and
// mounts the source directory to the container.
func (m *Terragrunt) WithSource(src *Directory, workdir string) *Terragrunt {
	if src != nil {
		m.Src = src
	}

	var workDirPath string
	if workdir == "" {
		workDirPath = fixtures.MntPrefix
	} else {
		workDirPath = filepath.Join(fixtures.MntPrefix, workdir)
	}

	m.Ctr = m.Ctr.
		WithMountedDirectory(fixtures.MntPrefix, m.Src).
		WithoutWorkdir().
		WithWorkdir(workDirPath)

	return m
}

// WithTerragruntCache creates a cache volume for Terragrunt for the .terragrunt-cache folder
func (m *Terragrunt) WithTerragruntCache() *Terragrunt {
	tgCache := dag.CacheVolume(".terragrunt-cache")
	ctr := m.Ctr.WithMountedCache(".terragrunt-cache", tgCache)

	m.Ctr = ctr

	return m
}

// WithEnvVar adds an environment variable to the container.
func (m *Terragrunt) WithEnvVar(name, value string, expand bool) *Terragrunt {
	m.Ctr = m.Ctr.WithEnvVariable(name, value, ContainerWithEnvVariableOpts{
		Expand: expand,
	})

	return m
}

// WithCMD sets the command to run in the container.
func (m *Terragrunt) WithCMD(cmd []string) *Terragrunt {
	m.Ctr = m.Ctr.WithFocus().WithExec(cmd)

	return m
}
