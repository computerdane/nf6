{ pkgs-nf6 }:

[
  (import ./server-api.nix { inherit pkgs-nf6; })
  (import ./server-db.nix { inherit pkgs-nf6; })
  (import ./server-git.nix { inherit pkgs-nf6; })
]
