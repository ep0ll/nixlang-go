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


// TestPrimOpValueLifetime exercises both nix_init_primop and the GC-backed
// callback handle lifetime: the Go PrimOp wrapper is closed while its Nix
// value still retains the underlying primop.
func TestPrimOpValueLifetime(t *testing.T) {
	ctx, _, state := setupEval(t)

	primop, err := value.NewPrimOp(
		ctx,
		func(_ any, _ *util.Context, _ *value.PrimOpState, args []*value.Value, result *value.Value) error {
			n, err := args[0].Int()
			if err != nil {
				return err
			}
			return result.SetInt(n + 1)
		},
		1, "goAddOne", []string{"x"}, "Add one to an integer.", nil,
	)
	if err != nil {
		t.Fatalf("NewPrimOp: %v", err)
	}

	fn, err := state.AllocValue()
	if err != nil {
		t.Fatalf("AllocValue(function): %v", err)
	}
	t.Cleanup(fn.Close)
	if err := fn.SetPrimOp(primop); err != nil {
		t.Fatalf("SetPrimOp: %v", err)
	}
	primop.Close()

	arg, err := state.AllocValue()
	if err != nil {
		t.Fatalf("AllocValue(argument): %v", err)
	}
	t.Cleanup(arg.Close)
	if err := arg.SetInt(41); err != nil {
		t.Fatalf("SetInt: %v", err)
	}

	out, err := state.Call(fn, arg)
	if err != nil {
		t.Fatalf("Call custom primop after closing wrapper: %v", err)
	}
	t.Cleanup(out.Close)
	got, err := out.Int()
	if err != nil {
		t.Fatalf("result Int: %v", err)
	}
	if got != 42 {
		t.Fatalf("custom primop result = %d, want 42", got)
	}
}

// TestCallAuto evaluates a curried attrset-argument function through the C API.
func TestCallAuto(t *testing.T) {
	_, _, state := setupEval(t)

	fn, err := state.EvalString(`{ a, b ? 1 }: a + b`, ".")
	if err != nil {
		t.Fatalf("EvalString(function): %v", err)
	}
	t.Cleanup(fn.Close)

	args, err := state.EvalString(`{ a = 41; }`, ".")
	if err != nil {
		t.Fatalf("EvalString(arguments): %v", err)
	}
	t.Cleanup(args.Close)

	out, err := state.CallAuto(fn, args)
	if err != nil {
		t.Fatalf("CallAuto: %v", err)
	}
	t.Cleanup(out.Close)
	got, err := out.Int()
	if err != nil {
		t.Fatalf("result Int: %v", err)
	}
	if got != 42 {
		t.Fatalf("CallAuto result = %d, want 42", got)
	}
}

// TestBuildersAndIndexedAttrs exercises the C list/attr builders and the lazy
// indexed attribute accessor.
func TestBuildersAndIndexedAttrs(t *testing.T) {
	_, _, state := setupEval(t)

	one, err := state.AllocValue()
	if err != nil {
		t.Fatalf("AllocValue(one): %v", err)
	}
	t.Cleanup(one.Close)
	if err := one.SetInt(1); err != nil {
		t.Fatalf("SetInt(one): %v", err)
	}

	two, err := state.AllocValue()
	if err != nil {
		t.Fatalf("AllocValue(two): %v", err)
	}
	t.Cleanup(two.Close)
	if err := two.SetInt(2); err != nil {
		t.Fatalf("SetInt(two): %v", err)
	}

	list, err := state.NewList(one, two)
	if err != nil {
		t.Fatalf("NewList: %v", err)
	}
	t.Cleanup(list.Close)
	if got, err := list.ListLen(); err != nil || got != 2 {
		t.Fatalf("ListLen = %d, %v; want 2", got, err)
	}

	attrs, err := state.NewAttrs(map[string]*value.Value{"one": one, "two": two})
	if err != nil {
		t.Fatalf("NewAttrs: %v", err)
	}
	t.Cleanup(attrs.Close)

	names, err := attrs.AttrNames()
	if err != nil {
		t.Fatalf("AttrNames: %v", err)
	}
	for i := range names {
		name, lazyValue, err := attrs.AttrAtLazy(i)
		if err != nil {
			t.Fatalf("AttrAtLazy(%d): %v", i, err)
		}
		got, err := lazyValue.Int()
		lazyValue.Close()
		if err != nil {
			t.Fatalf("AttrAtLazy(%d).Int: %v", i, err)
		}
		switch name {
		case "one":
			if got != 1 {
				t.Fatalf("attrs.one = %d, want 1", got)
			}
		case "two":
			if got != 2 {
				t.Fatalf("attrs.two = %d, want 2", got)
			}
		default:
			t.Fatalf("unexpected attribute name %q", name)
		}
	}
}
