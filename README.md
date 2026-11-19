# nixgo — Go bindings for the Nix C API

Idiomatic Go wrappers for the public, experimental Nix C API. The project links against Nix's exported C libraries with cgo; it does not reimplement the Nix evaluator or bind directly to Nix's internal C++ classes.

## Scope

This is not a complete C++ binding of the Nix codebase. Nix's C++ internals, CLI, daemon protocols, and APIs that are not exported by the public C API are outside the binding surface. The goal is to expose the published nix_api_*.h interfaces to Go, with low-level cgo packages and higher-level Go wrappers.

The feat/full-cgo-bindings branch expands coverage across the current libutil-c, libstore-c, libexpr-c, libfetchers-c, libflake-c, and libmain-c headers. It includes evaluator values and builders, custom primops, external-value callbacks, GC hooks, derivation inspection, store closure traversal, fetcher settings, flakes, and main-library setup. The scripts/check-capi.sh script audits the declared function names in the installed headers against the cgo call sites.

The upstream C API is experimental and may change between Nix revisions. Pin the Nix revision used to build the bindings and run the C API checks against that exact revision.

## Design

| Layer | Package path | Role |
|-------|--------------|------|
| Low-level cgo | internal/c/{util,store,expr,value,external,fetchers,flake,nixmain} | Thin bindings to Nix C functions |
| High-level | util, store, expr, value, external, fetchers, flake, nixmain | Typed Go APIs and resource wrappers |

The library evaluates expressions and talks to stores and flakes through Nix libraries. It does not shell out to nix build, nix develop, or other CLI commands.

## Requirements

- Go 1.22+
- CGO enabled for Nix integration
- Nix with the corresponding C API development headers and libraries
- pkg-config entries for nix-util-c, nix-store-c, nix-expr-c, nix-fetchers-c, nix-flake-c, and nix-main-c

Run nix develop, then ./scripts/check-capi.sh and go test -tags=nix ./...


### Flake output realisation

The flake package can evaluate output attributes, identify a derivation, and
realise it through the store API without shelling out to `nix build`. Select
outputs using a slice of attribute names (not a dotted string, so attribute names
containing dots remain unambiguous):

```go
outputs, err := locked.RealiseOutput(ctx, flakeSettings, evalState, st,
    []string{"packages", "x86_64-linux", "default"})
if err != nil {
    return err
}
defer func() {
    for _, output := range outputs {
        output.Close()
    }
}()

for outputName, path := range outputs {
    realPath, err := st.RealPath(path)
    if err != nil {
        return err
    }
    fmt.Printf("%s: %s\\n", outputName, realPath)
}
```

This builds the selected derivation and returns owned store-path handles for
its outputs. Call `Close` on each returned path. The complete runnable example
is in `examples/realise-flake`; it takes a flake reference followed by the
output attribute path.

## Packages

| Package | Capabilities |
|---------|--------------|
| util | Contexts, structured errors, settings, version, verbosity |
| store | Stores, paths, closure traversal, realise, derivations, copy |
| expr | EvalState builders, evaluation, function calls, derivation inspection |
| value | Typed values, lists, attrs, builders, retained values, custom primops and callback-scoped evaluator operations |
| external | Go-backed Nix external values and callback descriptors |
| fetchers | Fetcher settings |
| flake | Flake settings, reference parsing, locking, output evaluation and realisation |
| nixmain | Plugin initialization and log format |

## Resource lifetime

Close contexts, stores, states, values, descriptors, and other owned handles when finished. Finalizers are only a fallback. Values passed to primop and external-value callbacks are borrowed for the duration of the callback; use Value.Retain or the relevant Wrap API when a reference must outlive it. Custom primops receive a callback-scoped value.PrimOpState facade for allocating values, evaluating expressions, and calling functions. Values returned from that facade own their error contexts, but they still belong to the parent evaluator's memory arena: keep the originating expr.State open until every such value is closed. The public C API has no deterministic callback-userdata destructor for all Nix GC configurations, so Go callback handles attached to primops installed in values or the global builtin registry are retained for the process lifetime.

## License

LGPL-2.1-or-later (aligned with Nix).