// Package flake provides high-level access to Nix flakes:
// parsing references, locking, and reading output attributes.
//
// Evaluation of flake outputs does not run `nix build` or other CLI tools.
package flake

import (
	"fmt"

	cflake "github.com/ep0ll/nixlang-go/internal/c/flake"
	cvalue "github.com/ep0ll/nixlang-go/internal/c/value"
	"github.com/ep0ll/nixlang-go/expr"
	"github.com/ep0ll/nixlang-go/fetchers"
	"github.com/ep0ll/nixlang-go/util"
	"github.com/ep0ll/nixlang-go/value"
)

// Settings configures flake behavior.
type Settings struct {
	c *cflake.Settings
}

// Reference is a flake reference (how to fetch a flake).
type Reference struct {
	c *cflake.Reference
}

// LockFlags controls locking behavior.
type LockFlags struct {
	c *cflake.LockFlags
}

// ParseFlags controls reference parsing.
type ParseFlags struct {
	c *cflake.ParseFlags
}

// LockedFlake is a flake with a resolved lock.
type LockedFlake struct {
	c   *cflake.LockedFlake
	ctx *util.Context
}

// NewSettings creates default flake settings.
func NewSettings(ctx *util.Context) (*Settings, error) {
	s := cflake.NewSettings(ctx.Internal())
	if s == nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to create flake settings")
	}
	return &Settings{c: s}, nil
}

// Close frees settings.
func (s *Settings) Close() {
	if s != nil && s.c != nil {
		s.c.Free()
		s.c = nil
	}
}

// Internal returns the low-level handle.
func (s *Settings) Internal() *cflake.Settings {
	if s == nil {
		return nil
	}
	return s.c
}

// AddToEvalStateBuilder registers flake-related builtins (e.g. getFlake).
func (s *Settings) AddToEvalStateBuilder(ctx *util.Context, builder *expr.Builder) error {
	return ctx.Check(util.ErrorCode(s.c.AddToEvalStateBuilder(ctx.Internal(), builder.Internal())))
}

// NewParseFlags creates flags for parsing flake references.
func NewParseFlags(ctx *util.Context, settings *Settings) (*ParseFlags, error) {
	f := cflake.NewParseFlags(ctx.Internal(), settings.c)
	if f == nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to create parse flags")
	}
	return &ParseFlags{c: f}, nil
}

// Close frees parse flags.
func (f *ParseFlags) Close() {
	if f != nil && f.c != nil {
		f.c.Free()
		f.c = nil
	}
}

// SetBaseDirectory sets the base directory for relative flake refs.
func (f *ParseFlags) SetBaseDirectory(ctx *util.Context, dir string) error {
	return ctx.Check(util.ErrorCode(f.c.SetBaseDirectory(ctx.Internal(), dir)))
}

// NewLockFlags creates lock flags with defaults from settings.
func NewLockFlags(ctx *util.Context, settings *Settings) (*LockFlags, error) {
	f := cflake.NewLockFlags(ctx.Internal(), settings.c)
	if f == nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to create lock flags")
	}
	return &LockFlags{c: f}, nil
}

// Close frees lock flags.
func (f *LockFlags) Close() {
	if f != nil && f.c != nil {
		f.c.Free()
		f.c = nil
	}
}

// ModeCheck fails if the lock file would need an update.
func (f *LockFlags) ModeCheck(ctx *util.Context) error {
	return ctx.Check(util.ErrorCode(f.c.SetModeCheck(ctx.Internal())))
}

// ModeVirtual updates the lock in memory only.
func (f *LockFlags) ModeVirtual(ctx *util.Context) error {
	return ctx.Check(util.ErrorCode(f.c.SetModeVirtual(ctx.Internal())))
}

// ModeWriteAsNeeded updates the lock file on disk when needed.
func (f *LockFlags) ModeWriteAsNeeded(ctx *util.Context) error {
	return ctx.Check(util.ErrorCode(f.c.SetModeWriteAsNeeded(ctx.Internal())))
}

// AddInputOverride overrides a flake input.
func (f *LockFlags) AddInputOverride(ctx *util.Context, inputPath string, ref *Reference) error {
	return ctx.Check(util.ErrorCode(f.c.AddInputOverride(ctx.Internal(), inputPath, ref.c)))
}

// ParseReference parses a flake reference URI-like string.
func ParseReference(
	ctx *util.Context,
	fetchSettings *fetchers.Settings,
	settings *Settings,
	flags *ParseFlags,
	str string,
) (*Reference, string, error) {
	ref, frag, code := cflake.ParseReference(ctx.Internal(), fetchSettings.Internal(), settings.c, flags.c, str)
	if err := ctx.Check(util.ErrorCode(code)); err != nil {
		return nil, "", err
	}
	if ref == nil {
		return nil, "", fmt.Errorf("nix: failed to parse flake reference %q", str)
	}
	return &Reference{c: ref}, frag, nil
}

// Close frees a flake reference.
func (r *Reference) Close() {
	if r != nil && r.c != nil {
		r.c.Free()
		r.c = nil
	}
}

// Lock resolves and locks a flake reference.
func Lock(
	ctx *util.Context,
	fetchSettings *fetchers.Settings,
	settings *Settings,
	state *expr.State,
	flags *LockFlags,
	ref *Reference,
) (*LockedFlake, error) {
	lf := cflake.Lock(
		ctx.Internal(),
		fetchSettings.Internal(),
		settings.c,
		state.Internal(),
		flags.c,
		ref.c,
	)
	if lf == nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to lock flake")
	}
	return &LockedFlake{c: lf, ctx: ctx}, nil
}

// Close frees the locked flake.
func (lf *LockedFlake) Close() {
	if lf != nil && lf.c != nil {
		lf.c.Free()
		lf.c = nil
	}
}

// OutputAttrs returns the flake's outputs as a Nix attrset value.
// The caller must Close the returned value.
// This evaluates the flake outputs; it does not build packages.
func (lf *LockedFlake) OutputAttrs(ctx *util.Context, settings *Settings, state *expr.State) (*value.Value, error) {
	ptr := lf.c.GetOutputAttrs(ctx.Internal(), settings.c, state.Internal())
	if ptr == nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to get flake output attrs")
	}
	cv := cvalue.WrapPtr(ptr)
	if cv == nil {
		return nil, fmt.Errorf("nix: failed to wrap output attrs value")
	}
	return value.FromInternal(ctx, state.Internal(), cv), nil
}
