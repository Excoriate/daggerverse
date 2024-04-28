package main

import (
	"fmt"
	"path/filepath"
)

// Base sets the base image and version, and creates the base container.
// For Terragrunt, the default image is "alpine/terragrunt" and the default version is "1.7.0".
// Consider that the container isn't created based on the Terragrunt version, but on the Terraform version.
func (m *Terragrunt) Base(image, version string) *Terragrunt {
	if image == "" {
		image = tgContainerImageDefault
	}

	if version == "" {
		version = tfVersionDefault
	}

	ctrImage := fmt.Sprintf("%s:%s", image, version)

	c := dag.Container().From(ctrImage)

	m.Ctr = c

	return m
}

// WithSource sets the source directory if it's passed, and
// mounts the source directory to the container.
func (m *Terragrunt) WithSource(src *Directory, workdir string) *Terragrunt {
	if src != nil {
		m.Src = src
	}

	var workDirPath string
	if workdir == "" {
		workDirPath = workdirRootPath
	} else {
		workDirPath = filepath.Join(workdirRootPath, workdir)
	}

	m.Ctr = m.Ctr.
		WithMountedDirectory(workdirRootPath, m.Src).
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

// WithEnvVars sets the environment variables to the container.
func (m *Terragrunt) WithEnvVars(envVars []EnvVarDagger) *Terragrunt {
	for _, envVar := range envVars {
		m.Ctr = m.WithEnvVar(envVar.Name, envVar.Value, envVar.Expand).Ctr
	}

	return m
}
