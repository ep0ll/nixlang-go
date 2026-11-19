package value

import (
	"fmt"

	"github.com/ep0ll/nixlang-go/external"
	"github.com/ep0ll/nixlang-go/util"
)

// SetExternal initializes this value with a Nix C API external value.
func (v *Value) SetExternal(externalValue *external.Value) error {
	if err := v.ensureOpen(); err != nil {
		return err
	}
	if externalValue == nil {
		return fmt.Errorf("nix: external value is required")
	}
	return v.ctx.Check(util.ErrorCode(v.c.InitExternal(v.ctx.Internal(), externalValue.Pointer())))
}

// External returns an independently-referenced external value when the Nix
// value contains a C API external value. The caller must Close the result.
func (v *Value) External() (*external.Value, error) {
	if v == nil || v.c == nil {
		return nil, fmt.Errorf("nix: value is closed")
	}
	if err := v.ensureType(TypeExternal); err != nil {
		return nil, err
	}
	ptr := v.c.GetExternal(v.ctx.Internal())
	if ptr == nil {
		if err := v.ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: value is not a Nix C API external value")
	}
	return external.FromPointer(v.ctx, ptr)
}
