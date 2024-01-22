package main

import (
	"fmt"
	"main/constants"
	"main/utils"
	"os"
)

// IACTerraGruntOptions are options for IacTerragrunt.
type IACTerraGruntOptions struct {
	version       string
	shellCommands []string
	src           string
	module        string
}

type initIACTerraGruntOptFn func(*IACTerraGruntOptions) error

type IACTerraGruntBuilder struct {
	version       string
	shellCommands []string
	src           string
	module        string
	srcDir        *Directory
}

func (b *IACTerraGruntBuilder) Build(optFns ...func(fn *IACTerraGruntOptions) error) (*IacTerragrunt, error) {
	for _, option := range optFns {
		if err := option(&IACTerraGruntOptions{
			version:       b.version,
			src:           b.src,
			module:        b.module,
			shellCommands: b.shellCommands,
		}); err != nil {
			return nil, err
		}
	}

	mountPath := "/src"
	var iacTerraGrunt IacTerragrunt
	var workDir string

	if b.module == "" {
		workDir = "/src"
	} else {
		workDir = fmt.Sprintf("/%s/%s", mountPath, b.module)
	}

	fmt.Printf("source directory: %s\n", b.src)

	imageAddr := fmt.Sprintf("%s:%s", defaultContainerImage, b.version)
	fmt.Printf("imageAddr: %s\n", imageAddr)

	ctr := dag.Container().
		From(imageAddr).
		WithDirectory(mountPath, b.srcDir, ContainerWithDirectoryOpts{
			Exclude: []string{".git", ".terragrunt-cache/**", ".terraform/**", ".terraform.lock.hcl"},
		}).
		WithWorkdir(workDir)

	if len(b.shellCommands) > 0 {
		for _, cmd := range b.shellCommands {
			cmdBuilt := append(constants.ShellEntryPoint, cmd)
			ctr = ctr.WithExec(cmdBuilt)
		}
	}

	iacTerraGrunt.Ctr = ctr

	return &iacTerraGrunt, nil
}

func (b *IACTerraGruntBuilder) WithVersion(version string) func(*IACTerraGruntOptions) error {
	return func(options *IACTerraGruntOptions) error {
		if version == "" {
			version = "latest"
		}

		b.version = version

		return nil
	}
}

func (b *IACTerraGruntBuilder) WithModule(module string) func(*IACTerraGruntOptions) error {
	return func(options *IACTerraGruntOptions) error {
		// It'll accept the same directory or src as the module's path if it's not set.
		if module == "" {
			return nil
		}

		b.module = module
		return nil
	}
}

func (b *IACTerraGruntBuilder) WithHostDir(dir *Directory) func(*IACTerraGruntOptions) error {
	return func(options *IACTerraGruntOptions) error {
		if dir == nil {
			b.srcDir = dag.Host().Directory(".")
		}

		b.srcDir = dir
		return nil
	}
}

func (b *IACTerraGruntBuilder) WithSRC(src string) func(*IACTerraGruntOptions) error {
	return func(options *IACTerraGruntOptions) error {
		if src == "" {
			options.src, _ = os.Getwd()
		}

		fmt.Printf("src: %s\n", src)

		absSRC, err := utils.ConvertToAbs(src)
		fmt.Print("absSRC: ", absSRC, "\n")
		if err != nil {
			return err
		}

		if err := utils.DirExist(absSRC); err != nil {
			return err
		}

		if err := utils.DirHasFilesWithExtension(absSRC, ".hcl"); err != nil {
			return err
		}

		// FIXME: Validate if the path actually exists.
		b.src = absSRC

		return nil
	}
}

func (b *IACTerraGruntBuilder) WithShellCommands(cmd []string) func(*IACTerraGruntOptions) error {
	return func(options *IACTerraGruntOptions) error {
		if cmd == nil {
			b.shellCommands = []string{}
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
