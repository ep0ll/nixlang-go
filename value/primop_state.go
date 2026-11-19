package value

import (
	"fmt"
	"unsafe"

	cexpr "github.com/ep0ll/nixlang-go/internal/c/expr"
	cvalue "github.com/ep0ll/nixlang-go/internal/c/value"
	"github.com/ep0ll/nixlang-go/util"
)

// PrimOpState is a callback-scoped facade over the evaluator that invoked a
// custom primop. Its Context is borrowed and becomes invalid when the callback
// returns. Values allocated through it are tied to the same evaluator arena;
// keep the parent expr.State open while using those values.
type PrimOpState struct {
	ctx   *util.Context
	state *cexpr.EvalState
}

func (s *PrimOpState) ensureOpen() error {
	if s == nil || s.ctx == nil || s.ctx.Internal() == nil || s.state == nil || s.state.Ptr() == nil {
		return fmt.Errorf("nix: primop evaluator is no longer available")
	}
	return nil
}

// Close drops this facade's references. It does not free the evaluator state,
// which is owned by the caller's expr.State and may still be referenced by
// values retained from this callback.
func (s *PrimOpState) Close() {
	if s == nil {
		return
	}
	s.ctx = nil
	s.state = nil
}

// AllocValue allocates an uninitialized value in the callback evaluator.
func (s *PrimOpState) AllocValue() (*Value, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	ctx := util.NewContext()
	v := cvalue.Alloc(ctx.Internal(), s.state)
	if v == nil {
		err := ctx.Err()
		ctx.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to allocate primop value")
	}
	return &Value{c: v, state: s.state, ctx: ctx, ownsContext: true}, nil
}

// NewList constructs a list in the callback evaluator.
func (s *PrimOpState) NewList(elems ...*Value) (*Value, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	ctx := util.NewContext()
	v, err := NewList(ctx, s.state, elems)
	if err != nil {
		ctx.Close()
		return nil, err
	}
	v.ownsContext = true
	return v, nil
}

// NewAttrs constructs an attribute set in the callback evaluator.
func (s *PrimOpState) NewAttrs(attrs map[string]*Value) (*Value, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	ctx := util.NewContext()
	v, err := NewAttrs(ctx, s.state, attrs)
	if err != nil {
		ctx.Close()
		return nil, err
	}
	v.ownsContext = true
	return v, nil
}

// EvalString parses and evaluates an expression in the callback evaluator.
func (s *PrimOpState) EvalString(expression, path string) (*Value, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	v, err := s.AllocValue()
	if err != nil {
		return nil, err
	}
	code := cexpr.EvalFromString(s.ctx.Internal(), s.state, expression, path, v.c.Ptr())
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

// Call applies a Nix function to one argument in the callback evaluator.
func (s *PrimOpState) Call(fn, arg *Value) (*Value, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	if fn == nil || fn.c == nil || fn.c.Ptr() == nil || fn.state == nil || !fn.state.Same(s.state) {
		return nil, fmt.Errorf("nix: function must belong to the callback evaluator")
	}
	if arg == nil || arg.c == nil || arg.c.Ptr() == nil || arg.state == nil || !arg.state.Same(s.state) {
		return nil, fmt.Errorf("nix: argument must belong to the callback evaluator")
	}
	out, err := s.AllocValue()
	if err != nil {
		return nil, err
	}
	code := cexpr.ValueCall(s.ctx.Internal(), s.state, fn.c.Ptr(), arg.c.Ptr(), out.c.Ptr())
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

// CallMulti applies a curried Nix function to multiple arguments.
func (s *PrimOpState) CallMulti(fn *Value, args ...*Value) (*Value, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	if fn == nil || fn.c == nil || fn.c.Ptr() == nil || fn.state == nil || !fn.state.Same(s.state) {
		return nil, fmt.Errorf("nix: function must belong to the callback evaluator")
	}
	argPointers := make([]unsafe.Pointer, len(args))
	for i, arg := range args {
		if arg == nil || arg.c == nil || arg.c.Ptr() == nil || arg.state == nil || !arg.state.Same(s.state) {
			return nil, fmt.Errorf("nix: argument %d must belong to the callback evaluator", i)
		}
		argPointers[i] = arg.c.Ptr()
	}
	out, err := s.AllocValue()
	if err != nil {
		return nil, err
	}
	code := cexpr.ValueCallMulti(s.ctx.Internal(), s.state, fn.c.Ptr(), argPointers, out.c.Ptr())
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
