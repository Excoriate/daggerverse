// Package main provides the ModuleTemplateLight Dagger module and related functions.
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

// ModuleTemplateLight is a Dagger module.
//
// This module is used to create and manage containers.
type ModuleTemplateLight struct {
	// Ctr is the container to use as a base container.
	Ctr *dagger.Container
}

// New creates a new ModuleTemplateLight module.
//
// Parameters:
// - version: The version of the GoReleaser to use, e.g., "v1.22.0". Optional parameter.
// - image: The image to use as the base container. Optional parameter.
// - ctr: The container to use as a base container. Optional parameter.
// - envVarsFromHost: A list of environment variables to pass from the host to the container in a
// slice of strings. Optional parameter.
//
// Returns a pointer to a ModuleTemplateLight instance and an error, if any.
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
) (*ModuleTemplateLight, error) {
	//nolint:exhaustruct // It's 'okaysh' for now, I'll decide later what's going to be the pattern here.
	dagModule := &ModuleTemplateLight{}

	if ctr != nil {
		dagModule.Ctr = ctr
	} else {
		imageURL, imageURLErr := dagModule.getImageURL(image, version)
		if imageURLErr != nil {
			return nil, WrapError(imageURLErr, "failed to initialize Dagger module")
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
func (m *ModuleTemplateLight) Base(imageURL string) *ModuleTemplateLight {
	c := dag.
		Container().
		From(imageURL)

	m.Ctr = c

	return m
}

// getImageURL gets the image URL and validates it.
//
// If the image URL is not valid, it returns an error.
func (m *ModuleTemplateLight) getImageURL(image, version string) (string, error) {
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
