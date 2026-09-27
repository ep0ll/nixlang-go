# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
The upstream Nix C API is **experimental**; this binding tracks Nix master and
may change when header signatures change.

## [Unreleased]

### Added

- `store.Open` encodes `params map[string]string` as `const char ***` for Nix master
- Idempotent `Free`/`Decref` via `sync.Once` (store, path, derivation, value, context)
- `expr.EvalFile` — evaluate a `.nix` file without CLI
- `store.Path.String`, safer `Name`/`Clone`
- `util.NixError.Is`, `AsNixError`, `IsKeyError`, `IsRecoverable`
- Package rename: `main` → `nixmain`
- Integration tests gated by `//go:build cgo && nix`
- `scripts/check-capi.sh` and improved `flake.nix` / `BUILD.md`
- GitHub Actions unit-test workflow
- Full LGPL-2.1 `LICENSE` text
- `examples/flake` parse-only demo
- `value.TypeName`, `CopyFrom`, `RealiseString` (string context)
- Panic recovery in all cgo `//export` callbacks
- `flake` parse integration test (`cgo && nix`)

### Guarantees

- Evaluation and store/flake APIs use the C library only — no `nix build`,
  `nix develop`, or other CLI invocation.
