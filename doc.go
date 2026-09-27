// Package nixgo is the root module for production Go bindings to the Nix C API.
//
// Import the subpackages you need:
//
//	github.com/ep0ll/nixlang-go/util      — contexts, errors, settings
//	github.com/ep0ll/nixlang-go/store     — stores, paths, closures
//	github.com/ep0ll/nixlang-go/expr      — EvalState, EvalString, Call
//	github.com/ep0ll/nixlang-go/value     — typed Nix values
//	github.com/ep0ll/nixlang-go/fetchers  — fetcher settings
//	github.com/ep0ll/nixlang-go/flake     — flake refs, lock, outputs
//	github.com/ep0ll/nixlang-go/external  — foreign values
//	github.com/ep0ll/nixlang-go/main      — plugins, log format
//
// Low-level cgo bindings live under internal/c/* and are not part of the
// stable public API.
//
// This module evaluates Nix expressions and interacts with stores/flakes via
// the C library only. It does not execute `nix build`, `nix develop`, or other
// CLI commands.
package nixgo
