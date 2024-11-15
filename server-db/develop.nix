{
  dataDir ? "$HOME/.local/share/nf6-db-dev",
  postgresql,
  socketDir ? "/tmp",
  sql-scripts,
  writeShellApplication,
}:

with sql-scripts;

[
  (writeShellApplication {
    name = "dev-server-db";
    runtimeInputs = [ postgresql ];
    text = ''
      mkdir -p "${dataDir}"
      chmod 700 "${dataDir}"

      initdb -D "${dataDir}" || true
      postgres -D "${dataDir}" -k "${socketDir}"
    '';
  })
  (writeShellApplication {
    name = "dev-server-db-init";
    runtimeInputs = [ postgresql ];
    text = ''
      createdb -h "${socketDir}" nf6
      psql -h ${socketDir} -d nf6 -f "${init-tables-sql}"
      psql -h ${socketDir} -d nf6 -f "${init-api-user-sql}"
      psql -h ${socketDir} -d nf6 -f "${init-git-user-sql}"
    '';
  })
]
