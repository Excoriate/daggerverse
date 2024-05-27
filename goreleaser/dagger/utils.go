package main

import (
	"fmt"

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
			cfgFileResolved = goReleaserDefaultCfgFile
		}
	}

	return fmt.Sprintf("--config=%s", cfgFileResolved)
}
