package main

import (
	"fmt"
	"main/constants"
)

type IacTerragrunt struct {
	Ctr *Container
}

func (tg *IacTerragrunt) Run(module string, cmds []string) *Container {
	shellEntryPoint := constants.ShellEntryPoint

	if module != "" {
		tg.Ctr = tg.Ctr.WithWorkdir(fmt.Sprintf("/src/%s", module))
	}

	if len(cmds) == 0 {
		return tg.Ctr
	}

	for _, cmd := range cmds {
		cmdBuilt := append(shellEntryPoint, cmd)
		tg.Ctr = tg.Ctr.WithExec(cmdBuilt)
	}

	return tg.Ctr
}

// New returns a new IacTerragrunt instance that can be used to run.
// Terragrunt commands, and optionally, shell commands.
func (tg *IacTerragrunt) New(version string, src *Directory) (*IacTerragrunt, error) {
	b := &IACTerraGruntBuilder{}
	var apiOpts []func(fn *IACTerraGruntOptions) error

	// Adding options.
	apiOpts = append(apiOpts, b.WithVersion(version))
	//apiOpts = append(apiOpts, b.WithModule(module))
	//apiOpts = append(apiOpts, b.WithSRC(src))
	apiOpts = append(apiOpts, b.WithHostDir(src))

	// Build the damn thing.
	tg, err := b.Build(apiOpts...)
	if err != nil {
		return nil, err
	}

	return tg, nil
}
