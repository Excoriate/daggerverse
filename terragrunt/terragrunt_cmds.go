package main

import (
	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/cmdx"
	"github.com/Excoriate/daggerx/pkg/envvars"
)

type IACCommand interface {
	TgExecRemote(command string, gitRepoURL string, module string, gitBranch string, gitToken string, gitSSHSocket *dagger.Socket, envVars []string) (*dagger.Container, error)
	TgExec(command string, source *dagger.Directory, module string, gitSSHSocket *dagger.Socket, envVars []string, tfLogMode string, tgLogMode string, tfToken string) (*dagger.Container, error)
}

type IACToolTerragrunt struct {
	// Terragrunt is the terragrunt instance.
	// +private
	Terragrunt *Terragrunt
	// Entrypoint is the entrypoint to use when executing the terragrunt command.
	// +private
	Entrypoint string
}

// TgExec executes a given terragrunt command within a dagger container.
// It returns a pointer to the resulting dagger.Container or an error if the command is invalid or fails to execute.
//
//nolint:lll // It's okay, since the ignore pattern is included.
func (m *Terragrunt) TgExec(
	// command is the terragrunt command to execute. It's the actual command that comes after 'terragrunt'
	command string,
	// source is the source directory that includes the source code.
	// +defaultPath="/"
	// +ignore=[".terragrunt-cache", ".terraform", ".github", ".gitignore", ".git", "vendor", "node_modules", "build", "dist", "target", "tmp", "log"]
	source *dagger.Directory,
	// module is the module to execute or the terragrunt configuration where the terragrunt.hcl file is located.
	// +optional
	module string,
	// gitSSHSocket is the SSH socket to use when cloning the git repository.
	// +optional
	gitSSHSocket *dagger.Socket,
	// envVars are the environment variables to set in the container in the format of "key=value, key=value".
	// +optional
	envVars []string,
	// tfLogMode is the terraform log mode to use when executing the terragrunt command.
	// +optional
	tfLogMode string,
	// tgLogMode is the terragrunt log mode to use when executing the terragrunt command.
	// +optional
	tgLogMode string,
	// tfToken is the terraform token to use when executing the terragrunt command. It will form
	// an environment variable called TF_TOKEN_<token>
	// +optional
	tfToken string,
) (*dagger.Container, error) {
	if command == "" {
		return nil, WrapError(nil, "command is required, can't execute empty command")
	}

	if source == nil {
		return nil, WrapError(nil, "source is required, can't execute command without source")
	}

	if !isValidTerragruntCommand(command) {
		return nil, WrapErrorf(nil, "invalid terragrunt command: %s", command)
	}

	cmdAsSlice, cmdAsSliceErr := cmdx.GenerateDaggerCMDFromStr(command)
	if cmdAsSliceErr != nil {
		return nil, WrapErrorf(cmdAsSliceErr, "failed to generate dagger command from string: %s", command)
	}

	m.WithSource(source, module, containerUser)

	if envVars != nil {
		envVarsAsDagger, envVarsErr := envvars.ToDaggerEnvVarsFromSlice(envVars)
		if envVarsErr != nil {
			return nil, WrapErrorf(envVarsErr, "failed to convert environment variables to dagger environment variables: %s", envVars)
		}

		for _, envVar := range envVarsAsDagger {
			m.Ctr = m.Ctr.WithEnvVariable(envVar.Name, envVar.Value)
		}
	}

	return m.Ctr.
		WithExec(append([]string{"terragrunt"}, cmdAsSlice...)), nil
}

// isValidTerragruntCommand validates the given command against known terragrunt commands.
// It returns true if the command is valid, otherwise false.
func isValidTerragruntCommand(command string) bool {
	validCommands := map[string]bool{
		"apply": true, "destroy": true, "init": true, "plan": true, "output": true,
		"validate": true, "refresh": true, "show": true, "graph": true, "import": true,
		"state": true, "taint": true, "untaint": true, "workspace": true,
	}

	return validCommands[command]
}
