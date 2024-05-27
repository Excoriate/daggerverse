package main

import "github.com/Excoriate/daggerx/pkg/envvars"

func (m *Terragrunt) addEnvVarsToContainerFromSlice(envVars []string) *Container {
	envVarsParsed, _ := envvars.ToDaggerEnvVarsFromSlice(envVars)

	for _, envVar := range envVarsParsed {
		m.Ctr = m.Ctr.WithEnvVariable(envVar.Name, envVar.Value, ContainerWithEnvVariableOpts{
			Expand: envVar.Expand,
		})
	}

	return m.Ctr
}
