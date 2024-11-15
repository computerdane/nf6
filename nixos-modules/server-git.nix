{ pkgs-nf6 }:

{
  config,
  lib,
  pkgs,
  ...
}:

let
  cfg = config.services.nf6-git;
in
{
  options.services.nf6-git =
    with lib;
    with types;
    {
      enable = mkEnableOption "nf6 git server";
      settings = mkOption {
        description = "maps directly to YAML configuration for nf6-git-auth";
        type = attrs;
        default = { };
      };
      user = mkOption {
        description = "run as specified user";
        type = str;
        default = "git";
      };
      group = mkOption {
        description = "run as specified group";
        type = str;
        default = "git";
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
    };

  config =
    let
      configYaml = pkgs.writeText "config.yaml" (builtins.toJSON cfg.settings);
      dataDir = if (cfg.settings ? dataDir) then cfg.settings.dataDir else "/var/lib/nf6-git-auth/data";
      gitReposDir =
        if (cfg.settings ? gitReposPath) then cfg.settings.gitReposPath else "/var/lib/nf6-git/repos";
    in
    lib.mkIf cfg.enable {
      users.groups.${cfg.group} = { };
      users.users.${cfg.user} = {
        isNormalUser = true;
        group = cfg.group;
        packages = [ pkgs.git ];
      };

      systemd.services.nf6-git-auth = {
        wantedBy = [ "multi-user.target" ];
        path = [ pkgs-nf6.server-git-auth ];
        script = ''
          PG_PASS=$(cat "${cfg.postgresPasswordFile}")
          nf6-git-auth listen --config "${configYaml}" \
            --dbUrl "postgres://nf6_git:$PG_PASS@${cfg.postgresHost}/nf6" \
            --gitShell "${pkgs-nf6.server-git-shell}/bin/nf6-git-shell"
        '';
        serviceConfig = {
          User = cfg.user;
          Group = cfg.group;
          PrivateTmp = true;
          RemoveIPC = true;
        };
      };

      systemd.tmpfiles.settings."10-nf6-git-auth".${dataDir}.d = {
        user = cfg.user;
        group = cfg.group;
        mode = "0755";
      };

      systemd.tmpfiles.settings."10-nf6-git-repos".${gitReposDir}.d = {
        user = cfg.user;
        group = cfg.group;
        mode = "0755";
      };

      services.openssh = {
        authorizedKeysCommand = ''/bin/nf6-git-auth ask %u "%t %k"'';
        authorizedKeysCommandUser = "nobody";
      };

      systemd.services.copy-nf6-git-auth-to-bin = {
        wantedBy = [ "multi-user.target" ];
        script = ''
          cp "${pkgs-nf6.server-git-auth}/bin/nf6-git-auth" /bin/nf6-git-auth
        '';
      };
    };
}
