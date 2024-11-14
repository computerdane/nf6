{
  dataDir ? "$HOME/.local/share/nf-dev",
  go,
  writeShellApplication,
}:

writeShellApplication {
  name = "dev-client-cli";
  runtimeInputs = [ go ];
  text = ''
    go run ./client-cli/*.go --dataDir "${dataDir}" "$@"
  '';
}
