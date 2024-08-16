package main

import (
	"fmt"
	"strings"

	"github.com/Excoriate/daggerverse/module-template/internal/dagger"
)

const (
	terragruntLatestVersion          = "latest"
	terragruntReleaseURL             = "https://github.com/gruntwork-io/terragrunt/releases/download"
	terragruntDefaultVer             = "v0.66.8"
	terraformReleaseURLForTerragrunt = "https://releases.hashicorp.com/terraform"
	terraformDefaultVer              = "1.9.4"
)

// resolveTerragruntVersion resolves the provided version to "latest" if it's empty.
func resolveTerragruntVersion(version string) string {
	if version == "" || version == "latest" {
		return getLatestTerragruntVersion()
	}
	return version
}

// resolveTerraformVersionInTerragrunt resolves the provided version to "latest" if it's empty.
func resolveTerraformVersionInTerragrunt(version string) string {
	if version == "" {
		return terragruntLatestVersion
	}
	return version
}

// getLatestTerragruntVersion fetches the latest Terragrunt version.
// This is a placeholder function. In a real scenario, you'd implement logic to fetch the latest version.
func getLatestTerragruntVersion() string {
	// Placeholder: In reality, you'd fetch this from Terragrunt's releases page or API
	return terragruntDefaultVer
}

// getLatestTerraformVersionInTerragrunt fetches the latest Terraform version.
// This is a placeholder function. In a real scenario, you'd implement logic to fetch the latest version.
func getLatestTerraformVersionInTerragrunt() string {
	// Placeholder: In reality, you'd fetch this from Terraform's releases page or API
	return terraformDefaultVer
}

// WithTerragruntUbuntu sets up the container with Terragrunt and optionally Terraform on Ubuntu.
// It updates the package list, installs required dependencies, and then installs the specified versions of Terragrunt and Terraform.
//
// Parameters:
// version - The version of Terragrunt to install. If empty, "latest" will be installed.
// tfVersion - The version of Terraform to install. If empty, "latest" will be installed.
// skipTerraform - If true, Terraform installation will be skipped.
func (m *ModuleTemplate) WithTerragruntUbuntu(
	// version is the version of Terragrunt to install. If empty, it will be installed as "latest".
	// +optional
	version string,
	// tfVersion is the version of Terraform to install. If empty, it will be installed as "latest".
	// +optional
	tfVersion string,
	// skipTerraform if true, Terraform installation will be skipped.
	// +optional
	skipTerraform bool,
) *ModuleTemplate {
	version = resolveTerragruntVersion(version)
	if version == terragruntLatestVersion {
		version = getLatestTerragruntVersion()
	}

	tfVersion = resolveTerraformVersionInTerragrunt(tfVersion)
	if tfVersion == terragruntLatestVersion {
		tfVersion = getLatestTerraformVersionInTerragrunt()
	}

	m.Ctr = m.downloadRequiredUtilitiesTerragruntUbuntu()
	m.Ctr = m.downloadAndInstallTerragrunt(version)

	if !skipTerraform {
		m.Ctr = m.downloadAndInstallTerraformInTerragrunt(tfVersion)
	}

	return m.verifyTerragruntInstallation()
}

// WithTerragruntAlpine sets up the container with Terragrunt and optionally Terraform on Alpine Linux.
// It updates the package list, installs required dependencies, and then installs the specified versions of Terragrunt and Terraform.
//
// Parameters:
// version - The version of Terragrunt to install. If empty, "latest" will be installed.
// tfVersion - The version of Terraform to install. If empty, "latest" will be installed.
// skipTerraform - If true, Terraform installation will be skipped.
func (m *ModuleTemplate) WithTerragruntAlpine(
	// version is the version of Terragrunt to install. If empty, it will be installed as "latest".
	// +optional
	version string,
	// tfVersion is the version of Terraform to install. If empty, it will be installed as "latest".
	// +optional
	tfVersion string,
	// skipTerraform if true, Terraform installation will be skipped.
	// +optional
	skipTerraform bool,
) *ModuleTemplate {
	version = resolveTerragruntVersion(version)
	if version == terragruntLatestVersion {
		version = getLatestTerragruntVersion()
	}

	tfVersion = resolveTerraformVersionInTerragrunt(tfVersion)
	if tfVersion == terragruntLatestVersion {
		tfVersion = getLatestTerraformVersionInTerragrunt()
	}

	m.Ctr = m.downloadRequiredUtilitiesTerragruntAlpine()
	m.Ctr = m.downloadAndInstallTerragrunt(version)

	if !skipTerraform {
		m.Ctr = m.downloadAndInstallTerraformInTerragrunt(tfVersion)
	}

	return m.verifyTerragruntInstallation()
}

func (m *ModuleTemplate) downloadAndInstallTerragrunt(version string) *dagger.Container {
	// Remove "v" prefix if present
	version = strings.TrimPrefix(version, "v")
	terragruntURL := fmt.Sprintf("%s/v%s/terragrunt_linux_amd64", terragruntReleaseURL, version)
	return m.Ctr.
		WithExec([]string{"wget", "-q", "-O", "/usr/local/bin/terragrunt", terragruntURL}).
		WithExec([]string{"chmod", "+x", "/usr/local/bin/terragrunt"})
}

func (m *ModuleTemplate) downloadAndInstallTerraformInTerragrunt(version string) *dagger.Container {
	terraformURL := fmt.Sprintf("%s/%s/terraform_%s_linux_amd64.zip", terraformReleaseURLForTerragrunt, version, version)
	return m.Ctr.
		WithExec([]string{"wget", "-q", terraformURL}).
		WithExec([]string{"unzip", fmt.Sprintf("terraform_%s_linux_amd64.zip", version)}).
		WithExec([]string{"mv", "terraform", "/usr/local/bin/"}).
		WithExec([]string{"chmod", "+x", "/usr/local/bin/terraform"}).
		WithExec([]string{"rm", fmt.Sprintf("terraform_%s_linux_amd64.zip", version)})
}

func (m *ModuleTemplate) verifyTerragruntInstallation() *ModuleTemplate {
	m.Ctr = m.Ctr.
		WithExec([]string{"terragrunt", "--version"}).
		WithExec([]string{"terraform", "version"})
	return m
}

func (m *ModuleTemplate) downloadRequiredUtilitiesTerragruntAlpine() *dagger.Container {
	return m.Ctr.
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "add", "--no-cache", "curl", "wget", "unzip"})
}

func (m *ModuleTemplate) downloadRequiredUtilitiesTerragruntUbuntu() *dagger.Container {
	return m.Ctr.
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "curl", "wget", "unzip"})
}
