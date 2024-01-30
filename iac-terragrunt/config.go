package main

var (
	entryPointTerragrunt = []string{"terragrunt"}
)

const (
	defaultContainerVersion = "latest"
	defaultContainerImage   = "alpine/terragrunt"
	workDirDefault          = "/mounted"
)
