// Package main provides functionality for managing infrastructure-as-code toolkit versions.
package main

// Default versions for OpenTofu, Terraform, and Terragrunt.
const (
	defaultOpenTofuVersion   = "1.8.0"
	defaultTerraformVersion  = "1.9.1"
	defaultTerragruntVersion = "0.67.4"
)

// ToolkitCfg defines an interface for retrieving version information
// for OpenTofu, Terraform, and Terragrunt.
type ToolkitCfg interface {
	GetOpenTofuVersion() string
	GetTerraformVersion() string
	GetTerragruntVersion() string
}

// Toolkit represents a set of infrastructure-as-code tools with their respective versions.
type Toolkit struct {
	OpenTofuVersion   string
	TerraformVersion  string
	TerragruntVersion string
}

// NewToolkit creates a new Toolkit instance with default versions,
// which can be customized using functional options.
func NewToolkit(options ...func(*Toolkit)) *Toolkit {
	toolkit := &Toolkit{
		OpenTofuVersion:   defaultOpenTofuVersion,
		TerraformVersion:  defaultTerraformVersion,
		TerragruntVersion: defaultTerragruntVersion,
	}

	for _, option := range options {
		option(toolkit)
	}
	return toolkit
}

// WithOpenTofuVersion returns a function that sets the OpenTofu version.
// If an empty string is provided, it sets the default version.
func WithOpenTofuVersion(version string) func(*Toolkit) {
	return func(t *Toolkit) {
		if version == "" {
			t.OpenTofuVersion = defaultOpenTofuVersion
		} else {
			t.OpenTofuVersion = version
		}
	}
}

// WithTerraformVersion returns a function that sets the Terraform version.
// If an empty string is provided, it sets the default version.
func WithTerraformVersion(version string) func(*Toolkit) {
	return func(t *Toolkit) {
		if version == "" {
			t.TerraformVersion = defaultTerraformVersion
		} else {
			t.TerraformVersion = version
		}
	}
}

// WithTerragruntVersion returns a function that sets the Terragrunt version.
// If an empty string is provided, it sets the default version.
func WithTerragruntVersion(version string) func(*Toolkit) {
	return func(t *Toolkit) {
		if version == "" {
			t.TerragruntVersion = defaultTerragruntVersion
		} else {
			t.TerragruntVersion = version
		}
	}
}

// GetOpenTofuVersion returns the OpenTofu version.
func (t *Toolkit) GetOpenTofuVersion() string {
	return t.OpenTofuVersion
}

// GetTerraformVersion returns the Terraform version.
func (t *Toolkit) GetTerraformVersion() string {
	return t.TerraformVersion
}

// GetTerragruntVersion returns the Terragrunt version.
func (t *Toolkit) GetTerragruntVersion() string {
	return t.TerragruntVersion
}
