# nixgo — Production Go bindings for the Nix C API

Idiomatic, production-ready Go bindings for the experimental [Nix C API](https://nix.dev/manual/nix/latest/c-api.html) (master branch).

## Design

| Layer | Package path | Role |
|-------|--------------|------|
| Low-level cgo | `internal/c/{util,store,expr,value,external,fetchers,flake,main}` | Thin C bindings |
| High-level | `util`, `store`, `expr`, `value`, `external`, `fetchers`, `flake`, `main` | Safe Go APIs |

**This library evaluates Nix expressions and talks to stores/flakes programmatically.** It does **not** shell out to `nix build`, `nix develop`, or other CLI commands.

## Requirements

- Go 1.22+
- CGO enabled
- Nix master (or release shipping C API libs) with development headers
- `pkg-config` for: `nix-util-c`, `nix-store-c`, `nix-expr-c`, `nix-fetchers-c`, `nix-flake-c`, `nix-main-c`

```bash
nix develop   # from this flake
go test ./...
go run ./examples/eval
```

## Quick start

```go
package main

import (
	"fmt"
	"log"

	"github.com/ep0ll/nixlang-go/expr"
	"github.com/ep0ll/nixlang-go/store"
	"github.com/ep0ll/nixlang-go/util"
)

func main() {
	ctx := util.NewContext()
	defer ctx.Close()

	must(util.Init(ctx))
	must(store.Init(ctx))
	must(expr.Init(ctx))

	st, err := store.Open(ctx, "dummy://", nil)
	must(err)
	defer st.Close()

	state, err := expr.NewState(ctx, st, nil)
	must(err)
	defer state.Close()

	v, err := state.EvalString(`builtins.toJSON { answer = 6 * 7; }`, ".")
	must(err)
	defer v.Close()

	s, err := v.String()
	must(err)
	fmt.Println(s) // {"answer":42}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
```

## Packages

| Package | Capabilities |
|---------|--------------|
| **util** | Context, errors, settings, version, verbosity |
| **store** | Open store, paths, realise, closures, derivations, copy |
| **expr** | EvalState builder, EvalString, Call, CallMulti |
| **value** | Force, typed get/set, lists, attrs, builders |
| **fetchers** | Fetcher settings |
| **flake** | Settings, parse refs, lock modes, output attrs |
| **external** | Foreign values |
| **main** | Plugins, log format |

## Stability

The upstream C API is **experimental**. Pin the Nix revision via `flake.lock`.

## License

LGPL-2.1-or-later (aligned with Nix).
