package main

import (
	"context"

	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/cmdx"
	"github.com/Excoriate/daggerx/pkg/envvars"
)

const (
	terragruntEntrypoint = "terragrunt"
)

// Cmd is an interface that represents a command to be executed by Terragrunt.
type Cmd interface {
	// Exec executes a given command within a dagger container.
	// It returns a pointer to the resulting dagger.Container or an error if the command is invalid or fails to execute.
	Exec(command string, source *dagger.Directory, module string, envVars []string) (*dagger.Container, error)
	// validate checks if the provided command is a recognized IaC tool command.
	// It returns an error if the command is invalid, otherwise it returns nil.
	// It expects the command to be a valid IaC tool command.
	// It returns an error if the command is invalid or empty.
	validate(command string) error
	// getEntrypoint returns the entrypoint to use when executing the command.
	getEntrypoint() string
}

// TerragruntCmd represents a command to be executed by Terragrunt.
type TerragruntCmd struct {
	// Tg is the Terragrunt module.
	// +private
	Tg *Terragrunt
	// Entrypoint is the entrypoint to use when executing the Terragrunt command.
	Entrypoint string
	// Opts contains the options for the Terragrunt command.
	// +private
	Opts *TerragruntOptsConfig
	// Logs contains the configuration for logging.
	// +private
	Logs *LogsConfig
}

// getEntrypoint returns the entrypoint to use when executing the command.
func (c *TerragruntCmd) getEntrypoint() string {
	return terragruntEntrypoint
}

// Validate checks if the provided command is a recognized terragrunt or terraform command.
// It returns an error if the command is invalid, otherwise it returns nil.
// It expects the command to be a valid terragrunt command, or a valid terraform command.
// It returns an error if the command is invalid or empty.
func (m *TerragruntCmd) validate(command string) error {
	isNotTgCmdErr := validateTerragruntCommands(command)
	isNotTfCmdErr := validateMainTerraformCommands(command)
	isNotTfOtherCmdErr := validateTerraformOtherCommands(command)

	if isNotTgCmdErr != nil && isNotTfCmdErr != nil && isNotTfOtherCmdErr != nil {
		return WrapErrorf(nil, "invalid command: %s", command)
	}

	return nil
}

// Exec executes a given command within a dagger container.
// It returns the output of the command or an error if the command is invalid or fails to execute.
//
//nolint:lll // It's okay, since the ignore pattern is included.
func (m *Terragrunt) Exec(
	ctx context.Context,
	// command is the terragrunt command to execute. It's the actual command that comes after 'terragrunt'
	command string,
	// source is the source directory that includes the source code.
	// +defaultPath="/"
	// +ignore=[".terragrunt-cache", ".terraform", ".github", ".gitignore", ".git", "vendor", "node_modules", "build", "dist", "log"]
	source *dagger.Directory,
	// module is the module to execute or the terragrunt configuration where the terragrunt.hcl file is located.
	// +optional
	module string,
	// envVars is the environment variables to pass to the container.
	// +optional
	envVars []string,
) (*dagger.Container, error) {
	if err := m.Tg.validate(command); err != nil {
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
	m = m.WithSource(source, module, terragruntCtrUser)

	// Set the environment variables
	if envVars != nil {
		envVarsAsDagger, envVarsErr := envvars.ToDaggerEnvVarsFromSlice(envVars)
		if envVarsErr != nil {
			return nil, WrapErrorf(envVarsErr, "failed to convert environment variables to dagger environment variables: %s", envVars)
		}

		for _, envVar := range envVarsAsDagger {
			m.Ctr = m.
				Ctr.
				WithEnvVariable(envVar.Name, envVar.Value)
		}
	}

	return m.Ctr.
		WithExec(append([]string{m.Tg.getEntrypoint()}, cmdAsSlice...)), nil
}

// ExecCmd executes a given command within a dagger container.
// It returns the output of the command or an error if the command is invalid or fails to execute.
func (m *Terragrunt) ExecCmd(
	ctx context.Context,
	// command is the terragrunt command to execute. It's the actual command that comes after 'terragrunt'
	command string,
	// source is the source directory that includes the source code.
	// +defaultPath="/"
	// +ignore=[".terragrunt-cache", ".terraform", ".github", ".gitignore", ".git", "vendor", "node_modules", "build", "dist", "log"]
	source *dagger.Directory,
	// module is the module to execute or the terragrunt configuration where the terragrunt.hcl file is located.
	// +optional
	module string,
	// envVars is the environment variables to pass to the container.
	// +optional
	envVars []string,
) (string, error) {
	container, err := m.Exec(ctx, command, source, module, envVars)

	if err != nil {
		return "", WrapErrorf(err, "failed to execute terragrunt command: %s", command)
	}

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", WrapErrorf(err, "failed to get stdout from terragrunt command: %s", command)
	}

	return output, nil
}

