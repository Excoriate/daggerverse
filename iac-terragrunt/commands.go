package main

//import "github.com/excoriate/daggerverse/daggerx/pkg/terragrunt"

// Run executes a command in the container.
func (tg *IacTerragrunt) Run(
	// Cmds are the commands to execute. E.g.: "ls -lth, pwd"
	cmds []string,
	// Src is the source directory to mount in the container.
	src Optional[*Directory],
	// Stdout is a flag to enable or disable the standard output per command to execute.
	stdout Optional[bool],
	// Module is the working directory to use in the container.
	module Optional[string],
) (*Container, error) {
	if tg.SRC == nil && !src.isSet {
		return nil, &IacTerragruntCMDError{
			Message: "source directory cannot be empty, and it was not set in the constructor",
		}
	}

	if src.isSet {
		tg.SRC = src.value
	}

	if len(cmds) == 0 {
		return nil, &IacTerragruntCMDError{
			Message: "command cannot be empty",
		}
	}

	// Set the entry point to use shell instead of the default entry point.
	tg.Ctr = tg.WithEntrypoint(nil).Ctr

	// Set the source directory
	enableCacheOptional := toDaggerOptional(false)
	tg.Ctr = tg.WithSource(tg.SRC, enableCacheOptional, module).Ctr
	// Creating the commands, and setting them.
	daggerCMDs := buildShellCMDs(cmds)

	// Expose or not the standard output per command to execute.
	tg.Ctr = tg.WithCommands(daggerCMDs, stdout).Ctr

	return tg.Ctr, nil
}

// RunTG executes a terragrunt command
func (tg *IacTerragrunt) RunTG(
	// Cmds are the commands to execute. E.g.: "ls -lth, pwd"
	cmds []string,
	// Src is the source directory to mount in the container.
	src Optional[*Directory],
	// Module is the module directory to mount in the container.
	module string,
	// Stdout is a flag to enable or disable the standard output per command to execute.
	stdout Optional[bool],
) (*Container, error) {
	if module == "" {
		return nil, &IacTerragruntCMDError{
			Message: "module directory cannot be empty. Ensure that you're passing the module directory where the target terragrunt.hcl file is located.",
		}
	}

	workDirOptional := toDaggerOptional(module)

	// Set the entry point to use shell instead of the default entry point.
	tg.Ctr = tg.WithEntrypoint(nil).Ctr
	cmds = concatTerragruntInCommand(cmds)

	ctr, runErr := tg.Run(cmds, src, stdout, workDirOptional)

	if runErr != nil {
		return nil, &IacTerragruntCMDError{
			ErrWrapped: runErr,
			Message:    "failed to run terragrunt command",
		}
	}

	tg.Ctr = ctr
	tg.Ctr = tg.WithModule(module).Ctr

	return tg.Ctr, nil
}

// Init initializes the terragrunt module.
func (tg *IacTerragrunt) execTerragrunt(
	// TerragruntCMD is the terragrunt command to execute.
	terragruntCMD string,
	// SRC is the source directory to mount in the container.
	src Optional[*Directory],
	// Module is the Terragrunt module to initialize.
	module string,
	//args Optional[string],
	args Optional[[]string], // Arguments for the "init" command.
	// EnableCache is a flag to enable or disable the cache.
	enableCache Optional[bool],
) (*Container, error) {
	var cmd []string
	cmd = append(cmd, terragruntCMD)

	if args.isSet && len(args.value) > 0 {
		for _, arg := range args.value {
			cmd = append(cmd, arg)
		}
	}

	srcToUse := src.GetOr(tg.SRC)

	tg.Ctr = tg.WithSource(srcToUse, enableCache, toDaggerOptional(module)).Ctr
	tg.Ctr = tg.WithEntrypoint(entryPointTerragrunt).
		WithCommands(addCMDToDaggerCMD(cmd), toDaggerOptional(false)).Ctr

	return tg.Ctr, nil
}
