package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/Excoriate/daggerx/pkg/cmdbuilder"
)

// Terminal returns a terminal for the container.
//
// It returns a terminal for the container. It's meant to be used as a terminal for the module.
// Arguments:
// - None.
// Returns:
// - *Terminal: The terminal for the container.
func (m *Dagindag) Terminal(
	// src is the directory that contains all the source code, including the module directory.
	// +optional
	src *Directory,
	// envVars is a set of strings (e.g., "KEY=value,KEY=value") to use as environment variables. They're
	// used to set the environment variables for the container when it's required to pass multiple environment variables
	// in a single argument. E.g.: "GITHUB_TOKEN=token,GO_VERSION=1.22.0,MYVAR=myvar"
	// +optional
	envVars []string,
) *Terminal {
	if len(envVars) > 0 {
		mMut, _ := m.WithEnvVarsFromStrs(envVars)
		m.Ctr = mMut.Ctr
	}

	if src != nil {
		m.Ctr = m.WithSource(src).Ctr
	}

	return m.Ctr.Terminal()
}

// DagCLI Allows to execute the Dagger CLI with the given flags.
//
// It allows to execute the Dagger CLI with the given flags.
// Arguments:
// - src is the directory that contains all the source code, including the module directory.
// - envVars is a set of strings (e.g., "KEY=value,KEY=value") to use as environment variables. They're
// used to set the environment variables for the container when it's required to pass multiple environment variables
// - flags is a set of string representing the flags to pass to the Dagger CLI.
// Returns:
// - string: The output of the function.
// - error: An error if the function fails.
func (m *Dagindag) DagCLI(
	// src is the directory that contains all the source code, including the module directory.
	// +optional
	src *Directory,
	// envVars is a set of strings (e.g., "KEY=value,KEY=value") to use as environment variables. They're
	// used to set the environment variables for the container when it's required to pass multiple environment variables
	// in a single argument. E.g.: "GITHUB_TOKEN=token,GO_VERSION=1.22.0,MYVAR=myvar"
	// +optional
	envVars []string,
	// dagCMDs is a set of string representing the flags to pass to the Dagger CLI. It must be set in a single string, such as "--dag-cmds="call -m module""
	// +optional
	dagCMDs string,
) (string, error) {
	if len(envVars) > 0 {
		mMut, _ := m.WithEnvVarsFromStrs(envVars)
		m.Ctr = mMut.Ctr
	}

	if src != nil {
		m.Ctr = m.WithSource(src).Ctr
	}

	var cmdsToRun []string

	if dagCMDs != "" {
		dagCMDsAsCMDs, err := cmdbuilder.GenerateDaggerCMDFromStr(dagCMDs)
		if err != nil {
			return "", fmt.Errorf("failed to generate Dagger CLI commands: %w", err)
		}

		cmdsToRun = append(cmdsToRun, dagCMDsAsCMDs...)
	}

	// Pass the commands to the container that are going through the Dagger CLI.
	out, err := m.Ctr.WithExec(cmdsToRun, ContainerWithExecOpts{
		ExperimentalPrivilegedNesting: true,
	}).Stdout(context.Background())

	if err != nil {
		return "", fmt.Errorf("failed to execute Dagger CLI: %w", err)
	}

	return out, nil
}

// UseFn calls a module with the given function and arguments.
//
// It calls a module with the given function and arguments.
// Arguments:
// - module: The name of the module, normally it's represented by the github repository name. E.g., "github.com/<owner>/<repo>/module@<version>"
// - version: The version of the module to use, e.g., "v0.11.5" If it's not specified, the latest version is used.
// - fn: The function to call in the module. If it's not specified, the default --help function is called.
// - args: The arguments to pass to the function. It receives a slice of strings. E.g., []string{"arg1", "arg2"}
// Returns:
// - string: The output of the function.
// - error: An error if the function fails.
func (m *Dagindag) UseFn(
	//modName is the name of the module, normally it's represented by the github repository name. E.g., "github.com/<owner>/<repo>/module@<version>"
	modName string,
	// modArgs are the arguments that the module's constructor requires. It receives a single slice of strings separated by commas. E.g., "arg1,arg2,arg3=value"
	// +optional
	modArgs []string,
	// modVersion is the version of the module to use, e.g., "v0.11.5" If it's not specified, the latest version is used.
	// +optional
	modVersion string,
	// fn is the function to call in the module. If it's not specified, the default --help function is called.
	fn string,
	// fnArgs is the arguments to pass to the function. It receives a single string separated by commas. E.g., "arg1,arg2,arg3=value"
	// +optional
	fnArgs []string,
	// src is the directory that contains all the source code, including the module directory.
	// +optional
	src *Directory,
	// envVars is a set of strings (e.g., "KEY=value,KEY=value") to use as environment variables. They're
	// used to set the environment variables for the container when it's required to pass multiple environment variables
	// in a single argument. E.g.: "GITHUB_TOKEN=token,GO_VERSION=1.22.0,MYVAR=myvar"
	// +optional
	envVars []string,
) (string, error) {
	// Resolve the module path with the given module and version
	dagModPath := getDaggerModulePath(modName, modVersion)
	daggerCallCMDWithMod := fmt.Sprintf("%s %s", daggerCallCMD, strings.TrimSpace(dagModPath))

	// If the function is not specified, call the default --help function
	if fn == "" {
		fn = "--help"
	}

	// Set environment variables if provided
	if len(envVars) > 0 {
		mMut, _ := m.WithEnvVarsFromStrs(envVars)
		m.Ctr = mMut.Ctr
	}

	// Set the source directory if provided
	if src != nil {
		m.Ctr = m.WithSource(src).Ctr
	}

	// Add module arguments to the command
	if len(modArgs) > 0 {
		for _, arg := range modArgs {
			daggerCallCMDWithMod = fmt.Sprintf("%s %s", daggerCallCMDWithMod, strings.TrimSpace(arg))
		}
	}

	// Add function to the command
	daggerCallCMDWithMod = fmt.Sprintf("%s %s", daggerCallCMDWithMod, fn)

	// Add function arguments to the command
	if len(fnArgs) > 0 {
		for _, arg := range fnArgs {
			daggerCallCMDWithMod = fmt.Sprintf("%s %s", daggerCallCMDWithMod, strings.TrimSpace(arg))
		}
	}

	// Execute the final command (Placeholder for actual execution)
	finalCmd := fmt.Sprintf("%s", daggerCallCMDWithMod)
	fmt.Println("Executing command:", finalCmd)

	cmdInDaggerFmt, err := cmdbuilder.GenerateDaggerCMDFromStr(finalCmd)
	if err != nil {
		return "", fmt.Errorf("failed to generate Dagger CLI commands: %w", err)
	}

	out, err := m.Ctr.WithExec(cmdInDaggerFmt, ContainerWithExecOpts{
		ExperimentalPrivilegedNesting: true,
	}).Stdout(context.Background())

	if err != nil {
		return "", fmt.Errorf("failed to execute Dagger command for module %s, and function %s: %w", modName, fn, err)
	}

	return out, nil
}
