#!/usr/bin/env bash
# Verify that the local Nix C API headers match what the Go bindings expect.
# Run inside: nix develop
set -euo pipefail

need=(
  nix-util-c
  nix-store-c
  nix-expr-c
  nix-fetchers-c
  nix-flake-c
  nix-main-c
)

echo "== pkg-config packages =="
missing=0
for p in "${need[@]}"; do
  if pkg-config --exists "$p"; then
    echo "  OK  $p $(pkg-config --modversion "$p" 2>/dev/null || true)"
  else
    echo "  MISSING  $p"
    missing=1
  fi
done
if [ "$missing" -ne 0 ]; then
  echo "error: one or more nix-*-c packages not found; check PKG_CONFIG_PATH" >&2
  exit 1
fi

echo "== required symbols / headers =="
cflags=$(pkg-config --cflags nix-store-c)
# shellcheck disable=SC2086
hdr=$(echo $cflags | tr ' ' '\n' | sed -n 's/^-I//p' | head -1)
if [ -z "$hdr" ]; then
  echo "warn: could not resolve -I path from nix-store-c"
else
  for f in nix_api_util.h nix_api_store.h nix_api_expr.h nix_api_value.h; do
    if [ -f "$hdr/$f" ] || find "$hdr" -name "$f" 2>/dev/null | grep -q .; then
      echo "  OK  $f"
    else
      echo "  MISSING header $f under $hdr"
      missing=1
    fi
  done
fi

echo "== go build =="
export CGO_ENABLED=1
go build ./...
echo "  OK  go build ./..."

echo "== done =="
exit "$missing"
