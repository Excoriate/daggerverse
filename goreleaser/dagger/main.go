package main

import (
	"fmt"
	daggerArgs "github.com/excoriate/daggerverse/daggerx/pkg/args"
	"log/slog"
)

const (
	defaultImage   = "goreleaser/goreleaser-cross-base"
	defaultVersion = "v1.22.0"
	mntPrefix      = "/mnt"
	goreleaserCMD  = "goreleaser"
)

type Goreleaser struct {
	// Src is the directory that contains all the source code, including the module directory.
	Src *Directory

	// Ctr is the container to use as a base container.
	Ctr *Container

	// CfgFile is the configuration file to use.
	CfgFile *File
}

func New(
	// version is the version of the GoReleaser to use, e.g., "v1.22.0".
	// +default="v1.22.0"
	// +optional
	version string,

	// image is the image of the container to use.
	// +default="goreleaser/goreleaser-cross-base"
	// +optional
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

	if ctr == nil {
		g.Base(image, version)
	} else {
		slog.Info("This module has received a container as an argument. The container will be used as the base container.")
		g.Ctr = ctr
	}

	// Here, regardless of whether the container is set or not,
	// the WithGoCache method is called.
	g = g.WithGoCache()

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
func (g *Goreleaser) Base(image, version string) *Goreleaser {
	if image == "" {
		image = defaultImage
	}

	if version == "" {
		version = defaultVersion
	}

	slog.Info(fmt.Sprintf("Creating a new GoReleaser module with version %s", version))
	slog.Info(fmt.Sprintf("Creating a new GoReleaser module with image %s", image))

	c := dag.Container().From(image).
		WithMountedDirectory(mntPrefix, g.Src)

	g.Ctr = c
	return g
}

func (g *Goreleaser) Check(
	//// +optional
	//// +default=".goreleaser.yaml"
	//cfgFile string,
	// args is the arguments to pass to the 'goreleaser' command.
	// +optional
	args string,
) (*Container, error) {
	parsedArgs := daggerArgs.ParseArgsFromStrToSlice(args)
	cmdToRun := []string{goreleaserCMD, "check"}
	cmdToRun = append(cmdToRun, parsedArgs...)

	g.Ctr = g.Ctr.
		WithWorkdir(mntPrefix).
		WithExec(cmdToRun)

	return g.Ctr, nil
}
