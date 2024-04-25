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
	// Src is the directory that contains all the source code, including the module directory.
	Src *Directory
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
	// Src is the directory that contains all the source code,
	// including the module directory.
	src *Directory,
	// ctr is the container to use as a base container.
	// It's an optional parameter. If it's not set, it's going to create a new container.
	// +optional
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

	if src == nil {
		slog.Info("Src is not set, using the current module source directory")
		src = dag.CurrentModule().Source().Directory(".")
	}

	tt.Src = src

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
func (t *Terratest) Base(goVersion, tfVersion string) *Terratest {
	mod := dag.CacheVolume("gomodcache")
	build := dag.CacheVolume("gobuildcache")
	dotTerraform := dag.CacheVolume(".terraform")
	image := fmt.Sprintf("golang:%s", goVersion)

	c := dag.Container().
		From(image).
		WithMountedCache("/go/pkg/mod", mod).
		WithMountedCache("/root/.cache/go-build", build).
		WithMountedCache("/root/.terraform", dotTerraform).
		WithExec(t.getTFInstallCMD(tfVersion)).
		WithMountedDirectory(mntPrefix, t.Src)

	t.Ctr = c
	return t
}

func (t *Terratest) setWorkDir(testDirPath string) (*Container, error) {
	if testDirPath == "" {
		return nil, fmt.Errorf("the 'test' path cannot be empty")
	}

	if filepath.IsAbs(testDirPath) {
		return nil, fmt.Errorf("the 'test' path must be relative")
	}

	ttWorkdir := filepath.Join(mntPrefix, testDirPath)
	return t.Ctr.WithWorkdir(ttWorkdir), nil
}

// WithModule specifies the module to use in the Terraform module by the 'Src' directory.
func (t *Terratest) WithModule(src *Directory) *Terratest {
	t.Src = src
	return t
}

// WithContainer specifies the container to use in the Terraform module.
func (t *Terratest) WithContainer(ctr *Container) *Terratest {
	t.Ctr = ctr
	return t
}

func (t *Terratest) Run(
	// testDir is the directory that contains all the test code.
	testDir string,
	// args is the arguments to pass to the 'go test' command.
	// +optional
	args string,
) (*Container, error) {
	ctr, err := t.setWorkDir(testDir)
	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("The test directory resolved is: %s", ctr.Workdir))

	t.Ctr = ctr
	parsedArgs := parseArgsFromStrToSlice(args)

	cmdToRun := []string{"go", "test"}
	cmdToRun = append(cmdToRun, parsedArgs...)

	t.Ctr = t.Ctr.WithExec(cmdToRun).WithFocus()
	return t.Ctr, nil
}
