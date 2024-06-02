package main

import (
	"context"
	"fmt"
	"strings"
)

// Terminal returns a terminal for the container.
//
// It returns a terminal for the container. It's meant to be used as a terminal for the module.
// Arguments:
// - None.
// Returns:
// - *Terminal: The terminal for the container.
func (m *Gopublisher) Terminal(
	// src is the directory that contains all the source code, including the module directory.
	// +optional
	src *Directory,
	// envVars is a set of strings (e.g., "KEY=value,KEY=value") to use as environment variables. They're
	// used to set the environment variables for the container when it's required to pass multiple environment variables
	// in a single argument. E.g.: "GITHUB_TOKEN=token,GO_VERSION=1.22.0,MYVAR=myvar"
	// +optional
	envVars []string,
) *Terminal {
	if len(envVars) > 0 {
		mMut, _ := m.WithEnvVarsFromStrs(envVars)
		m.Ctr = mMut.Ctr
	}

	if src != nil {
		m.Ctr = m.WithSource(src).Ctr
	}

	return m.Ctr.Terminal()
}

// GoModPath returns the module path.
//
// It returns the module path. It's meant to be used as a terminal for the module.
// Arguments:
// - src: The directory that contains all the source code, including the module directory.
// Returns:
// - string: The module path.
func (m *Gopublisher) GoModPath(
	// src is the directory that contains all the source code, including the module directory.
	src *Directory,
) (string, error) {
	m.Ctr = m.WithSource(src).Ctr

	m.Ctr = m.Ctr.
		WithExec([]string{"go", "list", "-m"})

	modulePath, err := m.Ctr.Stdout(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get the Go module path: %w", err)
	}

	modulePath = strings.TrimSpace(modulePath)

	return modulePath, nil
}

// GoModVersion returns the module version by running git describe --tags --abbrev=0.
//
// It returns the module version by running git describe --tags --abbrev=0. It's meant to be used as a terminal for the module.
// Arguments:
// - src: The directory that contains all the source code, including the module directory.
// Returns:
// - string: The module version.
func (m *Gopublisher) GoModVersion(
	// src is the directory that contains all the source code, including the module directory.
	src *Directory,
) (string, error) {
	m.Ctr = m.WithSource(src).Ctr
	version, err := m.Ctr.WithExec([]string{"git", "describe", "--tags", "--abbrev=0"}).Stdout(context.Background())

	if err != nil {
		return "", fmt.Errorf("failed to get the Go module version: %w", err)
	}

	version = strings.TrimSpace(version)
	return version, nil
}

// GoModPublish publishes the module to the registry.
//
// It publishes the module to the registry. It's meant to be used as a terminal for the module.
// Arguments:
// - src: The directory that contains all the source code, including the module directory.
// - tag: The tag to use for the release.
// Returns:
// - string: The tag that was used for the release.
// - error: An error that occurred while publishing the module.
func (m *Gopublisher) GoModPublish(
	// src is the directory that contains all the source code, including the module directory.
	src *Directory,
	// tag is the tag to use for the release.
	// +optional
	tag string,
	// githubToken is the GitHub token to use for the release. It's required when publishing from GitHub Actions.
	// +optional
	githubToken string,
	// envVars is a set of strings (e.g., "KEY=value,KEY=value") to use as environment variables. They're
	// used to set the environment variables for the container when it's required to pass multiple environment variables
	// in a single argument. E.g.: "GITHUB_TOKEN=token,GO_VERSION=1.22.0,MYVAR=myvar"
	// +optional
	envVars []string,
) (string, error) {
	m.Ctr = m.WithSource(src).Ctr

	if githubToken != "" {
		m.Ctr = m.WithEnvVariable("GITHUB_TOKEN", githubToken, false).Ctr
	}

	if len(envVars) > 0 {
		mMut, err := m.WithEnvVarsFromStrs(envVars)
		if err != nil {
			return "", fmt.Errorf("failed to set environment variables: %w", err)
		}

		m.Ctr = mMut.Ctr
	}

	_, err := m.Ctr.
		WithExec([]string{"go", "mod", "tidy"}).
		Stdout(context.Background())

	if err != nil {
		return "", fmt.Errorf("failed to tidy the Go module: %w", err)
	}

	modulePath, err := m.GoModPath(src)
	if err != nil {
		return "", fmt.Errorf("failed to publish the module: %w", err)
	}

	var moduleVersion string

	if tag != "" {
		m.Ctr = m.Ctr.WithExec([]string{"git", "tag", tag})
		moduleVersion = tag
	} else {
		moduleVersion, err = m.GoModVersion(src)
		if err != nil {
			return "", fmt.Errorf("failed to publish the module: %w", err)
		}
	}

	//# Push the module version to the Go module proxy
	m.Ctr = m.WithEnvVariable("GOPROXY", goProxyURL, false).Ctr
	goModWithVersion := fmt.Sprintf("%s@%s", modulePath, moduleVersion)
	indexPkgInProxy := fmt.Sprintf("%s/%s", pkgRegistryURL, modulePath)

	_, err = m.Ctr.
		WithExec([]string{"go", "list", "-m", goModWithVersion}).
		WithExec([]string{"curl", indexPkgInProxy}).
		Stdout(context.Background())

	if err != nil {
		return "", fmt.Errorf("failed to publish the module: %s with version %s: %w", modulePath, moduleVersion, err)
	}

	successMessage := fmt.Sprintf(`
🎉 Successfully published module!
-----------------------------------

- 🌐 GO PKG URL: %s
- 📚 GO Module: %s
- 📦 Version: %s

🔎 Inspect your package module at the provided URL above.`, indexPkgInProxy, modulePath, moduleVersion)

	return successMessage, nil
}
