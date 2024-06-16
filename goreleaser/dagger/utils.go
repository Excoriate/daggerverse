package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Excoriate/daggerx/pkg/envvars"
)

func replaceEntryPointForShell(ctr *Container) *Container {
	return ctr.
		WithoutEntrypoint().
		WithEntrypoint(nil)
}

// addEnvVarsToContainerFromStr adds environment variables to a container.
func (m *Goreleaser) addEnvVarsToContainerFromStr(envVarsFromHost string) *Container {
	envVars, _ := envvars.ToDaggerEnvVarsFromStr(envVarsFromHost)

	for _, envVar := range envVars {
		m.Ctr = m.Ctr.WithEnvVariable(envVar.Name, envVar.Value, ContainerWithEnvVariableOpts{
			Expand: envVar.Expand,
		})
	}

	return m.Ctr
}

func (m *Goreleaser) addEnvVarsToContainerFromSlice(envVars []string) *Container {
	envVarsParsed, _ := envvars.ToDaggerEnvVarsFromSlice(envVars)

	for _, envVar := range envVarsParsed {
		m.Ctr = m.Ctr.WithEnvVariable(envVar.Name, envVar.Value, ContainerWithEnvVariableOpts{
			Expand: envVar.Expand,
		})
	}

	return m.Ctr
}

func (m *Goreleaser) resolveCfgArg(cfg string) string {
	var cfgFileResolved string
	if cfg != "" && cfg != goReleaserDefaultCfgFile {
		cfgFileResolved = cfg
	} else {
		if m.CfgFile != "" {
			cfgFileResolved = m.CfgFile
		} else {
			_, err := os.Stat(goReleaserDefaultCfgFile)
			if os.IsNotExist(err) {
				// check for .yml extension
				goReleaserDefaultCfgFileWithYMLExt := strings.ReplaceAll(goReleaserDefaultCfgFile, ".yaml", ".yml")
				// do the check again
				_, err := os.Stat(goReleaserDefaultCfgFileWithYMLExt)
				if err == nil {
					cfgFileResolved = goReleaserDefaultCfgFileWithYMLExt
				} else {
					cfgFileResolved = ""
				}
			} else {
				cfgFileResolved = goReleaserDefaultCfgFile
			}
		}
	}

	return fmt.Sprintf("--config=%s", cfgFileResolved)
}
