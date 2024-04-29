package main

import (
	"context"
	"fmt"
)

const (
	tgContainerImageDefault = "ghcr.io/terraform-linters/tflint"
	tfVersionDefault        = "1.7.0"
	workdirRootPath         = "/mnt"
	entrypointCMD           = "terragrunt"
)

type Terragrunt struct {
	// Src is the directory that contains all the source code, including the module directory.
	Src *Directory

	// Ctr is the container to use as a base container.
	Ctr *Container

	// tfVersion is the version of the Terraform to use, e.g., "1.0.0". For more information, visit HashiCorp's Terraform GitHub repository.
	// +optional
	// +default="1.7.0"
	tfVersion string
}

func New(
	// tfVersion is the version of the Terraform to use, e.g., "1.0.0". For more information, visit HashiCorp's Terraform GitHub repository.
	// +optional
	// +default="1.7.0"
	tfVersion string,
	// image is the image to use as the base container.
	// +optional
	// +default="alpine/terragrunt"
	image string,
	// src is the directory that contains all the source code, including the module directory.
	src *Directory,
	// Ctrl is the container to use as a base container.
	// +optional
	ctr *Container,
) *Terragrunt {
	g := &Terragrunt{
		Src: src,
	}

	if ctr != nil {
		g.Ctr = ctr
	} else {
		g.Base(image, tfVersion)
	}

	g = g.WithSource(src, "")

	return g
}

// Version returns the Terragrunt version.
// Consider to configure Terragrunt based on the target Terraform version.
func (m *Terragrunt) Version() (string, error) {
	m.Ctr = addCMDsToContainer([]string{entrypointCMD}, []string{"--version"}, m.Ctr)

	out, err := m.Ctr.
		Stdout(context.Background())

	return out, err
}

// Help returns the Terragrunt help.
func (m *Terragrunt) Help() (string, error) {
	m.Ctr = addCMDsToContainer([]string{entrypointCMD}, []string{"--help"}, m.Ctr)

	out, err := m.Ctr.
		Stdout(context.Background())

	return out, err
}

// AddCMD adds a command to the container.
// It supports environment variables and arguments.
func (m *Terragrunt) AddCMD(
	// module is the terragunt module to use that includes the terragrunt.hcl
	// the module should be relative to the mounted directory (--src).
	module string,
	// cmd is the command to run.
	cmd string,
	// envVars is the list of environment variables to pass from the host to the container.
	// the format of the environment variables passed from the host are slices of strings separted by "=", and commas.
	// E.g., []string{"HOST=localhost", "PORT=8080"}
	// +optional
	// +default=[]
	envVars []string,
	// args is the list of arguments to pass to the Terragrunt command.
	// +optional
	args string,
) (*Terragrunt, error) {
	envVarsDaggerFormat, err := toEnvVarsDagger(envVars)
	if err != nil {
		return nil, err
	}

	// Add the terragrunt module as the workdir.
	m.Ctr = m.WithSource(m.Src, module).Ctr

	// Add the environment variables to the container.
	m.Ctr = m.WithEnvVars(envVarsDaggerFormat).Ctr

	// Add the command to the container.
	m.Ctr = addCMDsToContainer([]string{entrypointCMD, cmd}, buildArgs(args), m.Ctr)

	return m, nil
}

// Run executes any Terragrunt command.
func (m *Terragrunt) Run(
	// module is the terragunt module to use that includes the terragrunt.hcl
	// the module should be relative to the mounted directory (--src).
	module string,
	// cmd is the command to run.
	cmd string,
	// envVars is the list of environment variables to pass from the host to the container.
	// the format of the environment variables passed from the host are slices of strings separted by "=", and commas.
	// E.g., []string{"HOST=localhost", "PORT=8080"}
	// +optional
	// +default=[]
	envVars []string,
	// args is the list of arguments to pass to the Terragrunt command.
	// +optional
	args string,
) (string, error) {
	tgCMD, err := m.AddCMD(module, cmd, envVars, args)
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
	// cmd is the command to run.
	cmd string,
	// envVars is the list of environment variables to pass from the host to the container.
	// the format of the environment variables passed from the host are slices of strings separted by "=", and commas.
	// E.g., []string{"HOST=localhost", "PORT=8080"}
	// +optional
	// +default=[]
	envVars []string,
	// args is the list of arguments to pass to the Terragrunt command.
	// +optional
	args string,
) (string, error) {
	runAllCMD := fmt.Sprintf("run-all %s", cmd)
	nonInteractive := "--terragrunt-non-interactive"
	allArgs := fmt.Sprintf("%s %s", nonInteractive, args)

	fullCMD := append([]string{entrypointCMD, runAllCMD}, buildArgs(allArgs)...)
	fullCMDAsStr := convertCMDToString(fullCMD)

	// Add the terragrunt module as the workdir.
	m.Ctr = m.WithSource(m.Src, module).Ctr

	// Add the environment variables to the container.
	envVarsDaggerFormat, err := toEnvVarsDagger(envVars)
	if err != nil {
		return "", err
	}

	m.Ctr = m.WithEnvVars(envVarsDaggerFormat).Ctr

	out, err := m.Ctr.
		WithEntrypoint(nil).
		// Somehow the run-all command is not working with the entrypoint, so I'm using the full command as a string.
		// and running it with the shell.
		WithExec([]string{"sh", "-c", fullCMDAsStr}).
		Stdout(context.Background())

	return out, err
}
