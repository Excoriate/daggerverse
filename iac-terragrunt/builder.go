package main

import (
	"fmt"
	"github.com/excoriate/daggerverse/daggercommon/pkg/constants"
	"strings"
)

// IACTerraGruntOptions are options for IACTerraGrunt.
type IACTerraGruntOptions struct {
	Version string
}

type initIACTerraGruntOptFn func(*IACTerraGruntOptions) error

type IACTerraGruntBuilder struct {
	version       string
	shellCommands []string
}

func (b *IACTerraGruntBuilder) Build(optFns ...func(fn *IACTerraGruntOptions) error) (*IACTerraGrunt, error) {
	for _, option := range optFns {
		if err := option(&IACTerraGruntOptions{
			Version: b.version,
		}); err != nil {
			return nil, err
		}
	}

	var iacTerraGrunt IACTerraGrunt
	ctr := dag.Container().
		From(fmt.Sprintf("%s:%s", defaultContainerImage, b.version))

	iacTerraGrunt.Ctr = ctr

	if len(b.shellCommands) > 0 {
		for _, cmd := range b.shellCommands {
			cmdBuilt := append(constants.ShellEntryPoint, cmd)
			ctr = ctr.WithExec(cmdBuilt)
		}
	}

	return &iacTerraGrunt, nil
}

func (b *IACTerraGruntBuilder) WithVersion(version string) func(*IACTerraGruntOptions) error {
	return func(options *IACTerraGruntOptions) error {
		if version == "" {
			version = "latest"
		}

		if strings.Contains(version, "v") {
			return fmt.Errorf("version should not contain 'v'")
		}

		options.Version = version

		return nil
	}
}

func (b *IACTerraGruntBuilder) WithShellCommands(cmd []string) func(*IACTerraGruntOptions) error {
	return func(options *IACTerraGruntOptions) error {
		if cmd == nil {
			return fmt.Errorf("commands should not be nil")
		}

		for _, c := range cmd {
			if c == "" {
				return fmt.Errorf("command should not be empty")
			}
		}

		b.shellCommands = cmd
		return nil
	}
}
