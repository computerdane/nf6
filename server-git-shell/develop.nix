{ go, writeShellApplication }:

writeShellApplication {
  name = "dev-server-git-shell";
  runtimeInputs = [ go ];
  text = ''
    go run ./server-git-shell/*.go "$@"
  '';
}
