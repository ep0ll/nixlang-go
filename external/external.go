// Package external exposes Nix external values to Go callbacks.
package external

import (
	"fmt"
	"unsafe"

	cext "github.com/ep0ll/nixlang-go/internal/c/external"
	"github.com/ep0ll/nixlang-go/util"
)

// Callbacks describes the behavior of a Nix external-value class.
type Callbacks = cext.Callbacks

// Desc describes the callback class used to create external values.
type Desc = cext.Desc

// Value is a GC-managed external value.
type Value struct {
	c *cext.Value
}

// NewDesc creates a callback descriptor. Call Close after the descriptor is
// no longer needed; active values keep its C descriptor alive until Nix GC.
func NewDesc(callbacks Callbacks) *Desc {
	return cext.NewDesc(callbacks)
}

// New creates an external value from Go-managed data.
func New(ctx *util.Context, desc *Desc, data any) (*Value, error) {
	if ctx == nil || ctx.Internal() == nil || desc == nil {
		return nil, fmt.Errorf("nix: context and external descriptor are required")
	}
	v := cext.Create(ctx.Internal(), desc, data)
	if v == nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to create external value")
	}
	return &Value{c: v}, nil
}

// FromPointer wraps a borrowed C ExternalValue pointer with an additional Nix GC reference.
func FromPointer(ctx *util.Context, ptr unsafe.Pointer) (*Value, error) {
	if ctx == nil || ctx.Internal() == nil || ptr == nil {
		return nil, fmt.Errorf("nix: context and external value pointer are required")
	}
	v := cext.Wrap(ctx.Internal(), ptr)
	if v == nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to retain external value")
	}
	return &Value{c: v}, nil
}

// Close releases this Go reference to the external value.
func (v *Value) Close() {
	if v != nil && v.c != nil {
		v.c.Free()
		v.c = nil
	}
}

// Content returns the opaque Go data associated with the external value.
func (v *Value) Content(ctx *util.Context) any {
	if v == nil || v.c == nil || ctx == nil || ctx.Internal() == nil {
		return nil
	}
	return v.c.Content(ctx.Internal())
}

// Pointer returns the underlying C ExternalValue pointer for value initialization.
func (v *Value) Pointer() unsafe.Pointer {
	if v == nil || v.c == nil {
		return nil
	}
	return v.c.Ptr()
}
