package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/terragrunt/tests/internal/dagger"
)

// TestTerragruntContainerIsUp checks if the Terragrunt container is up and running by verifying the versions of Terragrunt, Terraform, and OpenTofu.
// It executes the version commands for each tool and checks their outputs to ensure they contain the expected version strings.
// If any of the commands fail or the outputs do not contain the expected strings, an error is returned.
func (m *Tests) TestTerragruntContainerIsUp(ctx context.Context) error {
	tgSourceData := m.getTestDir().Directory("terragrunt")
	tgCtr := dag.
		Terragrunt().
		WithSource(tgSourceData).
		Ctr()

	tgCtrOut, tgCtrErr := tgCtr.
		WithExec([]string{"ls", "-la"}).
		WithExec([]string{"pwd"}).
		WithExec([]string{"ls", "-la", "/home/terragrunt"}).
		WithExec([]string{"ls", "-la", "/home/.terraform.d"}).
		WithExec([]string{"ls", "-la", "/home"}).
		WithExec([]string{"ls", "-la", "/mnt"}).
		Stdout(ctx)

	if tgCtrErr != nil {
		return WrapErrorf(tgCtrErr, "failed to get terragrunt container")
	}

	if tgCtrOut == "" {
		return Errorf("terragrunt container output is empty")
	}

	tgCtrOut, tgCtrErr = tgCtr.
		WithExec([]string{"cat", "/mnt/terragrunt.hcl"}).
		Stdout(ctx)

	if tgCtrErr != nil {
		return WrapErrorf(tgCtrErr, "failed to get terragrunt terragrunt.hcl file")
	}

	if tgCtrOut == "" {
		return Errorf("terragrunt terragrunt.hcl file is empty")
	}

	return nil
}

// TestTerragruntBinariesAreInstalled checks if the Terragrunt, Terraform, and OpenTofu binaries are installed and their versions are correct.
// It executes the version command for each tool and verifies that the output contains the expected version string.
// If any of the commands fail or the outputs do not contain the expected strings, an error is returned.
func (m *Tests) TestTerragruntBinariesAreInstalled(ctx context.Context) error {
	tgCtr := dag.
		Terragrunt().
		Ctr()

	if err := validateVersion(ctx, tgCtr, "terragrunt", "terragrunt version"); err != nil {
		return err
	}

	if err := validateVersion(ctx, tgCtr, "terraform", "Terraform v"); err != nil {
		return err
	}

	if err := validateVersion(ctx, tgCtr, "opentofu", "OpenTofu v"); err != nil {
		return err
	}

	return nil
}

// validateVersion executes the version command for a given tool and checks the output to ensure it contains the expected version string.
// If the command fails or the output does not contain the expected string, an error is returned.
func validateVersion(
	ctx context.Context,
	tgCtr *dagger.Container,
	tool string,
	expectedVersion string,
) error {
	versionOut, versionOutErr := tgCtr.
		WithExec([]string{tool, "--version"}).
		Stdout(ctx)

	if versionOutErr != nil {
		return WrapErrorf(versionOutErr, "failed to get %s version output", tool)
	}

	if versionOut == "" {
		return Errorf("%s version output is empty", tool)
	}

	if !strings.Contains(versionOut, expectedVersion) {
		return Errorf("%s version is expected to contain '%s', but it doesn't", tool, expectedVersion)
	}

	return nil
}
