package main

import "strings"

// isTerraformTokenNameValid checks if the provided token name is a valid Terraform token name.
// A valid Terraform token name must start with the prefix "TF_TOKEN_".
//
// Parameters:
// - tokenName: The name of the token to validate.
//
// Returns:
// - bool: True if the token name is valid, false otherwise.
func isTerraformTokenNameValid(tokenName string) error {
	tfTokenPrefix := "TF_TOKEN_"

	if !strings.HasPrefix(tokenName, tfTokenPrefix) {
		return WrapErrorf(nil, "invalid terraform token name: %s, should start with %s",
			tokenName, tfTokenPrefix)
	}

	return nil
}

// validateTerragruntCommands validates the provided Terragrunt command.
// The command must be one of the valid Terragrunt commands as per the Terragrunt documentation.
// Returns an error if the command is invalid or empty.
func validateTerragruntCommands(command string) error {
	if command == "" {
		return WrapError(nil, "command is required, can't validate empty command")
	}

	validCommands := map[string]bool{
		"run-all": true, "terragrunt-info": true, "validate-inputs": true, "graph-dependencies": true,
		"hclfmt": true, "hclvalidate": true, "aws-provider-patch": true, "render-json": true,
		"output-module-groups": true, "scaffold": true, "catalog": true, "graph": true,
		"version": true,
	}

	if !validCommands[command] {
		return WrapErrorf(nil, "this terragrunt command is not supported, check gruntwork's documentation: %s", command)
	}

	return nil
}

// validateTerraformLogLevel validates the provided terraform log level.
// The log level must be one of the valid
// log levels such as "TRACE", "DEBUG", "INFO", "WARN", "ERROR", or "JSON".
// Returns an error if the log level is invalid or empty.
//
//nolint:unused // The log level is not used in the function, but it is validated.
func validateTerraformLogLevel(logLevel string) error {
	if logLevel == "" {
		return WrapError(nil, "log level is required, can't validate empty log level")
	}

	validLogLevels := map[string]bool{
		"TRACE": true, "DEBUG": true, "INFO": true, "WARN": true, "ERROR": true, "JSON": true,
	}

	if !validLogLevels[logLevel] {
		return WrapErrorf(nil, "invalid terraform log level: %s", logLevel)
	}

	return nil
}

// validateTerragruntLogLevel validates the provided Terragrunt log level.
// The log level must be one of the valid log
// levels such as "stderr", "stdout", "error", "warn", "info", "debug", or "trace".
// Returns an error if the log level is invalid or empty.
//
//nolint:unused // The log level is not used in the function, but it is validated.
func validateTerragruntLogLevel(logLevel string) error {
	if logLevel == "" {
		return WrapError(nil, "log level is required, can't validate empty log level")
	}

	validLogLevels := map[string]bool{
		"stderr": true, "stdout": true, "error": true, "warn": true, "info": true, "debug": true, "trace": true,
	}

	if !validLogLevels[logLevel] {
		return WrapErrorf(nil, "invalid terragrunt log level: %s", logLevel)
	}

	return nil
}

// validateMainTerraformCommands validates the provided terraform main command.
// The command must be one of the valid main commands such as "init", "validate", "plan", "apply", or "destroy".
// Returns an error if the command is invalid or empty.
func validateMainTerraformCommands(command string) error {
	if command == "" {
		return WrapError(nil, "command is required, can't validate empty command")
	}

	validMainCommands := map[string]bool{
		"init": true, "validate": true, "plan": true, "apply": true, "destroy": true,
	}

	if !validMainCommands[command] {
		return WrapErrorf(nil, "invalid main terraform command: %s", command)
	}

	return nil
}

// validateTerraformOtherCommands validates the provided terraform other command.
// The command must be one of the valid other commands such
// as "console", "fmt", "force-unlock", "get", "graph", "import",
// "login", "logout", "metadata", "output", "providers", "refresh", "show", "state",
// "taint", "untaint", "version", or "workspace".
// Returns an error if the command is invalid or empty.
func validateTerraformOtherCommands(command string) error {
	if command == "" {
		return WrapError(nil, "command is required, can't validate empty command")
	}

	validOtherCommands := map[string]bool{
		"console": true, "fmt": true, "force-unlock": true, "get": true, "graph": true, "import": true,
		"login": true, "logout": true, "metadata": true, "output": true, "providers": true, "refresh": true,
		"show": true, "state": true, "taint": true, "untaint": true, "version": true, "workspace": true,
	}

	if !validOtherCommands[command] {
		return WrapErrorf(nil, "invalid terraform other command: %s", command)
	}

	return nil
}
