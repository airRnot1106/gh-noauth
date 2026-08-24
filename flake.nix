{
  description = "GitHub CLI tool without the auth command";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages.default = pkgs.buildGoModule (finalAttrs: {
          pname = "gh-noauth";
          version = self.shortRev or self.dirtyShortRev or "unknown";

          __structuredAttrs = true;

          src = ./.;

          vendorHash = "sha256-bVc4dhDapAp1YtO06C/nSrdxllpEhVFC2iZNPmjsJkI=";

          nativeBuildInputs = [
            pkgs.installShellFiles
            pkgs.makeWrapper
          ];

          # Using the Makefile (rather than buildGoModule's default `go build`)
          # keeps this in step with how upstream produces bin/gh. "nixpkgs" for
          # build.Date keeps the derivation reproducible instead of embedding
          # the actual build time.
          buildPhase = ''
            runHook preBuild
            make GO_LDFLAGS="-s -w -X github.com/cli/cli/v2/internal/build.Date=nixpkgs" GH_VERSION=${finalAttrs.version} bin/gh ${pkgs.lib.optionalString (pkgs.stdenv.buildPlatform.canExecute pkgs.stdenv.hostPlatform) "manpages"}
            runHook postBuild
          '';

          installPhase = ''
            runHook preInstall
            installBin bin/gh
            wrapProgram $out/bin/gh \
              --set-default GH_TELEMETRY false
          ''
          + pkgs.lib.optionalString (pkgs.stdenv.buildPlatform.canExecute pkgs.stdenv.hostPlatform) ''
            installManPage share/man/*/*.[1-9]

            installShellCompletion --cmd gh \
              --bash <($out/bin/gh completion -s bash) \
              --fish <($out/bin/gh completion -s fish) \
              --zsh <($out/bin/gh completion -s zsh)
          ''
          + ''
            runHook postInstall
          '';

          # The full suite includes tests that shell out to git and touch the
          # filesystem in ways the build sandbox does not allow; run `go test
          # ./...` separately (e.g. via `nix develop`) instead.
          doCheck = false;

          nativeInstallCheckInputs = [ pkgs.versionCheckHook ];
          doInstallCheck = true;

          meta = {
            description = "GitHub CLI tool without the auth command";
            license = pkgs.lib.licenses.mit;
            mainProgram = "gh";
          };
        });

        apps.default = flake-utils.lib.mkApp {
          drv = self.packages.${system}.default;
          name = "gh";
        };

        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.go
            pkgs.gopls
            pkgs.gotools
            pkgs.golangci-lint
            pkgs.git
          ];
        };
      }
    );
}
