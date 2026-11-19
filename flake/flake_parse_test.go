//go:build cgo && nix

package flake_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ep0ll/nixlang-go/expr"
	"github.com/ep0ll/nixlang-go/fetchers"
	"github.com/ep0ll/nixlang-go/flake"
	"github.com/ep0ll/nixlang-go/store"
	"github.com/ep0ll/nixlang-go/util"
)

func TestParseReference(t *testing.T) {
	ctx := util.NewContext()
	defer ctx.Close()
	if err := util.Init(ctx); err != nil {
		t.Fatal(err)
	}

	fs, err := fetchers.New(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer fs.Close()

	settings, err := flake.NewSettings(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer settings.Close()

	flags, err := flake.NewParseFlags(ctx, settings)
	if err != nil {
		t.Fatal(err)
	}
	defer flags.Close()

	ref, frag, err := flake.ParseReference(ctx, fs, settings, flags, "github:NixOS/nixpkgs")
	if err != nil {
		t.Fatal(err)
	}
	defer ref.Close()
	if ref == nil {
		t.Fatal("nil reference")
	}
	_ = frag
}

// TestLockVirtual is skipped by default: locking needs network/store.
// Enable manually when a store and network are available.
func TestLockVirtual(t *testing.T) {
	t.Skip("requires network and writable store; run manually with -run TestLockVirtual")
}

func TestRealiseOutputLocalFlake(t *testing.T) {
	if os.Getenv("NIXGO_TEST_FLAKE_REALISE") != "1" {
		t.Skip("set NIXGO_TEST_FLAKE_REALISE=1 to run a local-store flake realisation test")
	}
	ctx := util.NewContext()
	defer ctx.Close()
	if err := util.Init(ctx); err != nil {
		t.Fatal(err)
	}
	if err := store.Init(ctx); err != nil {
		t.Fatal(err)
	}
	if err := expr.Init(ctx); err != nil {
		t.Fatal(err)
	}
	st, err := store.Open(ctx, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	dir := t.TempDir()
	flakeNix := `{
	  outputs = { self }: {
	    packages.${builtins.currentSystem}.default = builtins.derivation {
	      name = "nixgo-flake-realise-test";
	      system = builtins.currentSystem;
	      builder = "/bin/sh";
	      args = [ "-c" "mkdir -p $out; echo nixgo > $out/result" ];
	    };
	  };
	}`
	if err := os.WriteFile(filepath.Join(dir, "flake.nix"), []byte(flakeNix), 0600); err != nil {
		t.Fatal(err)
	}

	fs, err := fetchers.New(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer fs.Close()
	settings, err := flake.NewSettings(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer settings.Close()
	builder, err := expr.NewBuilder(ctx, st)
	if err != nil {
		t.Fatal(err)
	}
	defer builder.Close()
	if err := builder.LoadEnvironment(); err != nil {
		t.Fatal(err)
	}
	if err := settings.AddToEvalStateBuilder(ctx, builder); err != nil {
		t.Fatal(err)
	}
	state, err := builder.Build()
	if err != nil {
		t.Fatal(err)
	}
	defer state.Close()

	parseFlags, err := flake.NewParseFlags(ctx, settings)
	if err != nil {
		t.Fatal(err)
	}
	defer parseFlags.Close()
	if err := parseFlags.SetBaseDirectory(ctx, dir); err != nil {
		t.Fatal(err)
	}
	ref, _, err := flake.ParseReference(ctx, fs, settings, parseFlags, "path:"+dir)
	if err != nil {
		t.Fatal(err)
	}
	defer ref.Close()

	lockFlags, err := flake.NewLockFlags(ctx, settings)
	if err != nil {
		t.Fatal(err)
	}
	defer lockFlags.Close()
	if err := lockFlags.ModeVirtual(ctx); err != nil {
		t.Fatal(err)
	}
	locked, err := flake.Lock(ctx, fs, settings, state, lockFlags, ref)
	if err != nil {
		t.Fatal(err)
	}
	defer locked.Close()

	outputs, err := locked.RealiseOutput(ctx, settings, state, st, []string{"packages", stateSystem(t, state), "default"})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		for _, p := range outputs {
			p.Close()
		}
	}()
	if outputs["out"] == nil {
		t.Fatalf("expected out output, got %v", outputs)
	}
	if !st.IsValid(outputs["out"]) {
		t.Fatal("realised output is not valid in the store")
	}
}

func stateSystem(t *testing.T, state *expr.State) string {
	t.Helper()
	v, err := state.EvalString("builtins.currentSystem", ".")
	if err != nil {
		t.Fatal(err)
	}
	defer v.Close()
	system, err := v.String()
	if err != nil {
		t.Fatal(err)
	}
	return system
}
