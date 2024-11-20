// Package main provides the AwsTagInspector Dagger module and related functions.
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
	"context"
	"fmt"

	"github.com/Excoriate/daggerverse/aws-tag-inspector/internal/dagger"

	"github.com/Excoriate/daggerx/pkg/envvars"
)

// AwsTagInspector is a Dagger module.
//
// This module is used to create and manage containers.
type AwsTagInspector struct {
	// Ctr is the container to use as a base container.
	Ctr *dagger.Container
	// Cfg is the configuration file to use for the container.
	// +private
	Cfg *inspectorConfig
}

// New creates a new AwsTagInspector module.
//
// Parameters:
// - version: The version of the GoReleaser to use, e.g., "v1.22.0". Optional parameter.
// - image: The image to use as the base container. Optional parameter.
// - ctr: The container to use as a base container. Optional parameter.
// - envVarsFromHost: A list of environment variables to pass from the host to the container in a
// slice of strings. Optional parameter.
//
// Returns a pointer to a AwsTagInspector instance and an error, if any.
func New(
	// ctr is the container to use as a base container.
	// +optional
	ctr *dagger.Container,
	// awsAccessKeyID is the AWS access key ID to use for the container.
	awsAccessKeyID *dagger.Secret,
	// awsSecretAccessKey is the AWS secret access key to use for the container.
	awsSecretAccessKey *dagger.Secret,
	// configPath is the path to the configuration file to use for the container.
	config *dagger.File,
	// awsRegion is the AWS region to use for the container.
	// +optional
	awsRegion string,
	// envVarsFromHost is a list of environment variables to pass from the host to the container in a slice of strings.
	// +optional
	envVarsFromHost []string,
) (*AwsTagInspector, error) {
	//nolint:exhaustruct // It's 'okaysh' for now, I'll decide later what's going to be the pattern here.
	dagModule := &AwsTagInspector{}

	// Ensure awsRegion has a default value, but only if it's empty
	if awsRegion == "" {
		awsRegion = "us-east-1"
	}

	// Validate input parameters before processing
	if awsAccessKeyID == nil || awsSecretAccessKey == nil {
		return nil, Errorf("AWS access key ID and secret access key are required")
	}

	if config == nil {
		return nil, Errorf("configuration file is required")
	}

	// Append AWS region to environment variables
	envVarsFromHost = append(envVarsFromHost, fmt.Sprintf("AWS_REGION=%s", awsRegion))

	if err := dagModule.setupEnvironmentVariables(envVarsFromHost); err != nil {
		return nil, WrapError(err,
			"environment variable setup failed, unable to configure environment variables")
	}

	cfgLoader := newCfg()

	// Improve error handling and provide more context
	cfg, cfgErr := cfgLoader.loadConfig(context.Background(), config)
	if cfgErr != nil {
		return nil, WrapError(cfgErr, "failed to load configuration file")
	}

	if cfg == nil {
		return nil, Errorf("configuration file is empty")
	}

	dagModule.Cfg = cfg

	return dagModule, nil
}

// setupEnvironmentVariables sets up the environment variables.
//
// If the environment variables are not passed, it returns nil.
// If the environment variables are passed, it sets the environment variables.
func (m *AwsTagInspector) setupEnvironmentVariables(envVarsFromHost []string) error {
	if len(envVarsFromHost) == 0 {
		return nil
	}

	envVars, err := envvars.ToDaggerEnvVarsFromSlice(envVarsFromHost)
	if err != nil {
		return WrapError(err, "failed to parse environment variables")
	}

	for _, envVar := range envVars {
		m.WithEnvironmentVariable(envVar.Name, envVar.Value, false)
	}

	return nil
}
