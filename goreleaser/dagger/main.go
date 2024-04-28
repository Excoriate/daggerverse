package main

import (
	"context"
	"fmt"
)

const (
	goReleaserDefaultVersion = "latest"
	goReleaserDefaultImage   = "goreleaser/goreleaser"
	mntPrefix                = "/mnt"
	goReleaserDefaultCfgFile = ".goreleaser.yaml"
)

type Goreleaser struct {
	// Src is the directory that contains all the source code, including the module directory.
	Src *Directory

	// Ctr is the container to use as a base container.
	Ctr *Container

	// CfgFile is the configuration file to use.
	CfgFile string

	// EnvVarsFromHost is a list of environment variables to pass from the host to the container.
	// Later on, in order to pass it to the container, it's going to be converted into a map.
	EnvVarsFromHost []string
}

func New(
	// version is the version of the GoReleaser to use, e.g., "v1.22.0".
	// +default="latest"
	// +optional
	version string,
	// image is the image to use as the base container.
	// +optional
	// +default="goreleaser/goreleaser"
	image string,
	// src is the directory that contains all the source code, including the module directory.
	src *Directory,
	// Ctrl is the container to use as a base container.
	// +optional
	ctr *Container,
) *Goreleaser {
	g := &Goreleaser{
		Src: src,
	}

	if ctr != nil {
		g.Ctr = ctr
	} else {
		g.Base(image, version)
	}

	// Here, regardless of whether the container is set or not,
	// the WithGoCache method is called.
	g = g.WithGoCache()

	return g
}

// WithCfgFile sets the configuration file to use.
// The default configuration file is ".goreleaser.yaml".
func (g *Goreleaser) WithCfgFile(
	// cfgFile is the configuration file to use.
	// +optional
	// default=".goreleaser.yaml"
	cfgFile string,
) *Goreleaser {
	g.CfgFile = cfgFile

	return g
}

// WithPrintEnv passes environment variables from the host to the container.
func (g *Goreleaser) WithPrintEnv(
	// envVars is a list of environment variables to pass from the host to the container.
	// Later on, in order to pass it to the container, it's going to be converted into a map.
	// +optional
	// +default=[]
	envVars []string,
) *Goreleaser {
	g.EnvVarsFromHost = envVars

	return g
}

// WithGoCache mounts the Go cache directories.
// The Go cache directories are:
// - /go/pkg/mod
// - /root/.cache/go-build
func (g *Goreleaser) WithGoCache() *Goreleaser {
	goModCache := dag.CacheVolume("gomodcache")
	goBuildCache := dag.CacheVolume("gobuildcache")

	ctr := g.Ctr.WithMountedCache("/go/pkg/mod", goModCache).
		WithMountedCache("/root/.cache/go-build", goBuildCache)

	g.Ctr = ctr

	return g
}

// Base sets the base image and version, and creates the base container.
// The default image is "goreleaser/goreleaser" and the default version is "latest".
func (g *Goreleaser) Base(image, version string) *Goreleaser {
	if image == "" {
		image = goReleaserDefaultImage
	}

	if version == "" {
		version = goReleaserDefaultVersion
	}

	ctrImage := fmt.Sprintf("%s:%s", image, version)

	c := dag.Container().From(ctrImage).
		WithEnvVariable("TINI_SUBREAPER", "true").
		WithWorkdir(mntPrefix).
		WithMountedDirectory(mntPrefix, g.Src)

	g.Ctr = c

	return g
}

// Version runs the 'goreleaser --version' command.
// It's equivalent to running 'goreleaser --version' in the terminal.
func (g *Goreleaser) Version() (string, error) {
	g.Ctr = addCMDsToContainer([]string{"--version"}, []string{}, g.Ctr)

	out, err := g.Ctr.
		Stdout(context.Background())

	return out, err
}

// Check runs the 'goreleaser check' command.
// It's equivalent to running 'goreleaser check' in the terminal.
func (g *Goreleaser) Check(
	// cfg is the configuration file to use.
	// +optional
	// default=".goreleaser.yaml"
	cfg string,
	// args is the arguments to pass to the 'goreleaser' command.
	// +optional
	args string,
) (string, error) {
	allArgs := buildArgs(args, g.resolveCfgArg(cfg))

	g.Ctr = addCMDsToContainer([]string{"check"}, allArgs, g.Ctr)

	out, err := g.Ctr.
		Stdout(context.Background())

	return out, err
}

