// Package value — extended type/copy/string-realise helpers.
package value

import (
	"fmt"

	cvalue "github.com/ep0ll/nixlang-go/internal/c/value"
	"github.com/ep0ll/nixlang-go/util"
)

// TypeName returns the Nix type name (e.g. "int", "attrs", "string").
func (v *Value) TypeName() string {
	return v.c.GetTypename(v.ctx.Internal())
}

// CopyFrom copies source into v (both must be allocated on the same state).
func (v *Value) CopyFrom(source *Value) error {
	code := v.c.CopyValue(v.ctx.Internal(), source.c)
	return v.ctx.Check(util.ErrorCode(code))
}

// RealiseString forces a string and returns its content plus store-path context count.
// The returned buffer is the realised string; pathCount is the number of context paths.
func (v *Value) RealiseString(allowIFD bool) (buf string, pathCount int, err error) {
	rs := cvalue.StringRealise(v.ctx.Internal(), v.state, v.c, allowIFD)
	if rs == nil {
		if e := v.ctx.Err(); e != nil {
			return "", 0, e
		}
		return "", 0, fmt.Errorf("nix: string realise failed")
	}
	defer rs.Free()
	return rs.Buffer(), rs.StorePathCount(), nil
}
