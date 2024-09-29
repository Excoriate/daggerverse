package main

import (
	"strings"

	"github.com/Excoriate/daggerverse/module-template/internal/dagger"
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
// - *ModuleTemplate: A pointer to the updated ModuleTemplate instance.
func (m *ModuleTemplate) WithGoPlatform(
	// platform is the target platform specified in the format "[os]/[platform]/[version]".
	// For example, "darwin/arm64/v7", "windows/amd64", "linux/arm64". If the platform
	// is not provided, it defaults to "linux/amd64".
	// +optional
	platform dagger.Platform,
) *ModuleTemplate {
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
// - *ModuleTemplate: A pointer to the updated ModuleTemplate instance.
func (m *ModuleTemplate) WithGoCgoEnabled() *ModuleTemplate {
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
// - *ModuleTemplate: A pointer to the updated ModuleTemplate instance.
func (m *ModuleTemplate) WithCgoDisabled() *ModuleTemplate {
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
// - *ModuleTemplate: A pointer to the updated ModuleTemplate instance.
func (m *ModuleTemplate) WithGoBuildCache(
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
) *ModuleTemplate {
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
// - *ModuleTemplate: A pointer to the updated ModuleTemplate instance.
func (m *ModuleTemplate) WithGoModCache(
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
) *ModuleTemplate {
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
//	ModuleTemplate.WithGoInstall([]string{"github.com/Excoriate/daggerverse@latest", "github.com/another/package@v1.0"})
//
// Returns:
// - *ModuleTemplate: A pointer to the updated ModuleTemplate instance.
func (m *ModuleTemplate) WithGoInstall(
	// pkgs are the URLs of the packages to install.
	pkgs []string,
) *ModuleTemplate {
	// if len(pkgs) == 0 {
	// 	// Optionally, handle empty URLs here or return an error
	// 	// For now, we'll just exit early.
	// 	return m
	// }

	// Concatenate WithExec arguments based on the package URLs
	args := []string{"go", "install"}
	for _, pkg := range pkgs {
		m.Ctr = m.Ctr.
			WithExec(append(args, pkg))
	}

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
//	ModuleTemplate.WithGoExec([]string{"build", "./..."}, "linux/amd64")
//
// Returns:
// - *ModuleTemplate: A pointer to the updated ModuleTemplate instance.
func (m *ModuleTemplate) WithGoExec(
	// args are the arguments to pass to the Go command.
	args []string,
	// platform is the target platform specified in the format "[os]/[platform]/[version]".
	// For example, "darwin/arm64/v7", "windows/amd64", "linux/arm64". If the platform
	// is not provided, it defaults to "linux/amd64".
	// +optional
	platform dagger.Platform,
) *ModuleTemplate {
	if platform == "" {
		platform = "linux/amd64"
	}

	m.WithGoPlatform(platform)

	// Run Go commands, such as "go build" or "go test"
	args = append([]string{"go"}, args...)

	m.Ctr = m.Ctr.WithExec(args)

	return m
}

// goBuildCmdOptions represents the options for configuring a Go build command.
type goBuildCmdOptions struct {
	// Package is the Go package to compile. If empty, the package in the current directory is compiled.
	Package string
	// Race enables data race detection.
	Race bool
	// LDFlags are arguments to pass on each go tool link invocation.
	LDFlags []string
	// Tags are additional build tags to consider satisfied during the build.
	Tags []string
	// TrimPath removes all file system paths from the resulting executable.
	TrimPath bool
	// RawArgs are additional arguments to pass to the build command.
	RawArgs []string
	// Platform is the target platform in "[os]/[platform]/[version]" format, defaults to "linux/amd64".
	Platform dagger.Platform
	// Output specifies the output binary name or path.
	Output string
	// OptimizeSize optimizes for size.
	OptimizeSize bool
	// Verbose enables verbose output.
	Verbose bool
	// BuildMode specifies the build mode.
	BuildMode string
}

// WithGoBuild configures and executes a Go build command within the container environment.
//
// This method builds a Go package based on specified parameters such as package name, race detection,
// build flags, tags, and output binary name.
//
// Params:
// - pkg (string) +optional: The Go package to compile. If empty, the package in
// the current directory is compiled.
// - race (bool) +optional: Enable data race detection if true.
// - ldflags ([]string) +optional: Arguments to pass on each go tool link invocation.
// - tags ([]string) +optional: A list of additional build tags to consider satisfied
// during the build.
// - trimpath (bool) +optional: Remove all file system paths from the resulting executable if true.
// - rawArgs ([]string) +optional: Additional arguments to pass to the build command.
// - platform (dagger.Platform) +optional: Target platform in "[os]/[platform]/[version]" format,
// defaults to "linux/amd64".
// - output (string) +optional: The output binary name or path.
// - optimizeSize (bool) +optional: Optimize for size if true.
// - verbose (bool) +optional: Verbose output if true.
// - buildMode (string) +optional: Specify the build mode.
//
// Returns:
// - *ModuleTemplate: A pointer to the updated ModuleTemplate instance.
func (m *ModuleTemplate) WithGoBuild(
	// pkg is the Go package to compile.
	// +optional
	pkg string,

	// race enables data race detection.
	// +optional
	race bool,

	// ldflags are arguments to pass on each go tool link invocation.
	// +optional
	ldflags []string,

	// tags are additional build tags to consider satisfied during the build.
	// +optional
	tags []string,

	// trimpath removes all file system paths from the resulting executable.
	// +optional
	trimpath bool,

	// rawArgs are additional arguments to pass to the build command.
	// +optional
	rawArgs []string,

	// platform is the target platform in "[os]/[platform]/[version]" format, defaults to "linux/amd64".
	// +optional
	platform dagger.Platform,

	// output specifies the output binary name or path.
	// +optional
	output string,

	// optimizeSize optimizes for size.
	// +optional
	optimizeSize bool,

	// verbose enables verbose output.
	// +optional
	verbose bool,

	// buildMode specifies the build mode.
	// +optional
	buildMode string,
) *ModuleTemplate { // Set the target platform if provided
	m.WithGoPlatform(platform)

	// Build the package
	args := m.buildGoArgs(goBuildCmdOptions{
		Package:      pkg,
		Race:         race,
		LDFlags:      ldflags,
		Tags:         tags,
		TrimPath:     trimpath,
		RawArgs:      rawArgs,
		Platform:     platform,
		Output:       output,
		OptimizeSize: optimizeSize,
		Verbose:      verbose,
		BuildMode:    buildMode,
	})

	return m.WithGoExec(args, platform)
}

// buildGoArgs constructs the arguments for the Go build command based on the provided goBuildCmdOptions.
//
// Parameters:
//   - opts: goBuildCmdOptions struct containing all build configuration options.
//
// Returns:
//   - []string: A slice of strings representing the arguments to pass to the Go build command.
func (m *ModuleTemplate) buildGoArgs(opts goBuildCmdOptions) []string {
	args := []string{"build"}

	if opts.Race {
		args = append(args, "-race")
	}

	if opts.Output != "" {
		args = append(args, "-o", opts.Output)
	}

	if opts.OptimizeSize {
		args = append(args, "-ldflags", "-w -s")
	} else if len(opts.LDFlags) > 0 {
		args = append(args, "-ldflags", strings.Join(opts.LDFlags, " "))
	}

	if len(opts.Tags) > 0 {
		args = append(args, "-tags", strings.Join(opts.Tags, ","))
	}

	if opts.TrimPath {
		args = append(args, "-trimpath")
	}

	if opts.Verbose {
		args = append(args, "-v")
	}

	if opts.BuildMode != "" {
		args = append(args, "-buildmode", opts.BuildMode)
	}

	args = append(args, opts.RawArgs...)

	if opts.Package != "" {
		args = append(args, opts.Package)
	}

	return args
}

// WithGoPrivate configures the GOPRIVATE environment variable for the container environment.
//
// This method sets the GOPRIVATE variable, which is used by the Go toolchain to identify
// private modules or repositories that should not be fetched from public proxies.
//
// Parameters:
// - privateHost (string): The hostname of the private host for the Go packages or modules.
//
// Returns:
// - *ModuleTemplate: A pointer to the updated ModuleTemplate instance.
func (m *ModuleTemplate) WithGoPrivate(
	// privateHost is the hostname of the private host for the Go packages or modules.
	privateHost string,
) *ModuleTemplate {
	// Set the GOPRIVATE environment variable within the container
	m.Ctr = m.Ctr.
		WithExec(
			[]string{"go", "env", "GOPRIVATE", privateHost},
			dagger.ContainerWithExecOpts{
				InsecureRootCapabilities: true,
			},
		).
		WithEnvVariable("GOPRIVATE", privateHost)

	return m
}

// WithGoGCCCompiler installs the GCC compiler and development tools in the container environment.
//
// This method uses the Alpine package manager (`apk`) to install the GCC compiler along
// with `musl-dev`, which is the development package of the standard C library on Alpine.
//
// Example usage:
//
//	ModuleTemplate.WithGCCCompiler()
//
// Returns:
// - *ModuleTemplate: A pointer to the updated ModuleTemplate instance.
func (m *ModuleTemplate) WithGoGCCCompiler() *ModuleTemplate {
	m.Ctr = m.Ctr.
		WithExec([]string{"apk", "add", "--no-cache", "gcc", "musl-dev"})

	return m
}

// WithGoTestSum installs the GoTestSum tool and optionally its dependency `tparse` in the container environment.
//
// This method installs `gotest.tools/gotestsum` using the specified version, and optionally installs
// `github.com/mfridman/tparse` using a specified version, unless the skipTParse flag is set to true.
//
// Parameters:
//   - goTestSumVersion (string) +optional: The version of GoTestSum to use, e.g., "v0.8.0".
//     If empty, it defaults to "latest".
//   - tParseVersion (string) +optional: The version of TParse to use, e.g., "v0.8.0".
//     If empty, it defaults to the same version as goTestSumVersion.
//   - skipTParse (bool) +optional: If true, TParse will not be installed. Default is false.
//
// Example:
//
//	m := &ModuleTemplate{}
//	m.WithGoTestSum("v0.8.0", "v0.7.0", false)   // Install specific versions
//	m.WithGoTestSum("", "", true)                // Install latest version of GoTestSum and skip TParse
//
// Returns:
// - *ModuleTemplate: A pointer to the updated ModuleTemplate instance.
func (m *ModuleTemplate) WithGoTestSum(
	// goTestSumVersion is the version of GoTestSum to use, e.g., "v0.8.0".
	// +optional
	goTestSumVersion string,

	// tParseVersion is the version of TParse to use, e.g., "v0.8.0".
	// +optional
	tParseVersion string,

	// skipTParse is a flag to indicate whether TParse should be skipped.
	// +optional
	skipTParse bool,
) *ModuleTemplate {
	if goTestSumVersion == "" {
		goTestSumVersion = "latest"
	}

	if tParseVersion == "" {
		tParseVersion = "latest"
	}

	pkgs := []string{"gotest.tools/gotestsum@" + goTestSumVersion}
	if !skipTParse {
		pkgs = append(pkgs, "github.com/mfridman/tparse@"+tParseVersion)
	}

	return m.WithGoInstall(pkgs)
}

// WithGoReleaser installs the GoReleaser tool in the container environment.
//
// This method installs `goreleaser` using the specified version. Check
// https://goreleaser.com/install/
//
// Parameters:
// - version (string): The version of GoReleaser to use, e.g., "v2@latest".
//
// Example:
//
//	m := &ModuleTemplate{}
//	m.WithGoReleaser("v2.2.0")  // Install GoReleaser version 2.2.0
//
// Returns:
// - *ModuleTemplate: A pointer to the updated ModuleTemplate instance.
func (m *ModuleTemplate) WithGoReleaser(
	// version is the version of GoReleaser to use, e.g., By default, it's set "v2@latest".
	// +optional
	version string,
) *ModuleTemplate {
	if version == "" {
		version = "v2@latest"
	}

	return m.
		WithGoInstall([]string{"github.com/goreleaser/goreleaser/" + version})
}

// WithGoLint installs the GoLint tool in the container environment.
//
// This method installs `golangci-lint` using the specified version.
//
// Parameters:
// - version (string): The version of GoLint to use, e.g., "v1.60.1".
//
// Example:
//
//	m := &ModuleTemplate{}
//	m.WithGoLint("v1.60.1")  // Install GoLint version 1.60.1
//
// Returns:
// - *ModuleTemplate: A pointer to the updated ModuleTemplate instance.
func (m *ModuleTemplate) WithGoLint(
	// version is the version of GoLint to use, e.g., "v1.60.1".
	version string,
) *ModuleTemplate {
	return m.WithGoInstall([]string{"github.com/golangci/golangci-lint/cmd/golangci-lint@" + version})
}
