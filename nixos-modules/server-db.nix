{ pkgs-nf6 }:

{
  config,
  lib,
  pkgs,
  ...
}:

let
  cfg = config.services.nf6-db;
in
{
  options.services.nf6-db =
    with lib;
    with types;
    {
      enable = mkEnableOption "nf6 database";
      apiUserPasswordFile = mkOption {
        description = "path to file containing password for api user";
        type = str;
      };
      gitUserPasswordFile = mkOption {
        description = "path to file containing password for git user";
        type = str;
      };
    };

  config = lib.mkIf cfg.enable {
    services.postgresql = {
      enable = true;
      ensureDatabases = [ "nf6" ];
    };

    systemd.services.nf6-db-init = {
      requires = [ "postgresql.service" ];
      after = [ "postgresql.service" ];
      wantedBy = [ "multi-user.target" ];
      path = [ pkgs.postgresql ];
      preStart = ''
        sleep 5
      '';
      script = ''
        PG_NF6_API_PASS=$(cat "${cfg.apiUserPasswordFile}")
        PG_NF6_GIT_PASS=$(cat "${cfg.gitUserPasswordFile}")

        cat "${pkgs-nf6.init-tables-sql}" >> /tmp/init.sql
        cat "${pkgs-nf6.init-api-user-sql}" >> /tmp/init.sql
        cat "${pkgs-nf6.init-git-user-sql}" >> /tmp/init.sql

        sed -i -e "s/PG_NF6_API_PASS/$PG_NF6_API_PASS/g" /tmp/init.sql
        sed -i -e "s/PG_NF6_GIT_PASS/$PG_NF6_GIT_PASS/g" /tmp/init.sql

        psql -d nf6 -f /tmp/init.sql

        rm -f /tmp/init.sql
      '';
      serviceConfig = {
        User = "postgres";
        Group = "postgres";
        PrivateTmp = true;
      };
    };
  };
}
