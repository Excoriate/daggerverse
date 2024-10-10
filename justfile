export NIXPKGS_ALLOW_UNFREE := "1"
export NOTHANKS := "1"

default:
  @just --list

# --------------------------------------------------
# Section: Nix Development Environment
# --------------------------------------------------
# This section contains recipes for setting up and managing
# the Nix development environment, including entering the environment,
# cleaning caches, and initializing the project.
# --------------------------------------------------

# Recipe to enter the Nix development environment 🧰
dev:
  @echo "Entering Nix development environment 🧰 ..."
  @nix develop --impure --extra-experimental-features nix-command --extra-experimental-features flakes

# Recipe to clean Go cache, Go modules cache, and Nix/DevEnv/DirEnv cache 🧹
clean-nix-cache:
  @echo "Cleaning Go cache 🧹 ..."
  @go clean -cache
  @echo "Cleaning Go modules cache 🧹 ..."
  @go clean -modcache
  @echo "Cleaning Nix/DevEnv/DirEnv cache 🧹 ..."
  @nix-collect-garbage -d

# Recipe to run pre-commit hooks 🔍
run-hooks:
  @echo "Running pre-commit hooks 🔍 ..."
  @nix-shell -p pre-commit --run "pre-commit run --all-files"
  @echo "Pre-commit hooks passed ✅"

# Recipe to bump version of a module 🔄
bump-version mod bump='minor':
    #!/usr/bin/env bash
    set -euo pipefail

    echo "🔄 Bumping version for {{mod}} module"

    # Verify that the module directory exists and contains a dagger.json file
    if [ ! -d "{{mod}}" ] || [ ! -f "{{mod}}/dagger.json" ]; then
        echo "❌ Module {{mod}} not found or dagger.json missing"
        exit 1
    fi

    # Get the latest tag for this module
    latest_tag=$(git describe --tags --abbrev=0 --match "{{mod}}/*" 2>/dev/null || echo "{{mod}}/v0.0.0")
    current_version=$(echo $latest_tag | sed 's/{{mod}}\/v//')

    # Calculate the new version
    new_version="v$(semver bump {{bump}} "v$current_version")"

    echo "🔢 Current version: v$current_version"
    echo "🆕 New version: $new_version"

    read -p "Proceed with version bump? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ Aborting"
        exit 1
    fi

    # Create and push the new tag
    new_tag="{{mod}}/$new_version"
    git tag -a "$new_tag" -m "Bump {{mod}} to $new_version"
    git push origin "$new_tag"

    echo "✅ Version bumped to $new_version and tag $new_tag created"
    echo "🚀 Tag has been pushed to the remote repository"



# Recipe to run all the tests in the target module 🧪
test mod *args: (reloadmod mod) (reloadtest mod)
  #!/usr/bin/env sh
  set -e
  echo "🚀 Running Dagger module tests..."
  echo "📦 Currently in {{mod}} module 🧪, path=`pwd`"
  test -d {{mod}}/tests || (echo "❌ Module not found" && exit 1)
  cd {{mod}}/tests && dagger functions
  cd {{mod}}/tests && dagger call test-all {{args}}

# Recipe to run all the examples in the target module 📄
examplesgo mod *args: (reloadmod mod)
  #!/usr/bin/env sh
  set -e
  echo "🚀 Running Dagger module examples (Go SDK)..."
  echo "📦 Currently in {{mod}} module 🧪, path=`pwd`"
  test -d {{mod}}/examples/go || (echo "❌ Module examples not found" && exit 1)
  cd {{mod}}/examples/go && dagger call all-recipes {{args}}

# Recipe to run GolangCI Lint 🧹
golint mod *args:
  #!/usr/bin/env sh
  set -e
  echo "Running Go (GolangCI)... 🧹 "
  test -d {{mod}} || (echo "❌ Module not found" && exit 1)
  echo "📦 Currently in {{mod}} module, path=`pwd`/{{mod}}"
  cd ./{{mod}} && nix-shell -p golangci-lint --run "golangci-lint run --config ../.golangci.yml {{args}}"
  echo "🧪 Checking now the tests project ..."
  cd ./{{mod}}/tests && nix-shell -p golangci-lint --run "golangci-lint run --config ../../.golangci.yml {{args}}"
  echo "📄 Checking now the examples project ..."
  cd ./{{mod}}/examples/go && nix-shell -p golangci-lint --run "golangci-lint run --config ../../../.golangci.yml {{args}}"

