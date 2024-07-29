# Module {{.module_name}} for Dagger

![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.0-0f0f19.svg?style=flat-square)

A simple [Dagger](https://dagger.io) _place the description of the module here_

## Configuration 🛠️

Through the [Dagger CLI](https://docs.dagger.io/cli/465058/install), or by using it directly within your module, you can configure the following options:

* ⚙️ `ctr`: The container to use as a base container. If not specified, a new container is created.
* ⚙️ `version`: The version of the Go image to use. Defaults to `latest`.
* ⚙️ `image`: The Go image to use. Defaults to `golang:alpine`.

---

## Features 🎨

| Command or functionality  | Command | Example                     | Status |
|---------------------------|---------|-----------------------------|--------|
| Add your feature **here** | **run** | `dagger call <my function>` | ✅      |


## Using the {{.module_name}} Module 🚀

_Place the description of the module here_

---

### Usage through the Dagger CLI 🚀

List all the functions available in the module:

  ```bash
  # enter into the module's directory
  cd {{.module_name}}

  # list all the functions available in the module
  dagger develop && dagger functions
```

Call a function:

  ```bash
  # call a function
  # dagger call <function-name> [arguments]
  dagger call github.com/excoriate/daggerverse/{{.module_name}}@version <function-name> [arguments]
```

---

## Testing 🧪

This module includes a [testing]({{.module_name_pkg}}/tests) module that aims to test the functionality of the {{.module_name}} module. The tests are written in Go and can be run using the following command:

```bash
## Run the tests using the just command
just test {{.module_name}}
```

## Developer Experience 🛠️

If you'd like to contribute, mostly we use [Just](https://just.systems) to automate tasks and [Nix](https://nixos.org) to manage the development environment. You can use the following commands to get started:

```bash
# initialize the pre-commit hooks
just init
# run CI or common things locally
just golint {{.module_name}}
# run the tests
just test {{.module_name}}
# Run the entire CI tasks locally
just cilocal {{.module_name}}
```

Additionally, this module brings a new [Daggerverse](https://daggerverse.dev/) functionality that allows to automatically generate the module's documentation using an special (sub) module called [**examples**/]({{.module_name_pkg}}/examples). This module contains a set of examples hat demonstrate how to use the module's functions. To generate the documentation

>NOTE: The `just` command entails the use of the [**Justfile**](https://just.systems) for task automation. If you don't have it, don't worry, you just need [Nix](https://nixos.org) to run the tasks using the `dev-shell` built-in command: `nix develop --impure --extra-experimental-features nix-command --extra-experimental-features flakes`
