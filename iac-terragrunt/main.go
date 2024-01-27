package main

import (
	"fmt"
	"path"
	"strings"
)

type IacTerragrunt struct {
	Ctr            *Container
	GlobalSRC      *Directory
	GlobalTGModule string
	workdir        string
}

// Container returns the container of IacTerragrunt.
func (tg *IacTerragrunt) Container() *Container {
	return tg.Ctr
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
	// module is the module directory where the terragrunt commands will be executed.
	// It's an optional parameter. If it's not set, it's going to use the default module directory (current directory).
	module Optional[string],
) (*IacTerragrunt, error) {
	var ctr *Container
	var versionResolved string
	var imageResolved string

	if version.isSet {
		versionResolved = strings.TrimSpace(version.value)
	} else {
		versionResolved = defaultContainerVersion
	}

	if image.isSet {
		imageResolved = strings.TrimSpace(image.value)
	} else {
		imageResolved = defaultContainerImage
	}

	baseImage := FormImageAddress(imageResolved, versionResolved)

	if container.isSet {
		ctr = container.value
	} else {
		ctr = dag.Container().From(baseImage)
	}

	tg := &IacTerragrunt{}

	if src.isSet {
		tg.GlobalSRC = src.value
	}

	if module.isSet {
		tg.GlobalTGModule = module.value
	}

	tg.workdir = workDirDefault

	tg.Ctr = ctr
	return tg, nil
}

// Run just run an arbitrary commands.
func (tg *IacTerragrunt) Run(cmds []string, src Optional[*Directory]) (*Container, error) {
	if len(cmds) == 0 {
		return nil, fmt.Errorf("command cannot be empty")
	}

	srcDirToMount, err := isEitherGlobalOrArgSetDir(src, tg.GlobalSRC)
	if err != nil {
		return nil, err
	}

	tg.Ctr = tg.Ctr.
		WithWorkdir(workDirDefault).
		WithMountedDirectory(workDirDefault, srcDirToMount).
		WithEntrypoint(nil)

	tg.Ctr = addCMDsToContainer(tg.Ctr, cmds)

	return tg.Ctr, nil
}

// RunTG just run an arbitrary Terragrunt commands. The entry point is 'Terragrunt'.
func (tg *IacTerragrunt) RunTG(cmds []string, src Optional[*Directory], module Optional[string]) (*Container, error) {
	if len(cmds) == 0 {
		return nil, fmt.Errorf("command cannot be empty")
	}

	srcDirToMount, err := isEitherGlobalOrArgSetDir(src, tg.GlobalSRC)
	if err != nil {
		return nil, err
	}

	var moduleToWorkOn string

	if !module.isSet {
		if tg.GlobalTGModule == "" {
			return nil, fmt.Errorf("module directory cannot be empty, and it was not set in the constructor")
		}
	} else {
		moduleToWorkOn = module.value
	}

	workDirWithModule := path.Join(workDirDefault, moduleToWorkOn)

	tg.Ctr = tg.Ctr.
		WithWorkdir(workDirWithModule).
		WithMountedDirectory(workDirDefault, srcDirToMount)

	tg.Ctr = addTGCommandsToContainer(tg.Ctr, cmds)

	return tg.Ctr, nil
}

// TGInit initializes a Terragrunt module.
func (tg *IacTerragrunt) TGInit(src Optional[*Directory], module Optional[string], args Optional[[]string]) (*Container, error) {
	srcDirToMount, err := isEitherGlobalOrArgSetDir(src, tg.GlobalSRC)
	if err != nil {
		return nil, err
	}

	moduleToWorkOn, err := isEitherGlobalOrArgSetString(module, tg.GlobalTGModule)
	if err != nil {
		return nil, err
	}

	workDirWithModule := path.Join(workDirDefault, moduleToWorkOn)

	tg.Ctr = tg.Ctr.
		WithWorkdir(workDirWithModule).
		WithMountedDirectory(workDirDefault, srcDirToMount)

	var argsToSet []string

	if args.isSet {
		argsToSet = args.value
	}

	tg.Ctr = addTGCommandsToContainer(tg.Ctr, []string{"init"}, argsToSet...)

	return tg.Ctr, nil
}
