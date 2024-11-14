{
  baseDir ? "$HOME/.nf6/server-api",
  go,
  writeShellApplication,
}:

writeShellApplication {
  name = "dev-server-api";
  runtimeInputs = [ go ];
  text = ''
    go run ./server-api/*.go -base-dir="${baseDir}" "$@"
  '';
}
