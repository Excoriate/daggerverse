package main

import "fmt"

func FormImageAddress(image, version string) string {
	if version == "" {
		version = "latest"
	}

	return fmt.Sprintf("%s:%s", image, version)
}

func addCMDsToContainer(container *Container, cmds []string) *Container {
	container = container.WithEntrypoint(nil)

	cmdBuilt := append([]string(nil), entryPointShell...) // Create a new slice
	cmdBuilt = append(cmdBuilt, cmds...)                  // Append the command
	container = container.WithExec(cmdBuilt)

	return container
}

func addTGCommandsToContainer(container *Container, cmds []string, args ...string) *Container {
	container = container.WithEntrypoint(nil)

	cmdBuilt := append([]string(nil), entruPointTerragrunt...)
	cmdBuilt = append(cmdBuilt, cmds...)
	cmdBuilt = append(cmdBuilt, args...)
	container = container.WithExec(cmdBuilt)

	return container
}

func isEitherGlobalOrArgSetDir(arg Optional[*Directory], global *Directory) (*Directory, error) {
	if !arg.isSet && global == nil {
		return nil, fmt.Errorf("source directory cannot be empty, and it was not set in the constructor")
	}

	if !arg.isSet {
		return global, nil
	}

	return arg.value, nil
}

func isEitherGlobalOrArgSetString(arg Optional[string], global string) (string, error) {
	if !arg.isSet && global == "" {
		return "", fmt.Errorf("source directory cannot be empty, and it was not set in the constructor")
	}

	if !arg.isSet {
		return global, nil
	}

	return arg.value, nil
}
