{ callPackage }:

callPackage ../build-go-sub-package.nix {
  subPackage = "server-router";
  pname = "nf6-router";
  version = "0.0.1";
}
