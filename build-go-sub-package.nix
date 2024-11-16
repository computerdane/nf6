{
  buildGoModule,
  buildInputs ? [ ],
  buildPackages,
  installShellFiles,
  pname,
  stdenv,
  subPackage,
  version,
}:

buildGoModule {
  inherit pname version buildInputs;
  src = ./.;
  subPackages = [ subPackage ];
  nativeBuildInputs = [ installShellFiles ];
  doCheck = false;
  postInstall =
    let
      emulator = stdenv.hostPlatform.emulator buildPackages;
    in
    ''
      mv $out/bin/${subPackage} $out/bin/${pname}
      installShellCompletion --cmd ${pname} \
        --bash <(${emulator} $out/bin/${pname} completion bash) \
        --fish <(${emulator} $out/bin/${pname} completion fish) \
        --zsh  <(${emulator} $out/bin/${pname} completion zsh)
    '';
  vendorHash = "sha256-qp2G3u95W58wQ8ikJDooST7wbqKpSkwakyfeSE00Wko=";
}
