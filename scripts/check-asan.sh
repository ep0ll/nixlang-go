#!/usr/bin/env bash
# Optional AddressSanitizer build of packages that link the Nix C API.
# Requires: clang or gcc with -fsanitize=address, and nix develop environment.
set -euo pipefail

export CGO_ENABLED=1
export CC="${CC:-gcc}"

# Prefer ASan when the toolchain supports it.
if ! echo 'int main(void){return 0;}' | "$CC" -fsanitize=address -x c - -o /tmp/nixgo-asan-probe 2>/dev/null; then
  echo "error: $CC does not support -fsanitize=address" >&2
  exit 1
fi
rm -f /tmp/nixgo-asan-probe

export CGO_CFLAGS="${CGO_CFLAGS:-} -fsanitize=address -fno-omit-frame-pointer"
export CGO_LDFLAGS="${CGO_LDFLAGS:-} -fsanitize=address"
export ASAN_OPTIONS="${ASAN_OPTIONS:-detect_leaks=0:halt_on_error=1}"

echo "== ASan go test -tags=nix (subset) =="
go test -tags=nix ./util ./value ./expr ./store ./flake -count=1
echo "OK"
