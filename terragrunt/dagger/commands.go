package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/Excoriate/daggerx/pkg/cmdbuilder"
)

// Version returns the Terragrunt version.
// Consider to configure Terragrunt based on the target Terraform version.
func (m *Terragrunt) Version() (string, error) {
	cmd := []string{tgCMD, "--version"}
	m.Ctr = m.WithCMD(cmd).Ctr

	out, err := m.Ctr.
		Stdout(context.Background())

	return out, err
}

// Help returns the Terragrunt help.
func (m *Terragrunt) Help() (string, error) {
	cmd := []string{tgCMD, "--help"}
	m.Ctr = m.WithCMD(cmd).Ctr

	out, err := m.Ctr.
		Stdout(context.Background())

	return out, err
}

// Run executes any Terragrunt command.
func (m *Terragrunt) Run(
	// module is the terragunt module to use that includes the terragrunt.hcl
	// the module should be relative to the mounted directory (--src).
	module string,
	// cmd is the Terragrunt command to run.
	cmd string,
	// envVars is the list of environment variables to pass from the host to the container.
	// the format of the environment variables passed from the host are slices of strings separted by "=", and commas.
	// E.g., []string{"HOST=localhost", "PORT=8080"}
	// +optional
	envVars []string,
	// args is the list of arguments to pass to the Terragrunt command.
	// +optional
	args string,
) (string, error) {
	tgCMD, err := m.addCMD(module, cmd, envVars, args)
	if err != nil {
		return "", err
	}

	out, err := tgCMD.Ctr.
		Stdout(context.Background())

	return out, err
}

// RunAll executes any Terragrunt run-all command.
func (m *Terragrunt) RunAll(
	// module is the terragunt module to use that includes the terragrunt.hcl
	// the module should be relative to the mounted directory (--src).
	module string,
	// cmd is the Terragrunt command to run.
	cmd string,
	// envVars is the list of environment variables to pass from the host to the container.
	// the format of the environment variables passed from the host are slices of strings separted by "=", and commas.
	// E.g., []string{"HOST=localhost", "PORT=8080"}
	// +optional
	envVars []string,
	// args is the list of arguments to pass to the Terragrunt command.
	// +optional
	args string,
) (string, error) {
	runAllCMD := fmt.Sprintf("run-all %s", cmd)
	allArgs := fmt.Sprintf("%s %s", tgNonInteractiveFlag, args)

	fullCMD := append([]string{tgCMD, runAllCMD}, cmdbuilder.BuildArgs(allArgs)...)
	fullCMDAsStr := strings.Join(fullCMD, " ")
	fullCMDAsStr = strings.ReplaceAll(fullCMDAsStr, ",", "")

	// Add the terragrunt module as the workdir.
	m.Ctr = m.WithSource(m.Src, module).Ctr

	//// Add the environment variables to the container.
	m.addEnvVarsToContainerFromSlice(envVars)

	out, err := m.Ctr.
		WithEntrypoint(nil).
		// Somehow the run-all command is not working with the entrypoint, so I'm using the full command as a string.
		// and running it with the shell.
		WithExec([]string{"sh", "-c", fullCMDAsStr}).
		Stdout(context.Background())

	return out, err
}
