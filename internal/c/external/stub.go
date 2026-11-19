//go:build !cgo

package external

import (
	"unsafe"

	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

type Callbacks struct {
	Print            func(data any, emit func(string) error)
	ShowType         func(data any) string
	TypeOf           func(data any) string
	CoerceToString   func(data any, coerceMore, copyToStore bool, addContext func(string) error) (string, error)
	Equal            func(left, right any) bool
	PrintValueAsJSON func(data any, state unsafe.Pointer, strict, copyToStore bool, addContext func(string) error) (string, error)
	PrintValueAsXML  func(data any, state unsafe.Pointer, strict, location bool, doc, stringContext, drvsSeen unsafe.Pointer, pos int, addContext func(string) error)
}

type Desc struct{}
type Value struct{}

func NewDesc(Callbacks) *Desc                                { return &Desc{} }
func (d *Desc) Close()                                       {}
func Create(ctx *cutil.Context, desc *Desc, data any) *Value { return nil }
func Wrap(ctx *cutil.Context, ptr unsafe.Pointer) *Value     { return nil }
func (v *Value) Free()                                       {}
func (v *Value) Content(ctx *cutil.Context) any              { return nil }
func (v *Value) Ptr() unsafe.Pointer                         { return nil }
