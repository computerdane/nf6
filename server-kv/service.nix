{
  config,
  pkgs,
  lib,
  ...
}:

let
  cfg = config.nf6-server.kv;
in
{
  options.nf6-server.kv.enable = lib.mkEnableOption "Enable the nf6 key-value store";

  config = lib.mkIf cfg.enable {
    systemd.services.valkey.serviceConfig = {
      DynamicUser = true;
      ExecStart = "${pkgs.valkey}/bin/valkey";
    };
  };
}
