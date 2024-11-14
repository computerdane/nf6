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
          buildInputs =
            (with pkgs; [
              buf-language-server
              go
              gopls
              openssl
              postgresql
              protobuf
              protoc-gen-go
              protoc-gen-go-grpc
            ])
            ++ [
              (pkgs.callPackage ./client-cli/develop.nix { })
              (pkgs.callPackage ./server-api/develop.nix { })
              (pkgs.callPackage ./server-db/develop.nix { })
            ];
        };
        packages = {
          client-cli = pkgs.callPackage ./client-cli/default.nix { };
          server-api = pkgs.callPackage ./server-api/default.nix { };

          init-sql = ./server-db/init.sql;
        };
      }
    );
}
