// Package fetchers provides low-level cgo bindings for nix_api_fetchers.h.
package fetchers

/*
#cgo pkg-config: nix-fetchers-c nix-util-c
#include <nix_api_fetchers.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"

	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

// Settings wraps nix_fetchers_settings.
type Settings struct {
	ptr *C.nix_fetchers_settings
}

// New creates default fetcher settings.
func New(ctx *cutil.Context) *Settings {
	p := C.nix_fetchers_settings_new(ctx.Ptr())
	if p == nil {
		return nil
	}
	s := &Settings{ptr: p}
	runtime.SetFinalizer(s, (*Settings).Free)
	return s
}

// Free releases the settings.
func (s *Settings) Free() {
	if s != nil && s.ptr != nil {
		C.nix_fetchers_settings_free(s.ptr)
		s.ptr = nil
		runtime.SetFinalizer(s, nil)
	}
}

// Ptr returns the C pointer.
func (s *Settings) Ptr() *C.nix_fetchers_settings {
	if s == nil {
		return nil
	}
	return s.ptr
}
