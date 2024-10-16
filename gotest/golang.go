package main

import (
	"github.com/Excoriate/daggerverse/gotest/internal/dagger"
	"github.com/containerd/containerd/platforms"
)

// WithGoPlatform sets the Go environment variables for the target platform.
//
// platform is the target platform specified in the format "[os]/[platform]/[version]".
// For example, "darwin/arm64/v7", "windows/amd64", "linux/arm64". If the platform
// is not provided, it defaults to "linux/amd64".
//
// Params:
// - platform (dagger.Platform): The target platform.
//
// Returns:
// - *Gotest: A pointer to the updated Gotest instance.
func (m *Gotest) WithGoPlatform(
	// platform is the target platform specified in the format "[os]/[platform]/[version]".
	// For example, "darwin/arm64/v7", "windows/amd64", "linux/arm64". If the platform
	// is not provided, it defaults to "linux/amd64".
	// +optional
	platform dagger.Platform,
) *Gotest {
	if platform == "" {
		platform = "linux/amd64"
	}

	//nolint:staticcheck // This is a constant string.
	parsedPlatform := platforms.MustParse(string(platform))

	m.Ctr = m.Ctr.
		WithEnvVariable("GOOS", parsedPlatform.OS).
		WithEnvVariable("GOARCH", parsedPlatform.Architecture).
		With(func(c *dagger.Container) *dagger.Container {
			if parsedPlatform.Variant != "" {
				return c.WithEnvVariable("GOARM", parsedPlatform.Variant)
			}

			return c
		})

	return m
}

// WithGoCgoEnabled enables CGO for the container environment.
//
// When CGO is enabled, the Go toolchain will allow the use of cgo, which
// means it can link against C libraries and send code to the C compiler.
//
// Returns:
// - *Gotest: A pointer to the updated Gotest instance.
func (m *Gotest) WithGoCgoEnabled() *Gotest {
	m.Ctr = m.Ctr.
		WithEnvVariable("CGO_ENABLED", "1")

	return m
}

// WithGCCCompilerInstalled installs the GCC compiler and musl-dev package
// in the container environment using the Alpine package manager (apk).
//
// This method is useful for enabling the Go toolchain to compile C code
// and link against C libraries, which is necessary for certain Go packages
// that rely on CGO.
//
// Returns:
// - *Gotest: A pointer to the updated Gotest instance.
func (m *Gotest) WithGCCCompilerInstalled() *Gotest {
	m.Ctr = m.Ctr.
		WithExec(
			[]string{"apk", "add", "--no-cache", "gcc", "musl-dev"},
		)

	return m
}

// WithCgoDisabled disables CGO for the container environment.
//
// When CGO is disabled, the Go toolchain will not permit the use of cgo,
// which means it will not link against C libraries or send code to the C compiler.
// This can be beneficial for creating fully static binaries or for environments
// where C dependencies are not available.
//
// Returns:
// - *Gotest: A pointer to the updated Gotest instance.
func (m *Gotest) WithCgoDisabled() *Gotest {
	m.Ctr = m.Ctr.WithEnvVariable("CGO_ENABLED", "0")

	return m
}

