{
  inputs,
  lib,
  system,
  ...
}:

let
  account = builtins.fromJSON ./account.json;
  global = builtins.fromJSON ./global.json;
  hosts = builtins.fromJSON ./hosts.json;
in
lib.mapAttrs (name: host: {
  inherit system;
  modules = [
    (
      { ... }:
      {
        services.openssh = {
          enable = true;
          settings = {
            PermitRootLogin = "no";
            PasswordAuthentication = false;
            KbdInteractiveAuthentication = false;
          };
        };

        users.users.root.openssh.authorizedKeys.keys = [ account.SshPubKey ];

        networking.useNetworkd = true;

        networking.wg-quick.interfaces.wgtest = {
          address = [ host.Addr6 ];
          dns = [
            "2606:4700:4700::1111"
            "2606:4700:4700::1001"
          ];
          privateKeyFile = "/run/nf6-secrets/wg.key";
          peers = [
            {
              allowedIPs = [ "::/0" ];
              endpoint = global.VipWgEndpoint;
              persistentKeepalive = 25;
              publicKey = global.VipWgPubKey;
            }
          ];
        };
      }
    )
  ];
}) hosts
