package value

/*
#cgo pkg-config: nix-expr-c nix-store-c nix-util-c
#include <nix_api_value.h>
#include <nix_api_expr.h>
#include <nix_api_util.h>
#include "primop_bridge.h"
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"

	cexpr "github.com/ep0ll/nixlang-go/internal/c/expr"
	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

//export cgoNixPrimOpCB
func cgoNixPrimOpCB(userdata unsafe.Pointer, context *C.nix_c_context, state *C.EvalState, args **C.nix_value, result *C.nix_value) {
	ctx := cutil.WrapBorrowedContext(unsafe.Pointer(context))
	defer func() {
		if recovered := recover(); recovered != nil {
			callbackFailure(ctx, recovered)
		}
	}()
	if userdata == nil {
		callbackFailure(ctx, "missing primop callback handle")
		return
	}
	handle := cgo.Handle(uintptr(C.nixgo_primop_handle_from_ptr(userdata)))
	callback, ok := handle.Value().(*primOpCallback)
	if !ok || callback == nil {
		callbackFailure(ctx, "unknown primop callback handle")
		return
	}

	values := make([]*Value, callback.arity)
	if callback.arity > 0 {
		if args == nil {
			callbackFailure(ctx, "Nix provided a nil primop argument array")
			return
		}
		raw := unsafe.Slice(args, callback.arity)
		for i, arg := range raw {
			values[i] = WrapBorrowedPtr(unsafe.Pointer(arg))
		}
	}
	callback.fn(
		callback.userData,
		ctx,
		cexpr.WrapState(unsafe.Pointer(state)),
		values,
		WrapBorrowedPtr(unsafe.Pointer(result)),
	)
}
