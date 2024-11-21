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
        vendorHash = "sha256-KK8Q3SMK3KSrCsHpdW4sxKNTSNGi0UW0YhsCaQcRxhc=";
      in
      {
        devShell = pkgs.mkShell {
          IS_DEV_SHELL = "1";

          buildInputs = with pkgs; [
            buf-language-server
            git
            go
            gopls
            openssh
            openssl
            postgresql
            protobuf
            protoc-gen-go
            protoc-gen-go-grpc
          ];
        };
        packages = {
          default =
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
                  --prefix PATH : ${with pkgs; lib.makeBinPath [ openssh ]}
                installShellCompletion --cmd nf \
                  --bash <($out/bin/nf completion bash) \
                  --fish <($out/bin/nf completion fish) \
                  --zsh  <($out/bin/nf completion zsh)
              '';
            };
        };
      }
    );
}
