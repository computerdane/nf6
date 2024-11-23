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
        modules = [
          (
            {
              lib,
              modulesPath,
              pkgs,
              ...
            }:
            {
              imports = [ (modulesPath + "/installer/cd-dvd/installation-cd-minimal.nix") ];

              services.openssh = {
                enable = true;
                settings = {
                  PermitRootLogin = lib.mkForce "no";
                  PasswordAuthentication = false;
                  KbdInteractiveAuthentication = false;
                };
              };

              users.users.setup = {
                isNormalUser = true;
                extraGroups = [ "wheel" ];
                openssh.authorizedKeys.keys = [ cfg.SshPubKey ];
              };

              networking = {
                useNetworkd = true;
                nftables.enable = true;
                nameservers = [
                  "1.1.1.1"
                  "1.0.0.1"
                ];
              };

              systemd.network.networks."10-wgnf6" = {
                name = "wgnf6";
                networkConfig.Address = cfg.HostAddr;
                routes = [
                  {
                    routeConfig = {
                      PreferredSource = cfg.HostAddr;
                      Destination = cfg.ServerGlobalPrefix6;
                    };
                  }
                ];
              };

              systemd.network.netdevs."10-wgnf6" = {
                netdevConfig = {
                  Kind = "wireguard";
                  Name = "wgnf6";
                };
                wireguardConfig.PrivateKeyFile = pkgs.writeText "wgnf6.key" cfg.WgPrivKey;
                wireguardPeers = [
                  {
                    wireguardPeerConfig = {
                      Endpoint = cfg.WgServerEndpoint;
                      AllowedIPs = [ cfg.HostAddr ];
                      PublicKey = cfg.WgServerWgPubKey;
                    };
                  }
                ];
              };

              nix.settings = {
                experimental-features = [
                  "nix-command"
                  "flakes"
                ];
                auto-optimise-store = true;
              };

              system.stateVersion = "24.05";
            }
          )
        ];
      };
    };
}
