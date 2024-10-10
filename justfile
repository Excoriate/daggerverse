export NIXPKGS_ALLOW_UNFREE := "1"
export NOTHANKS := "1"

default:
  @just --list

# Recipe to run your development environment commands 🧰
dev:
  @echo "Entering Nix development environment 🧰 ..."
  @nix develop --impure --extra-experimental-features nix-command --extra-experimental-features flakes

# Recipe to clean Go cache, Go modules cache, and Nix/DevEnv/DirEnv cache 🧹
clean-cache:
  @echo "Cleaning Go cache 🧹 ..."
  @go clean -cache
  @echo "Cleaning Go modules cache 🧹 ..."
  @go clean -modcache
  @echo "Cleaning Nix/DevEnv/DirEnv cache 🧹 ..."
  @nix-collect-garbage -d

# Recipe to initialize the project 🚀
init:
  @echo "Initializing the project 🚀 ..."
  @nix-shell -p pre-commit --run "pre-commit install --hook-type pre-commit"
  @echo "Pre-commit hook installed ✅"
  @nix-shell -p pre-commit --run "pre-commit install --hook-type pre-push"
  @echo "Pre-push hook installed ✅"
  @nix-shell -p pre-commit --run "pre-commit install --hook-type commit-msg"
  @echo "Commit-msg hook installed ✅"
  @nix-shell -p pre-commit --run "pre-commit autoupdate"
  @echo "Pre-commit hooks updated to the latest version 🔄"

# Recipe to run pre-commit hooks 🔍
precommit:
  @echo "Running pre-commit hooks 🔍 ..."
  @nix-shell -p pre-commit --run "pre-commit run --all-files"
  @echo "Pre-commit hooks passed ✅"

# Recipe to run Dagger module 📦
dc mod *args:
  #!/usr/bin/env sh
  set -e
  echo "🚀 Running Dagger module..."
  echo "📦 Currently in {{mod}} module, path=`pwd`"
  test -d {{mod}} || (echo "❌ Module not found" && exit 1)
  cd {{mod}} && dagger call {{args}}

# Recipe to run Dagger module tests 🧪
dct mod *args:
  #!/usr/bin/env sh
  set -e
  echo "🧪 Running Dagger module tests..."
  echo "🧪 Currently in {{mod}} module, path=`pwd`"
  test -d {{mod}}/tests || (echo "❌ Module not found" && exit 1)
  cd {{mod}}/tests && dagger call {{args}}

# Recipe to run Dagger module examples 📄
dce mod *args:
  #!/usr/bin/env sh
  set -e
  echo "📄 Running Dagger module examples ..."
  echo "🧪 Currently in {{mod}} module, path=`pwd`"
  test -d {{mod}}/examples/go || (echo "❌ Module examples not found" && exit 1)
  cd {{mod}}/examples/go && dagger call {{args}}

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

# Recipe to reload Dagger module (Dagger Develop) 🔄
reloadmod mod *args:
  #!/usr/bin/env sh
  set -e
  echo "🚀 Running Dagger development in a given module..."
  echo "📦 Currently in {{mod}} module, path=`pwd`"
  test -d {{mod}} || (echo "❌ Module not found" && exit 1)
  if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
  fi
  cd {{mod}} && dagger develop {{args}}
  echo "✅ Module reloaded successfully"

# Recipe to reload a Dagger module's tests (Dagger Develop) 🔄
reloadtest mod *args:
  #!/usr/bin/env sh
  set -e
  echo "🚀 Running Dagger development in a given module's tests..."
  echo "📦 Currently in {{mod}}/tests module, path=`pwd`"
  test -d {{mod}} || (echo "❌ Module not found" && exit 1)
  cd {{mod}}/tests && dagger develop {{args}}
  echo "✅ Module Tests reloaded successfully"

# Recipe to reload Dagger module and its underlying tests (Dagger Develop & Dagger Call/Functions) 🔄
reloadall mod *args:
  #!/usr/bin/env sh
  set -e
  echo "🔄 Reloading Dagger module and also the tests..."
  echo "📦 Currently in {{mod}} module, path=`pwd`"
  test -d {{mod}} || (echo "❌ Module not found" && exit 1)
  cd {{mod}} && dagger develop {{args}}
  cd tests && dagger develop {{args}}
  cd ../examples/go && dagger develop {{args}}
  echo "🚀 Module reloaded successfully"
  echo "🕵️ Inspecting the module..."
  cd .. && dagger call && dagger functions

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

