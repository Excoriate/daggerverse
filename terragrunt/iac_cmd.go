package main

import (
	"context"

	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
)

// Cmd defines the interface for executing commands in Infrastructure as Code (IaC) tools
// such as Terraform, Terragrunt, and OpenTofu within a Dagger container.
// Cmd defines the interface for executing commands in Infrastructure as Code (IaC) tools
// such as Terraform, Terragrunt, and OpenTofu within a Dagger container.
type Cmd interface {
	// Exec executes a given IaC command within a Dagger container.
	//
	// Parameters:
	// - ctx: The context for controlling the execution.
	// - command: The IaC command to execute.
	// - args: The arguments for the command.
	// - autoApprove: Flag to auto-approve prompts.
	// - source: The source directory for the command.
	// - module: The module to execute or the terragrunt configuration where the terragrunt.hcl file is located. +optional
	// - envVars: The environment variables to pass to the container. +optional
	// - secrets: The secrets to pass to the container. +optional
	// - tool: The tool to use for executing the command.
	//
	// Returns:
	// - *dagger.Container: Pointer to the resulting container.
	// - error: If execution fails.
	Exec(
		ctx context.Context,
		command string,
		args []string,
		autoApprove bool,
		source *dagger.Directory,
		module string,
		envVars []string,
		secrets []*dagger.Secret,
		tool string,
	) (*dagger.Container, error)

	// ExecCmd executes a given command within a Dagger container.
	//
	// Parameters:
	// - ctx: The context for controlling the execution.
	// - command: The command to execute.
	// - args: The arguments for the command.
	// - autoApprove: Flag to auto-approve prompts.
	// - source: The source directory for the command.
	// - module: The module to execute. +optional
	// - envVars: The environment variables to pass to the container. +optional
	// - secrets: The secrets to pass to the container. +optional
	// - tool: The tool to use for executing the command.
	//
	// Returns:
	// - *dagger.Container: Pointer to the resulting container.
	// - error: If execution fails.
	ExecCmd(
		ctx context.Context,
		command string,
		args []string,
		autoApprove bool,
		source *dagger.Directory,
		module string,
		envVars []string,
		secrets []*dagger.Secret,
		tool string,
	) (*dagger.Container, error)

	// validate checks if the provided command is recognized by the IaC tool.
	//
	// Returns an error for invalid or empty commands.
	validate(command string) error

	// getEntrypoint returns the executable entrypoint for the IaC tool.
	getEntrypoint() string
}
