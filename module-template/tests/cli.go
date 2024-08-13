package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/module-template/tests/internal/dagger"
)

// TestWithAWSCLIInAlpineContainer verifies that the AWS CLI is
// correctly installed
// inside an Alpine container, both in terms of binary location and version.
//
// It performs the following checks:
// 1. Verifies if the AWS CLI binary is installed in the OS
// filesystem using the `which aws` command.
// 2. Confirms the AWS CLI version by executing `aws --version`.
//
// If any of these checks fail, the function returns an error detailing the failure.
func (m *Tests) TestWithAWSCLIInAlpineContainer(ctx context.Context) error {
	targetModule := dag.ModuleTemplate()

	// Configure the target module to use AWS CLI in Alpine container
	targetModule = targetModule.WithAwscliinAlpineContainer()
	// Check if the AWS CLI binary is installed in the OS filesystem
	out, whichErr := targetModule.
		Ctr().
		WithExec([]string{"which", "aws"}).
		Stdout(ctx)

	if whichErr != nil {
		return WrapError(whichErr, "failed to get AWS CLI binary path")
	}

	if !strings.Contains(out, "/usr/bin/aws") {
		return Errorf("expected AWS CLI binary to be in "+
			"/usr/bin/aws, got %s", out)
	}

	// Check if the AWS CLI binary is installed and get its version
	outVersion, errVersion := targetModule.
		Ctr().
		WithExec([]string{"aws", "--version"}).
		Stdout(ctx)

	if errVersion != nil {
		return WrapError(errVersion, "failed to get AWS CLI version")
	}

	if !strings.Contains(outVersion, "aws-cli/2") {
		return Errorf("expected AWS CLI version to be 2, got %s", outVersion)
	}

	return nil
}

// TestWithAWSCLIInUbuntuContainer verifies that the AWS CLI is
// correctly installed inside an Ubuntu container, both in terms of
// binary location and version.
//
// It performs the following checks:
// 1. Verifies if the AWS CLI binary is installed in the OS
// filesystem using the `which aws` command.
// 2. Confirms the AWS CLI version by executing `aws --version`.
//
// If any of these checks fail, the function returns an error detailing the failure.
func (m *Tests) TestWithAWSCLIInUbuntuContainer(ctx context.Context) error {
	targetModule := dag.
		ModuleTemplate().
		BaseUbuntu()

	// Configure the target module to use AWS CLI in Ubuntu container
	targetModule = targetModule.
		WithAwscliinUbuntuContainer(dagger.
			ModuleTemplateWithAwscliinUbuntuContainerOpts{
			Architecture: "x86_64",
		})

	// Check if the AWS CLI binary is installed in the OS filesystem
	out, whichErr := targetModule.Ctr().
		WithExec([]string{"which", "aws"}).
		Stdout(ctx)

	if whichErr != nil {
		return WrapError(whichErr, "failed to get AWS CLI binary path")
	}

	if !strings.Contains(out, "/usr/local/bin/aws") {
		return Errorf("expected AWS CLI binary to "+
			"be in /usr/local/bin/aws, got %s", out)
	}

	// List on /usr/local/bin/aws
	lsOut, lsErr := targetModule.Ctr().
		WithExec([]string{"ls", "/usr/local/bin/aws"}).
		Stdout(ctx)

	if lsErr != nil {
		return WrapError(lsErr, "failed to list /usr/local/bin/aws")
	}

	if !strings.Contains(lsOut, "aws") {
		return Errorf("expected AWS CLI binary to be in /usr/local/bin/aws, got %s", lsOut)
	}

	return nil
}
