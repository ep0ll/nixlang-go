// Example: lock a flake and realise a selected derivation output through the
// Nix C API, without invoking the Nix CLI.
//
// Usage (inside nix develop):
//   go run -tags nix ./examples/realise-flake path:./. packages x86_64-linux default
//   go run -tags nix ./examples/realise-flake github:NixOS/nixpkgs#hello packages x86_64-linux hello
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ep0ll/nixlang-go/expr"
	"github.com/ep0ll/nixlang-go/fetchers"
	"github.com/ep0ll/nixlang-go/flake"
	"github.com/ep0ll/nixlang-go/store"
	"github.com/ep0ll/nixlang-go/util"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("usage: %s <flake-ref> <output-attr> [<output-attr> ...]", os.Args[0])
	}
	refString, attrPath := os.Args[1], os.Args[2:]

	ctx := util.NewContext()
	defer ctx.Close()
	must(util.Init(ctx))
	must(store.Init(ctx))
	must(expr.Init(ctx))

	st, err := store.Open(ctx, "", nil)
	must(err)
	defer st.Close()

	fetchSettings, err := fetchers.New(ctx)
	must(err)
	defer fetchSettings.Close()

	settings, err := flake.NewSettings(ctx)
	must(err)
	defer settings.Close()

	builder, err := expr.NewBuilder(ctx, st)
	must(err)
	defer builder.Close()
	must(builder.LoadEnvironment())
	must(settings.AddToEvalStateBuilder(ctx, builder))
	state, err := builder.Build()
	must(err)
	defer state.Close()

	parseFlags, err := flake.NewParseFlags(ctx, settings)
	must(err)
	defer parseFlags.Close()

	ref, _, err := flake.ParseReference(ctx, fetchSettings, settings, parseFlags, refString)
	must(err)
	defer ref.Close()

	lockFlags, err := flake.NewLockFlags(ctx, settings)
	must(err)
	defer lockFlags.Close()
	must(lockFlags.ModeVirtual(ctx))

	locked, err := flake.Lock(ctx, fetchSettings, settings, state, lockFlags, ref)
	must(err)
	defer locked.Close()

	outputs, err := locked.RealiseOutput(ctx, settings, state, st, attrPath)
	must(err)
	defer func() {
		for _, output := range outputs {
			output.Close()
		}
	}()

	for name, path := range outputs {
		realPath, err := st.RealPath(path)
		must(err)
		fmt.Printf("%s: %s\n", name, realPath)
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
