# 🧰 Gotoolbox Module for Dagger

A powerful [Dagger](https://dagger.io) module for Go development and CI/CD workflows.

## 🛠️ Configuration

Through the [Dagger CLI](https://docs.dagger.io/cli/465058/install) or within your module, you can configure:

- ⚙️ `ctr`: Base container (default: new container created)
- ⚙️ `version`: Go image version (default: `latest`)
- ⚙️ `image`: Go image (default: `golang:alpine`)

## 🌟 Features

| Function                            | Description                 | Example                                                                    |
| ----------------------------------- | --------------------------- | -------------------------------------------------------------------------- |
| 🐳 `Base`                           | Set base image and version  | `dagger call base --image-url golang:1.22-alpine`                          |
| 📦 `WithEnvironmentVariable`        | Set environment variable    | `dagger call with-environment-variable --name GO_ENV --value production`   |
| 📂 `WithSource`                     | Mount source directory      | `dagger call with-source --src . --workdir /app`                           |
| 🔒 `WithSecretAsEnvVar`             | Set secret as env var       | `dagger call with-secret-as-env-var --name API_KEY --secret mysecret`      |
| 📥 `WithDownloadedFile`             | Download and mount file     | `dagger call with-downloaded-file --url https://example.com/file.txt`      |
| 🔄 `WithClonedGitRepo`              | Clone and mount Git repo    | `dagger call with-cloned-git-repo --repo-url https://github.com/user/repo` |
| 🔄 `WithCacheBuster`                | Add cache busting           | `dagger call with-cache-buster`                                            |
| 🛠️ `WithGitInAlpineContainer`       | Install Git                 | `dagger call with-git-in-alpine-container`                                 |
| 🧰 `WithUtilitiesInAlpineContainer` | Install common utilities    | `dagger call with-utilities-in-alpine-container`                           |
| 🖥️ `WithGoPlatform`                 | Set Go platform             | `dagger call with-go-platform --platform linux/amd64`                      |
| 🔧 `WithGoCgoEnabled`               | Enable CGO                  | `dagger call with-go-cgo-enabled`                                          |
| 🚫 `WithCgoDisabled`                | Disable CGO                 | `dagger call with-cgo-disabled`                                            |
| 💾 `WithGoBuildCache`               | Configure Go build cache    | `dagger call with-go-build-cache`                                          |
| 📦 `WithGoModCache`                 | Configure Go mod cache      | `dagger call with-go-mod-cache`                                            |
| 📥 `WithGoInstall`                  | Install Go packages         | `dagger call with-go-install --pkgs github.com/user/pkg@latest`            |
| 🏃 `WithGoExec`                     | Execute Go command          | `dagger call with-go-exec --args build --args ./...`                       |
| 🛠️ `WithGoBuild`                    | Configure Go build          | `dagger call with-go-build --pkg ./cmd/app --race`                         |
| 🔒 `WithGoPrivate`                  | Set GOPRIVATE               | `dagger call with-go-private --private-host example.com`                   |
| 🔧 `WithGoGCCCompiler`              | Install GCC compiler        | `dagger call with-go-gcc-compiler`                                         |
| 📊 `WithGoTestSum`                  | Install GoTestSum           | `dagger call with-go-test-sum`                                             |
| 🚀 `WithGoReleaser`                 | Install GoReleaser          | `dagger call with-go-releaser`                                             |
| 🔍 `WithGoLint`                     | Install golangci-lint       | `dagger call with-go-lint --version v1.60.1`                               |
| 🖥️ `OpenTerminal`                   | Open interactive terminal   | `dagger call open-terminal`                                                |
| 🐚 `RunShell`                       | Run shell command           | `dagger call run-shell --cmd "echo Hello, World!"`                         |
| 🖨️ `PrintEnvVars`                   | Print environment variables | `dagger call print-env-vars`                                               |
| 🔍 `InspectEnvVar`                  | Inspect specific env var    | `dagger call inspect-env-var --key GO_VERSION`                             |
| 🏃 `RunGoCMD`                       | Run Go command              | `dagger call run-go-cmd --cmd test --cmd ./...`                            |
| 🏃 `RunAnyCmd`                      | Run any command             | `dagger call run-any-cmd --cmd go --cmd version`                           |

## Using the Gotoolbox Module 🚀

Refer to the examples in the [**{{.module_name_pkg}}/examples**]({{.module_name_pkg}}/examples) module to see how to use the module's functions.

---

### Usage through the Dagger CLI 🚀

List all the functions available in the module:

```bash
# enter into the module's directory
cd gotoolbox

# list all the functions available in the module
dagger develop && dagger functions
```

Call a function:

```bash
# call a function
# dagger call <function-name> [arguments]
dagger call github.com/excoriate/daggerverse/gotoolbox@version <function-name> [arguments]
```

---

## Testing 🧪

This module includes a [testing]({{.module_name_pkg}}/tests) module that aims to test the functionality of the Gotoolbox module. The tests are written in Go and can be run using the following command:

```bash
## Run the tests using the just command
just test gotoolbox
```

## Developer Experience 🛠️

If you'd like to contribute, mostly we use [Just](https://just.systems) to automate tasks and [Nix](https://nixos.org) to manage the development environment. You can use the following commands to get started:

```bash
# initialize the pre-commit hooks
just init
# run CI or common things locally
just golint gotoolbox
# run the tests
just test gotoolbox
# Run the entire CI tasks locally
just cilocal gotoolbox
```

### Examples (aka Recipes) 🍲

Additionally, this module brings a new [Daggerverse](https://daggerverse.dev/) functionality that allows to automatically generate the module's documentation using an special (sub) module called [**{{.module_name_pkg}}/examples/sdk**]({{.module_name_pkg}}/examples). This module contains a set of examples hat demonstrate how to use the module's functions.

To generate the documentation
It's important to notice that each **example** function in order to be rendered in the documentation, it must be preprocessed by module's name, in this case (camelCase) `gotoolbox`.

> NOTE: The `just` command entails the use of the [**Justfile**](https://just.systems) for task automation. If you don't have it, don't worry, you just need [Nix](https://nixos.org) to run the tasks using the `dev-shell` built-in command: `nix develop --impure --extra-experimental-features nix-command --extra-experimental-features flakes`
