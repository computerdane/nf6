{ callPackage }:

callPackage ../build-go-sub-package.nix {
  subPackage = "server-api";
  pname = "nf-api";
  version = "0.0.1";
}
