//go:build !cgo

package nixmain

import cutil "github.com/ep0ll/nixlang-go/internal/c/util"

func InitPlugins(ctx *cutil.Context) cutil.Err { return cutil.OK }
func SetLogFormat(ctx *cutil.Context, format string) cutil.Err {
	return cutil.OK
}
