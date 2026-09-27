// Package flake provides low-level cgo bindings for nix_api_flake.h.
package flake

/*
#cgo pkg-config: nix-flake-c nix-fetchers-c nix-expr-c nix-store-c nix-util-c
#include <nix_api_flake.h>
#include <nix_api_fetchers.h>
#include <nix_api_expr.h>
#include <nix_api_value.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"unsafe"

	cexpr "github.com/ep0ll/nixlang-go/internal/c/expr"
	cfetch "github.com/ep0ll/nixlang-go/internal/c/fetchers"
	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

// Settings wraps nix_flake_settings.
type Settings struct {
	ptr *C.nix_flake_settings
}

// Reference wraps nix_flake_reference.
type Reference struct {
	ptr *C.nix_flake_reference
}

// LockFlags wraps nix_flake_lock_flags.
type LockFlags struct {
	ptr *C.nix_flake_lock_flags
}

// LockedFlake wraps nix_locked_flake.
type LockedFlake struct {
	ptr *C.nix_locked_flake
}

// ParseFlags wraps nix_flake_reference_parse_flags.
type ParseFlags struct {
	ptr *C.nix_flake_reference_parse_flags
}

// NewSettings creates default flake settings.
func NewSettings(ctx *cutil.Context) *Settings {
	p := C.nix_flake_settings_new(ctx.Ptr())
	if p == nil {
		return nil
	}
	s := &Settings{ptr: p}
	runtime.SetFinalizer(s, func(s *Settings) { s.Free() })
	return s
}

func (s *Settings) Free() {
	if s != nil && s.ptr != nil {
		C.nix_flake_settings_free(s.ptr)
		s.ptr = nil
		runtime.SetFinalizer(s, nil)
	}
}

func (s *Settings) Ptr() *C.nix_flake_settings {
	if s == nil {
		return nil
	}
	return s.ptr
}

// AddToEvalStateBuilder registers flake builtins on the builder.
func (s *Settings) AddToEvalStateBuilder(ctx *cutil.Context, builder *cexpr.Builder) cutil.Err {
	return cutil.FromCErr(C.nix_flake_settings_add_to_eval_state_builder(ctx.Ptr(), s.ptr, builder.Ptr()))
}

// NewParseFlags creates parse flags.
func NewParseFlags(ctx *cutil.Context, settings *Settings) *ParseFlags {
	p := C.nix_flake_reference_parse_flags_new(ctx.Ptr(), settings.ptr)
	if p == nil {
		return nil
	}
	f := &ParseFlags{ptr: p}
	runtime.SetFinalizer(f, func(f *ParseFlags) { f.Free() })
	return f
}

func (f *ParseFlags) Free() {
	if f != nil && f.ptr != nil {
		C.nix_flake_reference_parse_flags_free(f.ptr)
		f.ptr = nil
		runtime.SetFinalizer(f, nil)
	}
}

func (f *ParseFlags) SetBaseDirectory(ctx *cutil.Context, dir string) cutil.Err {
	cd := C.CString(dir)
	defer C.free(unsafe.Pointer(cd))
	return cutil.FromCErr(C.nix_flake_reference_parse_flags_set_base_directory(
		ctx.Ptr(), f.ptr, cd, C.size_t(len(dir)),
	))
}

// NewLockFlags creates lock flags.
func NewLockFlags(ctx *cutil.Context, settings *Settings) *LockFlags {
	p := C.nix_flake_lock_flags_new(ctx.Ptr(), settings.ptr)
	if p == nil {
		return nil
	}
	f := &LockFlags{ptr: p}
	runtime.SetFinalizer(f, func(f *LockFlags) { f.Free() })
	return f
}

func (f *LockFlags) Free() {
	if f != nil && f.ptr != nil {
		C.nix_flake_lock_flags_free(f.ptr)
		f.ptr = nil
		runtime.SetFinalizer(f, nil)
	}
}

func (f *LockFlags) SetModeCheck(ctx *cutil.Context) cutil.Err {
	return cutil.FromCErr(C.nix_flake_lock_flags_set_mode_check(ctx.Ptr(), f.ptr))
}

func (f *LockFlags) SetModeVirtual(ctx *cutil.Context) cutil.Err {
	return cutil.FromCErr(C.nix_flake_lock_flags_set_mode_virtual(ctx.Ptr(), f.ptr))
}

func (f *LockFlags) SetModeWriteAsNeeded(ctx *cutil.Context) cutil.Err {
	return cutil.FromCErr(C.nix_flake_lock_flags_set_mode_write_as_needed(ctx.Ptr(), f.ptr))
}

func (f *LockFlags) AddInputOverride(ctx *cutil.Context, inputPath string, ref *Reference) cutil.Err {
	cip := C.CString(inputPath)
	defer C.free(unsafe.Pointer(cip))
	return cutil.FromCErr(C.nix_flake_lock_flags_add_input_override(ctx.Ptr(), f.ptr, cip, ref.ptr))
}

// Lock locks a flake reference.
func Lock(
	ctx *cutil.Context,
	fetch *cfetch.Settings,
	settings *Settings,
	state *cexpr.EvalState,
	flags *LockFlags,
	ref *Reference,
) *LockedFlake {
	p := C.nix_flake_lock(ctx.Ptr(), fetch.Ptr(), settings.ptr, state.Ptr(), flags.ptr, ref.ptr)
	if p == nil {
		return nil
	}
	lf := &LockedFlake{ptr: p}
	runtime.SetFinalizer(lf, func(lf *LockedFlake) { lf.Free() })
	return lf
}

func (lf *LockedFlake) Free() {
	if lf != nil && lf.ptr != nil {
		C.nix_locked_flake_free(lf.ptr)
		lf.ptr = nil
		runtime.SetFinalizer(lf, nil)
	}
}

// GetOutputAttrs returns flake outputs as an unsafe.Pointer to nix_value.
func (lf *LockedFlake) GetOutputAttrs(ctx *cutil.Context, settings *Settings, state *cexpr.EvalState) unsafe.Pointer {
	return unsafe.Pointer(C.nix_locked_flake_get_output_attrs(ctx.Ptr(), settings.ptr, state.Ptr(), lf.ptr))
}

// ParseReference parses a flake reference string.
func ParseReference(
	ctx *cutil.Context,
	fetch *cfetch.Settings,
	settings *Settings,
	flags *ParseFlags,
	str string,
) (ref *Reference, fragment string, err cutil.Err) {
	var out *C.nix_flake_reference
	var frag string
	ud := cutil.RegisterStringSlot(&frag)
	cs := C.CString(str)
	defer C.free(unsafe.Pointer(cs))
	e := C.nix_flake_reference_and_fragment_from_string(
		ctx.Ptr(), fetch.Ptr(), settings.ptr, flags.ptr,
		cs, C.size_t(len(str)),
		&out,
		C.nix_get_string_callback(C.cgoFlakeStringCB), ud,
	)
	err = cutil.FromCErr(e)
	if out != nil {
		ref = &Reference{ptr: out}
		runtime.SetFinalizer(ref, func(r *Reference) { r.Free() })
	}
	return ref, frag, err
}

func (r *Reference) Free() {
	if r != nil && r.ptr != nil {
		C.nix_flake_reference_free(r.ptr)
		r.ptr = nil
		runtime.SetFinalizer(r, nil)
	}
}

func (r *Reference) Ptr() *C.nix_flake_reference {
	if r == nil {
		return nil
	}
	return r.ptr
}
