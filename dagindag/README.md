# Dagger In Dagger Module 🧑‍🚀

![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.0-0f0f19.svg?style=flat-square)

A simple [Dagger](https://dagger.io) module that wraps Dagger in Dagger, which means you can call other modules, functions, or Dagger commands from within a Dagger module. It's specially useful for:

- ✅ Test other Dagger modules.
- ✅ Call Dagger commands from within a Dagger module.
- ✅ Publish your Dagger modules through CI.

## Configuration 🛠️

Through the [Dagger CLI](https://docs.dagger.io/cli/465058/install), or by using it directly reusing it within your module, you can configure the following options:

- ⚙️ `daggerVersion`: Is the version of Dagger to use. It's a string that represents the version of Dagger to use. It's optional, and if not provided, the latest version of Dagger will be used (Currently, **0.11.6**).
- ⚙️ `dockerVersion`: Is the version of Docker to use. It's a string that represents the version of Docker to use. It's optional, and if not provided, the latest version of Docker will be used (currently, **0.24.0**).

> **NOTE**: This configuration is available through the [module's constructor](https://docs.dagger.io/manuals/developer/go/520657/constructor/).

---

## Features 🎨

### Commands and Functionalities 📜

| Command or functionality | Command    | Example                       |
| ------------------------ | ---------- | ----------------------------- |
| Terminal                 | `terminal` | `dagger call terminal --help` |
| Dag (Dagger) CLI         | `dag-cli`  | `dagger call dag-cli` --      |

### Examples

#### Calling a module using the dag-cli command

```bash
dagger call dag-cli --dag-cmds="call -m github.com/shykes/daggerverse/hello hello"
```

### Open a terminal in the container

In this example, a new terminal is open, and from within it's possible to inspect environment variables, and run commands.

```bash
# Open a terminal in the container
dagger call terminal \
  --env-vars="TEST=123,TEST2=456"

# And then, from within the terminal, you can run:
printenv
```

### Call a module's function

In this example, the `hello` module is called, and the `hello` function is executed.

```bash
dagger call use-fn \
--mod-name=github.com/shykes/daggerverse/hello \
--fn="hello" \
--fn-args="--greeting=Yoooo, --name=Alex"
```

### API Reference 📚

For a more detailed information, just run `dagger call --help` and you will get the following output:

````txt
FUNCTIONS
  base                      Base sets the base container for the Dagindag module.
  ctr                       Ctr is the container to use as a base container.
  dag-cli                   DagCLI Allows to execute the Dagger CLI with the given flags.
  src                       Src is the directory that contains all the source code, including the module directory.
  terminal                  Terminal returns a terminal for the container.
  use-fn                    UseFn calls a module with the given function and arguments.
  with-dagger-entry-point   WithDaggerEntryPoint sets the Dagger CLI entry point.
  with-dagger-setup         WithDaggerSetup sets up the container with the Dagger engine.
  with-docker-service       WithDockerService sets up the container with the Docker service.
  with-env-variable         WithEnvVariable sets an environment variable.
  with-env-vars-from-strs   WithEnvVarsFromStrs sets the environment variables for the container.
  with-source               WithSource sets the source directory.

ARGUMENTS
      --ctr Container           ctr is the container to use as a base container.
      --dagger-version string   daggerVersion is the version of the Dagger engine to use, e.g., "v0.11.5
      --docker-version string   dockerVersion is the version of the Docker engine to use, e.g., "24.0
      --src Directory           src is the directory that contains all the source code, including the module directory.

OPTIONS
      --json            Present result as JSON
  -m, --mod string      Path to dagger.json config file for the module or a directory containing that file. Either local path (e.g. "/path/to/some/dir") or a github repo
                        (e.g. "github.com/dagger/dagger/path/to/some/subdir")
  -o, --output string   Path in the host to save the result to```
````
