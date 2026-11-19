// Package value provides primop registration for the Nix evaluator C API.
package value

/*
#cgo pkg-config: nix-expr-c nix-store-c nix-util-c
#include <nix_api_value.h>
#include <nix_api_expr.h>
#include <nix_api_util.h>
#include <stdlib.h>

typedef void (*nixgo_primop_callback)(
    void *, nix_c_context *, EvalState *, nix_value **, nix_value *);
extern void cgoNixPrimOpCB(void *, nix_c_context *, EvalState *, nix_value **, nix_value *);
*/
import "C"
import (
	"fmt"
	"runtime"
	"sync"
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

var (
	primOpMu        sync.RWMutex
	primOpCallbacks = make(map[uintptr]*primOpCallback)
	primOpOwners    = make(map[*PrimOp]uintptr)
	primOpNext      uintptr = 1
	registeredPrimOpIDs []uintptr
)

func registerPrimOpCallback(callback *primOpCallback) uintptr {
	primOpMu.Lock()
	defer primOpMu.Unlock()
	id := primOpNext
	primOpNext++
	primOpCallbacks[id] = callback
	return id
}

func primOpCallbackFor(id uintptr) *primOpCallback {
	primOpMu.RLock()
	defer primOpMu.RUnlock()
	return primOpCallbacks[id]
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
	if fn == nil || arity < 0 {
		_ = ctx.SetErr(cutil.ErrUnknown, "nix: primop requires a callback and non-negative arity")
		return nil
	}

	callbackID := registerPrimOpCallback(&primOpCallback{fn: fn, userData: userData, arity: arity})
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
		C.int(arity), cname, cargs, cdoc, unsafe.Pointer(callbackID),
	)
	if p == nil {
		primOpMu.Lock()
		delete(primOpCallbacks, callbackID)
		primOpMu.Unlock()
		return nil
	}
	result := &PrimOp{ptr: p}
	primOpMu.Lock()
	primOpOwners[result] = callbackID
	primOpMu.Unlock()
	runtime.SetFinalizer(result, func(p *PrimOp) { p.Free() })
	return result
}

// Register adds this primop to the global builtins registry for subsequently
// created evaluator states. Nix takes ownership of the primop's contents.
func (p *PrimOp) Register(ctx *cutil.Context) cutil.Err {
	if p == nil || p.ptr == nil {
		return cutil.ErrUnknown
	}
	err := cutil.FromCErr(C.nix_register_primop(ctx.Ptr(), p.ptr))
	if err != cutil.OK {
		return err
	}

	// The API explicitly requires releasing the caller's reference after the
	// contents have been moved into the global registry.
	decrefErr := cutil.FromCErr(C.nix_gc_decref(ctx.Ptr(), unsafe.Pointer(p.ptr)))
	p.ptr = nil
	primOpMu.Lock()
	id := primOpOwners[p]
	delete(primOpOwners, p)
	if id != 0 {
		registeredPrimOpIDs = append(registeredPrimOpIDs, id)
	}
	primOpMu.Unlock()
	runtime.SetFinalizer(p, nil)
	return decrefErr
}

// Free releases an unregistered primop. Registered callbacks are retained for
// the process lifetime because the global Nix registry may invoke them later.
func (p *PrimOp) Free() {
	if p == nil || p.ptr == nil {
		return
	}
	p.ptr = nil
	primOpMu.Lock()
	id := primOpOwners[p]
	delete(primOpOwners, p)
	if id != 0 {
		delete(primOpCallbacks, id)
	}
	primOpMu.Unlock()
	C.nix_gc_decref(nil, unsafe.Pointer(p.ptr))
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
