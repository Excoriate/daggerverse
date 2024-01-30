package main

var (
	entryPointShell      = []string{"sh", "-c"}
	entryPointTerragrunt = []string{"terragrunt"}

	//excludedFiles = []string{".git", ".terraform.lock.hcl"}
	//excludedDirs  = []string{".terragrunt-cache/**", ".terraform/**"}
)

const (
	defaultContainerVersion = "latest"
	defaultContainerImage   = "alpine/terragrunt"
	workDirDefault          = "/mounted"
)
