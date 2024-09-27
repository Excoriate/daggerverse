package main

import (
	"context"

	"github.com/Excoriate/daggerverse/module-template/internal/dagger"
)

const (
	validKeyPattern = `^[a-zA-Z0-9_]+$`
)

// OpenTerminal returns a terminal
//
// It returns a terminal for the container.
// Arguments:
// - None.
// Returns:
// - *Terminal: The terminal for the container.
func (m *ModuleTemplate) OpenTerminal() *dagger.Container {
	return m.Ctr.Terminal()
}

// RunShell runs a shell command in the container.
//
// It runs a shell command in the container and returns the output.
// Arguments:
// - cmd: The command to run in the container.
// Returns:
// - string: The output of the command.
// - error: An error if the command fails.
func (m *ModuleTemplate) RunShell(cmd string) (string, error) {
	out, err := m.Ctr.WithoutEntrypoint().WithExec([]string{"sh", "-c", cmd}).Stdout(context.Background())
	if err != nil {
		return "", WrapErrorf(err, "failed to run shell command: %s", cmd)
	}

	return out, nil
}
