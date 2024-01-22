# IAC Terragrunt 🏗️

## Description

A simple _Zenith_ [Dagger](https://dagger.io) module to manage your Terragrunt projects. Made with ❤️ by Alex T.

## Features

| Command or functionality                                | Is Implement? |
|---------------------------------------------------------|---------------|
| Provide a custom `version` for Terragrunt               | ✅             |
| Set an specific path for the target Terragrunt `module` | ✅             |

## Usage

### Dagger CLI

```sh
dagger -m iac-terragrunt call new \
--version="" --src=$(pwd) \
run \
--module="test/iac-terragrunt/testdata/terragrunt-module-simple" \
--cmds="ls -ltrah, cat terragrunt.hcl, terragrunt init, terragrunt plan"
```
