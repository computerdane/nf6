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
      }
    );
}
