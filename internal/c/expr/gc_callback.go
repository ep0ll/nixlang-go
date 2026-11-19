package expr

/*
#cgo pkg-config: nix-expr-c nix-util-c
#include <nix_api_expr.h>
#include <nix_api_util.h>
*/
import "C"
import "unsafe"

//export cgoExprGCFinalizer
func cgoExprGCFinalizer(obj unsafe.Pointer, userdata unsafe.Pointer) {
	defer func() { _ = recover() }()
	if userdata == nil {
		return
	}
	id := uintptr(userdata)
	gcFinalizerMu.Lock()
	fn := gcFinalizers[id]
	delete(gcFinalizers, id)
	gcFinalizerMu.Unlock()
	if fn != nil {
		fn(obj)
	}
}
