{
  cfgFile ? "$HOME/.config/nf6-git-auth-dev/config.yaml",
  go,
  writeShellApplication,
}:

writeShellApplication {
  name = "dev-server-git-auth";
  runtimeInputs = [ go ];
  text = ''
    go run ./server-git-auth/*.go --config "${cfgFile}" "$@"
  '';
}
