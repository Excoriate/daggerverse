package main

import "strings"

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
	// Tg holds the Terragrunt configuration.
	// +private
	Tg *Terragrunt
}

func newTfLogsConfigDagger(tg *Terragrunt, tfLog, tfLogCore, tfLogProvider, tfLogPath string) *LogsConfig {
	var daggers []TgConfigSetAsEnvVar

	cleanValue := func(value string) string {
		return strings.TrimSpace(strings.ReplaceAll(value, "\\", "\\\\"))
	}

	if tfLog == "" {
		tfLog = "INFO" // Default value
	}
	cleanedTfLog := cleanValue(tfLog)
	daggers = append(daggers, TgConfigSetAsEnvVar{
		EnvVarKey:      "TF_LOG",
		EnvVarValue:    cleanedTfLog,
		LogOptionValue: cleanedTfLog,
	})

	if tfLogCore == "" {
		tfLogCore = "INFO" // Default value
	}
	cleanedTfLogCore := cleanValue(tfLogCore)
	daggers = append(daggers, TgConfigSetAsEnvVar{
		EnvVarKey:      "TF_LOG_CORE",
		EnvVarValue:    cleanedTfLogCore,
		LogOptionValue: cleanedTfLogCore,
	})

	if tfLogProvider == "" {
		tfLogProvider = "INFO" // Default value
	}
	cleanedTfLogProvider := cleanValue(tfLogProvider)
	daggers = append(daggers, TgConfigSetAsEnvVar{
		EnvVarKey:      "TF_LOG_PROVIDER",
		EnvVarValue:    cleanedTfLogProvider,
		LogOptionValue: cleanedTfLogProvider,
	})

	if tfLogPath == "" {
		tfLogPath = "/var/log/terraform.log" // Default value
	}
	cleanedTfLogPath := cleanValue(tfLogPath)
	daggers = append(daggers, TgConfigSetAsEnvVar{
		EnvVarKey:      "TF_LOG_PATH",
		EnvVarValue:    cleanedTfLogPath,
		LogOptionValue: cleanedTfLogPath,
	})

	return &LogsConfig{TfLogs: daggers, Tg: tg}
}

func newTgLogsConfigDagger(tg *Terragrunt, tgLogLevel string, tgLogDisableColor bool, tgLogShowAbsPaths bool, tgDisableLogFormatting bool, tgForwardTfStdout bool) *LogsConfig {
	var daggers []TgConfigSetAsEnvVar

	cleanValue := func(value string) string {
		return strings.TrimSpace(strings.ReplaceAll(value, "\\", "\\\\"))
	}

	if tgLogLevel == "" {
		tgLogLevel = "info" // Default value
	}
	cleanedTgLogLevel := cleanValue(tgLogLevel)
	daggers = append(daggers, TgConfigSetAsEnvVar{
		EnvVarKey:      "TERRAGRUNT_LOG_LEVEL",
		EnvVarValue:    cleanedTgLogLevel,
		LogOptionValue: cleanedTgLogLevel,
	})

	if tgLogDisableColor {
		daggers = append(daggers, TgConfigSetAsEnvVar{
			EnvVarKey:      "TERRAGRUNT_LOG_DISABLE_COLOR",
			EnvVarValue:    "true",
			LogOptionValue: true,
		})
	} else {
		daggers = append(daggers, TgConfigSetAsEnvVar{
			EnvVarKey:      "TERRAGRUNT_LOG_DISABLE_COLOR",
			EnvVarValue:    "false",
			LogOptionValue: false,
		})
	}

	if tgLogShowAbsPaths {
		daggers = append(daggers, TgConfigSetAsEnvVar{
			EnvVarKey:      "TERRAGRUNT_LOG_SHOW_ABS_PATHS",
			EnvVarValue:    "true",
			LogOptionValue: true,
		})
	} else {
		daggers = append(daggers, TgConfigSetAsEnvVar{
			EnvVarKey:      "TERRAGRUNT_LOG_SHOW_ABS_PATHS",
			EnvVarValue:    "false",
			LogOptionValue: false,
		})
	}

	if tgDisableLogFormatting {
		daggers = append(daggers, TgConfigSetAsEnvVar{
			EnvVarKey:      "TERRAGRUNT_DISABLE_LOG_FORMATTING",
			EnvVarValue:    "true",
			LogOptionValue: true,
		})
	} else {
		daggers = append(daggers, TgConfigSetAsEnvVar{
			EnvVarKey:      "TERRAGRUNT_DISABLE_LOG_FORMATTING",
			EnvVarValue:    "false",
			LogOptionValue: false,
		})
	}

	if tgForwardTfStdout {
		daggers = append(daggers, TgConfigSetAsEnvVar{
			EnvVarKey:      "TERRAGRUNT_FORWARD_TF_STDOUT",
			EnvVarValue:    "true",
			LogOptionValue: true,
		})
	} else {
		daggers = append(daggers, TgConfigSetAsEnvVar{
			EnvVarKey:      "TERRAGRUNT_FORWARD_TF_STDOUT",
			EnvVarValue:    "false",
			LogOptionValue: false,
		})
	}

	return &LogsConfig{TgLogs: daggers, Tg: tg}
}

func (l *LogsConfig) WithTerraformLogsSetInContainer() *Terragrunt {
	for _, envVar := range l.TfLogs {
		l.Tg.Ctr = l.Tg.Ctr.WithEnvVariable(envVar.EnvVarKey, envVar.EnvVarValue)
	}

	return l.Tg
}

func (l *LogsConfig) WithTerragruntLogsSetInContainer() *Terragrunt {
	for _, envVar := range l.TgLogs {
		l.Tg.Ctr = l.Tg.Ctr.WithEnvVariable(envVar.EnvVarKey, envVar.EnvVarValue)
	}

	return l.Tg
}