# Recipe to create a new module using Daggy (a rust CLI tool) 🛠️
create mod with-ci='false' type='full':
  #!/usr/bin/env sh
  set -e
  echo "🚀 Creating a new {{type}} module of type {{type}}..."
  cd .daggerx/daggy && cargo build --release
  .daggerx/daggy/target/release/daggy --task=create --module={{mod}} --module-type={{type}}
  if [ "{{with-ci}}" = "true" ]; then just cilocal {{mod}}; fi

# Recipe to create a new light module using Daggy 🛠️
createlight mod with-ci='false' type='light':
  #!/usr/bin/env sh
  set -e
  echo "🚀 Creating a new {{type}} module of type {{type}}..."
  cd .daggerx/daggy && cargo build --release
  .daggerx/daggy/target/release/daggy --task=create --module={{mod}} --module-type={{type}}
  if [ "{{with-ci}}" = "true" ]; then just cilocal {{mod}}; fi

# Recipe to validate if the dagger module has the README.md file and the LICENSE file 📄
ci-module-docs mod:
  #!/usr/bin/env sh
  set -e
  echo "🔍 Validating the module documentation..."
  test -f {{mod}}/README.md || (echo "❌ README.md file not found" && exit 1)
  test -f {{mod}}/LICENSE || (echo "❌ LICENSE file not found" && exit 1)
  echo "✅ Module documentation is valid"

# Recipe for dagger call 📞
call mod *args:
  #!/usr/bin/env sh
  set -e
  echo "🚀 Running Dagger call..."
  echo "📦 Currently in {{mod}} module, path=`pwd`"
  test -d {{mod}} || (echo "❌ Module not found" && exit 1)
  cd {{mod}} && dagger call {{args}}

# Recipe for dagger call tests in a certain module 🧪
calltests mod *args: (reloadtest mod)
  #!/usr/bin/env sh
  set -e
  echo "🚀 Running Dagger call tests..."
  echo "🧪 Currently in {{mod}} module, path=`pwd`"
  test -d {{mod}}/tests || (echo "❌ Module not found" && exit 1)
  cd {{mod}}/tests && dagger functions
  cd {{mod}}/tests && dagger call {{args}}
# Recipe to run dagger develop in all modules 🔄
develop-all:
  #!/usr/bin/env sh
  set -e
  echo "🚀 Developing (or upgrading) all Dagger modules..."
  cd .daggerx/daggy && cargo build --release
  .daggerx/daggy/target/release/daggy --task=develop

# Recipe that wraps the dagger CLI in a certain module 📦
dag mod *args:
  #!/usr/bin/env sh
  set -e
  echo "🚀 Running Dagger CLI in a certain module..."
  echo "📦 Currently in {{mod}} module, path=`pwd`"
  test -d {{mod}} || (echo "❌ Module not found" && exit 1)
  cd {{mod}} && dagger {{args}}

# Recipe to call a certain function by a module's name, passing extra arguments optionally 📞
callfn mod *args:
  #!/usr/bin/env sh
  set -e
  echo "🔧 Calling a function in a certain module..."
  echo "📦 Currently in {{mod}} module, path=`pwd`"
  test -d {{mod}} || (echo "❌ Module not found" && exit 1)
  cd {{mod}} && dagger functions
  cd {{mod}} && dagger call {{args}}
# Recipe to list functions in a certain module 📄
listfns mod *args:
  #!/usr/bin/env sh
  set -e
  echo "📄 Listing functions in a certain module..."
  echo "📦 Currently in {{mod}} module, path=`pwd`"
  test -d {{mod}} || (echo "❌ Module not found" && exit 1)
  cd {{mod}} && dagger functions

# Recipe to run Daggy tests 🧪
daggy-tests:
  @echo "Running Daggy tests 🧪 ..."
  @cd .daggerx/daggy && cargo build --release
  @cd .daggerx/daggy && cargo test