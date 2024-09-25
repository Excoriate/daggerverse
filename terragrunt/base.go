package main

import (
	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/builderx"
)

const (
	configPresetAlpinePath = "config/presets/base-alpine.yaml"
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
func (m *Terragrunt) BaseApko(preset string) (*dagger.Container, error) {
	// Keyring configuration
	//	1. KeyURL:  "https://alpinelinux.org/keys/alpine-devel@lists.alpinelinux.org-4a6a0840.rsa.pub",
	//  2. KeyPath: "/etc/apk/keys/alpine-devel@lists.alpinelinux.org-4a6a0840.rsa.pub",
	keyRingCfg, err := builderx.GetKeyringInfoForPreset(preset)
	if err != nil {
		return nil, WrapError(err, "failed to get keyring info for preset")
	}

	// APKO preset, or configuration file path inside the container.
	apkoCfgFilePath, _ := builderx.GetApkoConfigOrPreset("", configPresetAlpinePath)
	apkoCfgFile := dag.CurrentModule().
		Source().
		File(apkoCfgFilePath)
	// APKO cache directory.
	apkoCacheDir := builderx.GetCacheDir("")
	// APKO Alpine key to mount into the container.
	keyRingAlpineKey := dag.HTTP(keyRingCfg.KeyURL)

	// Here, the APKO command is built.
	apkoBuildCmd, apkoBuildCmdErr := builderx.
		NewApkoBuilder().
		WithBuildArch(builderx.ArchX8664).
		WithBuildArch(builderx.ArchAarch64).
		WithKeyring(keyRingCfg.KeyPath).
		WithConfigFile(apkoCfgFilePath).
		WithOutputImage(apkoOutputTar).
		WithCacheDir(apkoCacheDir).
		BuildCommand()

	if apkoBuildCmdErr != nil {
		return nil, WrapError(apkoBuildCmdErr, "failed to build apko command")
	}

	// Let's build the basic container.
	ctr := dag.
		Container().
		From(apkoRepositoryURL)

	// Decorate container with APKO-related mounts.
	ctr = ctr.
		WithMountedFile(apkoCfgFilePath, apkoCfgFile).
		WithMountedFile(keyRingCfg.KeyPath, keyRingAlpineKey).
		WithMountedCache(apkoCacheDir, dag.CacheVolume("apko-cache"))

	ctr = ctr.WithExec(apkoBuildCmd)

	outputTar := ctr.File(apkoOutputTar)

	return dag.
		Container().
		Import(outputTar), nil
}
