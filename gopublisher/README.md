# GoPublisher Module

![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.0-0f0f19.svg?style=flat-square)


A simple [Dagger](https://dagger.io) module that [publishes Go packages and modules](https://go.dev/doc/modules/publishing) into the Golang public registry

## Configuration 🛠️

Through the [Dagger CLI](https://docs.dagger.io/cli/465058/install), or by using it directly reusing it within your module, you can configure the following options:

* ⚙️ `version`: This argument (to the module's constructor) dictates what's the base [Golang version](https://golang.org/dl/) to use. By default, it's set to `1.22.3`.
* ⚙️ `ctr`: The container if it's passed. If not, the module will use the default container.

> **NOTE**: This configuration is available through the [module's constructor](https://docs.dagger.io/manuals/developer/go/520657/constructor/).

---

## Features 🎨

### Commands and Functionalities  📜

| Command or functionality         | Command                     | Example                                                                                   | Status |
|----------------------------------|-----------------------------|-------------------------------------------------------------------------------------------|--------|
| Go Publisher GoModPath           | **go-mod-path**             | `dagger call --src="mydir/src" go-mod-path`                                               | ✅      |
| Go Publisher Publish             | **go-mod-publish**          | `dagger call --src="mydir/src" go-mod-publish --tag="v1.0.0"`                             | ✅      |
| Go Publisher GoModVersion        | **go-mod-version**          | `dagger call --src="mydir/src" go-mod-version`                                            | ✅      |
| Go Publisher Terminal            | **terminal**                | `dagger call --src="mydir/src" terminal`                                                  | ✅      |


### API Reference 📚

For a more detailed information, just run `dagger call --help` and you will get the following output:

```bash
USAGE
  dagger call [options] [arguments] <function>

EXAMPLES
  dagger call test
  dagger call build -o ./bin/myapp
  dagger call lint stdout

FUNCTIONS
  base                      Base sets the base container for gopublisher.
  ctr                       Ctr is the container to use as a base container for gopublisher, if it's passed, it's used as the base container.
  go-mod-path               GoModPath returns the module path.
  go-mod-publish            GoModPublish publishes the module to the registry.
  go-mod-version            GoModVersion returns the module version by running git describe --tags --abbrev=0.
  src                       Src is the directory that contains all the source code, including the module directory.
  terminal                  Terminal returns a terminal for the container.
  with-cgodisabled          WithCGODisabled sets the CGO_ENABLED environment variable to 0.
  with-curl                 WithCURL installs or setup the container with the curl binary.
  with-env-variable         WithEnvVariable sets an environment variable.
  with-env-vars-from-strs   WithEnvVarsFromStrs sets the environment variables for the container.
  with-git                  WithGit installs or setup the container with the git binary.
  with-source               WithSource sets the source directory.

ARGUMENTS
      --version string   version is the version of Go that the publisher module will use, e.g., "1.22.0".

OPTIONS
      --json            Present result as JSON
  -m, --mod string      Path to dagger.json config file for the module or a directory containing that file. Either local path (e.g.
                        "/path/to/some/dir") or a github repo (e.g. "github.com/dagger/dagger/path/to/some/subdir")
  -o, --output string   Path in the host to save the result to

INHERITED OPTIONS
  -d, --debug             show debug logs and full verbosity
      --progress string   progress output format (auto, plain, tty) (default "auto")
  -s, --silent            disable terminal UI and progress output
  -v, --verbose count     increase verbosity (use -vv or -vvv for more)
  ````

---

## Usage in CI (Like GitHub Actions) 🚀

```yaml
  publish-go-module:
    name: Publish Go Module with GoPublisher Dagger
    needs: release-please
    runs-on: ubuntu-latest
    if: ${{ needs.release-please.outputs.releases_created }}
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
      - name: Publish to Go.pkg.dev
        uses: dagger/dagger-for-github@v5
        with:
          verb: call
          module: github.com/Excoriate/daggerverse/gopublisher@v1.14.0
          args: go-mod-publish --src="."
          version: 0.11.6

```
