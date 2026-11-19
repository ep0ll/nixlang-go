// Package store implements filesystem-closure access from nix_api_store.h.
package store

/*
#cgo pkg-config: nix-store-c nix-util-c
#include <nix_api_store.h>
#include <nix_api_util.h>
#include <stdint.h>
#include "bridge.h"

typedef void (*nixgo_fs_closure_callback)(nix_c_context *, void *, const StorePath *);
extern void cgoStoreFSClosureCB(nix_c_context *, void *, const StorePath *);
*/
import "C"
import (
	"sync"
	"unsafe"

	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

var (
	fsClosureMu        sync.Mutex
	fsClosureCallbacks         = make(map[uintptr]func(*StorePath))
	fsClosureNext      uintptr = 1
)

func registerFSClosureCallback(cb func(*StorePath)) uintptr {
	fsClosureMu.Lock()
	defer fsClosureMu.Unlock()
	id := fsClosureNext
	fsClosureNext++
	fsClosureCallbacks[id] = cb
	return id
}

func takeFSClosureCallback(id uintptr) func(*StorePath) {
	fsClosureMu.Lock()
	defer fsClosureMu.Unlock()
	return fsClosureCallbacks[id]
}

func unregisterFSClosureCallback(id uintptr) {
	fsClosureMu.Lock()
	defer fsClosureMu.Unlock()
	delete(fsClosureCallbacks, id)
}

// GetFSClosure enumerates a store path's filesystem closure.
//
// The callback receives an independently-owned StorePath and must call Free
// when it is no longer needed. The callback is invoked synchronously and in
// unspecified order. flipDirection selects forward or reverse closure; the
// other flags control inclusion of derivations and outputs.
func (s *Store) GetFSClosure(
	ctx *cutil.Context,
	path *StorePath,
	flipDirection, includeOutputs, includeDerivers bool,
	cb func(*StorePath),
) cutil.Err {
	if s == nil || s.ptr == nil || path == nil || path.ptr == nil {
		return cutil.ErrUnknown
	}
	if cb == nil {
		cb = func(*StorePath) {}
	}
	id := registerFSClosureCallback(cb)
	defer unregisterFSClosureCallback(id)
	return cutil.FromCErr(C.nix_store_get_fs_closure(
		ctx.Ptr(), s.ptr, path.ptr,
		C.bool(flipDirection), C.bool(includeOutputs), C.bool(includeDerivers),
		unsafe.Pointer(C.nixgo_store_id_to_userdata(C.uintptr_t(id))), C.nixgo_fs_closure_callback(C.cgoStoreFSClosureCB),
	))
}
