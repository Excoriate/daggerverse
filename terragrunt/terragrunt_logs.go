package main

import (
	"strconv"
	"strings"

	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
)

// TfLogsConfig holds the configuration for Terraform logs.
type TfLogsConfig struct {
	// TfLog is the log level for Terraform.
	TfLog string
	// TfLogCore is the log level for the core Terraform components.
	TfLogCore string
	// TfLogProvider is the log level for Terraform providers.
	TfLogProvider string
	// TfLogPath is the file path where Terraform logs will be written.
	TfLogPath string
}

// TgLogsConfig holds the configuration for Terragrunt logs.
type TgLogsConfig struct {
	// TgLogLevel is the log level for Terragrunt.
	TgLogLevel string
	// TgLogDisableColor indicates whether to disable color in Terragrunt logs.
	TgLogDisableColor bool
	// TgLogShowAbsPaths indicates whether to show absolute paths in Terragrunt logs.
	TgLogShowAbsPaths bool
}

// LogsConfig holds the configuration for both Terraform and Terragrunt logs.
type LogsConfig struct {
	// TfLogs holds the configuration for Terraform logs.
	TfLogs []TgConfigSetAsEnvVar
	// TgLogs holds the configuration for Terragrunt logs.
	TgLogs []TgConfigSetAsEnvVar
}

func newTfLogsConfigDagger(tfLog, tfLogCore, tfLogProvider, tfLogPath string) *LogsConfig {
	l := &LogsConfig{}
	var daggers []TgConfigSetAsEnvVar

	cleanValue := func(value string) string {
		return strings.TrimSpace(strings.ReplaceAll(value, "\\", "\\\\"))
	}

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

	addStringFlag("TF_LOG", tfLog, "INFO")
	addStringFlag("TF_LOG_CORE", tfLogCore, "INFO")
	addStringFlag("TF_LOG_PROVIDER", tfLogProvider, "INFO")
	addStringFlag("TF_LOG_PATH", tfLogPath, "/var/log/terraform.log")

	l.TfLogs = daggers

	return l
}

func newTgLogsConfigDagger(tgLogLevel string, tgLogDisableColor bool, tgLogShowAbsPaths bool, tgLogDisableFormatting bool, tgForwardTfStdout bool) *LogsConfig {
	l := &LogsConfig{}
	var daggers []TgConfigSetAsEnvVar

	cleanValue := func(value string) string {
		return strings.TrimSpace(strings.ReplaceAll(value, "\\", "\\\\"))
	}

	addBoolFlag := func(key string, value bool) {
		daggers = append(daggers, TgConfigSetAsEnvVar{
			EnvVarKey:      key,
			EnvVarValue:    strconv.FormatBool(value),
			LogOptionValue: value,
		})
	}

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

	addStringFlag("TERRAGRUNT_LOG_LEVEL", tgLogLevel, "info")
	addBoolFlag("TERRAGRUNT_LOG_DISABLE_COLOR", tgLogDisableColor)
	addBoolFlag("TERRAGRUNT_LOG_SHOW_ABS_PATHS", tgLogShowAbsPaths)
	addBoolFlag("TERRAGRUNT_LOG_DISABLE_FORMATTING", tgLogDisableFormatting)
	addBoolFlag("TERRAGRUNT_FORWARD_TF_STDOUT", tgForwardTfStdout)

	l.TgLogs = daggers
	return l
}

func (l *LogsConfig) WithTerraformLogsSetInContainer(ctr *dagger.Container) *dagger.Container {
	for _, envVar := range l.TfLogs {
		ctr = ctr.
			WithEnvVariable(envVar.EnvVarKey, envVar.EnvVarValue)
	}

	return ctr
}

func (l *LogsConfig) WithTerragruntLogsSetInContainer(ctr *dagger.Container) *dagger.Container {
	for _, envVar := range l.TgLogs {
		ctr = ctr.
			WithEnvVariable(envVar.EnvVarKey, envVar.EnvVarValue)
	}

	return ctr
}
