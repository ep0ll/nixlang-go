//go:build !cgo

// Package expr provides stubs when CGO is disabled.
package expr

import (
	"unsafe"

	cstore "github.com/ep0ll/nixlang-go/internal/c/store"
	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

type EvalState struct{}
type Builder struct{}

func Init(ctx *cutil.Context) cutil.Err { return cutil.OK }

func NewBuilder(ctx *cutil.Context, store *cstore.Store) *Builder { return &Builder{} }
func (b *Builder) Load(ctx *cutil.Context) cutil.Err              { return cutil.OK }
func (b *Builder) SetLookupPath(ctx *cutil.Context, paths []string) cutil.Err {
	return cutil.OK
}
func (b *Builder) SetSetting(ctx *cutil.Context, key, value string) cutil.Err {
	return cutil.OK
}
func (b *Builder) Build(ctx *cutil.Context) *EvalState { return &EvalState{} }
func (b *Builder) Free()                               {}
func (b *Builder) Ptr() unsafe.Pointer                 { return nil }

func Create(ctx *cutil.Context, lookupPath []string, store *cstore.Store) *EvalState {
	return &EvalState{}
}
func (s *EvalState) Free()               {}
func (s *EvalState) Ptr() unsafe.Pointer { return nil }

func EvalFromString(ctx *cutil.Context, state *EvalState, expr, path string, value unsafe.Pointer) cutil.Err {
	return cutil.OK
}
func ValueCall(ctx *cutil.Context, state *EvalState, fn, arg, out unsafe.Pointer) cutil.Err {
	return cutil.OK
}
func ValueCallMulti(ctx *cutil.Context, state *EvalState, fn unsafe.Pointer, args []unsafe.Pointer, out unsafe.Pointer) cutil.Err {
	return cutil.OK
}
func ValueForce(ctx *cutil.Context, state *EvalState, value unsafe.Pointer) cutil.Err {
	return cutil.OK
}
func ValueForceDeep(ctx *cutil.Context, state *EvalState, value unsafe.Pointer) cutil.Err {
	return cutil.OK
}
