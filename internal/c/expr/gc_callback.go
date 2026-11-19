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
	tokenID := uintptr(userdata)
	gcFinalizerMu.Lock()
	fn, ok := gcFinalizers[tokenID]
	if ok {
		delete(gcFinalizers, tokenID)
		objectID := uintptr(obj)
		if gcFinalizerObjects[objectID] == tokenID {
			delete(gcFinalizerObjects, objectID)
		}
	}
	gcFinalizerMu.Unlock()
	if !ok {
		return
	}
	C.nixgo_expr_userdata_free(userdata)
	if fn != nil {
		fn(obj)
	}
}
