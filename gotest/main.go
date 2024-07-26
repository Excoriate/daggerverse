package main

import (
	"fmt"

	"github.com/Excoriate/daggerx/pkg/containerx"

	"github.com/Excoriate/daggerx/pkg/envvars"
	"github.com/excoriate/daggerverse/gotest/internal/dagger"
)

// Gotest is a module that provides functionality to run Go tests.
//
// The module can be used to run tests in a Go project.
type Gotest struct {
	// Src is the directory that contains all the source code, including the module directory.
	Src *dagger.Directory

	// Ctr is the container to use as a base container.
	Ctr *dagger.Container
}

// New creates a new Gotest module.
func New(
	// version is the version of the GoReleaser to use, e.g., "v1.22.0".
	// +optional
	version string,
	// image is the image to use as the base container.
	// +optional
	image string,
	// Ctrl is the container to use as a base container.
	// +optional
	ctr *dagger.Container,
	// EnvVarsFromHost is a list of environment variables to pass from the host to the container.
	// Later on, in order to pass it to the container, it's going to be converted into a map.
	// +optional
	envVarsFromHost string,
) (*Gotest, error) {
	//nolint:exhaustruct // It's 'okaysh' for now, I'll decide later what's going to be the pattern here.
	dagModule := &Gotest{}

	if ctr != nil {
		dagModule.Ctr = ctr
	} else {
		imageURL, err := containerx.GetImageURL(&containerx.NewBaseContainerOpts{
			Image:           image,
			Version:         version,
			FallbackImage:   goTestDefaultImage,
			FallBackVersion: goTestDefaultVersion,
		})

		if err != nil {
			return nil, fmt.Errorf("failed to get image URL: %w", err)
		}

		dagModule.Base(imageURL)
	}

	// If environment variables are passed in a string, with a format like "SOMETHING=SOMETHING,SOMETHING=SOMETHING",
	// they are converted into a map and then into a list of DaggerEnvVars.
	// Then, each environment variable is added to the container.
	if envVarsFromHost != "" {
		envVars, err := envvars.ToDaggerEnvVarsFromStr(envVarsFromHost)
		if err != nil {
			return nil, fmt.Errorf("failed to parse environment variables: %w", err)
		}

		for _, envVar := range envVars {
			dagModule.Ctr = dagModule.WithEnvironmentVariable(envVar.Name, envVar.Value, false).Ctr
		}
	}

	return dagModule, nil
}

// Base sets the base image and version, and creates the base container.
// The default image is "golang/alpine" and the default version is "latest".
//
//nolint:nolintlint,revive // This is a method that is used to set the base image and version.
func (m *Gotest) Base(imageURL string) *Gotest {
	c := dag.Container().From(imageURL)
	m.Ctr = c

	return m.WithCgoDisabled()
}
