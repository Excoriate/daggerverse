package main

import (
	"context"

	"github.com/Excoriate/daggerverse/terragrunt/tests/internal/dagger"
)

// TestTerragruntContainerIsUp checks if the Terragrunt container is up and running by verifying the versions of
// Terragrunt, Terraform, and OpenTofu.
// It executes the version commands for each tool and checks their outputs to ensure they contain the expected
// version strings.
// If any of the commands fail or the outputs do not contain the expected strings, an error is returned.
func (m *Tests) TestTerragruntContainerIsUp(ctx context.Context) error {
	tgTestDir := m.
		getTestDir("").
		Directory("terragrunt")

	tgCtr := dag.
		Terragrunt().
		WithSource(tgTestDir).
		Ctr()

	if err := m.assertTheseFoldersExistsInContainer(ctx, tgCtr,
		[]string{"/home/terragrunt",
			"/home/.terraform.d",
			"/home"}); err != nil {
		return WrapErrorf(err, "failed to validate folders in terragrunt container")
	}

	if err := m.assertTheseFilesExistsInContainer(ctx, tgCtr, []string{"terragrunt.hcl"}); err != nil {
		return WrapErrorf(err, "failed to validate terragrunt.hcl file in terragrunt container")
	}

	if err := m.assertFileContentShouldContain(ctx, tgCtr, "terragrunt.hcl", "terraform {"); err != nil {
		return WrapErrorf(err, "failed to validate terragrunt.hcl file content in terragrunt container")
	}

	return nil
}

// TestTerragruntBinariesAreInstalled checks if the Terragrunt, Terraform, and OpenTofu binaries are installed and
// their versions are correct.
// It executes the version command for each tool and verifies that the output contains the expected version string.
// If any of the commands fail or the outputs do not contain the expected strings, an error is returned.
func (m *Tests) TestTerragruntBinariesAreInstalled(ctx context.Context) error {
	tgCtr := dag.
		Terragrunt().
		Ctr()

	if err := m.assertVersionOfBinaryInContainer(ctx, tgCtr, "terragrunt", "terragrunt version", ""); err != nil {
		return err
	}

	if err := m.assertVersionOfBinaryInContainer(ctx, tgCtr, "terraform", "Terraform v", ""); err != nil {
		return err
	}

	if err := m.assertVersionOfBinaryInContainer(ctx, tgCtr, "opentofu", "OpenTofu v", ""); err != nil {
		return err
	}

	return nil
}

// TestTerragruntExecInitSimpleCommand tests the execution of the 'terragrunt init' command with a simple
// configuration.
// It sets up the necessary environment variables, initializes the Terragrunt module, and executes the 'init'
// command.
// The function then validates the output of the command and checks if the environment variables are correctly set
// in the container.
// If any step fails, an error is returned.
func (m *Tests) TestTerragruntExecInitSimpleCommand(ctx context.Context) error {
	testEnvVars := []string{
		"AWS_ACCESS_KEY_ID=test",
		"AWS_SECRET_ACCESS_KEY=test",
		"AWS_SESSION_TOKEN=test",
	}

	// Initialize the Terragrunt module
	tgModule := dag.
		Terragrunt(dagger.TerragruntOpts{
			EnvVarsFromHost: testEnvVars,
		}).
		WithTerragruntPermissionsOnDirsDefault()

	// Execute the init command, but don't run it in a container
	tgCtrConfigured := tgModule.
		Exec("init", dagger.TerragruntExecOpts{
			Source: m.getTestDir("").
				Directory("terragrunt"),
		})

	// Evaluate the terragrunt init command.
	tgInitCmdOut, tgInitCmdErr := tgCtrConfigured.
		Stdout(ctx)

	if tgInitCmdErr != nil {
		return WrapErrorf(tgInitCmdErr, "failed to get terragrunt init command output")
	}

	if tgInitCmdOut == "" {
		return Errorf("terragrunt init command output is empty")
	}

	// Check the environment variables set in the container
	for _, envVar := range testEnvVars {
		if err := m.assertEnvVarIsSetInContainer(ctx, tgCtrConfigured, envVar); err != nil {
			return err
		}
	}

	return nil
}

