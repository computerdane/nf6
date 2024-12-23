{ pkgs-nf6 }:

{
  config,
  lib,
  pkgs,
  ...
}:

let
  cfg = config.services.nf6;

  defaultSettings = {
    port = 6969;
    port-public = 6968;
    state-dir = "/var/lib/nf6-api/state";
  };

  defaultVipSettings = {
    state-dir = "/var/lib/nf6-vip/state";
  };

  initDbApiUserSql = pkgs.writeText "init-db-api-user.sql" ''
    create user nf6_api;

    grant usage on schema public to nf6_api;
    grant usage on all sequences in schema public to nf6_api;
    grant select, insert, update, delete on all tables in schema public to nf6_api;
  '';
in
{
  options.services.nf6 =
    with lib;
    with types;
    {
      enable = mkEnableOption "Nf6 API Server";
      settings = mkOption {
        description = "attrset mapping to YAML config for Nf6 API";
        type = attrs;
        default = { };
      };
      vipSettings = mkOption {
        description = "attrset mapping to YAML config for Nf6 VIP";
        type = attrs;
        default = { };
      };
      openFirewall = mkOption {
        description = "Whether or not to open firewall ports for the API server";
        type = bool;
        default = false;
      };
    };

  config =
    let
      settings = defaultSettings // cfg.settings;
      vipSettings = defaultVipSettings // cfg.vipSettings;
      configYaml = pkgs.writeText "config.yaml" (builtins.toJSON settings);
      vipConfigYaml = pkgs.writeText "config.yaml" (builtins.toJSON vipSettings);
    in
    lib.mkIf cfg.enable {
      services.postgresql = {
        enable = true;
        ensureDatabases = [ "nf6" ];
      };

      systemd.services.nf6-db-init = {
        requires = [ "postgresql.service" ];
        after = [ "postgresql.service" ];
        wantedBy = [ "multi-user.target" ];
        path = [ pkgs.postgresql ];
        script = ''
          sleep 5
          psql -d nf6 -f "${../db/init.sql}"
          psql -d nf6 -f "${initDbApiUserSql}"
        '';
        serviceConfig = {
          User = "postgres";
          Group = "postgres";
        };
      };

      networking.firewall =
        with settings;
        lib.mkIf cfg.openFirewall {
          allowedTCPPorts = [
            port
            port-public
          ];
          allowedUDPPorts = [
            port
            port-public
          ];
        };

      users.groups.nf6_api = { };
      users.users.nf6_api = {
        isNormalUser = true;
        group = "nf6_api";
      };

      systemd.services.nf6-api-public = {
        requires = [ "postgresql.service" ];
        after = [ "postgresql.service" ];
        wantedBy = [ "multi-user.target" ];
        path = [ pkgs-nf6.nf6-api ];
        script = ''
          sleep 5
          nf6-api serve-public --config-path "${configYaml}"
        '';
        serviceConfig = {
          User = "nf6_api";
          Group = "nf6_api";
          PrivateTmp = true;
        };
      };

      systemd.services.nf6-api = {
        requires = [ "postgresql.service" ];
        after = [ "postgresql.service" ];
        wantedBy = [ "multi-user.target" ];
        path = [ pkgs-nf6.nf6-api ];
        script = ''
          sleep 5
          nf6-api serve --config-path "${configYaml}"
        '';
        serviceConfig = {
          User = "nf6_api";
          Group = "nf6_api";
          PrivateTmp = true;
        };
      };

      systemd.tmpfiles.settings."10-nf6-api".${settings.state-dir}.d = {
        user = "nf6_api";
        group = "nf6_api";
        mode = "0755";
      };

      systemd.services.nf6-vip = {
        requires = [
          "postgresql.service"
          "nf6-api.service"
        ];
        after = [
          "postgresql.service"
          "nf6-api.service"
        ];
        wantedBy = [ "multi-user.target" ];
        path = [ pkgs-nf6.nf ];
        script = ''
          sleep 5
          nf vip --config-path "${vipConfigYaml}"
        '';
      };
    };
}
