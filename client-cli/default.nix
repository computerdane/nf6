{
  callPackage,
  git,
  openssh,
}:

callPackage ../build-go-sub-package.nix {
  subPackage = "client-cli";
  pname = "nf";
  version = "0.0.1";
  buildInputs = [
    git
    openssh
  ];
}
