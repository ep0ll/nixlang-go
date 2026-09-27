// Package util provides high-level utilities for the Nix C API:
// error contexts, settings, version, and verbosity.
package util

import (
	"errors"
	"fmt"

	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

// Context holds Nix error state for a sequence of API calls.
// Safe for sequential use from a single goroutine.
type Context struct {
	c *cutil.Context
}

// NewContext creates a new error context. Call Close when finished.
func NewContext() *Context {
	return &Context{c: cutil.NewContext()}
}

// Close frees the context. Idempotent.
func (ctx *Context) Close() {
	if ctx != nil && ctx.c != nil {
		ctx.c.Free()
		ctx.c = nil
	}
}

// Internal returns the low-level context for sibling packages.
func (ctx *Context) Internal() *cutil.Context {
	if ctx == nil {
		return nil
	}
	return ctx.c
}

// Code returns the last error code.
func (ctx *Context) Code() ErrorCode {
	if ctx == nil || ctx.c == nil {
		return OK
	}
	return ErrorCode(ctx.c.Code())
}

// Message returns the last error message.
func (ctx *Context) Message() string {
	if ctx == nil || ctx.c == nil {
		return ""
	}
	return ctx.c.Msg()
}

// Clear resets the error state.
func (ctx *Context) Clear() {
	if ctx != nil && ctx.c != nil {
		ctx.c.Clear()
	}
}

// Err converts the context into a Go error if the last call failed.
func (ctx *Context) Err() error {
	if ctx == nil || ctx.c == nil {
		return nil
	}
	code := ctx.c.Code()
	if code == cutil.OK {
		return nil
	}
	return &NixError{Code: ErrorCode(code), Msg: ctx.c.Msg()}
}

// Check returns an error if code is not OK.
func (ctx *Context) Check(code ErrorCode) error {
	if code == OK {
		return nil
	}
	if ctx != nil {
		if e := ctx.Err(); e != nil {
			return e
		}
	}
	return &NixError{Code: code, Msg: "nix error"}
}

// ErrorCode is a Nix C API error code.
type ErrorCode int

const (
	OK             ErrorCode = ErrorCode(cutil.OK)
	ErrUnknown     ErrorCode = ErrorCode(cutil.ErrUnknown)
	ErrOverflow    ErrorCode = ErrorCode(cutil.ErrOverflow)
	ErrKey         ErrorCode = ErrorCode(cutil.ErrKey)
	ErrNixError    ErrorCode = ErrorCode(cutil.ErrNixError)
	ErrRecoverable ErrorCode = ErrorCode(cutil.ErrRecoverable)
)

// NixError is a structured Nix API error.
type NixError struct {
	Code ErrorCode
	Msg  string
}

func (e *NixError) Error() string {
	if e.Msg != "" {
		return fmt.Sprintf("nix: %s (code %d)", e.Msg, e.Code)
	}
	return fmt.Sprintf("nix: error code %d", e.Code)
}

// Is reports whether target is a *NixError with the same Code.
// Enables errors.Is(err, &NixError{Code: util.ErrKey}).
func (e *NixError) Is(target error) bool {
	t, ok := target.(*NixError)
	if !ok || e == nil || t == nil {
		return false
	}
	return e.Code == t.Code
}

// AsNixError extracts *NixError from err via errors.As.
func AsNixError(err error) (*NixError, bool) {
	var ne *NixError
	if errors.As(err, &ne) {
		return ne, true
	}
	return nil, false
}

// IsRecoverable reports whether err is NIX_ERR_RECOVERABLE
// (e.g. a failed primop that the evaluator may retry).
func IsRecoverable(err error) bool {
	ne, ok := AsNixError(err)
	return ok && ne.Code == ErrRecoverable
}

// IsKeyError reports whether err is NIX_ERR_KEY
// (missing attr, out-of-range index, unknown setting).
func IsKeyError(err error) bool {
	ne, ok := AsNixError(err)
	return ok && ne.Code == ErrKey
}

// Init initializes libutil. Must be called before other Nix API use.
func Init(ctx *Context) error {
	return ctx.Check(ErrorCode(cutil.Init(ctx.Internal())))
}

// Version returns the linked Nix library version.
func Version() string {
	return cutil.Version()
}

// SetSetting sets a Nix configuration setting.
func SetSetting(ctx *Context, key, value string) error {
	return ctx.Check(ErrorCode(cutil.SettingSet(ctx.Internal(), key, value)))
}

// GetSetting retrieves a Nix configuration setting.
func GetSetting(ctx *Context, key string) (string, error) {
	v, code := cutil.SettingGet(ctx.Internal(), key)
	if err := ctx.Check(ErrorCode(code)); err != nil {
		return "", err
	}
	return v, nil
}

// SetVerbosity sets the global log verbosity.
func SetVerbosity(ctx *Context, level Verbosity) error {
	return ctx.Check(ErrorCode(cutil.SetVerbosity(ctx.Internal(), cutil.Verbosity(level))))
}

// Verbosity is a log verbosity level.
type Verbosity int

const (
	LvlError     Verbosity = Verbosity(cutil.LvlError)
	LvlWarn      Verbosity = Verbosity(cutil.LvlWarn)
	LvlNotice    Verbosity = Verbosity(cutil.LvlNotice)
	LvlInfo      Verbosity = Verbosity(cutil.LvlInfo)
	LvlTalkative Verbosity = Verbosity(cutil.LvlTalkative)
	LvlChatty    Verbosity = Verbosity(cutil.LvlChatty)
	LvlDebug     Verbosity = Verbosity(cutil.LvlDebug)
	LvlVomit     Verbosity = Verbosity(cutil.LvlVomit)
)
