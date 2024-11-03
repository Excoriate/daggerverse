// Package main provides the ModuleTemplate Dagger module and related functions.
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
	"github.com/Excoriate/daggerverse/module-template-light/internal/dagger"

	"github.com/Excoriate/daggerx/pkg/containerx"
	"github.com/Excoriate/daggerx/pkg/envvars"
)

// ModuleTemplate is a Dagger module.
//
// This module is used to create and manage containers.
type ModuleTemplate struct {
	// Ctr is the container to use as a base container.
	Ctr *dagger.Container
}

// New creates a new ModuleTemplate module.
//
// Parameters:
// - version: The version of the GoReleaser to use, e.g., "v1.22.0". Optional parameter.
// - image: The image to use as the base container. Optional parameter.
// - ctr: The container to use as a base container. Optional parameter.
// - envVarsFromHost: A list of environment variables to pass from the host to the container in a
// slice of strings. Optional parameter.
//
// Returns a pointer to a ModuleTemplate instance and an error, if any.
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
) (*ModuleTemplate, error) {
	//nolint:exhaustruct // It's 'okaysh' for now, I'll decide later what's going to be the pattern here.
	dagModule := &ModuleTemplate{}

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

// Base sets the base image and version, and creates the base container.
//
// The default image is "alpine/latest" and the default version is "latest".
//
//nolint:nolintlint,revive // This is a method that is used to set the base image and version.
func (m *ModuleTemplate) Base(imageURL string) *ModuleTemplate {
	c := dag.
		Container().
		From(imageURL)

	m.Ctr = c

	return m
}

// getImageURL gets the image URL and validates it.
//
// If the image URL is not valid, it returns an error.
func (m *ModuleTemplate) getImageURL(image, version string) (string, error) {
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

// setupContainer sets up the container.
//
// If the image is not passed, it sets up the default image.
// If the image is passed, it sets up the custom image.
func (m *ModuleTemplate) setupContainer(image, version string) error {
	if image != "" {
		return m.setupCustomImage(image)
	}

	return m.setupDefaultImage(version)
}

// setupCustomImage sets up the custom image.
//
// It validates the image URL and sets the base container.
func (m *ModuleTemplate) setupCustomImage(image string) error {
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
func (m *ModuleTemplate) setupDefaultImage(version string) error {
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
func (m *ModuleTemplate) setupEnvironmentVariables(envVarsFromHost []string) error {
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
