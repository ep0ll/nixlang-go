# Compatibility contract

The Nix C API is **experimental**. This library pins a specific upstream
revision and treats it as the supported ABI.

## Supported combination

| Component | Version / revision |
|-----------|--------------------|
| **nixlang-go** | main (see git tags for releases) |
| **Nix** | `github:NixOS/nix` @ `db7dc787a492` (see `flake.lock`) |
| **nixpkgs** (dev shell) | `nixos-unstable` @ `e158d9ed9b51` |
| **Go** | 1.22+ |
| **C API packages** | `nix-util-c`, `nix-store-c`, `nix-expr-c`, `nix-fetchers-c`, `nix-flake-c`, `nix-main-c` |

Pin source of truth: **`flake.lock`** in this repository.

```bash
# Inspect the locked Nix revision
nix flake metadata
# or:
python3 -c "import json; d=json.load(open('flake.lock')); print(d['nodes']['nix']['locked']['rev'])"
```

## Policy

1. **Do not** track a floating `master` tip in production. Always use a
   checkout whose `flake.lock` matches a known good revision.
2. When upgrading Nix, update `flake.lock`, run `./scripts/check-capi.sh`,
   and the full `go test -tags=nix ./...` suite before merging.
3. Breaking C header changes require a **new major** or explicit
   compatibility note in `CHANGELOG.md`.
4. Releases (when tagged) MUST list the exact Nix commit they were tested
   against.

## CI expectation

The `nix-capi` job **must** succeed against the locked revision. A failure
of `check-capi.sh` or of tagged integration tests fails the workflow —
there is no silent fallback to CGO-disabled stubs for the native job.

## Supported platforms

| OS | Arch | Unit (CGO=0) | Native (`-tags=nix`) | Notes |
|----|------|--------------|----------------------|-------|
| Linux | amd64 | ✅ CI | ✅ CI (`ubuntu-latest`) | Primary |
| Linux | arm64 | ✅ local | 🟡 best-effort | Same flake systems list |
| macOS | amd64 | ✅ CI unit | 🟡 best-effort | C API via `nix develop` |
| macOS | arm64 | ✅ CI unit | 🟡 best-effort | C API via `nix develop` |

Native integration is **required to pass on Linux amd64** against the pinned
Nix revision. Other platforms are supported at the source level (`flake.nix`
`systems`) but are not currently hard-gated in CI.

When adding a release, list which matrix cells were green.

## Dual validation (CI)

CI runs the native suite against **two** C API sources:

| Job | Source | Role |
|-----|--------|------|
| `nix-capi` | `flake.lock` → Nix input (pinned commit) | **Primary** ABI contract — must pass |
| `nix-capi-nixpkgs` | `nixpkgs#nixVersions.latest` (+ `.dev`) | Secondary — catches packaging / release-line drift |

Both must remain green for merges. When they disagree, treat the **flake pin**
as the source of truth and open an issue documenting the nixpkgs delta.

Optional `asan` job (`continue-on-error`) runs `./scripts/check-asan.sh` on
`main` pushes and `workflow_dispatch`.

## Concurrency

- A single `util.Context` is **not** safe for concurrent API calls.
- `Close` on Context, Store, State, and Value is idempotent and safe to
  call from multiple goroutines (uses `sync.Once` / nil-out).
- String callback slot registration (`internal/c/util`) is mutex-protected.
- Prefer one Context (and derived Store/State) per goroutine, or external
  synchronization around shared handles.

## Versioning and releases

Until the Nix C API stabilizes upstream:

1. **Pre-1.0** tags may break on Nix upgrades even for “minor” bumps when
   headers change.
2. Each release tag **must** record in the GitHub release notes:
   - `flake.lock` Nix `rev` (full hash)
   - nixpkgs `rev` used for secondary validation
   - Go version
   - Platforms where `nix-capi` and `nix-capi-nixpkgs` were green
3. Module path remains `github.com/ep0ll/nixlang-go`. Consumers should pin
   both this module and a matching Nix revision (or use this flake’s lock).
4. Changelog entries that touch C bindings should mention whether a Nix
   bump was required.
