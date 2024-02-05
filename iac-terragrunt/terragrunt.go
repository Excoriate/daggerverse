package main

// Init initializes the terragrunt module. It returns only the Container
func (tg *IacTerragrunt) Init(
	// SRC is the source directory to mount in the container.
	src Optional[*Directory],
	// Module is the Terragrunt module to initialize.
	module string,
	// args Optional[string],
	args Optional[[]string], // Arguments for the "init" command.
	// enableCacheVolume is a flag to enable or disable the cache.
	enableCacheVolume Optional[bool],
	// EnvVars are the environment variables to set in the container.
	envVars Optional[[]string],
	// secretVars are the secret variables to set in the container.
	secretVars Optional[[]string],
	// invalidateCache is a flag to enable or disable the cache.
	invalidateCache Optional[bool],
	// enableGitSSH
	gitSSH Optional[string],
	// stdOut
	stdout Optional[bool],
) *Container {
	c, _ := tg.execTerragrunt("init", src, module, args, enableCacheVolume, envVars, secretVars, invalidateCache, gitSSH, stdout)
	return c
}

// InitE initializes the terragrunt module, and returns the Container plus an error if any.
func (tg *IacTerragrunt) InitE(
	// SRC is the source directory to mount in the container.
	src Optional[*Directory],
	// Module is the Terragrunt module to initialize.
	module string,
	// args Optional[string],
	args Optional[[]string], // Arguments for the "init" command.
	// enableCacheVolume is a flag to enable or disable the cache.
	enableCacheVolume Optional[bool],
	// EnvVars are the environment variables to set in the container.
	envVars Optional[[]string],
	// secretVars are the secret variables to set in the container.
	secretVars Optional[[]string],
	// invalidateCache is a flag to enable or disable the cache.
	invalidateCache Optional[bool],
	// enableGitSSH
	gitSSH Optional[string],
	// stdOut
	stdout Optional[bool],
) (*Container, error) {
	return tg.execTerragrunt("init", src, module, args, enableCacheVolume, envVars, secretVars, invalidateCache, gitSSH, stdout)
}

// Plan plans the terragrunt module. It returns only the Container
func (tg *IacTerragrunt) Plan(
	// SRC is the source directory to mount in the container.
	src Optional[*Directory],
	// Module is the Terragrunt module to initialize.
	module string,
	// args Optional[string],
	args Optional[[]string], // Arguments for the "init" command.
	// enableCacheVolume is a flag to enable or disable the cache.
	enableCacheVolume Optional[bool],
	// EnvVars are the environment variables to set in the container.
	envVars Optional[[]string],
	// secretVars are the secret variables to set in the container.
	secretVars Optional[[]string],
	// invalidateCache is a flag to enable or disable the cache.
	invalidateCache Optional[bool],
	// enableGitSSH
	gitSSH Optional[string],
	// stdOut
	stdout Optional[bool],
) *Container {
	c, _ := tg.execTerragrunt("plan", src, module, args, enableCacheVolume, envVars, secretVars, invalidateCache, gitSSH, stdout)
	return c
}

// PlanE plans the terragrunt module, and returns the Container plus an error if any.
func (tg *IacTerragrunt) PlanE(
	// SRC is the source directory to mount in the container.
	src Optional[*Directory],
	// Module is the Terragrunt module to initialize.
	module string,
	// args Optional[string],
	args Optional[[]string], // Arguments for the "init" command.
	// enableCacheVolume is a flag to enable or disable the cache.
	enableCacheVolume Optional[bool],
	// EnvVars are the environment variables to set in the container.
	envVars Optional[[]string],
	// secretVars are the secret variables to set in the container.
	secretVars Optional[[]string],
	// invalidateCache is a flag to enable or disable the cache.
	invalidateCache Optional[bool],
	// enableGitSSH
	gitSSH Optional[string],
	// stdOut
	stdout Optional[bool],
) (*Container, error) {
	return tg.execTerragrunt("plan", src, module, args, enableCacheVolume, envVars, secretVars, invalidateCache, gitSSH, stdout)
}

