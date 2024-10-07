package main

import (
	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/cmdx"
	"github.com/Excoriate/daggerx/pkg/envvars"
)

const (
	terragruntEntrypoint = "terragrunt"
)

// TerragruntCmd represents a command to be executed by Terragrunt.
type TerragruntCmd struct {
	// This field is private and should not be accessed directly.
	// +private
	Terragrunt *Terragrunt
	// Entrypoint is the entrypoint to use when executing the Terragrunt command.
	// This field is private and should not be accessed directly.
	Entrypoint string
	// Opts contains the options for the Terragrunt command.
	// +private
	Opts *TerragruntOptsConfig
	// Logs contains the configuration for logging.
	// +private
	Logs *LogsConfig
}

func (c *TerragruntCmd) getEntrypoint() string {
	return c.Entrypoint
}

// Validate checks if the provided command is a recognized terragrunt or terraform command.
// It returns an error if the command is invalid, otherwise it returns nil.
// It expects the command to be a valid terragrunt command, or a valid terraform command.
// It returns an error if the command is invalid or empty.
func (m *TerragruntCmd) Validate(command string) error {
	isNotTgCmdErr := validateTerragruntCommands(command)
	isNotTfCmdErr := validateMainTerraformCommands(command)
	isNotTfOtherCmdErr := validateTerraformOtherCommands(command)

	if isNotTgCmdErr != nil && isNotTfCmdErr != nil && isNotTfOtherCmdErr != nil {
		return WrapErrorf(nil, "invalid command: %s", command)
	}

	return nil
}

// TgExec executes a given terragrunt command within a dagger container.
// It returns a pointer to the resulting dagger.Container or an error if the command is invalid or fails to execute.
//
//nolint:lll // It's okay, since the ignore pattern is included.
func (c *TerragruntCmd) TgExec(
	// command is the terragrunt command to execute. It's the actual command that comes after 'terragrunt'
	command string,
	// source is the source directory that includes the source code.
	// +defaultPath="/"
	// +ignore=[".terragrunt-cache", ".terraform", ".github", ".gitignore", ".git", "vendor", "node_modules", "build", "dist", "log"]
	source *dagger.Directory,
	// module is the module to execute or the terragrunt configuration where the terragrunt.hcl file is located.
	// +optional
	module string,
	// tfLog is the terraform log mode to use when executing the terragrunt command.
	// +optional
	tfLog string,
	// tfLogCore is the terraform log core mode to use when executing the terragrunt command.
	// +optional
	tfLogCore string,
	// tfLogProvider is the terraform log provider mode to use when executing the terragrunt command.
	// +optional
	tfLogProvider string,
	// tgLog is the terragrunt log mode to use when executing the terragrunt command.
	// +optional
	tgLog string,
	// tfLogPath is the path to the terraform log file to use when executing the terragrunt command.
	// +optional
	tfLogPath string,
	// tfToken is the terraform token to use when executing the terragrunt command. It will form
	// an environment variable called TF_TOKEN_<token>
	// +optional
	tgLogLevel string,
	// tgLogDisableColor is the flag to disable color in terragrunt logs.
	// +optional
	tgLogDisableColor bool,
	// tgLogShowAbsPaths is the flag to show absolute paths in terragrunt logs.
	// +optional
	tgLogShowAbsPaths bool,
	// tgLogPath is the path to the terragrunt log file to use when executing the terragrunt command.
	// +optional
	// tgLogDisableFormatting is the flag to disable formatting in terragrunt logs.
	// +optional
	tgLogDisableFormatting bool,
	// tgForwardTfStdout is the flag to forward terraform stdout to terragrunt stdout.
	// +optional
	tgForwardTfStdout bool,
	// tfToken is the terraform token to use when executing the terragrunt command. It will form
	// an environment variable called TF_TOKEN_<token>
	// +optional
	tfToken string,
	// envVars are the environment variables to set in the container in the format of "key=value, key=value".
	// +optional
	envVars []string,
) (*dagger.Container, error) {
	if err := c.Validate(command); err != nil {
		return nil, WrapErrorf(err, "failed to validate command: %s", command)
	}

	if source == nil {
		return nil, WrapError(nil, "source is required, can't execute command without source")
	}

	// Generate the command as a slice of strings
	cmdAsSlice, cmdAsSliceErr := cmdx.GenerateDaggerCMDFromStr(command)
	if cmdAsSliceErr != nil {
		return nil, WrapErrorf(cmdAsSliceErr, "failed to generate dagger command from string: %s", command)
	}

	// Mount the source directory, and set 'terragrunt' as the owner of the directory
	c.Terragrunt.WithSource(source, module, terragruntCtrUser)

	// Set the environment variables
	if envVars != nil {
		envVarsAsDagger, envVarsErr := envvars.ToDaggerEnvVarsFromSlice(envVars)
		if envVarsErr != nil {
			return nil, WrapErrorf(envVarsErr, "failed to convert environment variables to dagger environment variables: %s", envVars)
		}

		for _, envVar := range envVarsAsDagger {
			c.Terragrunt.Ctr = c.Terragrunt.
				Ctr.
				WithEnvVariable(envVar.Name, envVar.Value)
		}
	}

	// Set the logs configuration
	newTfLogsConfigDagger(c.Terragrunt, tfLog, tfLogCore, tfLogProvider, tfLogPath)
	newTgLogsConfigDagger(c.Terragrunt, tgLogLevel, tgLogDisableColor, tgLogShowAbsPaths, tgLogDisableFormatting, tgForwardTfStdout)

	// Set the logs configuration in the actual module's container.
	c.Logs.WithTerraformLogsSetInContainer()
	c.Logs.WithTerragruntLogsSetInContainer()

	// Set the terraform token if it is provided, and passes the validation.
	if tfToken != "" {
		parsedTfToken, err := parseTerraformToken(tfToken)
		if err != nil {
			return nil, WrapErrorf(err, "failed to validate terraform token: %s", tfToken)
		}

		c.Terragrunt.Ctr = c.
			Terragrunt.
			Ctr.
			WithEnvVariable(parsedTfToken.EnvVarKey, parsedTfToken.EnvVarValue)
	}

	return c.Terragrunt.Ctr.
		WithExec(append([]string{c.getEntrypoint()}, cmdAsSlice...)), nil
}
