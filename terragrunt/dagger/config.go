package main

import "fmt"

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

	c := dag.Container().From(ctrImage).
		WithWorkdir(workdirRootPath)

	m.Ctr = c

	return m
}

// WithSource sets the source directory if it's passed, and
// mounts the source directory to the container.
func (m *Terragrunt) WithSource(src *Directory) *Terragrunt {
	if src != nil {
		m.Src = src
	}

	m.Ctr = m.Ctr.WithMountedDirectory(workdirRootPath, m.Src)

	return m
}

// WithTerragruntCache creates a cache volume for Terragrunt for the .terragrunt-cache folder
func (m *Terragrunt) WithTerragruntCache() *Terragrunt {
	tgCache := dag.CacheVolume(".terragrunt-cache")
	ctr := m.Ctr.WithMountedCache(".terragrunt-cache", tgCache)

	m.Ctr = ctr

	return m
}
