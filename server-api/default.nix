{ callPackage }:

callPackage ../build-go-sub-package.nix {
  subPackage = "server-api";
  pname = "nf6-api";
  version = "0.0.1";
}
