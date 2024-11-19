# Terratest 🧪

![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.0-0f0f19.svg?style=flat-square)

A simple [Dagger](https://dagger.io) module to run tests on Terraform modules using [Terratest](https://terratest.gruntwork.io/).

## Configuration 🛠️

Through the [Dagger CLI](https://docs.dagger.io/cli/465058/install), or by using it directly reusing it within your module, you can configure the following options:

- ⚙️ `version`: The version of the underlying [Go](https://golang.org/) image to use. Defaults to `1.22.0-alpine3.19`.
- ⚙️ `tfVersion`: The version of Terraform to use. Defaults to `1.6.0`.

> NOTE: The terraform binary is installed on top of the Go image, so the `tfVersion` is the version of the Terraform binary to install.

---

## Features 🎨

| Command or functionality                                 | Command | Example                                                                 | Status |
| -------------------------------------------------------- | ------- | ----------------------------------------------------------------------- | ------ |
| Run a Terratest test with optional arguments             | **run** | `dagger call --src="." run --test-dir="test/demo-1"`                    | ✅     |
| Run a Terratest test with verbose and coverage reporting | **run** | `dagger call --src="." run --test-dir="test/demo-1" --args="-v -cover"` | ✅     |

> **NOTE**: The `run` command is used to execute Terratest tests within a specified directory by using the `--test-dir` argument. Additional go test arguments can be passed through the `--args=` flag to customize test execution. This feature allows for more detailed output or coverage analysis.

## Usage 🚀

To run Terratest tests from within your project directory:

```bash
dagger call --src="." run --test-dir="test/my-test-directory"
```

With additional arguments:

```bash
dagger call --src="." run --test-dir="test/my-test-directory" --args="-v -cover"
```

---
