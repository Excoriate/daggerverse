// Package main provides utility functions for working with Docker containers.
//
// This package includes constants and a function for getting the Docker-in-Docker image.
//
// Copyright: Excoriate
// License: Apache-2.0
package main

// DefaultBinaryVersion represents the default version for a binary.
type DefaultBinaryVersion string

// ContainerImage represents a container image.
type ContainerImage string

const (
	// OpenTofuDefaultVersion specifies the default version for OpenTofu.
	OpenTofuDefaultVersion DefaultBinaryVersion = "1.8.2"
	// TerragruntDefaultVersion specifies the default version for Terragrunt.
	TerragruntDefaultVersion DefaultBinaryVersion = "0.67.4"
	// TerraformDefaultVersion specifies the default version for Terraform.
	TerraformDefaultVersion DefaultBinaryVersion = "1.9.5"
)