// WithGoBuildCache configures a Go build cache for the container environment.
//
// This method sets up the cache volume for Go build artifacts, which speeds up
// the build process by reusing previously compiled packages and dependencies.
//
// Params:
// - cacheRoot (string) +optional: The path to the cache volume's root. If not provided,
// it defaults to "/root/.cache/go-build".
// - cache (*dagger.CacheVolume) +optional: The cache volume to use. If not provided, a
// default volume named "gobuildcache" is used.
// - source (*dagger.Directory) +optional: The directory to use as the source for the cache
// volume's root.
// - sharing (dagger.CacheSharingMode) +optional: The sharing mode of the cache volume. If not
// provided, it defaults to "shared".
//
// Returns:
// - *Gotest: A pointer to the updated Gotest instance.
func (m *Gotest) WithGoBuildCache(
	// cacheRoot is the path to the cache volume's root.
	// +optional
	cacheRoot string,
	// cache is the cache volume to use.
	// +optional
	cache *dagger.CacheVolume,
	// source is the identifier of the directory to use as the cache volume's root.
	// +optional
	source *dagger.Directory,
	// sharing is the Sharing mode of the cache volume.
	// +optional
	sharing dagger.CacheSharingMode,
) *Gotest {
	if cacheRoot == "" {
		cacheRoot = "/root/.cache/go-build"
	}

	if cache == nil {
		cache = dag.CacheVolume("gobuildcache")
	}

	if sharing == "" {
		sharing = dagger.Shared
	}

	m.Ctr = m.Ctr.WithMountedCache(
		cacheRoot,
		cache,
		dagger.ContainerWithMountedCacheOpts{
			Source:  source,
			Sharing: sharing,
		},
	)

	// Set the GOCACHE environment variable to the cache volume's root path.
	m.Ctr = m.Ctr.
		WithEnvVariable("GOCACHE", cacheRoot)

	return m
}

// WithGoModCache configures a Go module cache for the container environment.
//
// This method sets up the cache volume for Go modules, which speeds up
// the build process by reusing previously downloaded dependencies.
//
// Params:
// - cacheRoot (string): The path to the cache volume's root. If not provided, it defaults to "/go/pkg/mod".
// - cache (*dagger.CacheVolume) +optional: The cache volume to use. If not provided,
// a default volume named "gomodcache" is used.
// - source (*dagger.Directory) +optional: The directory to use as the
// source for the cache volume's root.
// - sharing (dagger.CacheSharingMode) +optional: The sharing mode of the cache volume. If
// not provided, it defaults to "shared".
//
// Returns:
// - *Gotest: A pointer to the updated Gotest instance.
func (m *Gotest) WithGoModCache(
	// cacheRoot is the path to the cache volume's root.
	// +optional
	cacheRoot string,
	// cache is the cache volume to use.
	// +optional
	cache *dagger.CacheVolume,
	// source is the identifier of the directory to use as the cache volume's root.
	// +optional
	source *dagger.Directory,
	// sharing is the Sharing mode of the cache volume.
	// +optional
	sharing dagger.CacheSharingMode,
) *Gotest {
	if cacheRoot == "" {
		cacheRoot = "/go/pkg/mod"
	}

	if cache == nil {
		cache = dag.CacheVolume("gomodcache")
	}

	if sharing == "" {
		sharing = dagger.Shared
	}

	// Set the GOCACHE environment variable to the cache volume's root path.
	m.Ctr = m.Ctr.
		WithEnvVariable("GOCACHE", cacheRoot)

	m.Ctr = m.Ctr.WithMountedCache(
		cacheRoot,
		cache,
		dagger.ContainerWithMountedCacheOpts{
			Source:  source,
			Sharing: sharing,
		},
	)

	return m
}

// WithGoInstall installs Go packages in the container environment.
//
// This method uses the Go toolchain to install packages from specified URLs.
// It performs the equivalent of running "go install [url]" inside the container.
//
// Params:
//   - pkgs ([]string): A slice of URLs for the Go packages to install. These should be
//     valid package URLs that the `go install` command can use.
//
// Example:
//
//	Gotest.WithGoInstall([]string{"github.com/Excoriate/daggerverse@latest",
//
// "github.com/another/package@v1.0"})
//
// Returns:
// - *Gotest: A pointer to the updated Gotest instance.
func (m *Gotest) WithGoInstall(
	// pkgs are the URLs of the packages to install.
	pkgs []string,
) *Gotest {
	if len(pkgs) == 0 {
		// Optionally, handle empty URLs here or return an error
		// For now, we'll just exit early.
		return m
	}

	// Concatenate WithExec arguments based on the package URLs
	args := []string{"go", "install"}
	for _, pkg := range pkgs {
		m.Ctr = m.Ctr.
			WithExec(append(args, pkg))
	}

	return m
}