// Apply applies the terragrunt module. It returns only the Container
func (tg *IacTerragrunt) Apply(
	// SRC is the source directory to mount in the container.
	src Optional[*Directory],
	// Module is the Terragrunt module to initialize.
	module string,
	// args Optional[string],
	args Optional[[]string], // Arguments for the "init" command.
	// enableCacheVolume is a flag to enable or disable the cache.
	enableCacheVolume Optional[bool],
	// EnvVars are the environment variables to set in the container.
	envVars Optional[[]string],
	// secretVars are the secret variables to set in the container.
	secretVars Optional[[]string],
	// invalidateCache is a flag to enable or disable the cache.
	invalidateCache Optional[bool],
	// enableGitSSH
	gitSSH Optional[string],
	// stdOut
	stdout Optional[bool],
) *Container {
	c, _ := tg.execTerragrunt("apply", src, module, args, enableCacheVolume, envVars, secretVars, invalidateCache, gitSSH, stdout)
	return c
}

// ApplyE applies the terragrunt module, and returns the Container plus an error if any.
func (tg *IacTerragrunt) ApplyE(
	// SRC is the source directory to mount in the container.
	src Optional[*Directory],
	// Module is the Terragrunt module to initialize.
	module string,
	// args Optional[string],
	args Optional[[]string], // Arguments for the "init" command.
	// enableCacheVolume is a flag to enable or disable the cache.
	enableCacheVolume Optional[bool],
	// EnvVars are the environment variables to set in the container.
	envVars Optional[[]string],
	// secretVars are the secret variables to set in the container.
	secretVars Optional[[]string],
	// invalidateCache is a flag to enable or disable the cache.
	invalidateCache Optional[bool],
	// enableGitSSH
	gitSSH Optional[string],
	// stdOut
	stdout Optional[bool],
) (*Container, error) {
	return tg.execTerragrunt("apply", src, module, args, enableCacheVolume, envVars, secretVars, invalidateCache, gitSSH, stdout)
}

// Destroy destroys the terragrunt module. It returns only the Container
func (tg *IacTerragrunt) Destroy(
	// SRC is the source directory to mount in the container.
	src Optional[*Directory],
	// Module is the Terragrunt module to initialize.
	module string,
	// args Optional[string],
	args Optional[[]string], // Arguments for the "init" command.
	// enableCacheVolume is a flag to enable or disable the cache.
	enableCacheVolume Optional[bool],
	// EnvVars are the environment variables to set in the container.
	envVars Optional[[]string],
	// secretVars are the secret variables to set in the container.
	secretVars Optional[[]string],
	// invalidateCache is a flag to enable or disable the cache.
	invalidateCache Optional[bool],
	// enableGitSSH
	gitSSH Optional[string],
	// stdOut
	stdout Optional[bool],
) *Container {
	c, _ := tg.execTerragrunt("destroy", src, module, args, enableCacheVolume, envVars, secretVars, invalidateCache, gitSSH, stdout)
	return c
}

// DestroyE destroys the terragrunt module, and returns the Container plus an error if any.
func (tg *IacTerragrunt) DestroyE(
	// SRC is the source directory to mount in the container.
	src Optional[*Directory],
	// Module is the Terragrunt module to initialize.
	module string,
	// args Optional[string],
	args Optional[[]string], // Arguments for the "init" command.
	// enableCacheVolume is a flag to enable or disable the cache.
	enableCacheVolume Optional[bool],
	// EnvVars are the environment variables to set in the container.
	envVars Optional[[]string],
	// secretVars are the secret variables to set in the container.
	secretVars Optional[[]string],
	// invalidateCache is a flag to enable or disable the cache.
	invalidateCache Optional[bool],
	// enableGitSSH
	gitSSH Optional[string],
	// stdOut
	stdout Optional[bool],
) (*Container, error) {
	return tg.execTerragrunt("destroy", src, module, args, enableCacheVolume, envVars, secretVars, invalidateCache, gitSSH, stdout)
}

