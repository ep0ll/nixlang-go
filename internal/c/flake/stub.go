//go:build !cgo

package flake

import (
	"unsafe"

	cexpr "github.com/ep0ll/nixlang-go/internal/c/expr"
	cfetch "github.com/ep0ll/nixlang-go/internal/c/fetchers"
	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

type Settings struct{}
type Reference struct{}
type LockFlags struct{}
type LockedFlake struct{}
type ParseFlags struct{}

func NewSettings(ctx *cutil.Context) *Settings { return &Settings{} }
func (s *Settings) Free()                      {}
func (s *Settings) Ptr() unsafe.Pointer        { return nil }
func (s *Settings) AddToEvalStateBuilder(ctx *cutil.Context, builder *cexpr.Builder) cutil.Err {
	return cutil.OK
}

func NewParseFlags(ctx *cutil.Context, settings *Settings) *ParseFlags {
	return &ParseFlags{}
}
func (f *ParseFlags) Free() {}
func (f *ParseFlags) SetBaseDirectory(ctx *cutil.Context, dir string) cutil.Err {
	return cutil.OK
}

func NewLockFlags(ctx *cutil.Context, settings *Settings) *LockFlags {
	return &LockFlags{}
}
func (f *LockFlags) Free() {}
func (f *LockFlags) SetModeCheck(ctx *cutil.Context) cutil.Err         { return cutil.OK }
func (f *LockFlags) SetModeVirtual(ctx *cutil.Context) cutil.Err       { return cutil.OK }
func (f *LockFlags) SetModeWriteAsNeeded(ctx *cutil.Context) cutil.Err { return cutil.OK }
func (f *LockFlags) AddInputOverride(ctx *cutil.Context, inputPath string, ref *Reference) cutil.Err {
	return cutil.OK
}

func Lock(
	ctx *cutil.Context,
	fetchSettings *cfetch.Settings,
	settings *Settings,
	state *cexpr.EvalState,
	flags *LockFlags,
	ref *Reference,
) *LockedFlake {
	return nil
}
func (lf *LockedFlake) Free() {}
func (lf *LockedFlake) GetOutputAttrs(ctx *cutil.Context, settings *Settings, state *cexpr.EvalState) unsafe.Pointer {
	return nil
}

func ParseReference(
	ctx *cutil.Context,
	fetchSettings *cfetch.Settings,
	settings *Settings,
	flags *ParseFlags,
	str string,
) (*Reference, string, cutil.Err) {
	return &Reference{}, "", cutil.OK
}
func (r *Reference) Free()               {}
func (r *Reference) Ptr() unsafe.Pointer { return nil }
