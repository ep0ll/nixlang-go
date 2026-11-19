//go:build !cgo

package value

import (
	"unsafe"

	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

type PrimOpFunc func(userData any, ctx *cutil.Context, state interface{}, args []*Value, result *Value)

func AllocPrimOp(ctx *cutil.Context, fn PrimOpFunc, arity int, name string, argNames []string, doc string, userData any) *PrimOp {
	return nil
}

func (p *PrimOp) Register(ctx *cutil.Context) cutil.Err { return cutil.ErrUnknown }
func (p *PrimOp) Free()                                {}
func (p *PrimOp) Ptr() unsafe.Pointer                  { return nil }
