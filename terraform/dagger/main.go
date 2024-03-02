package main

import (
	"fmt"
	"golang.org/x/exp/slog"
	"path/filepath"
	"strings"
)

const mntPrefix = "/mnt"

const (
	tfCmdInit     = "init"
	tfCmdPlan     = "plan"
	tfCmdApply    = "apply"
	tfCmdDestroy  = "destroy"
	tfCmdValidate = "validate"
	tfFmt         = "fmt"
)

type Terraform struct {
	// The Version of the Terraform to use, e.g., "0.12.24".
	Version string
	// Image of the container to use.
	Image string
	// Src is the directory that contains all the source code, including the module directory.
	Src *Directory
	// BaseCtr is the container to use as a base container.
	BaseCtr *Container
}

func New(
	// the Version of the Terraform to use, e.g., "0.12.24".
	// by default, it uses the latest Version.
	// +default="latest"
	// +optional
	version string,

	// Image of the container to use.
	// by default, it uses the official HashiCorp Terraform Image hashicorp/terraform.
	// +default="hashicorp/terraform"
	// +optional
	image string,

	// Src is the directory that contains all the source code,
	// including the module directory.
	src *Directory,
	// ctr is the container to use as a base container.
	// It's an optional parameter. If it's not set, it's going to create a new container.
	// +optional
	ctr *Container,
) *Terraform {
	tf := &Terraform{
		Version: version,
		Image:   image,
		Src:     src,
	}

	slog.Info(fmt.Sprintf("Terraform Version: %s", version))
	slog.Info(fmt.Sprintf("Terraform Image: %s", image))

	if ctr == nil {
		ctr = tf.Base().BaseCtr
	}

	if src == nil {
		slog.Info("Src is not set, using the current module source directory")
		src = dag.CurrentModule().Source().Directory(".")
	}

	tf.BaseCtr = ctr

	slog.Info("Terraform container created")

	return tf
}

// Base sets up the Container with a Terraform Image and cache volumes
func (t *Terraform) Base() *Terraform {
	tfCache := dag.CacheVolume(".terraform")
	image := fmt.Sprintf("%s:%s", t.Image, t.Version)
	c := dag.Container().
		From(image).
		WithMountedCache("/.terraform", tfCache)

	t.BaseCtr = c
	return t
}

// WithModule specifies the module to use in the Terraform module by the 'Src' directory.
func (t *Terraform) WithModule(src *Directory) *Terraform {
	t.Src = src
	return t
}

// WithContainer specifies the container to use in the Terraform module.
func (t *Terraform) WithContainer(ctr *Container) *Terraform {
	t.BaseCtr = ctr
	return t
}

func parseArgsFromStrToSlice(argStr string) []string {
	if argStr == "" {
		return []string{}
	}

	var parsedArgs []string
	// Split the string on the comma as a preliminary step.
	args := strings.Split(argStr, ",")
	for _, arg := range args {
		// Trim leading and trailing whitespace from each argument.
		arg = strings.TrimSpace(arg)
		// Handle special case for '-var' arguments.
		if strings.HasPrefix(arg, "-var ") {
			// Keep '-var' and its value together as a single string.
			parsedArgs = append(parsedArgs, arg)
		} else {
			// For other arguments, split on spaces assuming they are separate.
			parts := strings.Fields(arg)
			parsedArgs = append(parsedArgs, parts...)
		}
	}
	return parsedArgs
}

func (t *Terraform) setCommands(command string, args ...string) *Container {
	// Initialize the command slice with the main command (e.g., "plan")
	cmdWithArgs := []string{command}

	for _, arg := range args {
		if strings.Contains(arg, " ") && !strings.HasPrefix(arg, "'") && !strings.HasPrefix(arg, "\"") {
			parts := strings.SplitN(arg, " ", 2)
			cmdWithArgs = append(cmdWithArgs, parts[0], parts[1])
		} else {
			cmdWithArgs = append(cmdWithArgs, arg)
		}
	}

	t.BaseCtr = t.BaseCtr.WithExec(cmdWithArgs)

	return t.BaseCtr
}

func (t *Terraform) setTFModuleSRC(tfModPath string) *Container {
	tfWorkDir := filepath.Join(mntPrefix, tfModPath)
	slog.Info(fmt.Sprintf("The Terraform module directory resolved is: %s", tfWorkDir))

	return t.BaseCtr.
		WithMountedDirectory(mntPrefix, t.Src).
		WithWorkdir(tfWorkDir)
}

