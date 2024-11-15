{ pkgs-nf6 }:

{
  config,
  lib,
  pkgs,
  ...
}:

let
  cfg = config.services.nf6-api;
in
{
  options.services.nf6-api =
    with lib;
    with types;
    {
      enable = mkEnableOption "nf6 API server";
      settings = mkOption {
        description = "maps directly to YAML configuration";
        type = attrs;
        default = { };
      };
      user = mkOption {
        description = "run as specified user";
        type = str;
        default = "nf6-api";
      };
      group = mkOption {
        description = "run as specified group";
        type = str;
        default = "nf6-api";
      };
      postgresPasswordFile = mkOption {
        description = "path to file containing postgres password";
        type = str;
      };
      postgresHost = mkOption {
        description = "postgres host";
        type = str;
        default = "localhost";
      };
      openFirewall = mkOption {
        description = "whether to open ports in firewall";
        type = bool;
        default = false;
      };
    };

  config =
    let
      configYaml = pkgs.writeText "config.yaml" (builtins.toJSON cfg.settings);
      dataDir = if (cfg.settings ? dataDir) then cfg.settings.dataDir else "/var/lib/nf6-api/data";
      ports = [
        (if (cfg.settings ? portInsecure) then cfg.portInsecure else 6968)
        (if (cfg.settings ? portSecure) then cfg.portSecure else 6969)
      ];
    in
    lib.mkIf cfg.enable {
      networking.firewall = lib.mkIf cfg.openFirewall {
        allowedTCPPorts = ports;
        allowedUDPPorts = ports;
      };

      users.groups.${cfg.group} = { };
      users.users.${cfg.user} = {
        isNormalUser = true;
        group = cfg.group;
      };

      systemd.services.nf6-api = {
        wantedBy = [ "multi-user.target" ];
        path = [ pkgs-nf6.server-api ];
        script = ''
          PG_PASS=$(cat "${cfg.postgresPasswordFile}")
          nf6-api --config "${configYaml}" \
            --dbUrl "postgres://nf6_api:$PG_PASS@${cfg.postgresHost}/nf6"
        '';
        serviceConfig = {
          User = cfg.user;
          Group = cfg.group;
          PrivateTmp = true;
          RemoveIPC = true;
        };
      };

      systemd.tmpfiles.settings."10-nf6-api".${dataDir}.d = {
        user = cfg.user;
        group = cfg.group;
        mode = "0755";
      };
    };
}
