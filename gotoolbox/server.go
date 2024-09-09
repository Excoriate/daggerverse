package main

import (
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/Excoriate/daggerverse/gotoolbox/internal/dagger"

	"github.com/Excoriate/daggerx/pkg/containerx"
	"github.com/Excoriate/daggerx/pkg/fixtures"
)

// Default configuration constants for the GoServer package.
const (
	// defaultBinaryName is the default name of the binary to build and run inside the container.
	// If no name is provided by the user, this default will be used.
	defaultGoServerBinaryName        = "app"
	defaultGoServerProxy             = "https://proxy.golang.org,direct"
	defaultGoServerDNSResolver       = "8.8.8.8 8.8.4.4"
	defaultGoServerGarbageCollection = "100"
	defaultGoServerEnvironment       = "production"
	defaultGoServerDebugOptions      = "http2debug=1"
	defaultGoServerHTTPMaxConns      = "1000"
	defaultGoServerHTTPKeepAlive     = "1"
)

// GoServer represents a Go-based server configuration.
type GoServer struct {
	// serverBinaryName is the name of the binary to build.
	// +private
	ServerBinaryName string

	// Ctr is the container to use as a base container.
	// +private
	Ctr *dagger.Container

	// CompileArgs is the arguments to pass to the go build command.
	// +private
	CompileArgs []string

	// RunArgs is the arguments to pass to the go run command.
	// +private
	RunArgs []string
}

func (m *GoServer) getBinaryName() string {
	if m.ServerBinaryName == "" {
		m.ServerBinaryName = defaultGoServerBinaryName
	}

	return m.ServerBinaryName
}

func (m *GoServer) getGoBuildCMD() []string {
	return []string{"go", "build", "-o", m.getBinaryName()}
}

