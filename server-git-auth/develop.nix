{
  cfgFile ? "$HOME/.config/nf6-git-auth-dev/config.yaml",
  dataDir ? "$HOME/.local/share/nf6-git-auth-dev",
  go,
  writeShellApplication,
}:

writeShellApplication {
  name = "dev-server-git-auth";
  runtimeInputs = [ go ];
  text = ''
    go run ./server-git-auth/*.go --config "${cfgFile}" --dataDir "${dataDir}" "$@"
  '';
}