// WithTerraformLogOptions sets the terraform log options in the container.
func (m *Terragrunt) WithTerraformLogOptions(
	// tfLog is the terraform log mode to use when executing the terragrunt command.
	// +optional
	tfLog string,
	// tfLogCore is the terraform log core mode to use when executing the terragrunt command.
	// +optional
	tfLogCore string,
	// tfLogProvider is the terraform log provider mode to use when executing the terragrunt command.
	// +optional
	tfLogProvider string,
	// tfLogPath is the path to the terraform log file to use when executing the terragrunt command.
	// +optional
	tfLogPath string,
) *Terragrunt {
	l := newTfLogsConfigDagger(tfLog, tfLogCore, tfLogProvider, tfLogPath)
	m.Ctr = l.WithTerraformLogsSetInContainer(m.Ctr)
	return m
}

// WithTerragruntLogOptions sets the terragrunt log options in the container.
func (m *Terragrunt) WithTerragruntLogOptions(
	// tgLogLevel is the terragrunt log level to use when executing the terragrunt command.
	// +optional
	tgLogLevel string,
	// tgLogDisableColor is the flag to disable color in terragrunt logs.
	// +optional
	tgLogDisableColor bool,
	// tgLogShowAbsPaths is the flag to show absolute paths in terragrunt logs.
	// +optional
	tgLogShowAbsPaths bool,
	// tgLogDisableFormatting is the flag to disable formatting in terragrunt logs.
	// +optional
	tgLogDisableFormatting bool,
	// tgForwardTfStdout is the flag to forward terraform stdout to terragrunt stdout.
	// +optional
	tgForwardTfStdout bool,
) *Terragrunt {
	lCfg := newTgLogsConfigDagger(tgLogLevel, tgLogDisableColor, tgLogShowAbsPaths, tgLogDisableFormatting, tgForwardTfStdout)
	m.Ctr = lCfg.WithTerragruntLogsSetInContainer(m.Ctr)
	return m
}

