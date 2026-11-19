// Package flake provides high-level access to Nix flakes:
// parsing references, locking, and reading output attributes.
//
// Evaluation of flake outputs does not run `nix build` or other CLI tools.
package flake

import (
	"fmt"

	"github.com/ep0ll/nixlang-go/expr"
	"github.com/ep0ll/nixlang-go/fetchers"
	cflake "github.com/ep0ll/nixlang-go/internal/c/flake"
	cvalue "github.com/ep0ll/nixlang-go/internal/c/value"
	"github.com/ep0ll/nixlang-go/store"
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
	if ctx == nil || ctx.Internal() == nil {
		return nil, fmt.Errorf("nix: context is required")
	}
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
	if s == nil || s.c == nil || ctx == nil || ctx.Internal() == nil || builder == nil || builder.Internal() == nil {
		return fmt.Errorf("nix: flake settings, context and evaluator builder are required")
	}
	return ctx.Check(util.ErrorCode(s.c.AddToEvalStateBuilder(ctx.Internal(), builder.Internal())))
}

// NewParseFlags creates flags for parsing flake references.
func NewParseFlags(ctx *util.Context, settings *Settings) (*ParseFlags, error) {
	if ctx == nil || ctx.Internal() == nil || settings == nil || settings.c == nil {
		return nil, fmt.Errorf("nix: context and flake settings are required")
	}
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
	if f == nil || f.c == nil || ctx == nil || ctx.Internal() == nil {
		return fmt.Errorf("nix: parse flags and context are required")
	}
	return ctx.Check(util.ErrorCode(f.c.SetBaseDirectory(ctx.Internal(), dir)))
}

// NewLockFlags creates lock flags with defaults from settings.
func NewLockFlags(ctx *util.Context, settings *Settings) (*LockFlags, error) {
	if ctx == nil || settings == nil || settings.c == nil {
		return nil, fmt.Errorf("nix: context and flake settings are required")
	}
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
	if f == nil || f.c == nil || ctx == nil {
		return fmt.Errorf("nix: lock flags and context are required")
	}
	return ctx.Check(util.ErrorCode(f.c.SetModeCheck(ctx.Internal())))
}

// ModeVirtual updates the lock in memory only.
func (f *LockFlags) ModeVirtual(ctx *util.Context) error {
	if f == nil || f.c == nil || ctx == nil {
		return fmt.Errorf("nix: lock flags and context are required")
	}
	return ctx.Check(util.ErrorCode(f.c.SetModeVirtual(ctx.Internal())))
}

// ModeWriteAsNeeded updates the lock file on disk when needed.
func (f *LockFlags) ModeWriteAsNeeded(ctx *util.Context) error {
	if f == nil || f.c == nil || ctx == nil {
		return fmt.Errorf("nix: lock flags and context are required")
	}
	return ctx.Check(util.ErrorCode(f.c.SetModeWriteAsNeeded(ctx.Internal())))
}

// AddInputOverride overrides a flake input.
func (f *LockFlags) AddInputOverride(ctx *util.Context, inputPath string, ref *Reference) error {
	if f == nil || f.c == nil || ref == nil || ref.c == nil || ctx == nil || ctx.Internal() == nil {
		return fmt.Errorf("nix: lock flags, context and input reference are required")
	}
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
	if ctx == nil || ctx.Internal() == nil || fetchSettings == nil || fetchSettings.Internal() == nil ||
		settings == nil || settings.c == nil || flags == nil || flags.c == nil {
		return nil, "", fmt.Errorf("nix: fetcher settings, flake settings, parse flags and context are required")
	}
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
	if ctx == nil || fetchSettings == nil || fetchSettings.Internal() == nil ||
		settings == nil || settings.c == nil || state == nil || state.Internal() == nil ||
		flags == nil || flags.c == nil || ref == nil || ref.c == nil {
		return nil, fmt.Errorf("nix: context, fetcher settings, flake settings, state, lock flags and reference are required")
	}
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
	if lf == nil || lf.c == nil || ctx == nil || ctx.Internal() == nil || settings == nil || settings.c == nil || state == nil || state.Internal() == nil {
		return nil, fmt.Errorf("nix: locked flake, context, settings and evaluator state are required")
	}
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

// RealiseOutput evaluates a selected flake output attribute path, verifies that
// it is a derivation, and realises its derivation path through the supplied
// store. The returned output paths are owned by the caller and must be closed.
// Attribute paths are individual Nix attribute names, not dot-separated text;
// this permits names that contain dots.
func (lf *LockedFlake) RealiseOutput(
	ctx *util.Context,
	settings *Settings,
	state *expr.State,
	st *store.Store,
	attrPath []string,
) (map[string]*store.Path, error) {
	if lf == nil || lf.c == nil || ctx == nil || ctx.Internal() == nil ||
		settings == nil || settings.c == nil || state == nil || state.Internal() == nil ||
		st == nil || st.Internal() == nil {
		return nil, fmt.Errorf("nix: locked flake, context, settings, evaluator state and store are required")
	}
	if len(attrPath) == 0 {
		return nil, fmt.Errorf("nix: flake output attribute path must not be empty")
	}
	for i, name := range attrPath {
		if name == "" {
			return nil, fmt.Errorf("nix: flake output attribute %d is empty", i)
		}
	}

	outputs, err := lf.OutputAttrs(ctx, settings, state)
	if err != nil {
		return nil, err
	}
	defer outputs.Close()

	selected := outputs
	owned := make([]*value.Value, 0, len(attrPath))
	defer func() {
		for _, v := range owned {
			v.Close()
		}
	}()
	for _, name := range attrPath {
		next, attrErr := selected.Attr(name)
		if attrErr != nil {
			return nil, fmt.Errorf("nix: select flake output %q: %w", name, attrErr)
		}
		owned = append(owned, next)
		selected = next
	}

	drvPath, isDerivation, err := state.DerivationPath(selected, false)
	if err != nil {
		return nil, fmt.Errorf("nix: inspect flake output %v: %w", attrPath, err)
	}
	if !isDerivation || drvPath == nil {
		return nil, fmt.Errorf("nix: flake output %v is not a derivation", attrPath)
	}
	defer drvPath.Close()

	realised := make(map[string]*store.Path)
	var cloneErr error
	err = st.Realise(drvPath, func(name string, path *store.Path) {
		if path == nil {
			return
		}
		clone := path.Clone()
		if clone == nil {
			cloneErr = fmt.Errorf("nix: failed to clone realised output %q", name)
			return
		}
		if previous := realised[name]; previous != nil {
			previous.Close()
		}
		realised[name] = clone
	})
	if err == nil && cloneErr != nil {
		err = cloneErr
	}
	if err != nil {
		for _, path := range realised {
			path.Close()
		}
		return nil, fmt.Errorf("nix: realise flake output %v: %w", attrPath, err)
	}
	if len(realised) == 0 {
		return nil, fmt.Errorf("nix: realisation of flake output %v returned no outputs", attrPath)
	}
	return realised, nil
}
