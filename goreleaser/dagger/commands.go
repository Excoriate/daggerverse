package main

import (
	"context"

	"github.com/Excoriate/daggerx/pkg/cmdbuilder"

	"github.com/Excoriate/daggerx/pkg/merger"
)

const (
	goReleaserCMDClean = "--clean"
)

// Version runs the 'goreleaser --version' command.
// It's equivalent to running 'goreleaser --version' in the terminal.
func (m *Goreleaser) Version() (string, error) {
	m.Ctr = m.WithCMD([]string{"--version"}).Ctr

	out, err := m.Ctr.
		Stdout(context.Background())

	return out, err
}

// ShowEnvVars runs the 'printenv' command to show the environment variables.
func (m *Goreleaser) ShowEnvVars(
	// envVars is a map of environment (but as a slice) variables to pass from the host to the container.
	// +optional
	// +default=[]
	envVars []string,
) (string, error) {
	m.addEnvVarsToContainerFromSlice(envVars)
	m.Ctr = replaceEntryPointForShell(m.Ctr)
	m.Ctr = m.WithCMD([]string{"printenv"}).Ctr

	out, err := m.Ctr.
		Stdout(context.Background())

	return out, err
}

// Check runs the 'goreleaser check' command.
// It's equivalent to running 'goreleaser check' in the terminal.
func (m *Goreleaser) Check(
	// src is the source directory.
	src *Directory,
	// cfg is the configuration file to use.
	// +optional
	cfg string,
	// args is the arguments to pass to the 'goreleaser' command.
	// +optional
	args string,
) (string, error) {
	cfg = setToDefaultCfgIfEmpty(cfg)

	allArgs := cmdbuilder.BuildArgs(args, m.resolveCfgArg(cfg))

	m.Ctr = m.WithSource(src).Ctr
	m.Ctr = m.WithCMD(merger.MergeSlices([]string{"check"}, allArgs)).Ctr

	out, err := m.Ctr.
		Stdout(context.Background())

	return out, err
}

// Build runs the 'goreleaser build' command.
// It's equivalent to running 'goreleaser build' in the terminal.
func (m *Goreleaser) Build(
	// src is the directory that contains all the source code, including the module directory.
	src *Directory,
	// cfg is the configuration file to use.
	// +optional
	// default=".goreleaser.yaml"
	cfg string,
	// clean ensures that if there's a previous build, it will be cleaned.
	// +optional
	clean bool,
	// envVars is a list of environment variables to pass from the host to the container.
	// +optional
	envVars []string,
	// args is the arguments to pass to the 'goreleaser' command.
	// +optional
	args string,
) (string, error) {
	cfg = setToDefaultCfgIfEmpty(cfg)
	m.Ctr = m.WithSource(src).Ctr

	cfgFileArg := m.resolveCfgArg(cfg)
	cleanArg := ""

	if clean {
		cleanArg = goReleaserCMDClean
	}

	allArgs := cmdbuilder.BuildArgs(args, cfgFileArg, cleanArg)

	m.addEnvVarsToContainerFromSlice(envVars)
	m.Ctr = m.WithCMD(merger.MergeSlices([]string{"build"}, allArgs)).Ctr

	out, err := m.Ctr.
		Stdout(context.Background())

	return out, err
}

// Snapshot runs the 'goreleaser release' command.
// It's equivalent to running 'goreleaser release --snapshot' in the terminal.
func (m *Goreleaser) Snapshot(
	// cfg is the configuration file to use.
	// +optional
	cfg string,
	// autoSnapshot ensures that the snapshot is automatically set if the repository is dirty
	// +optional
	autoSnapshot bool,
	// clean ensures that if there's a previous build, it will be cleaned.
	// +optional
	// default=false
	clean bool,
	// envVars is a list of environment variables to pass from the host to the container.
	// +optional
	// +default=[]
	envVars []string,
	// args is the arguments to pass to the 'goreleaser' command.
	// +optional
	args string,
) (string, error) {
	cfgFileArg := m.resolveCfgArg(cfg)
	autoSnapshotArg := ""
	snapshotArg := "--snapshot"
	cleanArg := ""
	if autoSnapshot {
		autoSnapshotArg = "--auto-snapshot"
	}

	if clean {
		cleanArg = "--clean"
	}

	allArgs := cmdbuilder.BuildArgs(args, cfgFileArg, snapshotArg, autoSnapshotArg, cleanArg)

	m.addEnvVarsToContainerFromSlice(envVars)
	m.Ctr = m.WithCMD(merger.MergeSlices([]string{"release"}, allArgs)).Ctr

	out, err := m.Ctr.
		Stdout(context.Background())

	return out, err
}

// Release runs the 'goreleaser release' command.
// It's equivalent to running 'goreleaser release --snapshot' in the terminal.
func (m *Goreleaser) Release(
	// cfg is the configuration file to use.
	// +optional
	cfg string,
	// autoSnapshot ensures that the snapshot is automatically set if the repository is dirty
	// +optional
	autoSnapshot bool,
	// clean ensures that if there's a previous build, it will be cleaned.
	// +optional
	// default=false
	clean bool,
	// envVars is a list of environment variables to pass from the host to the container.
	// +optional
	// +default=[]
	envVars []string,
	// args is the arguments to pass to the 'goreleaser' command.
	// +optional
	args string,
) (string, error) {
	cfgFileArg := m.resolveCfgArg(cfg)
	autoSnapshotArg := ""
	cleanArg := ""
	if autoSnapshot {
		autoSnapshotArg = "--auto-snapshot"
	}

	if clean {
		cleanArg = "--clean"
	}

	allArgs := cmdbuilder.BuildArgs(args, cfgFileArg, autoSnapshotArg, cleanArg)

	m.addEnvVarsToContainerFromSlice(envVars)
	m.Ctr = m.WithCMD(merger.MergeSlices([]string{"release"}, allArgs)).Ctr

	out, err := m.Ctr.
		Stdout(context.Background())

	return out, err
}
