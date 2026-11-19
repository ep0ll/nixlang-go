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
    fmt.Printf("%s: %s\n", outputName, realPath)
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

## Nix expression builtins from Go

The `builtins` package exposes the Nix 2.35.6 evaluator builtin set through
ordinary Go APIs. It does not reimplement Nix builtin semantics: calls are
dispatched to the active Nix evaluator through the generic C API value-call
interface. You pass Go values, and the package converts them into evaluator
values before calling the actual builtin.

```go
b, err := builtins.New(state)
if err != nil {
    return err
}
defer b.Close()

// Generated, arity-specific Go wrapper.
sum, err := b.Add(20, 22)
if err != nil {
    return err
}
defer sum.Close()
n, err := sum.Int() // 42

// Go slices and maps become Nix lists and attribute sets.
joined, err := b.CallString("concatStringsSep", "/", []string{"usr", "local", "bin"})
answer, err := b.GetAttr("answer", map[string]any{"answer": 42})

// Use builtins.Path when a Go string must be a Nix path, rather than a string.
isPath, err := b.IsPath(builtins.Path("./src"))

// Wrap an evaluated Nix closure before passing it to higher-order builtins.
fnValue, err := state.EvalString("x: x + 1", ".")
defer fnValue.Close()
fn, err := builtins.NewFunction(fnValue)
defer fn.Close()
mapped, err := b.Map(fn, []any{1, 2, 3})
```

Supported Go arguments include `nil`, booleans, signed and unsigned integers
(fitting Nix's signed 64-bit integer range), floats, strings, `builtins.Path`,
slices and arrays, string-keyed maps, and `*value.Value` handles from the same
evaluator state. `[]byte` is treated as a Nix string. Function parameters are
passed as checked `builtins.Function` values created with
`builtins.NewFunction` from an evaluator value. Other arbitrary evaluator
values can be supplied as same-state `*value.Value` handles. Results are
returned as `*value.Value` so Nix types,
attribute sets, lists, and store-path string context are not silently flattened.

Use `b.Get(name)` for constants such as `currentTime`, `currentSystem`,
`nixVersion`, `storeDir`, `langVersion`, and `nixPath`; `b.Names()` and
`b.Has(name)` inspect the builtin set provided by the running evaluator.
`builtins.Registry()` describes the Nix 2.35.6 callable catalog, including
arity and experimental-feature gates. Runtime availability can differ when a
builtin is feature-gated, platform-dependent, or when the linked Nix version
differs from 2.35.6.

The generated methods use statically typed Go parameters wherever the
Nix argument has a clear shape: strings use `string`, integer operands use
`int64`, function operands use `builtins.Function`, and structured Nix
attribute sets use named Go structs such as `PathArgs`, `FetchClosureArgs`,
`ConvertHashArgs`, `DerivationAttrs`, and `FlakeRefAttrs`. Lists with a
known element type use typed slices, for example `[]string` and
`[]SearchPathEntry`. Lists with arbitrary Nix values use `[]any` (and
nested lists use `[][]any`) rather than accepting a scalar `any` parameter. A parameter is `any` only when Nix itself accepts
multiple distinct Nix value types or an arbitrary expression value (for
example, `HashFile(algorithm string, path any)` accepts either a Nix path or a
string path).

```go
// A string algorithm plus a Nix path or string path.
hash, err := b.HashFile("sha256", builtins.Path("./LICENSE"))

// A Nix attribute set represented by a typed Go struct.
converted, err := b.ConvertHash(builtins.ConvertHashArgs{
    Hash:         "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    HashAlgo:     "sha256",
    ToHashFormat: "sri",
})

// Higher-order functions take checked Nix function handles.
fnValue, err := state.EvalString("x: x + 1", ".")
fn, err := builtins.NewFunction(fnValue)
defer fn.Close()
mapped, err := b.Map(fn, []any{1, 2, 3})

// Structured path arguments, including a Nix function predicate.
path, err := b.Path(builtins.PathArgs{
    Path:   builtins.Path("./source"),
    Filter: fn,
})
```

The exported `builtins.Registry()` also records each Nix parameter name and
its generated Go parameter type in `Spec.Args` and `Spec.GoTypes`. Dynamic
`Builtins.Call` remains available for plugin-defined or otherwise unusual
calls, while the generated methods provide the statically typed API.

The catalog includes core `libexpr` primops, string-context and fetcher
primops, and the three flake primops (`getFlake`, `parseFlakeRef`, and
`flakeRefToString`). Flake primops are only available when flake settings
have been added to the evaluator. The optional `importNative` and `exec`
primops require Nix's native-code setting and are not available on Windows.
The internal `fetchFinalTree` primop is deliberately not exposed.
