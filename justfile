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
dc mod args:
  @echo "Running Dagger module..."
  @echo "Currently in {{mod}} module 📦"
  @test -d {{mod}}/dagger || (echo "Module not found" && exit 1)
  @cd {{mod}}/dagger && dagger call {{args}}

# Recipe to run Dagger module tests. It requires the module name and extra arguments.
dct mod:
  @echo "Running Dagger module tests..."
  @echo "Currently in {{mod}} module 🧪
  @test -d {{mod}}/dagger || (echo "Module not found" && exit 1)
  @cd {{mod}}/tests/dagger && dagger call {{args}}
