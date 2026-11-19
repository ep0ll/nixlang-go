//go:build !cgo

package expr

import (
	"unsafe"

	cstore "github.com/ep0ll/nixlang-go/internal/c/store"
	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

// RegisterGCFinalizer is unavailable when CGO is disabled.
func RegisterGCFinalizer(obj unsafe.Pointer, fn func(unsafe.Pointer)) {}

// GetDerivation is unavailable when CGO is disabled.
func GetDerivation(ctx *cutil.Context, state *EvalState, value unsafe.Pointer, ignoreAssertionFailures bool) *cstore.StorePath {
	return nil
}

// ValueAutoCallFunction is unavailable when CGO is disabled.
func ValueAutoCallFunction(ctx *cutil.Context, state *EvalState, autoArgs, fn, result unsafe.Pointer) cutil.Err {
	return cutil.ErrUnknown
}
