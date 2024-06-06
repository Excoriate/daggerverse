package main

import (
	"fmt"
	"log/slog"
	"path/filepath"
)

const (
	defaultImageVersion = "1.22.0-alpine3.19"
	defaultTfVersion    = "1.6.0"
	mntPrefix           = "/mnt"
)

type Terratest struct {
	// The Version of the Golang image that'll host the 'terratest' test
	Version string
	// TfVersion is the Version of the Terraform to use, e.g., "0.12.24".
	// by default, it uses the latest Version.
	TfVersion string
	// Image of the container to use.
	Image string
	// Ctr is the container to use as a base container.
	Ctr *Container
}

func New(
	// the Version of the Terraform to use, e.g., "0.12.24".
	// by default, it uses the latest Version.
	// +default="1.22.0-alpine3.19"
	// +optional
	version string,
	// the Version of the Terraform to use, e.g., "0.12.24".
	// by default, it uses the latest Version.
	// +default="1.6.0"
	// +optional
	tfVersion string,
	// Image of the container to use.
	// by default, it uses the official HashiCorp Terraform Image hashicorp/terraform.
	// +default="gcr.io/distroless/static-debian11"
	// +optional
	image string,
	// ctr is the container to use as a base container.
	// It's an optional parameter. If it's not set, it's going to create a new container.
	ctr *Container,
	// envVars is a string of environment variables in the form of "key1=value1,key2=value2"
	// +optional
	envVars string,
) *Terratest {
	tt := &Terratest{
		Version: version,
		Image:   image,
	}

	if tfVersion == "" {
		tfVersion = defaultTfVersion
	}

	if version == "" {
		version = defaultImageVersion
	}

	if ctr == nil {
		ctr = tt.Base(version, tfVersion).Ctr
	}

	tt.Ctr = ctr

	// A bit of a dirty hack to get environment variables into the container
	// from the string passed in the envVars parameter.
	if envVars != "" {
		envVarsParsed := tt.parseEnvVarsInStringMapAsMap(envVars)
		tt.Ctr = tt.setEnvVarsInContainer(envVarsParsed)
		slog.Info(fmt.Sprintf("Environment variables set: %v", envVarsParsed))
	}

	slog.Info("Terraform container created")

	return tt
}

// Base sets up the Container with a golang image and cache volumes
// version string
func (m *Terratest) Base(goVersion, tfVersion string) *Terratest {
	mod := dag.CacheVolume("gomodcache")
	build := dag.CacheVolume("gobuildcache")
	dotTerraform := dag.CacheVolume(".terraform")
	image := fmt.Sprintf("golang:%s", goVersion)

	c := dag.Container().
		From(image).
		WithMountedCache("/go/pkg/mod", mod).
		WithMountedCache("/root/.cache/go-build", build).
		WithMountedCache("/root/.terraform", dotTerraform).
		WithExec(m.getTFInstallCMD(tfVersion))

	m.Ctr = c
	return m
}

func (m *Terratest) setWorkDir(testDirPath string) (*Container, error) {
	if testDirPath == "" {
		return nil, fmt.Errorf("the 'test' path cannot be empty")
	}

	if filepath.IsAbs(testDirPath) {
		return nil, fmt.Errorf("the 'test' path must be relative")
	}

	ttWorkdir := filepath.Join(mntPrefix, testDirPath)
	return m.Ctr.WithWorkdir(ttWorkdir), nil
}

// WithContainer specifies the container to use in the Terraform module.
func (m *Terratest) WithContainer(ctr *Container) *Terratest {
	m.Ctr = ctr
	return m
}

func (m *Terratest) Run(
	// testDir is the directory that contains all the test code.
	testDir *Directory,
	// args is the arguments to pass to the 'go test' command.
	// +optional
	args string,
) (*Container, error) {
	m.Ctr = m.WithSource(testDir, mntPrefix).Ctr

	parsedArgs := parseArgsFromStrToSlice(args)

	cmdToRun := []string{"go", "test"}
	cmdToRun = append(cmdToRun, parsedArgs...)

	m.Ctr = m.Ctr.WithExec(cmdToRun).WithFocus()
	return m.Ctr, nil
}
