// Package main provides the Terragrunt Dagger module and related functions.
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger. The module demonstrates
// usage of arguments and return types using simple echo and grep commands. The functions
// can be called from the dagger CLI or from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.
package main

import (
	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/containerx"
	"github.com/Excoriate/daggerx/pkg/envvars"
)

// Terragrunt is a Dagger module.
//
// This module is used to create and manage containers.
type Terragrunt struct {
	// Ctr is the container to use as a base container.
	Ctr *dagger.Container
	// ApkoPackages is a list of packages to install with APKO.
	// +private
	ApkoPackages []string
}

// New creates a new Terragrunt module.
//
// Parameters:
// - version: The version of the GoReleaser to use, e.g., "v1.22.0". Optional parameter.
// - image: The image to use as the base container. Optional parameter.
// - ctr: The container to use as a base container. Optional parameter.
// - envVarsFromHost: A list of environment variables to pass from the host to the container in a
// slice of strings. Optional parameter.
//
// Returns a pointer to a Terragrunt instance and an error, if any.
func New(
	// ctr is the container to use as a base container.
	// +optional
	ctr *dagger.Container,
	// tgVersion is the Terragrunt version to use. Default is "0.66.0".
	// +optional
	tgVersion string,
	// imageURL is the URL of the image to use as the base container.
	// It should includes tags. E.g. "ghcr.io/devops-infra/docker-terragrunt:tf-1.9.5-ot-1.8.2-tg-0.67.4"
	// +optional
	imageURL string,
	// tfVersion is the Terraform version to use. Default is "1.9.1".
	// +optional
	tfVersion string,
	// openTofuVersion is the OpenTofu version to use. Default is "1.8.0".
	// +optional
	openTofuVersion string,
	// envVarsFromHost is a list of environment variables to pass from the host to the container in a slice of strings.
	// +optional
	envVarsFromHost []string,
	// enableAWSCLI is a boolean to enable or disable the installation of the AWS CLI.
	// +optional
	enableAWSCLI bool,
	// awscliVersion is the version of the AWS CLI to install. Ensure the version is listed in the Alpine packages.
	// +optional
	awscliVersion string,
	// extraPackages is a list of extra packages to install with APKO, from the Alpine packages repository.
	// +optional
	extraPackages []string,
) (*Terragrunt, error) {
	dagModule := &Terragrunt{
		ApkoPackages: []string{},
	}
	// Precedence:
	// 1. ctr
	// 2. imageURL
	// 3. built-in base image
	if ctr != nil {
		dagModule.Ctr = ctr

		return dagModule, nil
	}

	if enableAWSCLI {
		dagModule.WithAWSCLIPackage(awscliVersion)
	}

	dagModule.WithExtraPackages(extraPackages...)

	if imageURL != "" {
		isValid, err := containerx.ValidateImageURL(imageURL)
		if err != nil {
			return nil, WrapErrorf(err, "failed to validated image URL: %s", imageURL)
		}

		if !isValid {
			return nil, Errorf("the image URL %s is not valid", imageURL)
		}

		dagModule.Base(imageURL)
	} else {
		_, tgCtrErr := dagModule.BaseApko(dagModule.ApkoPackages)
		if tgCtrErr != nil {
			return nil, WrapError(tgCtrErr, "failed to create base image apko")
		}

		dagModule.WithTerragruntCacheConfiguration()
		dagModule.WithTerraformCacheConfiguration()
		dagModule.WithIACToolsInstalled(tgVersion, tfVersion, openTofuVersion)
		dagModule.WithTerragruntPermissions()
	}

	if len(envVarsFromHost) > 0 {
		if err := addEnvVars(dagModule, envVarsFromHost); err != nil {
			return nil, err
		}
	}

	return dagModule, nil
}

// addEnvVars adds environment variables from the host to the Terragrunt configuration.
// It parses the environment variables and adds them to the Terragrunt instance.
func addEnvVars(terragrunt *Terragrunt, envVarsFromHost []string) error {
	envVars, err := envvars.ToDaggerEnvVarsFromSlice(envVarsFromHost) // Parse environment variables from the host
	if err != nil {
		return WrapError(err, "failed to parse environment variables") // Return error if parsing fails
	}

	for _, envVar := range envVars {
		terragrunt.WithEnvironmentVariable(envVar.Name, envVar.Value, false) // Add each environment variable to Terragrunt
	}

	return nil
}
