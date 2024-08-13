package main

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
// This method installs the AWS CLI in an Ubuntu-based
// container using the 'apt-get' package manager.
// It is particularly useful for environments that need to
// interact with AWS services.
//
// Returns:
//   - *ModuleTemplate: The updated ModuleTemplate with the AWS CLI installed in the container.
func (m *ModuleTemplate) WithAWSCLIInUbuntuContainer() *ModuleTemplate {
	m.Ctr = m.Ctr.
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "awscli"})

	return m
}
