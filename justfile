export NIXPKGS_ALLOW_UNFREE := "1"
export NOTHANKS := "1"

default:
  @just --list

# Recipe to run your development environment commands
dev:
  @echo "Entering Nix development environment 🧰 ..."
  @nix develop --impure --extra-experimental-features nix-command --extra-experimental-features flakes

# Recipe to initialize the project
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

# Recipe to run pre-commit hooks
precommit:
  @echo "Running pre-commit hooks 🔍 ..."
  @nix-shell -p pre-commit --run "pre-commit run --all-files"
  @echo "Pre-commit hooks passed ✅"

# Recipe to run Dagger module. It requires the module name and extra arguments.
dc mod *args:
  @echo "Running Dagger module..."
  @echo "Currently in {{mod}} module 📦, path=`pwd`"
  @test -d {{mod}} || (echo "Module not found" && exit 1)
  @cd {{mod}} && dagger call {{args}}

# Recipe to run Dagger module tests. It requires the module name and extra arguments.
dct mod *args:
  @echo "Running Dagger module tests..."
  @echo "Currently in {{mod}} module 🧪, path=`pwd`"
  @test -d {{mod}}/tests || (echo "Module not found" && exit 1)
  @cd {{mod}}/tests && dagger call {{args}}

# Recipe to run Dagger module examples. It requires the module name and extra arguments.
dce mod *args:
  @echo "Running Dagger module examples ... 📄"
  @echo "Currently in {{mod}} module 🧪, path=`pwd`"
  @test -d {{mod}}/examples/go || (echo "Module examples not found" && exit 1)
  @cd {{mod}}/examples/go && dagger call {{args}}

# Recipe to reload Dagger module (Dagger Develop)
reloadmod mod:
  @echo "Running Dagger development in a given module..."
  @echo "Currently in {{mod}} module 📦, path=`pwd`"
  @test -d {{mod}} || (echo "Module not found" && exit 1)
  @cd {{mod}} && dagger develop
  @echo "Module reloaded successfully ✅"

# Recipe to reload Dagger module and its underlying tests (Dagger Develop & Dagger Call/Functions)
reloadall mod:
  @echo "Reloading Dagger module and also the tests..."
  @echo "Currently in {{mod}} module 🔄, path=`pwd`"
  @test -d {{mod}} || (echo "Module not found" && exit 1)
  @cd {{mod}} && dagger develop
  @cd {{mod}}/tests && dagger develop
  @cd {{mod}}/examples/go && dagger develop
  @echo "Module reloaded successfully 🚀"
  @echo "Inspecting the module... 🕵️"
  @cd {{mod}} && dagger call && dagger functions

# Recipe to run all the tests in the target module
test mod: (reloadmod mod)
  @echo "Running Dagger module tests..."
  @echo "Currently in {{mod}} module 🧪, path=`pwd`"
  @test -d {{mod}}/tests || (echo "Module not found" && exit 1)
  @cd {{mod}}/tests && dagger call test-all

# Recipe to run all the examples in the target module
examplesgo mod: (reloadmod mod)
  @echo "Running Dagger module examples (Go SDK)..."
  @echo "Currently in {{mod}} module 🧪, path=`pwd`"
  @test -d {{mod}}/examples/go || (echo "Module examples not found" && exit 1)
  @cd {{mod}}/examples/go && dagger call create-container
  @cd {{mod}}/examples/go && dagger call run-arbitrary-command
  @cd {{mod}}/examples/go && dagger call passed-env-vars
  @cd {{mod}}/examples/go && dagger call create-net-rc-file-for-github

# Recipe to run GolangCI Lint
golint mod:
  @echo "Running Go (GolangCI)... 🧹 "
  @test -d {{mod}} || (echo "Module not found" && exit 1)
  @echo "Currently in {{mod}} module 📦, path=`pwd`/{{mod}}"
  @nix-shell -p golangci-lint --run "golangci-lint run --config .golangci.yml ./{{mod}}"
  @echo "Checking now the tests 🧪 project ..."
  @nix-shell -p golangci-lint --run "golangci-lint run --config .golangci.yml ./{{mod}}/tests"
  @echo "Checking now the examples 📄 project ..."
  @nix-shell -p golangci-lint --run "golangci-lint run --config .golangci.yml ./{{mod}}/examples/go"

# Recipe to run the whole CI locally
cilocal mod: (reloadall mod) (golint mod) (test mod) (examplesgo mod) (ci-module-docs mod)
  @echo "Running the whole CI locally... 🚀"

# Recipe to create a new module using Daggy (a rust CLI tool)
create mod:
  @echo "Creating a new module..."
  @cd .daggerx/daggy && cargo build --release
  @.daggerx/daggy/target/release/daggy --task=create --module={{mod}}

# This recipe validate if the dagger module has the README.md file and the LICENSE file
ci-module-docs mod:
  @echo "Validating the module documentation..."
  @test -f {{mod}}/README.md || (echo "README.md file not found" && exit 1)
  @test -f {{mod}}/LICENSE || (echo "LICENSE file not found" && exit 1)
  @echo "Module documentation is valid ✅"

# recipe for dagger call
call mod *args:
  @echo "Running Dagger call..."
  @echo "Currently in {{mod}} module 📦, path=`pwd`"
  @test -d {{mod}} || (echo "Module not found" && exit 1)
  @cd {{mod}} && dagger call {{args}}

# Recipe for dagger call tests in a certain module, E.g.: just calltests modexample my-function
calltests mod *args:
  @echo "Running Dagger call tests..."
  @echo "Currently in {{mod}} module 🧪, path=`pwd`"
  @test -d {{mod}}/tests || (echo "Module not found" && exit 1)
  @cd {{mod}}/tests && dagger call {{args}}

# Recipe to run dagger develop in all modules
develop-all:
  @echo "Developing all Dagger modules..."
  @cd .daggerx/daggy && cargo build --release
  @.daggerx/daggy/target/release/daggy --task=develop
