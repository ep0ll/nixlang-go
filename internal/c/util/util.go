// Package util provides low-level cgo bindings for nix_api_util.h.
package util

/*
#cgo pkg-config: nix-util-c
#include <nix_api_util.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"sync"
	"unsafe"
)

// Err is the C API error code type.
type Err C.nix_err

const (
	OK             Err = C.NIX_OK
	ErrUnknown     Err = C.NIX_ERR_UNKNOWN
	ErrOverflow    Err = C.NIX_ERR_OVERFLOW
	ErrKey         Err = C.NIX_ERR_KEY
	ErrNixError    Err = C.NIX_ERR_NIX_ERROR
	ErrRecoverable Err = C.NIX_ERR_RECOVERABLE
)

// Verbosity mirrors nix_verbosity.
type Verbosity C.nix_verbosity

const (
	LvlError     Verbosity = C.NIX_LVL_ERROR
	LvlWarn      Verbosity = C.NIX_LVL_WARN
	LvlNotice    Verbosity = C.NIX_LVL_NOTICE
	LvlInfo      Verbosity = C.NIX_LVL_INFO
	LvlTalkative Verbosity = C.NIX_LVL_TALKATIVE
	LvlChatty    Verbosity = C.NIX_LVL_CHATTY
	LvlDebug     Verbosity = C.NIX_LVL_DEBUG
	LvlVomit     Verbosity = C.NIX_LVL_VOMIT
)

// Context wraps nix_c_context.
// Free is idempotent (sync.Once); always call Free/Close explicitly.
type Context struct {
	ptr  *C.nix_c_context
	once sync.Once
}

// NewContext allocates a new error context.
func NewContext() *Context {
	p := C.nix_c_context_create()
	c := &Context{ptr: p}
	runtime.SetFinalizer(c, func(c *Context) { c.Free() })
	return c
}

// Ptr returns the underlying C pointer (nil-safe).
func (c *Context) Ptr() *C.nix_c_context {
	if c == nil {
		return nil
	}
	return c.ptr
}

// Free releases the context. Idempotent; safe with concurrent finalizer.
func (c *Context) Free() {
	if c == nil {
		return
	}
	c.once.Do(func() {
		if c.ptr != nil {
			C.nix_c_context_free(c.ptr)
			c.ptr = nil
		}
		runtime.SetFinalizer(c, nil)
	})
}

// Code returns the last error code.
func (c *Context) Code() Err {
	if c == nil || c.ptr == nil {
		return OK
	}
	return Err(C.nix_err_code(c.ptr))
}

// Msg returns the last error message.
func (c *Context) Msg() string {
	if c == nil || c.ptr == nil {
		return ""
	}
	var n C.uint
	p := C.nix_err_msg(nil, c.ptr, &n)
	if p == nil {
		return ""
	}
	return C.GoStringN(p, C.int(n))
}

// Clear clears the error state.
func (c *Context) Clear() {
	if c != nil && c.ptr != nil {
		C.nix_clear_err(c.ptr)
	}
}

// SetErr sets an error message (e.g. from primop callbacks).
func (c *Context) SetErr(code Err, msg string) Err {
	cs := C.CString(msg)
	defer C.free(unsafe.Pointer(cs))
	return Err(C.nix_set_err_msg(c.Ptr(), C.nix_err(code), cs))
}

// Init initializes libutil.
func Init(ctx *Context) Err {
	return Err(C.nix_libutil_init(ctx.Ptr()))
}

// Version returns the Nix library version string.
func Version() string {
	return C.GoString(C.nix_version_get())
}

// SettingGet retrieves a configuration setting.
func SettingGet(ctx *Context, key string) (string, Err) {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	var out string
	ud := RegisterStringSlot(&out)
	err := Err(C.nix_setting_get(ctx.Ptr(), ckey, C.nix_get_string_callback(C.cgoUtilStringCB), ud))
	return out, err
}

// SettingSet sets a configuration setting.
func SettingSet(ctx *Context, key, value string) Err {
	ck := C.CString(key)
	cv := C.CString(value)
	defer C.free(unsafe.Pointer(ck))
	defer C.free(unsafe.Pointer(cv))
	return Err(C.nix_setting_set(ctx.Ptr(), ck, cv))
}

// SetVerbosity sets global verbosity.
func SetVerbosity(ctx *Context, level Verbosity) Err {
	return Err(C.nix_set_verbosity(ctx.Ptr(), C.nix_verbosity(level)))
}

// --- shared string callback slots ---

var (
	strMu    sync.Mutex
	strSlots = make(map[uintptr]*string)
	strNext  uintptr = 1
)

// RegisterStringSlot stores a *string and returns userdata for C callbacks.
func RegisterStringSlot(s *string) unsafe.Pointer {
	strMu.Lock()
	defer strMu.Unlock()
	id := strNext
	strNext++
	strSlots[id] = s
	return unsafe.Pointer(id)
}

// TakeStringSlot retrieves and removes a registered string slot.
func TakeStringSlot(p unsafe.Pointer) *string {
	strMu.Lock()
	defer strMu.Unlock()
	id := uintptr(p)
	s := strSlots[id]
	delete(strSlots, id)
	return s
}

// FromCErr converts C.nix_err to Err.
func FromCErr(e C.nix_err) Err { return Err(e) }

// CErr converts Err to C.nix_err.
func CErr(e Err) C.nix_err { return C.nix_err(e) }
