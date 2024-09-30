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
	// presetFile is the preset file to use for the Apko image. Either presetFile or presetFilePath must be provided.
	presetFile *dagger.File,
	// cacheDir is the cache directory to use for the Apko image.
	// +optional
	cacheDir string,
	// keyrings is the list of keyrings to use for the Apko image. They should be provided as path=url.
	// E.g.: /etc/apk/keys/alpine-devel@lists.alpinelinux.org-4a6a0840.rsa.pub=
	// https://alpinelinux.org/keys/alpine-devel@lists.alpinelinux.org-4a6a0840.rsa.pub
	// +optional
	keyrings []string,
	// buildArch specifies the architecture to build for.
	// +optional
	buildArch string,
	// buildContext is the build context directory.
	// +optional
	buildContext string,
	// debug enables debug mode for verbose output.
	// +optional
	debug bool,
	// noNetwork disables network access during the build.
	// +optional
	noNetwork bool,
	// repositoryAppend is a slice of additional repositories to append.
	// +optional
	repositoryAppend []string,
	// timestamp sets a specific timestamp for reproducible builds.
	// +optional
	timestamp string,
	// tags is a slice of additional tags for the output image.
	// +optional
	tags []string,
	// buildDate sets the build date for the APKO build.
	// +optional
	buildDate string,
	// lockfile sets the lockfile path for the APKO build.
	// +optional
	lockfile string,
	// offline enables offline mode for the APKO build.
	// +optional
	offline bool,
	// packageAppend adds extra packages to the APKO build.
	// +optional
	packageAppend []string,
	// sbom enables or disables SBOM generation.
	// +optional
	sbom bool,
	// sbomFormats sets the SBOM formats for the APKO build.
	// +optional
	sbomFormats []string,
	// sbomPath sets the SBOM output path for the APKO build.
	// +optional
	sbomPath string,
	// vcs enables or disables VCS detection.
	// +optional
	vcs bool,
	// logLevel sets the log level for the APKO build.
	// +optional
	logLevel string,
	// logPolicy sets the log policy for the APKO build.
	// +optional
	logPolicy []string,
	// workdir sets the working directory for the APKO build.
	// +optional
	workdir string,
) (*ModuleTemplateLight, error) {
	builder := m.initializeApkoBuilder(
		presetFilePath,
		cacheDir,
		keyrings,
		buildArch,
		buildContext,
		debug,
		noNetwork,
		repositoryAppend,
		timestamp,
		tags,
		buildDate,
		lockfile,
		offline,
		packageAppend,
		sbom,
		sbomFormats,
		sbomPath,
		vcs,
		logLevel,
		logPolicy,
		workdir,
	)

	cmd, err := builder.BuildCommand()
	if err != nil {
		return nil, fmt.Errorf("failed to build command: %w", err)
	}

	ctr := dag.Container().From(defaultApkoImage).WithMountedFile(presetFilePath, presetFile).WithExec(cmd)
	outputTar := ctr.File(defaultApkoTarball)

	m.Ctr = dag.Container().Import(outputTar)

	return m, nil
}

//nolint:cyclop,gocyclo // Cyclomatic complexity is high, but refactoring is not feasible at the moment.
func (m *ModuleTemplateLight) initializeApkoBuilder(
	presetFilePath, cacheDir string,
	keyrings []string, buildArch, buildContext string,
	debug, noNetwork bool, repositoryAppend []string,
	timestamp string, tags []string, buildDate, lockfile string,
	offline bool, packageAppend []string, sbom bool,
	sbomFormats []string, sbomPath string, vcs bool,
	logLevel string, logPolicy []string, workdir string,
) *builderx.ApkoBuilder {
	builder := builderx.
		NewApkoBuilder().
		WithOutputImage(defaultApkoTarball).
		WithConfigFile(presetFilePath)

	if cacheDir != "" {
		builder = builder.WithCacheDir(cacheDir)
	}

	for _, keyring := range keyrings {
		builder = builder.WithKeyring(keyring)
	}

	if buildArch != "" {
		builder = builder.WithArchitecture(buildArch)
	}

	if buildContext != "" {
		builder = builder.WithBuildContext(buildContext)
	}

	if debug {
		builder = builder.WithDebug()
	}

	if noNetwork {
		builder = builder.WithNoNetwork()
	}

	for _, repo := range repositoryAppend {
		builder = builder.WithRepositoryAppend(repo)
	}

	if timestamp != "" {
		builder = builder.WithTimestamp(timestamp)
	}

	for _, tag := range tags {
		builder = builder.WithTag(tag)
	}

	if buildDate != "" {
		builder = builder.WithBuildDate(buildDate)
	}

	if lockfile != "" {
		builder = builder.WithLockfile(lockfile)
	}

	if offline {
		builder = builder.WithOffline()
	}

	for _, pkg := range packageAppend {
		builder = builder.WithPackageAppend(pkg)
	}

	if sbom {
		builder = builder.WithSBOM(sbom)
	}

	if len(sbomFormats) > 0 {
		builder = builder.WithSBOMFormats(sbomFormats...)
	}

	if sbomPath != "" {
		builder = builder.WithSBOMPath(sbomPath)
	}

	if vcs {
		builder = builder.WithVCS(vcs)
	}

	if logLevel != "" {
		builder = builder.WithLogLevel(logLevel)
	}

	for _, policy := range logPolicy {
		builder = builder.WithLogPolicy(policy)
	}

	if workdir != "" {
		builder = builder.WithWorkdir(workdir)
	}

	return builder
}
