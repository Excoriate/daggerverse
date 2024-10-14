package main

// Tool represents the entrypoint to use when executing the command.
type Tool string

const (
	// TerragruntTool specifies the tool for Terragrunt.
	TerragruntTool Tool = "terragrunt"
	// TerraformTool specifies the tool for Terraform.
	TerraformTool Tool = "terraform"
	// OpentofuTool specifies the tool for Opentofu.
	OpentofuTool Tool = "opentofu"
)

// IsValidIACTool validates the IaC tool.
func IsValidIACTool(tool string) error {
	validTools := []Tool{TerragruntTool, TerraformTool, OpentofuTool}
	isValidTool := false

	for _, valid := range validTools {
		if tool == string(valid) {
			isValidTool = true

			break
		}
	}

	if !isValidTool {
		return WrapErrorf(nil, "invalid tool: %s. Must be one of: %v", tool, validTools)
	}

	return nil
}
