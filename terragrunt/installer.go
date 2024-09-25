package main

import (
	"fmt"
	"strings"

	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
)

// InstallerCfg defines the interface for configuring and executing the installation process.
type InstallerCfg interface {
	DownloadFrom(baseURL string)
	MoveTo(path string)
	Unzip(path string)
	SetPermissions(path string)
	ConfigureEnvVars(envVars map[string]string)
	ReadinessCommand(command string)
}

// Installer represents the configuration and state for installing a tool.
type Installer struct {
	name                  string
	baseURL               string
	installPath           string
	installPathWithBinary string
	cmds                  []string
	envVars               map[string]string
	readinessCmd          string
	unzipIsEnabled        bool
	unzipPath             string
	setExecPerms          bool
	execPerms             string
	cleanup               bool
}

// InstallerOption defines a function type for configuring an Installer.
type InstallerOption func(*Installer) error

// NewInstaller creates and returns a new Installer instance with the given name and base URL.
// It serves as the builder constructor for the Installer.
//
// Parameters:
//   - name: The name of the tool to be installed.
//   - baseURL: The URL from which the tool will be downloaded.
//
// Returns:
//   - *Installer: A pointer to the newly created Installer instance.
//   - error: An error if the input parameters are invalid.
func NewInstaller(name, baseURL string) (*Installer, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("the 'baseURL' can't be empty")
	}

	if name == "" {
		return nil, fmt.Errorf("the 'name' can't be empty")
	}

	return &Installer{
		name:                  name,
		baseURL:               baseURL,
		installPath:           "/usr/local/bin",
		installPathWithBinary: fmt.Sprintf("/usr/local/bin/%s", name),
		envVars:               make(map[string]string),
		setExecPerms:          true,
		execPerms:             "755",
	}, nil
}

// Install executes the installation process on the provided Dagger container.
// It applies all the configured options and sets up the commands in the container.
//
// Parameters:
//   - ctr: A pointer to the Dagger container where the installation will occur.
//   - opts: A variadic list of InstallerOption functions to configure the installation.
//
// Returns:
//   - *dagger.Container: A pointer to the modified Dagger container after installation.
//   - error: An error if any step of the installation process fails.
func (i *Installer) Install(ctr *dagger.Container, opts ...InstallerOption) (*dagger.Container, error) {
	for _, opt := range opts {
		if err := opt(i); err != nil {
			return nil, err
		}
	}

	// Download the file
	ctr = ctr.WithExec([]string{"wget", "-q", "-O", fmt.Sprintf("/tmp/%s", i.name), i.baseURL})

	// Unzip if enabled
	if i.unzipIsEnabled {
		ctr = ctr.WithExec([]string{"unzip", "-q", fmt.Sprintf("/tmp/%s", i.name), "-d", "/tmp"})
	}

	// Move file to install path
	moveCmd := []string{"mv"}
	if i.unzipIsEnabled {
		moveCmd = append(moveCmd, "/tmp/*")
	} else {
		moveCmd = append(moveCmd, fmt.Sprintf("/tmp/%s", i.name))
	}
	moveCmd = append(moveCmd, i.installPath)
	ctr = ctr.WithExec(moveCmd)

	// Set permissions (using sh -c to handle potential permission issues)
	// ctr = ctr.WithExec([]string{"sh", "-c", fmt.Sprintf("chmod %s %s", i.execPerms, i.installPathWithBinary)})
	// chmodCmd := []string{"chmod", i.execPerms, i.installPathWithBinary}
	// ctr = ctr.WithExec(chmodCmd)
	// Set environment variables
	for key, value := range i.envVars {
		ctr = ctr.WithEnvVariable(key, value)
	}

	// Add install path to PATH
	ctr = ctr.WithEnvVariable("PATH", fmt.Sprintf("%s:$PATH", i.installPath))

	// Execute additional commands
	for _, cmd := range i.cmds {
		ctr = ctr.WithExec(strings.Fields(cmd))
	}

	// Execute readiness command if set
	if i.readinessCmd != "" {
		ctr = ctr.WithExec(strings.Fields(i.readinessCmd))
	}

	// Verify installation
	ctr = ctr.WithExec([]string{i.name, "--version"})

	// Cleanup if enabled
	if i.cleanup {
		ctr = ctr.WithExec([]string{"rm", "-rf", fmt.Sprintf("/tmp/%s", i.name)})
	}

	return ctr, nil
}

// WithInstallerInstallPath sets the installation path.
func (i *Installer) WithInstallerInstallPath(path string) InstallerOption {
	return func(i *Installer) error {
		if path == "" {
			return fmt.Errorf("the 'path' can't be empty")
		}

		i.installPath = path
		return nil
	}
}

// WithInstallerUnzip enables unzipping.
func (i *Installer) WithInstallerUnzip() InstallerOption {
	return func(i *Installer) error {
		i.unzipIsEnabled = true
		return nil
	}
}

// WithInstallerSetPermissions sets permissions for the installed tool.
func (i *Installer) WithInstallerSetPermissions(permissions string) InstallerOption {
	return func(i *Installer) error {
		if permissions == "" {
			return fmt.Errorf("'permissions' can't be empty")
		}

		i.execPerms = permissions
		return nil
	}
}

// WithInstallerReadinessCommand sets a command to check if the tool is ready after installation.
func (i *Installer) WithInstallerReadinessCommand(command string) InstallerOption {
	return func(i *Installer) error {
		i.readinessCmd = command
		return nil
	}
}

// WithInstallerConfigureEnvVars adds environment variables to be set in the container.
func (i *Installer) WithInstallerConfigureEnvVars(envVars map[string]string) InstallerOption {
	return func(i *Installer) error {
		for key, value := range envVars {
			i.envVars[key] = value
		}
		return nil
	}
}

// WithInstallerCleanup enables cleanup after installation.
func (i *Installer) WithInstallerCleanup() InstallerOption {
	return func(i *Installer) error {
		i.cleanup = true
		return nil
	}
}
