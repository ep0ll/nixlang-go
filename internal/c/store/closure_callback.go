package store

/*
#cgo pkg-config: nix-store-c nix-util-c
#include <nix_api_store.h>
#include <nix_api_util.h>
#include "bridge.h"
*/
import "C"
import (
	"fmt"
	"unsafe"

	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

//export cgoStoreFSClosureCB
func cgoStoreFSClosureCB(context *C.nix_c_context, userdata unsafe.Pointer, borrowed *C.StorePath) {
	ctx := cutil.WrapBorrowedContext(unsafe.Pointer(context))
	var path *StorePath
	transferred := false
	defer func() {
		if recovered := recover(); recovered != nil {
			if !transferred && path != nil {
				path.Free()
			}
			_ = ctx.SetErr(cutil.ErrUnknown, fmt.Sprintf("panic in filesystem-closure callback: %v", recovered))
		}
	}()

	if userdata == nil || borrowed == nil {
		return
	}
	cb := takeFSClosureCallback(uintptr(C.nixgo_store_userdata_to_id(userdata)))
	if cb == nil {
		return
	}
	owned := C.nix_store_path_clone(borrowed)
	if owned == nil {
		_ = ctx.SetErr(cutil.ErrUnknown, "nix: failed to clone a filesystem-closure path")
		return
	}
	path = wrapOwnedStorePath(owned)
	cb(path)
	transferred = true
}
