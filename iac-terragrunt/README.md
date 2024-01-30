# IAC Terragrunt 🏗️

![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.9.5-0f0f19.svg?style=flat-square)


A simple _Zenith_ [Dagger](https://dagger.io) module to manage your Terragrunt projects. Made with ❤️ by Alex T.
ss

## Configuration 🛠️

Through the [Dagger CLI](https://docs.dagger.io/cli/465058/install), or by using it directly reusing it within your module, you can configure the following options:

* ⚙️ `version`: The version of Terragrunt to use. If not provided, the latest version will be used. Use the `--with-version` flag or the `WithVersion()` function to set it.
* 📁 `src`: The path to the source code of the module. Use the `--with-source` flag or the `WithSource()` function to set it.
* 🐳 `image`: The container image to use. If not provided, the default one will be used. Use the `--with-image` flag or the `WithImage()` function to set it. It defaults to the lightweight [alpine/terragrunt](https://hub.docker.com/r/alpine/terragrunt) image.
* 🚢 `container`: The dagger container that can be passed to the module. Use the `--with-container` flag or the `WithContainer()` function to set it.

>**NOTE**: This module uses constructors. So, each of these configurations can be set while calling the module. See the **initialization** documentation.

### Configuration on initialization 🏗️

The constructor is a nice way to configure your module. It allows you to set the configuration while calling the module.

```bash
# Calling the module with the CLI within this repository
cd iac-terragrunt
# Here, we're initializing the `src` when calling the module
dagger call -m . --src=$(pwd)/../ run \
--cmds="ls -ltrah, pwd" \
--module="test/iac-terragrunt/testdata/terragrunt-module-simple" \
--stdout
````

---

## Features 🎨

### IAC Terragrunt commands 🏗️

| Command or functionality                        | Command      | Example                                             | Status |
|-------------------------------------------------|--------------|-----------------------------------------------------|--------|
| Run terragrunt **init** (support arguments)     | **init**     | `init --module=mydir/module" --args=-backend=false` | ✅      |
| Run terragrunt **plan** (support arguments)     | **plan**     | `plan --module=mydir/module" --args=-`              | ✅      |
| Run terragrunt **apply** (support arguments)    | **apply**    | `apply --module=mydir/module" --args=-`             | ✅      |
| Run terragrunt **destroy** (support arguments)  | **destroy**  | `destroy --module=mydir/module" --args=`            | ✅      |
| Run terragrunt **validate** (support arguments) | **validate** | `validate --module=mydir/module" --args=`           | ✅      |
| Run terragrunt _hclfmt_ (support arguments)     | **hclfmt**   | `hclfmt --module=mydir/module" --args=`             | ✅      |
>**NOTE**: All these functions supports the flag `--args` to pass arguments to the command. For example, `--args=-backend=false` will be passed as `terragrunt init -backend=false`. The args supported are the native arguments that each terragrunt (technically terraform) command supports.

For composing this module, each of these commands are implemented as a function with an E variant, which means that it can be executed within an existing container. For example, `init` is implemented as `Init()` and `InitE()`, where the second one returns the **Container** and the **Error** in case of any.


### General purpose commands 📜

| Command or functionality                        | Command      | Example                                             | Status |
|-------------------------------------------------|--------------|-----------------------------------------------------|--------|
| Run arbitrary shell commands (one or many)      | **run**      | `run --cmds="ls -ltrah, pwd"`                       | ✅      |
| Run arbitrary terragrunt commands (one or many) | **run-tg**   | `run-tg --cmds="init -backend=false, plan"`         | ✅      |


## Usage

### Dagger CLI

```sh
dagger -m iac-terragrunt call new \
--version="" --src=$(pwd) \
run \
--module="test/iac-terragrunt/testdata/terragrunt-module-simple" \
--cmds="ls -ltrah, cat terragrunt.hcl, terragrunt init, terragrunt plan"
```
