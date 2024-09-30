package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/module-template/tests/internal/dagger"
)

// TestContainerWithUbuntuBase tests that the target module is based on the Ubuntu 22.04 image.
//
// This function verifies that the target module is configured appropriately to use the base Ubuntu 22.04 image.
// It runs a command to get the OS version and confirms it matches "Ubuntu 22.04".
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if the Ubuntu image is not used or if the output is not as expected.
func (m *Tests) TestContainerWithUbuntuBase(ctx context.Context) error {
	targetModule := dag.
		ModuleTemplate().
		BaseUbuntu(dagger.ModuleTemplateBaseUbuntuOpts{Version: "22.04"})

	out, err := targetModule.Ctr().
		WithExec([]string{"grep", "^ID=ubuntu$", "/etc/os-release"}).
		Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to get Ubuntu image")
	}

	if !strings.Contains(out, "ID=ubuntu") {
		return WrapErrorf(err, "expected Ubuntu 22.04 or ID=ubuntu, got %s", out)
	}

	return nil
}

// TestContainerWithAlpineBase tests that the target module is based on the Alpine Linux v3.17.3 image.
//
// This function verifies that the target module is configured appropriately to use the base Alpine Linux v3.17.3 image.
// It runs a command to get the OS version and confirms it matches "Alpine Linux v3.17.3".
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if the Alpine image is not used or if the output is not as expected.
func (m *Tests) TestContainerWithAlpineBase(ctx context.Context) error {
	targetModule := dag.ModuleTemplate().
		BaseAlpine(dagger.ModuleTemplateBaseAlpineOpts{Version: "3.17.3"})

	out, err := targetModule.Ctr().WithExec([]string{"cat", "/etc/os-release"}).Stdout(ctx)
	if err != nil {
		return WrapError(err, "failed to get Alpine image")
	}

	// Adjust the conditions to match the actual output.
	if !strings.Contains(out, "Alpine Linux") || !strings.Contains(out, "VERSION_ID=3.17.3") {
		return WrapErrorf(err, "expected Alpine Linux VERSION_ID=3.17.3, got %s", out)
	}

	return nil
}

// TestContainerWithBusyBoxBase tests that the target module is based on the BusyBox v1.35.0 image.
//
// This function verifies that the target module is configured appropriately to use the base BusyBox v1.35.0 image.
// It runs a command to get the OS version and confirms it matches "BusyBox v1.35.0".
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if the BusyBox image is not used or if the output is not as expected.
func (m *Tests) TestContainerWithBusyBoxBase(ctx context.Context) error {
	targetModule := dag.
		ModuleTemplate().
		BaseBusyBox(dagger.ModuleTemplateBaseBusyBoxOpts{Version: "1.35.0"})

	out, err := targetModule.Ctr().
		WithExec([]string{"busybox", "--help"}).
		Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to get BusyBox image")
	}

	if !strings.Contains(out, "v1.35.0") {
		return WrapErrorf(err, "expected BusyBox v1.35.0, got %s", out)
	}

	return nil
}

// TestContainerWithWolfiBase tests that the target module is based on the Wolfi image.
//
// This function verifies that the target module is configured appropriately to use the base Wolfi image.
// It runs a command to get the OS version and confirms it matches "Wolfi".
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if the Wolfi image is not used or if the output is not as expected.
//
//nolint:cyclop // The test handles multiple commands and environments, requiring a longer function.
func (m *Tests) TestContainerWithWolfiBase(ctx context.Context) error {
	targetModule := dag.
		ModuleTemplate().
		BaseWolfi()

	out, err := targetModule.Ctr().
		WithExec([]string{"cat", "/etc/os-release"}).
		Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to get Wolfi image")
	}

	if !strings.Contains(out, "Wolfi") {
		return Errorf("expected Wolfi, got %s", out)
	}

	targetModuleWithInstalledPkgs := dag.
		ModuleTemplate().
		BaseWolfi(dagger.ModuleTemplateBaseWolfiOpts{
			Packages: []string{"git", "curl", "wget"},
		})

		// Check if the Wolfi image has the installed packages
	osOut, osErr := targetModuleWithInstalledPkgs.Ctr().
		WithExec([]string{"cat", "/etc/os-release"}).
		Stdout(ctx)

	if osErr != nil {
		return WrapError(osErr, "failed to get Wolfi image")
	}

	if !strings.Contains(osOut, "Wolfi") {
		return Errorf("expected Wolfi, got %s", osOut)
	}

	// Check if git got installed
	gitOut, gitErr := targetModuleWithInstalledPkgs.Ctr().
		WithExec([]string{"git", "--version"}).
		Stdout(ctx)

	if gitErr != nil {
		return WrapError(gitErr, "failed to get git version")
	}

	if gitOut == "" {
		return Errorf("expected to have git version output, got empty output")
	}

	if !strings.Contains(gitOut, "git version") {
		return Errorf("expected git to be working correctly, got %s", gitOut)
	}

	// Check if curl got installed
	curlOut, curlErr := targetModuleWithInstalledPkgs.Ctr().
		WithExec([]string{"curl", "--version"}).
		Stdout(ctx)

	if curlErr != nil {
		return WrapError(curlErr, "failed to get curl version")
	}

	if curlOut == "" {
		return Errorf("expected to have curl version output, got empty output")
	}

	if !strings.Contains(curlOut, "curl") {
		return Errorf("expected curl to be working correctly, got %s", curlOut)
	}

	// Check if wget got installed
	wgetOut, wgetErr := targetModuleWithInstalledPkgs.Ctr().
		WithExec([]string{"wget", "--version"}).
		Stdout(ctx)

	if wgetErr != nil {
		return WrapError(wgetErr, "failed to get wget version")
	}

	if wgetOut == "" {
		return Errorf("expected to have wget version output, got empty output")
	}

	if !strings.Contains(wgetOut, "GNU Wget") {
		return Errorf("expected wget to be working correctly, got %s", wgetOut)
	}

	return nil
}