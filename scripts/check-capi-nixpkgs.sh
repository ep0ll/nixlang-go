#!/usr/bin/env bash
# Build/test against nixpkgs' packaged Nix C API (secondary validation).
# Primary pin remains flake.lock — see COMPAT.md.
set -euo pipefail

export CGO_ENABLED=1

NIX_DEV=$(nix build --no-link --print-out-paths 'nixpkgs#nixVersions.latest.dev')
export PKG_CONFIG_PATH="${NIX_DEV}/lib/pkgconfig${PKG_CONFIG_PATH:+:$PKG_CONFIG_PATH}"

# Pull runtime libs for linking
NIX_OUT=$(nix build --no-link --print-out-paths 'nixpkgs#nixVersions.latest')
export LD_LIBRARY_PATH="${NIX_OUT}/lib${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}"

echo "Nix: $(${NIX_OUT}/bin/nix --version 2>/dev/null || nix --version || true)"
echo "PKG_CONFIG_PATH=$PKG_CONFIG_PATH"
./scripts/check-capi.sh
go build -tags=nix ./...
go test -tags=nix ./util ./value ./expr ./store ./flake -count=1
echo "OK (nixpkgs C API)"
