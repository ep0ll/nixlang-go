// Example: evaluate a Nix expression without running any Nix CLI commands.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ep0ll/nixlang-go/expr"
	"github.com/ep0ll/nixlang-go/store"
	"github.com/ep0ll/nixlang-go/util"
)

func main() {
	ctx := util.NewContext()
	defer ctx.Close()

	if err := util.Init(ctx); err != nil {
		log.Fatal(err)
	}
	if err := store.Init(ctx); err != nil {
		log.Fatal(err)
	}
	if err := expr.Init(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Nix version:", util.Version())

	st, err := store.Open(ctx, "dummy://", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	state, err := expr.NewState(ctx, st, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer state.Close()

	exprStr := `builtins.toJSON { answer = 6 * 7; greeting = "hello from nixgo"; }`
	if len(os.Args) > 1 {
		exprStr = os.Args[1]
	}

	v, err := state.EvalString(exprStr, ".")
	if err != nil {
		log.Fatal(err)
	}
	defer v.Close()

	s, err := v.String()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
}
