// Package nixgo is the root module for Go bindings to the public Nix C API.
//
// Import the subpackages you need:
//
//	github.com/ep0ll/nixlang-go/util      — contexts, errors, settings
//	github.com/ep0ll/nixlang-go/store     — stores, paths, closures
//	github.com/ep0ll/nixlang-go/expr      — evaluator states and expression calls
//	github.com/ep0ll/nixlang-go/value     — typed values and custom primops
//	github.com/ep0ll/nixlang-go/external  — Go-backed external values
//	github.com/ep0ll/nixlang-go/fetchers  — fetcher settings
//	github.com/ep0ll/nixlang-go/flake     — flake references, locks, outputs
//	github.com/ep0ll/nixlang-go/nixmain   — plugin initialization and log format
//
// This module binds Nix's exported C API. It is not a binding of the complete
// Nix C++ implementation, internal C++ classes, CLI, or daemon protocols.
// The upstream C API is experimental; pin and test against a specific Nix
// revision.
//
// # Concurrency
//
// A util.Context, store.Store, and expr.State must not be shared across
// goroutines without external synchronization. Prefer one Context/State pair
// per goroutine. Callback registries are synchronized.
//
// # Resource lifetime
//
// Call Close() on Context, Store, State, Value, and related owned objects.
// Finalizers are a safety net only; do not rely on them for timely cleanup.
// A Value is allocated in its evaluator's memory arena. Keep the owning
// expr.State open until all Values and realised strings created from that state
// have been closed, including Values retained from a primop callback.
package nixgo
