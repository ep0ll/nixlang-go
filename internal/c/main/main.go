// Package main provides low-level cgo bindings for nix_api_main.h.
package main

/*
#cgo pkg-config: nix-main-c nix-util-c
#include <nix_api_main.h>
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

// InitPlugins loads plugins from the plugin-files setting.
func InitPlugins(ctx *cutil.Context) cutil.Err {
	return cutil.FromCErr(C.nix_init_plugins(ctx.Ptr()))
}

// SetLogFormat sets the log format by name.
func SetLogFormat(ctx *cutil.Context, format string) cutil.Err {
	cf := C.CString(format)
	defer C.free(unsafe.Pointer(cf))
	return cutil.FromCErr(C.nix_set_log_format(ctx.Ptr(), cf))
}
