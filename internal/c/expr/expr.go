// Package expr provides low-level cgo bindings for nix_api_expr.h.
package expr

/*
#cgo pkg-config: nix-expr-c nix-store-c nix-util-c
#include <nix_api_expr.h>
#include <nix_api_store.h>
#include <nix_api_util.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"unsafe"

	cstore "github.com/ep0ll/nixlang-go/internal/c/store"
	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

// EvalState wraps the Nix evaluator state.
type EvalState struct {
	ptr *C.EvalState
}

// Builder wraps nix_eval_state_builder.
type Builder struct {
	ptr *C.nix_eval_state_builder
}

// Init initializes libexpr.
func Init(ctx *cutil.Context) cutil.Err {
	return cutil.FromCErr(C.nix_libexpr_init(ctx.Ptr()))
}

// NewBuilder creates an EvalState builder for the given store.
func NewBuilder(ctx *cutil.Context, store *cstore.Store) *Builder {
	p := C.nix_eval_state_builder_new(ctx.Ptr(), store.Ptr())
	if p == nil {
		return nil
	}
	b := &Builder{ptr: p}
	runtime.SetFinalizer(b, func(b *Builder) { b.Free() })
	return b
}

// Load reads ambient configuration into the builder.
func (b *Builder) Load(ctx *cutil.Context) cutil.Err {
	return cutil.FromCErr(C.nix_eval_state_builder_load(ctx.Ptr(), b.ptr))
}

// SetLookupPath sets NIX_PATH-like lookup paths.
func (b *Builder) SetLookupPath(ctx *cutil.Context, paths []string) cutil.Err {
	cpaths := make([]*C.char, len(paths)+1)
	for i, p := range paths {
		cpaths[i] = C.CString(p)
		defer C.free(unsafe.Pointer(cpaths[i]))
	}
	cpaths[len(paths)] = nil
	var ptr **C.char
	if len(paths) > 0 {
		ptr = (**C.char)(unsafe.Pointer(&cpaths[0]))
	}
	return cutil.FromCErr(C.nix_eval_state_builder_set_lookup_path(ctx.Ptr(), b.ptr, ptr))
}

// SetSetting sets an evaluator setting on the builder.
func (b *Builder) SetSetting(ctx *cutil.Context, key, value string) cutil.Err {
	ck := C.CString(key)
	cv := C.CString(value)
	defer C.free(unsafe.Pointer(ck))
	defer C.free(unsafe.Pointer(cv))
	return cutil.FromCErr(C.nix_eval_state_builder_set_setting(ctx.Ptr(), b.ptr, ck, cv))
}

// Ptr returns the C builder pointer (for flake integration).
func (b *Builder) Ptr() *C.nix_eval_state_builder {
	if b == nil {
		return nil
	}
	return b.ptr
}

// Build creates an EvalState.
func (b *Builder) Build(ctx *cutil.Context) *EvalState {
	p := C.nix_eval_state_build(ctx.Ptr(), b.ptr)
	if p == nil {
		return nil
	}
	s := &EvalState{ptr: p}
	runtime.SetFinalizer(s, func(s *EvalState) { s.Free() })
	return s
}

// Free releases the builder.
func (b *Builder) Free() {
	if b != nil && b.ptr != nil {
		C.nix_eval_state_builder_free(b.ptr)
		b.ptr = nil
		runtime.SetFinalizer(b, nil)
	}
}

// Create is a convenience wrapper around builder defaults.
func Create(ctx *cutil.Context, lookupPath []string, store *cstore.Store) *EvalState {
	cpaths := make([]*C.char, len(lookupPath)+1)
	for i, p := range lookupPath {
		cpaths[i] = C.CString(p)
		defer C.free(unsafe.Pointer(cpaths[i]))
	}
	cpaths[len(lookupPath)] = nil
	var ptr **C.char
	if len(lookupPath) > 0 {
		ptr = (**C.char)(unsafe.Pointer(&cpaths[0]))
	}
	p := C.nix_state_create(ctx.Ptr(), ptr, store.Ptr())
	if p == nil {
		return nil
	}
	s := &EvalState{ptr: p}
	runtime.SetFinalizer(s, func(s *EvalState) { s.Free() })
	return s
}

// Free releases the EvalState.
func (s *EvalState) Free() {
	if s != nil && s.ptr != nil {
		C.nix_state_free(s.ptr)
		s.ptr = nil
		runtime.SetFinalizer(s, nil)
	}
}

// Ptr returns the C pointer.
func (s *EvalState) Ptr() *C.EvalState {
	if s == nil {
		return nil
	}
	return s.ptr
}

// EvalFromString parses and evaluates expr.
func EvalFromString(ctx *cutil.Context, state *EvalState, expr, path string, value unsafe.Pointer) cutil.Err {
	ce := C.CString(expr)
	cp := C.CString(path)
	defer C.free(unsafe.Pointer(ce))
	defer C.free(unsafe.Pointer(cp))
	return cutil.FromCErr(C.nix_expr_eval_from_string(
		ctx.Ptr(), state.ptr, ce, cp, (*C.nix_value)(value),
	))
}

// ValueCall applies fn to arg.
func ValueCall(ctx *cutil.Context, state *EvalState, fn, arg, value unsafe.Pointer) cutil.Err {
	return cutil.FromCErr(C.nix_value_call(
		ctx.Ptr(), state.ptr,
		(*C.nix_value)(fn), (*C.nix_value)(arg), (*C.nix_value)(value),
	))
}

// ValueCallMulti applies a curried function to multiple arguments.
func ValueCallMulti(ctx *cutil.Context, state *EvalState, fn unsafe.Pointer, args []unsafe.Pointer, value unsafe.Pointer) cutil.Err {
	if len(args) == 0 {
		return cutil.OK
	}
	cargs := make([]*C.nix_value, len(args))
	for i, a := range args {
		cargs[i] = (*C.nix_value)(a)
	}
	return cutil.FromCErr(C.nix_value_call_multi(
		ctx.Ptr(), state.ptr, (*C.nix_value)(fn),
		C.size_t(len(args)), &cargs[0], (*C.nix_value)(value),
	))
}

// ValueForce forces evaluation of a thunk.
func ValueForce(ctx *cutil.Context, state *EvalState, value unsafe.Pointer) cutil.Err {
	return cutil.FromCErr(C.nix_value_force(ctx.Ptr(), state.ptr, (*C.nix_value)(value)))
}

// ValueForceDeep deeply forces a value.
func ValueForceDeep(ctx *cutil.Context, state *EvalState, value unsafe.Pointer) cutil.Err {
	return cutil.FromCErr(C.nix_value_force_deep(ctx.Ptr(), state.ptr, (*C.nix_value)(value)))
}

// GCIncref increments GC refcount.
func GCIncref(ctx *cutil.Context, obj unsafe.Pointer) cutil.Err {
	return cutil.FromCErr(C.nix_gc_incref(ctx.Ptr(), obj))
}

// GCDecref decrements GC refcount.
func GCDecref(ctx *cutil.Context, obj unsafe.Pointer) cutil.Err {
	return cutil.FromCErr(C.nix_gc_decref(ctx.Ptr(), obj))
}

// GCNow triggers a GC cycle.
func GCNow() {
	C.nix_gc_now()
}
