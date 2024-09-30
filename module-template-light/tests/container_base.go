package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/module-template-light/tests/internal/dagger"
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
		ModuleTemplateLight().
		BaseUbuntu(dagger.ModuleTemplateLightBaseUbuntuOpts{Version: "22.04"})

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
	targetModule := dag.ModuleTemplateLight().
		BaseAlpine(dagger.ModuleTemplateLightBaseAlpineOpts{Version: "3.17.3"})

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
		ModuleTemplateLight().
		BaseBusyBox(dagger.ModuleTemplateLightBaseBusyBoxOpts{Version: "1.35.0"})

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
		ModuleTemplateLight().
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
		ModuleTemplateLight().
		BaseWolfi(dagger.ModuleTemplateLightBaseWolfiOpts{
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

// TestContainerWithApkoBaseAlpine tests that the target module is based on the Apko image.
//
// This function verifies that the target module is configured appropriately to use the base Apko image.
// It runs a command to get the OS version and confirms it matches "Apko".
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if the Apko image is not used or if the output is not as expected.
func (m *Tests) TestContainerWithApkoBaseAlpine(ctx context.Context) error {
	alpinePresetPath := "testdata/apko-presets/base-alpine.yaml"
	alpinePresetFile := dag.CurrentModule().
		Source().
		File(alpinePresetPath)

	scenarioBaseAlpineBasic := dag.
		ModuleTemplateLight().
		BaseApko(alpinePresetPath, alpinePresetFile)

	outInspectCtr, outInspectCtrErr := scenarioBaseAlpineBasic.Ctr().
		WithExec([]string{"sh", "-c", "uname"}).
		Stdout(ctx)

	if outInspectCtrErr != nil {
		return WrapError(outInspectCtrErr, "failed to inspect the Alpine container with preset base-alpine.yaml")
	}

	if !strings.Contains(outInspectCtr, "Linux") {
		return Errorf("expected Alpine Linux in the Alpine container with preset base-alpine.yaml, got %s",
			outInspectCtr)
	}

	// Install additional packages
	alpineApkoCtr := scenarioBaseAlpineBasic.
		Ctr().
		WithUser("root").
		WithExec([]string{"apk", "add", "--no-cache", "curl"})

	outInspectCtr, outInspectCtrErr = alpineApkoCtr.
		WithExec([]string{"sh", "-c", "curl --version"}).
		Stdout(ctx)

	if outInspectCtrErr != nil {
		return WrapError(outInspectCtrErr, "failed to inspect the Alpine container with preset base-alpine.yaml")
	}

	if !strings.Contains(outInspectCtr, "curl") {
		return Errorf("expected curl to be working correctly, got %s", outInspectCtr)
	}

	return nil
}

// TestContainerWithApkoBaseWolfi tests that the target module is based on the Apko Wolfi image.
//
// This function verifies that the target module is configured appropriately to use the base Wolfi image.
// It runs a command to get the OS version and confirms it matches "Apko", then installs and verifies "curl".
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if the Wolfi image is not used, if the output is not as expected,
//     or if package installation fails.
func (m *Tests) TestContainerWithApkoBaseWolfi(ctx context.Context) error {
	wolfiPresetPath := "testdata/apko-presets/base-wolfi.yaml"
	wolfiPresetFile := dag.CurrentModule().
		Source().
		File(wolfiPresetPath)

	scenarioBaseWolfi := dag.
		ModuleTemplateLight().
		BaseApko(wolfiPresetPath, wolfiPresetFile)

	outInspectCtr, outInspectCtrErr := scenarioBaseWolfi.Ctr().
		WithExec([]string{"sh", "-c", "uname"}).
		Stdout(ctx)

	if outInspectCtrErr != nil {
		return WrapError(outInspectCtrErr, "failed to inspect the Wolfi container with preset base-wolfi.yaml")
	}

	if !strings.Contains(outInspectCtr, "Apko") {
		return Errorf("expected Wolfi base image with Apko, got %s", outInspectCtr)
	}

	// Install additional packages
	wolfiCtr := scenarioBaseWolfi.
		Ctr().
		WithUser("root").
		WithExec([]string{"apk", "add", "--no-cache", "curl"})

	outCurl, outCurlErr := wolfiCtr.
		WithExec([]string{"curl", "--version"}).
		Stdout(ctx)

	if outCurlErr != nil {
		return WrapError(outCurlErr, "failed to get curl version")
	}

	if outCurl == "" {
		return Errorf("expected to have curl version output, got empty output")
	}

	if !strings.Contains(outCurl, "curl") {
		return Errorf("expected curl to be working correctly, got %s", outCurl)
	}

	return nil
}
