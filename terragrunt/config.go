package main

// Entrypoint represents the entrypoint to use when executing the command.
type Entrypoint string

const (
	// TerragruntEntrypoint specifies the entrypoint for Terragrunt.
	TerragruntEntrypoint Entrypoint = "terragrunt"
	// TerraformEntrypoint specifies the entrypoint for Terraform.
	TerraformEntrypoint Entrypoint = "terraform"
	// OpentofuEntrypoint specifies the entrypoint for Opentofu.
	OpentofuEntrypoint Entrypoint = "opentofu"
)

// IsValidIACTool validates the IaC tool.
func IsValidIACTool(tool string) error {
	validTools := []Entrypoint{TerragruntEntrypoint, TerraformEntrypoint, OpentofuEntrypoint}
	isValidTool := false

	for _, valid := range validTools {
		if tool == string(valid) {
			isValidTool = true

			break
		}
	}

	if !isValidTool {
		return WrapErrorf(nil, "invalid entrypoint: %s. Must be one of: %v", tool, validTools)
	}

	return nil
}
