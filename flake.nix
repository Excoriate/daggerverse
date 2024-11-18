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
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.treefmt-nix.flakeModule
        inputs.flake-root.flakeModule
      ];

      systems = ["x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin"];

      perSystem = { config, self', inputs', pkgs, system, ... }:
      let
        pkgs = import nixpkgs {
          inherit system;
          config = {
            allowUnfree = true;
          };
        };
      in
      {
        # Define development shell and other configurations
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            treefmt
            alejandra
            rustfmt
            nodePackages.prettier
            go
            terraform
            terragrunt
          ];
        };

        # Define the treefmt app
        apps.treefmt = {
          type = "app";
          program = "${pkgs.treefmt}/bin/treefmt";
        };

        treefmt.config = {
          projectRootFile = "flake.nix";
          programs = {
            alejandra.enable = true;
            gofmt.enable = true;
            prettier.enable = true;
            rustfmt.enable = true;
            terraform.enable = true;
            terragrunt.enable = true;
          };
          settings.formatter = {
            nix = {
              command = "alejandra";
              includes = ["*.nix"];
            };
            go = {
              command = "gofmt";
              includes = ["*.go"];
              excludes = [
                "*/internal/*"
                "dagger.gen.go"
              ];
            };
            prettier = {
              command = "prettier";
              options = ["--write"];
              includes = [
                "*.js"
                "*.jsx"
                "*.ts"
                "*.tsx"
                "*.json"
                "*.md"
                "*.markdown"
                "*.yaml"
                "*.yml"
              ];
            };
            rust = {
              command = "rustfmt";
              includes = ["*.rs"];
            };
            terraform = {
              command = "terraform";
              options = ["fmt"];
              includes = ["*.tf"];
            };
            terragrunt = {
              command = "terragrunt";
              options = ["hclfmt"];
              includes = ["*.hcl"];
            };
          };
        };
      };
    };
}
