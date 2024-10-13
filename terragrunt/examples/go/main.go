// A generated module for Go example functions
package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Excoriate/daggerverse/terragrunt/examples/go/internal/dagger"
	"github.com/sourcegraph/conc/pool"
)

// Go is a Dagger module that exemplifies the usage of the Terragrunt module.
//
// This module is used to create and manage containers.
type Go struct {
	TestDir *dagger.Directory
}

// New creates a new Tests instance.
//
// It's the initial constructor for the Tests struct.
func New() *Go {
	e := &Go{}

	e.TestDir = e.getTestDir()

	return e
}

var (
	errEnvVarsEmpty            = errors.New("env vars are empty")
	errEnvVarsDontMatch        = errors.New("expected env vars to be passed, got empty output")
	errNetRCFileNotFound       = errors.New("netrc file not found")
	errExpectedFoldersNotFound = errors.New("expected to have at least one folder, got empty output")
	errTerragruntPlanEmpty     = errors.New("terragrunt plan command output is empty")
	errCommandInitEmpty        = errors.New("command init output is empty")
	errCommandPlanEmpty        = errors.New("command plan output is empty")
	errCommandApplyEmpty       = errors.New("command apply output is empty")
	errCommandDestroyEmpty     = errors.New("command destroy output is empty")
)

// getTestDir returns the test directory.
//
// This helper method retrieves the test directory, which is typically located
// in the same directory as the test file and named "testdata".
//
// Returns:
//   - *dagger.Directory: A Dagger Directory object pointing to the "testdata" directory.
func (m *Go) getTestDir() *dagger.Directory {
	return dag.
		CurrentModule().
		Source().
		Directory("./testdata")
}

// AllRecipes executes all tests.
//
// AllRecipes is a helper method for tests, executing the built-in recipes and
// other specific functionalities of the Terragrunt module.
func (m *Go) AllRecipes(ctx context.Context) error {
	polTests := pool.New().WithErrors().WithContext(ctx)

	// Test different ways to configure the base container.
	polTests.Go(m.Terragrunt_PassedEnvVars)
	// From this point onwards, we're testing the specific functionality of the Terragrunt module.
	polTests.Go(m.Terragrunt_ExecPlanCommand)

	if err := polTests.Wait(); err != nil {
		return fmt.Errorf("there are some failed tests: %w", err)
	}

	return nil
}

// // BuiltInRecipes demonstrates how to run built-in recipes
// //
// // This method showcases the use of various built-in recipes provided by the Terragrunt
// // module, including creating a container, running an arbitrary command, and creating a .netrc
// // file for GitHub authentication.
// //
// // Parameters:
// //   - ctx: The context for controlling the function's timeout and cancellation.
// //
// // Returns:
// //   - An error if any of the internal methods fail, or nil otherwise.
// func (m *Go) BuiltInRecipes(ctx context.Context) error {
// 	// Pass environment variables to the Terragrunt module using Terragrunt_PassedEnvVars
// 	if err := m.Terragrunt_PassedEnvVars(ctx); err != nil {
// 		return fmt.Errorf("failed to pass environment variables: %w", err)
// 	}

// 	// Create a configured container using Terragrunt_CreateContainer
// 	if _, err := m.Terragrunt_CreateContainer(ctx); err != nil {
// 		return fmt.Errorf("failed to create container: %w", err)
// 	}

// 	// Run an arbitrary command in the container using Terragrunt_RunArbitraryCommand
// 	if _, err := m.Terragrunt_RunArbitraryCommand(ctx); err != nil {
// 		return fmt.Errorf("failed to run arbitrary command: %w", err)
// 	}

// 	// Create a .netrc file for GitHub using Terragrunt_CreateNetRcFileForGithub
// 	if _, err := m.Terragrunt_CreateNetRcFileForGithub(ctx); err != nil {
// 		return fmt.Errorf("failed to create netrc file: %w", err)
// 	}

// 	return nil
// }

// Terragrunt_PassedEnvVars demonstrates how to pass environment variables to the Terragrunt module.
//
// This method configures a Terragrunt module to use specific environment variables from the host.
func (m *Go) Terragrunt_PassedEnvVars(ctx context.Context) error {
	targetModule := dag.Terragrunt(dagger.TerragruntOpts{
		EnvVarsFromHost: []string{"SOMETHING=SOMETHING,SOMETHING=SOMETHING"},
	})

	out, err := targetModule.Ctr().
		WithExec([]string{"printenv"}).
		Stdout(ctx)

	if err != nil {
		return fmt.Errorf("failed when executing printenv: %w", err)
	}

	if out == "" {
		return errEnvVarsEmpty
	}

	if !strings.Contains(out, "SOMETHING") {
		return errEnvVarsDontMatch
	}

	return nil
}

