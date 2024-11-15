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
  vendorHash = "sha256-qp2G3u95W58wQ8ikJDooST7wbqKpSkwakyfeSE00Wko=";
}
