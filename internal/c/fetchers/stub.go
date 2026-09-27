//go:build !cgo

package fetchers

import (
	"unsafe"

	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

type Settings struct{}

func New(ctx *cutil.Context) *Settings  { return &Settings{} }
func (s *Settings) Free()               {}
func (s *Settings) Ptr() unsafe.Pointer { return nil }
