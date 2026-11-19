// Package expr implements advanced evaluator operations exposed by nix_api_expr.h.
package expr

/*
#cgo pkg-config: nix-expr-c nix-store-c nix-util-c
#include <nix_api_expr.h>
#include <nix_api_store.h>
#include <nix_api_value.h>
#include <nix_api_util.h>
#include <stdint.h>
#include "bridge.h"

typedef void (*nixgo_gc_finalizer_callback)(void *, void *);
extern void cgoExprGCFinalizer(void *, void *);
*/
import "C"
import (
	"sync"
	"unsafe"

	cstore "github.com/ep0ll/nixlang-go/internal/c/store"
	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

var (
	gcFinalizerMu      sync.Mutex
	gcFinalizers       = make(map[uintptr]func(unsafe.Pointer))
	gcFinalizerObjects = make(map[uintptr]uintptr)
)

// RegisterGCFinalizer registers a callback invoked when Nix's garbage collector
// collects obj. The callback runs on a Nix GC thread; it must be short and
// should not call back into the Nix evaluator. Registration replaces any
// finalizer already installed for the same object, matching the C API contract.
func RegisterGCFinalizer(obj unsafe.Pointer, fn func(unsafe.Pointer)) {
	if obj == nil || fn == nil {
		return
	}
	userdata := C.nixgo_expr_id_to_userdata(0)
	if userdata == nil {
		return
	}
	tokenID := uintptr(userdata)
	objectID := uintptr(obj)

	// Keep the registry update and C replacement ordered. A concurrently
	// executing old finalizer can then either take the old callback first or
	// observe that its token was superseded, without dereferencing freed data.
	gcFinalizerMu.Lock()
	if previous := gcFinalizerObjects[objectID]; previous != 0 {
		delete(gcFinalizers, previous)
		C.nixgo_expr_userdata_free(unsafe.Pointer(previous))
	}
	gcFinalizers[tokenID] = fn
	gcFinalizerObjects[objectID] = tokenID
	C.nix_gc_register_finalizer(obj, userdata,
		C.nixgo_gc_finalizer_callback(C.cgoExprGCFinalizer))
	gcFinalizerMu.Unlock()
}

// GetDerivation checks whether value is a derivation and returns its .drv path.
// A nil result with a clear context means value was not a derivation.
func GetDerivation(
	ctx *cutil.Context,
	state *EvalState,
	value unsafe.Pointer,
	ignoreAssertionFailures bool,
) *cstore.StorePath {
	p := C.nix_get_derivation(ctx.Ptr(), state.ptr, (*C.nix_value)(value), C.bool(ignoreAssertionFailures))
	if p == nil {
		return nil
	}
	return cstore.WrapPath(unsafe.Pointer(p))
}

// ValueAutoCallFunction invokes a function using named arguments from autoArgs.
// autoArgs may be nil; result must point to an allocated nix_value.
func ValueAutoCallFunction(
	ctx *cutil.Context,
	state *EvalState,
	autoArgs, fn, result unsafe.Pointer,
) cutil.Err {
	return cutil.FromCErr(C.nix_value_auto_call_function(
		ctx.Ptr(), state.ptr,
		(*C.nix_value)(autoArgs), (*C.nix_value)(fn), (*C.nix_value)(result),
	))
}
