package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/terragrunt/tests/internal/dagger"
)

// TestContainerBaseApkoDefault is a test function that verifies the default behavior of the Apko container.
// It runs a shell command "ls -l" in the context of the target module and checks the output for expected results.
// If the command fails or the output does not meet the expectations, it returns an error.
func (m *Tests) TestContainerBaseApkoDefault(ctx context.Context) error {
	// Initialize the target module with default Terragrunt options.
	targetModule := dag.Terragrunt()
	targetCtr := targetModule.Ctr()

	if err := checkLsOutput(ctx, targetCtr); err != nil {
		return err
	}

	if err := checkTerragruntInstalled(ctx, targetCtr); err != nil {
		return err
	}

	if err := checkTerraformInstalled(ctx, targetCtr); err != nil {
		return err
	}

	if err := checkOpentofuInstalled(ctx, targetCtr); err != nil {
		return err
	}

	return nil
}

func checkLsOutput(ctx context.Context, targetCtr *dagger.Container) error {
	simpleCmdOut, simpleCmdErr := targetCtr.WithExec([]string{"ls", "-l"}).Stdout(ctx)
	if simpleCmdErr != nil {
		return WrapError(simpleCmdErr, "failed to run shell command")
	}

	if simpleCmdOut == "" {
		return WrapError(simpleCmdErr, "expected to have at least one folder, got empty output")
	}

	if !strings.Contains(simpleCmdOut, "total") {
		return WrapErrorf(simpleCmdErr, "expected to have at least one folder, got %s", simpleCmdOut)
	}

	return nil
}

func checkTerragruntInstalled(ctx context.Context, targetCtr *dagger.Container) error {
	whichTgCmdOut, whichTgCmdErr := targetCtr.WithExec([]string{"which", "terragrunt"}).Stdout(ctx)
	if whichTgCmdErr != nil {
		return WrapError(whichTgCmdErr, "failed to check if terragrunt is installed")
	}

	if whichTgCmdOut == "" {
		return WrapError(whichTgCmdErr, "terragrunt is not installed")
	}

	if !strings.Contains(whichTgCmdOut, "/home/terragrunt/bin/terragrunt") {
		return WrapErrorf(whichTgCmdErr, "expected terragrunt to be installed in /home/terragrunt/bin/terragrunt, got %s",
			whichTgCmdOut)
	}

	return nil
}

func checkTerraformInstalled(ctx context.Context, targetCtr *dagger.Container) error {
	whichTfCmdOut, whichTfCmdErr := targetCtr.WithExec([]string{"which", "terraform"}).Stdout(ctx)
	if whichTfCmdErr != nil {
		return WrapError(whichTfCmdErr, "failed to check if terraform is installed")
	}

	if whichTfCmdOut == "" {
		return WrapError(whichTfCmdErr, "terraform is not installed")
	}

	if !strings.Contains(whichTfCmdOut, "/home/terragrunt/bin/terraform") {
		return WrapErrorf(whichTfCmdErr, "expected terraform to be installed in /home/terragrunt/bin/terraform, got %s",
			whichTfCmdOut)
	}

	return nil
}

func checkOpentofuInstalled(ctx context.Context, targetCtr *dagger.Container) error {
	whichOtCmdOut, whichOtCmdErr := targetCtr.WithExec([]string{"which", "opentofu"}).Stdout(ctx)
	if whichOtCmdErr != nil {
		return WrapError(whichOtCmdErr, "failed to check if opentofu is installed")
	}

	if whichOtCmdOut == "" {
		return WrapError(whichOtCmdErr, "opentofu is not installed")
	}

	if !strings.Contains(whichOtCmdOut, "/home/terragrunt/bin/opentofu") {
		return WrapErrorf(whichOtCmdErr, "expected opentofu to be installed in /home/terragrunt/bin/opentofu, got %s",
			whichOtCmdOut)
	}

	return nil
}

// TestContainerBaseWithPassedImage tests the initialization of a target module with a specified image
// and verifies the execution of basic shell commands within the context of the target module.
//
// This function performs the following steps:
// 1. Initializes the target module with default Terragrunt options, using the specified image URL.
// 2. Runs the shell command "ls -l" within the context of the target module and checks for errors and expected output.
// 3. Runs the "terragrunt --version" command within the context of the target module and verifies the version output.
//
// Parameters:
// - ctx: The context for controlling the execution of the function.
//
// Returns:
// - An error if any of the steps fail, otherwise nil.
func (m *Tests) TestContainerBaseWithPassedImage(ctx context.Context) error {
	// Initialize the target module with default Terragrunt options.
	targetModule := dag.Terragrunt(dagger.TerragruntOpts{
		ImageURL: "alpine/terragrunt:latest",
	})

	targetCtr := targetModule.Ctr()

	// Run the shell command "ls -l" in the context of the target module.
	simpleCmdOut, simpleCmdErr := targetCtr.WithExec([]string{"ls", "-l"}).Stdout(ctx)
	if simpleCmdErr != nil {
		return WrapError(simpleCmdErr, "failed to run shell command")
	}

	if simpleCmdOut == "" {
		return WrapError(simpleCmdErr, "expected to have at least one folder, got empty output")
	}

	// Run terragrunt --version
	terragruntVersionCmdOut, terragruntVersionCmdErr := targetCtr.WithExec([]string{"terragrunt", "--version"}).Stdout(ctx)
	if terragruntVersionCmdErr != nil {
		return WrapError(terragruntVersionCmdErr, "failed to run terragrunt --version")
	}

	if !strings.HasPrefix(terragruntVersionCmdOut, "terragrunt version v") {
		return WrapErrorf(terragruntVersionCmdErr,
			"expected terragrunt version to start with v0.66 and be greater than 0.66.0, got %s",
			terragruntVersionCmdOut)
	}

	return nil
}

