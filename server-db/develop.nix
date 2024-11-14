{
  dataDir ? "$HOME/.nf6/server-db",
  postgresql,
  socketDir ? "/tmp",
  writeShellApplication,
  writeText,
}:

let
  initSql = writeText "init.sql" (builtins.readFile ./init.sql);
in
writeShellApplication {
  name = "dev-server-db";
  runtimeInputs = [ postgresql ];
  text = ''
    mkdir -p "${dataDir}"
    chmod 700 "${dataDir}"

    initdb -D "${dataDir}" || true
    postgres -D "${dataDir}" -k "${socketDir}"

    createdb -h "${socketDir}" nf6
    psql -h ${socketDir} -d nf6 -f "${initSql}"
  '';
}
