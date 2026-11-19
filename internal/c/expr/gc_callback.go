package expr

/*
#cgo pkg-config: nix-expr-c nix-util-c
#include <nix_api_expr.h>
#include <nix_api_util.h>
#include "bridge.h"
*/
import "C"
import "unsafe"

//export cgoExprGCFinalizer
func cgoExprGCFinalizer(obj unsafe.Pointer, userdata unsafe.Pointer) {
	defer func() { _ = recover() }()
	if userdata == nil {
		return
	}
	id := uintptr(C.nixgo_expr_userdata_to_id(userdata))
	gcFinalizerMu.Lock()
	fn := gcFinalizers[id]
	delete(gcFinalizers, id)
	objectID := uintptr(obj)
	if gcFinalizerObjects[objectID] == id {
		delete(gcFinalizerObjects, objectID)
	}
	gcFinalizerMu.Unlock()
	if fn != nil {
		fn(obj)
	}
}
