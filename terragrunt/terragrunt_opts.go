package main

// TerragruntOptions holds the options for the Terragrunt CLI.
// Each field corresponds to a Terragrunt environment variable.
type TerragruntOptions struct {
	// The path to the Terragrunt configuration file.
	// Corresponds to the TERRAGRUNT_CONFIG environment variable.
	ConfigPath string

	// The path to the Terraform binary.
	// Corresponds to the TERRAGRUNT_TFPATH environment variable.
	TerraformPath string

	// The working directory for Terragrunt.
	// Corresponds to the TERRAGRUNT_WORKING_DIR environment variable.
	WorkingDir string

	// The log level for Terragrunt.
	// Corresponds to the TERRAGRUNT_LOG_LEVEL environment variable.
	LogLevel string

	// The IAM role to assume before running Terragrunt.
	// Corresponds to the TERRAGRUNT_IAM_ROLE environment variable.
	IamRole string

	// The IAM role session name to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_SESSION_NAME environment variable.
	IamRoleSessionName string

	// The IAM role duration to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_DURATION environment variable.
	IamRoleDuration string

	// The IAM role external ID to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_EXTERNAL_ID environment variable.
	IamRoleExternalID string

	// The IAM role policy to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_POLICY environment variable.
	IamRolePolicy string

	// The IAM role policy ARNs to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_POLICY_ARNS environment variable.
	IamRolePolicyArns string

	// The IAM role tags to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_TAGS environment variable.
	IamRoleTags string

	// The IAM role transitive tag keys to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_TRANSITIVE_TAG_KEYS environment variable.
	IamRoleTransitiveTagKeys string

	// The IAM role source identity to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_SOURCE_IDENTITY environment variable.
	IamRoleSourceIdentity string

	// The directory to download Terragrunt dependencies.
	// Corresponds to the TERRAGRUNT_DOWNLOAD_DIR environment variable.
	DownloadDir string

	// The source URL for Terragrunt.
	// Corresponds to the TERRAGRUNT_SOURCE environment variable.
	Source string

	// The source map for Terragrunt.
	// Corresponds to the TERRAGRUNT_SOURCE_MAP environment variable.
	SourceMap string

	// The flag to update the source before running Terragrunt.
	// Corresponds to the TERRAGRUNT_SOURCE_UPDATE environment variable.
	SourceUpdate bool

	// The flag to ignore dependency errors.
	// Corresponds to the TERRAGRUNT_IGNORE_DEPENDENCY_ERRORS environment variable.
	IgnoreDependencyErrors bool

	// The flag to ignore external dependencies.
	// Corresponds to the TERRAGRUNT_IGNORE_EXTERNAL_DEPENDENCIES environment variable.
	IgnoreExternalDependencies bool

	// The flag to include external dependencies.
	// Corresponds to the TERRAGRUNT_INCLUDE_EXTERNAL_DEPENDENCIES environment variable.
	IncludeExternalDependencies bool

	// The parallelism level for Terragrunt.
	// Corresponds to the TERRAGRUNT_PARALLELISM environment variable.
	Parallelism int

	// The flag to enable debug mode.
	// Corresponds to the TERRAGRUNT_DEBUG environment variable.
	Debug bool

	// The flag to disable color in logs.
	// Corresponds to the TERRAGRUNT_NO_COLOR environment variable.
	NoColor bool

	// The flag to check the configuration.
	// Corresponds to the TERRAGRUNT_CHECK environment variable.
	Check bool

	// The flag to enable diff mode.
	// Corresponds to the TERRAGRUNT_DIFF environment variable.
	Diff bool

	// The file for HCL formatting.
	// Corresponds to the TERRAGRUNT_HCLFMT_FILE environment variable.
	HclfmtFile string

	// The flag to validate HCL in JSON format.
	// Corresponds to the TERRAGRUNT_HCLVALIDATE_JSON environment variable.
	HclvalidateJson bool

	// The flag to show the config path in HCL validation.
	// Corresponds to the TERRAGRUNT_HCLVALIDATE_SHOW_CONFIG_PATH environment variable.
	HclvalidateShowConfigPath bool

	// The attribute to override.
	// Corresponds to the TERRAGRUNT_OVERRIDE_ATTR environment variable.
	OverrideAttr string

	// The directory for JSON output.
	// Corresponds to the TERRAGRUNT_JSON_OUT_DIR environment variable.
	JsonOutDir string

	// The flag to disable log formatting.
	// Corresponds to the TERRAGRUNT_DISABLE_LOG_FORMATTING environment variable.
	DisableLogFormatting bool

	// The flag to forward Terraform stdout.
	// Corresponds to the TERRAGRUNT_FORWARD_TF_STDOUT environment variable.
	ForwardTfStdout bool
}
