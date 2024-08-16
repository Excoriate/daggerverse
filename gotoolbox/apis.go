// Package main provides methods for setting up and managing a container environment.
// This includes setting environment variables, mounting directories and files,
// and configuring services such as Docker within the container.
//
// Copyright: Excoriate alex_torres@outlook.com
// License: MIT
package main

import (
	"context"
	"path/filepath"
	"time"

	"github.com/Excoriate/daggerverse/gotoolbox/internal/dagger"

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
func (m *Gotoolbox) WithEnvironmentVariable(
	// name is the name of the environment variable.
	name string,
	// value is the value of the environment variable.
	value string,
	// expand is whether to replace `${VAR}` or $VAR in the value according to the current
	// +optional
	expand bool,
) *Gotoolbox {
	m.Ctr = m.Ctr.WithEnvVariable(name, value, dagger.ContainerWithEnvVariableOpts{
		Expand: expand,
	})

	return m
}

// WithSource sets the source directory for the container.
//
// Parameters:
// - src: The directory that contains all the source code, including the module directory.
// - workdir: The working directory within the container. Optional parameter.
func (m *Gotoolbox) WithSource(
	// src is the directory that contains all the source code, including the module directory.
	src *dagger.Directory,
	// workdir is the working directory within the container. If not set it'll default to /mnt
	// +optional
	workdir string,
) *Gotoolbox {
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
func (m *Gotoolbox) WithContainer(
	ctr *dagger.Container,
) *Gotoolbox {
	m.Ctr = ctr

	return m
}

// WithFileMountedInContainer adds a file to the container.
//
// Parameters:
// - file: The file to add to the container.
// - dest: The destination path in the container. Optional parameter.
// - owner: The owner of the file. Optional parameter.
func (m *Gotoolbox) WithFileMountedInContainer(
	file *dagger.File,
	dest string,
	owner string,
) *Gotoolbox {
	path := filepath.Join(fixtures.MntPrefix, dest)
	if owner != "" {
		m.Ctr = m.Ctr.WithMountedFile(path, file, dagger.ContainerWithMountedFileOpts{
			Owner: owner,
		})

		return m
	}

	m.Ctr = m.Ctr.WithMountedFile(path, file)

	return m
}

// WithGitInAlpineContainer installs Git in the golang/alpine container.
//
// It installs Git in the golang/alpine container.
func (m *Gotoolbox) WithGitInAlpineContainer() *Gotoolbox {
	m.Ctr = m.Ctr.
		WithExec([]string{"apk", "add", "git"})

	return m
}

// WithUtilitiesInAlpineContainer installs common utilities in the golang/alpine container.
//
// It installs utilities such as curl, wget, and others that are commonly used.
func (m *Gotoolbox) WithUtilitiesInAlpineContainer() *Gotoolbox {
	m.Ctr = m.Ctr.
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "add", "curl", "wget", "bash", "jq", "yq"})

	return m
}

// WithSecretAsEnvVar sets an environment variable in the container using a secret.
//
// Parameters:
//   - name: The name of the environment variable (e.g., "API_KEY").
//   - secret: The secret containing the value of the environment variable.
//
// Returns:
//   - *Gotoolbox: The updated Gotoolbox with the environment variable set.
//
// Behavior:
//   - The secret value is expanded according to the current environment variables defined in the container.
func (m *Gotoolbox) WithSecretAsEnvVar(name string, secret *dagger.Secret) *Gotoolbox {
	secretValue, err := secret.Plaintext(context.Background())
	if err != nil {
		return nil
	}

	m.Ctr = m.Ctr.WithEnvVariable(name, secretValue, dagger.ContainerWithEnvVariableOpts{
		Expand: true,
	})

	return m
}

// WithCacheBuster sets a cache-busting environment variable in the container.
//
// This method sets an environment variable "CACHE_BUSTER" with a timestamp value in RFC3339Nano format.
// This can be useful for invalidating caches by providing a unique value.
//
// Returns:
//   - *Gotoolbox: The updated Gotoolbox with the cache-busting environment variable set.
func (m *Gotoolbox) WithCacheBuster() *Gotoolbox {
	m.Ctr = m.Ctr.
		WithEnvVariable("CACHE_BUSTER", time.
			Now().
			Format(time.RFC3339Nano))

	return m
}
