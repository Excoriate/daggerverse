package main

// WithGitInAlpineContainer installs Git in the golang/alpine container.
//
// It installs Git in the golang/alpine container.
func (m *Gotoolbox) WithGitInAlpineContainer() *Gotoolbox {
	m.Ctr = m.Ctr.
		WithExec([]string{"apk", "add", "git"})

	return m
}

// WithUtilitiesInAlpineContainer installs common utilities in the golang/alpine container.
//
// It installs utilities such as curl, wget, and others that are commonly used.
func (m *Gotoolbox) WithUtilitiesInAlpineContainer() *Gotoolbox {
	m.Ctr = m.Ctr.
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "add", "curl", "wget", "bash", "jq", "vim", "unzip", "yq"})

	return m
}
