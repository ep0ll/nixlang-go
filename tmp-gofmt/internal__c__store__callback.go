package store

/*
#cgo pkg-config: nix-store-c nix-util-c
#include <nix_api_store.h>
#include <nix_api_util.h>
#include <stdlib.h>
#include "bridge.h"
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"

	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

//export cgoStoreStringCB
func cgoStoreStringCB(start *C.char, n C.uint, userData unsafe.Pointer) {
	defer func() { recover() }()
	if userData == nil {
		return
	}
	s := cutil.TakeStringSlot(userData)
	if s != nil {
		*s = C.GoStringN(start, C.int(n))
	}
}

type realiseCallback struct {
	fn         func(outname string, out *StorePath)
	panicValue any
}

var (
	realiseMu    sync.Mutex
	realiseSlots = make(map[uintptr]*realiseCallback)
	realiseNext  uintptr = 1
)

func registerRealiseCB(cb func(outname string, out *StorePath)) uintptr {
	realiseMu.Lock()
	defer realiseMu.Unlock()
	id := realiseNext
	realiseNext++
	realiseSlots[id] = &realiseCallback{fn: cb}
	return id
}

func recordRealisePanic(id uintptr, recovered any) {
	realiseMu.Lock()
	defer realiseMu.Unlock()
	if callback := realiseSlots[id]; callback != nil && callback.panicValue == nil {
		callback.panicValue = recovered
	}
}

func unregisterRealiseCB(id uintptr) any {
	realiseMu.Lock()
	defer realiseMu.Unlock()
	callback := realiseSlots[id]
	delete(realiseSlots, id)
	if callback == nil {
		return nil
	}
	return callback.panicValue
}

//export cgoRealiseCB
func cgoRealiseCB(userdata unsafe.Pointer, outname *C.char, out *C.StorePath) {
	if userdata == nil {
		return
	}
	id := uintptr(C.nixgo_store_userdata_to_id(userdata))
	realiseMu.Lock()
	callback := realiseSlots[id]
	realiseMu.Unlock()
	if callback == nil || callback.fn == nil || callback.panicValue != nil {
		return
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			recordRealisePanic(id, recovered)
		}
	}()

	name := ""
	if outname != nil {
		name = C.GoString(outname)
	}
	var path *StorePath
	if out != nil {
		path = WrapBorrowedStorePath(out)
		defer path.Free()
	}
	callback.fn(name, path)
}
