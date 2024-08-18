package main

import (
	"github.com/Excoriate/daggerverse/module-template/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/envvars"
)

// CreateServiceFromContainer sets up a Go service from scratch with
// the provided container and exposes the specified ports.
//
// Parameters:
//
//	exposePorts []int: A list of ports to expose from the service (optional).
//	protocol string: The port protocol to use. Either TCP or UDP. Default is TCP (optional).
//	skipHealthcheck bool: Flag to skip the health check when run as a service (optional).
//	cacheVolume *dagger.CacheVolume: The cache volume to use for the
//	service. Handy to persist data and avoid the 10 seconds grace period (optional).
//	ctr *dagger.Container: The container to use as a base container, which will be exposed as a service.
//	enableDNSGoogle bool: Flag to enable Google DNS for the service (optional).
//
// Returns:
//
//	*dagger.Service: The configured service with the specified ports exposed.
func (m *ModuleTemplate) CreateServiceFromContainer(
	// exposePorts is a list of ports to expose from the service.
	// +optional
	exposePorts []int,
	// portProtocol is the port protocol to use. Either TCP or UDP. Default is TCP.
	// +optional
	portProtocol string,
	// Skip the health check when run as a service.
	// +optional
	skipHealthcheck bool,
	// cacheVolume is the cache volume to use for the service. It's handy to persist data, and
	// avoid the 10 seconds grace period (otherwise, it'll be restarted, and data will be lost).
	// +optional
	cacheVolume *dagger.CacheVolume,
	// cacheVolumePath is the path to the cache volume to use for the service. If it's not set,
	// it will default to "/data".
	// +optional
	cacheVolumePath string,
	// ctr is the container to use as a base container.
	// cacheVolumeSharing is the sharing mode of the cache volume. If not provided, it defaults to "shared".
	// +optional
	cacheVolumeSharing dagger.CacheSharingMode,
	// cacheVolumeOwner is the owner of the cache volume. If not provided, it defaults to "1000:1000".
	// +optional
	cacheVolumeOwner string,
	// ctr is the container to use as a base container for this service.
	ctr *dagger.Container,
	// enableDNSGoogle is a flag to enable Google DNS for the service.
	// +optional
	enableDNSGoogle bool,
	// envVars is a list of environment variables to pass from the host to the container.
	// +optional
	envVars []string,
) *dagger.Service {
	if portProtocol == "" {
		portProtocol = "TCP"
	}

	if cacheVolumePath == "" {
		cacheVolumePath = "/data"
	}

	if cacheVolumeSharing == "" {
		cacheVolumeSharing = dagger.Shared
	}

	if cacheVolumeOwner == "" {
		cacheVolumeOwner = "1000:1000"
	}

	for _, port := range exposePorts {
		ctr = ctr.WithExposedPort(port, dagger.ContainerWithExposedPortOpts{
			Protocol:                    dagger.NetworkProtocol(portProtocol),
			ExperimentalSkipHealthcheck: skipHealthcheck,
		})
	}

	if cacheVolume != nil {
		ctr = ctr.
			WithMountedCache(cacheVolumePath, cacheVolume, dagger.
				ContainerWithMountedCacheOpts{
				Sharing: cacheVolumeSharing,
				Owner:   cacheVolumeOwner,
			})
	}

	if enableDNSGoogle {
		ctr = ctr.
			WithEnvVariable("DNS_RESOLVER", "8.8.8.8 8.8.4.4")
	}

	if envVars != nil {
		envVarsAsDagger, _ := envvars.ToDaggerEnvVarsFromSlice(envVars)
		for _, envVar := range envVarsAsDagger {
			ctr = ctr.
				WithEnvVariable(envVar.Name, envVar.Value, dagger.ContainerWithEnvVariableOpts{
					Expand: envVar.Expand,
				})
		}
	}

	return ctr.AsService()
}
