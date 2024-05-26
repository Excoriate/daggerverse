package main

import (
	"github.com/Excoriate/daggerx/pkg/containerx"

	"github.com/Excoriate/daggerx/pkg/envvars"
)

type Gotest struct {
	// Src is the directory that contains all the source code, including the module directory.
	Src *Directory

	// Ctr is the container to use as a base container.
	Ctr *Container
}

// New creates a new instance of the GoTest module with the given version, image, source directory, and environment variables.
// If the version is not specified, the default version is used.
// If the image is not specified, the default image is used.
// If the container is not specified, a new container is created, and this container is considered
// as the base container.
func New(
	// version is the version of the GoReleaser to use, e.g., "v1.22.0".
	// +optional
	version string,
	// image is the image to use as the base container.
	// +optional
	image string,
	// Ctrl is the container to use as a base container.
	// +optional
	ctr *Container,
	// EnvVarsFromHost is a list of environment variables to pass from the host to the container.
	// Later on, in order to pass it to the container, it's going to be converted into a map.
	// +optional
	envVarsFromHost string,
) (*Gotest, error) {
	g := &Gotest{}

	if ctr != nil {
		g.Ctr = ctr
	} else {
		imageURL, err := containerx.GetImageURL(&containerx.NewBaseContainerOpts{
			Image:           image,
			Version:         version,
			FallbackImage:   goTestDefaultImage,
			FallBackVersion: goTestDefaultVersion,
		})

		if err != nil {
			return nil, err
		}

		g.Base(imageURL)
	}

	// If environment variables are passed in a string, with a format like "SOMETHING=SOMETHING,SOMETHING=SOMETHING",
	// they are converted into a map and then into a list of DaggerEnvVars.
	// Then, each environment variable is added to the container.
	if envVarsFromHost != "" {
		envVars, err := envvars.ToDaggerEnvVarsFromStr(envVarsFromHost)
		if err != nil {
			return nil, err
		}
		for _, envVar := range envVars {
			g.Ctr = g.WithEnvVariable(envVar.Name, envVar.Value, false).Ctr
		}
	}

	return g, nil
}

// Base sets the base image and version, and creates the base container.
// The default image is "golang/alpine" and the default version is "latest".
func (m *Gotest) Base(imageURL string) *Gotest {
	c := dag.Container().From(imageURL).
		WithEnvVariable("CGO_ENABLED", "0")

	m.Ctr = c

	return m
}
