{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
    utils.url = "github:numtide/flake-utils";
  };
  outputs =
    { nixpkgs, utils, ... }:
    utils.lib.eachDefaultSystem (
      system:
      let

        pkgs = import nixpkgs { inherit system; };
        vendorHash = "sha256-aK8VAY628aqy9L75LQg+M6YtbCuqF5P1rGjhxfXb8kE=";
        # vendorHash = pkgs.lib.fakeHash;

      in
      rec {
        devShell = pkgs.mkShell {
          IS_DEV_SHELL = "1";

          buildInputs = with pkgs; [
            buf-language-server
            git
            go
            gopls
            nix
            openssh
            openssl
            postgresql
            protobuf
            protoc-gen-go
            protoc-gen-go-grpc
          ];
        };

        nixosModules = import ./modules/server.nix { pkgs-nf6 = packages; };

        packages = rec {
          default = nf;
          nf =
            let
              version = "0.1";
              base = pkgs.buildGoModule {
                inherit version vendorHash;
                pname = "nf6-cli";
                src = ./.;
                subPackages = [ "cli" ];
              };
            in
            pkgs.stdenv.mkDerivation {
              inherit version;
              pname = "nf";
              nativeBuildInputs = [
                base
                pkgs.makeBinaryWrapper
                pkgs.installShellFiles
              ];
              dontUnpack = true;
              installPhase = ''
                mkdir -p $out/bin
                makeBinaryWrapper ${base}/bin/cli $out/bin/nf \
                  --prefix PATH : ${
                    with pkgs;
                    lib.makeBinPath [
                      nix
                      openssh
                    ]
                  }
                installShellCompletion --cmd nf \
                  --bash <($out/bin/nf completion bash) \
                  --fish <($out/bin/nf completion fish) \
                  --zsh  <($out/bin/nf completion zsh)
              '';
            };
          nf6-api =
            let
              version = "0.1";
              base = pkgs.buildGoModule {
                inherit version vendorHash;
                pname = "nf6-api";
                src = ./.;
                subPackages = [ "api" ];
              };
            in
            pkgs.stdenv.mkDerivation {
              inherit version;
              pname = "nf6-api";
              nativeBuildInputs = [
                base
                pkgs.makeBinaryWrapper
                pkgs.installShellFiles
              ];
              dontUnpack = true;
              installPhase = ''
                mkdir -p $out/bin
                makeBinaryWrapper ${base}/bin/api $out/bin/nf6-api
                installShellCompletion --cmd nf6-api \
                  --bash <($out/bin/nf completion bash) \
                  --fish <($out/bin/nf completion fish) \
                  --zsh  <($out/bin/nf completion zsh)
              '';
            };
        };
      }
    );
}
