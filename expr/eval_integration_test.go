//go:build cgo && nix

package expr_test

import (
	"strings"
	"testing"

	"github.com/ep0ll/nixlang-go/expr"
	"github.com/ep0ll/nixlang-go/store"
	"github.com/ep0ll/nixlang-go/util"
	"github.com/ep0ll/nixlang-go/value"
)

// setupEval opens a dummy store and creates an EvalState.
// Requires Nix C libraries (run via: nix develop && go test -tags=nix ./...).
func setupEval(t *testing.T) (*util.Context, *store.Store, *expr.State) {
	t.Helper()
	ctx := util.NewContext()
	t.Cleanup(ctx.Close)

	if err := util.Init(ctx); err != nil {
		t.Fatalf("util.Init: %v", err)
	}
	if err := store.Init(ctx); err != nil {
		t.Fatalf("store.Init: %v", err)
	}
	if err := expr.Init(ctx); err != nil {
		t.Fatalf("expr.Init: %v", err)
	}

	st, err := store.Open(ctx, "dummy://", nil)
	if err != nil {
		t.Fatalf("store.Open(dummy://): %v", err)
	}
	t.Cleanup(st.Close)

	state, err := expr.NewState(ctx, st, nil)
	if err != nil {
		t.Fatalf("expr.NewState: %v", err)
	}
	t.Cleanup(state.Close)

	return ctx, st, state
}

func TestEvalStringArithmetic(t *testing.T) {
	_, _, state := setupEval(t)

	v, err := state.EvalString(`1 + 2 * 3`, ".")
	if err != nil {
		t.Fatalf("EvalString: %v", err)
	}
	t.Cleanup(v.Close)

	n, err := v.Int()
	if err != nil {
		t.Fatalf("Int: %v", err)
	}
	if n != 7 {
		t.Fatalf("got %d, want 7", n)
	}
}

func TestEvalStringJSON(t *testing.T) {
	_, _, state := setupEval(t)

	v, err := state.EvalString(`builtins.toJSON { answer = 6 * 7; }`, ".")
	if err != nil {
		t.Fatalf("EvalString: %v", err)
	}
	t.Cleanup(v.Close)

	s, err := v.String()
	if err != nil {
		t.Fatalf("String: %v", err)
	}
	if !strings.Contains(s, `"answer":42`) && !strings.Contains(s, `"answer": 42`) {
		t.Fatalf("unexpected JSON: %q", s)
	}
}

func TestEvalAttrs(t *testing.T) {
	_, _, state := setupEval(t)

	v, err := state.EvalString(`{ a = 1; b = "hi"; }`, ".")
	if err != nil {
		t.Fatalf("EvalString: %v", err)
	}
	t.Cleanup(v.Close)

	if v.Type() != value.TypeAttrs {
		t.Fatalf("type = %s, want attrs", v.Type())
	}

	names, err := v.AttrNames()
	if err != nil {
		t.Fatalf("AttrNames: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("AttrNames = %v, want 2 names", names)
	}

	av, err := v.Attr("a")
	if err != nil {
		t.Fatalf("Attr(a): %v", err)
	}
	t.Cleanup(av.Close)
	n, err := av.Int()
	if err != nil {
		t.Fatalf("a.Int: %v", err)
	}
	if n != 1 {
		t.Fatalf("a = %d, want 1", n)
	}

	bv, err := v.Attr("b")
	if err != nil {
		t.Fatalf("Attr(b): %v", err)
	}
	t.Cleanup(bv.Close)
	s, err := bv.String()
	if err != nil {
		t.Fatalf("b.String: %v", err)
	}
	if s != "hi" {
		t.Fatalf("b = %q, want hi", s)
	}
}

func TestEvalList(t *testing.T) {
	_, _, state := setupEval(t)

	v, err := state.EvalString(`[1 2 3]`, ".")
	if err != nil {
		t.Fatalf("EvalString: %v", err)
	}
	t.Cleanup(v.Close)

	n, err := v.ListLen()
	if err != nil {
		t.Fatalf("ListLen: %v", err)
	}
	if n != 3 {
		t.Fatalf("len = %d, want 3", n)
	}

	e, err := v.ListAt(1)
	if err != nil {
		t.Fatalf("ListAt(1): %v", err)
	}
	t.Cleanup(e.Close)
	x, err := e.Int()
	if err != nil {
		t.Fatalf("element Int: %v", err)
	}
	if x != 2 {
		t.Fatalf("element = %d, want 2", x)
	}
}

func TestEvalCall(t *testing.T) {
	_, _, state := setupEval(t)

	fn, err := state.EvalString(`x: x + 1`, ".")
	if err != nil {
		t.Fatalf("EvalString fn: %v", err)
	}
	t.Cleanup(fn.Close)

	arg, err := state.AllocValue()
	if err != nil {
		t.Fatalf("AllocValue: %v", err)
	}
	t.Cleanup(arg.Close)
	if err := arg.SetInt(41); err != nil {
		t.Fatalf("SetInt: %v", err)
	}

	out, err := state.Call(fn, arg)
	if err != nil {
		t.Fatalf("Call: %v", err)
	}
	t.Cleanup(out.Close)

	n, err := out.Int()
	if err != nil {
		t.Fatalf("out.Int: %v", err)
	}
	if n != 42 {
		t.Fatalf("got %d, want 42", n)
	}
}

func TestEvalTypeMismatch(t *testing.T) {
	_, _, state := setupEval(t)

	v, err := state.EvalString(`"hello"`, ".")
	if err != nil {
		t.Fatalf("EvalString: %v", err)
	}
	t.Cleanup(v.Close)

	if _, err := v.Int(); err == nil {
		t.Fatal("expected type error for Int on string")
	}
}

func TestStoreURI(t *testing.T) {
	ctx, st, _ := setupEval(t)

	uri, err := st.URI()
	if err != nil {
		t.Fatalf("URI: %v", err)
	}
	if uri == "" {
		t.Fatal("empty store URI")
	}
	t.Logf("store URI: %s (version %s)", uri, util.Version())
	_ = ctx
}
