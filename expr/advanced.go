package expr

import (
	"fmt"
	"unsafe"

	cexpr "github.com/ep0ll/nixlang-go/internal/c/expr"
	"github.com/ep0ll/nixlang-go/store"
	"github.com/ep0ll/nixlang-go/util"
	"github.com/ep0ll/nixlang-go/value"
)

// CallAuto calls a function using named arguments from autoArgs. Pass nil to
// supply no named arguments. The result is evaluated before it is returned.
func (s *State) CallAuto(fn, autoArgs *value.Value) (*value.Value, error) {
	if s == nil || s.c == nil || s.ctx == nil || fn == nil ||
		fn.Internal() == nil || fn.Internal().Ptr() == nil || fn.StateInternal() == nil {
		return nil, fmt.Errorf("nix: evaluator state and function are required")
	}
	if !fn.StateInternal().Same(s.c) {
		return nil, fmt.Errorf("nix: function belongs to a different evaluator state")
	}
	var autoPtr unsafe.Pointer
	if autoArgs != nil {
		if autoArgs.Internal() == nil || autoArgs.Internal().Ptr() == nil ||
			autoArgs.StateInternal() == nil || !autoArgs.StateInternal().Same(s.c) {
			return nil, fmt.Errorf("nix: auto arguments must belong to the same evaluator state")
		}
		autoPtr = autoArgs.Internal().Ptr()
	}
	out, err := s.AllocValue()
	if err != nil {
		return nil, err
	}
	code := cexpr.ValueAutoCallFunction(
		s.ctx.Internal(), s.c, autoPtr, fn.Internal().Ptr(), out.Internal().Ptr(),
	)
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

// DerivationPath determines whether value is a derivation and, when it is,
// returns its .drv store path. The bool is false for non-derivation values.
func (s *State) DerivationPath(v *value.Value, ignoreAssertionFailures bool) (*store.Path, bool, error) {
	if s == nil || s.c == nil || s.ctx == nil || v == nil ||
		v.Internal() == nil || v.Internal().Ptr() == nil || v.StateInternal() == nil {
		return nil, false, fmt.Errorf("nix: evaluator state and value are required")
	}
	if !v.StateInternal().Same(s.c) {
		return nil, false, fmt.Errorf("nix: value belongs to a different evaluator state")
	}
	path := cexpr.GetDerivation(s.ctx.Internal(), s.c, v.Internal().Ptr(), ignoreAssertionFailures)
	if path == nil {
		if err := s.ctx.Err(); err != nil {
			return nil, false, err
		}
		return nil, false, nil
	}
	return store.PathFromInternal(path), true, nil
}
