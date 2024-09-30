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
	"strings"
	"time"

	"github.com/Excoriate/daggerverse/module-template-light/internal/dagger"

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
func (m *ModuleTemplateLight) WithEnvironmentVariable(
	// name is the name of the environment variable.
	name string,
	// value is the value of the environment variable.
	value string,
	// expand is whether to replace `${VAR}` or $VAR in the value according to the current
	// +optional
	expand bool,
) *ModuleTemplateLight {
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
func (m *ModuleTemplateLight) WithSource(
	// src is the directory that contains all the source code, including the module directory.
	src *dagger.Directory,
	// workdir is the working directory within the container. If not set it'll default to /mnt
	// +optional
	workdir string,
) *ModuleTemplateLight {
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
func (m *ModuleTemplateLight) WithContainer(
	ctr *dagger.Container,
) *ModuleTemplateLight {
	m.Ctr = ctr

	return m
}

// WithFileMountedInContainer adds a file to the container.
//
// Parameters:
// - file: The file to add to the container.
// - dest: The destination path in the container. Optional parameter.
// - owner: The owner of the file. Optional parameter.
func (m *ModuleTemplateLight) WithFileMountedInContainer(
	file *dagger.File,
	dest string,
	owner string,
) *ModuleTemplateLight {
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

// WithSecretAsEnvVar sets an environment variable in the container using a secret.
//
// Parameters:
//   - name: The name of the environment variable (e.g., "API_KEY").
//   - secret: The secret containing the value of the environment variable.
//
// Returns:
//   - *ModuleTemplateLight: The updated ModuleTemplateLight with the environment variable set.
//
// Behavior:
//   - The secret value is expanded according to the current environment variables defined in the container.
func (m *ModuleTemplateLight) WithSecretAsEnvVar(name string, secret *dagger.Secret) *ModuleTemplateLight {
	secretValue, err := secret.Plaintext(context.Background())
	if err != nil {
		return nil
	}

	m.Ctr = m.Ctr.WithEnvVariable(name, secretValue, dagger.ContainerWithEnvVariableOpts{
		Expand: true,
	})

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
//   - *ModuleTemplateLight: The updated ModuleTemplateLight with the downloaded file mounted in the container.
func (m *ModuleTemplateLight) WithDownloadedFile(
	// url is the URL of the file to download.
	url string,
	// destFileName is the name of the file to download. If not set, it'll default to the basename of the URL.
	// +optional
	destFileName string,
) *ModuleTemplateLight {
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

// WithClonedGitRepo clones a Git repository and mounts it as a directory in the container.
//
// This method downloads a Git repository and mounts it as a directory in the container. It supports optional
// authentication tokens for private repositories and can handle both GitHub and GitLab repositories.
//
// Parameters:
//   - repoURL: The URL of the git repository to clone (e.g., "https://github.com/user/repo").
//   - token: (Optional) The VCS token to use for authentication. If
//     not provided, the repository will be cloned without authentication.
//   - vcs: (Optional) The version control system (VCS) to use for
//     authentication. Defaults to "github". Supported values are "github" and "gitlab".
//
// Returns:
//   - *ModuleTemplateLight: The updated ModuleTemplateLight with the cloned repository mounted in the container.
func (m *ModuleTemplateLight) WithClonedGitRepo(
	repoURL string,
	// token is the VCS token to use for authentication. Optional parameter.
	// +optional
	token string,
	// vcs is the VCS to use for authentication. Optional parameter.
	// +optional
	vcs string,
) *ModuleTemplateLight {
	// Call the helper function to clone the repository.
	clonedRepo := m.CloneGitRepo(repoURL, token, vcs)

	// Mount the cloned repository as a directory inside the container.
	m.Ctr = m.Ctr.WithMountedDirectory(fixtures.MntPrefix, clonedRepo)

	return m
}

// WithCacheBuster sets a cache-busting environment variable in the container.
//
// This method sets an environment variable "CACHE_BUSTER" with a timestamp value in RFC3339Nano format.
// This can be useful for invalidating caches by providing a unique value.
//
// Returns:
//   - *ModuleTemplateLight: The updated ModuleTemplateLight with the cache-busting environment variable set.
func (m *ModuleTemplateLight) WithCacheBuster() *ModuleTemplateLight {
	m.Ctr = m.Ctr.
		WithEnvVariable("CACHE_BUSTER", time.
			Now().
			Format(time.RFC3339Nano))

	return m
}

// WithConfigFile mounts a configuration file into the container at the specified path.
//
// This method allows you to mount a configuration file into the container at a specified path.
// If no path is provided, it defaults to a predefined mount prefix from the fixtures package.
//
// Args:
//   - cfgPath (string): The path where the config file will be mounted. If empty, it defaults to fixtures.MntPrefix.
//     +optional
//   - cfgFile (*dagger.File): The config file to be mounted.
//
// Returns:
//   - *ModuleTemplateLight: The updated ModuleTemplateLight with the config file mounted in the container.
func (m *ModuleTemplateLight) WithConfigFile(
	// cfgPath is the path where the config file will be mounted.
	// +optional
	cfgPath string,
	// setEnvVar is a string that set an environment variable in the container with the config file path.
	// +optional
	setEnvVar string,
	// cfgFile is the config file to be mounted.
	cfgFile *dagger.File) *ModuleTemplateLight {
	if cfgPath == "" {
		cfgPath = fixtures.MntPrefix
	}

	m.Ctr = m.Ctr.
		WithMountedFile(cfgPath, cfgFile)

	if setEnvVar != "" {
		setEnvVar = strings.ToUpper(setEnvVar)
		m.Ctr = m.Ctr.WithEnvVariable(setEnvVar, cfgPath)
	}

	return m
}