// Validate validates the terragrunt module. It returns only the Container
func (tg *IacTerragrunt) Validate(
	// SRC is the source directory to mount in the container.
	src Optional[*Directory],
	// Module is the Terragrunt module to initialize.
	module string,
	// args Optional[string],
	args Optional[[]string], // Arguments for the "init" command.
	// enableCacheVolume is a flag to enable or disable the cache.
	enableCacheVolume Optional[bool],
	// EnvVars are the environment variables to set in the container.
	envVars Optional[[]string],
	// secretVars are the secret variables to set in the container.
	secretVars Optional[[]string],
	// invalidateCache is a flag to enable or disable the cache.
	invalidateCache Optional[bool],
	// enableGitSSH
	gitSSH Optional[string],
	// stdOut
	stdout Optional[bool],
) *Container {
	c, _ := tg.execTerragrunt("validate", src, module, args, enableCacheVolume, envVars, secretVars, invalidateCache, gitSSH, stdout)
	return c
}

// ValidateE validates the terragrunt module, and returns the Container plus an error if any.
func (tg *IacTerragrunt) ValidateE(
	// SRC is the source directory to mount in the container.
	src Optional[*Directory],
	// Module is the Terragrunt module to initialize.
	module string,
	// args Optional[string],
	args Optional[[]string], // Arguments for the "init" command.
	// enableCacheVolume is a flag to enable or disable the cache.
	enableCacheVolume Optional[bool],
	// EnvVars are the environment variables to set in the container.
	envVars Optional[[]string],
	// secretVars are the secret variables to set in the container.
	secretVars Optional[[]string],
	// invalidateCache is a flag to enable or disable the cache.
	invalidateCache Optional[bool],
	// enableGitSSH
	gitSSH Optional[string],
	// stdOut
	stdout Optional[bool],
) (*Container, error) {
	return tg.execTerragrunt("validate", src, module, args, enableCacheVolume, envVars, secretVars, invalidateCache, gitSSH, stdout)
}

// HCLFmt formats the terragrunt module. It returns only the Container
func (tg *IacTerragrunt) HCLFmt(
	// SRC is the source directory to mount in the container.
	src Optional[*Directory],
	// Module is the Terragrunt module to initialize.
	module string,
	// args Optional[string],
	args Optional[[]string], // Arguments for the "init" command.
	// enableCacheVolume is a flag to enable or disable the cache.
	enableCacheVolume Optional[bool],
	// EnvVars are the environment variables to set in the container.
	envVars Optional[[]string],
	// secretVars are the secret variables to set in the container.
	secretVars Optional[[]string],
	// invalidateCache is a flag to enable or disable the cache.
	invalidateCache Optional[bool],
	// enableGitSSH
	gitSSH Optional[string],
	// stdOut
	stdout Optional[bool],
) *Container {
	c, _ := tg.execTerragrunt("hclfmt", src, module, args, enableCacheVolume, envVars, secretVars, invalidateCache, gitSSH, stdout)
	return c
}

// HCLFmtE formats the terragrunt module, and returns the Container plus an error if any.
func (tg *IacTerragrunt) HCLFmtE(
	// SRC is the source directory to mount in the container.
	src Optional[*Directory],
	// Module is the Terragrunt module to initialize.
	module string,
	// args Optional[string],
	args Optional[[]string], // Arguments for the "init" command.
	// enableCacheVolume is a flag to enable or disable the cache.
	enableCacheVolume Optional[bool],
	// EnvVars are the environment variables to set in the container.
	envVars Optional[[]string],
	// secretVars are the secret variables to set in the container.
	secretVars Optional[[]string],
	// invalidateCache is a flag to enable or disable the cache.
	invalidateCache Optional[bool],
	// enableGitSSH
	gitSSH Optional[string],
	// stdOut
	stdout Optional[bool],
) (*Container, error) {
	return tg.execTerragrunt("hclfmt", src, module, args, enableCacheVolume, envVars, secretVars, invalidateCache, gitSSH, stdout)
}