// TestContainerBaseWithAWSClI tests the initialization of a target module with AWS CLI enabled
// and verifies the execution of basic shell commands within the context of the target module.
//
// This function performs the following steps:
// 1. Initializes the target module with default Terragrunt options, enabling AWS CLI.
// 2. Runs the shell command "ls -l" within the context of the target module and checks for errors and expected output.
// 3. Runs the "aws --version" command within the context of the target module and verifies the version output.
//
// Parameters:
// - ctx: The context for controlling the execution of the function.
//
// Returns:
// - An error if any of the steps fail, otherwise nil.
func (m *Tests) TestContainerBaseWithAWSClI(ctx context.Context) error {
	// Initialize the target module with default Terragrunt options, enabling AWS CLI.
	targetModule := dag.Terragrunt(dagger.TerragruntOpts{
		EnableAwscli: true,
	})

	targetCtr := targetModule.Ctr()

	// Run the shell command "ls -l" in the context of the target module.
	simpleCmdOut, simpleCmdErr := targetCtr.WithExec([]string{"ls", "-l"}).Stdout(ctx)
	if simpleCmdErr != nil {
		return WrapError(simpleCmdErr, "failed to run shell command")
	}

	if simpleCmdOut == "" {
		return WrapError(simpleCmdErr, "expected to have at least one folder, got empty output")
	}

	// Run aws --version
	awsVersionCmdOut, awsVersionCmdErr := targetCtr.WithExec([]string{"aws", "--version"}).Stdout(ctx)
	if awsVersionCmdErr != nil {
		return WrapError(awsVersionCmdErr, "failed to run aws --version")
	}

	if !strings.HasPrefix(awsVersionCmdOut, "aws-cli/2.") {
		return WrapErrorf(awsVersionCmdErr, "expected aws version to start with aws-cli/2.x, got %s", awsVersionCmdOut)
	}

	return nil
}

// TestContainerBaseApkoWithCustomVersions tests the initialization of a target module with custom versions
// for Terragrunt, Terraform, and OpenTofu, with AWS CLI enabled, and verifies the execution of basic shell commands
// within the context of the target module.
//
// This function performs the following steps:
// 1. Initializes the target module with default Terragrunt options, enabling
// AWS CLI and setting custom versions for Terragrunt, Terraform, and OpenTofu.
// 2. Runs the shell command "ls -l" within the context of the target module and checks for errors and expected output.
// 3. Runs the "terragrunt --version" command within the context of the target module and verifies the version output.
// 4. Runs the "terraform --version" command within the context of the target module and verifies the version output.
// 5. Runs the "tofu --version" command within the context of the target module and verifies the version output.
//
// Parameters:
// - ctx: The context for controlling the execution of the function.
//
// Returns:
// - An error if any of the steps fail, otherwise nil.
func (m *Tests) TestContainerBaseApkoWithCustomVersions(ctx context.Context) error {
	// Initialize the target module with default Terragrunt options, enabling AWS CLI.
	targetModule := dag.Terragrunt(dagger.TerragruntOpts{
		EnableAwscli:    true,
		TgVersion:       "0.66.0",
		TfVersion:       "1.4.6",
		OpenTofuVersion: "1.6.3",
	})

	// Run the shell command "ls -l" in the context of the target module.
	targetCtr := targetModule.Ctr()
	simpleCmdOut, simpleCmdErr := targetCtr.WithExec([]string{"ls", "-l"}).Stdout(ctx)

	if simpleCmdErr != nil {
		return WrapError(simpleCmdErr, "failed to run shell command")
	}

	if simpleCmdOut == "" {
		return WrapError(simpleCmdErr, "expected to have at least one folder, got empty output")
	}

	// Run terragrunt --version
	tgVersionCmdOut, tgVersionCmdErr := targetCtr.WithExec([]string{"terragrunt", "--version"}).Stdout(ctx)
	if tgVersionCmdErr != nil {
		return WrapError(tgVersionCmdErr, "failed to run terragrunt --version")
	}

	if !strings.Contains(tgVersionCmdOut, "v0.66.0") {
		return WrapErrorf(tgVersionCmdErr, "expected terragrunt version to be v0.66.0, got %s", tgVersionCmdOut)
	}

	// Run terraform --version
	tfVersionCmdOut, tfVersionCmdErr := targetCtr.WithExec([]string{"terraform", "--version"}).Stdout(ctx)

	if tfVersionCmdErr != nil {
		return WrapError(tfVersionCmdErr, "failed to run terraform --version")
	}

	if !strings.Contains(tfVersionCmdOut, "v1.4.6") {
		return WrapErrorf(tfVersionCmdErr, "expected terraform version to be v1.4.6, got %s", tfVersionCmdOut)
	}

	// Run tofu --version
	tofuVersionCmdOut, tofuVersionCmdErr := targetCtr.WithExec([]string{"opentofu", "--version"}).Stdout(ctx)

	if tofuVersionCmdErr != nil {
		return WrapError(tofuVersionCmdErr, "failed to run opentofu --version")
	}

	if !strings.Contains(tofuVersionCmdOut, "v1.6.3") {
		return WrapErrorf(tofuVersionCmdErr, "expected opentofu version to be v1.6.3, got %s", tofuVersionCmdOut)
	}

	return nil
}
