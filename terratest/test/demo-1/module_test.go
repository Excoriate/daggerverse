package terratest_tf_module_1

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTerraformHelloWorldExample(t *testing.T) {
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../tf-module-1",
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndPlan(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	randomId := terraform.Output(t, terraformOptions, "random_id")
	randomPassword := terraform.Output(t, terraformOptions, "random_password")
	assert.NotEmptyf(t, randomId, "random_id should not be empty")
	assert.NotEmptyf(t, randomPassword, "random_password should not be empty")
}
