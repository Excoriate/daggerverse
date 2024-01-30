package main

type DaggerCMD [][]string

func buildShellCMDs(cmds []string) DaggerCMD {
	var cmdsBuilt DaggerCMD

	for _, cmd := range cmds {
		cmdsBuilt = append(cmdsBuilt, []string{"sh", "-c", cmd})
	}

	return cmdsBuilt
}

func concatTerragruntInCommand(cmds []string) []string {
	var terragruntCmds []string

	for _, cmd := range cmds {
		terragruntCmds = append(terragruntCmds, "terragrunt "+cmd)
	}

	return terragruntCmds
}

func addCMDToDaggerCMD(cmd []string) [][]string {
	return [][]string{cmd}
}
