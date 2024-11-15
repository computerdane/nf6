{ callPackage, git }:

callPackage ../build-go-sub-package.nix {
  subPackage = "server-git-shell";
  pname = "nf6-git-shell";
  version = "0.0.1";
  buildInputs = [ git ];
}
