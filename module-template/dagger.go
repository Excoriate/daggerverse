package main

import (
	"strings"

	"github.com/Excoriate/daggerverse/module-template/internal/dagger"
)

// getDaggerInstallCMDByVersion returns the command to install the Dagger engine.
//
// The command is a shell script that sets the DAGGER_VERSION environment variable
// and then downloads and runs the Dagger install script for the specific version.
//
// Example:
//
//	getDaggerInstallCMDByVersion("v0.12.1")
//	=> `cd / && DAGGER_VERSION="v0.12.1" curl -L https://dl.dagger.io/dagger/install.sh | DAGGER_VERSION="v0.12.1" sh`
func getDaggerInstallCMDByVersion(version string) string {
	return strings.Join([]string{
		"cd /",
		"&&",
		"DAGGER_VERSION=\"" + version + "\"",
		"curl -L https://dl.dagger.io/dagger/install.sh |",
		"DAGGER_VERSION=\"" + version + "\"",
		"sh",
	}, " ")
}

// WithDaggerCLIAlpine sets up the Dagger CLI entry point for Alpine within the ModuleTemplate.
//
// Parameters:
//   - version: The version of the Dagger Engine to use, e.g., "v0.12.1".
//
// This method performs the following steps:
//  1. Generates a shell command to install the Dagger CLI using the specified version.
//  2. Executes the installation command within the Alpine container context.
//  3. Sets the DAGGER_VERSION environment variable in the container.
//
// Returns:
//   - *ModuleTemplate: Returns the modified ModuleTemplate instance with the Dagger CLI configured.
func (m *ModuleTemplate) WithDaggerCLIAlpine(version string) *ModuleTemplate {
	daggerInstallCommand := getDaggerInstallCMDByVersion(version)
	installDaggerCLI := []string{"sh", "-c", daggerInstallCommand}

	m.Ctr = m.Ctr.
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec(installDaggerCLI).
		WithEnvVariable("DAGGER_VERSION", version,
			dagger.ContainerWithEnvVariableOpts{
				Expand: false,
			})

	return m
}

// WithDaggerCLIUbuntu sets up the Dagger CLI entry point for Ubuntu within the ModuleTemplate.
//
// Parameters:
//   - version: The version of the Dagger Engine to use, e.g., "v0.12.1".
//
// This method performs the following steps:
//  1. Updates package lists and installs curl.
//  2. Generates a shell command to install the Dagger CLI using the specified version.
//  3. Executes the installation command within the Ubuntu container context.
//  4. Sets the DAGGER_VERSION environment variable in the container.
//
// Returns:
//   - *ModuleTemplate: Returns the modified ModuleTemplate instance with the Dagger CLI configured.
func (m *ModuleTemplate) WithDaggerCLIUbuntu(version string) *ModuleTemplate {
	daggerInstallCommand := getDaggerInstallCMDByVersion(version)
	installDaggerCLI := []string{"bash", "-c", daggerInstallCommand}

	m.Ctr = m.Ctr.
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "curl"}).
		WithExec(installDaggerCLI).
		WithEnvVariable("DAGGER_VERSION", version,
			dagger.ContainerWithEnvVariableOpts{
				Expand: false,
			})

	return m
}

// WithDaggerDockerService sets up the container with the Docker service.
//
// Arguments:
//   - version: The version of the Docker engine to use, e.g., "v20.10.17".
//     If empty, a default version will be used.
//
// Returns:
//   - *dagger.Service: A Dagger service configured with Docker.
func (m *ModuleTemplate) WithDaggerDockerService(version string) *dagger.Service {
	if version == "" {
		version = dockerVersionDefault
	}

	dindImage := getDockerInDockerImage(version)
	dockerPort := 2375

	return dag.Container().
		From(dindImage).
		WithMountedCache(
			"/var/lib/docker",
			dag.CacheVolume(version+"-docker-lib"),
			dagger.ContainerWithMountedCacheOpts{
				Sharing: dagger.Private,
			}).
		WithExposedPort(dockerPort).
		WithExec([]string{
			"dockerd",
			"--host=tcp://0.0.0.0:2375",
			"--host=unix:///var/run/docker.sock",
			"--tls=false",
		}, dagger.ContainerWithExecOpts{
			InsecureRootCapabilities: true,
		}).
		AsService()
}
