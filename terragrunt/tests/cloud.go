package main

import (
	"context"
	"strings"
)

// TestWithAWSKeys tests the setting of AWS keys as environment variables within the target module's container.
//
// This method creates secrets for AWS Access Key ID and AWS Secret Access Key, sets these secrets
// as environment variables in the target module's container, and verifies if the expected environment
// variables are set by running the `printenv` command within the container.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if there is an issue creating secrets, setting environment variables,
//     executing commands in the container, or if the output does not contain the expected environment variables.
func (m *Tests) TestWithAWSKeys(ctx context.Context) error {
	targetModule := dag.Terragrunt()

	awsKeyID := dag.SetSecret("AWS_ACCESS_KEY_ID", "AKIAIOSFODNN7EXAMPLE")
	awsSecretAccessKey := dag.SetSecret("AWS_SECRET_ACCESS_KEY", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")

	// With required AWS keys only.
	targetModule = targetModule.
		WithAwskeys(awsKeyID, awsSecretAccessKey)

	out, err := targetModule.Ctr().
		WithExec([]string{"printenv"}).
		Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to get AWS Keys environment variables")
	}

	if !strings.Contains(out, "AWS_ACCESS_KEY_ID") {
		return WrapErrorf(err, "expected AWS_ACCESS_KEY_ID to be set, got %s", out)
	}

	if !strings.Contains(out, "AWS_SECRET_ACCESS_KEY") {
		return WrapErrorf(err, "expected AWS_SECRET_ACCESS_KEY to be set, got %s", out)
	}

	// Check if the content of the env vars is correct.
	if !strings.Contains(out, "AKIAIOSFODNN7EXAMPLE") {
		return WrapErrorf(err, "expected AKIAIOSFODNN7EXAMPLE to be set, got %s", out)
	}

	if !strings.Contains(out, "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY") {
		return WrapErrorf(err, "expected wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY to be set, got %s", out)
	}

	return nil
}

// TestWithAzureCredentials tests the setting of Azure credentials as
// environment variables within the target module's container.
//
// This method creates secrets for Azure Client ID, Azure Client Secret, and Azure Tenant ID,
// sets these secrets as environment variables in the target module's container, and verifies if the expected
// environment variables are set by running the `printenv` command within the container.
//
// Arguments:
// - ctx (context.Context): The context for the test execution.
//
// Returns:
//   - error: Returns an error if there is an issue creating secrets, setting environment variables,
//     executing commands in the container, or if the output does not contain the expected environment variables.
func (m *Tests) TestWithAzureCredentials(ctx context.Context) error {
	targetModule := dag.Terragrunt()

	azureClientID := dag.SetSecret("AZURE_CLIENT_ID", "00000000-0000-0000-0000-000000000000")
	azureClientSecret := dag.SetSecret("AZURE_CLIENT_SECRET", "00000000-0000-0000-0000-000000000000")
	azureTenantID := dag.SetSecret("AZURE_TENANT_ID", "00000000-0000-0000-0000-000000000000")

	// With required Azure keys only.
	targetModule = targetModule.
		WithAzureCredentials(azureClientID, azureClientSecret, azureTenantID)

	out, err := targetModule.Ctr().
		WithExec([]string{"printenv"}).
		Stdout(ctx)

	if err != nil {
		return WrapError(err, "failed to get Azure Keys environment variables")
	}

	if !strings.Contains(out, "AZURE_CLIENT_ID") {
		return WrapErrorf(err, "expected AZURE_CLIENT_ID to be set, got %s", out)
	}

	if !strings.Contains(out, "AZURE_CLIENT_SECRET") {
		return WrapErrorf(err, "expected AZURE_CLIENT_SECRET to be set, got %s", out)
	}

	if !strings.Contains(out, "AZURE_TENANT_ID") {
		return WrapErrorf(err, "expected AZURE_TENANT_ID to be set, got %s", out)
	}

	return nil
}
