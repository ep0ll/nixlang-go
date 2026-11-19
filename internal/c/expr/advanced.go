// Package expr implements advanced evaluator operations exposed by nix_api_expr.h.
package expr

/*
#cgo pkg-config: nix-expr-c nix-store-c nix-util-c
#include <nix_api_expr.h>
#include <nix_api_store.h>
#include <nix_api_value.h>
#include <nix_api_util.h>

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
	gcFinalizerMu sync.Mutex
	gcFinalizers = make(map[uintptr]func(unsafe.Pointer))
	gcFinalizerNext uintptr = 1
)

// RegisterGCFinalizer registers a callback invoked when Nix's garbage collector
// collects obj. The callback runs on a Nix GC thread; it must be short and
// should not call back into the Nix evaluator. Registration replaces any
// finalizer already installed for the same object, matching the C API contract.
func RegisterGCFinalizer(obj unsafe.Pointer, fn func(unsafe.Pointer)) {
	if obj == nil || fn == nil {
		return
	}
	gcFinalizerMu.Lock()
	id := gcFinalizerNext
	gcFinalizerNext++
	gcFinalizers[id] = fn
	gcFinalizerMu.Unlock()
	C.nix_gc_register_finalizer(obj, unsafe.Pointer(id),
		C.nixgo_gc_finalizer_callback(C.cgoExprGCFinalizer))
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
