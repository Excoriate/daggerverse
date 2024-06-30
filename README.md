<h1 align="center">
  <img alt="logo" src="docs/logo/daggerverse-logo-nobackground.png" width="450px"/><br/>
</h1>

<h1 align="center">Daggerverse Modules 📦</h1>
---

| Module                                         | Status | What it does?                                                                |
|------------------------------------------------|--------|------------------------------------------------------------------------------|
| [Terraform](terraform/README.md)               | ✅      | Run [Terraform](https://www.terraform.io) commands.                          |
| [Terratest](terratest/README.md)               | ✅      | Run [Terratest](https://terratest.gruntwork.io) commands.                    |
| [GitLab CICD Vars](gitlab-cicd-vars/README.md) | ✅      | Manage GitLab CI/CD variables.                                               |
| [GoReleaser](goreleaser/README.md)             | ✅      | Wraps [GoReleaser](https://goreleaser.com) functionality as a dagger module. |
| [TFLint](tflint/README.md)                     | ✅      | Run [TFLint](https://github.com/terraform-linters/tflint) commands.          |
| [GoTest](gotest/README.md)                     | ✅      | A simplify way to run Go Tests (using Go Test, and/or GoTestsum).            |
| [Terragrunt](terragrunt/README.md)             | ✅      | A simple [Terragrunt](https://terragrunt.gruntwork.io) module.               |

---

## Contributions 🤝

This is a mono-repo, and each module is a separate Go module. To contribute to a module, first of all read the [contribution guidelines](./CONTRIBUTING.md).

## Pre-requisites 📋

- [Go](https://golang.org)
- [Nix](https://nixos.org) (optional, mostly for maintainers)
- [Just](https://github.com/casey/just) (optional, mostly for maintainers)

### What about new modules? 🤔

New modules can be generated using **Daggy**, a [Rust](https://www.rust-lang.org) CLI tool that generates the boilerplate code for a new module. To use Daggy and create a new module, just execute:

```bash
# It generates a new module with the name <module-name>
just create <module-name>

# Run the CI on the new module
just cilocal <module-name>
```

>**NOTE**: See the [Module Template](./module-template) for more information for the new module structure, and the boilerplate code that's generated.
