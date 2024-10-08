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
func (c *TerragruntCmd) validate(command string) error {
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
//
//nolint:lll // It's okay, since the ignore pattern is included.
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
	container, err := m.Exec(command, source, module, envVars)

	if err != nil {
		return "", WrapErrorf(err, "failed to execute terragrunt command: %s", command)
	}

	output, err := container.
		Stdout(ctx)

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
	lCfg := newTgLogsConfigDagger(
		tgLogLevel,
		tgLogDisableColor,
		tgLogShowAbsPaths,
		tgLogDisableFormatting,
		tgForwardTfStdout,
	)
	m.Ctr = lCfg.WithTerragruntLogsSetInContainer(m.Ctr)

	return m
}

// WithTerragruntOptions sets various Terragrunt options in the container.
// This function allows you to configure the Terragrunt environment by setting
// multiple parameters such as the configuration file path, Terraform binary path,
// working directory, log level, IAM role details, download directory, source URL,
// source map, and a flag to update the source before running Terragrunt.
// Each parameter corresponds to a specific Terragrunt environment variable.
//
//nolint:funlen // It's okay, it's not complex, just long due to Dagger's limitations.
func (m *Terragrunt) WithTerragruntOptions(
	// the path to the terragrunt configuration file.
	// corresponds to the terragrunt_config environment variable.
	configPath string,
	// the path to the terraform binary.
	// corresponds to the terragrunt_tfpath environment variable.
	terraformPath string,
	// the working directory for terragrunt.
	// corresponds to the terragrunt_working_dir environment variable.
	workingDir string,
	// the log level for terragrunt.
	// corresponds to the terragrunt_log_level environment variable.
	logLevel string,
	// the iam role to assume before running terragrunt.
	// corresponds to the terragrunt_iam_role environment variable.
	iamRole string,
	// the iam role session name to use when assuming the iam role.
	// corresponds to the terragrunt_iam_role_session_name environment variable.
	iamRoleSessionName string,
	// the iam role duration to use when assuming the iam role.
	// corresponds to the terragrunt_iam_role_duration environment variable.
	iamRoleDuration string,
	// the iam role external id to use when assuming the iam role.
	// corresponds to the terragrunt_iam_role_external_id environment variable.
	iamRoleExternalID string,
	// the iam role policy to use when assuming the iam role.
	// corresponds to the terragrunt_iam_role_policy environment variable.
	iamRolePolicy string,
	// the iam role policy arns to use when assuming the iam role.
	// corresponds to the terragrunt_iam_role_policy_arns environment variable.
	iamRolePolicyArns string,
	// the iam role tags to use when assuming the iam role.
	// corresponds to the terragrunt_iam_role_tags environment variable.
	iamRoleTags string,
	// the iam role transitive tag keys to use when assuming the iam role.
	// corresponds to the terragrunt_iam_role_transitive_tag_keys environment variable.
	iamRoleTransitiveTagKeys string,
	// the iam role source identity to use when assuming the iam role.
	// corresponds to the terragrunt_iam_role_source_identity environment variable.
	iamRoleSourceIdentity string,
	// the directory to download terragrunt dependencies.
	// corresponds to the terragrunt_download_dir environment variable.
	downloadDir string,
	// the source url for terragrunt.
	// corresponds to the terragrunt_source environment variable.
	source string,
	// the source map for terragrunt.
	// corresponds to the terragrunt_source_map environment variable.
	sourceMap string,
	// the flag to update the source before running terragrunt.
	// corresponds to the terragrunt_source_update environment variable.
	sourceUpdate bool,
	// the flag to ignore dependency errors.
	// corresponds to the terragrunt_ignore_dependency_errors environment variable.
	ignoreDependencyErrors bool,
	// the flag to ignore external dependencies.
	// corresponds to the terragrunt_ignore_external_dependencies environment variable.
	ignoreExternalDependencies bool,
	// the flag to include external dependencies.
	// corresponds to the terragrunt_include_external_dependencies environment variable.
	includeExternalDependencies bool,
	// the parallelism level for terragrunt.
	// corresponds to the terragrunt_parallelism environment variable.
	parallelism int,
	// the flag to enable debug mode.
	// corresponds to the terragrunt_debug environment variable.
	debug bool,
	// the flag to disable color in logs.
	// corresponds to the terragrunt_no_color environment variable.
	noColor bool,
	// the flag to check the configuration.
	// corresponds to the terragrunt_check environment variable.
	check bool,
	// the flag to enable diff mode.
	// corresponds to the terragrunt_diff environment variable.
	diff bool,
	// the file for hcl formatting.
	// corresponds to the terragrunt_hclfmt_file environment variable.
	hclfmtFile string,
	// hclValidateJSON is the flag to validate hcl in json format.
	// corresponds to the terragrunt_hclvalidate_json environment variable.
	hclValidateJSON bool,
	// hclValidateShowConfigPath is the flag to show the config path in hcl validation.
	// corresponds to the terragrunt_hclvalidate_show_config_path environment variable.
	hclValidateShowConfigPath bool,
	// overrideAttr is the attribute to override.
	// corresponds to the terragrunt_override_attr environment variable.
	overrideAttr string,
	// jsonOutDir is the directory for json output.
	// corresponds to the terragrunt_json_out_dir environment variable.
	jsonOutDir string,
	// the flag to disable log formatting.
	// corresponds to the terragrunt_disable_log_formatting environment variable.
	disableLogFormatting bool,
	// the flag to forward terraform stdout.
	// corresponds to the terragrunt_forward_tf_stdout environment variable.
	forwardTfStdout bool,
	// the flag to disable auto init.
	// corresponds to the terragrunt_no_auto_init environment variable.
	noAutoInit bool,
	// the flag to disable auto retry.
	// corresponds to the terragrunt_no_auto_retry environment variable.
	noAutoRetry bool,
	// the flag to disable interactive mode.
	// corresponds to the terragrunt_non_interactive environment variable.
	nonInteractive bool,
	// the flag to exclude directories.
	// corresponds to the terragrunt_exclude_dir environment variable.
	excludeDir string,
	// the flag to include directories.
	// corresponds to the terragrunt_include_dir environment variable.
	includeDir string,
	// the flag to enable strict include.
	// corresponds to the terragrunt_strict_include environment variable.
	strictInclude bool,
	// the flag to enable strict validate.
	// corresponds to the terragrunt_strict_validate environment variable.
	strictValidate bool,
	// the flag to ignore dependency order.
	// corresponds to the terragrunt_ignore_dependency_order environment variable.
	ignoreDependencyOrder bool,
	// the flag to use partial parse config cache.
	// corresponds to the terragrunt_use_partial_parse_config_cache environment variable.
	usePartialParseConfigCache bool,
	// the flag to fail on state bucket creation.
	// corresponds to the terragrunt_fail_on_state_bucket_creation environment variable.
	failOnStateBucketCreation bool,
	// the flag to disable bucket update.
	// corresponds to the terragrunt_disable_bucket_update environment variable.
	disableBucketUpdate bool,
	// the flag to disable command validation.
	// corresponds to the terragrunt_disable_command_validation environment variable.
	disableCommandValidation bool,
) *Terragrunt {
	m.Tg.Opts = newTerragruntOptionsDagger(
		m,
		configPath,
		terraformPath,
		workingDir,
		logLevel,
		iamRole,
		iamRoleSessionName,
		iamRoleDuration,
		iamRoleExternalID,
		iamRolePolicy,
		iamRolePolicyArns,
		iamRoleTags,
		iamRoleTransitiveTagKeys,
		iamRoleSourceIdentity,
		downloadDir,
		source,
		sourceMap,
		sourceUpdate,
		ignoreDependencyErrors,
		ignoreExternalDependencies,
		includeExternalDependencies,
		parallelism,
		debug,
		noColor,
		check,
		diff,
		hclfmtFile,
		hclValidateJSON,
		hclValidateShowConfigPath,
		overrideAttr,
		jsonOutDir,
		disableLogFormatting,
		forwardTfStdout,
		noAutoInit,
		noAutoRetry,
		nonInteractive,
		excludeDir,
		includeDir,
		strictInclude,
		strictValidate,
		ignoreDependencyOrder,
		usePartialParseConfigCache,
		failOnStateBucketCreation,
		disableBucketUpdate,
		disableCommandValidation,
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
