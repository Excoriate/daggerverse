package main

import "strings"

// ParseArgsFromStrToSlice parses a string of arguments in the form of "arg1, arg2, arg3"
func parseArgsFromStrToSlice(argStr string) []string {
	if argStr == "" {
		return []string{}
	}

	var parsedArgs []string
	// Split the string on the comma as a preliminary step.
	args := strings.Split(argStr, ",")
	for _, arg := range args {
		// Trim leading and trailing whitespace from each argument.
		arg = strings.TrimSpace(arg)
		// Handle special case for '-var' arguments.
		if strings.HasPrefix(arg, "-var ") {
			// Keep '-var' and its value together as a single string.
			parsedArgs = append(parsedArgs, arg)
		} else {
			// For other arguments, split on spaces assuming they are separate.
			parts := strings.Fields(arg)
			parsedArgs = append(parsedArgs, parts...)
		}
	}
	return parsedArgs
}

// parseEnvVarsInStringMapAsMap parses a string of environment variables in the form of "key1=value1,key2=value2"
func (t *Terratest) parseEnvVarsInStringMapAsMap(envVarsStr string) map[string]string {
	envVars := make(map[string]string)
	if envVarsStr == "" {
		return envVars
	}

	envVarsList := strings.Split(envVarsStr, ",")
	for _, envVar := range envVarsList {
		envVar = strings.TrimSpace(envVar)
		// Split on the first equals sign to separate the key and value.
		parts := strings.SplitN(envVar, "=", 2)
		envVars[parts[0]] = parts[1]
	}
	return envVars
}

// setEnvVarsInContainer sets the environment variables for the Terraform container.
func (t *Terratest) setEnvVarsInContainer(envVars map[string]string) *Container {
	for key, value := range envVars {
		t.Ctr = t.Ctr.WithEnvVariable(key, value)
	}
	return t.Ctr
}
