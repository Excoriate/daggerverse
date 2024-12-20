package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
)

// TerragruntOptsConfig holds the configuration and options for Terragrunt.
type TerragruntOptsConfig struct {
	// TgOpts holds the Terragrunt options.
	// +private
	TgOpts []TgConfigSetAsEnvVar
}

// newTerragruntOptionsDagger creates a new TerragruntOptsConfig with the provided parameters.
// It sets various Terragrunt options based on the provided arguments and environment variables.
//
//nolint:funlen // It's okay, it's not complex, just long due to Dagger's limitations.
func newTerragruntOptionsDagger(
	// The path to the Terragrunt configuration file.
	// Corresponds to the TERRAGRUNT_CONFIG environment variable.
	configPath string,
	// The path to the Terraform binary.
	// Corresponds to the TERRAGRUNT_TFPATH environment variable.
	terraformPath string,
	// The working directory for Terragrunt.
	// Corresponds to the TERRAGRUNT_WORKING_DIR environment variable.
	workingDir string,
	// The log level for Terragrunt.
	// Corresponds to the TERRAGRUNT_LOG_LEVEL environment variable.
	logLevel string,
	// The IAM role to assume before running Terragrunt.
	// Corresponds to the TERRAGRUNT_IAM_ROLE environment variable.
	iamRole string,
	// The IAM role session name to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_SESSION_NAME environment variable.
	iamRoleSessionName string,
	// The IAM role duration to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_DURATION environment variable.
	iamRoleDuration int,
	// The IAM role external ID to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_EXTERNAL_ID environment variable.
	iamRoleExternalID string,
	// The IAM role policy to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_POLICY environment variable.
	iamRolePolicy string,
	// The IAM role policy ARNs to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_POLICY_ARNS environment variable.
	iamRolePolicyArns string,
	// The IAM role tags to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_TAGS environment variable.
	iamRoleTags string,
	// The IAM role transitive tag keys to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_TRANSITIVE_TAG_KEYS environment variable.
	iamRoleTransitiveTagKeys string,
	// The IAM role source identity to use when assuming the IAM role.
	// Corresponds to the TERRAGRUNT_IAM_ROLE_SOURCE_IDENTITY environment variable.
	iamRoleSourceIdentity string,
	// The directory to download Terragrunt dependencies.
	// Corresponds to the TERRAGRUNT_DOWNLOAD_DIR environment variable.
	downloadDir string,
	// The source URL for Terragrunt.
	// Corresponds to the TERRAGRUNT_SOURCE environment variable.
	source string,
	// The source map for Terragrunt.
	// Corresponds to the TERRAGRUNT_SOURCE_MAP environment variable.
	sourceMap string,
	// The flag to update the source before running Terragrunt.
	// Corresponds to the TERRAGRUNT_SOURCE_UPDATE environment variable.
	sourceUpdate bool,
	// The flag to ignore dependency errors.
	// Corresponds to the TERRAGRUNT_IGNORE_DEPENDENCY_ERRORS environment variable.
	ignoreDependencyErrors bool,
	// The flag to ignore external dependencies.
	// Corresponds to the TERRAGRUNT_IGNORE_EXTERNAL_DEPENDENCIES environment variable.
	ignoreExternalDependencies bool,
	// The flag to include external dependencies.
	// Corresponds to the TERRAGRUNT_INCLUDE_EXTERNAL_DEPENDENCIES environment variable.
	includeExternalDependencies bool,
	// The parallelism level for Terragrunt.
	// Corresponds to the TERRAGRUNT_PARALLELISM environment variable.
	parallelism int,
	// The flag to enable debug mode.
	// Corresponds to the TERRAGRUNT_DEBUG environment variable.
	debug bool,
	// The flag to disable color in logs.
	// Corresponds to the TERRAGRUNT_NO_COLOR environment variable.
	noColor bool,
	// The flag to check the configuration.
	// Corresponds to the TERRAGRUNT_CHECK environment variable.
	check bool,
	// The flag to enable diff mode.
	// Corresponds to the TERRAGRUNT_DIFF environment variable.
	diff bool,
	// The file for HCL formatting.
	// Corresponds to the TERRAGRUNT_HCLFMT_FILE environment variable.
	hclfmtFile string,
	// hclValidateJSON is the flag to validate HCL in JSON format.
	// Corresponds to the TERRAGRUNT_HCLVALIDATE_JSON environment variable.
	hclValidateJSON bool,
	// hclValidateShowConfigPath is the flag to show the config path in HCL validation.
	// Corresponds to the TERRAGRUNT_HCLVALIDATE_SHOW_CONFIG_PATH environment variable.
	hclValidateShowConfigPath bool,
	// overrideAttr is the attribute to override.
	// Corresponds to the TERRAGRUNT_OVERRIDE_ATTR environment variable.
	overrideAttr string,
	// The directory for JSON output.
	// Corresponds to the TERRAGRUNT_JSON_OUT_DIR environment variable.
	jsonOutDir string,
	// The flag to disable log formatting.
	// Corresponds to the TERRAGRUNT_DISABLE_LOG_FORMATTING environment variable.
	disableLogFormatting bool,
	// The flag to forward Terraform stdout.
	// Corresponds to the TERRAGRUNT_FORWARD_TF_STDOUT environment variable.
	forwardTfStdout bool,
	// The flag to disable auto init.
	// Corresponds to the TERRAGRUNT_NO_AUTO_INIT environment variable.
	noAutoInit bool,
	// The flag to disable auto retry.
	// Corresponds to the TERRAGRUNT_NO_AUTO_RETRY environment variable.
	noAutoRetry bool,
	// The flag to disable interactive mode.
	// Corresponds to the TERRAGRUNT_NON_INTERACTIVE environment variable.
	nonInteractive bool,
	// The flag to exclude directories.
	// Corresponds to the TERRAGRUNT_EXCLUDE_DIR environment variable.
	excludeDir string,
	// The flag to include directories.
	// Corresponds to the TERRAGRUNT_INCLUDE_DIR environment variable.
	includeDir string,
	// The flag to enable strict include.
	// Corresponds to the TERRAGRUNT_STRICT_INCLUDE environment variable.
	strictInclude bool,
	// The flag to enable strict validate.
	// Corresponds to the TERRAGRUNT_STRICT_VALIDATE environment variable.
	strictValidate bool,
	// The flag to ignore dependency order.
	// Corresponds to the TERRAGRUNT_IGNORE_DEPENDENCY_ORDER environment variable.
	ignoreDependencyOrder bool,
	// The flag to use partial parse config cache.
	// Corresponds to the TERRAGRUNT_USE_PARTIAL_PARSE_CONFIG_CACHE environment variable.
	usePartialParseConfigCache bool,
	// The flag to fail on state bucket creation.
	// Corresponds to the TERRAGRUNT_FAIL_ON_STATE_BUCKET_CREATION environment variable.
	failOnStateBucketCreation bool,
	// The flag to disable bucket update.
	// Corresponds to the TERRAGRUNT_DISABLE_BUCKET_UPDATE environment variable.
	disableBucketUpdate bool,
	// The flag to disable command validation.
	// Corresponds to the TERRAGRUNT_DISABLE_COMMAND_VALIDATION environment variable.
	disableCommandValidation bool,
	// The duration for IAM role assumption.
	// Corresponds to the TERRAGRUNT_IAM_ASSUME_ROLE_DURATION environment variable.
	iamAssumeRoleDuration int,
) *TerragruntOptsConfig {
	var daggers []TgConfigSetAsEnvVar

	cleanValue := func(value string) string {
		return strings.TrimSpace(strings.ReplaceAll(value, "\\", "\\\\"))
	}

	addStringFlag := func(key, value, defaultValue string) {
		if value != "" && value != defaultValue {
			daggers = append(daggers, TgConfigSetAsEnvVar{
				EnvVarKey:   key,
				EnvVarValue: cleanValue(value),
			})
		}
	}

	addBoolFlag := func(key string, value bool) {
		if value {
			daggers = append(daggers, TgConfigSetAsEnvVar{
				EnvVarKey:   key,
				EnvVarValue: "true",
			})
		}
	}

	addIntFlag := func(key string, value, defaultValue int) {
		if value != defaultValue {
			daggers = append(daggers, TgConfigSetAsEnvVar{
				EnvVarKey:   key,
				EnvVarValue: strconv.Itoa(value),
			})
		}
	}

	addMapFlag := func(key, value string) {
		if value != "" {
			// Assume the value is a comma-separated list of key=value pairs
			pairs := strings.Split(value, ",")
			formattedPairs := make([]string, 0, len(pairs))

			for _, pair := range pairs {
				kv := strings.SplitN(pair, "=", 2)
				if len(kv) == 2 {
					formattedPairs = append(formattedPairs, fmt.Sprintf("%s=%s", cleanValue(kv[0]), cleanValue(kv[1])))
				}
			}

			if len(formattedPairs) > 0 {
				daggers = append(daggers, TgConfigSetAsEnvVar{
					EnvVarKey:   key,
					EnvVarValue: strings.Join(formattedPairs, ","),
				})
			}
		}
	}

	// Add all flags
	addStringFlag("TERRAGRUNT_CONFIG", configPath, "terragrunt.hcl")
	addStringFlag("TERRAGRUNT_TFPATH", terraformPath, "terraform")
	addStringFlag("TERRAGRUNT_WORKING_DIR", workingDir, ".")
	addStringFlag("TERRAGRUNT_LOG_LEVEL", logLevel, "info")
	addStringFlag("TERRAGRUNT_IAM_ROLE", iamRole, "")
	addIntFlag("TERRAGRUNT_IAM_ASSUME_ROLE_DURATION", iamRoleDuration, 3600) // Default to 1 hour
	addStringFlag("TERRAGRUNT_IAM_ASSUME_ROLE_SESSION_NAME", iamRoleSessionName, "terragrunt-session")
	addStringFlag("TERRAGRUNT_IAM_ROLE_EXTERNAL_ID", iamRoleExternalID, "")
	addStringFlag("TERRAGRUNT_IAM_ROLE_POLICY", iamRolePolicy, "")
	addStringFlag("TERRAGRUNT_IAM_ROLE_POLICY_ARNS", iamRolePolicyArns, "")
	addStringFlag("TERRAGRUNT_IAM_ROLE_TAGS", iamRoleTags, "")
	addStringFlag("TERRAGRUNT_IAM_ROLE_TRANSITIVE_TAG_KEYS", iamRoleTransitiveTagKeys, "")
	addStringFlag("TERRAGRUNT_IAM_ROLE_SOURCE_IDENTITY", iamRoleSourceIdentity, "")
	addStringFlag("TERRAGRUNT_DOWNLOAD_DIR", downloadDir, ".terragrunt-cache")
	addStringFlag("TERRAGRUNT_SOURCE", source, "")

	addBoolFlag("TERRAGRUNT_SOURCE_UPDATE", sourceUpdate)
	addBoolFlag("TERRAGRUNT_IGNORE_DEPENDENCY_ERRORS", ignoreDependencyErrors)
	addBoolFlag("TERRAGRUNT_IGNORE_EXTERNAL_DEPENDENCIES", ignoreExternalDependencies)
	addBoolFlag("TERRAGRUNT_INCLUDE_EXTERNAL_DEPENDENCIES", includeExternalDependencies)
	addBoolFlag("TERRAGRUNT_DEBUG", debug)
	addBoolFlag("TERRAGRUNT_NO_COLOR", noColor)
	addBoolFlag("TERRAGRUNT_CHECK", check)
	addBoolFlag("TERRAGRUNT_DIFF", diff)
	addBoolFlag("TERRAGRUNT_HCLVALIDATE_JSON", hclValidateJSON)
	addBoolFlag("TERRAGRUNT_HCLVALIDATE_SHOW_CONFIG_PATH", hclValidateShowConfigPath)
	addBoolFlag("TERRAGRUNT_DISABLE_LOG_FORMATTING", disableLogFormatting)
	addBoolFlag("TERRAGRUNT_FORWARD_TF_STDOUT", forwardTfStdout)

	addIntFlag("TERRAGRUNT_PARALLELISM", parallelism, 0)
	addIntFlag("TERRAGRUNT_IAM_ASSUME_ROLE_DURATION", iamAssumeRoleDuration, 0)

	addStringFlag("TERRAGRUNT_HCLFMT_FILE", hclfmtFile, "")
	addBoolFlag("TERRAGRUNT_NO_AUTO_INIT", noAutoInit)
	addBoolFlag("TERRAGRUNT_NO_AUTO_RETRY", noAutoRetry)
	addBoolFlag("TERRAGRUNT_NON_INTERACTIVE", nonInteractive)
	addStringFlag("TERRAGRUNT_EXCLUDE_DIR", excludeDir, "")
	addStringFlag("TERRAGRUNT_INCLUDE_DIR", includeDir, "")
	addBoolFlag("TERRAGRUNT_STRICT_INCLUDE", strictInclude)
	addBoolFlag("TERRAGRUNT_STRICT_VALIDATE", strictValidate)
	addBoolFlag("TERRAGRUNT_IGNORE_DEPENDENCY_ORDER", ignoreDependencyOrder)
	addStringFlag("TERRAGRUNT_OVERRIDE_ATTR", overrideAttr, "")
	addStringFlag("TERRAGRUNT_JSON_OUT", jsonOutDir, "")
	addBoolFlag("TERRAGRUNT_USE_PARTIAL_PARSE_CONFIG_CACHE", usePartialParseConfigCache)
	addBoolFlag("TERRAGRUNT_FAIL_ON_STATE_BUCKET_CREATION", failOnStateBucketCreation)
	addBoolFlag("TERRAGRUNT_DISABLE_BUCKET_UPDATE", disableBucketUpdate)
	addBoolFlag("TERRAGRUNT_DISABLE_COMMAND_VALIDATION", disableCommandValidation)

	// Handle TERRAGRUNT_SOURCE_MAP as a map-like string
	addMapFlag("TERRAGRUNT_SOURCE_MAP", sourceMap)

	return &TerragruntOptsConfig{TgOpts: daggers}
}

// WithTerragruntOptionsSetInContainer sets the environment variables for the Terragrunt container.
//
// This method iterates over the TgOpts slice in the TerragruntOptsConfig, which contains the environment
// variable key-value pairs, and sets each environment variable in the Terragrunt container using the
// WithEnvVariable method.
//
// Parameters:
// - ctr: A pointer to the dagger.Container in which the environment variables will be set.
//
// Returns:
// - A pointer to the modified dagger.Container with the environment variables set.
//
// Example usage:
//
//	config := &TerragruntOptsConfig{...}
//	updatedConfig := config.WithTerragruntOptionsSetInContainer(ctr)
func (c *TerragruntOptsConfig) WithTerragruntOptionsSetInContainer(ctr *dagger.Container) *dagger.Container {
	for _, envVar := range c.TgOpts {
		ctr = ctr.
			WithEnvVariable(envVar.EnvVarKey, envVar.EnvVarValue)
	}

	return ctr
}
