//go:build !cgo

package util

// ErrorInfoMsg is unavailable when CGO is disabled.
func ErrorInfoMsg(ctx, readContext *Context) (string, Err) { return "", ErrUnknown }

// ErrorName is unavailable when CGO is disabled.
func ErrorName(ctx, readContext *Context) (string, Err) { return "", ErrUnknown }
