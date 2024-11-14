{
  buildGoModule,
  buildInputs ? [ ],
  installShellFiles,
  pname,
  subPackage,
  version,
}:

buildGoModule {
  inherit pname version buildInputs;
  src = ./.;
  subPackages = [ subPackage ];
  nativeBuildInputs = [ installShellFiles ];
  postInstall = ''
    mv "$out/bin/${subPackage}" "$out/bin/${pname}"
    installShellCompletion --cmd ${pname} \
      --bash <($out/bin/${pname} completion bash) \
      --fish <($out/bin/${pname} completion fish) \
      --zsh <($out/bin/${pname} completion zsh)
  '';
  vendorHash = "sha256-yW1i3deY6P3xYllGFdiy93bA6PtXuwPvCgMZXDInyjU=";
}
