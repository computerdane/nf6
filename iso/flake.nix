{
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
  outputs =
    { nixpkgs, ... }:
    let
      cfg = builtins.fromJSON (builtins.readFile ./config.json);
    in
    {
      nixosConfigurations.nf6 = nixpkgs.lib.nixosSystem {
        system = cfg.HostSystem;
        modules = [
          (
            { modulesPath, pkgs, ... }:
            {
              imports = [ (modulesPath + "/installer/cd-dvd/installation-cd-minimal.nix") ];

              services.openssh = {
                enable = true;
                settings = {
                  PermitRootLogin = "yes";
                  PasswordAuthentication = false;
                  KbdInteractiveAuthentication = false;
                };
              };

              users.users.root.openssh.authorizedKeys.keys = [ cfg.AccountSshPubKey ];

              networking.useNetworkd = true;
              networking.wg-quick.interfaces.wgnf6 = {
                address = [ cfg.HostAddr6 ];
                dns = [
                  "2606:4700:4700::1111"
                  "2606:4700:4700::1001"
                ];
                postUp = "${pkgs.iproute2}/bin/ip -6 rule add from ${cfg.HostAddr6} lookup 51820";
                postDown = "${pkgs.iproute2}/bin/ip -6 rule del from ${cfg.HostAddr6}";
                privateKey = cfg.HostWgPrivKey;
                peers = [
                  {
                    allowedIPs = [ "::/0" ];
                    endpoint = cfg.VipWgEndpoint;
                    persistentKeepalive = 25;
                    publicKey = cfg.VipWgPubKey;
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

              systemd.services.wg-quick-wgnf6.requires = [ "systemd-networkd-wait-online.service" ];
              systemd.services.wg-quick-wgnf6.after = [ "systemd-networkd-wait-online.service" ];

              system.stateVersion = "24.05";
            }
          )
        ];
      };
    };
}
