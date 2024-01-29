package main

type buildDaggerCMDsOptions struct {
	entryPoint []string // E.g.: []string{"sh", "-c"}
	cmds       []string // E.g.: []string{"terragrunt plan", "terragrunt init"}
}

type DaggerCMD [][]string

func buildDaggerCMDs(cmds []buildDaggerCMDsOptions) DaggerCMD {
	var results DaggerCMD
	for _, c := range cmds {
		for _, command := range c.cmds {
			execCmd := append(c.entryPoint, command)
			results = append(results, execCmd)
		}
	}
	return results
}
