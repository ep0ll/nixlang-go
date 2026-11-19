// Package util provides the remaining error-inspection operations in nix_api_util.h.
package util

/*
#cgo pkg-config: nix-util-c
#include <nix_api_util.h>
*/
import "C"

// ErrorInfoMsg retrieves the structured Nix error message stored in readContext.
func ErrorInfoMsg(ctx, readContext *Context) (string, Err) {
	var out string
	ud := RegisterStringSlot(&out)
	err := Err(C.nix_err_info_msg(ctx.Ptr(), readContext.Ptr(),
		C.nix_get_string_callback(C.cgoUtilStringCB), ud))
	return out, err
}

// ErrorName retrieves the structured Nix error name stored in readContext.
func ErrorName(ctx, readContext *Context) (string, Err) {
	var out string
	ud := RegisterStringSlot(&out)
	err := Err(C.nix_err_name(ctx.Ptr(), readContext.Ptr(),
		C.nix_get_string_callback(C.cgoUtilStringCB), ud))
	return out, err
}
