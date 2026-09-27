// Package expr provides high-level evaluation of Nix expressions.
//
// It never invokes the Nix CLI (no `nix build`, `nix develop`, etc.).
// Evaluation is performed entirely through the C API.
package expr

import (
	"fmt"
	"os"
	"unsafe"

	cexpr "github.com/ep0ll/nixlang-go/internal/c/expr"
	cvalue "github.com/ep0ll/nixlang-go/internal/c/value"
	"github.com/ep0ll/nixlang-go/store"
	"github.com/ep0ll/nixlang-go/util"
	"github.com/ep0ll/nixlang-go/value"
)

// State is a Nix language evaluator instance.
type State struct {
	c   *cexpr.EvalState
	ctx *util.Context
	st  *store.Store
}

// Builder constructs an EvalState with custom settings.
type Builder struct {
	c   *cexpr.Builder
	ctx *util.Context
	st  *store.Store
}

// Init initializes libexpr. Call after util.Init and store.Init.
func Init(ctx *util.Context) error {
	return ctx.Check(util.ErrorCode(cexpr.Init(ctx.Internal())))
}

// NewBuilder creates a builder for a customized evaluator.
func NewBuilder(ctx *util.Context, st *store.Store) (*Builder, error) {
	b := cexpr.NewBuilder(ctx.Internal(), st.Internal())
	if b == nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to create eval state builder")
	}
	return &Builder{c: b, ctx: ctx, st: st}, nil
}

// LoadEnvironment loads settings from the ambient environment and config files.
func (b *Builder) LoadEnvironment() error {
	return b.ctx.Check(util.ErrorCode(b.c.Load(b.ctx.Internal())))
}

// SetLookupPath sets paths used for <...> lookups (like NIX_PATH).
func (b *Builder) SetLookupPath(paths []string) error {
	return b.ctx.Check(util.ErrorCode(b.c.SetLookupPath(b.ctx.Internal(), paths)))
}

// SetSetting sets an evaluator setting (e.g. "pure-eval").
func (b *Builder) SetSetting(key, value string) error {
	return b.ctx.Check(util.ErrorCode(b.c.SetSetting(b.ctx.Internal(), key, value)))
}

// Internal returns the low-level builder (for flake integration).
func (b *Builder) Internal() *cexpr.Builder {
	return b.c
}

// Build creates the EvalState. Call Close on the builder after Build.
func (b *Builder) Build() (*State, error) {
	s := b.c.Build(b.ctx.Internal())
	if s == nil {
		if err := b.ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to build eval state")
	}
	return &State{c: s, ctx: b.ctx, st: b.st}, nil
}

// Close frees the builder.
func (b *Builder) Close() {
	if b != nil && b.c != nil {
		b.c.Free()
		b.c = nil
	}
}

// NewState creates an evaluator with optional lookup path.
func NewState(ctx *util.Context, st *store.Store, lookupPath []string) (*State, error) {
	s := cexpr.Create(ctx.Internal(), lookupPath, st.Internal())
	if s == nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to create eval state")
	}
	return &State{c: s, ctx: ctx, st: st}, nil
}

// Close frees the evaluator state.
func (s *State) Close() {
	if s != nil && s.c != nil {
		s.c.Free()
		s.c = nil
	}
}

// Internal returns the low-level EvalState.
func (s *State) Internal() *cexpr.EvalState {
	if s == nil {
		return nil
	}
	return s.c
}

// AllocValue allocates an uninitialized value bound to this state.
func (s *State) AllocValue() (*value.Value, error) {
	cv := cvalue.Alloc(s.ctx.Internal(), s.c)
	if cv == nil {
		if err := s.ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to allocate value")
	}
	return value.FromInternal(s.ctx, s.c, cv), nil
}

// EvalString parses and evaluates a Nix expression string.
// path is the logical file path used to resolve relative paths (use "." if N/A).
// This only evaluates; it does not build or run any Nix CLI command.
func (s *State) EvalString(exprStr, path string) (*value.Value, error) {
	v, err := s.AllocValue()
	if err != nil {
		return nil, err
	}
	code := cexpr.EvalFromString(s.ctx.Internal(), s.c, exprStr, path, v.Internal().Ptr())
	if err := s.ctx.Check(util.ErrorCode(code)); err != nil {
		v.Close()
		return nil, err
	}
	if err := v.Force(); err != nil {
		v.Close()
		return nil, err
	}
	return v, nil
}

// EvalFile reads path and evaluates its contents as a Nix expression.
// Relative imports resolve against the file's directory.
// Evaluation only — does not run nix build/develop or any CLI command.
func (s *State) EvalFile(path string) (*value.Value, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("nix: read %s: %w", path, err)
	}
	return s.EvalString(string(data), path)
}

// Call applies a Nix function to a single argument.
func (s *State) Call(fn, arg *value.Value) (*value.Value, error) {
	out, err := s.AllocValue()
	if err != nil {
		return nil, err
	}
	code := cexpr.ValueCall(s.ctx.Internal(), s.c, fn.Internal().Ptr(), arg.Internal().Ptr(), out.Internal().Ptr())
	if err := s.ctx.Check(util.ErrorCode(code)); err != nil {
		out.Close()
		return nil, err
	}
	if err := out.Force(); err != nil {
		out.Close()
		return nil, err
	}
	return out, nil
}

// CallMulti applies a curried function to multiple arguments.
func (s *State) CallMulti(fn *value.Value, args ...*value.Value) (*value.Value, error) {
	out, err := s.AllocValue()
	if err != nil {
		return nil, err
	}
	argPtrs := make([]unsafe.Pointer, len(args))
	for i, a := range args {
		argPtrs[i] = a.Internal().Ptr()
	}
	code := cexpr.ValueCallMulti(s.ctx.Internal(), s.c, fn.Internal().Ptr(), argPtrs, out.Internal().Ptr())
	if err := s.ctx.Check(util.ErrorCode(code)); err != nil {
		out.Close()
		return nil, err
	}
	if err := out.Force(); err != nil {
		out.Close()
		return nil, err
	}
	return out, nil
}

// Store returns the store associated with this state.
func (s *State) Store() *store.Store {
	return s.st
}
