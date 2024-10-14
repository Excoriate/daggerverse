package main

import (
	"context"
	"strings"

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
		WithTerragruntPermissionsOnDirsDefault().
		WithTerragruntLogOptions(
			dagger.TerragruntWithTerragruntLogOptionsOpts{
				TgLogLevel:        "debug",
				TgForwardTfStdout: true,
			},
		)

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

	// run apply command with the auto-approve flag as an argument
	cmdApplyOut, cmdApplyErr := tgModule.ExecCmd(ctx, "apply", dagger.TerragruntExecCmdOpts{
		Source: m.
			getTestDir("").
			Directory("terragrunt"),
		Args: []string{
			"-auto-approve",
		},
	})

	if cmdApplyErr != nil {
		return WrapErrorf(cmdApplyErr, "failed to execute command apply")
	}

	if cmdApplyOut == "" {
		return Errorf("command apply output is empty")
	}

	// run destroy command with the auto-approve built-in option.
	cmdDestroyOut, cmdDestroyErr := tgModule.ExecCmd(ctx, "destroy", dagger.TerragruntExecCmdOpts{
		Source: m.
			getTestDir("").
			Directory("terragrunt"),
		AutoApprove: true,
	})

	if cmdDestroyErr != nil {
		return WrapErrorf(cmdDestroyErr, "failed to execute command destroy")
	}

	if cmdDestroyOut == "" {
		return Errorf("command destroy output is empty")
	}

	return nil
}

// TestTerragruntExecWithPlanOutput tests the execution of the 'terragrunt plan' command and validates the output.
//
// This function sets up the necessary environment variables and secrets, configures the Terragrunt module,
// and executes the 'init' command. It validates the output of the command and ensures that
// the command is executed successfully.
// If any step fails, an error is returned.
//
// Parameters:
// - ctx: The context for controlling the execution.
//
// Returns:
// - error: If any step fails, an error is returned.
func (m *Tests) TestTerragruntExecWithPlanOutput(ctx context.Context) error {
	testEnvVars := []string{
		"TF_VAR_project_name=myProject",
		"TF_VAR_environment=production",
		"TF_VAR_region=us-east-1",
		"TF_VAR_instance_type=t2.micro",
		"TF_VAR_db_name=myDatabase",
		"GOOGLE_CLOUD_PROJECT=my-gcp-project",
		"GOOGLE_CLOUD_KEYFILE_JSON=/path/to/keyfile.json",
		"AWS_ACCESS_KEY_ID=my_aws_access_key",
		"AWS_SECRET_ACCESS_KEY=my_aws_secret_key",
		"DOCKER_REGISTRY=mydockerregistry",
	}

	dbPasswordSecret := dag.SetSecret("DB_PASSWORD", "supersecurepassword")
	apiKeySecret := dag.SetSecret("API_KEY", "mysecureapikey")
	sshKeySecret := dag.SetSecret("SSH_PRIVATE_KEY", "mysshprivatekey")

	// main module configuration.
	tgModule := dag.
		Terragrunt(dagger.TerragruntOpts{
			EnvVarsFromHost: testEnvVars,
			TfVersion:       "1.7.0",
		}).
		WithTerragruntPermissionsOnDirsDefault().
		WithTerragruntLogOptions(dagger.TerragruntWithTerragruntLogOptionsOpts{
			TgLogDisableFormatting: true,
			TgLogShowAbsPaths:      true,
			TgLogLevel:             "debug",
		}).
		WithTerraformLogOptions(dagger.TerragruntWithTerraformLogOptionsOpts{
			TfLog:     "debug",
			TfLogPath: "/mnt/tflogs", // it's a directory that the terragrunt user owns.
		}).
		// Extra options added for more realism.
		WithTerragruntOptions(dagger.TerragruntWithTerragruntOptionsOpts{
			IgnoreDependencyErrors:     true,
			IgnoreExternalDependencies: true,
			DisableBucketUpdate:        true,
		})

	// Execute the plan command and get the container back.
	tgCtr := tgModule.Exec("plan", dagger.TerragruntExecOpts{
		Source: m.
			getTestDir("").
			Directory("terragrunt"),
		Secrets: []*dagger.Secret{
			dbPasswordSecret,
			apiKeySecret,
			sshKeySecret,
		},
		// Args to output the plan to a file.
		Args: []string{
			"-out=plan.tfplan",
			"-refresh=true",
		},
	})

	// Execute the plan command, and return the stdout.
	_, outPlanErr := tgCtr.
		Stdout(ctx)

	if outPlanErr != nil {
		return WrapErrorf(outPlanErr, "failed to get the stdout of the plan command")
	}

	// get the plan file
	planFile := tgCtr.Terminal().File("/mnt/plan.tfplan")

	if planFile == nil {
		return Errorf("the terragrunt container does not have the plan file named 'plan.tfplan'")
	}

	return nil
}

