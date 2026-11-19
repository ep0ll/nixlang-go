package util

import cutil "github.com/ep0ll/nixlang-go/internal/c/util"

// WrapInternal exposes a borrowed low-level context to sibling high-level
// packages. Closing the returned context never frees a Nix-owned context.
func WrapInternal(ctx *cutil.Context) *Context {
	if ctx == nil {
		return nil
	}
	return &Context{c: ctx, owned: false}
}
