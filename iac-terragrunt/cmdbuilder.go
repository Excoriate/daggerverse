package main

type buildDaggerCMDsOptions struct {
	entryPoint []string // E.g.: []string{"sh", "-c"}
	cmds       []string // E.g.: []string{"terragrunt plan", "terragrunt init"}
}

type DaggerCMD [][]string

func BuilderDaggerCMDs(cmds []string, entryPoint []string) DaggerCMD {
	var results DaggerCMD
	for _, cmd := range cmds {
		execCmd := append([]string(nil), entryPoint...)
		execCmd = append(execCmd, cmd)
		results = append(results, execCmd)
	}
	return results
}
