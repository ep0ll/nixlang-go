{
  description = "Production Go bindings for the Nix C API";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
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
          nixPkg = nix.packages.${system}.nix;
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              go
              pkg-config
              gcc
            ];
            buildInputs = [
              nixPkg.dev
            ];
            shellHook = ''
              export CGO_ENABLED=1
              echo "nixgo dev shell: Nix C API from $(${nixPkg}/bin/nix --version 2>/dev/null || echo master)"
            '';
          };
        });
    };
}
