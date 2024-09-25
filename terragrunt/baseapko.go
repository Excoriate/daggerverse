package main

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/fixtures"
)

// Constants
const (
	baseImagePresetPath = "config/presets"
	apkoRepositoryURL   = "cgr.dev/chainguard/apko"
	apkoOutputTar       = "image.tar"
)

// BaseContainerPreset represents a preset for a base container.
// It defines the type of base image to be used for building containers.
type BaseContainerPreset string

const (
	// Alpine represents the Alpine Linux base image preset.
	// Alpine is a lightweight Linux distribution that's commonly used in containers.
	Alpine BaseContainerPreset = "alpine"

	// Wolfi represents the Wolfi Linux base image preset.
	// Wolfi is a community Linux OS designed for containers, offering enhanced security features.
	// It's optimized for use with apko and melange for building OCI images.
	Wolfi BaseContainerPreset = "wolfi"
)

// BaseImageApkoBuilder is the builder for BaseImageApko
type BaseImageApko struct {
	name BaseContainerPreset
	cfg  *dagger.File
}

type BaseImageApkoOptionFn func(*BaseImageApko) error

// NewBaseImageApkoBuilder creates a new BaseImageApkoBuilder
func NewBaseImageApko(opts ...BaseImageApkoOptionFn) (*BaseImageApko, error) {
	b := &BaseImageApko{}

	for _, opt := range opts {
		if err := opt(b); err != nil {
			return nil, err
		}
	}
	return b, nil
}

// Start of Selection
// WithApkoPreset configures the BaseImageApko builder by setting the name of the base image.
// The provided name must be either "alpine" or "wolfi". If the name is empty or does not
// match one of the allowed values, an error is returned.
//
// Parameters:
//   - name: A string representing the name of the base image. Must be "alpine" or "wolfi".
//
// Returns:
//   - A BaseImageApkoOptionFn that applies the name configuration to a BaseImageApko instance.
func WithApkoPreset(name string) BaseImageApkoOptionFn {
	return func(b *BaseImageApko) error {
		if name == "" {
			return errors.New("name is required, either 'alpine' or 'wolfi'")
		}

		b.cfg = dag.CurrentModule().
			Source().
			File(fmt.Sprintf("%s/%s.yaml", baseImagePresetPath, name))

		b.name = BaseContainerPreset(name)

		return nil
	}
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

// BuildImage builds the Apko image based on the specified configuration preset.
//
// This function retrieves the configuration preset file, sets up the necessary
// cache directories, and executes the Apko build command within a Dagger container.
// It uses either Alpine or Wolfi specific configurations based on the preset name.
//
// Returns:
//   - A pointer to a dagger.Container containing the built Apko image.
//   - An error if the build process fails.
func (b *BaseImageApko) BuildImage() (*dagger.Container, error) {
	cfgPath := b.getCfgPresetFilenameInCtr()
	cacheDir := b.getAPkoCacheDir()
	imageOutputTar := b.getAPkoOutputTarFilePathInCtr()

	keyring, err := b.getKeyringInfo()
	if err != nil {
		return nil, err
	}

	key := dag.HTTP(keyring.KeyURL)

	buildArgs := []string{
		"apko",
		"build",
		"--keyring-append", keyring.KeyPath,
		"--arch", "x86_64",
		"--arch", "aarch64",
		cfgPath,
		"latest",
		imageOutputTar,
		"--cache-dir",
		cacheDir,
	}

	ctr := dag.
		Container().
		From(apkoRepositoryURL).
		WithMountedFile(cfgPath, b.cfg).
		WithMountedFile(keyring.KeyPath, key).
		WithMountedCache(cacheDir,
			dag.CacheVolume("apko-cache"))

	ctr = ctr.WithExec(buildArgs)
	outputTar := ctr.File(imageOutputTar)

	return dag.
		Container().
		Import(outputTar), nil
}

type keyringInfo struct {
	KeyURL  string
	KeyPath string
}

func (b *BaseImageApko) getKeyringInfo() (keyringInfo, error) {
	switch b.name {
	case "alpine":
		return keyringInfo{
			KeyURL:  "https://alpinelinux.org/keys/alpine-devel@lists.alpinelinux.org-4a6a0840.rsa.pub",
			KeyPath: "/etc/apk/keys/alpine-devel@lists.alpinelinux.org-4a6a0840.rsa.pub",
		}, nil
	case "wolfi":
		return keyringInfo{
			KeyURL:  "https://packages.wolfi.dev/os/wolfi-signing.rsa.pub",
			KeyPath: "/etc/apk/keys/wolfi-signing.rsa.pub",
		}, nil
	default:
		return keyringInfo{}, fmt.Errorf("unsupported preset: %s", b.name)
	}
}
