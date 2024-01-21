package main

import "github.com/excoriate/daggerverse/daggercommon/pkg/constants"

type IACTerraGrunt struct {
	Ctr *Container
}

func (tg *IACTerraGrunt) RunShell(cmd string) *Container {
	shellEntryPoint := constants.ShellEntryPoint
	cmdBuilt := append(shellEntryPoint, cmd)

	tg.Ctr = tg.Ctr.WithExec(cmdBuilt)
	return tg.Ctr
}
