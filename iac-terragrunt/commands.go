package main

import (
	"context"
	"fmt"
)

// Exec executes a command in the container.
func (tg *IacTerragrunt) Exec(cmds []string, entryPointOverride Optional[[]string], src Optional[*Directory]) (*Container, error) {
	if tg.SRC == nil {
		return nil, &IacTerragruntCMDError{
			Message: "source directory cannot be empty, and it was not set in the constructor",
		}
	}

	entryPointToSet := entryPointOverride.GetOr(entryPointShell)

	if len(cmds) == 0 {
		return nil, &IacTerragruntCMDError{
			Message: "command cannot be empty",
		}
	}

	srcDirToMount := src.GetOr(tg.SRC)

	tg.Ctr = tg.Ctr.
		WithWorkdir(workDirDefault).
		WithMountedDirectory(workDirDefault, srcDirToMount)

	var cmdsToExec [][]string

	for _, cmd := range cmds {
		cmdOptions := buildDaggerCMDsOptions{
			entryPoint: entryPointToSet,
			cmds:       []string{cmd},
		}

		cmdsToExec = append(cmdsToExec, buildDaggerCMDs([]buildDaggerCMDsOptions{cmdOptions})...)
	}

	for _, cmd := range cmdsToExec {
		tg.Ctr = tg.Ctr.WithExec(cmd)
	}

	return tg.Ctr, nil
}

// Run executes a command in the container and prints the output.
func (tg *IacTerragrunt) Run(ctx context.Context, cmds []string, src Optional[*Directory]) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	var emptyEntryPoint Optional[[]string]

	ctr, err := tg.Exec(cmds, emptyEntryPoint, src)
	if err != nil {
		return "", &IacTerragruntCMDError{
			ErrWrapped: err,
			Message:    "failed to execute command and print the output",
		}
	}

	output, outErr := ctr.Stdout(ctx)
	fmt.Printf("output: %s\n", output)
	if outErr != nil {
		return "", &IacTerragruntCMDError{
			ErrWrapped: outErr,
			Message:    "failed to print the output",
		}
	}

	return output, nil
}
