package main

import (
	"strconv"
	"strings"
)

type TerragruntOptsConfig struct {
	// Tg holds the Terragrunt configuration.
	// +private
	Tg *Terragrunt
	// TgOpts holds the Terragrunt options.
	// +private
	TgOpts []TgConfigSetAsEnvVar
}

func newTerragruntOptionsDagger(
	tg *Terragrunt,
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
) *TerragruntOptsConfig {
	var daggers []TgConfigSetAsEnvVar
	cleanValue := func(value string) string {
		return strings.TrimSpace(strings.ReplaceAll(value, "\\", "\\\\"))
	}

	// Helper function to add boolean flags
	addBoolFlag := func(key string, value bool) {
		daggers = append(daggers, TgConfigSetAsEnvVar{
			EnvVarKey:      key,
			EnvVarValue:    strconv.FormatBool(value),
			LogOptionValue: value,
		})
	}

	// Helper function to add string flags
	addStringFlag := func(key, value, defaultValue string) {
		if value == "" {
			value = defaultValue
		}
		cleanedValue := cleanValue(value)
		daggers = append(daggers, TgConfigSetAsEnvVar{
			EnvVarKey:      key,
			EnvVarValue:    cleanedValue,
			LogOptionValue: cleanedValue,
		})
	}

	// Add all flags
	addStringFlag("TERRAGRUNT_CONFIG", ConfigPath, "terragrunt.hcl")
	addStringFlag("TERRAGRUNT_TFPATH", TerraformPath, "terraform")
	addStringFlag("TERRAGRUNT_WORKING_DIR", WorkingDir, ".")
	addStringFlag("TERRAGRUNT_LOG_LEVEL", LogLevel, "info")
	addStringFlag("TERRAGRUNT_IAM_ROLE", IamRole, "")
	addStringFlag("TERRAGRUNT_IAM_ASSUME_ROLE_DURATION", IamRoleDuration, "")
	addStringFlag("TERRAGRUNT_IAM_ASSUME_ROLE_SESSION_NAME", IamRoleSessionName, "")
	addStringFlag("TERRAGRUNT_IAM_ROLE_EXTERNAL_ID", IamRoleExternalID, "")
	addStringFlag("TERRAGRUNT_IAM_ROLE_POLICY", IamRolePolicy, "")
	addStringFlag("TERRAGRUNT_IAM_ROLE_POLICY_ARNS", IamRolePolicyArns, "")
	addStringFlag("TERRAGRUNT_IAM_ROLE_TAGS", IamRoleTags, "")
	addStringFlag("TERRAGRUNT_IAM_ROLE_TRANSITIVE_TAG_KEYS", IamRoleTransitiveTagKeys, "")
	addStringFlag("TERRAGRUNT_IAM_ROLE_SOURCE_IDENTITY", IamRoleSourceIdentity, "")
	addStringFlag("TERRAGRUNT_DOWNLOAD_DIR", DownloadDir, ".terragrunt-cache")
	addStringFlag("TERRAGRUNT_SOURCE", Source, "")
	addStringFlag("TERRAGRUNT_SOURCE_MAP", SourceMap, "")

	addBoolFlag("TERRAGRUNT_SOURCE_UPDATE", SourceUpdate)
	addBoolFlag("TERRAGRUNT_IGNORE_DEPENDENCY_ERRORS", IgnoreDependencyErrors)
	addBoolFlag("TERRAGRUNT_IGNORE_EXTERNAL_DEPENDENCIES", IgnoreExternalDependencies)
	addBoolFlag("TERRAGRUNT_INCLUDE_EXTERNAL_DEPENDENCIES", IncludeExternalDependencies)
	addBoolFlag("TERRAGRUNT_DEBUG", Debug)
	addBoolFlag("TERRAGRUNT_NO_COLOR", NoColor)
	addBoolFlag("TERRAGRUNT_CHECK", Check)
	addBoolFlag("TERRAGRUNT_DIFF", Diff)
	addBoolFlag("TERRAGRUNT_HCLVALIDATE_JSON", HclvalidateJson)
	addBoolFlag("TERRAGRUNT_HCLVALIDATE_SHOW_CONFIG_PATH", HclvalidateShowConfigPath)
	addBoolFlag("TERRAGRUNT_DISABLE_LOG_FORMATTING", DisableLogFormatting)
	addBoolFlag("TERRAGRUNT_FORWARD_TF_STDOUT", ForwardTfStdout)

	// Add Parallelism with default
	if Parallelism == 0 {
		Parallelism = 10
	}
	daggers = append(daggers, TgConfigSetAsEnvVar{
		EnvVarKey:      "TERRAGRUNT_PARALLELISM",
		EnvVarValue:    strconv.Itoa(Parallelism),
		LogOptionValue: Parallelism,
	})

	// Add HclfmtFile
	addStringFlag("TERRAGRUNT_HCLFMT_FILE", HclfmtFile, "")

	addBoolFlag("TERRAGRUNT_NO_AUTO_INIT", NoAutoInit)
	addBoolFlag("TERRAGRUNT_NO_AUTO_RETRY", NoAutoRetry)
	addBoolFlag("TERRAGRUNT_NON_INTERACTIVE", NonInteractive)
	addStringFlag("TERRAGRUNT_EXCLUDE_DIR", ExcludeDir, "")
	addStringFlag("TERRAGRUNT_INCLUDE_DIR", IncludeDir, "")
	addBoolFlag("TERRAGRUNT_STRICT_INCLUDE", StrictInclude)
	addBoolFlag("TERRAGRUNT_STRICT_VALIDATE", StrictValidate)
	addBoolFlag("TERRAGRUNT_IGNORE_DEPENDENCY_ORDER", IgnoreDependencyOrder)
	addStringFlag("TERRAGRUNT_OVERRIDE_ATTR", OverrideAttr, "")
	addStringFlag("TERRAGRUNT_JSON_OUT", JsonOutDir, "")
	addBoolFlag("TERRAGRUNT_USE_PARTIAL_PARSE_CONFIG_CACHE", UsePartialParseConfigCache)
	addBoolFlag("TERRAGRUNT_FAIL_ON_STATE_BUCKET_CREATION", FailOnStateBucketCreation)
	addBoolFlag("TERRAGRUNT_DISABLE_BUCKET_UPDATE", DisableBucketUpdate)
	addBoolFlag("TERRAGRUNT_DISABLE_COMMAND_VALIDATION", DisableCommandValidation)

	return &TerragruntOptsConfig{TgOpts: daggers, Tg: tg}
}

func (c *TerragruntOptsConfig) WithTerragruntOptionsSetInContainer() *TerragruntOptsConfig {
	for _, envVar := range c.TgOpts {
		c.Tg.Ctr = c.
			Tg.
			Ctr.
			WithEnvVariable(envVar.EnvVarKey, envVar.EnvVarValue)
	}

	return c
}