# Recipe to run the whole CI locally 🚀
cilocal mod: (reloadall mod) (golint mod) (test mod) (examplesgo mod) (ci-module-docs mod)
  #!/usr/bin/env sh
  set -e
  echo "🚀 Running the whole CI locally... 🚀"

# Recipe to validate if the dagger module has the README.md file and the LICENSE file 📄
ci-module-docs mod:
  #!/usr/bin/env sh
  set -e
  echo "🔍 Validating the module documentation..."
  test -f {{mod}}/README.md || (echo "❌ README.md file not found" && exit 1)
  test -f {{mod}}/LICENSE || (echo "❌ LICENSE file not found" && exit 1)
  echo "✅ Module documentation is valid"

# --------------------------------------------------
# Section: Dagger Functions
# --------------------------------------------------
# This section contains recipes for calling functions
# in a certain module.
# --------------------------------------------------

# Recipe to call an specific function from the examples/go project in a certain module 📞
callfnexample mod *args: (check-dagger-pre-requisites mod) (reloadexamples mod)
  @echo "🔧 Calling a function in the 📄 examples/go module [{{mod}}/examples/go]..."
  @echo "📦 Currently in [{{mod}}/examples/go] module, path=`pwd`"
  @cd {{mod}}/examples/go && dagger call {{args}}

# Recipe for dagger call tests in a certain module 🧪
callfntest mod *args: (check-dagger-pre-requisites mod) (reloadtest mod)
  @echo "🔨 Calling a function {{args}} in the 🧪 test module [{{mod}}/tests]..."
  @echo "📦 Currently in [{{mod}}/tests] module, path=`pwd`/{{mod}}/tests"
  @cd {{mod}}/tests && dagger functions
  @cd {{mod}}/tests && dagger call {{args}}

# Recipe to call a certain function by a module's name, passing extra arguments optionally 📞
callfn mod *args: (check-dagger-pre-requisites mod) (reloadmod mod)
  @echo "🔨 Calling a function {{args}} in the module [{{mod}}]..."
  @echo "📂 Currently in [{{mod}}] module, path=`pwd`/{{mod}}"
  @cd {{mod}} && dagger functions
  @cd {{mod}} && dagger call {{args}}

# Recipe to list functions in a certain module 📄
listfns mod *args: (check-dagger-pre-requisites mod) (reloadmod mod)
  @echo "📄 Retrieving available functions for the module..."
  @echo "📦 Currently in [{{mod}}] module, path=`pwd`/{{mod}}"
  @cd {{mod}} && dagger functions

# Recipe to list functions in a test 🧪 module
listfnstest mod *args: (check-dagger-pre-requisites mod) (reloadtest mod)
  @echo "📄 Retrieving available functions for the module..."
  @echo "📦 Currently in [{{mod}}/tests] module, path=`pwd`/{{mod}}/tests"
  @cd {{mod}}/tests && dagger functions {{args}}

# Recipe to list functions in a examples/go 📄 module
listfnsexamples mod *args: (check-dagger-pre-requisites mod) (reloadexamples mod)
  @echo "📄 Retrieving available functions for the module..."
  @echo "📦 Currently in [{{mod}}/examples/go] module, path=`pwd`"
  @cd {{mod}}/examples/go && dagger functions {{args}}

# --------------------------------------------------
# Section: Dagger Maintenance
# --------------------------------------------------
# This section contains recipes for various maintenance
# operations related to Dagger modules.
# --------------------------------------------------

# Recipe that wraps the dagger CLI in a certain module 📦
dagcli mod *args: (check-dagger-pre-requisites mod)
  @echo "🚀 Running Dagger CLI in a certain module..."
  @echo "📦 Currently in [{{mod}}] module, path=`pwd`/{{mod}}"
  @cd {{mod}} && dagger {{args}}

# Recipe to run dagger develop and if the engine gots updated, update the modules 🔄
update-all: (daggy-compile) (check-docker-or-podman)
  @echo "🚀 Developing (or upgrading) all Dagger modules..."
  @.daggerx/daggy/target/release/daggy --task=develop

# --------------------------------------------------
# Section: Creating new modules
# --------------------------------------------------
# This section contains recipes for various operations
# related to Dagger modules, such as calling functions,
# listing functions, and running tests.
# --------------------------------------------------

# Recipe to create a new module using Daggy (a rust CLI tool) 🛠️
create mod with-ci='false' type='full': (daggy-compile) (check-docker-or-podman)
  @echo "🚀 Creating a new {{type}} module of type {{type}}..."
  @.daggerx/daggy/target/release/daggy --task=create --module={{mod}} --module-type={{type}}
  @if [ "{{with-ci}}" = "true" ]; then just cilocal {{mod}}; fi

