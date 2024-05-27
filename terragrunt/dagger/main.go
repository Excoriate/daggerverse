package main

import (
	"fmt"

	"github.com/Excoriate/daggerx/pkg/merger"

	"github.com/Excoriate/daggerx/pkg/cmdbuilder"

	"github.com/Excoriate/daggerx/pkg/containerx"
)

type Terragrunt struct {
	// Src is the directory that contains all the source code, including the module directory.
	Src *Directory

	// Ctr is the container to use as a base container.
	Ctr *Container

	// tfVersion is the version of the Terraform to use, e.g., "1.0.0". For more information, visit HashiCorp's Terraform GitHub repository.
	// +optional
	tfVersion string
}

func New(
	// tfVersion is the version of the Terraform to use, e.g., "1.0.0". For more information, visit HashiCorp's Terraform GitHub repository.
	// +optional
	tfVersion string,
	// image is the image to use as the base container.
	// +optional
	image string,
	// src is the directory that contains all the source code, including the module directory.
	src *Directory,
	// Ctrl is the container to use as a base container.
	// +optional
	ctr *Container,
	// envVarsFromHost is the list of environment variables to pass from the host to the container.
	// the format of the environment variables passed from the host are slices of strings separted by "=", and commas.
	// E.g., []string{"HOST=localhost", "PORT=8080"}
	// +optional
	envVarsFromHost []string,
) (*Terragrunt, error) {
	m := &Terragrunt{
		Src: src,
	}

	if ctr != nil {
		m.Ctr = ctr
	} else {
		imageURL, err := containerx.GetImageURL(&containerx.NewBaseContainerOpts{
			Image:           image,
			FallBackVersion: tfVersionDefault,
			FallbackImage:   tgContainerImageDefault,
			Version:         tfVersion,
		})

		if err != nil {
			return nil, fmt.Errorf("failed to create the image URL: %w", err)
		}

		m.Base(imageURL)
	}

	m.addEnvVarsToContainerFromSlice(envVarsFromHost)

	return m, nil
}

// Base sets the base image and version, and creates the base container.
// For Terragrunt, the default image is "alpine/terragrunt" and the default version is "1.7.0".
// Consider that the container isn't created based on the Terragrunt version, but on the Terraform version.
func (m *Terragrunt) Base(imageURL string) *Terragrunt {
	c := dag.Container().From(imageURL)
	m.Ctr = c

	return m
}

// addCMD adds a command to the container.
// It supports environment variables and arguments.
func (m *Terragrunt) addCMD(
	// module is the terragunt module to use that includes the terragrunt.hcl
	// the module should be relative to the mounted directory (--src).
	module string,
	// cmd is the command to run.
	cmd string,
	// envVars is the list of environment variables to pass from the host to the container.
	// the format of the environment variables passed from the host are slices of strings separted by "=", and commas.
	// E.g., []string{"HOST=localhost", "PORT=8080"}
	// +optional
	// +default=[]
	envVars []string,
	// args is the list of arguments to pass to the Terragrunt command.
	// +optional
	args string,
) (*Terragrunt, error) {
	// Add the terragrunt module as the workdir.
	m.Ctr = m.WithSource(m.Src, module).Ctr

	// Add the environment variables to the container.
	m.addEnvVarsToContainerFromSlice(envVars)

	// Add the command to the container.
	cmdWithArgs := merger.MergeSlices([]string{tgCMD, cmd}, cmdbuilder.BuildArgs(args))
	m.Ctr = m.WithCMD(cmdWithArgs).Ctr

	return m, nil
}