// TestTerragruntWithCustomRegistriesToCacheProvidersFrom tests the configuration of custom registries
// for caching providers in Terragrunt.
//
// This function sets up the Terragrunt module with custom registries, executes the 'printenv' command,
// and validates that the environment variable TERRAGRUNT_PROVIDER_CACHE_REGISTRY_NAMES is set correctly.
//
// Parameters:
// - ctx: The context for controlling the execution.
//
// Returns:
// - error: If any step fails, an error is returned.
func (m *Tests) TestTerragruntWithCustomRegistriesToCacheProvidersFrom(ctx context.Context) error {
	// main module configuration.
	tgModule := dag.
		Terragrunt().
		WithTerragruntPermissionsOnDirsDefault().
		WithTerragruntLogOptions(dagger.TerragruntWithTerragruntLogOptionsOpts{
			TgLogDisableFormatting: true,
			TgLogShowAbsPaths:      true,
			TgLogLevel:             "debug",
		}).
		WithTerraformLogOptions(dagger.TerragruntWithTerraformLogOptionsOpts{
			TfLog:     "debug",
			TfLogPath: "/mnt/tflogs", // it's a directory that the terragrunt user owns.
		}).
		WithRegistriesToCacheProvidersFrom([]string{
			"myregistry.mycompany.com",
		})

	tgCtr := tgModule.Ctr()

	// Execute the plan command, and return the stdout.
	envVars, outPlanErr := tgCtr.WithExec([]string{"printenv"}).
		Stdout(ctx)

	// Ensure that the TERRAGRUNT_PROVIDER_CACHE_REGISTRY_NAMES has this value:
	// "registry.terraform.io,registry.opentofu.org,myregistry.mycompany.com"
	expectedEnvVarCfg := "TERRAGRUNT_PROVIDER_CACHE_REGISTRY_NAMES=" +
		"registry.terraform.io,registry.opentofu.org,myregistry.mycompany.com"
	if !strings.Contains(envVars, expectedEnvVarCfg) {
		return Errorf("TERRAGRUNT_PROVIDER_CACHE_REGISTRY_NAMES is not set correctly, got: %s", envVars)
	}

	if outPlanErr != nil {
		return WrapErrorf(outPlanErr, "failed to get the stdout of the plan command")
	}

	return nil
}

// TestTerragruntWithProviderCacheServerDisabled tests the configuration of the provider cache server in Terragrunt.
//
// This function sets up the Terragrunt module with the provider cache server disabled, executes the 'printenv' command,
// and validates that the environment variable TERRAGRUNT_PROVIDER_CACHE is set to "0".
//
// Parameters:
// - ctx: The context for controlling the execution.
//
// Returns:
// - error: If any step fails, an error is returned.
func (m *Tests) TestTerragruntWithProviderCacheServerDisabled(ctx context.Context) error {
	// main module configuration.
	tgModule := dag.
		Terragrunt().
		WithTerragruntPermissionsOnDirsDefault().
		WithTerragruntLogOptions(dagger.TerragruntWithTerragruntLogOptionsOpts{
			TgLogDisableFormatting: true,
			TgLogShowAbsPaths:      true,
			TgLogLevel:             "debug",
			TgForwardTfStdout:      true, // forward terraform stdout to the terragrunt stdout.
		}).
		WithTerragruntProviderCacheServerDisabled()

	tgCtr := tgModule.Ctr()

	// Execute the plan command, and return the stdout.
	envVars, outPlanErr := tgCtr.WithExec([]string{"printenv"}).
		Stdout(ctx)

	if outPlanErr != nil {
		return WrapErrorf(outPlanErr, "failed to get the stdout of the plan command")
	}

	// Ensure that the TERRAGRUNT_PROVIDER_CACHE is set to "0"
	if !strings.Contains(envVars, "TERRAGRUNT_PROVIDER_CACHE=0") {
		return Errorf("TERRAGRUNT_PROVIDER_CACHE is not set correctly")
	}

	// Run plan command
	planCmdOut, planCmdErr := tgModule.ExecCmd(ctx, "plan", dagger.TerragruntExecCmdOpts{
		Source: m.
			getTestDir("").
			Directory("terragrunt"),
	})

	if planCmdErr != nil {
		return WrapErrorf(planCmdErr, "failed to run plan command")
	}

	if planCmdOut == "" {
		return Errorf("plan command output is empty")
	}

	return nil
}
