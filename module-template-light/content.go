package main

import (
	"context"
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

// CloneGitRepoHTTPS clones a Git repository into a Dagger Directory.
//
// Parameters:
//   - repoURL: The URL of the git repository to clone (e.g., "https://github.com/user/repo").
//   - token: (Optional) The VCS token to use for authentication. If
//     not provided, the repository will be cloned without authentication.
//   - vcs: (Optional) The version control system (VCS) to use for
//     authentication. Defaults to "github". Supported values are "github" and "gitlab".
//   - authHeader: (Optional) The authentication header to use for authentication. If
//     not provided, the repository will be cloned without authentication.
//   - returnDir: (Optional) A string that indicates the directory path of the repository to return.
//   - branch: (Optional) The branch to checkout. If not provided, the default branch will be checked out.
//   - keepGitDir: (Optional) A boolean that indicates if the .git directory should be kept. Defaults to false.
//   - tag: (Optional) The tag to checkout. If not provided, no tag will be checked out.
//   - commit: (Optional) The commit to checkout. If not provided, the latest commit will be checked out.
//
// Returns:
//   - *dagger.Directory: A directory object representing the cloned repository.
//
// If a token is provided, it will be securely set using Dagger's
// secret mechanism and used for authentication during the clone operation.
func (m *ModuleTemplateLight) CloneGitRepoHTTPS(
	// repoURL is the URL of the git repo to clone.
	repoURL string,
	// token is the VCS token to use for authentication. Optional parameter.
	// +optional
	token *dagger.Secret,
	// vcs is the VCS to use for authentication. Optional parameter.
	// +optional
	vcs string,
	// authHeader is the authentication header to use for authentication. Optional parameter.
	// +optional
	authHeader *dagger.Secret,
	// returnDir is a string that indicates the directory path of the repository to return.
	// +optional
	returnDir string,
	// branch is the branch to checkout. Optional parameter.
	// +optional
	branch string,
	// discardGitDir is a boolean that indicates if the .git directory should be discarded. Optional parameter.
	// +optional
	discardGitDir bool,
	// tag is the tag to checkout. Optional parameter.
	// +optional
	tag string,
	// commit is the commit to checkout. Optional parameter.
	// +optional
	commit string,
) *dagger.Directory {
	// Default to GitHub if no VCS is specified.
	if vcs == "" {
		vcs = "github"
	}

	gitCloneOpts := dagger.GitOpts{}

	// If discardGitDir is true, set KeepGitDir to false. This changed with 0.13.4
	if discardGitDir {
		gitCloneOpts.KeepGitDir = false
	}

	// Initialize the Git clone request.
	gitCloneRequest := dag.Git(repoURL, gitCloneOpts)

	// If a token is provided, set it as a secret and attach it to the clone request.
	if token != nil {
		gitCloneRequest = gitCloneRequest.
			WithAuthToken(dag.
				SetSecret(m.getTokenNameByVcs(vcs), m.getTokenValueBySecret(token)))
	}

	if authHeader != nil {
		gitCloneRequest = gitCloneRequest.
			WithAuthHeader(authHeader)
	}

	var repoContent *dagger.Directory

	switch {
	case commit != "":
		repoContent = gitCloneRequest.Commit(commit).Tree()
	case tag != "":
		repoContent = gitCloneRequest.Tag(tag).Tree()
	case branch != "":
		repoContent = gitCloneRequest.Branch(branch).Tree()
	default:
		repoContent = gitCloneRequest.Head().Tree()
	}

	if returnDir != "" {
		return repoContent.
			Directory(returnDir)
	}

	return repoContent
}

// CloneGitRepoSSH clones a git repository using SSH for authentication.
// It allows specifying various options such as the branch, tag, or commit to checkout,
// whether to keep the .git directory, and the directory path of the repository to return.
//
// Parameters:
//
//	repoURL: The URL of the git repo to clone.
//	sshSocket: The SSH socket to use for authentication.
//	sshKnownHosts: The known hosts to use for authentication. Optional parameter.
//	returnDir: A string that indicates the directory path of the repository to return. Optional parameter.
//	branch: The branch to checkout. Optional parameter.
//	keepGitDir: A boolean that indicates if the .git directory should be kept. Optional parameter.
//	tag: The tag to checkout. Optional parameter.
//	commit: The commit to checkout. Optional parameter.
//
// Returns:
//
//	*dagger.Directory: The directory of the cloned repository.
func (m *ModuleTemplateLight) CloneGitRepoSSH(
	// repoURL is the URL of the git repo to clone.
	repoURL string,
	// sshAuthSocket is the SSH socket to use for authentication.
	sshAuthSocket *dagger.Socket,
	// sshKnownHosts is the known hosts to use for authentication. Optional parameter.
	// +optional
	sshKnownHosts string,
	// returnDir is a string that indicates the directory path of the repository to return.
	// +optional
	returnDir string,
	// branch is the branch to checkout. Optional parameter.
	// +optional
	branch string,
	// discardGitDir is a boolean that indicates if the .git directory should be discarded. Optional parameter.
	// +optional
	discardGitDir bool,
	// tag is the tag to checkout. Optional parameter.
	// +optional
	tag string,
	// commit is the commit to checkout. Optional parameter.
	// +optional
	commit string,
) *dagger.Directory {
	gitCloneOpts := dagger.GitOpts{
		SSHAuthSocket: sshAuthSocket,
	}

	// If discardGitDir is true, set KeepGitDir to false. This changed with 0.13.4
	if discardGitDir {
		gitCloneOpts.KeepGitDir = false
	}

	if sshKnownHosts != "" {
		gitCloneOpts.SSHKnownHosts = sshKnownHosts
	}

	gitCloneRequest := dag.Git(repoURL, gitCloneOpts)

	var repoContent *dagger.Directory

	switch {
	case commit != "":
		repoContent = gitCloneRequest.Commit(commit).Tree()
	case tag != "":
		repoContent = gitCloneRequest.Tag(tag).Tree()
	case branch != "":
		repoContent = gitCloneRequest.Branch(branch).Tree()
	default:
		repoContent = gitCloneRequest.Head().Tree()
	}

	if returnDir != "" {
		return repoContent.Directory(returnDir)
	}

	return repoContent
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
	plainTxtToken, err := secret.Plaintext(context.TODO())
	if err != nil {
		return ""
	}

	return plainTxtToken
}
