// Package main provides methods for setting up and managing a container environment.
// This includes setting environment variables, mounting directories and files,
// and configuring services such as Docker within the container.
//
// Copyright: Excoriate alex_torres@outlook.com
// License: MIT
package main

import (
	"path/filepath"
	"strings"
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
	// owner is the owner of the directory. If not set it'll default to terragrunt
	// +optional
	owner string,
) *Terragrunt {
	var ctr *dagger.Container

	// Set the owner if provided, otherwise use default mounting
	mountOpts := dagger.ContainerWithMountedDirectoryOpts{}
	if owner != "" {
		mountOpts.Owner = owner
		ctr = m.Ctr.WithMountedDirectory(fixtures.MntPrefix, src, mountOpts)
	} else {
		ctr = m.Ctr.WithMountedDirectory(fixtures.MntPrefix, src)
	}

	// Set the working directory if provided, otherwise use default
	if workdir != "" {
		ctr = ctr.WithWorkdir(filepath.Join(fixtures.MntPrefix, workdir))
	} else {
		ctr = ctr.WithWorkdir(fixtures.MntPrefix)
	}

	m.Ctr = ctr

	return m
}

// WithUserAsOwnerOfDirs sets the specified user (and optionally group) as the owner of the given directories within the container.
//
// This method iterates over the provided list of directories and executes the "chown" command
// within the container to change the ownership of each directory to the specified user (and optionally group).
//
// Parameters:
//   - user: The user to set as the owner of the directories. This should be a valid user within the container.
//   - group: The group to set as the owner of the directories. This is an optional parameter.
//   - dirs: A slice of strings representing the directories to set the owner of. Each directory path should be
//     specified relative to the container's filesystem.
//
// Returns:
// - *Terragrunt: The updated Terragrunt instance with the ownership of the specified directories changed.
func (m *Terragrunt) WithUserAsOwnerOfDirs(
	// user is the user to set as the owner of the directories.
	user string,
	// group is the group to set as the owner of the directories.
	// +optional
	group string,
	// dirs is the directories to set the owner of.
	dirs []string,
	// configureAsRoot is whether to configure the directories as root, and then it'll use the given user.
	// +optional
	configureAsRoot bool,
) *Terragrunt {
	if configureAsRoot {
		m.Ctr = m.Ctr.WithUser("root")
	}

	for _, dir := range dirs {
		if group != "" {
			m.Ctr = m.Ctr.WithExec([]string{"chown", "-R", user + ":" + group, dir})
		} else {
			m.Ctr = m.Ctr.WithExec([]string{"chown", "-R", user, dir})
		}
	}

	if configureAsRoot {
		m.Ctr = m.Ctr.WithUser(user)
	}

	return m
}

// WithUserWithPermissionsOnDirs sets the specified permissions on the given directories within the container.
//
// This method iterates over the provided list of directories and executes the "chmod" command
// within the container to change the permissions of each directory to the specified mode.
//
// Parameters:
//   - user: The user to set as the owner of the directories. This should be a valid user within the container.
//   - mode: The permissions to set on the directories. This should be a valid mode string (e.g., "0777").
//   - dirs: A slice of strings representing the directories to set the permissions of. Each directory path should be
//     specified relative to the container's filesystem.
//   - configureAsRoot: Whether to configure the directories as root, and then it'll use the given mode.
//
// Returns:
// - *Terragrunt: The updated Terragrunt instance with the permissions of the specified directories changed.
func (m *Terragrunt) WithUserWithPermissionsOnDirs(
	// user is the user to set as the owner of the directories.
	// +optional
	user string,
	// mode is the permissions to set on the directories.
	mode string,
	// dirs is the directories to set the permissions of.
	dirs []string,
	// configureAsRoot is whether to configure the directories as root, and then it'll use the given mode.
	// +optional
	configureAsRoot bool,
) *Terragrunt {
	if configureAsRoot {
		m.Ctr = m.Ctr.WithUser("root")
	}

	for _, dir := range dirs {
		m.Ctr = m.Ctr.WithExec([]string{"chmod", "-R", mode, dir})
	}

	if configureAsRoot && user != "" {
		m.Ctr = m.Ctr.WithUser(user)
	}

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

// WithCachedDirectory mounts a cache volume in the container.
//
// Parameters:
// - path: The path in the container where the cache volume will be mounted.
// - cacheVolume: The cache volume to mount.
func (m *Terragrunt) WithCachedDirectory(
	// path is the path in the container where the cache volume will be mounted.
	path string,
	// enablePrefixWithMountPath is whether to enable the prefix with the mount path.
	// +optional
	enablePrefixWithMountPath bool,
	// setEnvVarWithCacheDirValue is the value to set the cache directory in the container.
	// +optional
	setEnvVarWithCacheDirValue string,
	// cacheSharingMode is the sharing mode of the cache volume.
	// +optional
	cacheSharingMode dagger.CacheSharingMode,
	// cacheVolumeRootDir is the root directory of the cache volume.
	// +optional
	cacheVolumeRootDir *dagger.Directory,
	// cacheVolumeOwner is the owner of the cache volume.
	// +optional
	cacheVolumeOwner string,
) *Terragrunt {
	// Define the cache volume
	cacheVolume := dag.CacheVolume(path)

	// Define the path in the container
	mountPath := path
	if enablePrefixWithMountPath {
		mountPath = filepath.Join(fixtures.MntPrefix, path)
	}

	cacheMountOpts := dagger.ContainerWithMountedCacheOpts{}
	if cacheVolumeRootDir != nil {
		cacheMountOpts.Source = cacheVolumeRootDir
	}

	if cacheVolumeOwner != "" {
		cacheMountOpts.Owner = cacheVolumeOwner
	}

	if cacheSharingMode != "" {
		cacheMountOpts.Sharing = cacheSharingMode
	}

	// Mount the cache volume in the container
	m.Ctr = m.Ctr.
		WithMountedCache(mountPath, cacheVolume, cacheMountOpts)

	// Set the environment variable if provided
	if setEnvVarWithCacheDirValue != "" {
		envVarName := strings.ToUpper(strings.TrimSpace(setEnvVarWithCacheDirValue))
		m.Ctr = m.
			Ctr.
			WithEnvVariable(envVarName, mountPath)
	}

	return m
}
