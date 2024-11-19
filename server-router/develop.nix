{
  cfgFile ? "$HOME/.config/nf6-router/config.yaml",
  go,
  writeShellApplication,
}:

writeShellApplication {
  name = "dev-server-router";
  runtimeInputs = [ go ];
  text = ''
    go run ./server-router/*.go --config "${cfgFile}" "$@"
  '';
}
