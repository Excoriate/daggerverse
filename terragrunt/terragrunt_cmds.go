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
	// configPath is the path to the terragrunt configuration file.
	// corresponds to the TERRAGRUNT_CONFIG environment variable.
	// +optional
	configPath string,
	// terraformPath is the path to the terraform binary.
	// corresponds to the TERRAGRUNT_TFPATH environment variable.
	// +optional
	terraformPath string,
	// workingDir is the working directory for terragrunt.
	// corresponds to the TERRAGRUNT_WORKING_DIR environment variable.
	// +optional
	workingDir string,
	// logLevel is the log level for terragrunt.
	// corresponds to the TERRAGRUNT_LOG_LEVEL environment variable.
	// +optional
	logLevel string,
	// iamRole is the iam role to assume before running terragrunt.
	// corresponds to the TERRAGRUNT_IAM_ROLE environment variable.
	// +optional
	iamRole string,
	// iamRoleSessionName is the iam role session name to use when assuming the iam role.
	// corresponds to the TERRAGRUNT_IAM_ROLE_SESSION_NAME environment variable.
	// +optional
	iamRoleSessionName string,
	// iamRoleDuration is the iam role duration to use when assuming the iam role.
	// corresponds to the TERRAGRUNT_IAM_ROLE_DURATION environment variable.
	// +optional
	iamRoleDuration string,
	// iamRoleExternalID is the iam role external id to use when assuming the iam role.
	// corresponds to the TERRAGRUNT_IAM_ROLE_EXTERNAL_ID environment variable.
	// +optional
	iamRoleExternalID string,
	// iamRolePolicy is the iam role policy to use when assuming the iam role.
	// corresponds to the TERRAGRUNT_IAM_ROLE_POLICY environment variable.
	// +optional
	iamRolePolicy string,
	// iamRolePolicyArns is the iam role policy arns to use when assuming the iam role.
	// corresponds to the TERRAGRUNT_IAM_ROLE_POLICY_ARNS environment variable.
	// +optional
	iamRolePolicyArns string,
	// iamRoleTags is the iam role tags to use when assuming the iam role.
	// corresponds to the TERRAGRUNT_IAM_ROLE_TAGS environment variable.
	// +optional
	iamRoleTags string,
	// iamRoleTransitiveTagKeys is the iam role transitive tag keys to use when assuming the iam role.
	// corresponds to the TERRAGRUNT_IAM_ROLE_TRANSITIVE_TAG_KEYS environment variable.
	// +optional
	iamRoleTransitiveTagKeys string,
	// iamRoleSourceIdentity is the iam role source identity to use when assuming the iam role.
	// corresponds to the TERRAGRUNT_IAM_ROLE_SOURCE_IDENTITY environment variable.
	// +optional
	iamRoleSourceIdentity string,
	// downloadDir is the directory to download terragrunt dependencies.
	// corresponds to the TERRAGRUNT_DOWNLOAD_DIR environment variable.
	// +optional
	downloadDir string,
	// source is the source url for terragrunt.
	// corresponds to the TERRAGRUNT_SOURCE environment variable.
	// +optional
	source string,
	// sourceMap is the source map for terragrunt.
	// corresponds to the TERRAGRUNT_SOURCE_MAP environment variable.
	// +optional
	sourceMap string,
	// sourceUpdate is the flag to update the source before running terragrunt.
	// corresponds to the TERRAGRUNT_SOURCE_UPDATE environment variable.
	// +optional
	sourceUpdate bool,
	// ignoreDependencyErrors is the flag to ignore dependency errors.
	// corresponds to the TERRAGRUNT_IGNORE_DEPENDENCY_ERRORS environment variable.
	// +optional
	ignoreDependencyErrors bool,
	// ignoreExternalDependencies is the flag to ignore external dependencies.
	// corresponds to the TERRAGRUNT_IGNORE_EXTERNAL_DEPENDENCIES environment variable.
	// +optional
	ignoreExternalDependencies bool,
	// includeExternalDependencies is the flag to include external dependencies.
	// corresponds to the TERRAGRUNT_INCLUDE_EXTERNAL_DEPENDENCIES environment variable.
	// +optional
	includeExternalDependencies bool,
	// parallelism is the parallelism level for terragrunt.
	// corresponds to the TERRAGRUNT_PARALLELISM environment variable.
	// +optional
	parallelism int,
	// debug is the flag to enable debug mode.
	// corresponds to the TERRAGRUNT_DEBUG environment variable.
	// +optional
	debug bool,
	// noColor is the flag to disable color in logs.
	// corresponds to the TERRAGRUNT_NO_COLOR environment variable.
	// +optional
	noColor bool,
	// check is the flag to check the configuration.
	// corresponds to the TERRAGRUNT_CHECK environment variable.
	// +optional
	check bool,
	// diff is the flag to enable diff mode.
	// corresponds to the TERRAGRUNT_DIFF environment variable.
	// +optional
	diff bool,
	// hclfmtFile is the file for hcl formatting.
	// corresponds to the TERRAGRUNT_HCLFMT_FILE environment variable.
	// +optional
	hclfmtFile string,
	// hclValidateJSON is the flag to validate hcl in json format.
	// corresponds to the TERRAGRUNT_HCLVALIDATE_JSON environment variable.
	// +optional
	hclValidateJSON bool,
	// hclValidateShowConfigPath is the flag to show the config path in hcl validation.
	// corresponds to the TERRAGRUNT_HCLVALIDATE_SHOW_CONFIG_PATH environment variable.
	// +optional
	hclValidateShowConfigPath bool,
	// overrideAttr is the attribute to override.
	// corresponds to the TERRAGRUNT_OVERRIDE_ATTR environment variable.
	// +optional
	overrideAttr string,
	// jsonOutDir is the directory for json output.
	// corresponds to the TERRAGRUNT_JSON_OUT_DIR environment variable.
	// +optional
	jsonOutDir string,
	// disableLogFormatting is the flag to disable log formatting.
	// corresponds to the TERRAGRUNT_DISABLE_LOG_FORMATTING environment variable.
	// +optional
	disableLogFormatting bool,
	// forwardTfStdout is the flag to forward terraform stdout.
	// corresponds to the TERRAGRUNT_FORWARD_TF_STDOUT environment variable.
	// +optional
	forwardTfStdout bool,
	// noAutoInit is the flag to disable auto init.
	// corresponds to the TERRAGRUNT_NO_AUTO_INIT environment variable.
	// +optional
	noAutoInit bool,
	// noAutoRetry is the flag to disable auto retry.
	// corresponds to the TERRAGRUNT_NO_AUTO_RETRY environment variable.
	// +optional
	noAutoRetry bool,
	// nonInteractive is the flag to disable interactive mode.
	// corresponds to the TERRAGRUNT_NON_INTERACTIVE environment variable.
	// +optional
	nonInteractive bool,
	// excludeDir is the flag to exclude directories.
	// corresponds to the TERRAGRUNT_EXCLUDE_DIR environment variable.
	// +optional
	excludeDir string,
	// includeDir is the flag to include directories.
	// corresponds to the TERRAGRUNT_INCLUDE_DIR environment variable.
	// +optional
	includeDir string,
	// strictInclude is the flag to enable strict include.
	// corresponds to the TERRAGRUNT_STRICT_INCLUDE environment variable.
	// +optional
	strictInclude bool,
	// strictValidate is the flag to enable strict validate.
	// corresponds to the TERRAGRUNT_STRICT_VALIDATE environment variable.
	// +optional
	strictValidate bool,
	// ignoreDependencyOrder is the flag to ignore dependency order.
	// corresponds to the TERRAGRUNT_IGNORE_DEPENDENCY_ORDER environment variable.
	// +optional
	ignoreDependencyOrder bool,
	// usePartialParseConfigCache is the flag to use partial parse config cache.
	// corresponds to the TERRAGRUNT_USE_PARTIAL_PARSE_CONFIG_CACHE environment variable.
	// +optional
	usePartialParseConfigCache bool,
	// failOnStateBucketCreation is the flag to fail on state bucket creation.
	// corresponds to the TERRAGRUNT_FAIL_ON_STATE_BUCKET_CREATION environment variable.
	// +optional
	failOnStateBucketCreation bool,
	// disableBucketUpdate is the flag to disable bucket update.
	// corresponds to the TERRAGRUNT_DISABLE_BUCKET_UPDATE environment variable.
	// +optional
	disableBucketUpdate bool,
	// disableCommandValidation is the flag to disable command validation.
	// corresponds to the TERRAGRUNT_DISABLE_COMMAND_VALIDATION environment variable.
	// +optional
	disableCommandValidation bool,
) *Terragrunt {
	tgOpts := newTerragruntOptionsDagger(
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

	m.Ctr = tgOpts.WithTerragruntOptionsSetInContainer(m.Ctr)

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
	// ctx is the context to use when executing the terragrunt command.
	ctx context.Context,
	// tfToken is the value of the terraform token to use when executing the terragrunt command.
	tfToken *dagger.Secret,
) (*Terragrunt, error) {
	tfTokenName, err := tfToken.Name(ctx)
	if err != nil {
		return nil, WrapError(err, "failed to get the name of the terraform token passed as a secret")
	}

	if err := isTerraformTokenNameValid(tfTokenName); err != nil {
		return nil, WrapError(err, "failed to validate terraform token")
	}

	tfTokenValueAsTxt, err := tfToken.Plaintext(ctx)
	if err != nil {
		return nil, WrapError(err, "failed to get the value of the terraform token passed as a secret")
	}

	if tfTokenValueAsTxt == "" {
		return nil, WrapError(nil, "terraform token value is empty")
	}

	// Set the parsed terraform token as an environment variable in the container.
	m.Ctr = m.
		Ctr.
		WithSecretVariable(tfTokenName, tfToken)

	return m, nil
}
