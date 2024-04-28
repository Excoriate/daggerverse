package main

import (
	"context"
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

	// tgVersion is the version of the Terragrunt to use, e.g., "v0.55.20". For more information, visit's Gruntwork's Terragrunt GitHub repository.
	// +default="v0.55.20"
	// +optional
	tgVersion string

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

	g = g.WithSource(src)

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
