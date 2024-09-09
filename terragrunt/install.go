package main

// WithGitInAlpineContainer installs Git in the golang/alpine container.
//
// It installs Git in the golang/alpine container.
func (m *Terragrunt) WithGitInAlpineContainer() *Terragrunt {
	m.Ctr = m.Ctr.
		WithExec([]string{"apk", "add", "git"})

	return m
}

// WithGitInUbuntuContainer installs Git in the Ubuntu-based container.
//
// This method installs Git in the Ubuntu-based container.
//
// Returns:
//   - *Terragrunt: The updated Terragrunt with Git installed in the container.
func (m *Terragrunt) WithGitInUbuntuContainer() *Terragrunt {
	m.Ctr = m.Ctr.
		WithExec([]string{"apt-get", "update", "-y"}).
		WithExec([]string{"apt-get", "install", "-y", "git"})

	return m
}

// WithUtilitiesInAlpineContainer installs common utilities in the golang/alpine container.
//
// It installs utilities such as curl, wget, and others that are commonly used.
func (m *Terragrunt) WithUtilitiesInAlpineContainer() *Terragrunt {
	m.Ctr = m.Ctr.
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "add", "curl", "wget", "bash", "jq", "vim", "unzip", "yq"})

	return m
}

// WithUtilitiesInUbuntuContainer installs common utilities in the Ubuntu-based container.
//
// This method updates the package lists for upgrades and installs the specified utilities
// such as curl, wget, bash, jq, and vim in the Ubuntu-based container.
//
// Returns:
//   - *Terragrunt: The updated Terragrunt with the utilities installed in the container.
func (m *Terragrunt) WithUtilitiesInUbuntuContainer() *Terragrunt {
	m.Ctr = m.Ctr.
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "curl", "wget", "bash", "jq", "vim", "unzip", "yq"})

	return m
}
