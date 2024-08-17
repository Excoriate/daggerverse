package main

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/Excoriate/daggerx/pkg/fixtures"

	"github.com/Excoriate/daggerverse/module-template/tests/internal/dagger"
)

func (m *Tests) TestServiceSimple(ctx context.Context) error {
	// Initialize the target module.
	targetModule := dag.ModuleTemplate()

	// Create a container for Go, so a Go service that exposes a port can be created.
	workdir := filepath.
		Join(fixtures.MntPrefix, "goolang-service")

		// Configure the goCtr that's going to be later on the GoServer
	goCtr := dag.Container().
		From("golang:alpine").
		WithMountedDirectory(fixtures.MntPrefix, m.TestDir).
		WithWorkdir(workdir).
		WithExec([]string{"go", "build", "-o", "goolang-service"}).
		WithExec([]string{"./goolang-service"})

	// Configure the container as a service.
	goServer, goServerErr := targetModule.
		WithServiceFromContainer(goCtr, dagger.
			ModuleTemplateWithServiceFromContainerOpts{
			ExposePorts:     []int{8080},
			EnableDnsgoogle: true,
		}).Start(ctx)

	defer goServer.Stop(ctx)

	if goServerErr != nil {
		return WrapError(goServerErr, "failed to start the GoServer")
	}

	// Initialize the go service in the go server (Dagger service)
	clientCtr := dag.Container().
		From("alpine:latest").
		// Install curl
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithServiceBinding("go-server", goServer)

		// Hit the service's API endpoint
	apiOut, apiErr := clientCtr.
		Terminal().
		WithExec([]string{"curl", "-s", "http://localhost:8080/products"}).
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
