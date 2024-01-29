package main

import (
	"fmt"
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

// WithSource returns the Terragrunt container with source as a mounted directory.
func (tg *IacTerragrunt) WithSource(source *Directory, enableCache Optional[bool]) *IacTerragrunt {
	tg.SRC = source

	cachePathInContainer := fmt.Sprintf("%s/.terragrunt-cache", workDirDefault)
	cacheVolume := dag.CacheVolume("terragrunt-cache")

	if enableCache.GetOr(true) {
		tg.Ctr = tg.Ctr.WithMountedCache(cachePathInContainer, cacheVolume)
	}

	tg.Ctr = tg.Ctr.
		WithWorkdir(workDirDefault).
		WithMountedDirectory(workDirDefault, source)

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

	tg := &IacTerragrunt{}

	tg.SRC = src.GetOr(dag.Host().Directory("."))

	tg.Ctr = ctr
	return tg, nil
}

// RunTG just run an arbitrary Terragrunt commands. The entry point is 'Terragrunt'.
//func (tg *IacTerragrunt) RunTG(cmds []string, src Optional[*Directory], module Optional[string]) (*Container, error) {
//	if len(cmds) == 0 {
//		return nil, fmt.Errorf("command cannot be empty")
//	}
//
//	srcDirToMount, err := isEitherGlobalOrArgSetDir(src, tg.GlobalSRC)
//	if err != nil {
//		return nil, err
//	}
//
//	var moduleToWorkOn string
//
//	if !module.isSet {
//		if tg.GlobalTGModule == "" {
//			return nil, fmt.Errorf("module directory cannot be empty, and it was not set in the constructor")
//		}
//	} else {
//		moduleToWorkOn = module.value
//	}
//
//	workDirWithModule := path.Join(workDirDefault, moduleToWorkOn)
//
//	tg.Ctr = tg.Ctr.
//		WithWorkdir(workDirWithModule).
//		WithMountedDirectory(workDirDefault, srcDirToMount)
//
//	tg.Ctr = addTGCommandsToContainer(tg.Ctr, cmds)
//
//	return tg.Ctr, nil
//}

// TGInit initializes a Terragrunt module.
//func (tg *IacTerragrunt) TGInit(src Optional[*Directory], module Optional[string], args Optional[[]string]) (*Container, error) {
//	srcDirToMount, err := isEitherGlobalOrArgSetDir(src, tg.GlobalSRC)
//	if err != nil {
//		return nil, err
//	}
//
//	moduleToWorkOn, err := isEitherGlobalOrArgSetString(module, tg.GlobalTGModule)
//	if err != nil {
//		return nil, err
//	}
//
//	workDirWithModule := path.Join(workDirDefault, moduleToWorkOn)
//
//	tg.Ctr = tg.Ctr.
//		WithWorkdir(workDirWithModule).
//		WithMountedDirectory(workDirDefault, srcDirToMount)
//
//	var argsToSet []string
//
//	if args.isSet {
//		argsToSet = args.value
//	}
//
//	tg.Ctr = addTGCommandsToContainer(tg.Ctr, []string{"init"}, argsToSet...)
//
//	return tg.Ctr, nil
//}
