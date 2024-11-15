{
  cfgFile ? "$HOME/.config/nf6-dev/config.yaml",
  dataDir ? "$HOME/.local/share/nf6-dev",
  go,
  writeShellApplication,
}:

writeShellApplication {
  name = "dev-client-cli";
  runtimeInputs = [ go ];
  text = ''
    go run ./client-cli/*.go --config "${cfgFile}" --dataDir "${dataDir}" "$@"
  '';
}