# Recipe to create a new light module using Daggy 🛠️
createlight mod with-ci='false' type='light': (daggy-compile) (check-docker-or-podman)
  @echo "🚀 Creating a new {{type}} module of type {{type}}..."
  @./.daggerx/daggy/target/release/daggy --task=create --module={{mod}} --module-type={{type}}
  @if [ "{{with-ci}}" = "true" ]; then just cilocal {{mod}}; fi

# --------------------------------------------------
# Section: Daggy Operations
# --------------------------------------------------
# This section contains recipes for compiling and testing
# the Daggy tool.
# --------------------------------------------------

# Recipe to run Daggy tests 🧪
daggy-tests: (daggy-compile)
  @echo "Running Daggy tests 🧪 ..."
  @cd .daggerx/daggy && cargo test

# Recipe to compile Daggy 🔄
daggy-compile:
  @echo "Compiling Daggy 🔄 ..."
  @cd .daggerx/daggy && cargo build --release
  @echo "Daggy compiled successfully 🔄"

# --------------------------------------------------
# Section: Reloading Dagger Modules
# --------------------------------------------------
# This section contains recipes for reloading Dagger
# modules and their tests.
# --------------------------------------------------

# Recipe to reload Dagger module (Dagger Develop) 🔄
reloadmod mod *args: (check-dagger-pre-requisites mod)
  @echo "🚀 Running Dagger development in a given module..."
  @echo "📦 Currently in [{{mod}}] module, path=`pwd`/{{mod}}"
  @cd {{mod}} && dagger develop {{args}}
  @echo "✅ Module reloaded successfully"

# Recipe to reload a Dagger module's tests (Dagger Develop) 🔄
reloadtest mod *args: (check-dagger-pre-requisites mod)
  @echo "🚀 Running Dagger development in a given module's tests..."
  @echo "📦 Currently in [{{mod}}/tests] module, path=`pwd`/{{mod}}/tests"
  @cd {{mod}}/tests && dagger develop {{args}}
  @echo "✅ Module Tests reloaded successfully"

# Recipe to reload the Dagger module's examples (examples/go) 🔄
reloadexamples mod *args: (check-dagger-pre-requisites mod)
  @echo "🚀 Reloading the module's examples..."
  @echo "📦 Currently in {{mod}}/examples/go module, path=`pwd`"
  @test -d {{mod}}/examples/go || (echo "❌ Module examples not found" && exit 1)
  @cd {{mod}}/examples/go && dagger develop {{args}}
  @echo "🚀 Module's examples reloaded successfully"

# Recipe to reload Dagger module and its underlying tests (Dagger Develop & Dagger Call/Functions) 🔄
reloadall mod *args: (reloadmod mod) (reloadtest mod) (reloadexamples mod)
  @echo "🔄 Reloading all the module, tests and examples... [{{mod}}]"
  @echo "🚀 Module reloaded successfully"


# --------------------------------------------------
# Section: Utilities
# --------------------------------------------------
# This section contains recipes for utilities that
# are used to validate if Docker or Podman is running
# and to validate if a given directory is a Dagger module.
# --------------------------------------------------

# Recipe to check if Docker or Podman is running 🔄
check-docker-or-podman:
  #!/usr/bin/env sh
  set -e

  if command -v docker > /dev/null 2>&1; then
    if ! docker info > /dev/null 2>&1; then
      echo "❌ Docker is installed but not running. Please start Docker and try again."
      exit 1
    else
      echo "✅ Docker is running."
    fi
  elif command -v podman > /dev/null 2>&1; then
    if ! podman info > /dev/null 2>&1; then
      echo "❌ Podman is installed but not running. Please start Podman and try again."
      exit 1
    else
      echo "✅ Podman is running."
    fi
  else
    echo "❌ Neither Docker nor Podman is installed. Please install one of them and try again."
    exit 1
  fi

# Recipe that validate if it's an actual Dagger module
is-dagger-module mod:
  @echo "🔍 Validating if [{{mod}}] is a Dagger module..."
  @test -d {{mod}} || (echo "❌ Module not found at path=`pwd`/{{mod}}" && exit 1)
  @test -f {{mod}}/dagger.json || (echo "❌ dagger.json not found in module at path=`pwd`/{{mod}}. Not a Dagger module." && exit 1)
  @echo "✅ [{{mod}}] is a Dagger module"

# Recipe that check Dagger pre-requisites
check-dagger-pre-requisites mod: (check-docker-or-podman) (is-dagger-module mod)