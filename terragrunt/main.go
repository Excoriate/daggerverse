// Package main provides the Terragrunt Dagger module and related functions.
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger. The module demonstrates
// usage of arguments and return types using simple echo and grep commands. The functions
// can be called from the dagger CLI or from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.
package main

import (
	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"

	"github.com/Excoriate/daggerx/pkg/envvars"
)

// Terragrunt is a Dagger module.
//
// This module is used to create and manage containers.
type Terragrunt struct {
	// Ctr is the container to use as a base container.
	Ctr *dagger.Container
	// BaseImage is the base image to use as the base container.
	// +private
	BaseImage BaseImageApko
	// Toolkit is a struct that contains the versions of the tools that are going to be used.
	// +private
	Toolkit *Toolkit
}

// New creates a new Terragrunt module.
//
// Parameters:
// - version: The version of the GoReleaser to use, e.g., "v1.22.0". Optional parameter.
// - image: The image to use as the base container. Optional parameter.
// - ctr: The container to use as a base container. Optional parameter.
// - envVarsFromHost: A list of environment variables to pass from the host to the container in a
// slice of strings. Optional parameter.
//
// Returns a pointer to a Terragrunt instance and an error, if any.
func New(
	// ctr is the container to use as a base container.
	// +optional
	ctr *dagger.Container,
	// tgVersion is the Terragrunt version to use. Default is "0.66.0".
	// +optional
	tgVersion string,
	// imageURL is the URL of the image to use as the base container.
	// It should includes tags. E.g. "ghcr.io/devops-infra/docker-terragrunt:tf-1.9.5-ot-1.8.2-tg-0.67.4"
	// +optional
	imageURL string,
	// tfVersion is the Terraform version to use. Default is "1.9.1".
	// +optional
	tfVersion string,
	// openTofuVersion is the OpenTofu version to use. Default is "1.8.0".
	// +optional
	openTofuVersion string,
	// envVarsFromHost is a list of environment variables to pass from the host to the container in a slice of strings.
	// +optional
	envVarsFromHost []string,
	// enableApko is a flag to enable Apko as a mechanism to build the container image. Default is false.
	// +optional
	enableApko bool,
) (*Terragrunt, error) {
	//nolint:exhaustruct // It's 'okaysh' for now, I'll decide later what's going to be the pattern here.
	dagModule := &Terragrunt{
		// Toolkit is a struct that contains the versions of the tools that are going to be used.
		Toolkit: NewToolkit(
			WithOpenTofuVersion(openTofuVersion),
			WithTerraformVersion(tfVersion),
			WithTerragruntVersion(tgVersion),
		),
	}

	if ctr != nil {
		dagModule.Ctr = ctr
	} else {
		if imageURL == "" {
			dagModule.Base(imageURL)
		}
	}

	// If environment variables are passed in a string, with a format like "SOMETHING=SOMETHING,SOMETHING=SOMETHING",
	// they are converted into a map and then into a list of DaggerEnvVars.
	// Then, each environment variable is added to the container.
	if len(envVarsFromHost) > 0 {
		envVars, err := envvars.ToDaggerEnvVarsFromSlice(envVarsFromHost)
		if err != nil {
			return nil, WrapError(err, "failed to parse environment variables")
		}

		for _, envVar := range envVars {
			dagModule.WithEnvironmentVariable(envVar.Name, envVar.Value, false)
		}
	}

	return dagModule, nil
}

// Base sets the base image and version, and creates the base container.
//
// The default image is "alpine/latest" and the default version is "latest".
//
//nolint:nolintlint,revive // This is a method that is used to set the base image and version.
func (m *Terragrunt) Base(imageURL string) *Terragrunt {
	c := dag.Container().
		From(imageURL)

	m.Ctr = c

	return m
}