func (m *GoServer) getExecBinaryCMD() []string {
	return []string{"./" + m.getBinaryName()}
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
		Image:   "golang",
		Version: "1.23-alpine",
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

// WithBinaryName sets the name of the binary to be built for the GoServer.
//
// This method allows specifying a custom name for the binary to be built and used
// within the server container. If the binary name is not explicitly set, the default
// binary name ("app") will be used.
//
// Parameters:
//
//	binaryName string: The name of the binary to build.
//
// Returns:
//
//	*GoServer: An instance of GoServer configured with the specified binary name.
func (m *GoServer) WithBinaryName(
	// binaryName is the name of the binary to build.
	binaryName string,
) *GoServer {
	m.ServerBinaryName = binaryName

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

// WithExposePort sets the port to expose from the service.
//
// This method allows setting the port to expose from the service.
//
// Parameters:
//
//	port int: The port to expose from the service.
//
// Returns:
//
//	*GoServer: An instance of GoServer configured with the provided port.
func (m *GoServer) WithExposePort(
	// ports is a list of ports to expose from the service.
	port int,
	// skipHealthcheck is a flag to skip the health check when run as a service.
	// +optional
	skipHealthcheck bool,
) *GoServer {
	m.Ctr = m.Ctr.WithExposedPort(port, dagger.ContainerWithExposedPortOpts{
		Protocol:                    "TCP",
		ExperimentalSkipHealthcheck: skipHealthcheck,
	})

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

// WithCompileOptions executes a user-defined compilation command within the container.
//
// This method allows appending custom arguments to the "go build -o app" command for compiling
// the Go server. It always defaults to "go build -o app" but allows additional arguments
// such as verbose output, module mode, and build tags to customize the build process.
//
// Parameters:
//
//	extraArgs []string: (optional) Extra arguments to append to the "go build -o app" command.
//	verbose bool: (optional) Flag to enable verbose output (adds the -v flag to the command).
//	mod string: (optional) Sets the module download mode (adds the -mod flag to the command).
//	tags string: (optional) A comma-separated list of build tags to consider
//	satisfied during the build (adds the -tags flag to the command).
//
// Returns:
//
//	*GoServer: An instance of GoServer configured with the custom compilation command and any additional flags.
func (m *GoServer) WithCompileOptions(
	// extraArgs are additional arguments to append to the go build command.
	// +optional
	extraArgs []string,
	// verbose is a flag to enable verbose output.
	// +optional
	verbose bool,
	// mod is the module download mode for the build.
	// +optional
	mod string,
	// tags is a comma-separated list of build tags to consider satisfied during the build.
	// +optional
	tags string,
) *GoServer {
	// Initialize the empty slice to store the compilation arguments
	var cmd []string
	// Add verbose flag if set
	if verbose {
		cmd = append(cmd, "-v")
	}

	// Add mod flag if set
	if mod != "" {
		cmd = append(cmd, "-mod="+mod)
	}

	// Add tags flag if set
	if tags != "" {
		cmd = append(cmd, "-tags="+tags)
	}
	// Add any extra arguments
	if len(extraArgs) > 0 {
		cmd = append(cmd, extraArgs...)
	}

	// Store the compilation arguments in the GoServer instance
	m.CompileArgs = cmd

	return m
}

// WithRunOptions adds extra arguments to the command when running the Go server binary.
//
// This method allows additional arguments to be appended when running the server binary.
// Normally, the server is run as "./binary", but if the binary that represents the Go server
// receives flags or additional arguments, these should be passed through this function.
//
// Parameters:
//
//	runCmd []string: A list of additional command-line arguments to pass
//	when executing the server binary.
//
// Returns:
//
//	*GoServer: An instance of GoServer configured to run the binary with the specified arguments.
func (m *GoServer) WithRunOptions(
	// runCmd is a list of additional command-line arguments to pass when executing the server binary.
	runFlags []string,
) *GoServer {
	m.RunArgs = runFlags

	return m
}

// WithGoProxy sets the Go Proxy URL for the service.
//
// Parameters:
//
//	goproxy string: URL for the Go Proxy. If empty, defaults to "https://proxy.golang.org,direct".
//
// Returns:
//
//	*GoServer: The GoServer instance for method chaining.
func (m *GoServer) WithGoProxy(
	// goproxy is the Go Proxy URL for the service.
	// +optional
	goproxy string) *GoServer {
	if goproxy == "" {
		goproxy = defaultGoServerProxy
	}

	m.Ctr = m.
		Ctr.
		WithEnvVariable("GOPROXY", goproxy)

	return m
}

// WithDNSResolver sets the DNS resolver for the service.
//
// Parameters:
//
//	dnsResolver string: DNS resolver used by the service. If empty, defaults to "8.8.8.8 8.8.4.4".
//
// Returns:
//
//	*GoServer: The GoServer instance for method chaining.
func (m *GoServer) WithDNSResolver(
	// dnsResolver is the DNS resolver used by the service.
	// +optional
	dnsResolver string) *GoServer {
	if dnsResolver == "" {
		dnsResolver = defaultGoServerDNSResolver
	}

	m.Ctr = m.
		Ctr.
		WithEnvVariable("DNS_RESOLVER", dnsResolver)

	return m
}

// WithGarbageCollectionSettings sets the garbage collection optimization for the Go runtime.
//
// Parameters:
//
//	gogc string: Garbage collection optimization value. If empty, defaults to "100".
//
// Returns:
//
//	*GoServer: The GoServer instance for method chaining.
func (m *GoServer) WithGarbageCollectionSettings(
	// gogc is the garbage collection optimization for the Go runtime.
	// +optional
	gogc string) *GoServer {
	if gogc == "" {
		gogc = defaultGoServerGarbageCollection
	}

	m.Ctr = m.
		Ctr.
		WithEnvVariable("GOGC", gogc)

	return m
}

// WithGoEnvironment sets the Go environment for the service.
//
// Parameters:
//
//	goEnv string: The Go environment. If empty, defaults to "production".
//
// Returns:
//
//	*GoServer: The GoServer instance for method chaining.
func (m *GoServer) WithGoEnvironment(
	// goEnv is the Go environment.
	// +optional
	goEnv string) *GoServer {
	if goEnv == "" {
		goEnv = defaultGoServerEnvironment
	}

	m.Ctr = m.
		Ctr.
		WithEnvVariable("GO_ENV", goEnv)

	return m
}

// WithMaxProcs sets the maximum number of CPU cores to use.
//
// Parameters:
//
//	goMaxProcs int: Maximum number of CPU cores to use. If 0, defaults to the number of available CPU cores.
//
// Returns:
//
//	*GoServer: The GoServer instance for method chaining.
func (m *GoServer) WithMaxProcs(
	// goMaxProcs is the maximum number of CPU cores to use.
	// +optional
	goMaxProcs int) *GoServer {
	if goMaxProcs == 0 {
		goMaxProcs = runtime.NumCPU()
	}

	m.Ctr = m.
		Ctr.
		WithEnvVariable("GOMAXPROCS", strconv.Itoa(goMaxProcs))

	return m
}

// WithDebugOptions sets debug options for the Go runtime.
//
// Parameters:
//
//	goDebug string: Debug options for Go. If empty, defaults to "http2debug=1".
//
// Returns:
//
//	*GoServer: The GoServer instance for method chaining.
func (m *GoServer) WithDebugOptions(
	// goDebug is the debug options for the Go runtime.
	// +optional
	goDebug string) *GoServer {
	if goDebug == "" {
		goDebug = defaultGoServerDebugOptions
	}

	m.Ctr = m.Ctr.
		WithEnvVariable("GODEBUG", goDebug)

	return m
}

// WithRuntimeThreadLock optimizes the number of threads in system calls.
//
// Parameters:
//
//	runtimeLockOsThread string: Optimizes number of threads in system calls. If empty, defaults to "1".
//
// Returns:
//
//	*GoServer: The GoServer instance for method chaining.
func (m *GoServer) WithRuntimeThreadLock(
	// runtimeLockOsThread  Optimizes number of threads in system calls.
	// +optional
	runtimeLockOsThread string) *GoServer {
	if runtimeLockOsThread == "" {
		runtimeLockOsThread = "1"
	}

	m.Ctr = m.
		Ctr.
		WithEnvVariable("RUNTIME_LOCKOSTHREAD", runtimeLockOsThread)

	return m
}

// WithHTTPSettings sets HTTP-related settings for the service.
//
// Parameters:
//
//	maxConns string: Maximum number of concurrent HTTP connections. If empty, defaults to "1000".
//	keepAlive string: Enables HTTP Keep-Alive. If empty, defaults to "1".
//
// Returns:
//
//	*GoServer: The GoServer instance for method chaining.
func (m *GoServer) WithHTTPSettings(
	// maxConns is the maximum number of concurrent HTTP connections.
	// +optional
	maxConns string,
	// keepAlive is the HTTP Keep-Alive setting.
	// +optional
	keepAlive string) *GoServer {
	if maxConns == "" {
		maxConns = defaultGoServerHTTPMaxConns
	}

	if keepAlive == "" {
		keepAlive = defaultGoServerHTTPKeepAlive
	}

	m.Ctr = m.
		Ctr.
		WithEnvVariable("HTTP_MAX_CONNS", maxConns)

	m.Ctr = m.
		Ctr.
		WithEnvVariable("HTTP_KEEP_ALIVE", keepAlive)

	return m
}

// InitService sets up a basic Go service with the provided container and exposes the specified ports.
//
// Returns:
//
//	*dagger.Service: The configured service with the specified ports exposed.
func (m *GoServer) InitService(
	// ctr is the container to use as a base container.
	// +optional
	ctr *dagger.Container,
) *dagger.Service {
	if ctr != nil {
		return ctr.AsService()
	}

	return m.
		InitContainer().
		AsService()
}

// InitContainer initializes the container with the default build and run commands.
//
// This method sets up the container for the Go server. If custom compilation
// or run commands have not been provided, it defaults to building the server
// binary using "go build" and running the server binary directly.
//
// Returns:
//
//	*dagger.Container: The initialized container with the default or custom commands.
func (m *GoServer) InitContainer() *dagger.Container {
	goBuildCMD := m.getGoBuildCMD()
	goExecCMD := m.getExecBinaryCMD()

	if len(m.CompileArgs) > 0 {
		goBuildCMD = append(goBuildCMD, m.CompileArgs...)
	}

	if len(m.RunArgs) > 0 {
		goExecCMD = append(goExecCMD, m.RunArgs...)
	}

	m.Ctr = m.Ctr.WithExec(goBuildCMD).WithExec(goExecCMD)

	return m.Ctr
}
