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
            { lib, modulesPath, ... }:
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

              users.users.install = {
                isNormalUser = true;
                extraGroups = [ "wheel" ];
                openssh.authorizedKeys.keys = [ cfg.SshPubKey ];
              };

              networking.useNetworkd = true;

              networking.wg-quick.interfaces.wgtest = {
                address = [ cfg.HostAddr6 ];
                dns = [ "2606:4700:4700::1111" ];
                privateKey = cfg.WgPrivKey;
                peers = [
                  {
                    allowedIPs = [ "::/0" ];
                    endpoint = cfg.WgServerEndpoint;
                    persistentKeepalive = 25;
                    publicKey = cfg.WgServerWgPubKey;
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
