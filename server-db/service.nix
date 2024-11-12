{ config, lib, ... }:

let
  cfg = config.nf6-server.db;
in
{
  options.nf6-server.db.enable = lib.mkEnableOption "Enable the nf6 database";

  config = lib.mkIf cfg.enable {
    services.postgresql = {
      enable = true;
      ensureDatabases = [ "nf6" ];
      initialScript = ./init.sql;
    };
  };
}
