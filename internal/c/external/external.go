// Package external provides safe Go callbacks for the Nix external-value C API.
package external

/*
#cgo pkg-config: nix-expr-c nix-store-c nix-util-c
#include <nix_api_external.h>
#include <nix_api_expr.h>
#include <nix_api_util.h>
#include <stdlib.h>
#include "bridge.h"
*/
import "C"
import (
	"runtime"
	"runtime/cgo"
	"sync"
	"unsafe"

	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

// Callbacks implements the callback surface of NixCExternalValueDesc.
// Pointer arguments are borrowed from Nix and must not be retained after the
// callback returns. Errors from string conversion callbacks cause conversion
// to fail according to the Nix external-value contract.
type Callbacks struct {
	Print func(data any, emit func(string) error)
	ShowType func(data any) string
	TypeOf func(data any) string
	CoerceToString func(data any, coerceMore, copyToStore bool, addContext func(string) error) (string, error)
	Equal func(left, right any) bool
	PrintValueAsJSON func(data any, state unsafe.Pointer, strict, copyToStore bool, addContext func(string) error) (string, error)
	PrintValueAsXML func(data any, state unsafe.Pointer, strict, location bool, doc, stringContext, drvsSeen unsafe.Pointer, pos int, addContext func(string) error)
}

type Desc struct {
	mu sync.Mutex
	ptr *C.NixCExternalValueDesc
	callbacks Callbacks
	refs int
	closed bool
}

type externalInstance struct {
	desc *Desc
	data any
}

// Value wraps a Nix ExternalValue reference.
type Value struct {
	ptr *C.ExternalValue
	once sync.Once
}

// NewDesc creates an external-value descriptor. Keep it open while creating
// values; Close prevents new values but defers freeing its C descriptor until
// all values created from it have been collected by Nix.
func NewDesc(callbacks Callbacks) *Desc {
	ptr := C.nixgo_external_desc_new()
	if ptr == nil {
		return nil
	}
	desc := &Desc{ptr: ptr, callbacks: callbacks}
	runtime.SetFinalizer(desc, func(desc *Desc) { desc.Close() })
	return desc
}

func (d *Desc) retain() *C.NixCExternalValueDesc {
	if d == nil {
		return nil
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed || d.ptr == nil {
		return nil
	}
	d.refs++
	return d.ptr
}

func (d *Desc) release() {
	if d == nil {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.refs > 0 {
		d.refs--
	}
	if d.closed && d.refs == 0 && d.ptr != nil {
		C.nixgo_external_desc_free(d.ptr)
		d.ptr = nil
	}
}

// Close prevents creating additional values and releases the descriptor once
// the last Nix external value created from it is finalized.
func (d *Desc) Close() {
	if d == nil {
		return
	}
	d.mu.Lock()
	d.closed = true
	if d.refs == 0 && d.ptr != nil {
		C.nixgo_external_desc_free(d.ptr)
		d.ptr = nil
	}
	d.mu.Unlock()
	runtime.SetFinalizer(d, nil)
}

// Create creates an external value whose opaque data is kept in a cgo.Handle.
// Nix never retains a Go pointer; the handle is released by a Nix GC finalizer.
func Create(ctx *cutil.Context, desc *Desc, data any) *Value {
	descPtr := desc.retain()
	if descPtr == nil {
		return nil
	}
	handle := cgo.NewHandle(&externalInstance{desc: desc, data: data})
	self := C.nixgo_external_handle_to_ptr(C.uintptr_t(handle))
	ptr := C.nix_create_external_value(ctx.Ptr(), descPtr, self)
	if ptr == nil {
		handle.Delete()
		desc.release()
		return nil
	}
	C.nixgo_external_register_finalizer(unsafe.Pointer(ptr), C.uintptr_t(handle))
	value := &Value{ptr: ptr}
	runtime.SetFinalizer(value, func(value *Value) { value.Free() })
	return value
}

// Wrap takes an additional GC reference to a borrowed ExternalValue pointer.
func Wrap(ctx *cutil.Context, ptr unsafe.Pointer) *Value {
	if ptr == nil || C.nix_gc_incref(ctx.Ptr(), ptr) != C.NIX_OK {
		return nil
	}
	value := &Value{ptr: (*C.ExternalValue)(ptr)}
	runtime.SetFinalizer(value, func(value *Value) { value.Free() })
	return value
}

// Free decrements this reference to the Nix external value.
func (v *Value) Free() {
	if v == nil {
		return
	}
	v.once.Do(func() {
		if v.ptr != nil {
			C.nix_gc_decref(nil, unsafe.Pointer(v.ptr))
			v.ptr = nil
		}
		runtime.SetFinalizer(v, nil)
	})
}

// Content extracts the original Go data from a C API external value.
func (v *Value) Content(ctx *cutil.Context) any {
	if v == nil || v.ptr == nil {
		return nil
	}
	ptr := C.nix_get_external_value_content(ctx.Ptr(), v.ptr)
	instance := externalInstanceFor(ptr)
	if instance == nil {
		return nil
	}
	return instance.data
}

// Ptr returns the underlying C ExternalValue pointer.
func (v *Value) Ptr() unsafe.Pointer {
	if v == nil {
		return nil
	}
	return unsafe.Pointer(v.ptr)
}
