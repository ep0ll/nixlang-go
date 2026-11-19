# Pure-Go Nix Expression Engine: Port Scope and Acceptance Criteria

## Non-negotiable boundary

The Nix expression language implementation itself must be fully implemented in Go. No C++/C parser, lexer, AST, evaluator, thunk engine, or expression-language semantics may be called through cgo. Expression parsing and evaluation must work in pure Go and be buildable/distributable as a Go module (including zip archives) without requiring Nix shared libraries.

CGO is optional integration for functionality outside the expression-language core, including store operations, fetchers, selected builtins/primops, derivation realization, flake operations, and native plugins where practical. Such adapters must be behind Go interfaces and optional build tags; the pure-Go evaluator must build and run without CGO or Nix installed.

## Required expression-core coverage

Port all language syntax, AST nodes, lexical/parser behavior, and evaluation semantics from the pinned upstream Nix revision, including:

- Full lexer/parser grammar, source positions, diagnostics, path and URI tokenization, escaped and indented strings, interpolation, comments, and parser edge cases.
- All expression forms, operators, precedence/associativity, attribute paths, selection/defaults, with, let, recursive attribute sets, inherit, conditionals, assertions, lambdas, pattern arguments (including defaults and ellipsis), function application, and list/attribute operations.
- Nix lazy evaluation: thunks, memoization, recursive scope, cycle/black-hole detection, strictness/forcing rules, call-by-need, and evaluation ordering.
- Nix value model and semantics: null, booleans, integers, floats, strings with context, paths, lists, attribute sets, functions, builtins/primops, and derivation-related values; exact coercion, equality, ordering, merge, and type behavior.
- Language-version and feature gates, experimental syntax/features, evaluation settings, purity/restricted mode, and compatible error/context propagation.
- Import/scoped import, search-path lookup, and parser/evaluation caches. File access and external effects are delegated through interfaces; core expression parsing/evaluation remains Go.

## Cgo / external-library allowance

Do not reimplement external systems unnecessarily. Prefer existing Go libraries where appropriate. Optional cgo adapters may bridge Nix C APIs for store/fetcher/flake/native builtins or primops when this achieves compatibility, but:

1. No cgo dependency is allowed in lexer/parser/AST/evaluator/thunk/value semantics.
2. The Go expression core defines public contracts; adapters implement those contracts.
3. CGO_ENABLED=0 builds must compile and run the expression core and pure builtins.
4. When optional native support is disabled/unavailable, failures must be explicit and typed—not silently approximate Nix behavior.

## Compatibility strategy

Pin and record the exact upstream Nix revision and language version. Build a differential harness that runs the same corpus against upstream Nix and Go, comparing parse success/failure, values, laziness/forcing behavior, errors (normalized only for unavoidable platform/source-path differences), string contexts, and trace/side-effect behavior. Add regression cases for every mismatch.

Acceptance requires complete syntax and semantics coverage, not merely common-expression coverage or an initial passing subset. Maintain grammar-feature and builtin/primop matrices tied to upstream source, each with implementation status and tests. Include upstream Nix tests where licensing and fixture requirements permit, plus generated and fuzz tests.

## Proposed package architecture

- syntax: lexer, parser, source locations, diagnostics, AST; pure Go.
- value: Nix value representation, string contexts, attrs/lists, function/primop values; pure Go.
- eval: environments, lazy thunks, evaluator, strictness/forcing, errors, settings, imports; pure Go.
- builtins: Go implementations and registry/contracts; pure implementations first, optional adapters for external behavior.
- store, fetchers, flake, native: Go contracts with optional implementations (including cgo adapters) kept outside expression core.
- cmd/nixgo: CLI and diagnostics; optional zip-friendly distribution.

## Completion gates

- Full grammar matrix implemented and tested.
- Full evaluator semantic matrix implemented and tested.
- Builtin/primop registry matches the target upstream revision; pure Go where feasible and optional native bridge where necessary.
- Differential suite has no unexplained mismatches for the upstream suite/corpus.
- CGO_ENABLED=0 go test ./... passes for the pure-Go core and tests.
- Native integrations have separately gated tests and clear availability detection.
- Zip distribution contains source/module and evaluates expressions without Nix installed for core functionality.
