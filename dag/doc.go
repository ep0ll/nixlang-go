// Package dag defines a validated directed acyclic graph of Nix store
// derivation nodes.
//
// Design (see COMPAT.md / project goals):
//
//	Go high-level API (store / expr / flake)
//	        ↓
//	Derivation / StorePath inspection
//	        ↓
//	dag.Graph extraction
//	        ↓
//	optional protobuf export (future)
//
// This package does not invoke the Nix CLI. Graph construction uses the
// in-process C API bindings (or pure data already obtained from them).
package dag
