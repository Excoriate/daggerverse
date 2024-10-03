package main

import (
	"github.com/Excoriate/daggerverse/terragrunt/internal/dagger"
	"github.com/Excoriate/daggerx/pkg/cmdx"
)

// terragruntCmd is the command to execute terragrunt.
var (
	terragruntCmd = []string{"terragrunt"}
)

// TgExec executes a given terragrunt command within a dagger container.
// It returns a pointer to the resulting dagger.Container or an error if the command is invalid or fails to execute.
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
	// gitRepoURL is the git repository to clone and use as a source directory for this Terragrunt execution.
	// +optional
	gitRepoURL string,
	// gitBranch is the branch to use when cloning the git repository.
	// +optional
	gitBranch string,
	// gitToken is the token to use when cloning the git repository.
	// +optional
	gitToken string,
	// gitSSHSocket is the SSH socket to use when cloning the git repository.
	// +optional
	gitSSHSocket *dagger.Socket,
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

	return m.
		Ctr.
		Terminal().
		WithExec(append(terragruntCmd, cmdAsSlice...)), nil
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
