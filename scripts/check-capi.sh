#!/usr/bin/env bash
# Verify installed Nix C API packages, headers, declared function coverage, and build.
# Run from the repository root, preferably inside nix develop.
set -euo pipefail

packages=(nix-util-c nix-store-c nix-expr-c nix-fetchers-c nix-flake-c nix-main-c)
headers=(nix_api_util.h nix_api_store.h nix_api_store/store_path.h nix_api_store/derivation.h nix_api_expr.h nix_api_value.h nix_api_external.h nix_api_fetchers.h nix_api_flake.h nix_api_main.h)

echo '== pkg-config packages ==' 
for package in "${packages[@]}"; do
  if pkg-config --exists "$package"; then
    echo "  OK  $package $(pkg-config --modversion "$package" 2>/dev/null || true)"
  else
    echo "  MISSING  $package"
    exit 1
  fi
done

echo '== C API headers ==' 
include_dirs=()
while IFS= read -r flag; do
  case "$flag" in -I*) include_dirs+=("${flag#-I}") ;; esac
done < <(pkg-config --cflags "${packages[@]}" | tr ' ' '\n' | sort -u)
if [ "${#include_dirs[@]}" -eq 0 ]; then
  echo 'error: pkg-config did not return any include directories' >&2
  exit 1
fi

header_paths=()
for header in "${headers[@]}"; do
  found=''
  for dir in "${include_dirs[@]}"; do
    if [ -f "$dir/$header" ]; then found="$dir/$header"; break; fi
  done
  if [ -z "$found" ]; then
    echo "error: missing $header in pkg-config include directories" >&2
    exit 1
  fi
  header_paths+=("$found")
  echo "  OK  $header"
done

echo '== declared C API function coverage ==' 
missing=0
for header in "${header_paths[@]}"; do
  while IFS= read -r symbol; do
    [ -z "$symbol" ] && continue
    if ! grep -R --include='*.go' -F -q "C.$symbol" internal/c; then
      echo "  MISSING binding: $symbol ($header)"
      missing=1
    fi
  done < <(perl -0777 -ne 's{/\*.*?\*/}{}gs; s{//[^\n]*}{}g; while (/\b(nix_[A-Za-z0-9_]+)\s*\(/g) { print "$1\n" }' "$header" | sort -u)
done
if [ "$missing" -ne 0 ]; then exit 1; fi
echo '  OK  every declared nix_* function is referenced by a cgo binding'

echo '== CGO build ==' 
export CGO_ENABLED=1
go build ./...
echo '  OK  go build ./...'
echo '== done ==' 
