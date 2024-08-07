// Package main provides utility functions for working with Docker containers.
//
// This package includes constants and a function for getting the Docker-in-Docker image.
//
// Copyright: Excoriate
// License: Apache-2.0
package main

import "fmt"

const (
	// defaultContainerVersion specifies the default version for the container.
	defaultContainerVersion = "latest"
	// defaultContainerImage specifies the default image for the container.
	defaultContainerImage = "alpine"
	// dockerVersionDefault specifies the default Docker version.
	dockerVersionDefault = "24.0"
)

// getDockerInDockerImage returns the Docker-in-Docker image with the given version.
//
// If the version is not provided, it defaults to dockerVersionDefault.
//
// Example:
//
//	getDockerInDockerImage("20.10.17") => "docker:20.10.17-dind"
func getDockerInDockerImage(version string) string {
	if version == "" {
		version = dockerVersionDefault
	}

	return fmt.Sprintf("docker:%s-dind", version)
}
