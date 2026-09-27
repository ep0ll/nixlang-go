// Package value — typename, copy, and string-realise bindings.
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
	"unsafe"

	cexpr "github.com/ep0ll/nixlang-go/internal/c/expr"
	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

// GetTypename returns the Nix type name string (e.g. "int", "attrs").
// The C API allocates the string; this function frees it after copying.
func (v *Value) GetTypename(ctx *cutil.Context) string {
	if v == nil || v.ptr == nil {
		return ""
	}
	p := C.nix_get_typename(ctx.Ptr(), v.ptr)
	if p == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(p))
	return C.GoString(p)
}

// CopyValue copies source into v (both must be allocated values).
func (v *Value) CopyValue(ctx *cutil.Context, source *Value) cutil.Err {
	if v == nil || source == nil {
		return cutil.ErrUnknown
	}
	return cutil.FromCErr(C.nix_copy_value(ctx.Ptr(), v.ptr, source.ptr))
}

// RealisedString wraps nix_realised_string (string with store-path context).
type RealisedString struct {
	ptr *C.nix_realised_string
}

// StringRealise forces a string value and collects derivation/store-path context.
// isIFD controls allow-import-from-derivation behavior.
func StringRealise(ctx *cutil.Context, state *cexpr.EvalState, v *Value, isIFD bool) *RealisedString {
	if v == nil || v.ptr == nil {
		return nil
	}
	p := C.nix_string_realise(ctx.Ptr(), state.Ptr(), v.ptr, C.bool(isIFD))
	if p == nil {
		return nil
	}
	return &RealisedString{ptr: p}
}

// Free releases a realised string.
func (r *RealisedString) Free() {
	if r != nil && r.ptr != nil {
		C.nix_realised_string_free(r.ptr)
		r.ptr = nil
	}
}

// Buffer returns the realised string contents.
func (r *RealisedString) Buffer() string {
	if r == nil || r.ptr == nil {
		return ""
	}
	start := C.nix_realised_string_get_buffer_start(r.ptr)
	n := C.nix_realised_string_get_buffer_size(r.ptr)
	if start == nil || n == 0 {
		return ""
	}
	return C.GoStringN(start, C.int(n))
}

// StorePathCount returns how many store paths are in the string context.
func (r *RealisedString) StorePathCount() int {
	if r == nil || r.ptr == nil {
		return 0
	}
	return int(C.nix_realised_string_get_store_path_count(r.ptr))
}

// StorePathAt returns the store path at index (borrowed; do not free).
func (r *RealisedString) StorePathAt(index int) unsafe.Pointer {
	if r == nil || r.ptr == nil {
		return nil
	}
	p := C.nix_realised_string_get_store_path(r.ptr, C.size_t(index))
	return unsafe.Pointer(p)
}