// WithTerragruntOptions sets various Terragrunt options in the container.
// This function allows you to configure the Terragrunt environment by setting
// multiple parameters such as the configuration file path, Terraform binary path,
// working directory, log level, IAM role details, download directory, source URL,
// source map, and a flag to update the source before running Terragrunt.
// Each parameter corresponds to a specific Terragrunt environment variable.
func (m *Terragrunt) WithTerragruntOptions(
	// The path to the Terragrunt configuration file.
	// Corresponds to the TERRAGRUNT_CONFIG environment variable.
	ConfigPath string,
	// The path to the Terraform binary.
	// Corresponds to the TERRAGRUNT_TFPATH environment variable.
	TerraformPath string,
	// The working directory for Terragrunt.
	// Corresponds to the TERRAGRUNT_WORKING_DIR environment variable.
	WorkingDir string,
	// The log level for Terragrunt.
	// Corresponds to the TERRAGRUNT_LOG_LEVEL environment variable.
	LogLevel string,
	// The IAM role to assume before running Terragrunt.
	// Corresponds to the TERRAGRUNT_IAM_ROLE environment variable.
	IamRole string,
	// The IAM role session name to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_SESSION_NAME environment variable.
	IamRoleSessionName string,
	// The IAM role duration to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_DURATION environment variable.
	IamRoleDuration string,
	// The IAM role external ID to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_EXTERNAL_ID environment variable.
	IamRoleExternalID string,
	// The IAM role policy to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_POLICY environment variable.
	IamRolePolicy string,
	// The IAM role policy ARNs to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_POLICY_ARNS environment variable.
	IamRolePolicyArns string,
	// The IAM role tags to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_TAGS environment variable.
	IamRoleTags string,
	// The IAM role transitive tag keys to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_TRANSITIVE_TAG_KEYS environment variable.
	IamRoleTransitiveTagKeys string,
	// The IAM role source identity to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_SOURCE_IDENTITY environment variable.
	IamRoleSourceIdentity string,
	// The directory to download Terragrunt dependencies.
	// Corresponds to the TERRAGRUNT_DOWNLOAD_DIR environment variable.
	DownloadDir string,
	// The source URL for Terragrunt.
	// Corresponds to the TERRAGRUNT_SOURCE environment variable.
	Source string,
	// The source map for Terragrunt.
	// Corresponds to the TERRAGRUNT_SOURCE_MAP environment variable.
	SourceMap string,
	// The flag to update the source before running Terragrunt.
	// Corresponds to the TERRAGRUNT_SOURCE_UPDATE environment variable.
	SourceUpdate bool,
	// The flag to ignore dependency errors.
	// Corresponds to the TERRAGRUNT_IGNORE_DEPENDENCY_ERRORS environment variable.
	IgnoreDependencyErrors bool,
	// The flag to ignore external dependencies.
	// Corresponds to the TERRAGRUNT_IGNORE_EXTERNAL_DEPENDENCIES environment variable.
	IgnoreExternalDependencies bool,
	// The flag to include external dependencies.
	// Corresponds to the TERRAGRUNT_INCLUDE_EXTERNAL_DEPENDENCIES environment variable.
	IncludeExternalDependencies bool,
	// The parallelism level for Terragrunt.
	// Corresponds to the TERRAGRUNT_PARALLELISM environment variable.
	Parallelism int,
	// The flag to enable debug mode.
	// Corresponds to the TERRAGRUNT_DEBUG environment variable.
	Debug bool,
	// The flag to disable color in logs.
	// Corresponds to the TERRAGRUNT_NO_COLOR environment variable.
	NoColor bool,
	// The flag to check the configuration.
	// Corresponds to the TERRAGRUNT_CHECK environment variable.
	Check bool,
	// The flag to enable diff mode.
	// Corresponds to the TERRAGRUNT_DIFF environment variable.
	Diff bool,
	// The file for HCL formatting.
	// Corresponds to the TERRAGRUNT_HCLFMT_FILE environment variable.
	HclfmtFile string,
	// The flag to validate HCL in JSON format.
	// Corresponds to the TERRAGRUNT_HCLVALIDATE_JSON environment variable.
	HclvalidateJson bool,
	// The flag to show the config path in HCL validation.
	// Corresponds to the TERRAGRUNT_HCLVALIDATE_SHOW_CONFIG_PATH environment variable.
	HclvalidateShowConfigPath bool,
	// The attribute to override.
	// Corresponds to the TERRAGRUNT_OVERRIDE_ATTR environment variable.
	OverrideAttr string,
	// The directory for JSON output.
	// Corresponds to the TERRAGRUNT_JSON_OUT_DIR environment variable.
	JsonOutDir string,
	// The flag to disable log formatting.
	// Corresponds to the TERRAGRUNT_DISABLE_LOG_FORMATTING environment variable.
	DisableLogFormatting bool,
	// The flag to forward Terraform stdout.
	// Corresponds to the TERRAGRUNT_FORWARD_TF_STDOUT environment variable.
	ForwardTfStdout bool,
	// The flag to disable auto init.
	// Corresponds to the TERRAGRUNT_NO_AUTO_INIT environment variable.
	NoAutoInit bool,
	// The flag to disable auto retry.
	// Corresponds to the TERRAGRUNT_NO_AUTO_RETRY environment variable.
	NoAutoRetry bool,
	// The flag to disable interactive mode.
	// Corresponds to the TERRAGRUNT_NON_INTERACTIVE environment variable.
	NonInteractive bool,
	// The flag to exclude directories.
	// Corresponds to the TERRAGRUNT_EXCLUDE_DIR environment variable.
	ExcludeDir string,
	// The flag to include directories.
	// Corresponds to the TERRAGRUNT_INCLUDE_DIR environment variable.
	IncludeDir string,
	// The flag to enable strict include.
	// Corresponds to the TERRAGRUNT_STRICT_INCLUDE environment variable.
	StrictInclude bool,
	// The flag to enable strict validate.
	// Corresponds to the TERRAGRUNT_STRICT_VALIDATE environment variable.
	StrictValidate bool,
	// The flag to ignore dependency order.
	// Corresponds to the TERRAGRUNT_IGNORE_DEPENDENCY_ORDER environment variable.
	IgnoreDependencyOrder bool,
	// The flag to use partial parse config cache.
	// Corresponds to the TERRAGRUNT_USE_PARTIAL_PARSE_CONFIG_CACHE environment variable.
	UsePartialParseConfigCache bool,
	// The flag to fail on state bucket creation.
	// Corresponds to the TERRAGRUNT_FAIL_ON_STATE_BUCKET_CREATION environment variable.
	FailOnStateBucketCreation bool,
	// The flag to disable bucket update.
	// Corresponds to the TERRAGRUNT_DISABLE_BUCKET_UPDATE environment variable.
	DisableBucketUpdate bool,
	// The flag to disable command validation.
	// Corresponds to the TERRAGRUNT_DISABLE_COMMAND_VALIDATION environment variable.
	DisableCommandValidation bool,
) *Terragrunt {
	m.Tg.Opts = newTerragruntOptionsDagger(
		m,
		ConfigPath,
		TerraformPath,
		WorkingDir,
		LogLevel,
		IamRole,
		IamRoleSessionName,
		IamRoleDuration,
		IamRoleExternalID,
		IamRolePolicy,
		IamRolePolicyArns,
		IamRoleTags,
		IamRoleTransitiveTagKeys,
		IamRoleSourceIdentity,
		DownloadDir,
		Source,
		SourceMap,
		SourceUpdate,
		IgnoreDependencyErrors,
		IgnoreExternalDependencies,
		IncludeExternalDependencies,
		Parallelism,
		Debug,
		NoColor,
		Check,
		Diff,
		HclfmtFile,
		HclvalidateJson,
		HclvalidateShowConfigPath,
		OverrideAttr,
		JsonOutDir,
		DisableLogFormatting,
		ForwardTfStdout,
		NoAutoInit,
		NoAutoRetry,
		NonInteractive,
		ExcludeDir,
		IncludeDir,
		StrictInclude,
		StrictValidate,
		IgnoreDependencyOrder,
		UsePartialParseConfigCache,
		FailOnStateBucketCreation,
		DisableBucketUpdate,
		DisableCommandValidation,
	)

	m.Tg.Opts.WithTerragruntOptionsSetInContainer()
	return m
}

// WithTerraformToken sets the terraform token in the container.
//
// This method takes a terraform token as input, validates it, and sets it as an environment variable
// in the container. The terraform token is used when executing the terragrunt command.
//
// Parameters:
// - tfToken: A string representing the terraform token to use.
//
// Returns:
// - *Terragrunt: A pointer to the updated Terragrunt instance.
// - error: An error if the terraform token validation fails.
func (m *Terragrunt) WithTerraformToken(
	// tfToken is the terraform token to use when executing the terragrunt command.
	tfToken string,
) (*Terragrunt, error) {
	// Parse and validate the terraform token.
	parsedTfToken, err := parseTerraformToken(tfToken)
	if err != nil {
		return nil, WrapErrorf(err, "failed to validate terraform token: %s", tfToken)
	}

	// Set the parsed terraform token as an environment variable in the container.
	m.Ctr = m.Ctr.WithEnvVariable(parsedTfToken.EnvVarKey, parsedTfToken.EnvVarValue)
	return m, nil
}
