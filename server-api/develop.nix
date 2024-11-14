{
  baseDir ? "$HOME/.local/share/nfapi-dev",
  go,
  writeShellApplication,
}:

writeShellApplication {
  name = "dev-server-api";
  runtimeInputs = [ go ];
  text = ''
    go run ./server-api/*.go --dataDir "${baseDir}" "$@"
  '';
}
