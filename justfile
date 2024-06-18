export NIXPKGS_ALLOW_UNFREE := "1"

default:
  @just --list

# Recipe to run your development environment commands
dev:
  @echo "Entering Nix development environment 🧰 ..."
  @nix develop --impure --extra-experimental-features nix-command --extra-experimental-features flakes

# Recipe to run pre-commit hooks
precommit:
  @echo "Running pre-commit hooks..."
  pre-commit run --all-files

# Recipe to find an specific nix package.
# Recipe to search for a specific nix package in fuzzy mode
nix-search pkg:
  @echo "Searching for packages in nixpkgs..."
  @nix search nixpkgs --json {{pkg}} --extra-experimental-features nix-command --extra-experimental-features flakes | jq -r 'to_entries[] | "\(.value.name) - \(.value.description)"'

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
  @nix-shell -p golangci-lint --run "golangci-lint run --config .golangci.yml ./{{mod}}/dagger" --verbose
  @echo "Checking now the tests 🧪 project ..."
  @nix-shell -p golangci-lint --run "golangci-lint run --config .golangci.yml ./{{mod}}/tests/dagger" --verbose

cilocal mod: (reloadall mod) (golint mod) (test mod)
  @echo "Running the whole CI locally... 🚀"
