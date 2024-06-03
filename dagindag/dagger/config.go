package main

import (
	"fmt"
	"strings"
)

const (
	cmdUpdateAndInstallCURL = "apk update && apk add curl"
	daggerDefaultVersion    = "v0.11.6"
	daggerCallCMD           = "call -m"
	dockerVersionDefault    = "24.0"
)

var (
	daggerCLIEntryPoint = []string{"bin/dagger"}
)

// getDaggerInstallCMDByVersion returns the command to install the Dagger engine with the given version.
//
// The command is a shell script that sets the environment variable DAGGER_VERSION to the given version
// and then downloads and runs the Dagger install script.
//
// Example:
//
//	getDaggerInstallCMDByVersion("v0.11.5") => `cd / && DAGGER_VERSION="v0.11.5" curl -L https://dl.dagger.io/dagger/install.sh | DAGGER_VERSION="v0.11.5" sh`
func getDaggerInstallCMDByVersion(version string) string {
	return fmt.Sprintf(`cd / && DAGGER_VERSION="%s" curl -L https://dl.dagger.io/dagger/install.sh | DAGGER_VERSION="%s" sh`, version, version)
}

// getDaggerCallCMD returns the command to call the Dagger engine with the given module.
//
// The command is a shell script that calls the Dagger engine with the given module.
//
// Example:
//
//	getDaggerCallCMD("module") => `dagger call -m module`
func getDaggerCallCMD(module string) string {
	return fmt.Sprintf("%s%s", daggerCallCMD, strings.TrimSpace(module))
}

// getDaggerModulePath returns the module path with the given module and version.
//
// If the version is empty, the module path is returned as "main".
//
// Example:
//
//	getDaggerModulePath("github.com/owner/repo/module", "v0.11.5") => "github.com/owner/repo/module@v0.11.5"
func getDaggerModulePath(module, version string) string {
	if version == "" {
		version = "main"
	}

	return fmt.Sprintf("%s@%s", module, version)
}

// getDockerInDockerImage returns the Docker-in-Docker image with the given version.
//
// Example:
//
//	getDockerInDockerImage("20.10.17") => "docker:20.10.17-dind"

func getDockerInDockerImage(version string) string {
	if version == "" {
		version = dockerVersionDefault
	}

	return fmt.Sprintf("docker:%s-dind", version)
}
