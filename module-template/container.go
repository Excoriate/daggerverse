package main

import (
	"fmt"

	"github.com/Excoriate/daggerverse/module-template/internal/dagger"
)

const (
	defaultAlpineImage        = "alpine"
	defaultUbuntuImage        = "ubuntu"
	defaultBusyBoxImage       = "busybox"
	defaultImageVersionLatest = "latest"
)

// BaseAlpine sets the base image to an Alpine Linux image and creates the base container.
//
// Parameters:
// - version: The version of the Alpine image to use. Optional parameter. Defaults to "latest".
//
// Returns a pointer to the ModuleTemplate instance.
func (m *ModuleTemplate) BaseAlpine(
	// version is the version of the Alpine image to use, e.g., "3.17.3".
	// +optional
	version string,
) *ModuleTemplate {
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
// Returns a pointer to the ModuleTemplate instance.
func (m *ModuleTemplate) BaseUbuntu(
	// version is the version of the Ubuntu image to use, e.g., "22.04".
	// +optional
	version string,
) *ModuleTemplate {
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
// Returns a pointer to the ModuleTemplate instance.
func (m *ModuleTemplate) BaseBusyBox(
	// version is the version of the BusyBox image to use, e.g., "1.35.0".
	// +optional
	version string,
) *ModuleTemplate {
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
// Returns a pointer to the ModuleTemplate instance.
func (m *ModuleTemplate) BaseWolfi(
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
) *ModuleTemplate {
	if version == "" {
		version = defaultImageVersionLatest
	}

	imageURL := fmt.Sprintf("%s:%s", "cgr.dev/chainguard/wolfi-base", version)

	m.Ctr = dag.
		Container().
		From(imageURL)

	// Install default and additional packages
	m.Ctr = m.Ctr.
		WithExec([]string{"apk", "add", "--no-cache"}).
		WithExec(packages)

	// Apply overlays
	for _, overlay := range overlays {
		m.Ctr = m.Ctr.
			WithDirectory("/", overlay.Rootfs())
	}

	return m
}
