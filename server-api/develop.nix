{
  cfgFile ? "$HOME/.config/nfapi-dev/config.yaml",
  baseDir ? "$HOME/.local/share/nfapi-dev",
  go,
  writeShellApplication,
}:

writeShellApplication {
  name = "dev-server-api";
  runtimeInputs = [ go ];
  text = ''
    go run ./server-api/*.go --config "${cfgFile}" --dataDir "${baseDir}" "$@"
  '';
}
