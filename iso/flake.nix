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
              config,
              lib,
              modulesPath,
              ...
            }:
            # let
            #   wgPrivKey = pkgs.writeText "wgnf6.key" cfg.WgPrivKey;
            # in
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
                # useNetworkd = true;
                # nftables.enable = true;
                nameservers = [
                  "1.1.1.1"
                  "1.0.0.1"
                ];
                wg-quick.interfaces.wgnf6 = {
                  dns = config.networking.nameservers;
                  privateKey = cfg.WgPrivKey;
                  peers = [
                    {
                      endpoint = cfg.WgServerEndpoint;
                      allowedIPs = [
                        "::/0"
                        "0.0.0.0/0"
                      ];
                      publicKey = cfg.WgServerWgPubKey;
                    }
                  ];
                };
              };

              # systemd.network.networks."10-wgnf6" = {
              #   name = "wgnf6";
              #   networkConfig.Address = cfg.HostAddr;
              #   routes = [
              #     {
              #       routeConfig = {
              #         PreferredSource = cfg.HostAddr;
              #         Destination = cfg.ServerGlobalPrefix6;
              #       };
              #     }
              #   ];
              # };

              # systemd.network.netdevs."10-wgnf6" = {
              #   netdevConfig = {
              #     Kind = "wireguard";
              #     Name = "wgnf6";
              #   };
              #   wireguardConfig.PrivateKeyFile = pkgs.writeText "wgnf6.key" cfg.WgPrivKey;
              #   wireguardPeers = [
              #     {
              #       wireguardPeerConfig = {
              #         Endpoint = cfg.WgServerEndpoint;
              #         AllowedIPs = [
              #           "::/0"
              #           "0.0.0.0/0"
              #         ];
              #         PublicKey = cfg.WgServerWgPubKey;
              #       };
              #     }
              #   ];
              # };

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
