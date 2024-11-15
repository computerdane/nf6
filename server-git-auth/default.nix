{ callPackage }:

callPackage ../build-go-sub-package.nix {
  subPackage = "server-git-auth";
  pname = "nf6-git-auth";
  version = "0.0.1";
}
