// Package value provides typed access to Nix values and custom primops.
package value

import (
	"errors"
	"fmt"

	cexpr "github.com/ep0ll/nixlang-go/internal/c/expr"
	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
	cvalue "github.com/ep0ll/nixlang-go/internal/c/value"
	"github.com/ep0ll/nixlang-go/util"
)

// PrimOpFunc implements a Nix builtin. The state pointer and argument values
// are borrowed for the duration of the callback. Call Retain on any argument
// that must outlive the callback. Return an error to raise a Nix evaluation
// error; a *util.NixError preserves its error code.
type PrimOpFunc func(
	userData any,
	ctx *util.Context,
	state *PrimOpState,
	args []*Value,
	result *Value,
) error

// PrimOp is a custom Nix builtin that can be registered for subsequently
// created evaluator states.
type PrimOp struct {
	c *cvalue.PrimOp
}

// NewPrimOp allocates a custom builtin. Register it before creating evaluator
// states that should see it. Close an unregistered primop to release it.
func NewPrimOp(
	ctx *util.Context,
	fn PrimOpFunc,
	arity int,
	name string,
	argNames []string,
	doc string,
	userData any,
) (*PrimOp, error) {
	if ctx == nil || ctx.Internal() == nil {
		return nil, fmt.Errorf("nix: context is required")
	}
	if fn == nil {
		return nil, fmt.Errorf("nix: primop callback is required")
	}
	callback := func(
		data any,
		callbackContext *cutil.Context,
		state *cexpr.EvalState,
		args []*cvalue.Value,
		result *cvalue.Value,
	) {
		wrappedContext := util.WrapInternal(callbackContext)
		wrappedArgs := make([]*Value, len(args))
		for i, arg := range args {
			wrappedArgs[i] = FromInternal(wrappedContext, state, arg)
		}
		wrappedResult := FromInternal(wrappedContext, state, result)
		callbackState := &PrimOpState{ctx: wrappedContext, state: state}
		defer callbackState.Close()
		if err := fn(data, wrappedContext, callbackState, wrappedArgs, wrappedResult); err != nil {
			code := cutil.ErrUnknown
			var nixErr *util.NixError
			if errors.As(err, &nixErr) {
				code = cutil.Err(nixErr.Code)
			}
			_ = callbackContext.SetErr(code, err.Error())
		}
	}
	primop := cvalue.AllocPrimOp(
		ctx.Internal(), callback, arity, name, argNames, doc, userData,
	)
	if primop == nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to allocate primop %q", name)
	}
	return &PrimOp{c: primop}, nil
}

// Register adds the builtin to the process-wide registry. It affects only
// evaluator states created after registration.
func (p *PrimOp) Register(ctx *util.Context) error {
	if p == nil || p.c == nil {
		return fmt.Errorf("nix: primop is closed")
	}
	if ctx == nil || ctx.Internal() == nil {
		return fmt.Errorf("nix: context is required")
	}
	return ctx.Check(util.ErrorCode(p.c.Register(ctx.Internal())))
}

// Close releases an unregistered primop. It is safe after Register.
func (p *PrimOp) Close() {
	if p != nil && p.c != nil {
		p.c.Free()
		p.c = nil
	}
}

// RegisterPrimOp allocates and immediately registers a custom builtin.
func RegisterPrimOp(
	ctx *util.Context,
	fn PrimOpFunc,
	arity int,
	name string,
	argNames []string,
	doc string,
	userData any,
) error {
	primop, err := NewPrimOp(ctx, fn, arity, name, argNames, doc, userData)
	if err != nil {
		return err
	}
	defer primop.Close()
	return primop.Register(ctx)
}


// SetPrimOp initializes this value as a custom primop function. This is
// separate from RegisterPrimOp: registration exposes a builtin through the
// global builtins attrset for subsequently created evaluator states.
func (v *Value) SetPrimOp(primop *PrimOp) error {
	if err := v.ensureOpen(); err != nil {
		return err
	}
	if primop == nil || primop.c == nil {
		return fmt.Errorf("nix: primop is required")
	}
	return v.ctx.Check(util.ErrorCode(v.c.InitPrimOp(v.ctx.Internal(), primop.c)))
}
