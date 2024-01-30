package main

import (
	"fmt"
)

type IacTerragrunt struct {
	Ctr *Container
	SRC *Directory
	//Ctx context.Context
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

// WithSource returns the Terragrunt container with source as a mounted directory.
func (tg *IacTerragrunt) WithSource(source *Directory, enableCache Optional[bool], workDir Optional[string]) *IacTerragrunt {
	cachePathInContainer := fmt.Sprintf("%s/.terragrunt-cache", workDirDefault)
	cacheVolume := dag.CacheVolume("terragrunt-cache")

	if enableCache.GetOr(true) {
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

// WithCommands returns the Terragrunt container with the given commands.
func (tg *IacTerragrunt) WithCommands(cmds DaggerCMD, withFocus Optional[bool]) *IacTerragrunt {
	if len(cmds) == 0 {
		return tg
	}

	withFocusIsSet := withFocus.GetOr(false)

	for _, cmd := range cmds {
		if withFocusIsSet {
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
	//// ctx context.Context is the context of the command.
	//ctx Optional[context.Context],
	// Version (image tag) to use from the official image repository as a base container.
	// It's an optional parameter. If it's not set, it's going to use the default version (latest).
	version Optional[string],
	// Image to use as a base container.
	// It's an optional parameter. If it's not set, it's going to use the default image 'alpine/terragrunt'.
	image Optional[string],
	// Container to use as a base container.
	// It's an optional parameter. If it's not set, it's going to create a new container.
	container Optional[*Container],
	//enableCache bool,
	// src *Directory is the directory that contains all the source code,
	// including the module directory.
	src Optional[*Directory],
) (*IacTerragrunt, error) {
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

	//ctxSet := ctx.GetOr(context.Background())

	tg := &IacTerragrunt{
		SRC: src.GetOr(dag.Host().Directory(".")),
		Ctr: ctr,
		//Ctx: ctxSet,
		//Ctx: context.Background(),
	}

	return tg, nil
}
