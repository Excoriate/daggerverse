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
	"path/filepath"

	"github.com/Excoriate/daggerverse/aws-tag-inspector/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/apkox"
	"github.com/Excoriate/daggerx/pkg/fixtures"
)

// ApkoKeyRingInfo represents the keyring information for Apko.
type ApkoKeyRingInfo apkox.KeyringInfo

const (
	defaultApkoImage   = "cgr.dev/chainguard/apko"
	defaultApkoTarball = "image.tar"
)

// AwsTagInspector is a Dagger module.
//
// This module is used to create and manage containers.
type AwsTagInspector struct {
	// Cfg is the configuration file to use for the container.
	// +private
	Cfg *inspectorConfig
	// AWSClient is the AWS client to use for the container.
	// +private
	AWSClient *AWSClient
	// Ctr is the main container to use for the module.
	// +private
	Ctr *dagger.Container
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
	// ctx is the context for the new function
	// +optional
	ctx context.Context,
	// awsAccessKeyID is the AWS access key ID to use for the container.
	awsAccessKeyID *dagger.Secret,
	// awsSecretAccessKey is the AWS secret access key to use for the container.
	awsSecretAccessKey *dagger.Secret,
	// configPath is the path to the configuration file to use for the container.
	// +optional
	config *dagger.File,
	// awsRegion is the AWS region to use for the container.
	// +optional
	awsRegion string,
	// envVarsFromHost is a list of environment variables to pass from the host to the container in a slice of strings.
	// +optional
	envVarsFromHost []string,
) (*AwsTagInspector, error) {
	//nolint:exhaustruct // It's 'okaysh' for now, I'll decide later what's going to be the pattern here.
	dagModule := &AwsTagInspector{
		Cfg: &inspectorConfig{},
	}

	awsClient, awsClientErr := dagModule.setupAWSCredentials(ctx, awsAccessKeyID, awsSecretAccessKey, awsRegion)
	if awsClientErr != nil {
		return nil, awsClientErr
	}

	dagModule.AWSClient = awsClient

	// Only process configuration if a config file is provided
	if config != nil {
		cfgLoader := newCfg()

		// Improve error handling and provide more context
		cfg, cfgErr := cfgLoader.loadConfig(ctx, config)
		if cfgErr != nil {
			return nil, WrapError(cfgErr, "failed to load configuration file")
		}

		if cfg == nil {
			return nil, Errorf("configuration file is empty")
		}

		dagModule.Cfg = cfg
	}

	dagModule.BaseContainer()

	dagModule.WithSource(dag.CurrentModule().
		Source(),
		fixtures.MntPrefix,
	)

	return dagModule, nil
}

// setupAWSCredentials validates and sets up AWS credentials for the module.
//
// Parameters:
// - ctx: The context for the operation.
// - awsAccessKeyID: The AWS access key ID secret.
// - awsSecretAccessKey: The AWS secret access key secret.
// - awsRegion: The AWS region to use (defaults to us-east-1 if empty).
//
// Returns:
// - *AWSClient: A configured AWS client.
// - error: Any error encountered during credential setup.
func (m *AwsTagInspector) setupAWSCredentials(
	// ctx is the context for the setupAWSCredentials function
	ctx context.Context,
	// awsAccessKeyID is the AWS access key ID secret
	awsAccessKeyID *dagger.Secret,
	// awsSecretAccessKey is the AWS secret access key secret
	awsSecretAccessKey *dagger.Secret,
	// awsRegion is the AWS region to use (defaults to us-east-1 if empty)
	awsRegion string,
) (*AWSClient, error) {
	// Ensure awsRegion has a default value, but only if it's empty
	if awsRegion == "" {
		awsRegion = "us-east-1"
	}

	// Validate input parameters before processing
	if awsAccessKeyID == nil || awsSecretAccessKey == nil {
		return nil, Errorf("AWS access key ID and secret access key are required")
	}

	accessKeyValue, accessKeyErr := awsAccessKeyID.Plaintext(ctx)
	if accessKeyErr != nil {
		return nil, WrapError(accessKeyErr, "failed to get AWS access key ID")
	}

	secretAccessKeyValue, secretAccessKeyErr := awsSecretAccessKey.Plaintext(ctx)
	if secretAccessKeyErr != nil {
		return nil, WrapError(secretAccessKeyErr, "failed to get AWS secret access key")
	}

	awsClient, awsClientErr := NewAWSClient(ctx, AWSClientConfig{
		Region:          awsRegion,
		AccessKeyID:     accessKeyValue,
		SecretAccessKey: secretAccessKeyValue,
	})

	if awsClientErr != nil {
		return nil, WrapError(awsClientErr, "failed to create AWS client")
	}

	return awsClient, nil
}

// BaseContainer sets the base image to an Apko image and creates the base container.
//
// Returns a pointer to the Gopkgpublisher instance.
func (m *AwsTagInspector) BaseContainer() (*AwsTagInspector, error) {
	apkoCfgFilePath := "config/presets/base-alpine.yaml"
	apkoCfgFile := dag.CurrentModule().
		Source().
		File(apkoCfgFilePath)

	apkoCfgFilePathMounted := filepath.Join(fixtures.MntPrefix, apkoCfgFilePath)

	apkoCtr := dag.Container().
		From(defaultApkoImage).
		WithMountedFile(apkoCfgFilePathMounted, apkoCfgFile)

	apkoBuildCmd := []string{
		"apko",
		"build",
		apkoCfgFilePathMounted,
		"latest",
		defaultApkoTarball,
		"--cache-dir",
		"/var/cache/apko",
	}

	apkoCtr = apkoCtr.
		WithExec(apkoBuildCmd)

	outputTar := apkoCtr.
		File(defaultApkoTarball)

	m.Ctr = dag.
		Container().
		Import(outputTar)

	return m, nil
}
