package main

import "github.com/Excoriate/daggerx/pkg/fixtures"

// WithGoCache mounts the Go cache directories.
// The Go cache directories are:
// - /go/pkg/mod
// - /root/.cache/go-build
func (m *Goreleaser) WithGoCache() *Goreleaser {
	goModCache := dag.CacheVolume("gomodcache")
	goBuildCache := dag.CacheVolume("gobuildcache")

	ctr := m.Ctr.WithMountedCache("/go/pkg/mod", goModCache).
		WithMountedCache("/root/.cache/go-build", goBuildCache)

	m.Ctr = ctr

	return m
}

// WithSource sets the source directory.
func (m *Goreleaser) WithSource(src *Directory) *Goreleaser {
	m.Src = src
	m.Ctr = m.Ctr.WithWorkdir(fixtures.MntPrefix).
		WithMountedDirectory(fixtures.MntPrefix, src)

	return m
}

// WithCfgFile sets the configuration file to use.
// The default configuration file is ".goreleaser.yaml".
func (m *Goreleaser) WithCfgFile(
	// cfgFile is the configuration file to use.
	// +optional
	cfgFile string,
) *Goreleaser {
	m.CfgFile = setToDefaultCfgIfEmpty(cfgFile)

	return m
}

// WithCMD sets the command to run.
func (m *Goreleaser) WithCMD(cmd []string) *Goreleaser {
	m.Ctr = m.Ctr.
		WithFocus().
		WithExec(cmd)
	return m
}

// WithEnvVar sets an environment variable.
func (m *Goreleaser) WithEnvVar(key, value string, expand bool) *Goreleaser {
	m.Ctr = m.Ctr.WithEnvVariable(key, value, ContainerWithEnvVariableOpts{
		Expand: expand,
	})

	return m
}

// WithSecret sets a Secret as environment variable.
func (m *Goreleaser) WithSecret(name string, secret *Secret) *Goreleaser {
	m.Ctr = m.Ctr.WithSecretVariable(name, secret)
	return m
}
