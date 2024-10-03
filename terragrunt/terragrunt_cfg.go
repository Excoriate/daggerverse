package main

import (
	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/fixtures"
)

// WithTerragruntPermissions sets the necessary permissions for the Terragrunt directories.
// It ensures that the specified user and group own the directories and sets the appropriate permissions.
// The directories include:
// - /home/terragrunt
// - /home/.terraform.d
// - /home
// - fixtures.MntPrefix
// Additionally, it sets the permissions to 0777 for the following directories:
// - /home
// - fixtures.MntPrefix.
func (m *Terragrunt) WithTerragruntPermissions() *Terragrunt {
	return m.WithUserAsOwnerOfDirs(containerUser, containerGroup, []string{
		"/home/terragrunt",
		"/home/.terraform.d",
		"/home",
		fixtures.MntPrefix,
	}, true).
		WithUserWithPermissionsOnDirs(containerUser, "0777", []string{
			"/home",
			fixtures.MntPrefix,
		}, true)
}

// WithTerragruntCacheConfiguration configures the cache directory for Terragrunt.
// It sets up the cache directory at /home/terragrunt/.terragrunt-providers-cache
// and assigns the environment variable TERRAGRUNT_PROVIDER_CACHE_DIR to this path.
func (m *Terragrunt) WithTerragruntCacheConfiguration() *Terragrunt {
	return m.
		WithCachedDirectory("/home/terragrunt/.terragrunt-providers-cache",
			false, "TERRAGRUNT_PROVIDER_CACHE_DIR",
			dagger.Shared,
			nil,
			"terragrunt",
		)
}

// WithTerraformCacheConfiguration configures the cache directories for Terraform.
// It sets up the following cache directories:
// - /home/.terraform.d/plugin-cache with the environment variable TF_PLUGIN_CACHE_DIR.
// - /home/.terraform.d/plugins without any specific environment variable.
func (m *Terragrunt) WithTerraformCacheConfiguration() *Terragrunt {
	return m.
		WithCachedDirectory("/home/.terraform.d/plugin-cache", false, "TF_PLUGIN_CACHE_DIR",
			dagger.Shared,
			nil,
			"terragrunt",
		).
		WithCachedDirectory("/home/.terraform.d/plugins", false, "",
			dagger.Shared,
			nil,
			"terragrunt",
		)
}

// WithIACToolsInstalled ensures that the specified versions of Terragrunt, Terraform, and OpenTofu are installed.
// If any of the provided version strings are empty, it defaults to the predefined versions for each tool.
// The function performs the following steps:
// 1. Checks if the provided Terragrunt version is empty. If it is, it assigns the defaultTerragruntVersion.
// 2. Checks if the provided Terraform version is empty. If it is, it assigns the defaultTerraformVersion.
// 3. Checks if the provided OpenTofu version is empty. If it is, it assigns the defaultOpenTofuVersion.
// 4. Calls WithTerragruntInstalled with the determined Terragrunt version to ensure Terragrunt is installed.
// 5. Calls WithTerraformInstalled with the determined Terraform version to ensure Terraform is installed.
// 6. Calls WithOpenTofuInstalled with the determined OpenTofu version to ensure OpenTofu is installed.
// The function returns the modified Terragrunt instance with the specified tools installed.
func (m *Terragrunt) WithIACToolsInstalled(tgVersion, tfVersion, openTofuVersion string) *Terragrunt {
	if tgVersion == "" {
		tgVersion = defaultTerragruntVersion
	}

	if tfVersion == "" {
		tfVersion = defaultTerraformVersion
	}

	if openTofuVersion == "" {
		openTofuVersion = defaultOpenTofuVersion
	}

	return m.
		WithTerragruntInstalled(tgVersion).
		WithTerraformInstalled(tfVersion).
		WithOpenTofuInstalled(openTofuVersion)
}
