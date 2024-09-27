package main

import "github.com/Excoriate/daggerx/pkg/installerx"

// WithAWSCLIInAlpineContainer installs the AWS CLI in the
// Alpine-based container.
// This method installs the AWS CLI in a golang/alpine container
// using the 'apk' package manager.
// It is particularly useful for environments that need to
// interact with AWS services.
//
// Returns:
//   - *ModuleTemplate: The updated ModuleTemplate with the AWS CLI installed in the container.
func (m *ModuleTemplate) WithAWSCLIInAlpineContainer() *ModuleTemplate {
	m.Ctr = m.Ctr.
		WithExec([]string{"apk", "add", "--no-cache", "aws-cli"})

	return m
}

// WithAWSCLIInUbuntuContainer installs the AWS CLI in the Ubuntu-based container.
//
// This method installs the AWS CLI in an Ubuntu-based container following the
// official AWS installation steps.
//
// Args:
//   - architecture (string): The architecture for which the AWS CLI should be downloaded.
//     Valid values are "x86_64" and "aarch64". Default is "x86_64".
//
// Returns:
//   - *ModuleTemplate: The updated ModuleTemplate with the AWS CLI installed in the container.
func (m *ModuleTemplate) WithAWSCLIInUbuntuContainer(
	// architecture is the architecture for which the AWS CLI should be downloaded.
	// Valid values are "x86_64" and "aarch64". Default is "x86_64".
	// +optional
	architecture string) *ModuleTemplate {
	awsCLIInstallation := installerx.GetAwsCliInstallCommand(architecture)

	m.Ctr = m.Ctr.
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "unzip", "curl", "sudo"}).
		WithExec([]string{"sh", "-c", awsCLIInstallation})
		// WithExec([]string{"curl", "-L", url, "-o", "awscliv2.zip"}).
		// WithExec([]string{"unzip", "awscliv2.zip"}).
		// WithExec([]string{"sudo", "./aws/install"}).
		// WithExec([]string{"rm", "-rf", "awscliv2.zip", "aws"})

	return m
}
