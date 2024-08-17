package main

import "github.com/Excoriate/daggerverse/module-template/internal/dagger"

// WithServiceFromContainer sets up a Go service from scratch with
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
func (m *ModuleTemplate) WithServiceFromContainer(
	exposePorts []int,
	// protocol is the port protocol to use. Either TCP or UDP. Default is TCP.
	// +optional
	protocol string,
	// Skip the health check when run as a service.
	// +optional
	skipHealthcheck bool,
	// cacheVolume is the cache volume to use for the service. It's handy to persist data, and
	// avoid the 10 seconds grace period (otherwise, it'll be restarted, and data will be lost).
	// +optional
	cacheVolume *dagger.CacheVolume,
	// ctr is the container to use as a base container.
	ctr *dagger.Container,
	// enableDNSGoogle is a flag to enable Google DNS for the service.
	// +optional
	enableDNSGoogle bool,
) *dagger.Service {
	if protocol == "" {
		protocol = "TCP"
	}

	for _, port := range exposePorts {
		ctr = ctr.WithExposedPort(port, dagger.ContainerWithExposedPortOpts{
			Protocol:                    dagger.NetworkProtocol(protocol),
			ExperimentalSkipHealthcheck: skipHealthcheck,
		})
	}

	if cacheVolume != nil {
		ctr = ctr.
			WithMountedCache("/data", cacheVolume, dagger.
				ContainerWithMountedCacheOpts{
				Sharing: dagger.Shared,
			})
	}

	if enableDNSGoogle {
		ctr = ctr.
			WithEnvVariable("DNS_RESOLVER", "8.8.8.8 8.8.4.4")
	}

	return ctr.AsService()
}
