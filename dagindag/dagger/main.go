package main

import (
	"context"
	"fmt"
)

type Dagindag struct {
	// Ctr is the container to use as a base container.
	Ctr *Container
	// Src is the directory that contains all the source code, including the module directory.
	Src *Directory
}

// New creates a new instance of the Dagindag module with the given version.
//
// If the version is not specified, the default version is used.
// The default version is "latest".
func New(
	// daggerVersion is the version of the Dagger engine to use, e.g., "v0.11.5
	// +optional
	daggerVersion string,
	// dockerVersion is the version of the Docker engine to use, e.g., "24.0
	// +optional
	dockerVersion string,
	//ctr is the container to use as a base container.
	// +optional
	ctr *Container,
	// src is the directory that contains all the source code, including the module directory.
	// +optional
	src *Directory,
) (*Dagindag, error) {
	m := &Dagindag{}

	if ctr != nil {
		m.Ctr = ctr
	} else {
		c, err := m.Base(daggerVersion, dockerVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to set the base container: %w", err)
		}

		m.Ctr = c.Ctr
	}

	if src != nil {
		return m.WithSource(src), nil
	}

	return m, nil
}

// Base sets the base container for the Dagindag module.
//
// The base container is set to the Ubuntu container with the "lunar" tag.
// This container is used as the base container for all the other containers
func (m *Dagindag) Base(
	// daggerVersion is the version of the Dagger engine to use, e.g., "v0.11.5
	// +optional
	daggerVersion,
	// dockerVersion is the version of the Docker engine to use, e.g., "24.0
	// +optional
	dockerVersion string,
) (*Dagindag, error) {
	if daggerVersion == "" {
		// FIXME: There's a library and an API available in the daggerx package, for fetching the latest version of Dagger. Use it.
		daggerVersion = daggerDefaultVersion
	}

	// Configure the very base container. I thought initially to use ubuntu:lunar, but since docker-in-docker
	// is also lightweight (and it's based in alpine), I decided to use it as the base container.
	c := dag.Container().
		From(getDockerInDockerImage(dockerVersion))
	m.Ctr = c

	// Setup Dagger.
	m.Ctr = m.WithDaggerSetup(daggerVersion).
		WithDaggerCLIEntryPoint().Ctr

	// Creating the docker configuration for the dockerd and the service binding.
	dockerd := m.WithDockerService(dockerVersion)
	dockerHost, err := dockerd.Endpoint(context.Background(), ServiceEndpointOpts{
		Scheme: "tcp",
	})

	if err != nil {
		return nil, err
	}

	m.Ctr = m.Ctr.
		WithServiceBinding("docker", dockerd).
		WithEnvVariable("DOCKER_HOST", dockerHost)

	return m, nil
}
