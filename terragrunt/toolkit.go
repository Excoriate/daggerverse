// Package main provides functionality for managing infrastructure-as-code toolkit versions.
package main

import "github.com/Excoriate/daggerx/pkg/installerx"

// Default versions for OpenTofu, Terraform, and Terragrunt.
const (
	defaultOpenTofuVersion   = "1.8.0"
	defaultTerraformVersion  = "1.9.5"
	defaultTerragruntVersion = "0.67.4"
)

// WithTerragruntInstalled installs the specified version of Terragrunt.
// If no version is specified, it defaults to the version defined in defaultTerragruntVersion.
// The function returns a pointer to the updated Terragrunt instance.
func (m *Terragrunt) WithTerragruntInstalled(
	// version is the version of Terragrunt to install.
	// +optional
	version string,
) *Terragrunt {
	if version == "" {
		version = defaultTerragruntVersion
	}

	installTgCmd := installerx.GetTerragruntInstallCommand(installerx.TerragruntInstallParams{
		Version:    version,
		InstallDir: "/home/terragrunt/bin",
	})

	m.Ctr = m.Ctr.WithExec([]string{"bash", "-c", installTgCmd})

	return m
}

// WithTerraformInstalled installs the specified version of Terraform.
// If no version is specified, it defaults to the version defined in defaultTerraformVersion.
// The function returns a pointer to the updated Terragrunt instance.
func (m *Terragrunt) WithTerraformInstalled(
	// version is the version of Terraform to install.
	// +optional
	version string,
) *Terragrunt {
	if version == "" {
		version = defaultTerraformVersion
	}

	installTfCmd := installerx.GetTerraformInstallCommand(installerx.TerraformInstallParams{
		Version:    version,
		InstallDir: "/home/terragrunt/bin",
	})

	m.Ctr = m.Ctr.WithExec([]string{"bash", "-c", installTfCmd})

	return m
}

// WithOpenTofuInstalled installs the specified version of OpenTofu.
// If no version is specified, it defaults to the version defined in defaultOpenTofuVersion.
// The function returns a pointer to the updated Terragrunt instance.
func (m *Terragrunt) WithOpenTofuInstalled(
	// version is the version of OpenTofu to install.
	// +optional
	version string,
) *Terragrunt {
	if version == "" {
		version = defaultOpenTofuVersion
	}

	installOpenTofuCmd := installerx.GetOpenTofuInstallCommand(installerx.OpenTofuInstallParams{
		Version:    version,
		InstallDir: "/home/terragrunt/bin",
	})

	m.Ctr = m.Ctr.WithExec([]string{"bash", "-c", installOpenTofuCmd})

	return m
}
