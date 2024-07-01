// Package main provides methods for setting up and managing a container environment.
// This includes setting environment variables, mounting directories and files,
// and configuring services such as Docker within the container.
//
// Copyright: Excoriate alex_torres@outlook.com
// License: MIT
package main

import (
	"path/filepath"

	"github.com/Excoriate/daggerx/pkg/fixtures"
)

// WithEnvironmentVariable sets an environment variable in the container.
//
// Parameters:
//   - name: The name of the environment variable (e.g., "HOST").
//   - value: The value of the environment variable (e.g., "localhost").
//   - expand: Whether to replace `${VAR}` or $VAR in the value according to the current
//     environment variables defined in the container (e.g., "/opt/bin:$PATH").
//     Optional parameter.
func (m *Ctrutils) WithEnvironmentVariable(
	// name is the name of the environment variable.
	name string,
	// value is the value of the environment variable.
	value string,
	// expand is whether to replace `${VAR}` or $VAR in the value according to the current
	// +optional
	expand bool,
) *Ctrutils {
	m.Ctr = m.Ctr.WithEnvVariable(name, value, ContainerWithEnvVariableOpts{
		Expand: expand,
	})

	return m
}

// WithSource sets the source directory for the container.
//
// Parameters:
// - src: The directory that contains all the source code, including the module directory.
// - workdir: The working directory within the container. Optional parameter.
func (m *Ctrutils) WithSource(
	// src is the directory that contains all the source code, including the module directory.
	src *Directory,
	// workdir is the working directory within the container. If not set it'll default to /mnt
	// +optional
	workdir string,
) *Ctrutils {
	ctr := m.Ctr.WithMountedDirectory(fixtures.MntPrefix, src)

	if workdir != "" {
		ctr = ctr.WithWorkdir(filepath.Join(fixtures.MntPrefix, workdir))
	} else {
		ctr = ctr.WithWorkdir(fixtures.MntPrefix)
	}

	m.Ctr = ctr

	return m
}

// WithContainer sets the container to be used.
//
// Parameters:
// - ctr: The container to run the command in. If passed, it will override the container set in the Dagger instance.
func (m *Ctrutils) WithContainer(
	ctr *Container,
) *Ctrutils {
	m.Ctr = ctr

	return m
}

// WithDockerService sets up the container with the Docker service.
//
// It sets up the container with the Docker service.
// Parameters:
//   - dockerVersion: The version of the Docker engine to use, e.g., "v20.10.17".
//     Optional parameter. If not provided, a default version is used.
func (m *Ctrutils) WithDockerService(
	dockerVersion string,
) *Service {
	if dockerVersion == "" {
		dockerVersion = dockerVersionDefault
	}

	dindImage := getDockerInDockerImage(dockerVersion)
	dockerPort := 2375

	return dag.Container().
		From(dindImage).
		WithMountedCache(
			"/var/lib/docker",
			dag.CacheVolume(dockerVersion+"-docker-lib"),
			ContainerWithMountedCacheOpts{
				Sharing: Private,
			}).
		WithExposedPort(dockerPort).
		WithExec([]string{
			"dockerd",
			"--host=tcp://0.0.0.0:2375",
			"--host=unix:///var/run/docker.sock",
			"--tls=false",
		}, ContainerWithExecOpts{
			InsecureRootCapabilities: true,
		}).
		AsService()
}

// WithFileMountedInContainer adds a file to the container.
//
// Parameters:
// - file: The file to add to the container.
// - dest: The destination path in the container. Optional parameter.
// - owner: The owner of the file. Optional parameter.
func (m *Ctrutils) WithFileMountedInContainer(
	file *File,
	dest string,
	owner string,
) *Ctrutils {
	path := filepath.Join(fixtures.MntPrefix, dest)
	if owner != "" {
		m.Ctr = m.Ctr.WithMountedFile(path, file, ContainerWithMountedFileOpts{
			Owner: owner,
		})

		return m
	}

	m.Ctr = m.Ctr.WithMountedFile(path, file)

	return m
}
