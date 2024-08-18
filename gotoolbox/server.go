package main

import (
	"fmt"
	"path/filepath"

	"github.com/Excoriate/daggerverse/gotoolbox/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/containerx"
	"github.com/Excoriate/daggerx/pkg/fixtures"
)

// Default configuration constants for the GoServer package.
const (
	// defaultBinaryName is the default name of the binary to build and run inside the container.
	// If no name is provided by the user, this default will be used.
	defaultGoServerBinaryName = "app"
)

// GoServer represents a Go-based server configuration.
type GoServer struct {
	// serverBinaryName is the name of the binary to build.
	// +private
	ServerBinaryName string

	// Ctr is the container to use as a base container.
	// +private
	Ctr *dagger.Container
}

func (m *GoServer) setDefaults() *GoServer {
	if m.ServerBinaryName == "" {
		m.ServerBinaryName = defaultGoServerBinaryName
	}

	return m
}

// NewGoServer initializes and returns a new instance of GoServer with the given service name and port.
//
// Parameters:
//
// serviceName string: The name of the service to be created (optional, defaults to "go-server").
// port int: The port to expose from the service.
//
// Returns:
//
// *GoServer: An instance of GoServer configured with a container created from the default image and version.
func (m *Gotoolbox) NewGoServer(
	// ctr is the container to use as a base container. If it's not set, it'll create a new container.
	// +optional
	ctr *dagger.Container,
) *GoServer {
	if ctr != nil {
		m.Ctr = ctr
		return &GoServer{Ctr: m.Ctr}
	}
	// Get the default container image URL
	imageURL, _ := containerx.GetImageURL(&containerx.
		NewBaseContainerOpts{
		Image:           defaultContainerImage,
		Version:         defaultContainerVersion,
		FallBackVersion: defaultContainerVersion,
		FallbackImage:   defaultContainerImage,
	})

	// Return a new GoServer instance with configured service name and container
	return &GoServer{
		Ctr: dag.Container().From(imageURL),
	}
}

// WithServerData configures the GoServer to use a cache volume at a specified path
// with specified sharing mode and ownership.
//
// This method mounts a cache volume inside the container at the provided path,
// with specified sharing mode and ownership details. If any of these parameters
// are not provided, default values will be used.
//
// Parameters:
//
// path string: (optional) The path to the cache volume's root. Defaults to "/data" if not provided.
// share dagger.CacheSharingMode: (optional) The sharing mode of the cache volume. Defaults to "shared" if not provided.
// owner string: (optional) The owner of the cache volume. Defaults to "1000:1000" if not provided.
//
// Returns:
//
// *GoServer: An instance of GoServer configured with the specified cache volume settings.
func (m *GoServer) WithServerData(
	// path is the path to the cache volume's root. If not provided, it defaults to "/data".
	// +optional
	path string,
	// share is the sharing mode of the cache volume. If not provided, it defaults to "shared".
	// +optional
	share dagger.CacheSharingMode,
	// owner is the owner of the cache volume. If not provided, it defaults to "1000:1000".
	// +optional
	owner string,
	// workdir is the working directory within the container. If not set it'll default to /mnt
	// +optional
	workdir string,
) *GoServer {
	// Set default values if not provided
	if path == "" {
		path = "/data"
	}

	if share == "" {
		share = dagger.Shared
	}

	if owner == "" {
		owner = "1000:1000"
	}

	// Create and configure cache volume
	cacheVolume := dag.CacheVolume("server-data")
	ctr := m.Ctr.WithMountedCache(path, cacheVolume, dagger.ContainerWithMountedCacheOpts{
		Sharing: share,
		Owner:   owner,
	})

	if workdir != "" {
		ctr = ctr.WithWorkdir(filepath.Join(fixtures.MntPrefix, workdir))
	}

	// Update the container configuration in the GoServer
	m.Ctr = ctr

	return m
}

// WithPreBuiltContainer configures the GoServer to use a pre-existing container as its base.
//
// This method allows setting an already created container as the base for the GoServer,
// overriding any previously set container.
//
// Parameters:
//
//	ctr *dagger.Container: The container to use as a base container.
//
// Returns:
//
//	*GoServer: An instance of GoServer configured with the provided container.
func (m *GoServer) WithPreBuiltContainer(
	// ctr is the container to use as a base container.
	ctr *dagger.Container,
) *GoServer {
	m.Ctr = ctr
	return m
}

// WithExposePorts sets the port to expose from the service.
//
// This method allows setting the port to expose from the service.
//
// Parameters:
//
//	ports []int: A list of ports to expose from the service.
//
// Returns:
//
//	*GoServer: An instance of GoServer configured with the provided port.
func (m *GoServer) WithExposePorts(
	// ports is a list of ports to expose from the service.
	ports []int,
) *GoServer {
	for _, port := range ports {
		m.Ctr = m.Ctr.WithExposedPort(port, dagger.ContainerWithExposedPortOpts{
			Protocol:                    "TCP",
			ExperimentalSkipHealthcheck: false,
		})
	}

	return m
}

// WithSource mounts the source directory inside the container and sets the working directory.
//
// This method configures the GoServer to mount the provided source directory at a fixed
// mount point and optionally set a specific working directory within the container. If
// the working directory is not provided, it defaults to the mount point.
//
// Parameters:
//
//	src *dagger.Directory: The directory containing all the source code, including the module directory.
//	workdir string: (optional) The working directory within the container, defaults to "/mnt".
//
// Returns:
//
//	*GoServer: An instance of GoServer configured with the provided source directory and working directory.
func (m *GoServer) WithSource(
	// src is the directory that contains all the source code, including the module directory.
	src *dagger.Directory,
	// workdir is the working directory within the container. If not set it'll default to /mnt
	// +optional
	workdir string,
) *GoServer {
	// Mount the source directory at the fixed mount point
	ctr := m.Ctr.
		WithMountedDirectory(fixtures.MntPrefix, src)

	// Set the working directory, defaulting to the mount point if not provided
	if workdir != "" {
		ctr = ctr.WithWorkdir(filepath.Join(fixtures.MntPrefix, workdir))
	} else {
		ctr = ctr.WithWorkdir(fixtures.MntPrefix)
	}

	// Update the container configuration in the GoServer
	m.Ctr = ctr

	return m
}

// Init sets up a Go service from scratch with the provided container
// and exposes the specified ports.
//
// Parameters:
//
//	exposePorts []int: A list of ports to expose from the service (optional).
//	ctr *dagger.Container: The container to use as a base container, which will be exposed as a service.
//
// Returns:
//
//	*dagger.Service: The configured service with the specified ports exposed.
func (m *GoServer) Init() *dagger.Service {
	m.setDefaults()

	envVars := map[string]string{
		"GOPROXY":      "https://proxy.golang.org,direct",
		"DNS_RESOLVER": "8.8.8.8 8.8.4.4",
	}

	for key, value := range envVars {
		m.Ctr = m.Ctr.WithEnvVariable(key, value)
	}

	return m.Ctr.
		WithExec([]string{"go", "build", "-o", m.ServerBinaryName}).
		WithExec([]string{fmt.Sprintf("./%s", m.ServerBinaryName)}).AsService()
}
