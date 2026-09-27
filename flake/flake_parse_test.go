//go:build cgo && nix

package flake_test

import (
	"testing"

	"github.com/ep0ll/nixlang-go/fetchers"
	"github.com/ep0ll/nixlang-go/flake"
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
