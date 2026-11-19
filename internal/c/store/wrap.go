// Package store contains helpers for wrapping pointers returned by other Nix C API modules.
package store

/*
#cgo pkg-config: nix-store-c nix-util-c
#include <nix_api_store.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// WrapPath takes ownership of a StorePath pointer returned by another Nix C API module.
// The caller must not use or free the pointer separately after wrapping it.
func WrapPath(p unsafe.Pointer) *StorePath {
	if p == nil {
		return nil
	}
	return wrapOwnedStorePath((*C.StorePath)(p))
}

func wrapOwnedStorePath(p *C.StorePath) *StorePath {
	if p == nil {
		return nil
	}
	path := &StorePath{ptr: p, owned: true}
	runtime.SetFinalizer(path, func(path *StorePath) { path.Free() })
	return path
}

// ClonePath clones a borrowed StorePath returned by another Nix C API object.
// The returned wrapper owns the cloned reference.
func ClonePath(p unsafe.Pointer) *StorePath {
	if p == nil {
		return nil
	}
	cloned := C.nix_store_path_clone((*C.StorePath)(p))
	if cloned == nil {
		return nil
	}
	return wrapOwnedStorePath(cloned)
}
