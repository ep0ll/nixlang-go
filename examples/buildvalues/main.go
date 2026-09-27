// Example: build Nix values in Go and call a Nix function (eval only, no CLI).
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

	// Evaluate a function in Nix, then call it with a Go-built argument.
	fn, err := state.EvalString(`x: x * 2 + 1`, ".")
	must(err)
	defer fn.Close()

	arg, err := state.EvalString(`21`, ".")
	must(err)
	defer arg.Close()

	out, err := state.Call(fn, arg)
	must(err)
	defer out.Close()

	n, err := out.Int()
	must(err)
	fmt.Println("21 * 2 + 1 =", n)

	// Attrset via evaluation (builders are also available on value package).
	attrs, err := state.EvalString(`{ answer = 6 * 7; label = "nixlang-go"; }`, ".")
	must(err)
	defer attrs.Close()

	a, err := attrs.Attr("answer")
	must(err)
	defer a.Close()
	ans, err := a.Int()
	must(err)
	fmt.Println("attrs.answer =", ans)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
