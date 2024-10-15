package main

import (
	"context"

	"github.com/Excoriate/daggerverse/terragrunt/tests/internal/dagger"
)

// TestTfExecInitSimpleCommand tests the execution of the 'terraform init' command with a simple
// configuration.
// It sets up the necessary environment variables, initializes the Terragrunt module, and executes the 'init'
// command.
// The function then validates the output of the command and checks if the environment variables are correctly set
// in the container.
// If any step fails, an error is returned.
func (m *Tests) TestTfExecInitSimpleCommand(ctx context.Context) error {
	testEnvVars := []string{
		"AWS_ACCESS_KEY_ID=AKIAEXAMPLE",
		"AWS_SECRET_ACCESS_KEY=secretKey12345",
		"AWS_SESSION_TOKEN=sessionToken67890",
		"AWS_REGION=us-west-2",
	}

	// Initialize the Terragrunt module
	tgModule := dag.
		Terragrunt(dagger.TerragruntOpts{
			EnvVarsFromHost: testEnvVars,
			TfVersion:       "1.9.1",
		}).
		WithSource(m.getTestDir("").
			Directory("terraform/tf-module-1"),
		).
		WithTerragruntLogOptions(
			dagger.TerragruntWithTerragruntLogOptionsOpts{
				TgLogLevel:        "debug",
				TgForwardTfStdout: true,
			},
		).
		WithTerragruntPermissionsOnDirsDefault().
		WithTerraformCommand("init").
		WithTerraformCommand("plan").
		WithTerraformCommand("apply", dagger.TerragruntWithTerraformCommandOpts{
			AutoApprove: true,
		})

	// Execute the init command, but don't run it in a container
	tfCtr := tgModule.Ctr()

	// // Evaluate the terraform init command.
	tfPlanOut, tfPlanErr := tfCtr.
		Stdout(ctx)

	if tfPlanErr != nil {
		return WrapErrorf(tfPlanErr, "failed to get terraform plan command output")
	}

	if tfPlanOut == "" {
		return Errorf("terraform plan command output is empty")
	}

	// Check the environment variables set in the container
	for _, envVar := range testEnvVars {
		if err := m.assertEnvVarIsSetInContainer(ctx, tfCtr, envVar); err != nil {
			return err
		}
	}

	return nil
}
