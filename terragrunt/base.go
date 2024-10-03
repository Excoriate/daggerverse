package main

import (
	"path/filepath"

	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/apkox"
	"github.com/Excoriate/daggerx/pkg/fixtures"
)

const (
	// Path to the APKO configuration preset for Alpine base image.
	configPresetAlpinePath = "config/presets/base-alpine.yaml"
	// Name of the output tar file generated by APKO.
	apkoOutputTar = "image.tar"
	// apkoRepositoryURL is the URL of the APKO repository.
	apkoRepositoryURL = "cgr.dev/chainguard/apko"
	containerUser     = "terragrunt"
	containerGroup    = "terragrunt"
)

// Base sets the base image and version, and creates the base container.
//
// The default image is "alpine/latest" and the default version is "latest".
//
//nolint:nolintlint,revive // This is a method that is used to set the base image and version.
func (m *Terragrunt) Base(imageURL string) *Terragrunt {
	c := dag.Container().
		From(imageURL)

	m.Ctr = c

	return m
}

// BaseApko sets up a base container using an APKO preset configuration.
//
// This function performs the following steps:
// 1. Retrieves keyring information for the given preset.
// 2. Obtains the APKO configuration file path.
// 3. Sets up the APKO cache directory.
// 4. Retrieves the Alpine key to be mounted into the container.
// 5. Builds the APKO command with the specified parameters.
// 6. Creates and decorates the container with APKO-related mounts and executes the APKO build command.
//
// Parameters:
// - preset: A string representing the APKO preset to be used.
//
// Returns:
// - *dagger.Container: A pointer to the created and configured container.
// - error: An error object if any step fails, otherwise nil.
// See: https://github.com/Excoriate/daggerx/tree/main/pkg/builderx
func (m *Terragrunt) BaseApko(extraPackages []string) (*dagger.Container, error) {
	apkoPresetFileToMount := filepath.Join(fixtures.MntPrefix, configPresetAlpinePath)
	apkoPresetFile := dag.CurrentModule().
		Source().
		File(configPresetAlpinePath)

	// APKO Alpine key to mount into the container.
	apkoCacheDir := filepath.Join(fixtures.MntPrefix, "var", "cache", "apko")

	// Here, the APKO command is built.
	apkoCmdBuider := apkox.
		NewApkoBuilder().
		WithConfigFile(apkoPresetFileToMount). // Path of the preset file mounted into the container.
		WithOutputImage(apkoOutputTar).
		WithCacheDir(apkoCacheDir)

	for _, pkg := range extraPackages {
		apkoCmdBuider.WithPackageAppend(pkg)
	}

	apkoBuildCmd, apkoBuildCmdErr := apkoCmdBuider.
		BuildCommand()

	if apkoBuildCmdErr != nil {
		return nil, WrapError(apkoBuildCmdErr, "failed to build apko command")
	}

	// Builder container with APKO preset file mounted.
	builderCtr := dag.
		Container().
		From(apkoRepositoryURL).
		WithMountedFile(apkoPresetFileToMount, apkoPresetFile).
		WithMountedCache(apkoCacheDir, dag.CacheVolume("apko-cache")) // Create a cache volume for APKO.

	builderCtr = builderCtr.WithExec(apkoBuildCmd)

	outputTar := builderCtr.File(apkoOutputTar)

	// Terragrunt container with the output tar mounted.
	tgCtr := dag.
		Container().
		Import(outputTar)

	m.Ctr = tgCtr

	return tgCtr, nil
}
