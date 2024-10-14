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
		WithTerragruntPermissionsOnDirsDefault().
		WithTerragruntLogOptions(
			dagger.TerragruntWithTerragruntLogOptionsOpts{
				TgLogLevel:        "debug",
				TgForwardTfStdout: true,
			},
		)

	// Execute the init command, but don't run it in a container
	tfInitCtr := tgModule.
		TfExec("init", dagger.TerragruntTfExecOpts{
			Source: m.getTestDir("").
				Directory("terraform"),
		})

	tfPlanCtr := tgModule.WithContainer(tfInitCtr).
		TfExec("plan")

	// Evaluate the terraform init command.
	tfPlanOut, tfPlanErr := tfPlanCtr.WithExec([]string{"ls"}).Terminal().
		Stdout(ctx)

	if tfPlanErr != nil {
		return WrapErrorf(tfPlanErr, "failed to get terraform plan command output")
	}

	if tfPlanOut == "" {
		return Errorf("terraform plan command output is empty")
	}

	// Check the environment variables set in the container
	for _, envVar := range testEnvVars {
		if err := m.assertEnvVarIsSetInContainer(ctx, tfPlanCtr, envVar); err != nil {
			return err
		}
	}

	return nil
}
