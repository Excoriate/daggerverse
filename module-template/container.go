package main

import "fmt"

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
