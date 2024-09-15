package main

const (
	// BaseImageDefaultRepositoryURL is the default repository URL for the base image.
	// This is the repository for the Docker Terragrunt image.
	BaseImageDefaultRepositoryURL = "ghcr.io/devops-infra/docker-terragrunt/docker-terragrunt"
	// BaseImageDefaultTag is the default tag for the base image.
	// For this https://github.com/devops-infra/docker-terragrunt repository, the tags
	// are the same as the ones used in the Dockerfile, and formed by the Terraform version
	// and the Terragrunt version.
	BaseImageDefaultTag = "tf-1.9.5-ot-1.8.2-tg-0.67.4"
)

// BaseImage is the base image for the Terragrunt container.
type BaseImage struct {
	// ImageRepositoryURL is the URL of the image repository.
	// +private
	ImageRepositoryURL string
	// ImageTag is the tag of the image.
	// +private
	ImageTag string
}

type BaseImageCfg interface {
	GetImageRepositoryURL() (string, error)
	GetImageTag() string
	GetOpeTofuVersion() string
	GetTerraformVersion() string
	GetTerragruntVersion() string
}
