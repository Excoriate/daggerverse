<h1 align="center">
  <img alt="logo" src="docs/logo/daggerverse-logo-nobackground.png" width="450px"/><br/>
</h1>

## <h1 align="center">Daggerverse Modules 📦</h1>

[![🏗️ CI CodeGen Daggy](https://github.com/Excoriate/daggerverse/actions/workflows/ci-daggy-codegen.yml/badge.svg)](https://github.com/Excoriate/daggerverse/actions/workflows/ci-daggy-codegen.yml)[![CI module-template 🧹](https://github.com/Excoriate/daggerverse/actions/workflows/ci-mod-module-template.yaml/badge.svg)](https://github.com/Excoriate/daggerverse/actions/workflows/ci-mod-module-template.yaml)

| Module                                         | Status | What it does?                                                                   |
| ---------------------------------------------- | ------ | ------------------------------------------------------------------------------- |
| [Terraform](terraform/README.md)               | ✅     | 🌍 Run [Terraform](https://www.terraform.io) commands.                          |
| [Terratest](terratest/README.md)               | ✅     | 🧪 Run [Terratest](https://terratest.gruntwork.io) commands.                    |
| [GitLab CICD Vars](gitlab-cicd-vars/README.md) | ✅     | ⚙️ Manage GitLab CI/CD variables.                                               |
| [GoReleaser](goreleaser/README.md)             | ✅     | 🚀 Wraps [GoReleaser](https://goreleaser.com) functionality as a dagger module. |
| [TFLint](tflint/README.md)                     | ✅     | 🔍 Run [TFLint](https://github.com/terraform-linters/tflint) commands.          |
| [GoTest](gotest/README.md)                     | ✅     | 🔋 A batteries-included way to run Go Tests (using Go Test, and/or GoTestsum).  |
| [Terragrunt](terragrunt/README.md)             | ✅     | 🔋 A batteries-included way to run Terragrunt commands.                         |
| [Gotoolbox](gotoolbox/README.md)               | ✅     | 🛠️ A toolbox for various Go utilities.                                          |

---

## Contributions 🤝

This is a mono-repo, and each module is a separate Go module. To contribute to a module, first of all read the [contribution guidelines](./CONTRIBUTING.md).

## Tools 🛠️

- [Go](https://golang.org)
- [Nix](https://nixos.org) (optional, mostly for maintainers)
- [Just](https://github.com/casey/just) (optional, mostly for maintainers)

## Scaffolding 🧰

New modules can be generated using **Daggy**, a [Rust](https://www.rust-lang.org) CLI tool that generates the boilerplate code for a new module. To use Daggy and create a new module, just execute:

### Dagger Module Types

- **Full**: A full-featured module with all the built-in functions and examples.

```bash
just create <module-name>
```

- **Light**: A minimal version of the module with the most important built-in functions.

```bash
just createlight <module-name>
```

### Dagger Module Structure 🧱

#### Module Template (Full)

```text
module-template
├── .gitattributes
├── .gitignore
├── LICENSE
├── README.md
├── apis.go
├── clis.go
├── cloud.go
├── commands.go
├── config.go
├── container_base.go
├── content.go
├── dagger.go
├── dagger.json
├── envvars.go
├── err.go
├── examples
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
├── golang.go
├── http.go
├── iac_terraform.go
├── iac_terragrunt.go
├── install.go
├── main.go
├── server_go.go
├── tests
│   ├── .gitattributes
│   ├── .gitignore
│   ├── apis.go
│   ├── cli.go
│   ├── cloud.go
│   ├── container_base.go
│   ├── dagger.go
│   ├── dagger.json
│   ├── err.go
│   ├── go.mod
│   ├── go.sum
│   ├── golang.go
│   ├── http.go
│   ├── iac_terraform.go
│   ├── iac_terragrunt.go
│   ├── install.go
│   ├── main.go
│   ├── server_go.go
│   ├── testdata
│   │   ├── apko-presets
│   │   │   ├── base-alpine.yaml
│   │   │   └── base-wolfi.yaml
│   │   ├── common
│   │   │   ├── README.md
│   │   │   └── test-file.yml
│   │   ├── golang-server-http
│   │   │   ├── Dockerfile
│   │   │   ├── Makefile
│   │   │   ├── go.mod
│   │   │   ├── go.sum
│   │   │   └── main.go
│   │   └── golang
│   │       └── main.go
│   └── vcs.go
└── vcs.go
```

#### Module Template (Light)

```text
module-template-light
├── .gitattributes
├── .gitignore
├── LICENSE
├── README.md
├── apis.go
├── commands.go
├── config.go
├── container_base.go
├── content.go
├── dagger.json
├── err.go
├── examples
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
├── golang.go
├── install.go
├── main.go
└── tests
    ├── .gitattributes
    ├── .gitignore
    ├── apis.go
    ├── container_base.go
    ├── dagger.json
    ├── err.go
    ├── go.mod
    ├── go.sum
    ├── golang.go
    ├── install.go
    ├── main.go
    └── testdata
        ├── apko-presets
        │   ├── base-alpine.yaml
        │   └── base-wolfi.yaml
        ├── common
        │   ├── README.md
        │   └── test-file.yml
        ├── golang-server-http
        │   ├── Dockerfile
        │   ├── Makefile
        │   ├── go.mod
        │   ├── go.sum
        │   └── main.go
        └── golang
            └── main.go
```

## Testing and CI 🧪

Currently, the following checks are executed on each module:

- [x] Run [GolangCI-Lint](https://golangci.com/) on the module, tests, and examples
- [x] Run `dagger call` on the module, tests, and examples
- [x] Run all the tests in the `tests/` Dagger module.
- [x] Run all the recipes in the `examples/go` Dagger module.

To run the CI checks locally, just execute:

```bash
just ci <module-name>
```

To run only the tests in your module, just execute:

```bash
just test <module-name>
```

To run only the Go lint checks in your module, just execute:

```bash
just lintall <module-name>
```
