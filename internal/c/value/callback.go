package value

/*
#include <nix_api_util.h>
*/
import "C"
import (
	"unsafe"

	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

//export cgoValueStringCB
func cgoValueStringCB(start *C.char, n C.uint, userData unsafe.Pointer) {
	defer func() { recover() }()
	if userData == nil {
		return
	}
	s := cutil.TakeStringSlot(userData)
	if s != nil {
		*s = C.GoStringN(start, C.int(n))
	}
}
