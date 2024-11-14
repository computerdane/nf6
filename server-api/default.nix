{ callPackage }:

callPackage ../build-go-sub-package.nix {
  subPackage = "server-api";
  pname = "nfapi";
  version = "0.0.1";
}
