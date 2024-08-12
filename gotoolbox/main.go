// Package main provides the Gotoolbox Dagger module and related functions.
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
	"fmt"

	"github.com/Excoriate/daggerverse/gotoolbox/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/containerx"
	"github.com/Excoriate/daggerx/pkg/envvars"
)

// Gotoolbox is a Dagger module.
//
// This module is used to create and manage containers.
type Gotoolbox struct {
	// Ctr is the container to use as a base container.
	Ctr *dagger.Container
}

// New creates a new Gotoolbox module.
//
// Parameters:
// - version: The version of the GoReleaser to use, e.g., "v1.22.0". Optional parameter.
// - ctr: The container to use as a base container. Optional parameter.
// - envVarsFromHost: A list of environment variables to pass from the host to the container in a
// slice of strings. Optional parameter.
//
// Returns a pointer to a Gotoolbox instance and an error, if any.
func New(
	// version is the version of the container image to use.
	// +optional
	version string,
	// ctr is the container to use as a base container.
	// +optional
	ctr *dagger.Container,
	// envVarsFromHost is a list of environment variables to pass from the host to the container in a slice of strings.
	// +optional
	envVarsFromHost []string,
) (*Gotoolbox, error) {
	//nolint:exhaustruct // It's 'okaysh' for now, I'll decide later what's going to be the pattern here.
	dagModule := &Gotoolbox{}

	if ctr != nil {
		dagModule.Ctr = ctr
	} else {
		if version != "" {
			version = fmt.Sprintf("%s-alpine3.19", version)
		}

		imageURL, err := containerx.GetImageURL(&containerx.NewBaseContainerOpts{
			Version:         version,
			FallbackImage:   defaultContainerImage,
			FallBackVersion: defaultContainerVersion,
		})

		if err != nil {
			return nil, fmt.Errorf("failed to get image URL: %w", err)
		}

		dagModule.Base(imageURL)
	}

	// If environment variables are passed in a string, with a format like "SOMETHING=SOMETHING,SOMETHING=SOMETHING",
	// they are converted into a map and then into a list of DaggerEnvVars.
	// Then, each environment variable is added to the container.
	if len(envVarsFromHost) > 0 {
		envVars, err := envvars.ToDaggerEnvVarsFromSlice(envVarsFromHost)
		if err != nil {
			return nil, fmt.Errorf("failed to parse environment variables: %w", err)
		}

		for _, envVar := range envVars {
			dagModule.
				WithEnvironmentVariable(envVar.Name, envVar.Value, false)
		}
	}

	return dagModule, nil
}

// Base sets the base image and version, and creates the base container.
//
// The default image is "alpine/latest" and the default version is "latest".
//
//nolint:nolintlint,revive // This is a method that is used to set the base image and version.
func (m *Gotoolbox) Base(imageURL string) *Gotoolbox {
	c := dag.Container().From(imageURL)
	m.Ctr = c

	return m.
		WithGitInAlpineContainer().
		WithUtilitiesInAlpineContainer()
}
