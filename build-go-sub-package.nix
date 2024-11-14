{
  buildGoModule,
  subPackage,
  pname,
  version,
}:

buildGoModule {
  inherit pname version;
  src = ./.;
  subPackages = [ subPackage ];
  postInstall = ''
    mv "$out/bin/${subPackage}" "$out/bin/${pname}"
  '';
  vendorHash = "sha256-8eT95F+qWykInrc+s1HluoacrumGVFmoxbGELWadwSw=";
}
