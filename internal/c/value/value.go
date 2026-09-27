// Package value provides low-level cgo bindings for nix_api_value.h.
package value

/*
#cgo pkg-config: nix-expr-c nix-store-c nix-util-c
#include <nix_api_value.h>
#include <nix_api_expr.h>
#include <nix_api_util.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"sync"
	"unsafe"

	cexpr "github.com/ep0ll/nixlang-go/internal/c/expr"
	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

// ValueType mirrors the C enum.
type ValueType C.ValueType

const (
	TypeThunk    ValueType = C.NIX_TYPE_THUNK
	TypeInt      ValueType = C.NIX_TYPE_INT
	TypeFloat    ValueType = C.NIX_TYPE_FLOAT
	TypeBool     ValueType = C.NIX_TYPE_BOOL
	TypeString   ValueType = C.NIX_TYPE_STRING
	TypePath     ValueType = C.NIX_TYPE_PATH
	TypeNull     ValueType = C.NIX_TYPE_NULL
	TypeAttrs    ValueType = C.NIX_TYPE_ATTRS
	TypeList     ValueType = C.NIX_TYPE_LIST
	TypeFunction ValueType = C.NIX_TYPE_FUNCTION
	TypeExternal ValueType = C.NIX_TYPE_EXTERNAL
	TypeFailed   ValueType = C.NIX_TYPE_FAILED
)

// Value wraps a GC-managed nix_value.
// Decref is idempotent; always call Close/Decref explicitly — finalizers are a backup only.
type Value struct {
	ptr  *C.nix_value
	once sync.Once
}

// ListBuilder wraps ListBuilder.
type ListBuilder struct {
	ptr *C.ListBuilder
}

// BindingsBuilder wraps BindingsBuilder.
type BindingsBuilder struct {
	ptr *C.BindingsBuilder
}

// PrimOp wraps PrimOp.
type PrimOp struct {
	ptr *C.PrimOp
}

// WrapPtr wraps an existing C nix_value pointer (takes ownership of one ref).
// Accepts unsafe.Pointer so callers in other packages can pass C pointers.
func WrapPtr(p unsafe.Pointer) *Value {
	if p == nil {
		return nil
	}
	v := &Value{ptr: (*C.nix_value)(p)}
	runtime.SetFinalizer(v, func(v *Value) { v.Decref() })
	return v
}

// Alloc allocates a new value bound to state.
func Alloc(ctx *cutil.Context, state *cexpr.EvalState) *Value {
	p := C.nix_alloc_value(ctx.Ptr(), state.Ptr())
	if p == nil {
		return nil
	}
	return WrapPtr(unsafe.Pointer(p))
}

// Ptr returns unsafe.Pointer for cross-package use.
func (v *Value) Ptr() unsafe.Pointer {
	if v == nil {
		return nil
	}
	return unsafe.Pointer(v.ptr)
}

// CPtr returns the typed C pointer.
func (v *Value) CPtr() *C.nix_value {
	if v == nil {
		return nil
	}
	return v.ptr
}

// Incref increments the GC refcount.
func (v *Value) Incref(ctx *cutil.Context) cutil.Err {
	return cutil.FromCErr(C.nix_value_incref(ctx.Ptr(), v.ptr))
}

// Decref decrements the GC refcount. Idempotent; safe with concurrent finalizer.
func (v *Value) Decref() {
	if v == nil {
		return
	}
	v.once.Do(func() {
		if v.ptr != nil {
			C.nix_value_decref(nil, v.ptr)
			v.ptr = nil
		}
		runtime.SetFinalizer(v, nil)
	})
}

// Type returns the value type.
func (v *Value) Type(ctx *cutil.Context) ValueType {
	return ValueType(C.nix_get_type(ctx.Ptr(), v.ptr))
}

// GetBool extracts a boolean.
func (v *Value) GetBool(ctx *cutil.Context) bool {
	return bool(C.nix_get_bool(ctx.Ptr(), v.ptr))
}

// GetInt extracts an integer.
func (v *Value) GetInt(ctx *cutil.Context) int64 {
	return int64(C.nix_get_int(ctx.Ptr(), v.ptr))
}

// GetFloat extracts a float.
func (v *Value) GetFloat(ctx *cutil.Context) float64 {
	return float64(C.nix_get_float(ctx.Ptr(), v.ptr))
}

// GetString extracts a string via callback.
func (v *Value) GetString(ctx *cutil.Context) (string, cutil.Err) {
	var out string
	ud := cutil.RegisterStringSlot(&out)
	err := cutil.FromCErr(C.nix_get_string(ctx.Ptr(), v.ptr, C.nix_get_string_callback(C.cgoValueStringCB), ud))
	return out, err
}

// GetPathString returns a path as string.
func (v *Value) GetPathString(ctx *cutil.Context) string {
	p := C.nix_get_path_string(ctx.Ptr(), v.ptr)
	if p == nil {
		return ""
	}
	return C.GoString(p)
}

// ListSize returns the list length.
func (v *Value) ListSize(ctx *cutil.Context) uint {
	return uint(C.nix_get_list_size(ctx.Ptr(), v.ptr))
}

// AttrsSize returns the attrset size.
func (v *Value) AttrsSize(ctx *cutil.Context) uint {
	return uint(C.nix_get_attrs_size(ctx.Ptr(), v.ptr))
}

// ListByIdx returns element at index (caller Decref).
func (v *Value) ListByIdx(ctx *cutil.Context, state *cexpr.EvalState, idx uint) *Value {
	p := C.nix_get_list_byidx(ctx.Ptr(), v.ptr, state.Ptr(), C.uint(idx))
	return WrapPtr(unsafe.Pointer(p))
}

// ListByIdxLazy returns element without forcing.
func (v *Value) ListByIdxLazy(ctx *cutil.Context, state *cexpr.EvalState, idx uint) *Value {
	p := C.nix_get_list_byidx_lazy(ctx.Ptr(), v.ptr, state.Ptr(), C.uint(idx))
	return WrapPtr(unsafe.Pointer(p))
}

// AttrByName returns attribute by name.
func (v *Value) AttrByName(ctx *cutil.Context, state *cexpr.EvalState, name string) *Value {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	p := C.nix_get_attr_byname(ctx.Ptr(), v.ptr, state.Ptr(), cn)
	return WrapPtr(unsafe.Pointer(p))
}

// AttrByNameLazy returns attribute without forcing.
func (v *Value) AttrByNameLazy(ctx *cutil.Context, state *cexpr.EvalState, name string) *Value {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	p := C.nix_get_attr_byname_lazy(ctx.Ptr(), v.ptr, state.Ptr(), cn)
	return WrapPtr(unsafe.Pointer(p))
}

// HasAttr reports whether an attribute exists.
func (v *Value) HasAttr(ctx *cutil.Context, state *cexpr.EvalState, name string) bool {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	return bool(C.nix_has_attr_byname(ctx.Ptr(), v.ptr, state.Ptr(), cn))
}

// AttrByIdx returns attribute by index and its name.
func (v *Value) AttrByIdx(ctx *cutil.Context, state *cexpr.EvalState, idx uint) (name string, val *Value) {
	var cname *C.char
	p := C.nix_get_attr_byidx(ctx.Ptr(), v.ptr, state.Ptr(), C.uint(idx), &cname)
	if cname != nil {
		name = C.GoString(cname)
	}
	return name, WrapPtr(unsafe.Pointer(p))
}

// AttrNameByIdx returns only the attribute name at index.
func (v *Value) AttrNameByIdx(ctx *cutil.Context, state *cexpr.EvalState, idx uint) string {
	p := C.nix_get_attr_name_byidx(ctx.Ptr(), v.ptr, state.Ptr(), C.uint(idx))
	if p == nil {
		return ""
	}
	return C.GoString(p)
}

// InitNull initializes to null.
func (v *Value) InitNull(ctx *cutil.Context) cutil.Err {
	return cutil.FromCErr(C.nix_init_null(ctx.Ptr(), v.ptr))
}

// InitBool initializes a boolean.
func (v *Value) InitBool(ctx *cutil.Context, b bool) cutil.Err {
	return cutil.FromCErr(C.nix_init_bool(ctx.Ptr(), v.ptr, C.bool(b)))
}

// InitInt initializes an integer.
func (v *Value) InitInt(ctx *cutil.Context, i int64) cutil.Err {
	return cutil.FromCErr(C.nix_init_int(ctx.Ptr(), v.ptr, C.int64_t(i)))
}

// InitFloat initializes a float.
func (v *Value) InitFloat(ctx *cutil.Context, f float64) cutil.Err {
	return cutil.FromCErr(C.nix_init_float(ctx.Ptr(), v.ptr, C.double(f)))
}

// InitString initializes a string.
func (v *Value) InitString(ctx *cutil.Context, s string) cutil.Err {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return cutil.FromCErr(C.nix_init_string(ctx.Ptr(), v.ptr, cs))
}

// InitPathString initializes a path.
func (v *Value) InitPathString(ctx *cutil.Context, state *cexpr.EvalState, path string) cutil.Err {
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))
	return cutil.FromCErr(C.nix_init_path_string(ctx.Ptr(), state.Ptr(), v.ptr, cs))
}

// InitApply sets a thunk for function application.
func (v *Value) InitApply(ctx *cutil.Context, fn, arg *Value) cutil.Err {
	return cutil.FromCErr(C.nix_init_apply(ctx.Ptr(), v.ptr, fn.ptr, arg.ptr))
}

// MakeListBuilder creates a list builder.
func MakeListBuilder(ctx *cutil.Context, state *cexpr.EvalState, capacity uint) *ListBuilder {
	p := C.nix_make_list_builder(ctx.Ptr(), state.Ptr(), C.size_t(capacity))
	if p == nil {
		return nil
	}
	b := &ListBuilder{ptr: p}
	runtime.SetFinalizer(b, func(b *ListBuilder) { b.Free() })
	return b
}

// Insert inserts a value at index.
func (b *ListBuilder) Insert(ctx *cutil.Context, index uint, v *Value) cutil.Err {
	return cutil.FromCErr(C.nix_list_builder_insert(ctx.Ptr(), b.ptr, C.uint(index), v.ptr))
}

// Free frees the list builder.
func (b *ListBuilder) Free() {
	if b != nil && b.ptr != nil {
		C.nix_list_builder_free(b.ptr)
		b.ptr = nil
		runtime.SetFinalizer(b, nil)
	}
}

// MakeList finalizes the builder into value.
func (v *Value) MakeList(ctx *cutil.Context, b *ListBuilder) cutil.Err {
	return cutil.FromCErr(C.nix_make_list(ctx.Ptr(), b.ptr, v.ptr))
}

// MakeBindingsBuilder creates an attrs builder.
func MakeBindingsBuilder(ctx *cutil.Context, state *cexpr.EvalState, capacity uint) *BindingsBuilder {
	p := C.nix_make_bindings_builder(ctx.Ptr(), state.Ptr(), C.size_t(capacity))
	if p == nil {
		return nil
	}
	b := &BindingsBuilder{ptr: p}
	runtime.SetFinalizer(b, func(b *BindingsBuilder) { b.Free() })
	return b
}

// Insert inserts a binding.
func (b *BindingsBuilder) Insert(ctx *cutil.Context, name string, v *Value) cutil.Err {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	return cutil.FromCErr(C.nix_bindings_builder_insert(ctx.Ptr(), b.ptr, cn, v.ptr))
}

// Free frees the bindings builder.
func (b *BindingsBuilder) Free() {
	if b != nil && b.ptr != nil {
		C.nix_bindings_builder_free(b.ptr)
		b.ptr = nil
		runtime.SetFinalizer(b, nil)
	}
}

// MakeAttrs finalizes the builder into value.
func (v *Value) MakeAttrs(ctx *cutil.Context, b *BindingsBuilder) cutil.Err {
	return cutil.FromCErr(C.nix_make_attrs(ctx.Ptr(), v.ptr, b.ptr))
}
