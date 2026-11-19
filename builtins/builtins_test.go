//go:build cgo && nix

package builtins_test

import (
	"strings"
	"testing"

	nixbuiltins "github.com/ep0ll/nixlang-go/builtins"
	"github.com/ep0ll/nixlang-go/expr"
	"github.com/ep0ll/nixlang-go/store"
	"github.com/ep0ll/nixlang-go/util"
	"github.com/ep0ll/nixlang-go/value"
)

func newEvaluator(t *testing.T) *expr.State {
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
		t.Fatalf("store.Open: %v", err)
	}
	t.Cleanup(st.Close)

	state, err := expr.NewState(ctx, st, nil)
	if err != nil {
		t.Fatalf("expr.NewState: %v", err)
	}
	t.Cleanup(state.Close)
	return state
}

func TestGoArgumentsAndGeneratedWrappers(t *testing.T) {
	state := newEvaluator(t)
	b, err := nixbuiltins.New(state)
	if err != nil {
		t.Fatalf("builtins.New: %v", err)
	}
	t.Cleanup(b.Close)

	sum, err := b.Add(int(20), int64(22))
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	t.Cleanup(sum.Close)
	if got, err := sum.Int(); err != nil || got != 42 {
		t.Fatalf("Add result = %d, %v; want 42", got, err)
	}

	joined, err := b.CallString("concatStringsSep", "/", []string{"usr", "local", "bin"})
	if err != nil {
		t.Fatalf("CallString(concatStringsSep): %v", err)
	}
	if joined != "usr/local/bin" {
		t.Fatalf("joined = %q; want %q", joined, "usr/local/bin")
	}

	isPath, err := b.IsPath(nixbuiltins.Path("./src"))
	if err != nil {
		t.Fatalf("IsPath: %v", err)
	}
	if isPath != nil {
		t.Cleanup(isPath.Close)
		got, err := isPath.Bool()
		if err != nil || !got {
			t.Fatalf("IsPath result = %t, %v; want true", got, err)
		}
	} else {
		t.Fatal("IsPath returned a nil value")
	}

	answer, err := b.GetAttr("answer", map[string]any{"answer": 42})
	if err != nil {
		t.Fatalf("GetAttr: %v", err)
	}
	t.Cleanup(answer.Close)
	if got, err := answer.Int(); err != nil || got != 42 {
		t.Fatalf("GetAttr result = %d, %v; want 42", got, err)
	}
}

func TestExistingNixFunctionArgument(t *testing.T) {
	state := newEvaluator(t)
	b, err := nixbuiltins.New(state)
	if err != nil {
		t.Fatalf("builtins.New: %v", err)
	}
	t.Cleanup(b.Close)

	fn, err := state.EvalString("x: x + 1", ".")
	if err != nil {
		t.Fatalf("EvalString: %v", err)
	}
	t.Cleanup(fn.Close)

	mapped, err := b.Map(fn, []int{1, 2, 3})
	if err != nil {
		t.Fatalf("Map: %v", err)
	}
	t.Cleanup(mapped.Close)

	if got, err := mapped.ListLen(); err != nil || got != 3 {
		t.Fatalf("Map result length = %d, %v; want 3", got, err)
	}
	second, err := mapped.ListAt(1)
	if err != nil {
		t.Fatalf("ListAt: %v", err)
	}
	t.Cleanup(second.Close)
	if got, err := second.Int(); err != nil || got != 3 {
		t.Fatalf("Map result[1] = %d, %v; want 3", got, err)
	}
}

func TestRegistryAndRuntimeDiscovery(t *testing.T) {
	state := newEvaluator(t)
	b, err := nixbuiltins.New(state)
	if err != nil {
		t.Fatalf("builtins.New: %v", err)
	}
	t.Cleanup(b.Close)

	present, err := b.Has("add")
	if err != nil || !present {
		t.Fatalf("Has(add) = %t, %v; want true", present, err)
	}

	names, err := b.Names()
	if err != nil {
		t.Fatalf("Names: %v", err)
	}
	if len(names) == 0 {
		t.Fatal("Names returned an empty builtin set")
	}

	specs := nixbuiltins.Registry()
	foundAdd := false
	for _, spec := range specs {
		if spec.Name == "add" && spec.Arity == 2 {
			foundAdd = len(spec.Args) == 2 &&
				len(spec.GoTypes) == 2 &&
				spec.GoTypes[0] == "any" &&
				spec.GoTypes[1] == "any"
			break
		}
	}
	if !foundAdd {
		t.Fatal("registry does not describe add/2")
	}

	add, err := b.Get("add")
	if err != nil {
		t.Fatalf("Get(add): %v", err)
	}
	t.Cleanup(add.Close)
	if add.Type() != value.TypeFunction {
		t.Fatalf("Get(add) type = %s; want function", add.Type())
	}
}


func TestStaticallyTypedBuiltinArguments(t *testing.T) {
	state := newEvaluator(t)
	b, err := nixbuiltins.New(state)
	if err != nil {
		t.Fatalf("builtins.New: %v", err)
	}
	t.Cleanup(b.Close)

	hash, err := b.HashString("sha256", "hello")
	if err != nil {
		t.Fatalf("HashString: %v", err)
	}
	t.Cleanup(hash.Close)
	hashString, err := hash.String()
	if err != nil {
		t.Fatalf("HashString result: %v", err)
	}
	if hashString == "" {
		t.Fatal("HashString returned an empty hash")
	}

	converted, err := b.ConvertHash(nixbuiltins.ConvertHashArgs{
		Hash:         "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		HashAlgo:     "sha256",
		ToHashFormat: "sri",
	})
	if err != nil {
		t.Fatalf("ConvertHash: %v", err)
	}
	t.Cleanup(converted.Close)
	convertedString, err := converted.String()
	if err != nil {
		t.Fatalf("ConvertHash result: %v", err)
	}
	if !strings.HasPrefix(convertedString, "sha256-") {
		t.Fatalf("ConvertHash result = %q; want SRI sha256 hash", convertedString)
	}

	joined, err := b.ConcatStringsSep("/", []string{"usr", "local", "bin"})
	if err != nil {
		t.Fatalf("ConcatStringsSep: %v", err)
	}
	t.Cleanup(joined.Close)
	if got, err := joined.String(); err != nil || got != "usr/local/bin" {
		t.Fatalf("ConcatStringsSep result = %q, %v; want usr/local/bin", got, err)
	}

	appended, err := b.AppendContext("plain", nixbuiltins.StringContext{})
	if err != nil {
		t.Fatalf("AppendContext: %v", err)
	}
	t.Cleanup(appended.Close)
	if got, err := appended.String(); err != nil || got != "plain" {
		t.Fatalf("AppendContext result = %q, %v; want plain", got, err)
	}
}
