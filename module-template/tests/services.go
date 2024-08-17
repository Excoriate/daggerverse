package main

import (
	"context"

	"github.com/Excoriate/daggerverse/module-template/tests/internal/dagger"
)

func (m *Tests) TestServiceSimple(ctx context.Context) error {
	// Create a container for Go, so a Go service that exposes a port can be created.
	goCtr := dag.Container().
		From("golang:alpine")

		// Initialize the target module.
	targetModule := dag.ModuleTemplate()

	// Configure the container as a service.
	goSvc := targetModule.
		WithServiceFromContainer(goCtr, dagger.
			ModuleTemplateWithServiceFromContainerOpts{
			ExposePorts: []int{8080},
		})

	return nil
}
