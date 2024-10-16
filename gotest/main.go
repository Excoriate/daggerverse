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
	// version (image tag) to use for the container image.
	// +optional
	version string,
	// image is the imageURL (without the version) that's going to be used for his base container.
	// +optional
	image string,
	// ctr is the container to use as a base container.
	// +optional
	ctr *dagger.Container,
	// envVarsFromHost is a list of environment variables to pass from the host to the container in a slice of strings.
	// +optional
	envVarsFromHost []string,
) (*Gotest, error) {
	//nolint:exhaustruct // It's 'okaysh' for now, I'll decide later what's going to be the pattern here.
	dagModule := &Gotest{}

	if ctr != nil {
		dagModule.Ctr = ctr

		return dagModule, nil
	}

	if image != "" {
		isValid, err := containerx.ValidateImageURL(image)
		if err != nil {
			return nil, WrapErrorf(err, "failed to validate image URL passed with value %s", image)
		}

		if !isValid {
			return nil, Errorf("the image URL %s is not valid", image)
		}

		dagModule.Base(image)
	} else {

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
func (m *Gotest) Base(imageURL string) *Gotest {
	c := dag.
		Container().
		From(imageURL)

	m.Ctr = c

	return m
}

// getImageURL gets the image URL and validates it.
//
// If the image URL is not valid, it returns an error.
func (m *Gotest) getImageURL(image, version string) (string, error) {
	imageURL, err := containerx.GetImageURL(&containerx.NewBaseContainerOpts{
		Image:           image,
		Version:         version,
		FallbackImage:   defaultContainerImage,
		FallBackVersion: defaultContainerVersion,
	})

	if err != nil {
		return "", WrapErrorf(err, "failed to get image URL from image %s and version %s", image, version)
	}

	isValid, invalidImageErr := containerx.ValidateImageURL(imageURL)
	if invalidImageErr != nil {
		return "", WrapErrorf(invalidImageErr, "failed to validate image URL %s", imageURL)
	}

	if !isValid {
		return "", Errorf("the image URL %s is not valid", imageURL)
	}

	return imageURL, nil
}
