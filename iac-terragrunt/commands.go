package main

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
	// EnvVars are the environment variables to set in the container.
	envVars Optional[[]string],
	// secretVars are the secret variables to set in the container.
	secretVars Optional[[]string],
	// invalidateCache is a flag to enable or disable the cache.
	invalidateCache Optional[bool],
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

	// Convert slices to map, and inject it as environment variables.
	if envVars.isSet {
		envVarsMap, err := convertSliceToMap(envVars.value)
		if err != nil {
			return nil, &IacTerragruntCMDError{
				ErrWrapped: err,
				Message:    "failed to convert slice to map",
			}
		}

		for key, value := range envVarsMap {
			tg.Ctr = tg.WithEnvVar(key, value, false).Ctr
		}
	}

	if secretVars.isSet {
		secretVarsMap, err := convertSliceToMap(secretVars.value)
		if err != nil {
			return nil, &IacTerragruntCMDError{
				ErrWrapped: err,
				Message:    "failed to convert slice to map",
			}
		}

		for key, value := range secretVarsMap {
			tg.Ctr = tg.WithSecret(key, value).Ctr
		}
	}

	// Set the source directory
	enableCacheOptional := toDaggerOptional(false)
	tg.Ctr = tg.WithSource(tg.SRC, enableCacheOptional, module).Ctr
	// Creating the commands, and setting them.
	daggerCMDs := buildShellCMDs(cmds)

	// Expose or not the standard output per command to execute.
	stdoutValue := stdout.GetOr(false)
	tg.Ctr = tg.WithCommands(daggerCMDs, stdoutValue).Ctr

	// Invalidate the cache if the flag is set.
	if invalidateCache.isSet && invalidateCache.value {
		tg.Ctr = tg.WithCacheInvalidation().Ctr
	}

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
	// EnvVars are the environment variables to set in the container.
	envVars Optional[[]string],
	// secretVars are the secret variables to set in the container.
	secretVars Optional[[]string],
	// invalidateCache is a flag to enable or disable the cache.
	invalidateCache Optional[bool],
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

	ctr, runErr := tg.Run(cmds, src, stdout, workDirOptional, envVars, secretVars, invalidateCache)

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
	// args Optional[string],
	args Optional[[]string], // Arguments for the "init" command.
	// EnableCacheVolume is a flag to enable or disable the cache.
	enableCacheVolume Optional[bool],
	// EnvVars are the environment variables to set in the container.
	envVars Optional[[]string],
	// secretVars are the secret variables to set in the container.
	secretVars Optional[[]string],
	// invalidateCache is a flag to enable or disable the cache.
	invalidateCache Optional[bool],
	// enableGitSSH
	gitSSH Optional[string],
	//stdOut
	stdout Optional[bool],
) (*Container, error) {
	var cmd []string
	cmd = append(cmd, terragruntCMD)

	if args.isSet && len(args.value) > 0 {
		cmd = append(cmd, args.value...)
	}

	srcToUse := src.GetOr(tg.SRC)

	if envVars.isSet {
		envVarsMap, err := convertSliceToMap(envVars.value)
		if err != nil {
			return nil, &IacTerragruntCMDError{
				ErrWrapped: err,
				Message:    "failed to convert slice to map",
			}
		}

		for key, value := range envVarsMap {
			tg.Ctr = tg.WithEnvVar(key, value, false).Ctr
		}
	}

	if secretVars.isSet {
		secretVarsMap, err := convertSliceToMap(secretVars.value)
		if err != nil {
			return nil, &IacTerragruntCMDError{
				ErrWrapped: err,
				Message:    "failed to convert slice to map",
			}
		}

		for key, value := range secretVarsMap {
			tg.Ctr = tg.WithSecret(key, value).Ctr
		}
	}

	stdoutValue := stdout.GetOr(false)

	tg.Ctr = tg.WithSource(srcToUse, enableCacheVolume, toDaggerOptional(module)).Ctr
	tg.Ctr = tg.WithEntrypoint(entryPointTerragrunt).
		WithCommands(addCMDToDaggerCMD(cmd), stdoutValue).Ctr

	if invalidateCache.isSet && invalidateCache.value {
		tg.Ctr = tg.WithCacheInvalidation().Ctr
	}

	if gitSSH.isSet {
		tg.Ctr = tg.WithGitSSHConfig(gitSSH.value).Ctr
	}

	return tg.Ctr, nil
}
