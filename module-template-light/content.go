package main

import (
	"path/filepath"

	"github.com/Excoriate/daggerverse/module-template-light/internal/dagger"
)

// DownloadFile downloads a file from the specified URL.
//
// Parameters:
//   - url: The URL of the file to download.
//   - destFileName: The name of the file to download. Optional parameter.
//     If not set, it'll default to the basename of the URL.
//
// Returns:
//   - *dagger.File: The downloaded file.
//
// Functionality:
//
// This method downloads a file from the provided URL. If the destination file
// name is not specified, it defaults to the basename of the URL. The downloaded
// file is then returned as a *dagger.File.
func (m *ModuleTemplateLight) DownloadFile(
	// url is the URL of the file to download.
	url string,
	// destFileName is the name of the file to download. If not set, it'll default to the basename of the URL.
	// +optional
	destFileName string,
) *dagger.File {
	fileName := filepath.Base(url)
	if destFileName != "" {
		fileName = destFileName
	}

	fileDownloaded := dag.HTTP(url).WithName(fileName)

	return fileDownloaded
}

// CloneGitRepo clones a Git repository into a Dagger Directory.
//
// Parameters:
//   - repoURL: The URL of the git repository to clone (e.g., "https://github.com/user/repo").
//   - token: (Optional) The VCS token to use for authentication. If
//     not provided, the repository will be cloned without authentication.
//   - vcs: (Optional) The version control system (VCS) to use for
//     authentication. Defaults to "github". Supported values are "github" and "gitlab".
//
// Returns:
//   - *dagger.Directory: A directory object representing the cloned repository.
//
// If a token is provided, it will be securely set using Dagger's
// secret mechanism and used for authentication during the clone operation.
func (m *ModuleTemplateLight) CloneGitRepo(
	// repoURL is the URL of the git repo to clone.
	repoURL string,
	// token is the VCS token to use for authentication. Optional parameter.
	// +optional
	token *dagger.Secret,
	// vcs is the VCS to use for authentication. Optional parameter.
	// +optional
	vcs string,
) *dagger.Directory {
	// Default to GitHub if no VCS is specified.
	if vcs == "" {
		vcs = "github"
	}

	// Initialize the Git clone request.
	gitCloneRequest := dag.Git(repoURL)

	// If a token is provided, set it as a secret and attach it to the clone request.
	if token != nil {
		gitCloneRequest = gitCloneRequest.
			WithAuthToken(dag.
				SetSecret(m.getTokenNameByVcs(vcs), m.getTokenValueBySecret(token)))
	}

	// Perform the Git clone operation and return the resulting directory.
	return gitCloneRequest.
		Head().
		Tree()
}

// getTokenNameByVcs returns the appropriate token name based on the VCS type.
// If the VCS is "gitlab", it returns "GITLAB_TOKEN". Otherwise, it returns "GITHUB_TOKEN".
func (m *ModuleTemplateLight) getTokenNameByVcs(vcs string) string {
	if vcs == "gitlab" {
		return "GITLAB_TOKEN"
	}

	return "GITHUB_TOKEN"
}

// getTokenValueBySecret retrieves the plaintext value of the provided secret.
// If an error occurs while retrieving the plaintext, it returns an empty string.
func (m *ModuleTemplateLight) getTokenValueBySecret(secret *dagger.Secret) string {
	plainTxtToken, err := secret.Plaintext(nil)
	if err != nil {
		return ""
	}

	return plainTxtToken
}
