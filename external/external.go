// Package external provides support for foreign (external) Nix values.
//
// Use this to expose Go-managed data to the Nix evaluator via primops
// or custom value types.
package external

import (
	"fmt"

	cext "github.com/ep0ll/nixlang-go/internal/c/external"
	"github.com/ep0ll/nixlang-go/util"
)

// Desc describes callbacks for an external value class.
// Keep the Desc alive for as long as any value created with it exists.
type Desc = cext.Desc

// Value is a GC-managed external value.
type Value struct {
	c *cext.Value
}

// New creates an external value with the given descriptor and opaque pointer.
func New(ctx *util.Context, desc *Desc, data any) (*Value, error) {
	v := cext.Create(ctx.Internal(), desc, data)
	if v == nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to create external value")
	}
	return &Value{c: v}, nil
}

// Close releases the GC reference.
func (v *Value) Close() {
	if v != nil && v.c != nil {
		v.c.Free()
		v.c = nil
	}
}

// Content returns the opaque data associated with the external value.
func (v *Value) Content(ctx *util.Context) any {
	return v.c.Content(ctx.Internal())
}
