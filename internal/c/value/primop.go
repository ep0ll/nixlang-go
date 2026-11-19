// Package value provides primop registration for the Nix evaluator C API.
package value

/*
#cgo pkg-config: nix-expr-c nix-store-c nix-util-c
#include <nix_api_value.h>
#include <nix_api_expr.h>
#include <nix_api_util.h>
#include <stdlib.h>
#include "primop_bridge.h"

typedef void (*nixgo_primop_callback)(
    void *, nix_c_context *, EvalState *, nix_value **, nix_value *);
extern void cgoNixPrimOpCB(void *, nix_c_context *, EvalState *, nix_value **, nix_value *);
*/
import "C"
import (
	"fmt"
	"runtime"
	"runtime/cgo"
	"unsafe"

	cexpr "github.com/ep0ll/nixlang-go/internal/c/expr"
	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

// PrimOpFunc implements a Nix primop. Arguments and result are borrowed values
// valid only for the duration of the callback. Use Value.Retain to keep an
// argument beyond the callback. Report failures with ctx.SetErr.
type PrimOpFunc func(userData any, ctx *cutil.Context, state *cexpr.EvalState, args []*Value, result *Value)

type primOpCallback struct {
	fn       PrimOpFunc
	userData any
	arity    int
}

// AllocPrimOp creates a C API primop. The returned primop is not visible to
// evaluator states until Register succeeds. Call Free if it is not registered.
func AllocPrimOp(
	ctx *cutil.Context,
	fn PrimOpFunc,
	arity int,
	name string,
	argNames []string,
	doc string,
	userData any,
) *PrimOp {
	if fn == nil || arity < 0 || arity > int(^uint32(0)>>1) || len(argNames) != arity || name == "" {
		if ctx != nil {
			_ = ctx.SetErr(cutil.ErrUnknown, "nix: primop requires a name, callback, non-negative arity and one argument name per parameter")
		}
		return nil
	}

	handle := cgo.NewHandle(&primOpCallback{fn: fn, userData: userData, arity: arity})
	self := C.nixgo_primop_handle_to_ptr(C.uintptr_t(handle))
	if self == nil {
		handle.Delete()
		return nil
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var cdoc *C.char
	if doc != "" {
		cdoc = C.CString(doc)
		defer C.free(unsafe.Pointer(cdoc))
	}

	var cargs **C.char
	if len(argNames) > 0 {
		values := make([]*C.char, len(argNames)+1)
		for i, name := range argNames {
			values[i] = C.CString(name)
		}
		defer func() {
			for _, value := range values {
				if value != nil {
					C.free(unsafe.Pointer(value))
				}
			}
		}()
		cargs = (**C.char)(unsafe.Pointer(&values[0]))
	}

	p := C.nix_alloc_primop(
		ctx.Ptr(), C.nixgo_primop_callback(C.cgoNixPrimOpCB),
		C.int(arity), cname, cargs, cdoc, self,
	)
	if p == nil {
		C.nixgo_primop_handle_free(self)
		handle.Delete()
		return nil
	}
	result := &PrimOp{ptr: p, handle: uintptr(handle), userdata: self}
	runtime.SetFinalizer(result, func(p *PrimOp) { p.Free() })
	return result
}

// Register adds this primop to the global builtins registry for subsequently
// created evaluator states. Nix takes ownership of the primop's contents.
func (p *PrimOp) Register(ctx *cutil.Context) cutil.Err {
	if p == nil || p.ptr == nil || p.retained {
		return cutil.ErrUnknown
	}
	err := cutil.FromCErr(C.nix_register_primop(ctx.Ptr(), p.ptr))
	if err != cutil.OK {
		return err
	}

	// The API explicitly requires releasing the caller's reference after the
	// contents have been moved into the global registry. The callback handle
	// remains live for the process lifetime because the global registry owns it.
	decrefErr := cutil.FromCErr(C.nix_gc_decref(ctx.Ptr(), unsafe.Pointer(p.ptr)))
	p.ptr = nil
	// The C API moves the primop into its process-wide registry. The callback
	// handle must remain live for as long as that registry can call it.
	p.handle = 0
	p.userdata = nil // The process-wide registry now owns the callback handle.
	runtime.SetFinalizer(p, nil)
	return decrefErr
}

// Free releases an unregistered primop and its callback handle. Registered
// callbacks are retained for the process lifetime because the global Nix
// registry may invoke them later.
func (p *PrimOp) Free() {
	if p == nil {
		return
	}
	ptr := p.ptr
	if ptr == nil {
		runtime.SetFinalizer(p, nil)
		return
	}
	p.ptr = nil
	handle := cgo.Handle(p.handle)
	userdata := p.userdata
	p.handle = 0
	p.userdata = nil

	// A retained primop has already transferred its callback handle to Nix
	// and deliberately cleared both fields. An unretained primop cannot be
	// invoked after this wrapper is released, so its handle can be deleted now.
	C.nix_gc_decref(nil, unsafe.Pointer(ptr))
	if handle != 0 {
		handle.Delete()
	}
	if userdata != nil {
		C.nixgo_primop_handle_free(userdata)
	}
	runtime.SetFinalizer(p, nil)
}

// Ptr returns the underlying C primop pointer.
func (p *PrimOp) Ptr() unsafe.Pointer {
	if p == nil {
		return nil
	}
	return unsafe.Pointer(p.ptr)
}

// callbackFailure converts a Go panic into a Nix callback error.
func callbackFailure(ctx *cutil.Context, recovered any) {
	if recovered != nil && ctx != nil {
		_ = ctx.SetErr(cutil.ErrUnknown, fmt.Sprintf("Go primop callback panic: %v", recovered))
	}
}
