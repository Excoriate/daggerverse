package main

import (
	"fmt"

	"github.com/Excoriate/daggerverse/module-template/internal/dagger"
)

const (
	latestVersion       = "latest"
	terraformReleaseURL = "https://releases.hashicorp.com/terraform"
	hashicorpGPGURL     = "https://apt.releases.hashicorp.com/gpg"
	hashicorpRepoURL    = "https://apt.releases.hashicorp.com"
)

// resolveTerraformVersion resolves the provided version to "latest" if it's empty.
func resolveTerraformVersion(version string) string {
	if version == "" {
		return latestVersion
	}

	return version
}

// getLatestTerraformVersion fetches the latest Terraform version.
// This is a placeholder function. In a real scenario, you'd implement logic to fetch the latest version.
func getLatestTerraformVersion() string {
	// Placeholder: In reality, you'd fetch this from Terraform's releases page or API
	return "1.9.4" // Example latest version
}

// WithTerraformUbuntu sets up the container with Terraform on Ubuntu.
// It updates the package list, installs required dependencies, and then installs the specified version of Terraform.
//
// Parameters:
// version - The version of Terraform to install. If empty, "latest" will be installed.
func (m *ModuleTemplate) WithTerraformUbuntu(
	// version is the version of Terraform to install. If empty, it will be installed as "latest".
	// +optional
	version string,
) *ModuleTemplate {
	version = resolveTerraformVersion(version)
	if version == latestVersion {
		version = getLatestTerraformVersion()
	}

	m.Ctr = m.downloadRequiredUtilitiesUbuntu()

	m.Ctr = m.Ctr.
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "gnupg", "software-properties-common"}).
		WithExec(m.getGPGKeyCommand()).
		WithExec(m.getAddRepoCommand()).
		WithExec([]string{"apt-get", "update"})

	m.Ctr = m.installTerraformUbuntu(version)

	return m.verifyTerraformInstallation()
}

func (m *ModuleTemplate) getGPGKeyCommand() []string {
	command := fmt.Sprintf("wget -O- %s | gpg --dearmor | "+
		"tee /usr/share/keyrings/hashicorp-archive-keyring.gpg > /dev/null", hashicorpGPGURL)

	return []string{"bash", "-c", command}
}

func (m *ModuleTemplate) getAddRepoCommand() []string {
	distro := "$(lsb_release -cs)"
	repoLine := fmt.Sprintf(
		`echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] %s %s main"`,
		hashicorpRepoURL, distro,
	)

	command := repoLine + " | tee /etc/apt/sources.list.d/hashicorp.list"

	return []string{"bash", "-c", command}
}

func (m *ModuleTemplate) installTerraformUbuntu(version string) *dagger.Container {
	return m.Ctr.
		WithExec([]string{"apt-get", "install", "-y", "terraform=" + version})
}

// WithTerraformAlpine sets up the container with Terraform on Alpine Linux.
// It updates the package list, installs required dependencies, and then installs the
// specified version of Terraform.
//
// Parameters:
// version - The version of Terraform to install. If empty, "latest" will be installed.
func (m *ModuleTemplate) WithTerraformAlpine(
	// version is the version of Terraform to install. If empty, it will be installed as "latest".
	// +optional
	version string,
) *ModuleTemplate {
	version = resolveTerraformVersion(version)
	if version == latestVersion {
		version = getLatestTerraformVersion()
	}

	m.Ctr = m.downloadRequiredUtilitiesAlpine()

	m.Ctr = m.downloadAndInstallTerraform(version)

	return m.verifyTerraformInstallation()
}

func (m *ModuleTemplate) downloadAndInstallTerraform(version string) *dagger.Container {
	terraformURL := fmt.Sprintf("%s/%s/terraform_%s_linux_amd64.zip",
		terraformReleaseURL, version, version)

	zipFile := fmt.Sprintf("terraform_%s_linux_amd64.zip", version)

	return m.Ctr.
		WithExec([]string{"wget", terraformURL}).
		WithExec([]string{"unzip", zipFile}).
		WithExec([]string{"rm", zipFile}).
		WithExec([]string{"mv", "terraform", "/usr/local/bin"})
}

func (m *ModuleTemplate) verifyTerraformInstallation() *ModuleTemplate {
	m.Ctr = m.Ctr.WithExec([]string{"terraform", "version"})

	return m
}

func (m *ModuleTemplate) downloadRequiredUtilitiesAlpine() *dagger.Container {
	return m.Ctr.
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "add", "curl", "wget", "bash", "unzip", "yq"})
}

func (m *ModuleTemplate) downloadRequiredUtilitiesUbuntu() *dagger.Container {
	return m.Ctr.
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "curl", "wget", "bash", "unzip", "yq"})
}
