{
  description = "Production Go bindings for the Nix C API (master)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    # Track Nix master for the experimental C API. Run `nix flake lock` and
    # commit flake.lock so builds pin a known revision.
    nix.url = "github:NixOS/nix";
    nix.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, nix }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
    in
    {
      devShells = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          nixPkg =
            if nix.packages ? ${system} && nix.packages.${system} ? nix
            then nix.packages.${system}.nix
            else if nix.packages ? ${system} && nix.packages.${system} ? default
            then nix.packages.${system}.default
            else pkgs.nixVersions.latest;
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              go
              pkg-config
              gcc
            ];
            buildInputs = [
              nixPkg
            ] ++ pkgs.lib.optionals (nixPkg ? dev) [ nixPkg.dev ];

            shellHook = ''
              export CGO_ENABLED=1
              if [ -n "${nixPkg}" ]; then
                for d in "${nixPkg}/lib/pkgconfig" "${nixPkg.dev or nixPkg}/lib/pkgconfig"; do
                  if [ -d "$d" ]; then
                    export PKG_CONFIG_PATH="$d''${PKG_CONFIG_PATH:+:$PKG_CONFIG_PATH}"
                  fi
                done
              fi
              echo "nixlang-go dev shell"
              echo "  nix: $(${nixPkg}/bin/nix --version 2>/dev/null || echo unknown)"
              echo "  pkg-config nix-util-c: $(pkg-config --modversion nix-util-c 2>/dev/null || echo NOT FOUND)"
              echo "  Build: go build ./..."
              echo "  Test:  go test -tags=nix ./..."
            '';
          };
        });
    };
}
