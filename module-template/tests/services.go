package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/module-template/tests/internal/dagger"
)

// TestCreateService tests the creation and deployment of a Go service container.
// It performs the following steps:
// 1. Configures the base Go container using the `golang:alpine` image.
// 2. Builds the Go service inside the container.
// 3. Creates and starts a service from the configured Go container with exposed ports.
// 4. Runs a client container to hit the service's API endpoint and verifies the response.
//
// Args:
//
//	ctx (context.Context): The context for managing the lifecycle of the request.
//
// Returns:
//
//	error: Returns error if any of the operations fail. Otherwise, returns nil.
func (m *Tests) TestCreateService(ctx context.Context) error {
	// Configure the goCtr that's going to be later on the GoServer
	configuredModule := dag.ModuleTemplate(
		dagger.ModuleTemplateOpts{
			Ctr: dag.Container().
				From("golang:alpine"),
		}).
		WithSource(m.TestDir, dagger.ModuleTemplateWithSourceOpts{
			Workdir: "/goolang-service",
		})

		// Base container for this Go service.
	goCtr := configuredModule.Ctr().
		WithExec([]string{"ls", "-l"}).
		WithExec([]string{"go", "build", "-o", "gosvc"}).
		WithExec([]string{"./gosvc"})

	// Configure the container as a service.
	goServer, goServerErr := dag.
		ModuleTemplate().
		// Configure the service.
		CreateServiceFromContainer(goCtr, dagger.
			ModuleTemplateCreateServiceFromContainerOpts{
			ExposePorts:     []int{8080},
			EnableDnsgoogle: true,
		}).
		Start(ctx)

	defer goServer.Stop(ctx)

	if goServerErr != nil {
		return WrapError(goServerErr, "failed to start the GoServer on port 8080")
	}

	// Initialize the go service in the go server (Dagger service)
	clientCtr := dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithServiceBinding("go-server", goServer)

	// Hit the service's API endpoint
	apiOut, apiErr := clientCtr.
		WithExec([]string{"curl", "-s", "go-server:8080/products"}).
		Stdout(ctx)

	if apiErr != nil {
		return WrapError(apiErr, "failed to get API response")
	}

	if apiOut == "" {
		return Errorf("expected to have API response, got empty output")
	}

	if !strings.Contains(apiOut, "Product A") {
		return Errorf("expected to have API response containing 'Product A', got %s", apiOut)
	}

	return nil
}