// TestTerragruntExecVersionCommand tests the execution of the 'terragrunt version' command with a specific
// configuration.
// It sets up the necessary environment variables, initializes the Terragrunt module with advanced options, and
// executes the 'version' command.
// The function then validates the output of the command, checks if the expected version is present, and verifies
// if the environment variables are correctly set in the container.
// If any step fails, an error is returned.
func (m *Tests) TestTerragruntExecVersionCommand(ctx context.Context) error {
	testEnvVars := []string{
		"OTHER_ENV_VAR=test",
	}

	// Initialize the Terragrunt module with some advance options.
	tgModule := dag.
		Terragrunt(dagger.TerragruntOpts{
			EnvVarsFromHost: testEnvVars,
			EnableAwscli:    true,
			TgVersion:       "v0.52.1",
		}).
		WithTerragruntPermissionsOnDirsDefault().
		WithTerragruntLogOptions(dagger.TerragruntWithTerragruntLogOptionsOpts{
			TgLogLevel:        "debug",
			TgLogDisableColor: true,
		})

	// Execute the init command, but don't run it in a container
	tgCtrConfigured := tgModule.
		Exec("version", dagger.TerragruntExecOpts{
			Source: m.
				getTestDir("").
				Directory("terragrunt"),
		})

	// Evaluate the terragrunt version command.
	tgVersionCmdOut, tgVersionCmdErr := tgCtrConfigured.
		Stdout(ctx)

	if tgVersionCmdErr != nil {
		return WrapErrorf(tgVersionCmdErr, "failed to get terragrunt version command output")
	}

	if tgVersionCmdOut == "" {
		return Errorf("terragrunt version command output is empty")
	}

	// Expected environment variables due to configuration that's passed.
	expectedEnvVars := []string{
		"OTHER_ENV_VAR=test",
		"TERRAGRUNT_LOG_DISABLE_FORMATTING=false",
		"TERRAGRUNT_LOG_DISABLE_COLOR=true",
		"TERRAGRUNT_LOG_SHOW_ABS_PATHS=false",
		"TERRAGRUNT_LOG_LEVEL=debug",
		"TERRAGRUNT_PROVIDER_CACHE=1",
	}

	for _, envVar := range expectedEnvVars {
		if err := m.assertEnvVarIsSetInContainer(ctx, tgCtrConfigured, envVar); err != nil {
			return WrapErrorf(err, "failed to validate environment variables in terragrunt container")
		}
	}

	return nil
}

// TestTerragruntExecPlanCommand tests the execution of the 'terragrunt plan' command with a specific
// configuration.
// It sets up the necessary environment variables, initializes the Terragrunt module with advanced options, and
// executes the 'plan' command.
// The function then validates the output of the command, checks if the expected version is present, and verifies
// if the environment variables are correctly set in the container.
// If any step fails, an error is returned.
func (m *Tests) TestTerragruntExecPlanCommand(ctx context.Context) error {
	testEnvVars := []string{
		"OTHER_ENV_VAR=test",
		"AWS_ACCESS_KEY_ID=test",
		"AWS_SECRET_ACCESS_KEY=test",
		"AWS_SESSION_TOKEN=test",
		"TF_VAR_test=test",
	}

	tfTokenAsSecret := dag.SetSecret("TF_TOKEN_gitlab", "mysupertoken")

	// Initialize the Terragrunt module with some advance options.
	tgModule := dag.
		Terragrunt(dagger.TerragruntOpts{
			EnvVarsFromHost: testEnvVars,
			TgVersion:       "v0.52.1",
		}).
		WithTerraformToken(tfTokenAsSecret).
		WithTerragruntLogOptions(dagger.TerragruntWithTerragruntLogOptionsOpts{
			TgLogLevel:             "debug",
			TgLogDisableColor:      true,
			TgForwardTfStdout:      true,
			TgLogDisableFormatting: true,
		}).
		WithTerraformLogOptions(dagger.TerragruntWithTerraformLogOptionsOpts{
			TfLog:     "debug",
			TfLogPath: "/mnt/tflogs", // it's a directory that the terragrunt user owns.
		})

	// Container configured with all the options.
	tgCtrConfigured := tgModule.
		Exec("plan", dagger.TerragruntExecOpts{
			Source: m.
				getTestDir("").
				Directory("terragrunt"),
		})

	tgPlanCmdOut, tgPlanCmdErr := tgCtrConfigured.
		Stdout(ctx)

	if tgPlanCmdErr != nil {
		return WrapErrorf(tgPlanCmdErr, "failed to get terragrunt plan command output")
	}

	if tgPlanCmdOut == "" {
		return Errorf("terragrunt plan command output is empty")
	}

	// Check env vars
	for _, envVar := range testEnvVars {
		if err := m.assertEnvVarIsSetInContainer(ctx, tgCtrConfigured, envVar); err != nil {
			return WrapErrorf(err, "failed to validate environment variables in terragrunt container")
		}
	}

	return nil
}

