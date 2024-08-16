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
// - image: The image to use as the base container. Optional parameter.
// - ctr: The container to use as a base container. Optional parameter.
// - envVarsFromHost: A list of environment variables to pass from the host to the container in a
// slice of strings. Optional parameter.
//
// Returns a pointer to a Gotoolbox instance and an error, if any.
func New(
	// version is the version of the container image to use.
	// +optional
	version string,
	// image is the container image to use.
	// +optional
	image string,
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
		imageURL, err := containerx.GetImageURL(&containerx.NewBaseContainerOpts{
			Image:           image,
			Version:         version,
			FallbackImage:   defaultContainerImage,
			FallBackVersion: defaultContainerVersion,
		})

		if err != nil {
			return nil, WrapError(err, "failed to get image URL")
		}

		dagModule.Base(imageURL)
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
func (m *Gotoolbox) Base(imageURL string) *Gotoolbox {
	c := dag.Container().From(imageURL)
	m.Ctr = c

	return m
}

const (
	defaultAlpineImage  = "alpine"
	defaultUbuntuImage  = "ubuntu"
	defaultBusyBoxImage = "busybox"
)

// BaseAlpine sets the base image to an Alpine Linux image and creates the base container.
//
// Parameters:
// - version: The version of the Alpine image to use. Optional parameter. Defaults to "latest".
//
// Returns a pointer to the Gotoolbox instance.
func (m *Gotoolbox) BaseAlpine(
	// version is the version of the Alpine image to use, e.g., "3.17.3".
	// +optional
	version string,
) *Gotoolbox {
	if version == "" {
		version = "latest"
	}

	imageURL := fmt.Sprintf("%s:%s", defaultAlpineImage, version)

	return m.Base(imageURL)
}

// BaseUbuntu sets the base image to an Ubuntu Linux image and creates the base container.
//
// Parameters:
// - version: The version of the Ubuntu image to use. Optional parameter. Defaults to "latest".
//
// Returns a pointer to the Gotoolbox instance.
func (m *Gotoolbox) BaseUbuntu(
	// version is the version of the Ubuntu image to use, e.g., "22.04".
	// +optional
	version string,
) *Gotoolbox {
	if version == "" {
		version = "latest"
	}

	imageURL := fmt.Sprintf("%s:%s", defaultUbuntuImage, version)

	return m.Base(imageURL)
}

// BaseBusyBox sets the base image to a BusyBox Linux image and creates the base container.
//
// Parameters:
// - version: The version of the BusyBox image to use. Optional parameter. Defaults to "latest".
//
// Returns a pointer to the Gotoolbox instance.
func (m *Gotoolbox) BaseBusyBox(
	// version is the version of the BusyBox image to use, e.g., "1.35.0".
	// +optional
	version string,
) *Gotoolbox {
	if version == "" {
		version = "latest"
	}

	imageURL := fmt.Sprintf("%s:%s", defaultBusyBoxImage, version)

	return m.Base(imageURL)
}