// ShowEnvVars runs the 'printenv' command to show the environment variables.
func (g *Goreleaser) ShowEnvVars(
	// envVars is a map of environment (but as a slice) variables to pass from the host to the container.
	// +optional
	// +default=[]
	envVars []string,
) (string, error) {
	allEnvVars := mergeSlices(envVars, g.EnvVarsFromHost)
	envVarsAsMap, err := toEnvVars(allEnvVars)
	if err != nil {
		return "", err
	}

	g.Ctr = addEnvVarsToContainer(envVarsAsMap, g.Ctr)
	g.Ctr = replaceEntryPointForShell(g.Ctr)
	g.Ctr = addCMDsToContainer([]string{"printenv"}, []string{}, g.Ctr)

	out, err := g.Ctr.
		Stdout(context.Background())

	return out, err
}

// Build runs the 'goreleaser build' command.
// It's equivalent to running 'goreleaser build' in the terminal.
func (g *Goreleaser) Build(
	// cfg is the configuration file to use.
	// +optional
	// default=".goreleaser.yaml"
	cfg string,
	// clean ensures that if there's a previous build, it will be cleaned.
	// +optional
	clean bool,
	// envVars is a list of environment variables to pass from the host to the container.
	// +optional
	// +default=[]
	envVars []string,
	// args is the arguments to pass to the 'goreleaser' command.
	// +optional
	args string,
) (string, error) {
	cfgFileArg := g.resolveCfgArg(cfg)
	cleanArg := ""
	envVarsMap, err := toEnvVars(envVars)
	if err != nil {
		return "", err
	}

	if clean {
		cleanArg = "--clean"
	}

	allArgs := buildArgs(args, cfgFileArg, cleanArg)

	g.Ctr = addEnvVarsToContainer(envVarsMap, g.Ctr)
	g.Ctr = addCMDsToContainer([]string{"build"}, allArgs, g.Ctr)

	out, err := g.Ctr.
		Stdout(context.Background())

	return out, err
}

// Snapshot runs the 'goreleaser release' command.
// It's equivalent to running 'goreleaser release --snapshot' in the terminal.
func (g *Goreleaser) Snapshot(
	// cfg is the configuration file to use.
	// +optional
	// default=".goreleaser.yaml"
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
	cfgFileArg := g.resolveCfgArg(cfg)
	autoSnapshotArg := ""
	snapshotArg := "--snapshot"
	cleanArg := ""
	envVarsMap, err := toEnvVars(envVars)
	if err != nil {
		return "", err
	}

	if autoSnapshot {
		autoSnapshotArg = "--auto-snapshot"
	}

	if clean {
		cleanArg = "--clean"
	}

	allArgs := buildArgs(args, cfgFileArg, snapshotArg, autoSnapshotArg, cleanArg)

	g.Ctr = addEnvVarsToContainer(envVarsMap, g.Ctr)
	g.Ctr = addCMDsToContainer([]string{"release"}, allArgs, g.Ctr)

	out, err := g.Ctr.
		Stdout(context.Background())

	return out, err
}

// Release runs the 'goreleaser release' command.
// It's equivalent to running 'goreleaser release --snapshot' in the terminal.
func (g *Goreleaser) Release(
	// cfg is the configuration file to use.
	// +optional
	// default=".goreleaser.yaml"
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
	cfgFileArg := g.resolveCfgArg(cfg)
	autoSnapshotArg := ""
	cleanArg := ""
	envVarsMap, err := toEnvVars(envVars)
	if err != nil {
		return "", err
	}

	if autoSnapshot {
		autoSnapshotArg = "--auto-snapshot"
	}

	if clean {
		cleanArg = "--clean"
	}

	allArgs := buildArgs(args, cfgFileArg, autoSnapshotArg, cleanArg)

	g.Ctr = addEnvVarsToContainer(envVarsMap, g.Ctr)
	g.Ctr = addCMDsToContainer([]string{"release"}, allArgs, g.Ctr)

	out, err := g.Ctr.
		Stdout(context.Background())

	return out, err
}
