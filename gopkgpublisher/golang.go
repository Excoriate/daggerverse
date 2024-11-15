package main

import (
	"github.com/Excoriate/daggerverse/gopkgpublisher/internal/dagger"
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
// - *Gopkgpublisher: A pointer to the updated Gopkgpublisher instance.
func (m *Gopkgpublisher) WithGoPlatform(
	// platform is the target platform specified in the format "[os]/[platform]/[version]".
	// For example, "darwin/arm64/v7", "windows/amd64", "linux/arm64". If the platform
	// is not provided, it defaults to "linux/amd64".
	// +optional
	platform dagger.Platform,
) *Gopkgpublisher {
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

// WithGoGCCCompilerAlpine installs the GCC compiler and musl-dev package
// in the container environment using the Alpine package manager (apk).
//
// This method is useful for enabling the Go toolchain to compile C code
// and link against C libraries, which is necessary for certain Go packages
// that rely on CGO.
//
// Returns:
// - *Gopkgpublisher: A pointer to the updated Gopkgpublisher instance.
func (m *Gopkgpublisher) WithGoGCCCompilerAlpine() *Gopkgpublisher {
	m.Ctr = m.Ctr.
		WithExec(
			[]string{"apk", "add", "--no-cache", "gcc", "musl-dev"},
		)

	return m
}

// WithGoGCCCompilerUbuntu installs the GCC compiler and musl-dev package
// in the container environment using the Ubuntu package manager (apt-get).
//
// This method is useful for enabling the Go toolchain to compile C code
// and link against C libraries, which is necessary for certain Go packages
// that rely on CGO.
//
// Returns:
// - *Gopkgpublisher: A pointer to the updated Gopkgpublisher instance.
func (m *Gopkgpublisher) WithGoGCCCompilerUbuntu() *Gopkgpublisher {
	m.Ctr = m.Ctr.
		WithExec(
			[]string{"apt-get", "update", "-y"},
		).
		WithExec(
			[]string{"apt-get", "install", "-y", "gcc", "musl-dev"},
		)

	return m
}

// WithGoCgoEnabled enables CGO for the container environment.
//
// When CGO is enabled, the Go toolchain will allow the use of cgo, which
// means it can link against C libraries and send code to the C compiler.
//
// Returns:
// - *Gopkgpublisher: A pointer to the updated Gopkgpublisher instance.
func (m *Gopkgpublisher) WithGoCgoEnabled() *Gopkgpublisher {
	m.Ctr = m.Ctr.WithEnvVariable("CGO_ENABLED", "1")

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
// - *Gopkgpublisher: A pointer to the updated Gopkgpublisher instance.
func (m *Gopkgpublisher) WithCgoDisabled() *Gopkgpublisher {
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
// - *Gopkgpublisher: A pointer to the updated Gopkgpublisher instance.
func (m *Gopkgpublisher) WithGoBuildCache(
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
) *Gopkgpublisher {
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
// - *Gopkgpublisher: A pointer to the updated Gopkgpublisher instance.
func (m *Gopkgpublisher) WithGoModCache(
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
) *Gopkgpublisher {
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

// WithGoExec runs a Go command in the container environment.
//
// This method allows you to execute arbitrary Go commands, such as "go build" or "go test",
// for a specified platform. If no platform is provided, it defaults to "linux/amd64".
//
// Params:
//   - args ([]string): The arguments to pass to the Go command. For example, ["build", "./..."]
//     to run a Go build command on all packages.
//   - platform (dagger.Platform) +optional: The target platform specified in the format
//     "[os]/[platform]/[version]". For example, "darwin/arm64/v7", "windows/amd64", "linux/arm64".
//     If the platform is not provided, it defaults to "linux/amd64".
//
// Example:
//
//	Gopkgpublisher.WithGoExec([]string{"build", "./..."}, "linux/amd64")
//
// Returns:
// - *Gopkgpublisher: A pointer to the updated Gopkgpublisher instance.
func (m *Gopkgpublisher) WithGoExec(
	// args are the arguments to pass to the Go command.
	args []string,
	// platform is the target platform specified in the format "[os]/[platform]/[version]".
	// For example, "darwin/arm64/v7", "windows/amd64", "linux/arm64". If the platform
	// is not provided, it defaults to "linux/amd64".
	// +optional
	platform dagger.Platform,
) *Gopkgpublisher {
	if platform == "" {
		platform = "linux/amd64"
	}

	m.WithGoPlatform(platform)

	// Run Go commands, such as "go build" or "go test"
	args = append([]string{"go"}, args...)

	m.Ctr = m.Ctr.WithExec(args)

	return m
}
