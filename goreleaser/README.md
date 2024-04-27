# GoReleaser Module for Dagger

![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.0-0f0f19.svg?style=flat-square)


A simple [Dagger](https://dagger.io) module that wraps [GoReleaser](https://goreleaser.com) functionality to create and publish Go binaries.

## Configuration 🛠️

Through the [Dagger CLI](https://docs.dagger.io/cli/465058/install), or by using it directly reusing it within your module, you can configure the following options:

* ⚙️ `src`: The path to the Go source code directory.
* ⚙️ `ctr`: The container if it's passed. If not, the module will use the default container.
* ⚙️ `cfgFile`: The path to the GoReleaser configuration file. By default, the module will use the `.goreleaser.yml` file in the root of the project.
* ⚙️ `envVarsFromHost`: A slice of environment variables to pass from the host to the container.

---

## Features 🎨

| Command or functionality | Command      | Example                                                                             | Status |
|--------------------------|--------------|-------------------------------------------------------------------------------------|--------|
| GoReleaser Check         | **check**    | `dagger call --src="mydir/src" check`                                               | ✅      |
| GoReleaser Snapshot      | **snapshot** | `dagger call --src="mydir/src" snapshot --args="--snapshot"`                        | ✅      |
| GoReleaser Release       | **release**  | `dagger call --src="mydir/src" release --args="--rm-dist,--release-notes=notes.md"` | ✅      |

> **NOTE**: Commands support custom arguments using the `--args` flag. Arguments should be provided as a comma-separated string if the command requires customization.

## Using the GoReleaser Module

This module allows you to integrate GoReleaser into your Dagger pipelines easily. To use it, simply specify the source directory and the desired command. The module can execute various GoReleaser functions, including checking configurations, creating snapshots, and executing releases. It's ideal for automating your Go project's release process in a clean, manageable way.

### Additional Information

- **GoReleaser**: GoReleaser is a tool designed to speed up the delivery of Go applications by automating the building and releasing process. More details can be found in the [GoReleaser Documentation](https://goreleaser.com).
- **Dagger**: Dagger is a portable devkit for CI/CD pipelines, allowing you to define your pipeline as code and execute it anywhere. For more on Dagger, visit [Dagger.io](https://dagger.io).

---

## Usage 🚀

  ```bash
dagger call -m github.com/Excoriate/daggerverse/goreleaser@<version> \
check --src="mydir/src"

dagger -m github.com/Excoriate/daggerverse/goreleaser@v1.9.0 call --src="../" release \
--cfg=test/goreleaser/simple-go-app/.goreleaser.yaml --clean --env-vars="GITHUB_TOKEN=$GITHUB_TOKEN" --auto-snapshot

```
