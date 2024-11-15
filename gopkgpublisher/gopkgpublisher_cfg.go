package main

import "github.com/Excoriate/daggerverse/gopkgpublisher/internal/dagger"

const (
	ctrUser  = "gopkgpublisher"
	ctrGroup = "gopkgpublisher"
)

// WithGoPackage adds the specified Go package version to the Gopkgpublisher's list of ApkoPackages.
// If the version string is empty, it appends the default "go" package to the list.
// If a version is provided, it appends "go=<version>" to the list, where <version> is the specified version string.
// This method returns the updated Gopkgpublisher instance, allowing for method chaining.
// WithGoPackages adds a Go package version to the Gopkgpublisher's list of ApkoPackages.
//
// If the provided version string is empty, it appends the default "go" package to the list.
// If a version is specified, it appends "go=<version>" to the list, where <version> is the specified version string.
//
// Parameters:
//
//	version (string): The version of the Go package to add. If empty, the default "go" package is added.
//
// Returns:
//
//	*Gopkgpublisher: The updated Gopkgpublisher instance, allowing for method chaining.
func (m *Gopkgpublisher) WithGoPackage(version string) *Gopkgpublisher {
	if version == "" {
		m.ApkoPackages = append(m.ApkoPackages, "go="+defaultGoVersion)
	} else {
		m.ApkoPackages = append(m.ApkoPackages, "go="+version)
	}

	return m
}

// WithExtraPackages adds extra packages to the APKO packages list.
// This function allows adding multiple packages at once.
//
// Parameters:
//
//	packages - a variadic parameter representing the list of packages to add.
//
// Returns:
//
//	A pointer to the updated Gopkgpublisher instance.
func (m *Gopkgpublisher) WithExtraPackages(packages ...string) *Gopkgpublisher {
	m.ApkoPackages = append(m.ApkoPackages, packages...)

	return m
}

// WithCacheConfiguration configures all necessary Go caches in the paths
// defined in the APKO configuration, using the gopkgpublisher user.
// WithCacheConfiguration configures the necessary Go cache directories for the Gopkgpublisher instance.
// It sets up the following cache directories:
//   - Go build cache: Used for caching build artifacts to speed up subsequent builds.
//   - Go modules cache: Used for caching downloaded Go modules to avoid re-downloading them.
//   - Go test cache: Used for caching test results to speed up test execution.
//
// This method utilizes the WithCachedDirectory function to configure each cache directory with the
// appropriate settings, including the user ownership and shared access. It also sets the GOPATH
// environment variable to the specified path for Go workspace management.
//
// Returns:
//
//	*Gopkgpublisher: The updated Gopkgpublisher instance, allowing for method chaining.
func (m *Gopkgpublisher) WithCacheConfiguration() *Gopkgpublisher {
	const (
		goBuildCacheDir = "/home/gopkgpublisher/.cache/go-build" // Directory for Go build cache
		goModCacheDir   = "/home/gopkgpublisher/go/pkg/mod"      // Directory for Go modules cache
		goTestCacheDir  = "/home/gopkgpublisher/.cache/go-test"  // Directory for Go test cache
	)

	return m.
		// Configure Go build cache
		WithCachedDirectory(
			goBuildCacheDir,
			false,         // Do not create if it does not exist
			"GOCACHE",     // Environment variable name for Go build cache
			dagger.Shared, // Shared access for the cache
			nil,           // No additional options
			ctrUser,       // User owning the cache directory
		).
		// Configure Go modules cache
		WithCachedDirectory(
			goModCacheDir,
			false,         // Do not create if it does not exist
			"GOMODCACHE",  // Environment variable name for Go modules cache
			dagger.Shared, // Shared access for the cache
			nil,           // No additional options
			ctrUser,       // User owning the cache directory
		).
		// Configure Go test cache
		WithCachedDirectory(
			goTestCacheDir,
			false,         // Do not create if it does not exist
			"GOTESTCACHE", // Environment variable name for Go test cache
			dagger.Shared, // Shared access for the cache
			nil,           // No additional options
			ctrUser,       // User owning the cache directory
		).
		// Set the GOPATH environment variable
		WithEnvironmentVariable("GOPATH", "/home/gopkgpublisher/go", false) // Do not overwrite if it exists
}

// WithDirPermissionsConfiguration sets the ownership and permissions for
// the default directories used by the Gopkgpublisher instance. This method
// ensures that the specified user and group have the appropriate ownership
// and permissions for the directories necessary for the operation of the
// Gopkgpublisher. The permissions are set to "0777", allowing read, write,
// and execute access for all users, which may be suitable for development
// environments but should be reviewed for production use.
//
// The default directories configured are:
//   - /home/gopkgpublisher: The home directory for the gopkgpublisher user.
//   - /home/gopkgpublisher/bin: The directory for executable binaries.
//   - mnt: A mount point for additional resources.
//   - /home/gopkgpublisher/.cache/go-build: The cache directory for Go build
//     artifacts.
//   - /home/gopkgpublisher/go/pkg/mod: The directory for Go module cache.
//   - /home/gopkgpublisher/.cache/go-test: The cache directory for Go test
//     artifacts.
//
// Returns a pointer to the updated Gopkgpublisher instance.
func (m *Gopkgpublisher) WithDirPermissionsConfiguration() *Gopkgpublisher {
	defaultDirs := []string{
		"/home/gopkgpublisher",
		"/home/gopkgpublisher/bin",
		"mnt",
		"/home/gopkgpublisher/.cache/go-build",
		"/home/gopkgpublisher/go/pkg/mod",
		"/home/gopkgpublisher/.cache/go-test",
	}

	return m.
		WithUserAsOwnerOfDirs(ctrUser, ctrGroup, defaultDirs, true).
		WithUserWithPermissionsOnDirs(ctrUser, "0777", defaultDirs, true)
}
