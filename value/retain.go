package value

import (
	"fmt"

	"github.com/ep0ll/nixlang-go/util"
)

// Retain creates an independently-owned reference to this Nix value.
// Use it when keeping a value obtained in a primop callback beyond the
// callback's lifetime.
func (v *Value) Retain() (*Value, error) {
	if v == nil || v.c == nil {
		return nil, fmt.Errorf("nix: value is closed")
	}
	retained, code := v.c.Retain(v.ctx.Internal())
	if err := v.ctx.Check(util.ErrorCode(code)); err != nil {
		return nil, err
	}
	if retained == nil {
		return nil, fmt.Errorf("nix: failed to retain value")
	}
	return FromInternal(v.ctx, v.state, retained), nil
}
