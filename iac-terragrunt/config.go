package main

import "github.com/excoriate/daggerverse/daggercommon/pkg/container"

const (
	defaultContainerVersion = "latest"
	defaultContainerImage   = "alpine/terragrunt"
)

func GetDefaultImageAddr() string {
	return container.FormImageAddress(defaultContainerImage, defaultContainerVersion)
}
