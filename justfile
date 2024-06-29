export NIXPKGS_ALLOW_UNFREE := "1"

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
  @test -d {{mod}}/dagger || (echo "Module not found" && exit 1)
  @cd {{mod}}/dagger && dagger call {{args}}

# Recipe to run Dagger module tests. It requires the module name and extra arguments.
dct mod *args:
  @echo "Running Dagger module tests..."
  @echo "Currently in {{mod}} module 🧪, path=`pwd`"
  @test -d {{mod}}/tests/dagger || (echo "Module not found" && exit 1)
  @cd {{mod}}/tests/dagger && dagger call {{args}}

# Recipe to reload Dagger module (Dagger Develop)
reloadmod mod:
  @echo "Running Dagger development in a given module..."
  @echo "Currently in {{mod}} module 📦, path=`pwd`"
  @test -d {{mod}}/dagger || (echo "Module not found" && exit 1)
  @cd {{mod}}/dagger && dagger develop
  @echo "Module reloaded successfully ✅"

# Recipe to reload Dagger module and its underlying tests (Dagger Develop & Dagger Call/Functions)
reloadall mod:
  @echo "Reloading Dagger module and also the tests..."
  @echo "Currently in {{mod}} module 🔄, path=`pwd`"
  @test -d {{mod}}/dagger || (echo "Module not found" && exit 1)
  @cd {{mod}}/dagger && dagger develop
  @cd {{mod}}/tests/dagger && dagger develop
  @echo "Module reloaded successfully 🚀"
  @echo "Inspecting the module... 🕵️"
  @cd {{mod}}/dagger && dagger call && dagger functions

# Recipe to run all the tests in the target module
test mod: (reloadmod mod)
  @echo "Running Dagger module tests..."
  @echo "Currently in {{mod}} module 🧪, path=`pwd`"
  @test -d {{mod}}/tests/dagger || (echo "Module not found" && exit 1)
  @cd {{mod}}/tests/dagger && dagger call test-all

# Recipe to run GolangCI Lint
golint mod:
  @echo "Running Go (GolangCI)... 🧹 "
  @test -d {{mod}}/dagger || (echo "Module not found" && exit 1)
  @echo "Currently in {{mod}} module 📦, path=`pwd`/{{mod}}/dagger"
  @nix-shell -p golangci-lint --run "golangci-lint run --config .golangci.yml ./{{mod}}/dagger"
  @echo "Checking now the tests 🧪 project ..."
  @nix-shell -p golangci-lint --run "golangci-lint run --config .golangci.yml ./{{mod}}/tests/dagger"

# Recipe to run the whole CI locally
cilocal mod: (reloadall mod) (golint mod) (test mod) (ci-module-docs mod)
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
  @test -d {{mod}}/dagger || (echo "Module not found" && exit 1)
  @cd {{mod}}/dagger && dagger call {{args}}

# Recipe for dagger call tests in a certain module, E.g.: just calltests modexample my-function
calltests mod *args:
  @echo "Running Dagger call tests..."
  @echo "Currently in {{mod}} module 🧪, path=`pwd`"
  @test -d {{mod}}/tests/dagger || (echo "Module not found" && exit 1)
  @cd {{mod}}/tests/dagger && dagger call {{args}}
