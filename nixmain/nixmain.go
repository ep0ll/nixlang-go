// Package nixmain provides high-level access to libmain utilities:
// plugin loading and log format configuration.
//
// Named nixmain (not "main") to avoid collision with package main binaries.
package nixmain

import (
	cnixmain "github.com/ep0ll/nixlang-go/internal/c/nixmain"
	"github.com/ep0ll/nixlang-go/util"
)

// InitPlugins loads plugins listed in the plugin-files setting.
// Call once after init functions and settings are applied.
func InitPlugins(ctx *util.Context) error {
	code := cnixmain.InitPlugins(ctx.Internal())
	return ctx.Check(util.ErrorCode(code))
}

// SetLogFormat sets the Nix log format (e.g. "raw", "bar", "internal-json").
func SetLogFormat(ctx *util.Context, format string) error {
	code := cnixmain.SetLogFormat(ctx.Internal(), format)
	return ctx.Check(util.ErrorCode(code))
}
