// Package main provides methods for setting up and managing a container environment.
// This includes setting environment variables, mounting directories and files,
// and configuring services such as Docker within the container.
//
// Copyright: Excoriate alex_torres@outlook.com
// License: MIT
package main

import (
	"path/filepath"
	"time"

	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"

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
func (m *Terragrunt) WithEnvironmentVariable(
	// name is the name of the environment variable.
	name string,
	// value is the value of the environment variable.
	value string,
	// expand is whether to replace `${VAR}` or $VAR in the value according to the current
	// +optional
	expand bool,
) *Terragrunt {
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
func (m *Terragrunt) WithSource(
	// src is the directory that contains all the source code, including the module directory.
	src *dagger.Directory,
	// workdir is the working directory within the container. If not set it'll default to /mnt
	// +optional
	workdir string,
) *Terragrunt {
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
func (m *Terragrunt) WithContainer(
	ctr *dagger.Container,
) *Terragrunt {
	m.Ctr = ctr

	return m
}

// WithFileMountedInContainer adds a file to the container.
//
// Parameters:
// - file: The file to add to the container.
// - dest: The destination path in the container. Optional parameter.
// - owner: The owner of the file. Optional parameter.
func (m *Terragrunt) WithFileMountedInContainer(
	file *dagger.File,
	dest string,
	owner string,
) *Terragrunt {
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

// WithDownloadedFile downloads a file from the specified URL and mounts it in the container.
//
// Parameters:
//   - url: The URL of the file to download.
//   - destDir: The directory within the container where the file will be downloaded. Optional parameter.
//     If not provided, it defaults to the predefined mount prefix.
//
// Returns:
//   - *Terragrunt: The updated Terragrunt with the downloaded file mounted in the container.
func (m *Terragrunt) WithDownloadedFile(
	// url is the URL of the file to download.
	url string,
	// destFileName is the name of the file to download. If not set, it'll default to the basename of the URL.
	// +optional
	destFileName string,
) *Terragrunt {
	// Extract the filename from the last part of the URL.
	fileName := filepath.Base(url)
	if destFileName != "" {
		fileName = destFileName
	}

	// Download the file
	fileDownloaded := dag.HTTP(url).WithName(fileName)

	// Define the path in the container
	destFilePath := filepath.Join(fixtures.MntPrefix, fileName)

	// Mount the file in the container
	m.Ctr = m.
		Ctr.
		WithMountedFile(destFilePath, fileDownloaded)

	return m
}

// WithCacheBuster sets a cache-busting environment variable in the container.
//
// This method sets an environment variable "CACHE_BUSTER" with a timestamp value in RFC3339Nano format.
// This can be useful for invalidating caches by providing a unique value.
//
// Returns:
//   - *Terragrunt: The updated Terragrunt with the cache-busting environment variable set.
func (m *Terragrunt) WithCacheBuster() *Terragrunt {
	m.Ctr = m.Ctr.
		WithEnvVariable("CACHE_BUSTER", time.
			Now().
			Format(time.RFC3339Nano))

	return m
}
