package main

import (
	"context"

	"github.com/Excoriate/daggerverse/gotoolbox/internal/dagger"
)

// RunGo runs a Go command within a given context.
//
// cmd is the Go command to run, everything after the 'go' command.
// src is the optional source directory for the container.
//
// It returns the standard output of the executed command or an error if something goes wrong.
func (m *Gotoolbox) RunGo(
	// cmd is the Go command to run, everything after the 'go' command.
	cmd string,
	// src is the optional source directory for the container.
	// +optional
	src *dagger.Directory) (string, error) {
	tb := m

	// Conditionally include the source if 'src' is provided
	if src != nil {
		tb = m.WithSource(src, "")
	}

	tb = tb.WithGoExec([]string{cmd}, "")

	return tb.Ctr.Stdout(context.Background())
}