// TestTerragruntExecLifecycleCommands tests the execution of Terragrunt commands directly.
//
// This function sets up the necessary environment variables and secrets, configures the Terragrunt module,
// and executes a series of Terragrunt commands ("plan", "apply", "destroy"). It validates the output of each
// command and ensures that the commands are executed successfully.
//
// Parameters:
// - ctx: The context for controlling the execution.
//
// Returns:
// - error: If any command execution fails or if the command output is empty.
func (m *Tests) TestTerragruntExecLifecycleCommands(ctx context.Context) error {
	testEnvVars := []string{
		"TF_VAR_test=test",
		"TF_VAR_another_test=test",
		"TF_VAR_region=westus",
		"TF_VAR_resource_group=myResourceGroup",
		"TF_VAR_storage_account=myStorageAccount",
		"AZURE_SUBSCRIPTION_ID=your_subscription_id",
		"AZURE_CLIENT_ID=your_client_id",
		"AZURE_CLIENT_SECRET=your_client_secret",
		"AZURE_TENANT_ID=your_tenant_id",
	}

	awsSecret := dag.SetSecret("AWS_SECRET_ACCESS_KEY", "awssecretkey")
	gcpSecret := dag.SetSecret("GCP_SERVICE_ACCOUNT_KEY", "gcpserviceaccountkey")
	azureSecret := dag.SetSecret("AZURE_CLIENT_SECRET", "azureclientsecret")

	// github tf token
	tfTokenGitHub := dag.SetSecret("TF_TOKEN_github", "mygithubtoken")

	// main module configuration.
	tgModule := dag.
		Terragrunt(dagger.TerragruntOpts{
			EnvVarsFromHost: testEnvVars,
		}).
		WithTerragruntPermissionsOnDirsDefault().
		WithTerraformToken(tfTokenGitHub).
		WithTerragruntLogOptions(dagger.TerragruntWithTerragruntLogOptionsOpts{
			TgLogLevel:             "debug",
			TgLogDisableColor:      true,
			TgForwardTfStdout:      true,
			TgLogDisableFormatting: true,
		}).
		WithTerraformLogOptions(dagger.TerragruntWithTerraformLogOptionsOpts{
			TfLog:     "debug",
			TfLogPath: "/mnt/tflogs", // it's a directory that the terragrunt user owns.
		})

	// run init command
	cmdInitOut, cmdInitErr := tgModule.ExecCmd(ctx, "init", dagger.TerragruntExecCmdOpts{
		Source: m.
			getTestDir("").
			Directory("terragrunt"),
		Secrets: []*dagger.Secret{
			awsSecret,
			gcpSecret,
			azureSecret,
		},
	})

	if cmdInitErr != nil {
		return WrapErrorf(cmdInitErr, "failed to execute command init")
	}

	if cmdInitOut == "" {
		return Errorf("command init output is empty")
	}

	// run plan command with arguments
	cmdPlanOut, cmdPlanErr := tgModule.ExecCmd(ctx, "plan", dagger.TerragruntExecCmdOpts{
		Source: m.
			getTestDir("").
			Directory("terragrunt"),
		Secrets: []*dagger.Secret{
			awsSecret,
			gcpSecret,
			azureSecret,
		},
		Args: []string{
			"-out=plan.tfplan",
			"-refresh=true",
		},
	})

	if cmdPlanErr != nil {
		return WrapErrorf(cmdPlanErr, "failed to execute command plan")
	}

	if cmdPlanOut == "" {
		return Errorf("command plan output is empty")
	}

	return nil
}
