package main

import "main/container"

const (
	defaultContainerVersion = "latest"
	defaultContainerImage   = "alpine/terragrunt"
)

func GetDefaultImageAddr() string {
	return container.FormImageAddress(defaultContainerImage, defaultContainerVersion)
}