// Init initializes the Terraform module.
func (t *Terraform) Init(
	// The tfmod is the Terraform module to use.
	tfmod string,
	// args are the n number of arguments to pass to the Terraform init command.
	// +optional
	args string,
) (*Container, error) {
	t.BaseCtr = t.setTFModuleSRC(tfmod)
	parsedArgs := parseArgsFromStrToSlice(args)
	t.BaseCtr = t.setCommands(tfCmdInit, parsedArgs...)
	//return t.BaseCtr.Stdout(ctx)
	return t.BaseCtr, nil
}

// Plan creates an execution plan for the Terraform module.
func (t *Terraform) Plan(
	// The tfmod is the Terraform module to use.
	tfmod string,
	// args are the n number of arguments to pass to the Terraform plan command.
	// +optional
	args string,
	// initArgs are the n number of arguments to pass to the Terraform init command.
	// +optional
	initArgs string,
) (*Container, error) {
	t.BaseCtr = t.setTFModuleSRC(tfmod)
	parsedInitArgs := parseArgsFromStrToSlice(initArgs)
	// Set the init command
	t.BaseCtr = t.setCommands(tfCmdInit, parsedInitArgs...)
	// Set the plan command
	parsedArgs := parseArgsFromStrToSlice(args)
	t.BaseCtr = t.setCommands(tfCmdPlan, parsedArgs...)

	return t.BaseCtr, nil
}

// Apply creates an execution plan for the Terraform module.
func (t *Terraform) Apply(
	// The tfmod is the Terraform module to use.
	tfmod string,
	// args are the n number of arguments to pass to the Terraform apply command.
	// +optional
	args string,
	// initArgs are the n number of arguments to pass to the Terraform init command.
	// +optional
	initArgs string,
) (*Container, error) {
	t.BaseCtr = t.setTFModuleSRC(tfmod)
	parsedInitArgs := parseArgsFromStrToSlice(initArgs)
	// Set the init command
	t.BaseCtr = t.setCommands(tfCmdInit, parsedInitArgs...)
	// Set the plan command
	parsedArgs := parseArgsFromStrToSlice(args)
	t.BaseCtr = t.setCommands(tfCmdApply, parsedArgs...)
	return t.BaseCtr, nil
}

// Destroy creates an execution plan for the Terraform module.
func (t *Terraform) Destroy(
	// The tfmod is the Terraform module to use.
	tfmod string,
	// args are the n number of arguments to pass to the Terraform destroy command.
	// +optional
	args string,
	// initArgs are the n number of arguments to pass to the Terraform init command.
	// +optional
	initArgs string,
) (*Container, error) {
	t.BaseCtr = t.setTFModuleSRC(tfmod)
	parsedInitArgs := parseArgsFromStrToSlice(initArgs)
	// Set the init command
	t.BaseCtr = t.setCommands(tfCmdInit, parsedInitArgs...)
	// Set the plan command
	parsedArgs := parseArgsFromStrToSlice(args)
	t.BaseCtr = t.setCommands(tfCmdDestroy, parsedArgs...)
	return t.BaseCtr, nil
}

// Validate creates an execution plan for the Terraform module.
func (t *Terraform) Validate(
	// The tfmod is the Terraform module to use.
	tfmod string,
	// args are the n number of arguments to pass to the Terraform validate command.
	// +optional
	args string,
	// initArgs are the n number of arguments to pass to the Terraform init command.
	// +optional
	initArgs string,
) (*Container, error) {
	t.BaseCtr = t.setTFModuleSRC(tfmod)
	parsedInitArgs := parseArgsFromStrToSlice(initArgs)
	// Set the init command
	t.BaseCtr = t.setCommands(tfCmdInit, parsedInitArgs...)
	// Set the plan command
	parsedArgs := parseArgsFromStrToSlice(args)
	t.BaseCtr = t.setCommands(tfCmdValidate, parsedArgs...)
	return t.BaseCtr, nil
}

// Format creates an execution plan for the Terraform module.
func (t *Terraform) Format(
	// The tfmod is the Terraform module to use.
	tfmod string,
	// args are the n number of arguments to pass to the Terraform fmt command.
	// +optional
	args string,
) (*Container, error) {
	t.BaseCtr = t.setTFModuleSRC(tfmod)
	parsedArgs := parseArgsFromStrToSlice(args)
	t.BaseCtr = t.setCommands(tfFmt, parsedArgs...)
	return t.BaseCtr, nil
}
