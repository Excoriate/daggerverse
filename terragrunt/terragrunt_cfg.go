package main

import (
	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/fixtures"
)

const (
	terragruntCacheDir  = "/home/terragrunt/.terragrunt-providers-cache"
	terraformCacheDir   = "/home/.terraform.d/plugin-cache"
	terraformPluginsDir = "/home/.terraform.d/plugins"
	terragruntHomeDir   = "/home/terragrunt"
	terraformDir        = "/home/.terraform.d"
	homeDir             = "/home"
	varLogDir           = "/var/log"
	mntPrefixDefault    = fixtures.MntPrefix
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
	defaultDirs := []string{
		terragruntHomeDir,
		terraformDir,
		homeDir,
		varLogDir,
		mntPrefixDefault,
	}

	return m.
		WithUserAsOwnerOfDirs(terragruntCtrUser, terragruntCtrGroup, defaultDirs, true).
		WithUserWithPermissionsOnDirs(terragruntCtrUser, "0777", defaultDirs, true)
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
	defaultDirs := []string{
		terragruntHomeDir,
		terraformDir,
		homeDir,
		varLogDir,
		mntPrefixDefault,
	}

	dirsToOwn = append(defaultDirs, dirsToOwn...)
	dirsToHaveWritePermissions = append(defaultDirs, dirsToHaveWritePermissions...)

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

// WithSSHAuthForTerraformModules configures SSH authentication for Terraform modules with Git SSH sources.
//
// This function mounts an SSH authentication socket into the container, enabling Terraform to authenticate
// when fetching modules from Git repositories using SSH URLs (e.g., git@github.com:org/repo.git).
//
// Parameters:
//   - sshAuthSocket: The SSH authentication socket to mount in the container.
//   - socketPath: The path where the SSH socket will be mounted in the container.
//   - owner: Optional. The owner of the mounted socket in the container.
//
// Returns:
//   - *Terragrunt: The updated Terragrunt instance with SSH authentication configured for Terraform modules.
func (m *Terragrunt) WithSSHAuthForTerraformModules(
	// sshAuthSocket is the SSH socket to use for authentication.
	sshAuthSocket *dagger.Socket,
	// socketPath is the path where the SSH socket will be mounted in the container.
	socketPath string,
	// owner is the owner of the mounted socket in the container. Optional parameter.
	// +optional
	owner string,
) *Terragrunt {
	socketOpts := dagger.ContainerWithUnixSocketOpts{}

	if owner != "" {
		socketOpts.Owner = owner
	}

	m.Ctr = m.Ctr.
		WithUnixSocket(socketPath, sshAuthSocket, socketOpts)

	return m
}

// WithTerragruntProviderCacheServerDisabled disables the Terragrunt provider cache server.
// It sets the environment variable TERRAGRUNT_PROVIDER_CACHE to "0".
// WithTerragruntProviderCacheServerDisabled disables the Terragrunt provider cache server.
//
// By default, it's enabled, but in some cases, you may want to disable it.
//
// Returns:
//   - *Terragrunt: The updated Terragrunt instance with the provider cache server disabled.
func (m *Terragrunt) WithTerragruntProviderCacheServerDisabled() *Terragrunt {
	m.Ctr = m.Ctr.
		WithoutEnvVariable("TERRAGRUNT_PROVIDER_CACHE").
		WithEnvVariable("TERRAGRUNT_PROVIDER_CACHE", "1")

	return m
}
