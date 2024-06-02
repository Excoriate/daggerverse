package main

import (
	"fmt"

	"github.com/Excoriate/daggerx/pkg/containerx"
)

type Precommit struct {
	// Src is the directory that contains all the source code, including the module directory.
	Src *Directory

	// Ctr is the container to use as a base container for pre-commit, if it's passed, it's used as the base container.
	Ctr *Container
}

func New(
	// version is the version of the Precommit to use, e.g., "v1.22.0".
	// +optional
	version string,
	// Src is the directory that contains all the source code, including the module directory.
	src *Directory,
	// Ctr is the container to use as a base container for pre-commit, if it's passed, it's used as the base container.
	// +optional
	ctr *Container,
	// EnvVarsFromHost is a list of environment variables to pass from the host to the container.
	// +optional
	envVarsFromHost []string,
) (*Precommit, error) {
	m := &Precommit{
		Src: src,
	}

	imageURL, err := containerx.GetImageURL(&containerx.NewBaseContainerOpts{
		Image:           DefaultImage,
		Version:         version,
		FallBackVersion: DefaultVersion,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get image URL: %w", err)
	}

	m.Ctr = m.Base(imageURL).Ctr

	return nil, nil
}

// Base sets the base container for pre-commit.
func (m *Precommit) Base(imageURL string) *Precommit {
	c := dag.Container().From(imageURL)

	m.Ctr = c
	return m
}
