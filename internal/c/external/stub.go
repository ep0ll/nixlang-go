//go:build !cgo

package external

import (
	"unsafe"

	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

type Desc struct{}
type Value struct{}

func Create(ctx *cutil.Context, desc *Desc, data any) *Value { return nil }
func (v *Value) Free()                                      {}
func (v *Value) Content(ctx *cutil.Context) any             { return nil }
func (v *Value) Ptr() unsafe.Pointer                        { return nil }
