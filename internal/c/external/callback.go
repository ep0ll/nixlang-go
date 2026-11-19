package external

/*
#cgo pkg-config: nix-expr-c nix-store-c nix-util-c
#include <stdbool.h>
#include <nix_api_external.h>
#include <nix_api_expr.h>
#include <nix_api_util.h>
#include <stdlib.h>
#include "bridge.h"
*/
import "C"
import (
	"fmt"
	"runtime/cgo"
	"unsafe"
)

func externalInstanceFor(ptr unsafe.Pointer) (instance *externalInstance) {
	if ptr == nil {
		return nil
	}
	defer func() {
		if recover() != nil {
			instance = nil
		}
	}()
	handle := cgo.Handle(uintptr(C.nixgo_external_handle_from_ptr(ptr)))
	instance, _ = handle.Value().(*externalInstance)
	return instance
}

func cExternalStringContext(ptr *C.nix_string_context) func(string) error {
	return func(value string) error {
		ctx := C.nix_c_context_create()
		if ctx == nil {
			return fmt.Errorf("nix: failed to allocate callback context")
		}
		defer C.nix_c_context_free(ctx)
		cv := C.CString(value)
		defer C.free(unsafe.Pointer(cv))
		err := C.nix_external_add_string_context(ctx, ptr, cv)
		if err != C.NIX_OK {
			return fmt.Errorf("nix: external string context: %s", C.GoString(C.nix_err_msg(nil, ctx, nil)))
		}
		return nil
	}
}

func cExternalEmit(printer *C.nix_printer) func(string) error {
	return func(value string) error {
		ctx := C.nix_c_context_create()
		if ctx == nil {
			return fmt.Errorf("nix: failed to allocate callback context")
		}
		defer C.nix_c_context_free(ctx)
		cv := C.CString(value)
		defer C.free(unsafe.Pointer(cv))
		err := C.nix_external_print(ctx, printer, cv)
		if err != C.NIX_OK {
			return fmt.Errorf("nix: external print: %s", C.GoString(C.nix_err_msg(nil, ctx, nil)))
		}
		return nil
	}
}

func setExternalString(result *C.nix_string_return, value string) {
	cv := C.CString(value)
	defer C.free(unsafe.Pointer(cv))
	C.nix_set_string_return(result, cv)
}

//export cgoExternalPrint
func cgoExternalPrint(self unsafe.Pointer, printer *C.nix_printer) {
	defer func() { _ = recover() }()
	instance := externalInstanceFor(self)
	if instance != nil && instance.desc.callbacks.Print != nil {
		instance.desc.callbacks.Print(instance.data, cExternalEmit(printer))
	}
}

//export cgoExternalShowType
func cgoExternalShowType(self unsafe.Pointer, result *C.nix_string_return) {
	defer func() { _ = recover() }()
	instance := externalInstanceFor(self)
	if instance != nil && instance.desc.callbacks.ShowType != nil {
		setExternalString(result, instance.desc.callbacks.ShowType(instance.data))
	}
}

//export cgoExternalTypeOf
func cgoExternalTypeOf(self unsafe.Pointer, result *C.nix_string_return) {
	defer func() { _ = recover() }()
	instance := externalInstanceFor(self)
	if instance != nil && instance.desc.callbacks.TypeOf != nil {
		setExternalString(result, instance.desc.callbacks.TypeOf(instance.data))
	}
}

//export cgoExternalCoerceToString
func cgoExternalCoerceToString(
	self unsafe.Pointer,
	stringContext *C.nix_string_context,
	coerceMore C.int,
	copyToStore C.int,
	result *C.nix_string_return,
) {
	defer func() { _ = recover() }()
	instance := externalInstanceFor(self)
	if instance == nil || instance.desc.callbacks.CoerceToString == nil {
		setExternalString(result, "")
		return
	}
	value, err := instance.desc.callbacks.CoerceToString(
		instance.data, coerceMore != 0, copyToStore != 0, cExternalStringContext(stringContext),
	)
	if err != nil {
		value = ""
	}
	setExternalString(result, value)
}

//export cgoExternalEqual
func cgoExternalEqual(self unsafe.Pointer, other unsafe.Pointer) (equal C.int) {
	defer func() {
		if recover() != nil {
			equal = 0
		}
	}()
	left := externalInstanceFor(self)
	right := externalInstanceFor(other)
	if left == nil || right == nil || left.desc.callbacks.Equal == nil {
		return 0
	}
	if left.desc.callbacks.Equal(left.data, right.data) {
		return 1
	}
	return 0
}

//export cgoExternalPrintValueAsJSON
func cgoExternalPrintValueAsJSON(
	self unsafe.Pointer,
	state *C.EvalState,
	strict C.bool,
	stringContext *C.nix_string_context,
	copyToStore C.bool,
	result *C.nix_string_return,
) {
	defer func() { _ = recover() }()
	instance := externalInstanceFor(self)
	if instance == nil || instance.desc.callbacks.PrintValueAsJSON == nil {
		setExternalString(result, "")
		return
	}
	value, err := instance.desc.callbacks.PrintValueAsJSON(
		instance.data, unsafe.Pointer(state), bool(strict), bool(copyToStore), cExternalStringContext(stringContext),
	)
	if err != nil {
		value = ""
	}
	setExternalString(result, value)
}

//export cgoExternalPrintValueAsXML
func cgoExternalPrintValueAsXML(
	self unsafe.Pointer,
	state *C.EvalState,
	strict C.int,
	location C.int,
	doc unsafe.Pointer,
	stringContext *C.nix_string_context,
	drvsSeen unsafe.Pointer,
	pos C.int,
) {
	defer func() { _ = recover() }()
	instance := externalInstanceFor(self)
	if instance != nil && instance.desc.callbacks.PrintValueAsXML != nil {
		instance.desc.callbacks.PrintValueAsXML(
			instance.data, unsafe.Pointer(state), strict != 0, location != 0,
			doc, unsafe.Pointer(stringContext), drvsSeen, int(pos), cExternalStringContext(stringContext),
		)
	}
}

//export cgoExternalFinalizer
func cgoExternalFinalizer(_ unsafe.Pointer, userdata unsafe.Pointer) {
	defer func() { _ = recover() }()
	if userdata == nil {
		return
	}
	defer C.nixgo_external_handle_free(userdata)
	handle := cgo.Handle(uintptr(C.nixgo_external_handle_from_ptr(userdata)))
	instance, _ := handle.Value().(*externalInstance)
	handle.Delete()
	if instance != nil && instance.desc != nil {
		instance.desc.release()
	}
}
