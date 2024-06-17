# GoTest Module for Dagger

![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.0-0f0f19.svg?style=flat-square)

A simple [Dagger](https://dagger.io) module that wraps Go testing functionality to run Go tests within a containerized environment.

## Configuration 🛠️

Through the [Dagger CLI](https://docs.dagger.io/cli/465058/install), or by using it directly within your module, you can configure the following options:

* ⚙️ `src`: The path to the Go source code directory.
* ⚙️ `ctr`: The container to use as a base container. If not specified, a new container is created.
* ⚙️ `version`: The version of the Go image to use. Defaults to `latest`.
* ⚙️ `image`: The Go image to use. Defaults to `golang:alpine`.
* ⚙️ `envVarsFromHost`: Environment variables to pass from the host to the container.

---

## Features 🎨

| Command or functionality           | Command       | Example                                                                                                                                                        | Status |
|------------------------------------|---------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------|--------|
| Run Go Tests                       | **run**       | `dagger call run-go-test --packages="./..." --enableVerbose=true --race=true --src="my-code/src"`                                                              | ✅      |
| Run GoTestSum                      | **gotestsum** | `dagger call run-go-test-sum --src="mydir/src" --race=true --testFlags="-json" --goTestSumFlags="--format=short-verbose" --format="short" --enablePretty=true` | ✅      |


## Using the GoTest Module

This module allows you to integrate Go testing into your Dagger pipelines easily. To use it, simply specify the source directory and the desired command. The module can execute various Go testing functions, including running tests and using gotestsum for advanced test result formatting.

### Additional Information

* **GoTestSum**: GoTestSum is a tool for running and summarizing Go tests. More details can be found in the [GoTestSum Documentation](https://github.com/gotestyourself/gotestsum).
* **Dagger**: Dagger is a portable devkit for CI/CD pipelines, allowing you to define your pipeline as code and execute it anywhere. For more on Dagger, visit [Dagger.io](https://dagger.io).

---

## Usage 🚀

  ```bash
  dagger call run-go-test --src="../../test/testdata/gotest" \
--test-flags="--json" --enable-cache;

 dagger call --verbose run-go-test-sum --src="../../test/testdata/gotest" \
--enable-pretty --enable-cache;
```
