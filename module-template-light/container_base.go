package main

import (
	"fmt"

	"github.com/Excoriate/daggerverse/module-template-light/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/builderx"
)

// ApkoKeyRingInfo represents the keyring information for Apko.
type ApkoKeyRingInfo builderx.KeyringInfo

const (
	defaultAlpineImage        = "alpine"
	defaultUbuntuImage        = "ubuntu"
	defaultBusyBoxImage       = "busybox"
	defaultImageVersionLatest = "latest"
	defaultWolfiImage         = "cgr.dev/chainguard/wolfi-base"
	// Apko specifics.
	defaultApkoImage   = "cgr.dev/chainguard/apko"
	defaultApkoTarball = "image.tar"
)

var (
	// DefaultKeyringCfgAlpine is the default keyring configuration for Alpine.
	//nolint:gochecknoglobals // DefaultKeyringCfgAlpine is a global variable and is acceptable in this context.
	DefaultKeyringCfgAlpine = &ApkoKeyRingInfo{
		KeyURL:  "https://alpinelinux.org/keys/alpine-devel@lists.alpinelinux.org-4a6a0840.rsa.pub",
		KeyPath: "/etc/apk/keys/alpine-devel@lists.alpinelinux.org-4a6a0840.rsa.pub",
	}
	//nolint:gochecknoglobals // DefaultKeyringCfgWolfi is a global variable and is acceptable in this context.
	// DefaultKeyringCfgWolfi is the default keyring configuration for Wolfi.
	DefaultKeyringCfgWolfi = &ApkoKeyRingInfo{
		KeyURL:  "https://packages.wolfi.dev/os/wolfi-signing.rsa.pub",
		KeyPath: "/etc/apk/keys/wolfi-signing.rsa.pub",
	}
)

// BaseAlpine sets the base image to an Alpine Linux image and creates the base container.
//
// Parameters:
// - version: The version of the Alpine image to use. Optional parameter. Defaults to "latest".
//
// Returns a pointer to the ModuleTemplateLight instance.
func (m *ModuleTemplateLight) BaseAlpine(
	// version is the version of the Alpine image to use, e.g., "3.17.3".
	// +optional
	version string,
) *ModuleTemplateLight {
	if version == "" {
		version = defaultImageVersionLatest
	}

	imageURL := fmt.Sprintf("%s:%s", defaultAlpineImage, version)

	return m.Base(imageURL)
}

// BaseUbuntu sets the base image to an Ubuntu Linux image and creates the base container.
//
// Parameters:
// - version: The version of the Ubuntu image to use. Optional parameter. Defaults to "latest".
//
// Returns a pointer to the ModuleTemplateLight instance.
func (m *ModuleTemplateLight) BaseUbuntu(
	// version is the version of the Ubuntu image to use, e.g., "22.04".
	// +optional
	version string,
) *ModuleTemplateLight {
	if version == "" {
		version = defaultImageVersionLatest
	}

	imageURL := fmt.Sprintf("%s:%s", defaultUbuntuImage, version)

	return m.Base(imageURL)
}

// BaseBusyBox sets the base image to a BusyBox Linux image and creates the base container.
//
// Parameters:
// - version: The version of the BusyBox image to use. Optional parameter. Defaults to "latest".
//
// Returns a pointer to the ModuleTemplateLight instance.
func (m *ModuleTemplateLight) BaseBusyBox(
	// version is the version of the BusyBox image to use, e.g., "1.35.0".
	// +optional
	version string,
) *ModuleTemplateLight {
	if version == "" {
		version = defaultImageVersionLatest
	}

	imageURL := fmt.Sprintf("%s:%s", defaultBusyBoxImage, version)

	return m.Base(imageURL)
}

// BaseWolfi sets the base image to a Wolfi Linux image and creates the base container.
//
// Parameters:
// - version: The version of the Wolfi image to use. Optional parameter. Defaults to "latest".
// - packages: Additional packages to install. Optional parameter.
// - overlays: Overlay images to merge on top of the base. Optional parameter.
//
// Returns a pointer to the ModuleTemplateLight instance.
func (m *ModuleTemplateLight) BaseWolfi(
	// version is the version of the Wolfi image to use, e.g., "latest".
	// +optional
	version string,
	// packages is the list of additional packages to install.
	// +optional
	packages []string,
	// overlays are images to merge on top of the base.
	// See https://twitter.com/ibuildthecloud/status/1721306361999597884
	// +optional
	overlays []*dagger.Container,
) *ModuleTemplateLight {
	if version == "" {
		version = defaultImageVersionLatest
	}

	imageURL := fmt.Sprintf("%s:%s", defaultWolfiImage, version)

	m.Ctr = dag.
		Container().
		From(imageURL)

	// Default apk add command
	command := []string{"apk", "add", "--no-cache"}

	// Concatenate additional packages to the command
	if len(packages) > 0 {
		command = append(command, packages...)
	}

	// Install default and additional packages
	m.Ctr = m.Ctr.
		WithExec(command)

	// Apply overlays
	for _, overlay := range overlays {
		m.Ctr = m.Ctr.
			WithDirectory("/", overlay.Rootfs())
	}

	return m
}

