package main

import (
	"fmt"

	"github.com/Excoriate/daggerx/pkg/envvars"
	"github.com/Excoriate/daggerx/pkg/fixtures"
)

// WithEnvVariable sets an environment variable.
//
// The name of the environment variable (e.g., "HOST").
//
// The value of the environment variable (e.g., "localhost").
//
// Replace `${VAR}` or $VAR in the value according to the current environment
// variables defined in the container (e.g., "/opt/bin:$PATH").
// +optional
func (m *Dagindag) WithEnvVariable(
	// The name of the environment variable (e.g., "HOST").
	name string,

	// The value of the environment variable (e.g., "localhost").
	value string,

	// Replace `${VAR}` or $VAR in the value according to the current environment
	// variables defined in the container (e.g., "/opt/bin:$PATH").
	// +optional
	expand bool,
) *Dagindag {
	m.Ctr = m.Ctr.WithEnvVariable(name, value, ContainerWithEnvVariableOpts{
		Expand: expand,
	})

	return m
}

// WithSource sets the source directory.
//
// The source directory is the directory that contains all the source code, including the module directory.
func (m *Dagindag) WithSource(
	// src is the directory that contains all the source code, including the module directory.
	src *Directory,
) *Dagindag {
	m.Src = src

	m.Ctr = m.Ctr.WithWorkdir(fixtures.MntPrefix).
		WithMountedDirectory(fixtures.MntPrefix, m.Src)

	return m
}

// WithEnvVarsFromStrs sets the environment variables for the container.
//
// It sets the environment variables for the container. It's meant to be used as a terminal for the module.
// Arguments:
// - envVars: A set of strings (e.g., "KEY=value,KEY=value") to use as environment variables. They're
// used to set the environment variables for the container when it's required to pass multiple environment variables
// in a single argument. E.g.: "GITHUB=token,GO_VERSION=1.22.0,MYVAR=myvar"
func (m *Dagindag) WithEnvVarsFromStrs(envVars []string) (*Dagindag, error) {
	envVarsParsed, err := envvars.ToDaggerEnvVarsFromSlice(envVars)
	if err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	for _, envVar := range envVarsParsed {
		m.Ctr = m.WithEnvVariable(envVar.Name, envVar.Value, false).Ctr
	}

	return m, nil
}

// WithDaggerSetup sets up the container with the Dagger engine.
//
// It sets up the container with the Dagger engine.
// Arguments:
// - daggerVersion: The version of the Dagger engine to use, e.g., "v0.11.6
func (m *Dagindag) WithDaggerSetup(
	// daggerVersion is the version of the Dagger engine to use, e.g., "v0.11.6
	daggerVersion string) *Dagindag {
	updatedOSCMD := []string{"sh", "-c", cmdUpdateAndInstallCURL}
	installDaggerCLI := []string{"sh", "-c", getDaggerInstallCMDByVersion(daggerVersion)}
	m.Ctr = m.Ctr.WithExec(updatedOSCMD).
		WithExec(installDaggerCLI).
		WithEnvVariable("DAGGER_VERSION", daggerVersion, ContainerWithEnvVariableOpts{
			Expand: false,
		})

	return m
}

// WithDaggerCLIEntryPoint sets the Dagger CLI entry point.
//
// It sets the Dagger CLI entry point to /bin/dagger.
// Arguments:
// - daggerCLIEntryPoint is the Dagger CLI entry point.
func (m *Dagindag) WithDaggerCLIEntryPoint() *Dagindag {
	m.Ctr = m.Ctr.WithEntrypoint(daggerCLIEntryPoint)

	return m
}

// WithDockerService sets up the container with the Docker service.
//
// It sets up the container with the Docker service.
// Arguments:
// - dockerVersion: The version of the Docker engine to use, e.g., "v20.10.17
func (m *Dagindag) WithDockerService(
	// dockerVersion is the version of the Docker engine to use, e.g., "v20.10.17"
	// +optional
	dockerVersion string,
) *Service {
	if dockerVersion == "" {
		dockerVersion = dockerVersionDefault
	}

	dindImage := getDockerInDockerImage(dockerVersion)
	dockerPort := 2375

	return dag.Container().
		From(dindImage).
		WithMountedCache(
			"/var/lib/docker",
			dag.CacheVolume(dockerVersion+"-docker-lib"),
			ContainerWithMountedCacheOpts{
				Sharing: Private,
			}).
		WithExposedPort(dockerPort).
		WithExec([]string{
			"dockerd",
			"--host=tcp://0.0.0.0:2375",
			"--host=unix:///var/run/docker.sock",
			"--tls=false",
		}, ContainerWithExecOpts{
			InsecureRootCapabilities: true,
		}).
		AsService()
}
