package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/terragrunt/tests/internal/dagger"
)

// TestIACWithTerraformAlpine checks the installation, PATH inclusion, and functionality of Terraform
// on an Alpine Linux environment across various Terraform versions.
//
// It verifies Terraform installation by running `terraform version`, ensures
// Terraform is in the system PATH by using `which terraform`, and tests that
// Terraform commands work correctly by running `terraform --help`.
//
// Args:
//
//	ctx: Context to control execution.
//	version: List of Terraform versions to test. An empty string implies the latest version.
//
// Returns:
//
//	An error if any of the Terraform checks fail for any version; otherwise, nil.
func (m *Tests) TestIACWithTerraformAlpine(ctx context.Context) error {
	// List of Terraform versions to test. An empty string implies the latest version.
	versions := []string{"", "1.0.0", "1.9.4", "1.8.0"}

	for _, version := range versions {
		// Initialize the module with the specified Terraform version on Alpine Linux.
		targetModule := dag.
			Terragrunt().
			WithTerraformAlpine(dagger.
				TerragruntWithTerraformAlpineOpts{
				Version: version,
			})

		// Verify the installation of Terraform.
		if err := m.verifyTerraformInstallation(ctx, targetModule, version); err != nil {
			return err
		}

		// Verify that Terraform is in the system PATH.
		if err := m.verifyTerraformInPath(ctx, targetModule, version); err != nil {
			return err
		}

		// Verify the functionality of Terraform by running the help command.
		if err := m.verifyTerraformHelp(ctx, targetModule, version); err != nil {
			return err
		}
	}

	return nil
}

func (m *Tests) verifyTerraformInstallation(ctx context.Context, module *dagger.Terragrunt, version string) error {
	versionOut, err := module.
		Ctr().
		WithExec([]string{"terraform", "version"}).Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to run terraform version command for version "+version)
	}

	if versionOut == "" {
		return Errorf("expected to have terraform version output, got empty output for version %s", version)
	}

	if !strings.Contains(versionOut, "Terraform v") {
		return Errorf("expected Terraform to be working correctly, got %s for version %s", versionOut, version)
	}

	if version != "" && !strings.Contains(versionOut, "Terraform v"+version) {
		return Errorf("expected Terraform version %s, got %s", version, versionOut)
	}

	return nil
}

func (m *Tests) verifyTerraformInPath(ctx context.Context, module *dagger.Terragrunt, version string) error {
	pathOut, err := module.Ctr().WithExec([]string{"which", "terraform"}).Stdout(ctx)
	if err != nil {
		return WrapError(err, "failed to find terraform in PATH for version "+version)
	}

	if pathOut == "" {
		return Errorf("expected to have terraform in PATH, got empty output for version %s", version)
	}

	// Trim any whitespace from the pathOut
	pathOut = strings.TrimSpace(pathOut)

	// Use the path to check if in the filesystem
	_, err = module.Ctr().WithExec([]string{"ls", "-l", pathOut}).Stdout(ctx)
	if err != nil {
		return WrapError(err, "failed to find terraform in filesystem for version "+version)
	}

	return nil
}

func (m *Tests) verifyTerraformHelp(ctx context.Context, module *dagger.Terragrunt, version string) error {
	helpOut, err := module.Ctr().WithExec([]string{"terraform", "--help"}).Stdout(ctx)
	if err != nil {
		return WrapError(err, "failed to run terraform help command for version "+version)
	}

	if helpOut == "" {
		return Errorf("expected to have terraform help output, got empty output for version %s", version)
	}

	if !strings.Contains(helpOut, "Usage: terraform") {
		return Errorf("expected Terraform to be working correctly, got %s for version %s", helpOut, version)
	}

	return nil
}

// TestIACWithTerraformUbuntu checks the installation, PATH inclusion, and functionality of Terraform
// on an Ubuntu Linux environment across various Terraform versions.
//
// It verifies Terraform installation by running `terraform version`, ensures Terraform is in
// the system PATH by using `which terraform`, and tests that Terraform commands work correctly
// by running `terraform --help`.
//
// Args:
//
//	ctx: Context to control execution.
//	version: List of Terraform versions to test. An empty string implies the latest version.
//
// Returns:
//
//	An error if any of the Terraform checks fail for any version; otherwise, nil.
func (m *Tests) TestIACWithTerraformUbuntu(ctx context.Context) error {
	// Set up Ubuntu container.
	ubuntuCtr := dag.Container().From("ubuntu:latest")

	// List of Terraform versions to test. An empty string implies the latest version.
	versions := []string{"", "1.0.0", "1.9.0", "1.8.0"}

	for _, version := range versions {
		// Initialize the module with the specified Terraform version on Ubuntu.
		targetModule := dag.
			Terragrunt(dagger.TerragruntOpts{
				Ctr: ubuntuCtr,
			}).
			WithTerraformUbuntu(dagger.TerragruntWithTerraformUbuntuOpts{
				Version: version,
			})

		// Verify the installation of Terraform.
		if err := m.verifyTerraformInstallation(ctx, targetModule, version); err != nil {
			return err
		}

		// Verify that Terraform is in the system PATH.
		if err := m.verifyTerraformInPath(ctx, targetModule, version); err != nil {
			return err
		}

		// Verify the functionality of Terraform by running the help command.
		if err := m.verifyTerraformHelp(ctx, targetModule, version); err != nil {
			return err
		}
	}

	return nil
}
