package main

import (
	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/fixtures"
)

const (
	terragruntCacheDir  = "/home/terragrunt/.terragrunt-providers-cache"
	terraformCacheDir   = "/home/.terraform.d/plugin-cache"
	terraformPluginsDir = "/home/.terraform.d/plugins"
)

var (
	//nolint:gochecknoglobals // This is a global variable that is used to set the default permissions.
	terragruntPermissionsOnDirsDefault = []string{
		"/home/terragrunt",
		"/home/.terraform.d",
		"/home",
		"/var/log",
		fixtures.MntPrefix,
	}
)

// WithTerragruntPermissionsOnDirsDefault sets the default permissions for the Terragrunt directories.
// It ensures that the specified user and group own the directories and sets the appropriate permissions.
// The default directories include:
// - /home/terragrunt
// - /home/.terraform.d
// - /home
// - /var/log
// - fixtures.MntPrefix
// Additionally, it sets the permissions to 0777 for the following directories:
// - /home
// - /var/log
// - fixtures.MntPrefix
//
// Returns:
// - *Terragrunt: Updated instance with the default permissions set.
func (m *Terragrunt) WithTerragruntPermissionsOnDirsDefault() *Terragrunt {
	return m.
		WithUserAsOwnerOfDirs(terragruntCtrUser, terragruntCtrGroup, terragruntPermissionsOnDirsDefault, true).
		WithUserWithPermissionsOnDirs(terragruntCtrUser, "0777", terragruntPermissionsOnDirsDefault, true)
}

// WithTerragruntPermissionsOnDirs sets the necessary permissions for the Terragrunt directories.
// It ensures that the specified user and group own the directories and sets the appropriate permissions.
// The default directories include:
// - /home/terragrunt
// - /home/.terraform.d
// - /home
// - /var/log
// - fixtures.MntPrefix
// Additionally, it sets the permissions to 0777 for the following directories:
// - /home
// - /var/log
// - fixtures.MntPrefix
// If dirsToOwn and dirsToHaveWritePermissions are provided, they will be appended to the respective default lists.
func (m *Terragrunt) WithTerragruntPermissionsOnDirs(
	// dirsToOwn are the directories to set the permissions to 0777.
	// +optional
	dirsToOwn []string,
	// dirsToHaveWritePermissions are the directories to have write permissions.
	// +optional
	dirsToHaveWritePermissions []string,
) *Terragrunt {
	dirsToOwn = append(terragruntPermissionsOnDirsDefault, dirsToOwn...)
	dirsToHaveWritePermissions = append(dirsToHaveWritePermissions, dirsToHaveWritePermissions...)

	return m.WithUserAsOwnerOfDirs(terragruntCtrUser, terragruntCtrGroup, dirsToOwn, true).
		WithUserWithPermissionsOnDirs(terragruntCtrUser, "0777", dirsToHaveWritePermissions, true)
}

// WithTerragruntCacheConfiguration configures the cache directory for Terragrunt.
// It sets up the cache directory at /home/terragrunt/.terragrunt-providers-cache
// and assigns the environment variable TERRAGRUNT_PROVIDER_CACHE_DIR to this path.
func (m *Terragrunt) WithTerragruntCacheConfiguration() *Terragrunt {
	return m.
		WithCachedDirectory(terragruntCacheDir,
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
		WithCachedDirectory(terraformCacheDir, false, "TF_PLUGIN_CACHE_DIR",
			dagger.Shared,
			nil,
			"terragrunt",
		).
		WithCachedDirectory(terraformPluginsDir, false, "",
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
