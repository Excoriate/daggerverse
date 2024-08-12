# Module {{.module_name}} for Dagger


A simple [Dagger](https://dagger.io) _place the description of the module here_

## Configuration 🛠️

Through the [Dagger CLI](https://docs.dagger.io/cli/465058/install), or by using it directly within your module, you can configure the following options:

* ⚙️ `ctr`: The container to use as a base container. If not specified, a new container is created.
* ⚙️ `version`: The version of the Go image to use. Defaults to `latest`.
* ⚙️ `image`: The Go image to use. Defaults to `golang:alpine`.

### Structure 🏗️
```text
{{.module_name_pkg}} // main module
├── .gitattributes
├── .gitignore
├── LICENSE
├── README.md
├── apis.go
├── cloud.go
├── commands.go
├── common.go
├── config.go
├── dagger.json
├── examples // Sub modules that represent examples of the module's functions with each SDK
│   └── go
│       ├── .gitattributes
│       ├── .gitignore
│       ├── dagger.json
│       ├── go.mod
│       ├── go.sum
│       ├── main.go
│       └── testdata
│           └── common
│               ├── README.md
│               └── test-file.yml
├── go.mod
├── go.sum
├── main.go
└── tests // Sub module that represent tests of the module's functions 
    ├── .gitattributes
    ├── .gitignore
    ├── dagger.json
    ├── go.mod
    ├── go.sum
    ├── main.go
    └── testdata
        └── common
            ├── README.md
            └── test-file.yml

```
>NOTE: This structure comes out of the box if it's generated through **Daggy**. Just run `just create <module-name>` and you'll get the structure.

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

### Examples (aka Recipes) 🍲

Additionally, this module brings a new [Daggerverse](https://daggerverse.dev/) functionality that allows to automatically generate the module's documentation using an special (sub) module called [**{{.module_name_pkg}}/examples/sdk**]({{.module_name_pkg}}/examples). This module contains a set of examples hat demonstrate how to use the module's functions. 

To generate the documentation
It's important to notice that each **example** function in order to be rendered in the documentation, it must be preprocessed by module's name, in this case (camelCase) `{{.module_name}}`.

>NOTE: The `just` command entails the use of the [**Justfile**](https://just.systems) for task automation. If you don't have it, don't worry, you just need [Nix](https://nixos.org) to run the tasks using the `dev-shell` built-in command: `nix develop --impure --extra-experimental-features nix-command --extra-experimental-features flakes`
