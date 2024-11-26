{
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
  outputs =
    { nixpkgs, ... }:
    let
      cfg = builtins.fromJSON (builtins.readFile ./config.json);
    in
    {
      nixosConfigurations.nf6 = nixpkgs.lib.nixosSystem {
        system = cfg.System;
        modules = [ ];
      };
    };
}
