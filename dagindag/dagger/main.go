package main

import (
	"fmt"

	"github.com/Excoriate/daggerx/pkg/cmdbuilder"
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
	// version is the version of the Dagger engine to use, e.g., "v0.11.5
	// +optional
	version string,
	// src is the directory that contains all the source code, including the module directory.
	// +optional
	src *Directory,
) (*Dagindag, error) {
	m := &Dagindag{}

	return m, nil
}

// Base sets the base container for the Dagindag module.
//
// The base container is set to the Ubuntu container with the "lunar" tag.
// This container is used as the base container for all the other containers
func (m *Dagindag) Base() *Dagindag {
	c := dag.Container().From("ubuntu:lunar")
	m.Ctr = c
	return m
}

func (m *Dagindag) WithDaggerCLI(version string) (*Dagindag, error) {
	cmdSetup, err := cmdbuilder.GenerateShCommand(getDaggerInstallCMDByVersion(version))
	if err != nil {
		return nil, fmt.Errorf("failed to generate Dagger setup command: %w, command: %s", err, cmdSetup)
	}

	if version == "" {
		version = daggerDefaultVersion
	}

	cmdDaggerInstall, cmdDaggerErr := cmdbuilder.GenerateShCommand(getDaggerInstallCMDByVersion(version))
	if cmdDaggerErr != nil {
		return nil, fmt.Errorf("failed to generate Dagger install command: %w, command: %s", cmdDaggerErr, cmdDaggerInstall)
	}

	m.Ctr = m.Ctr.WithExec(cmdSetup)

}
