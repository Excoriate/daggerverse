<h1 align="center">Daggerverse</h1>

<p align="center">
  <img src="docs/img/daggerverse-logo.jpg" alt="daggerverse-logo.png">
Set of Dagger modules
that serve different purposes;
with a slight deviation for infrastructure automation;
made with ❤️.

</p>



| Module                                     | Status | What it does?                                                                |
|--------------------------------------------|--------|------------------------------------------------------------------------------|
| [IAC Terragrunt](iac-terragrunt/README.md) | ✅      | Run [Terragrunt](https://terragrunt.gruntwork.io) commands.                  |
| [Terraform](terraform/README.md)           | ✅      | Run [Terraform](https://www.terraform.io) commands.                          |
| [Terratest](terratest/README.md)           | ✅      | Run [Terratest](https://terratest.gruntwork.io) commands.                    |
| [GitLab CICD Vars](gitlab-cicd-vars/README.md) | ✅  | Manage GitLab CI/CD variables.                                              |
| [GoReleaser](goreleaser/README.md)         | ✅      | Wraps [GoReleaser](https://goreleaser.com) functionality as a dagger module. |


>**NOTE**: ⚠️ These modules are experimental, feel free to open an issue for any requests or bug report.

---

## How to contribute 🤔 ?

First, read the [contribution guidelines](./CONTRIBUTING.md). Then, if you're already a ninja, it means you enjoy using [Taskfile](https://taskfile.dev) so just run:

```sh
# This initialises the hooks, and ensure you're always using their latest version.
task pc-init

# This is just a check. It will run all the checks on the codebase.
task pc-run
```

The current workflows in [GitHub Actions](./.github/workflows) will do the rest ;).
