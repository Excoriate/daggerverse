{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    treefmt-nix.url = "github:numtide/treefmt-nix";
    rust-overlay.url = "github:oxalica/rust-overlay";
  };

  outputs = inputs@{ flake-parts, nixpkgs, rust-overlay, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.treefmt-nix.flakeModule
      ];

      systems = [ "x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin" ];

      perSystem = { config, self', inputs', pkgs, system, ... }:
        let
          overlays = [
            (import rust-overlay)
          ];
          pkgs = import nixpkgs {
            inherit system overlays;
            config.allowUnfree = true;
          };
        in
        {
          devShells.default = pkgs.mkShell {
            name = "dev-environment";

            packages = with pkgs; [
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
            ];

            shellHook = ''
              export RUST_SRC_PATH=${pkgs.rust.packages.stable.rustPlatform.rustLibSrc}
              export GOROOT=${pkgs.go}/share/go
              echo "🌟 Welcome to the Daggerverse development environment! 🚀"
              echo "Happy coding! 💻"
            '';
          };

          treefmt.config = {
            inherit (config.flake-root) projectRootFile;
            package = pkgs.treefmt;

            programs = {
              alejandra.enable = true;
              rustfmt.enable = true;
              prettier = {
                enable = true;
                # Configure Prettier to format JSON and YAML files
                include = [ "**/*.{js,jsx,ts,tsx,json,yml,yaml}" ];
              };
              gofmt.enable = true;
              terraform-fmt = {
                enable = true;
                command = "${pkgs.terraform}/bin/terraform fmt -";
                include = [ "**/*.tf" ];
              };
              terragrunt-fmt = {
                enable = true;
                command = "${pkgs.terragrunt}/bin/terragrunt hclfmt";
                include = [ "**/*.hcl" ];
              };
            };
          };

          formatter = config.treefmt.build.wrapper;
        };
    };
}
