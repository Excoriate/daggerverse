// Package main provides utility functions for working with Docker containers.
//
// This package includes constants and a function for getting the Docker-in-Docker image.
//
// Copyright: Excoriate
// License: Apache-2.0
package main

type DefaultBinaryVersion string
type ContainerImage string

const (
	// OpenTofuDefaultVersion specifies the default version for OpenTofu.
	OpenTofuDefaultVersion DefaultBinaryVersion = "1.0.0"
	// TerragruntDefaultVersion specifies the default version for Terragrunt.
	TerragruntDefaultVersion DefaultBinaryVersion = "0.35.0"
	// TerraformDefaultVersion specifies the default version for Terraform.
	TerraformDefaultVersion DefaultBinaryVersion = "1.0.11"
	// TerragruntAlpineDefaultVersion specifies the default version for the Terragrunt Alpine container image.
	TerragruntAlpineDefaultVersion ContainerImage = "0.35.0-alpine"
)
