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
        sql-scripts = {
          init-api-user-sql = ./server-db/init-api-user.sql;
          init-git-user-sql = ./server-db/init-git-user.sql;
          init-tables-sql = ./server-db/init-tables.sql;
        };
        pkgs-nf6 = {
          client-cli = pkgs.callPackage ./client-cli/default.nix { };
          server-api = pkgs.callPackage ./server-api/default.nix { };
          server-git-auth = pkgs.callPackage ./server-git-auth/default.nix { };
          server-git-shell = pkgs.callPackage ./server-git-shell/default.nix { };
        } // sql-scripts;
      in
      {
        devShell = pkgs.mkShell {
          buildInputs =
            (with pkgs; [
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
            ])
            ++ pkgs.lib.flatten [
              (pkgs.callPackage ./client-cli/develop.nix { })
              (pkgs.callPackage ./server-api/develop.nix { })
              (pkgs.callPackage ./server-git-auth/develop.nix { })
              (pkgs.callPackage ./server-git-shell/develop.nix { })
              (pkgs.callPackage ./server-router/develop.nix { })

              (pkgs.callPackage ./server-db/develop.nix { inherit sql-scripts; })
            ];
        };
        packages = pkgs-nf6 // {
          nixosModules = import ./nixos-modules/all-modules.nix { inherit pkgs-nf6; };
        };
      }
    );
}
