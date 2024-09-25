package main

import (
	"fmt"
	"strings"

	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
)

const (
	terragruntReleaseURL = "https://github.com/gruntwork-io/terragrunt/releases/download"
	terraformReleaseURL  = "https://releases.hashicorp.com/terraform"
	openTofuReleaseURL   = "https://github.com/opentofu/opentofu/releases/download"
)

// InstallTool installs the specified tool (Terragrunt, Terraform, or OpenTofu) in the container.
// It returns the modified container.
//
// Parameters:
//   - ctr: The input container
//   - tool: The tool to install ("terragrunt", "terraform", or "opentofu")
//   - version: The version to install (use "latest" for the latest version)
//
// Returns:
//   - *dagger.Container: The modified container with the tool installed
//   - error: Any error that occurred during the installation process
func InstallTool(ctr *dagger.Container, tool, version string) (*dagger.Container, error) {
	version = strings.TrimPrefix(version, "v")
	if version == "latest" {
		var err error
		version, err = getLatestVersion(tool)
		if err != nil {
			return nil, fmt.Errorf("failed to get latest version for %s: %v", tool, err)
		}
	}

	// Switch to root user for installation
	ctr = ctr.WithUser("root")

	// Ensure /usr/local/bin exists
	ctr = ctr.WithExec([]string{"mkdir", "-p", "/usr/local/bin"})

	switch tool {
	case "terragrunt":
		url := fmt.Sprintf("%s/v%s/terragrunt_linux_amd64", terragruntReleaseURL, version)
		ctr = ctr.
			WithExec([]string{"curl", "-L", "-o", "/usr/local/bin/terragrunt", url}).
			WithExec([]string{"chmod", "+x", "/usr/local/bin/terragrunt"})
	case "terraform":
		url := fmt.Sprintf("%s/%s/terraform_%s_linux_amd64.zip", terraformReleaseURL, version, version)
		ctr = ctr.
			WithExec([]string{"curl", "-L", "-o", "/tmp/terraform.zip", url}).
			WithExec([]string{"unzip", "-d", "/usr/local/bin", "/tmp/terraform.zip"}).
			WithExec([]string{"chmod", "+x", "/usr/local/bin/terraform"}).
			WithExec([]string{"rm", "/tmp/terraform.zip"})
	case "opentofu":
		url := fmt.Sprintf("%s/v%s/tofu_%s_linux_amd64.zip", openTofuReleaseURL, version, version)
		ctr = ctr.
			WithExec([]string{"curl", "-L", "-o", "/tmp/tofu.zip", url}).
			WithExec([]string{"unzip", "-d", "/usr/local/bin", "/tmp/tofu.zip"}).
			WithExec([]string{"chmod", "+x", "/usr/local/bin/tofu"}).
			WithExec([]string{"rm", "/tmp/tofu.zip"})
	default:
		return nil, fmt.Errorf("unsupported tool: %s", tool)
	}

	// Switch back to the default user
	ctr = ctr.WithUser("")

	// Verify installation
	ctr = ctr.WithExec([]string{tool, "--version"})

	return ctr, nil
}

// getLatestVersion fetches the latest version for the specified tool.
// This is a placeholder function. In a real scenario, you'd implement logic to fetch the latest version.
func getLatestVersion(tool string) (string, error) {
	// Placeholder: In reality, you'd fetch this from the tool's releases page or API
	switch tool {
	case "terragrunt":
		return "0.67.4", nil
	case "terraform":
		return "1.9.4", nil
	case "opentofu":
		return "1.5.0", nil
	default:
		return "", fmt.Errorf("unsupported tool: %s", tool)
	}
}
