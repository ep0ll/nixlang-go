//go:build !cgo

// Package value provides stubs when CGO is disabled.
package value

import (
	"sync"
	"unsafe"

	cexpr "github.com/ep0ll/nixlang-go/internal/c/expr"
	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

type ValueType int

const (
	TypeThunk ValueType = iota
	TypeInt
	TypeFloat
	TypeBool
	TypeString
	TypePath
	TypeNull
	TypeAttrs
	TypeList
	TypeFunction
	TypeExternal
	TypeFailed
)

type Value struct {
	once sync.Once
}

type ListBuilder struct{}
type BindingsBuilder struct{}
type PrimOp struct{}
type RealisedString struct{}

func WrapPtr(p unsafe.Pointer) *Value {
	if p == nil {
		return nil
	}
	return &Value{}
}

func Alloc(ctx *cutil.Context, state *cexpr.EvalState) *Value { return &Value{} }

func (v *Value) Ptr() unsafe.Pointer                 { return nil }
func (v *Value) CPtr() unsafe.Pointer                { return nil }
func (v *Value) Incref(ctx *cutil.Context) cutil.Err { return cutil.OK }
func (v *Value) Decref() {
	if v != nil {
		v.once.Do(func() {})
	}
}
func (v *Value) Type(ctx *cutil.Context) ValueType   { return TypeNull }
func (v *Value) GetBool(ctx *cutil.Context) bool     { return false }
func (v *Value) GetInt(ctx *cutil.Context) int64     { return 0 }
func (v *Value) GetFloat(ctx *cutil.Context) float64 { return 0 }
func (v *Value) GetString(ctx *cutil.Context) (string, cutil.Err) {
	return "", cutil.OK
}
func (v *Value) GetPathString(ctx *cutil.Context) string { return "" }
func (v *Value) ListSize(ctx *cutil.Context) uint        { return 0 }
func (v *Value) AttrsSize(ctx *cutil.Context) uint       { return 0 }
func (v *Value) ListByIdx(ctx *cutil.Context, state *cexpr.EvalState, idx uint) *Value {
	return nil
}
func (v *Value) ListByIdxLazy(ctx *cutil.Context, state *cexpr.EvalState, idx uint) *Value {
	return nil
}
func (v *Value) AttrByName(ctx *cutil.Context, state *cexpr.EvalState, name string) *Value {
	return nil
}
func (v *Value) AttrByNameLazy(ctx *cutil.Context, state *cexpr.EvalState, name string) *Value {
	return nil
}
func (v *Value) HasAttr(ctx *cutil.Context, state *cexpr.EvalState, name string) bool {
	return false
}
func (v *Value) AttrByIdx(ctx *cutil.Context, state *cexpr.EvalState, idx uint) (string, *Value) {
	return "", nil
}
func (v *Value) AttrNameByIdx(ctx *cutil.Context, state *cexpr.EvalState, idx uint) string {
	return ""
}
func (v *Value) InitNull(ctx *cutil.Context) cutil.Err             { return cutil.OK }
func (v *Value) InitBool(ctx *cutil.Context, b bool) cutil.Err     { return cutil.OK }
func (v *Value) InitInt(ctx *cutil.Context, i int64) cutil.Err     { return cutil.OK }
func (v *Value) InitFloat(ctx *cutil.Context, f float64) cutil.Err { return cutil.OK }
func (v *Value) InitString(ctx *cutil.Context, s string) cutil.Err { return cutil.OK }
func (v *Value) InitPathString(ctx *cutil.Context, state *cexpr.EvalState, path string) cutil.Err {
	return cutil.OK
}
func (v *Value) InitApply(ctx *cutil.Context, fn, arg *Value) cutil.Err { return cutil.OK }
func (v *Value) MakeList(ctx *cutil.Context, b *ListBuilder) cutil.Err  { return cutil.OK }
func (v *Value) MakeAttrs(ctx *cutil.Context, b *BindingsBuilder) cutil.Err {
	return cutil.OK
}
func (v *Value) GetTypename(ctx *cutil.Context) string { return "null" }
func (v *Value) CopyValue(ctx *cutil.Context, source *Value) cutil.Err {
	return cutil.OK
}

func MakeListBuilder(ctx *cutil.Context, state *cexpr.EvalState, capacity uint) *ListBuilder {
	return &ListBuilder{}
}
func (b *ListBuilder) Insert(ctx *cutil.Context, index uint, v *Value) cutil.Err {
	return cutil.OK
}
func (b *ListBuilder) Free() {}

func MakeBindingsBuilder(ctx *cutil.Context, state *cexpr.EvalState, capacity uint) *BindingsBuilder {
	return &BindingsBuilder{}
}
func (b *BindingsBuilder) Insert(ctx *cutil.Context, name string, v *Value) cutil.Err {
	return cutil.OK
}
func (b *BindingsBuilder) Free() {}

func StringRealise(ctx *cutil.Context, state *cexpr.EvalState, v *Value, isIFD bool) *RealisedString {
	return nil
}
func (r *RealisedString) Free()               {}
func (r *RealisedString) Buffer() string      { return "" }
func (r *RealisedString) StorePathCount() int { return 0 }
func (r *RealisedString) StorePathAt(index int) unsafe.Pointer {
	return nil
}
