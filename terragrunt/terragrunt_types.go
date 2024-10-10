package main

// TgConfigSetAsEnvVar represents a configuration set as an environment variable.
// EnvVarKey is the key of the environment variable.
// EnvVarValue is the value of the environment variable.
// LogOptionValue is an optional value for logging purposes.
type TgConfigSetAsEnvVar struct {
	EnvVarKey      string
	EnvVarValue    string
	LogOptionValue interface{}
}
