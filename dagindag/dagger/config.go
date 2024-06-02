package main

import (
	"fmt"
	"strings"
)

const (
	cmdUpdateAndInstallCURL = "apt update && apt install -y curl"
	daggerInstallScript     = "https://dl.dagger.io/dagger/install.sh"
	daggerDefaultVersion    = "v0.11.5"
	daggerCallCMD           = "dagger call -m " // intentionally left an empty space at the end
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
