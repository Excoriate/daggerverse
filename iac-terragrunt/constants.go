package main

var (
	entryPointShell      = []string{"sh", "-c"}
	entruPointTerragrunt = []string{"terragrunt"}

	excludedFiles = []string{".git", ".terraform.lock.hcl"}
	excludedDirs  = []string{".terragrunt-cache/**", ".terraform/**"}
)

const (
	defaultContainerVersion = "latest"
	defaultContainerImage   = "alpine/terragrunt"
	workDirDefault          = "/mounted"
)

func buildCommand(entryPoint, commands []string) []string {
	if entryPoint == nil {
		entryPoint = entryPointShell
	}

	var cmdBuilt []string

	for _, command := range commands {
		cmdBuilt = append([]string(nil), entryPoint...)
		cmdBuilt = append(cmdBuilt, command)
	}

	return cmdBuilt
}
