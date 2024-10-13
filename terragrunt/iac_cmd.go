package main

import (
	"context"

	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
)

// Cmd is an interface that represents a command to be executed by Terragrunt.
// Cmd defines the interface for executing commands in Infrastructure as Code (IaC) tools
// such as Terraform, Terragrunt, and OpenTofu within a Dagger container.
type Cmd interface {
	// Exec executes a given IaC command within a Dagger container.
	// Returns a pointer to the resulting dagger.Container or an error.
	Exec(
		ctx context.Context,
		command string,
		args []string,
		autoApprove bool,
		source *dagger.Directory,
		module string,
		envVars []string,
		secrets []*dagger.Secret,
	) (*dagger.Container, error)

	// validate checks if the provided command is recognized by the IaC tool.
	// Returns an error for invalid or empty commands.
	validate(command string) error

	// getEntrypoint returns the executable entrypoint for the IaC tool.
	getEntrypoint() string
}
