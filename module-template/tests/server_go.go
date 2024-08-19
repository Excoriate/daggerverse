package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/module-template/tests/internal/dagger"
)

// TestGoWithGoServerSimple is an end-to-end integration test for running a Go server
// and testing its API endpoint.
//
// This function sets up and starts a Go server using the Gotoolbox,
// then initializes a client container to send an HTTP request to the server.
//
// ctx: The context for the test execution, to control cancellation and deadlines.
//
// Returns an error if starting the server or fetching the API response fails.
func (m *Tests) TestGoWithGoServerSimple(ctx context.Context) error {
	// Create the Go server using the ModuleTemplate, specifying the port and source directory.
	goServer := dag.ModuleTemplate(dagger.ModuleTemplateOpts{
		Ctr: dag.Container().From("golang:1.23-alpine"),
	}).NewGoServer()

	// Configure the Go server.
	goServer = goServer.
		WithSource(m.TestDir, dagger.ModuleTemplateGoServerWithSourceOpts{
			Workdir: "golang-server-http",
		}).WithExposePort(8080)

	// Initialize the clientCtr container with necessary tools and bind it to the Go server.
	clientCtr := dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithServiceBinding("golang-server", goServer.InitService())

	// Make a request to the Go server's API endpoint and capture the response.
	apiOut, apiErr := clientCtr.
		WithExec([]string{"curl", "-s", "-v", "golang-server:8080/products"}).
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
