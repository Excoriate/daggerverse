// Package main provides the Gopkgpublisher Dagger module and related functions.
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
	"github.com/Excoriate/daggerverse/gopkgpublisher/internal/dagger"

	"github.com/Excoriate/daggerx/pkg/containerx"
	"github.com/Excoriate/daggerx/pkg/envvars"
)

// Gopkgpublisher is a Dagger module.
//
// This module is used to create and manage containers.
type Gopkgpublisher struct {
	// Ctr is the container to use as a base container.
	Ctr *dagger.Container
	// ApkoPackages is the list of packages to install with APKO.
	ApkoPackages []string
	// Publish is the Publish module that handles the publishing process.
	// +private
	Publish *Publish
}

// New creates a new Gopkgpublisher module.
//
// Parameters:
// - version: The version of the GoReleaser to use, e.g., "v1.22.0". Optional parameter.
// - image: The image to use as the base container. Optional parameter.
// - ctr: The container to use as a base container. Optional parameter.
// - envVarsFromHost: A list of environment variables to pass from the host to the container in a
// slice of strings. Optional parameter.
//
// Returns a pointer to a Gopkgpublisher instance and an error, if any.
func New(
	// version is the Go version to use for the container image. Default to "go-1.22.8-r0"
	// +optional
	version string,
	// imageURL is the imageURL (including the tag) that's going to be used for his base container.
	// +optional
	imageURL string,
	// ctr is the container to use as a base container.
	// +optional
	ctr *dagger.Container,
	// envVarsFromHost is a list of environment variables to pass from the host to the container in a slice of strings.
	// +optional
	envVarsFromHost []string,
	// extraPackages is a list of extra packages to install with APKO, from the Alpine packages repository.
	// +optional
	extraPackages []string,
) (*Gopkgpublisher, error) {
	//nolint:exhaustruct // It's 'okaysh' for now, I'll decide later what's going to be the pattern here.
	dagModule := &Gopkgpublisher{
		ApkoPackages: extraPackages,
		Publish:      &Publish{},
	}

	if ctr != nil {
		dagModule.Ctr = ctr

		return dagModule, nil
	}

	if imageURL != "" {
		if err := dagModule.setupCustomImage(imageURL); err != nil {
			return nil, WrapError(err, "failed to setup custom image")
		}
	} else {
		dagModule.WithGoPackage(version)

		// Build base image using APKO
		if _, err := dagModule.BaseApko(); err != nil {
			return nil, WrapError(err, "failed to create base image with apko")
		}

		// Configure cache and permissions only after container is initialized
		if dagModule.Ctr != nil {
			dagModule.
				WithCacheConfiguration().
				WithDirPermissionsConfiguration()
		}
	}

	if err := dagModule.setupEnvironmentVariables(envVarsFromHost); err != nil {
		return nil, err
	}

	return dagModule, nil
}

// Base sets the base image and version, and creates the base container.
//
// The default image is "alpine/latest" and the default version is "latest".
//
//nolint:nolintlint,revive // This is a method that is used to set the base image and version.
func (m *Gopkgpublisher) Base(imageURL string) *Gopkgpublisher {
	c := dag.
		Container().
		From(imageURL)

	m.Ctr = c

	return m
}

// setupCustomImage sets up the custom image.
//
// It validates the image URL and sets the base container.
func (m *Gopkgpublisher) setupCustomImage(image string) error {
	isValid, err := containerx.ValidateImageURL(image)

	if err != nil {
		return WrapErrorf(err, "failed to validate image URL passed with value %s", image)
	}

	if !isValid {
		return Errorf("the image URL %s is not valid", image)
	}

	m.Base(image)

	return nil
}

// setupEnvironmentVariables sets up the environment variables.
//
// If the environment variables are not passed, it returns nil.
// If the environment variables are passed, it sets the environment variables.
func (m *Gopkgpublisher) setupEnvironmentVariables(envVarsFromHost []string) error {
	if len(envVarsFromHost) == 0 {
		return nil
	}

	envVars, err := envvars.ToDaggerEnvVarsFromSlice(envVarsFromHost)
	if err != nil {
		return WrapError(err, "failed to parse environment variables")
	}

	for _, envVar := range envVars {
		m.WithEnvironmentVariable(envVar.Name, envVar.Value, false)
	}

	return nil
}
