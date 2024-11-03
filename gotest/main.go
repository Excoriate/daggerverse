// Package main provides the Gotest Dagger module for container management.
//
// This Dagger module is tailored for running Go test commands within containerized environments.
// It enables users to execute Go tests efficiently by defining a base container,
// passing environment variables from the host,
// and managing container configurations seamlessly. The Gotest module exemplifies how to
// leverage Dagger's capabilities
// for executing Go tests in various workflows.
//
// Functions in this module can be invoked from the Dagger CLI or through SDKs,
// allowing for flexible integration into CI/CD pipelines. This module is designed to be
// extensible and adaptable for diverse testing scenarios in Go applications.
package main

import (
	"github.com/Excoriate/daggerverse/gotest/internal/dagger"

	"github.com/Excoriate/daggerx/pkg/containerx"
	"github.com/Excoriate/daggerx/pkg/envvars"
)

// Gotest is a Dagger module.
//
// This module is used to create and manage containers.
type Gotest struct {
	// Ctr is the container to use as a base container.
	Ctr *dagger.Container
}

// New creates a new Gotest module.
//
// Parameters:
// - version: The version of the GoReleaser to use, e.g., "v1.22.0". Optional parameter.
// - image: The image to use as the base container. Optional parameter.
// - ctr: The container to use as a base container. Optional parameter.
// - envVarsFromHost: A list of environment variables to pass from the host to the container in a
// slice of strings. Optional parameter.
//
// Returns a pointer to a Gotest instance and an error, if any.
func New(
	// version is the version of the GoReleaser to use, e.g., "v1.22.0". Optional parameter.
	// +optional
	version string,
	// image is the image to use as the base container. Optional parameter.
	// +optional
	image string,
	// ctr is the container to use as a base container. Optional parameter.
	// +optional
	ctr *dagger.Container,
	// envVarsFromHost is a list of environment variables to pass from the host to the container in a
	// slice of strings. Optional parameter.
	// +optional
	envVarsFromHost []string,
) (*Gotest, error) {
	dagModule := &Gotest{}

	if ctr != nil {
		dagModule.Ctr = ctr

		return dagModule, nil
	}

	if err := dagModule.setupContainer(image, version); err != nil {
		return nil, err
	}

	if err := dagModule.setupEnvironmentVariables(envVarsFromHost); err != nil {
		return nil, err
	}

	return dagModule, nil
}

// setupContainer sets up the container.
//
// If the image is not passed, it sets up the default image.
// If the image is passed, it sets up the custom image.
func (m *Gotest) setupContainer(image, version string) error {
	if image != "" {
		return m.setupCustomImage(image)
	}

	return m.setupDefaultImage(version)
}

// setupCustomImage sets up the custom image.
//
// It validates the image URL and sets the base container.
func (m *Gotest) setupCustomImage(image string) error {
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

// setupDefaultImage sets up the default image.
//
// If the version is not passed, it sets the default version.
// If the version is passed, it sets the custom version.
func (m *Gotest) setupDefaultImage(version string) error {
	if version == "" {
		version = defaultContainerVersion
	}

	imageURL, err := m.getImageURL(defaultContainerImage, version)

	if err != nil {
		return WrapErrorf(err, "failed to get image URL from image %s and version %s",
			defaultContainerImage, version)
	}

	m.Base(imageURL)

	return nil
}

// setupEnvironmentVariables sets up the environment variables.
//
// If the environment variables are not passed, it returns nil.
// If the environment variables are passed, it sets the environment variables.
func (m *Gotest) setupEnvironmentVariables(envVarsFromHost []string) error {
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
