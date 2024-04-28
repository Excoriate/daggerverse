package main

import (
	"context"
	"fmt"
)

const (
	tfLintDefaultImage   = "ghcr.io/terraform-linters/tflint"
	tfLintDefaultVersion = "v0.50.3"
	workdirRootPath      = "/mnt"
)

type Tflint struct {
	// Src is the directory that contains all the source code, including the module directory.
	Src *Directory

	// Ctr is the container to use as a base container.
	Ctr *Container
}

func New(
	// version is the version of the TFLint to use, e.g., "v0.50.3". For more information, see https://github.com/terraform-linters/tflint
	// +default="v0.50.3"
	// +optional
	version string,
	// image is the image to use as the base container.
	// +optional
	// +default="ghcr.io/terraform-linters/tflint"
	image string,
	// src is the directory that contains all the source code, including the module directory.
	src *Directory,
	// Ctrl is the container to use as a base container.
	// +optional
	ctr *Container,
) *Tflint {
	g := &Tflint{
		Src: src,
	}

	if ctr != nil {
		g.Ctr = ctr
	} else {
		g.Base(image, version)
	}

	g = g.WithSource(src)

	return g
}

// Base sets the base image and version, and creates the base container.
func (m *Tflint) Base(image, version string) *Tflint {
	if image == "" {
		image = tfLintDefaultImage
	}

	if version == "" {
		version = tfLintDefaultVersion
	}

	ctrImage := fmt.Sprintf("%s:%s", image, version)

	c := dag.Container().From(ctrImage).
		WithWorkdir(workdirRootPath)

	m.Ctr = c

	return m
}

// WithSource sets the source directory if it's passed, and
// mounts the source directory to the container.
func (m *Tflint) WithSource(src *Directory) *Tflint {
	if src != nil {
		m.Src = src
	}

	m.Ctr = m.Ctr.WithMountedDirectory(workdirRootPath, m.Src)

	return m
}

// Version runs the 'tflint --version' command.
func (m *Tflint) Version() (string, error) {
	m.Ctr = addCMDsToContainer([]string{"--version"}, []string{}, m.Ctr)

	out, err := m.Ctr.
		Stdout(context.Background())

	return out, err
}

// Run executes any tflint command.
func (m *Tflint) Run(
	// cfg is the configuration file to use.
	// +optional
	cfg string,
	// args is the arguments to pass to the tfLint command.
	// +optional
	args string,
) (string, error) {
	cfgFileArg := ""
	if cfg != "" {
		cfgFileArg = fmt.Sprintf("--config=%s", cfg)
	}

	allArgs := buildArgs(args, cfgFileArg)
	m.Ctr = addCMDsToContainer(allArgs, []string{}, m.Ctr)

	out, err := m.Ctr.
		Stdout(context.Background())

	return out, err
}

// WithInit adds the 'init' command to the container.
func (m *Tflint) WithInit(
	// cfg is the configuration file to use.
	// +optional
	cfg string,
	// args is the arguments to pass to the tfLint init command.
	// +optional
	args string,
) *Tflint {
	cfgFileArg := ""
	if cfg != "" {
		cfgFileArg = fmt.Sprintf("--config=%s", cfg)
	}

	allArgs := buildArgs(args, cfgFileArg)

	m.Ctr = addCMDsToContainer([]string{"init"}, allArgs, m.Ctr)
	return m
}

// RunInit executes the 'init' command.
func (m *Tflint) RunInit(
	// cfg is the configuration file to use.
	// +optional
	cfg string,
	// args is the arguments to pass to the tfLint init command.
	// +optional
	args string,
) (string, error) {
	cfgFileArg := ""
	if cfg != "" {
		cfgFileArg = fmt.Sprintf("--config=%s", cfg)
	}

	allArgs := buildArgs(args, cfgFileArg)

	m.Ctr = addCMDsToContainer([]string{"--init"}, allArgs, m.Ctr)

	out, err := m.Ctr.
		Stdout(context.Background())

	return out, err
}

func (m *Tflint) Lint(
	// init specifies whether to run the 'init' command before running the 'lint' command.
	// +optional
	// +default=false
	init bool,
	// cfg is the configuration file to use.
	// +optional
	cfg string,
	// args is the arguments to pass to the tfLint init command.
	// +optional
	args string,
) *Tflint {
	initArg := ""
	cfgFileArg := ""
	if cfg != "" {
		cfgFileArg = fmt.Sprintf("--config=%s", cfg)
	}

	if init {
		initArg = "--init"
	}

	allArgs := buildArgs(args, cfgFileArg, initArg)

	m.Ctr = addCMDsToContainer([]string{}, allArgs, m.Ctr)

	return m
}

// RunLint executes the 'init' command.
// It's equivalent to running 'tflint' in the terminal.
func (m *Tflint) RunLint(
	// init specifies whether to run the 'init' command before running the 'lint' command.
	// +optional
	// +default=false
	init bool,
	// cfg is the configuration file to use.
	// +optional
	cfg string,
	// args is the arguments to pass to the tfLint init command.
	// +optional
	args string,
) (string, error) {
	initArg := ""
	cfgFileArg := ""
	if cfg != "" {
		cfgFileArg = fmt.Sprintf("--config=%s", cfg)
	}

	if init {
		initArg = "--init"
	}

	allArgs := buildArgs(args, cfgFileArg, initArg)

	m.Ctr = addCMDsToContainer([]string{}, allArgs, m.Ctr)

	out, err := m.Ctr.
		Stdout(context.Background())

	return out, err
}
