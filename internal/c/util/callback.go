package util

/*
#include <nix_api_util.h>
*/
import "C"
import "unsafe"

//export cgoUtilStringCB
func cgoUtilStringCB(start *C.char, n C.uint, userData unsafe.Pointer) {
	defer func() { recover() }()
	if userData == nil {
		return
	}
	s := TakeStringSlot(userData)
	if s != nil {
		*s = C.GoStringN(start, C.int(n))
	}
}
