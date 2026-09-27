# Building nixgo

## Prerequisites

1. Nix **master** (or a release that ships the C API packages)
2. Development outputs: `nix.dev` / individual `nix-*-c` packages
3. `pkg-config`, `gcc`, Go 1.22+

## Reproducible shell

```bash
nix develop
export CGO_ENABLED=1
go build ./...

# Pure unit tests (no Nix C libraries required)
go test ./util ./value

# Integration tests against the Nix C API (dummy store, eval only)
go test -tags=nix ./...

go run ./examples/eval '1 + 1'
```

## Linking

Each `internal/c/*` package uses:

```
#cgo pkg-config: nix-util-c
#cgo pkg-config: nix-store-c nix-util-c
#cgo pkg-config: nix-expr-c nix-store-c nix-util-c
#cgo pkg-config: nix-fetchers-c nix-util-c
#cgo pkg-config: nix-flake-c nix-fetchers-c nix-expr-c nix-store-c nix-util-c
#cgo pkg-config: nix-main-c nix-util-c
```

If `pkg-config` cannot find these, set:

```bash
export PKG_CONFIG_PATH=/path/to/nix/lib/pkgconfig:$PKG_CONFIG_PATH
```

## Notes

- The C API is experimental; pin the Nix revision in `flake.lock`.
- Store open `params` map is reserved; configure via URI or `util.SetSetting`.
- Always call `Close()` on Context, Store, State, Value, and related objects.
