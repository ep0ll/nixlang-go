package value

import (
	"fmt"

	"github.com/ep0ll/nixlang-go/util"
)

// Retain creates an independently-owned reference to this Nix value.
// Use it when keeping a value obtained in a primop callback beyond the
// callback's lifetime.
func (v *Value) Retain() (*Value, error) {
	if err := v.ensureOpen(); err != nil {
		return nil, err
	}

	ctx := v.ctx
	ownsContext := false
	// A borrowed callback context is invalidated when the callback returns.
	// A retained wrapper's private context must also not be shared with another
	// retained value, because either wrapper may be closed independently.
	if ctx.IsBorrowed() || v.ownsContext {
		ctx = util.NewContext()
		ownsContext = true
	}

	retained, code := v.c.Retain(ctx.Internal())
	if err := ctx.Check(util.ErrorCode(code)); err != nil {
		if ownsContext {
			ctx.Close()
		}
		return nil, err
	}
	if retained == nil {
		if ownsContext {
			ctx.Close()
		}
		return nil, fmt.Errorf("nix: failed to retain value")
	}
	return &Value{c: retained, state: v.state, ctx: ctx, ownsContext: ownsContext}, nil
}
