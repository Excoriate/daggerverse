package main

import "github.com/excoriate/daggerverse/daggerx/pkg/terragrunt"

// Run executes a command in the container.
func (tg *IacTerragrunt) Run(
	// Cmds are the commands to execute. E.g.: "ls -lth, pwd"
	cmds []string,
	// EntryPointOverride is the entry point to use. If it's not set, it will use the default 'sh -c'.
	entryPointOverride Optional[[]string],
	// Src is the source directory to mount in the container.
	src Optional[*Directory],
	// WithFocus is a flag to enable or disable the standard output per command to execute.
	withFocus Optional[bool],
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

	// Set the entry point
	entryPointToSet := entryPointOverride.GetOr(entryPointShell)

	if len(cmds) == 0 {
		return nil, &IacTerragruntCMDError{
			Message: "command cannot be empty",
		}
	}

	// Set the source directory
	enableCacheOptional := toDaggerOptional(false)
	tg.Ctr = tg.WithSource(tg.SRC, enableCacheOptional, module).Ctr
	// Creating the commands, and setting them.
	daggerCMDs := BuilderDaggerCMDs(cmds, entryPointToSet)

	// Expose or not the standard output per command to execute.
	withFocusOptional := toDaggerOptional(withFocus.GetOr(false))
	tg.Ctr = tg.WithCommands(daggerCMDs, withFocusOptional).Ctr

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
) (*Container, error) {
	withFocusSetInTrue := toDaggerOptional(true)
	entryPointTGOptional := toDaggerOptional(entryPointTerragrunt)

	// New validtor.
	_ = terragrunt.NewValidator()

	if module == "" {
		return nil, &IacTerragruntCMDError{
			Message: "module directory cannot be empty. Ensure that you're passing the module directory where the target terragrunt.hcl file is located.",
		}
	}

	workDirOptional := toDaggerOptional(module)
	ctr, runErr := tg.Run(cmds, entryPointTGOptional, src, withFocusSetInTrue, workDirOptional)
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
