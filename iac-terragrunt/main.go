package main

import (
	"fmt"
	"time"
)

type IacTerragrunt struct {
	Ctr *Container
	SRC *Directory
}

// Container returns the container of IacTerragrunt.
func (tg *IacTerragrunt) Container() *Container {
	return tg.Ctr
}

// WithVersion returns the Terragrunt container with a given Terragrunt version.
func (tg *IacTerragrunt) WithVersion(version string) *IacTerragrunt {
	if tg.Ctr == nil {
		tg.Ctr = dag.Container()
	}

	tg.Ctr = tg.Ctr.
		From(fmt.Sprintf("%s:%s", defaultContainerImage, version))
	return tg
}

// WithContainer returns the Terragrunt container with the given container.
func (tg *IacTerragrunt) WithContainer(ctr *Container) *IacTerragrunt {
	tg.Ctr = ctr
	return tg
}

// WithModule returns the Terragrunt container with the given Terragrunt module.
func (tg *IacTerragrunt) WithModule(module string) *IacTerragrunt {
	tg.Ctr = tg.Ctr.
		WithWorkdir(fmt.Sprintf("%s/%s", workDirDefault, module))

	return tg
}

// WithEntrypoint returns the Terragrunt container with the given entry point.
func (tg *IacTerragrunt) WithEntrypoint(entryPoint []string) *IacTerragrunt {
	tg.Ctr = tg.Ctr.WithEntrypoint(entryPoint)
	return tg
}

// WithEnvVar returns the Terragrunt container with the given environment variable.
func (tg *IacTerragrunt) WithEnvVar( // The name of the environment variable (e.g., "HOST").
	name string,

	// The value of the environment variable (e.g., "localhost").
	value string,

	// Replace `${VAR}` or $VAR in the value according to the current environment
	// variables defined in the container (e.g., "/opt/bin:$PATH").
	// +optional
	expand bool,
) *IacTerragrunt {
	return &IacTerragrunt{
		Ctr: tg.Ctr.WithEnvVariable(name, value, ContainerWithEnvVariableOpts{
			Expand: expand,
		}),
	}
}

// WithSource returns the Terragrunt container with source as a mounted directory.
func (tg *IacTerragrunt) WithSource(source *Directory, enableCache Optional[bool], workDir Optional[string]) *IacTerragrunt {
	cachePathInContainer := fmt.Sprintf("%s/.terragrunt-cache", workDirDefault)
	cacheVolume := dag.CacheVolume("terragrunt-cache")

	if enableCache.GetOr(false) {
		tg.Ctr = tg.Ctr.
			WithMountedCache(cachePathInContainer, cacheVolume)
	}

	var workDirToSet string

	if !workDir.isSet {
		workDirToSet = workDirDefault
	} else {
		workDirToSet = fmt.Sprintf("%s/%s", workDirDefault, workDir.value)
	}

	tg.Ctr = tg.Ctr.
		WithWorkdir(workDirToSet).
		WithMountedDirectory(workDirDefault, source)

	return tg
}

// WithSecret returns the Terragrunt container with the given secrets.
func (tg *IacTerragrunt) WithSecret(name, value string) *IacTerragrunt {
	secret := dag.SetSecret(name, value)
	tg.Ctr = tg.Ctr.WithSecretVariable(name, secret)

	return tg
}

// WithGitSSHConfig returns the Terragrunt container with the given Git SSH configuration.
func (tg *IacTerragrunt) WithGitSSHConfig(sshAuthSock string) *IacTerragrunt {
	tg.Ctr = tg.WithEnvVar("GIT_SSH_COMMAND", "ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=accept-new", false).
		WithEnvVar("SSH_AUTH_SOCK", sshAuthSock, false).Ctr

	return tg
}

// WithCacheInvalidation returns the Terragrunt container with cache invalidation.
func (tg *IacTerragrunt) WithCacheInvalidation() *IacTerragrunt {
	tg.Ctr = tg.WithEnvVar("CACHEBUSTER", time.Now().String(), false).Ctr
	return tg
}

// WithCommands returns the Terragrunt container with the given commands.
func (tg *IacTerragrunt) WithCommands(cmds DaggerCMD, withFocus bool) *IacTerragrunt {
	if len(cmds) == 0 {
		return tg
	}

	for _, cmd := range cmds {
		if withFocus {
			tg.Ctr = tg.Ctr.
				WithFocus().
				WithExec(cmd)
		} else {
			tg.Ctr = tg.Ctr.
				WithExec(cmd)
		}
	}

	return tg
}

// New creates a new instance of IacTerragrunt.
// If no image is specified, the default image will be used.
// If no version is specified, the default version will be used.
// If no container is specified, a new container will be created.
func New(
	// Version (image tag) to use from the official image repository as a base container.
	// It's an optional parameter. If it's not set, it's going to use the default version (latest).
	version Optional[string],
	// Image to use as a base container.
	// It's an optional parameter. If it's not set, it's going to use the default image 'alpine/terragrunt'.
	image Optional[string],
	// Container to use as a base container.
	// It's an optional parameter. If it's not set, it's going to create a new container.
	container Optional[*Container],
	// src *Directory is the directory that contains all the source code,
	// including the module directory.
	src Optional[*Directory],
) *IacTerragrunt {
	var ctr *Container
	var versionResolved string
	var imageResolved string

	versionResolved = version.GetOr(defaultContainerVersion)
	imageResolved = image.GetOr(defaultContainerImage)

	baseImage := getContainerImage(imageResolved, versionResolved)

	if container.isSet {
		ctr = container.value
	} else {
		ctr = dag.Container().From(baseImage)
	}

	tg := &IacTerragrunt{
		SRC: src.GetOr(dag.CurrentModule().Source().Directory(".")),
		Ctr: ctr,
	}

	return tg
}
