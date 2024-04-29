# Terragrunt Module for Dagger

![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.0-0f0f19.svg?style=flat-square)

This module provides a Dagger interface for managing Terragrunt operations, facilitating the execution of Terragrunt commands within Dagger-driven environments. It utilizes Terragrunt, a thin wrapper for Terraform that provides extra tools for managing multiple Terraform configurations.

## Configuration 🛠️

Configure this module using the Dagger CLI, specifying paths and versions to tailor the Terragrunt operations to your specific needs:

* ⚙️ **`src`**: Path to the directory containing the Terraform code.
* ⚙️ **`ctr`**: Specifies the container to use. If not provided, a default Terragrunt container based on Alpine is used.
* ⚙️ **`tfVersion`**: The version of Terraform to use (default `1.7.0`).
> **Note**: This module doesnt' set the Terragrunt version, as it is managed by the Terragrunt container. It's set to be compatible with the Terraform version specified.

---

## Features 🎨

| Command or functionality | Command     | Example                                                                                              | Status |
|--------------------------|-------------|------------------------------------------------------------------------------------------------------|--------|
| Terragrunt Version       | **version** | `dagger call --src="mydir/src" version`                                                              | ✅      |
| Terragrunt Help          | **help**    | `dagger call --src="mydir/src" help`                                                                 | ✅      |
| Run Terragrunt Command   | **run**     | `dagger call --src="mydir/src" run plan --envVars="AWS_ACCESS_KEY_ID=xxx,AWS_SECRET_ACCESS_KEY=yyy"` | ✅      |
| Run Terragrunt Run-all   | **run-all** | `dagger call --src="mydir/src" run-all apply --args="--terragrunt-non-interactive"`                  | ✅      |

> **NOTE**: Supports passing custom arguments using the `--args` flag as a comma-separated string, and environment variables with `--envVars` for the Terragrunt command execution.

---

## Usage 🚀

The following examples illustrate how to use the Terragrunt module in Dagger for various Terragrunt commands. Each example demonstrates how to specify the source directory, the specific module, the command to execute, and any additional arguments:

### Basic Usage

```bash
# Run a Terragrunt 'plan' command on a specific module with additional arguments
dagger call --src="../test/testdata" run \
--module=terragrunt/with-dependencies \
--cmd="plan" \
--args="-compact-warnings, -no-color, -lock=false"
```

### Using the GitHub Module
```bash
dagger -m github.com/Excoriate/daggerverse/terragrunt@v1.10.0 call --src="../test/testdata" run \
--module=terragrunt/with-dependencies \
--cmd="plan" \
--args="-compact-warnings, -no-color, -lock=false"

# Run 'apply' for all configurations in a non-interactive mode using the GitHub module
dagger -m github.com/Excoriate/daggerverse/terragrunt@v1.10.0 call --src="../test/testdata" run-all \
--module=terragrunt/with-dependencies \
--cmd="apply" \
--args="-compact-warnings, -no-color, -lock=false, -auto-approve"
```

>**NOTE**: Each module is delivered with a sample project that can be used to test these modules. The ones used in this usage are located in the `test/testdata` directory.
