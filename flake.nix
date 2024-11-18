{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    treefmt-nix.url = "github:numtide/treefmt-nix";
    rust-overlay.url = "github:oxalica/rust-overlay";
    flake-root.url = "github:srid/flake-root";
    nix-direnv.url = "github:nix-community/nix-direnv";
    pre-commit-hooks.url = "github:cachix/pre-commit-hooks.nix";
  };

  outputs = inputs @ {
    flake-parts,
    nixpkgs,
    rust-overlay,
    pre-commit-hooks,
    ...
  }:
    flake-parts.lib.mkFlake {inherit inputs;} {
      imports = [
        inputs.treefmt-nix.flakeModule
        inputs.flake-root.flakeModule
      ];

      systems = ["x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin"];

      perSystem = {
        config,
        self',
        inputs',
        pkgs,
        system,
        ...
      }: let
        overlays = [
          (import rust-overlay)
        ];
        pkgs = import nixpkgs {
          inherit system overlays;
          config.allowUnfree = true;
        };

        # Pre-commit hooks configuration
        pre-commit-check = pre-commit-hooks.lib.${system}.run {
          src = ./.;
          # Add excludes to skip unnecessary files/dirs
          excludes = [
            "^(.direnv/.*)"
            "^(target/.*)"
            "^(dist/.*)"
            "^(result.*)"
            "^(.git/.*)"
          ];

          hooks = {
            treefmt = {
              enable = true;
              # Only run on changed files
              pass_filenames = true;
            };

            rustfmt = {
              enable = true;
              entry = "${pkgs.rustfmt}/bin/rustfmt";
              types = ["rust"];
              # Only run on changed files
              pass_filenames = true;
            };

            golangci-lint = {
              enable = true;
              entry = "${pkgs.golangci-lint}/bin/golangci-lint run";
              types = ["go"];
              # Only run on changed files
              pass_filenames = true;
              # Add performance optimizations
              args = [
                "--fast"
                "--max-same-issues=0"
                "--max-issues-per-linter=0"
                "--modules-download-mode=readonly"
              ];
            };

            yamllint = {
              enable = true;
              entry = "${pkgs.yamllint}/bin/yamllint";
              types = ["yaml"];
              # Only run on changed files
              pass_filenames = true;
            };
          };
        };
      in {
        devShells.default = pkgs.mkShell {
          name = "dev-environment";
          shell = "${pkgs.bash}/bin/bash --noprofile --norc";

          packages = with pkgs; [
            # Direnv
            nix-direnv

            # Rust
            rust-bin.stable.latest.default
            cargo
            rustc
            rustfmt
            clippy

            # Go tools
            go
            golangci-lint

            # Terraform and Terragrunt
            terraform
            terragrunt

            # Required tools
            just
            git
            semver-tool
            jq
            yq-go
            moreutils
            yamllint
            yamlfmt
            pre-commit

            # Add treefmt here
            treefmt
          ];

          shellHook = ''
            export RUST_SRC_PATH=${pkgs.rust.packages.stable.rustPlatform.rustLibSrc}
            export GOROOT=${pkgs.go}/share/go

            # Cache Go modules
            export GOMODCACHE="$PWD/.direnv/go/pkg/mod"
            export GOCACHE="$PWD/.direnv/go/cache"

            # Only run pre-commit on changed files
            export PRE_COMMIT_FILES_CHANGED=$(git diff --name-only HEAD)

            # Run pre-commit checks only if there are changes
            if [ -n "$PRE_COMMIT_FILES_CHANGED" ]; then
              ${pre-commit-check.shellHook}
            fi

            echo "🌟 Welcome to the Daggerverse development environment! 🚀"
            echo "Happy coding! 💻"
          '';
        };

        # Add pre-commit check to checks
        checks.pre-commit-check = pre-commit-check;

        treefmt = {
          projectRootFile = config.flake-root.projectRootFile;
          programs = {
            alejandra.enable = true;
            rustfmt.enable = true;
            prettier.enable = true;
            gofmt.enable = true;
            terraform.enable = true;
          };
        };
      };
    };
}