// Terragrunt_OpenTerminal demonstrates how to open an interactive terminal session
// within a Terragrunt module container.
//
// This function showcases the initialization and configuration of a
// Terragrunt module container using various options like enabling Cgo,
// utilizing build cache, and including a GCC compiler.
//
// Parameters:
//   - None
//
// Returns:
//   - *dagger.Container: A configured Dagger container with an open terminal.
//
// Usage:
//
//	This function can be used to interactively debug or inspect the
//	container environment during test execution.
func (m *Go) Terragrunt_OpenTerminal() *dagger.Container {
	// Configure the Terragrunt module container with necessary options
	targetModule := dag.Terragrunt()

	// Retrieve and discard standard output
	_, _ = targetModule.Ctr().
		Stdout(context.Background())

	// Open and return the terminal session in the configured container
	return targetModule.Ctr().
		Terminal()
}

// Terragrunt_RunArbitraryCommand runs an arbitrary shell command in the test container.
//
// This function demonstrates how to execute a shell command within the container
// using the Terragrunt module.
//
// Parameters:
//
//	ctx - context for controlling the function lifetime.
//
// Returns:
//
//	A string containing the output of the executed command, or an error if the command fails or if the output is empty.
func (m *Go) Terragrunt_RunArbitraryCommand(ctx context.Context) (string, error) {
	targetModule := dag.Terragrunt().WithSource(m.TestDir)

	// Execute the 'ls -l' command
	out, err := targetModule.
		Ctr().
		WithExec([]string{"ls", "-l"}).
		Stdout(ctx)

	if err != nil {
		return "", fmt.Errorf("failed to run shell command: %w", err)
	}

	if out == "" {
		return "", errExpectedFoldersNotFound
	}

	return out, nil
}

// Terragrunt_ExecPlanCommand executes the 'terragrunt plan' command with specific configurations.
//
// This function demonstrates how to:
// - Set up environment variables for Terragrunt execution
// - Initialize the Terragrunt module with advanced logging and token options
// - Execute the 'plan' command
// - Validate the command output and environment variable setup
//
// It uses the Dagger Terragrunt module to perform these operations within a container.
//
// Parameters:
// - ctx: context.Context for controlling the function's lifetime
//
// Returns:
//   - error: If any step in the process fails, including command execution,
//     output validation, or environment variable checks
func (m *Go) Terragrunt_ExecPlanCommand(ctx context.Context) error {
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
				getTestDir().
				Directory("terragrunt"),
		})

	tgPlanCmdOut, tgPlanCmdErr := tgCtrConfigured.
		Stdout(ctx)

	if tgPlanCmdErr != nil {
		return fmt.Errorf("failed to execute terragrunt plan command: %w", tgPlanCmdErr)
	}

	if tgPlanCmdOut == "" {
		return errTerragruntPlanEmpty
	}

	return nil
}

// Terragrunt_ExecLifecycleCommands executes Terragrunt lifecycle commands.
//
// This function sets up the necessary environment variables and secrets, configures the Terragrunt module,
// and executes a series of Terragrunt commands ("init", "plan", "apply", "destroy"). It validates the output
// of each command and ensures that the commands are executed successfully.
//
// Parameters:
// - ctx: The context for controlling the execution.
//
// Returns:
// - error: If any command execution fails or if the command output is empty.
func (m *Go) Terragrunt_ExecLifecycleCommands(ctx context.Context) error {
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
			getTestDir().
			Directory("terragrunt"),
		Secrets: []*dagger.Secret{
			awsSecret,
			gcpSecret,
			azureSecret,
		},
	})

	if cmdInitErr != nil {
		return fmt.Errorf("failed to execute command init: %w", cmdInitErr)
	}

	if cmdInitOut == "" {
		return errCommandInitEmpty
	}

	// run plan command with arguments
	cmdPlanOut, cmdPlanErr := tgModule.ExecCmd(ctx, "plan", dagger.TerragruntExecCmdOpts{
		Source: m.
			getTestDir().
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
		return fmt.Errorf("failed to execute command plan: %w", cmdPlanErr)
	}

	if cmdPlanOut == "" {
		return errCommandPlanEmpty
	}

	// run apply command with the auto-approve flag as an argument
	cmdApplyOut, cmdApplyErr := tgModule.ExecCmd(ctx, "apply", dagger.TerragruntExecCmdOpts{
		Source: m.
			getTestDir().
			Directory("terragrunt"),
		Args: []string{
			"-auto-approve",
		},
	})

	if cmdApplyErr != nil {
		return fmt.Errorf("failed to execute command apply: %w", cmdApplyErr)
	}

	if cmdApplyOut == "" {
		return errCommandApplyEmpty
	}

	// run destroy command with the auto-approve built-in option.
	cmdDestroyOut, cmdDestroyErr := tgModule.ExecCmd(ctx, "destroy", dagger.TerragruntExecCmdOpts{
		Source: m.
			getTestDir().
			Directory("terragrunt"),
		AutoApprove: true,
	})

	if cmdDestroyErr != nil {
		return fmt.Errorf("failed to execute command destroy: %w", cmdDestroyErr)
	}

	if cmdDestroyOut == "" {
		return errCommandDestroyEmpty
	}

	return nil
}