// BaseApko sets the base image to an Apko image and creates the base container.
//
// Parameters:
// - preset: The preset to use for the Apko image. Optional parameter. Defaults to "alpine".
//
// Returns a pointer to the ModuleTemplateLight instance.
func (m *ModuleTemplateLight) BaseApko(
	// presetFilePath is the path to the preset file. Either presetFile or presetFilePath must be provided.
	presetFilePath string,
	// cacheDir is the cache directory to use for the Apko image.
	// +optional
	cacheDir string,
	// keyrings is the list of keyrings to use for the Apko image. If they aren't provided, it'll
	// be omitted completely. E.g.:https://alpinelinux.org/keys/alpine-devel@lists.alpinelinux.org-4a6a0840.rsa.pub
	// +optional
	enableDefaultKeyringAlpine bool,
	// enableDefaultKeyringWolfi is a flag to enable the default Wolfi keyring.
	// +optional
	enableDefaultKeyringWolfi bool,
	// enableArchAarch64 is a flag to enable the aarch64 architecture.
	// +optional
	enableArchAarch64 bool,
	// enableArchX8664 is a flag to enable the x86_64 architecture.
	// +optional
	enableArchX8664 bool,
	// overrideTarballName is the name of the tarball to use for the Apko image. By default, it's "image.tar".
	// +optional
	overrideTarballName string,
) (*ModuleTemplateLight, error) {
	// Handling the preset file
	if presetFilePath == "" {
		return nil, NewError("presetFilePath must be provided")
	}

	presetFile := dag.
		CurrentModule().
		Source().
		File(presetFilePath)

	// Creating builder container.
	ctr := dag.
		Container().
		From(defaultApkoImage)

	// Creating the APKO command builder (helper)
	apkoBuilder := builderx.
		NewApkoBuilder()

	// Mounting the preset file
	ctr = ctr.
		WithMountedFile(presetFilePath, presetFile)

	// Cache options - opinionated (performance vs. disk usage). Resolve later, or decide...
	apkoCacheDir := builderx.GetCacheDir("")
	if cacheDir != "" {
		apkoCacheDir = cacheDir
	}

	ctr = ctr.
		WithMountedCache(apkoCacheDir, dag.CacheVolume("apko-cache"))

	// Mounting the default Alpine keyring or Wolfi, for convenience.
	if enableDefaultKeyringAlpine {
		apkoBuilder = apkoBuilder.
			WithKeyring(DefaultKeyringCfgAlpine.KeyPath)

		ctr = ctr.
			WithMountedFile(DefaultKeyringCfgAlpine.KeyPath, dag.HTTP(DefaultKeyringCfgAlpine.KeyURL))
	}

	if enableDefaultKeyringWolfi {
		apkoBuilder = apkoBuilder.
			WithKeyring(DefaultKeyringCfgWolfi.KeyPath)

		ctr = ctr.
			WithMountedFile(DefaultKeyringCfgWolfi.KeyPath, dag.HTTP(DefaultKeyringCfgWolfi.KeyURL))
	}

	// Enabling the architectures. These two options are set through flags, so we can
	// control which architectures we want to build for.
	if enableArchAarch64 {
		apkoBuilder = apkoBuilder.
			WithBuildArch(builderx.ArchAarch64)
	}

	if enableArchX8664 {
		apkoBuilder = apkoBuilder.
			WithBuildArch(builderx.ArchX8664)
	}

	// Overriding the tarball name
	outputTarResolved := defaultApkoTarball
	if overrideTarballName != "" {
		outputTarResolved = overrideTarballName
	}

	apkoBuildCmd, apkoBuildCmdErr := apkoBuilder.
		BuildCommand()

	if apkoBuildCmdErr != nil {
		return nil, WrapError(apkoBuildCmdErr, "failed to build apko command")
	}

	ctr = ctr.
		WithExec(apkoBuildCmd)

	outputTar := ctr.File(outputTarResolved)

	m.Ctr = dag.
		Container().
		Import(outputTar)

	return m, nil
}
