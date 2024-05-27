package main

import (
	"github.com/Excoriate/daggerx/pkg/containerx"
)

type Goreleaser struct {
	// Src is the directory that contains all the source code, including the module directory.
	Src *Directory

	// Ctr is the container to use as a base container.
	Ctr *Container

	// CfgFile is the configuration file to use.
	CfgFile string
}

func New(
	// version is the version of the GoReleaser to use, e.g., "v1.22.0".
	// +default="latest"
	// +optional
	version string,
	// image is the image to use as the base container.
	// +optional
	// +default="goreleaser/goreleaser"
	image string,
	// Ctrl is the container to use as a base container.
	// +optional
	ctr *Container,
	// envVarsFromHost is a list of environment variables to pass from the host to the container.
	// Later on, in order to pass it to the container, it's going to be converted into a map.
	// +optional
	envVarsFromHost string,
) (*Goreleaser, error) {
	g := &Goreleaser{}

	if ctr != nil {
		g.Ctr = ctr
	} else {
		imageURL, err := containerx.GetImageURL(&containerx.NewBaseContainerOpts{
			Image:           image,
			Version:         version,
			FallbackImage:   goReleaserDefaultImage,
			FallBackVersion: goReleaserDefaultVersion,
		})

		if err != nil {
			return nil, err
		}

		g.Ctr = g.Base(imageURL).Ctr
	}

	g.addEnvVarsToContainerFromStr(envVarsFromHost)

	// Here, regardless of whether the container is set or not,
	// the WithGoCache method is called.
	g = g.WithGoCache()

	return g, nil
}

// Base sets the base image and version, and creates the base container.
// The default image is "goreleaser/goreleaser" and the default version is "latest".
func (m *Goreleaser) Base(imageURL string) *Goreleaser {
	c := dag.Container().From(imageURL).
		WithEnvVariable("TINI_SUBREAPER", "true")

	m.Ctr = c

	return m
}
