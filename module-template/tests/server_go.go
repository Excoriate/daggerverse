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

// TestGoWithGoServerAdvanced is a comprehensive end-to-end integration test for running a Go server
// with advanced configuration options and testing its API endpoint.
//
// The function sets up and starts a Go server using the Gotoolbox, then configures advanced settings
// such as debugging, HTTP configuration, DNS resolver, and garbage collection settings.
//
// A client container is initialized to send an HTTP request to the server's API endpoint, and the response
// is verified for expected content.
//
// Parameters:
//   - ctx: Context for the test execution, used to control cancellation and deadlines.
//
// Returns:
//   - An error if starting the server, configuring it, or fetching the API response fails.
func (m *Tests) TestGoWithGoServerAdvanced(ctx context.Context) error {
	// Initialize the Go server using the Gotoolbox, specifying the port and source directory.
	goServer := dag.ModuleTemplate(
		dagger.ModuleTemplateOpts{
			Ctr: dag.Container().From("golang:1.23-alpine"),
		},
	).
		NewGoServer()

		// Add the source directory to the GoServer.
	goServer = goServer.WithSource(m.TestDir, dagger.ModuleTemplateGoServerWithSourceOpts{
		Workdir: "golang-server-http",
	})

	// Expose the GoServer on port 8080.
	goServer = goServer.WithExposePort(8080)

	// Configure compile options for the GoServer.
	goServer = goServer.WithCompileOptions(dagger.ModuleTemplateGoServerWithCompileOptionsOpts{
		Verbose: true,
	})

	// Configure advanced options for the GoServer.
	goServer = goServer.
		WithDebugOptions().
		WithHttpsettings().
		WithDnsresolver().
		WithBinaryName("myservice").
		WithRunOptions([]string{"--flag", "testing this API from Dagger"}).
		WithGarbageCollectionSettings()

	goServerService := goServer.InitService()

	// Initialize a Go client container to send an HTTP request to the server.
	clientCtr := dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithServiceBinding("go-server", goServerService)

	// Make a request to the Go server's API endpoint and capture the response.
	apiOut, apiErr := clientCtr.
		WithExec([]string{"curl", "-s", "-v", "go-server:8080/products"}).
		Stdout(ctx)

		// Check for errors in API response.
	if apiErr != nil {
		return WrapError(apiErr, "failed to get API response")
	}

	// Validate that we received a non-empty API response.
	if apiOut == "" {
		return Errorf("expected to have API response, got empty output")
	}

	// Validate that the API response contains the expected content.
	if !strings.Contains(apiOut, "Product A") {
		return Errorf("expected to have API response containing 'Product A', got %s", apiOut)
	}

	clientCtr2 := dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithServiceBinding("go-server", goServerService)

	apiOut2, apiErr2 := clientCtr2.
		WithExec([]string{"curl", "-s", "-v", "go-server:8080/flag"}).
		Stdout(ctx)

	if apiErr2 != nil {
		return WrapError(apiErr2, "failed to get API response")
	}

	if apiOut2 == "" {
		return Errorf("expected to have API response, got empty output")
	}

	if !strings.Contains(apiOut2, "testing this API from Dagger") {
		return Errorf("expected to have API response "+
			"containing 'testing this API from Dagger', got %s", apiOut2)
	}

	return nil
}
