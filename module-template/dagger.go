package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/module-template/internal/dagger"
	"golang.org/x/mod/semver"
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
			"--dns", "8.8.8.8",
			"--dns", "8.8.4.4",
			"--host=tcp://0.0.0.0:2375",
			"--host=unix:///var/run/docker.sock",
			"--tls=false",
		}, dagger.ContainerWithExecOpts{
			InsecureRootCapabilities: true,
		}).
		AsService()
}

func (m *ModuleTemplate) validateDaggerVersion(dagVersion string) error {
	if dagVersion == "" {
		return WrapError(nil, "empty dagVersion")
	}

	// If version is lower than 0.12.0, it's not supported.
	if semver.Compare(dagVersion, "0.12.0") < 0 {
		return Errorf("unsupported dagger version %s, it must be greater "+
			"than or equal to 0.12.0", dagVersion)
	}

	return nil
}

// SetupDaggerInDagger sets up the Dagger CLI inside a Docker-in-Docker context.
//
// This method validates the specified Dagger Engine version, installs the Dagger CLI
// in an Alpine-based container, replaces the existing container with a Docker-in-Docker
// container, and sets up the Docker service within the Dagger context.
//
// Parameters:
//   - dagVersion: The version of the Dagger Engine to use, e.g., "v0.12.5".
//   - dockerVersion: The version of the Docker Engine to use, e.g., "24.0". This is optional.
//
// Returns:
//   - *ModuleTemplate: The modified ModuleTemplate instance with Dagger and Docker configured.
//   - error: An error if the setup process fails at any step.
func (m *ModuleTemplate) SetupDaggerInDagger(
	dagVersion string, // The version of the Dagger Engine to use.
	dockerVersion string, // The version of the Docker Engine to use. +optional
) (*ModuleTemplate, error) {
	// Validate the specified Dagger Engine version.
	if err := m.validateDaggerVersion(dagVersion); err != nil {
		return nil, WrapErrorf(err, "failed to validate dagger version %s", dagVersion)
	}

	// Setup Docker service within the Dagger context.
	dockerd := m.WithDaggerDockerService(dockerVersion)
	dockerHost, dockerHostErr := dockerd.Endpoint(context.Background(),
		dagger.ServiceEndpointOpts{
			Scheme: "tcp",
		})

	if dockerHostErr != nil {
		return nil, WrapError(dockerHostErr, "failed to get docker host")
	}

	// Get the Docker-in-Docker image
	dindImage := getDockerInDockerImage(dockerVersion)

	// Set up the container with the Docker-in-Docker image
	m.Ctr = dag.Container().From(dindImage)

	// Install necessary packages
	m.Ctr = m.Ctr.WithExec([]string{"apk", "add", "--no-cache", "git", "curl"})

	// Configure DNS
	m.Ctr = m.Ctr.
		WithEnvVariable("GODEBUG", "netdns=go").
		WithEnvVariable("DNS_RESOLVER", "8.8.8.8 8.8.4.4")

	// Set up the Dagger CLI in an Alpine-based container.
	m.WithDaggerCLIAlpine(dagVersion)

	// Configure Git
	m.Ctr = m.Ctr.WithExec([]string{"git", "config", "--global", "http.https://gopkg.in.followRedirects", "true"})

	// Bind the Docker service and set the DOCKER_HOST environment variable.
	m.Ctr = m.Ctr.
		WithServiceBinding("docker", dockerd).
		WithEnvVariable("DOCKER_HOST", dockerHost)

	// Set the entrypoint to the Dagger binary.
	m.Ctr = m.Ctr.WithEntrypoint([]string{"/bin/dagger"})

	return m, nil
}