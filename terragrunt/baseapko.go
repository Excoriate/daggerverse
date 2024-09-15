package main

import (
	"fmt"
	"path/filepath"

	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/fixtures"
)

// baseImagePresetPath is the path to the base image presets configuration.
const (
	baseImagePresetPath = "config/presets"
	apkoRepositoryURL   = "cgr.dev/chainguard/apko"
	apkoOutputTar       = "image.tar"
)

// BaseContainerPreset represents a preset for a base container.
type BaseContainerPreset string
type BaseContainerPresetPath string

const (
	// Alpine represents the Alpine Linux base container preset.
	Alpine BaseContainerPreset = "alpine"
	// Wolfi represents the Wolfi base container preset.
	Wolfi BaseContainerPreset = "wolfi"
)

// Binaries represents a package to be installed in the base image.
type Binaries struct {
	// BaseURL is the URL where the package can be downloaded.
	BaseURL string
	// Version is the version of the package to install.
	// If omitted, the latest version is installed.
	Version string
}

// BaseImageApko represents a base image with a name and a list of packages to be installed.
type BaseImageApko struct {
	// Name is the name of the base image.
	Name string
	// Cfg is the preset or configuration file to use for the base image. For more
	// documentation about its format, please refer to the apko documentation:
	// https://github.com/chainguard-images/apko or https://github.com/chainguard-dev/apko/blob/main/docs/apko_file.md
	Cfg *dagger.File
	// Binaries is a list of binaries to be installed in the base image.
	Binaries []Binaries
}

func (b *BaseImageApko) getCfgPresetFile(name string) *dagger.File {
	return dag.CurrentModule().
		Source().
		File(fmt.Sprintf("%s/%s.yaml", baseImagePresetPath, name))
}

func (b *BaseImageApko) getCfgPresetFilenameInCtr() string {
	return filepath.Join(fixtures.MntPrefix, "config.yaml")
}

func (b *BaseImageApko) getAPkoCacheDir() string {
	return filepath.Join(fixtures.MntPrefix, "var", "cache", "apko")
}

func (b *BaseImageApko) getAPkoOutputTarFilePathInCtr() string {
	return filepath.Join(fixtures.MntPrefix, apkoOutputTar)
}

// WithApkoCfgPresetFile sets the configuration preset file for the base image.
// The configuration preset file is defined as a .yaml file and is used to configure
// the base image. For more information about its format, please refer to the apko documentation:
// https://github.com/chainguard-images/apko or
// https://github.com/chainguard-dev/apko/blob/main/docs/apko_file.md
//
// Parameters:
//   - cfg: The base config (or preset) defined as a .yaml file to use for the base image.
//     This parameter is optional.
//
// Returns:
//   - A pointer to a dagger.Container with the applied configuration preset file.
func (b *BaseImageApko) WithApkoCfgPresetFile(
	// cfg is the base config (or preset) defined as a .yaml file
	// to use for the base image. For more information about its format, please
	// refer to the apko documentation:
	// https://github.com/chainguard-images/apko or
	// https://github.com/chainguard-dev/apko/blob/main/docs/apko_file.md
	// +optional
	cfg string) *dagger.Container {
	if cfg == "" {
		b.Cfg = b.getCfgPresetFile("alpine")
	} else {
		b.Cfg = b.getCfgPresetFile(cfg)
	}

	return nil
}

func (b *BaseImageApko) BuildApko(
	// cfgPreset is the name of the base image to build. if not provided, it'll use
	// Wolfi as the default configuration. Valid values are 'alpine' and 'wolfi'.
	// If you pass a custom configuration file, it must be a valid apko configuration, and should
	// be a .yaml file (placed inside the Context of the module).
	// +optional
	cfgPreset string,
	// binaries is a list of binaries to be installed in the base image.
	// +optional
	binaries []Binaries,
) (*dagger.Container, error) {
	// It handles already nullable cfg.
	cfg := b.getCfgPresetFile(cfgPreset)
	cfgPath := b.getCfgPresetFilenameInCtr()
	cacheDir := b.getAPkoCacheDir()
	imageOutputTar := b.getAPkoOutputTarFilePathInCtr()

	buildArgs := []string{
		"apko",
		"build",
		cfgPath,
		"latest",
		"--cache-dir",
		b.getAPkoCacheDir(),
	}

	// Base container from default repository URL.
	ctr := dag.
		Container().
		From(apkoRepositoryURL)

	// Mount in the container, the preset file, which's going to be treated as the config.yaml file.
	ctr = ctr.WithMountedFile(cfgPath, cfg).
		WithMountedCache(cacheDir,
			dag.CacheVolume("apko-cache"))

	// Generate the .tar output of the image built.
	outputTar := ctr.
		WithExec(buildArgs).
		File(imageOutputTar)

	return dag.
		Container().
		Import(outputTar), nil
}
