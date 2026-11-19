package value

/*
#cgo pkg-config: nix-expr-c nix-store-c nix-util-c
#include <nix_api_value.h>
#include <nix_api_expr.h>
#include <nix_api_util.h>
*/
import "C"
import (
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
	callback := primOpCallbackFor(uintptr(userdata))
	if callback == nil {
		callbackFailure(ctx, "unknown or released primop callback handle")
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
