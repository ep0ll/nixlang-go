// Package external provides low-level cgo bindings for nix_api_external.h.
package external

/*
#cgo pkg-config: nix-expr-c nix-util-c
#include <nix_api_external.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"unsafe"

	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

// Desc is an alias for the C descriptor (populate function pointers via C).
// Full callback wiring requires stable CGO exports; this package provides
// the create/get surface for use from higher layers that supply a filled desc.
type Desc struct {
	ptr *C.NixCExternalValueDesc
}

// Value wraps ExternalValue.
type Value struct {
	ptr *C.ExternalValue
}

// Create creates an external value. desc must remain valid for the value lifetime.
// data is stored as an opaque pointer (caller manages lifetime).
func Create(ctx *cutil.Context, desc *Desc, data any) *Value {
	if desc == nil || desc.ptr == nil {
		return nil
	}
	// Store data in a heap box for the C side.
	box := &data
	p := C.nix_create_external_value(ctx.Ptr(), desc.ptr, unsafe.Pointer(box))
	if p == nil {
		return nil
	}
	v := &Value{ptr: p}
	runtime.SetFinalizer(v, (*Value).Free)
	return v
}

// Free decrements the GC refcount.
func (v *Value) Free() {
	if v != nil && v.ptr != nil {
		C.nix_gc_decref(nil, unsafe.Pointer(v.ptr))
		v.ptr = nil
		runtime.SetFinalizer(v, nil)
	}
}

// Content extracts the opaque pointer.
func (v *Value) Content(ctx *cutil.Context) any {
	p := C.nix_get_external_value_content(ctx.Ptr(), v.ptr)
	if p == nil {
		return nil
	}
	box := (*any)(p)
	return *box
}

// Ptr returns the C pointer.
func (v *Value) Ptr() *C.ExternalValue {
	if v == nil {
		return nil
	}
	return v.ptr
}
