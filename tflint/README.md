# TFLint Module for Dagger

![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.0-0f0f19.svg?style=flat-square)

A streamlined [Dagger](https://dagger.io) module that encapsulates [TFLint](https://github.com/terraform-linters/tflint) functionality, enabling static analysis of Terraform code directly within your Dagger workflows.

## Configuration 🛠️

Configure the module using the Dagger CLI or by integrating it into your own modules. Configuration options include:

* ⚙️ **`src`**: The path to the Terraform code directory.
* ⚙️ **`ctr`**: Specifies the container to use. If not provided, a default TFLint container is used.

---

## Features 🎨

| Command or functionality | Command    | Example                                                                                                   | Status |
|--------------------------|------------|-----------------------------------------------------------------------------------------------------------|--------|
| TFLint Version Check     | **version**| `dagger call --src="mydir/src" version`                                                                   | ✅      |
| TFLint Initialize        | **init**   | `dagger call --src="mydir/src" init --cfg="custom.tflint.hcl"`                                            | ✅      |
| TFLint Lint              | **lint**   | `dagger call --src="mydir/src" lint --init --cfg=".tflint.hcl" --args="--enable-rule=terraform_unused_declarations"` | ✅      |

> **NOTE**: The `--args` flag allows passing custom arguments as a comma-separated string, applicable if the command supports additional customization.

## Using the TFLint Module

Easily integrate TFLint into your Dagger pipelines to perform static code analysis on your Terraform configurations. Specify the source directory and the command you wish to run. This module supports various TFLint functionalities, such as version checking, initializing configurations, and running lint checks.

### Additional Information

* **TFLint**: A Terraform linter focused on possible errors, best practices, and deprecated syntax. More details can be found in the [TFLint Documentation](https://github.com/terraform-linters/tflint).
* **Dagger**: A portable development kit for CI/CD pipelines that allows defining and executing pipelines as code anywhere. Learn more at [Dagger.io](https://dagger.io).

By using this module, you can enhance the quality and consistency of your Terraform code, leveraging TFLint's capabilities within your existing Dagger-based workflows.

---

## Usage 🚀

```bash
# Check TFLint version
dagger call -m github.com/Excoriate/daggerverse/tflint@master version --src="mydir/src"

# Initialize TFLint
dagger call -m github.com/Excoriate/daggerverse/tflint@master init --src="mydir/src" --cfg="custom.tflint.hcl"

# Run TFLint with custom configurations and initializations
dagger call -m github.com/Excoriate/daggerverse/tflint@master lint --src="mydir/src" --init --cfg=".tflint.hcl" --args="--enable-rule=terraform_unused_declarations"
