package util

import cutil "github.com/ep0ll/nixlang-go/internal/c/util"

// ErrorInfoMessage returns the detailed error message attached to the current
// context. It is useful after an API call returns ErrNixError.
func (ctx *Context) ErrorInfoMessage() (string, error) {
	if ctx == nil || ctx.c == nil {
		return "", nil
	}
	message, code := cutil.ErrorInfoMsg(ctx.c, ctx.c)
	if err := ctx.Check(ErrorCode(code)); err != nil {
		return "", err
	}
	return message, nil
}

// ErrorName returns the structured Nix error name attached to the current
// context. It is useful after an API call returns ErrNixError.
func (ctx *Context) ErrorName() (string, error) {
	if ctx == nil || ctx.c == nil {
		return "", nil
	}
	name, code := cutil.ErrorName(ctx.c, ctx.c)
	if err := ctx.Check(ErrorCode(code)); err != nil {
		return "", err
	}
	return name, nil
}
