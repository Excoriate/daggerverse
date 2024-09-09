package main

import (
	"context"
	"fmt"

	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
)

// WithNewNetrcFileGitHub creates a new .netrc file with the GitHub credentials.
//
// The .netrc file is created in the root directory of the container.
func (m *Terragrunt) WithNewNetrcFileGitHub(
	username string,
	password string,
) *Terragrunt {
	machineCMD := "machine github.com\nlogin " + username + "\npassword " + password + "\n"

	m.Ctr = m.Ctr.WithNewFile(netRcRootPath, machineCMD)

	return m
}

// WithNewNetrcFileAsSecretGitHub creates a new .netrc file with the GitHub credentials.
//
// The .netrc file is created in the root directory of the container.
// The argument 'password' is a secret that is not exposed in the logs.
func (m *Terragrunt) WithNewNetrcFileAsSecretGitHub(username string, password *dagger.Secret) *Terragrunt {
	passwordTxtValue, _ := password.Plaintext(context.Background())
	machineCMD := fmt.Sprintf("machine github.com\nlogin %s\npassword %s\n", username, passwordTxtValue)
	//nolint:exhaustruct // This is a method that is used to set the base image and version.
	m.Ctr = m.Ctr.WithNewFile(netRcRootPath, machineCMD)

	return m
}

// WithNewNetrcFileGitLab creates a new .netrc file with the GitLab credentials.
//
// The .netrc file is created in the root directory of the container.
func (m *Terragrunt) WithNewNetrcFileGitLab(
	username string,
	password string,
) *Terragrunt {
	machineCMD := "machine gitlab.com\nlogin " + username + "\npassword " + password + "\n"

	m.Ctr = m.Ctr.WithNewFile(netRcRootPath, machineCMD)

	return m
}

// WithNewNetrcFileAsSecretGitLab creates a new .netrc file with the GitLab credentials.
//
// The .netrc file is created in the root directory of the container.
// The argument 'password' is a secret that is not exposed in the logs.
func (m *Terragrunt) WithNewNetrcFileAsSecretGitLab(username string, password *dagger.Secret) *Terragrunt {
	passwordTxtValue, _ := password.Plaintext(context.Background())
	machineCMD := fmt.Sprintf("machine gitlab.com\nlogin %s\npassword %s\n", username, passwordTxtValue)

	//nolint:exhaustruct // This is a method that is used to set the base image and version.
	m.Ctr = m.Ctr.WithNewFile(netRcRootPath, machineCMD)

	return m
}
