package store

/*
#cgo pkg-config: nix-store-c nix-util-c
#include <nix_api_store.h>
#include <nix_api_util.h>
#include <stdlib.h>
*/
import "C"
import (
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

var (
	realiseMu    sync.Mutex
	realiseSlots = make(map[uintptr]func(outname string, out *StorePath))
	realiseNext  uintptr = 1
)

func registerRealiseCB(cb func(outname string, out *StorePath)) uintptr {
	realiseMu.Lock()
	defer realiseMu.Unlock()
	id := realiseNext
	realiseNext++
	realiseSlots[id] = cb
	return id
}

func unregisterRealiseCB(id uintptr) {
	realiseMu.Lock()
	defer realiseMu.Unlock()
	delete(realiseSlots, id)
}

//export cgoRealiseCB
func cgoRealiseCB(userdata unsafe.Pointer, outname *C.char, out *C.StorePath) {
	defer func() { recover() }()
	id := uintptr(userdata)
	realiseMu.Lock()
	cb := realiseSlots[id]
	realiseMu.Unlock()
	if cb == nil {
		return
	}
	name := ""
	if outname != nil {
		name = C.GoString(outname)
	}
	var sp *StorePath
	if out != nil {
		sp = &StorePath{ptr: out}
	}
	cb(name, sp)
}
