package store

/*
#cgo pkg-config: nix-store-c nix-util-c
#include <nix_api_store.h>
#include <nix_api_util.h>
*/
import "C"
import "unsafe"

//export cgoStoreFSClosureCB
func cgoStoreFSClosureCB(_ *C.nix_c_context, userdata unsafe.Pointer, borrowed *C.StorePath) {
	defer func() { _ = recover() }()
	if userdata == nil || borrowed == nil {
		return
	}
	cb := takeFSClosureCallback(uintptr(userdata))
	if cb == nil {
		return
	}
	owned := C.nix_store_path_clone(borrowed)
	if owned == nil {
		return
	}
	cb(wrapOwnedStorePath(owned))
}
