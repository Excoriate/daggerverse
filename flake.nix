{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    treefmt-nix.url = "github:numtide/treefmt-nix";
    rust-overlay.url = "github:oxalica/rust-overlay";

    # Add this for direnv support
    flake-root.url = "github:srid/flake-root";
    nix-direnv.url = "github:nix-community/nix-direnv";
  };

  outputs = inputs@{ flake-parts, nixpkgs, rust-overlay, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.treefmt-nix.flakeModule
        inputs.flake-root.flakeModule
      ];

      systems = [ "x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin" "x86_64-windows" ];

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
                include = [ "**/*.{js,jsx,ts,tsx,json}" ];
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
              yamllint = {
                enable = true;
                command = "${pkgs.yamllint}/bin/yamllint -c .yamllint.yml";
                include = [ "**/*.yaml" "**/*.yml" ];
              };
              yamlfmt = {
                enable = true;
                command = "${pkgs.yamlfmt}/bin/yamlfmt -w";
                include = [ "**/*.yaml" "**/*.yml" ];
              };
            };
          };

          formatter = config.treefmt.build.wrapper;
        };
    };
}
