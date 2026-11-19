//go:build !cgo

// Package util provides stubs when CGO is disabled so pure unit tests compile.
package util

import (
	"fmt"
	"sync"
	"unsafe"
)

type Err int

const (
	OK             Err = 0
	ErrUnknown     Err = -1
	ErrOverflow    Err = -2
	ErrKey         Err = -3
	ErrNixError    Err = -4
	ErrRecoverable Err = -5
)

type Verbosity int

const (
	LvlError Verbosity = iota
	LvlWarn
	LvlNotice
	LvlInfo
	LvlTalkative
	LvlChatty
	LvlDebug
	LvlVomit
)

type Context struct {
	once sync.Once
}

func NewContext() *Context { return &Context{} }

func WrapBorrowedContext(unsafe.Pointer) *Context { return &Context{} }

func (c *Context) Ptr() unsafe.Pointer { return nil }

func (c *Context) Free() {
	if c != nil {
		c.once.Do(func() {})
	}
}

func (c *Context) Code() Err { return OK }

func (c *Context) Msg() string { return "" }

func (c *Context) Clear() {}

func (c *Context) SetErr(code Err, msg string) Err { return code }

func Init(ctx *Context) Err { return OK }

func Version() string { return "stub-no-cgo" }

func SettingGet(ctx *Context, key string) (string, Err) {
	return "", ErrKey
}

func SettingSet(ctx *Context, key, value string) Err { return OK }

func SetVerbosity(ctx *Context, level Verbosity) Err { return OK }

var (
	strMu    sync.Mutex
	strSlots         = make(map[uintptr]*string)
	strNext  uintptr = 1
)

func RegisterStringSlot(s *string) unsafe.Pointer {
	strMu.Lock()
	defer strMu.Unlock()
	token := new(byte)
	strSlots[unsafe.Pointer(token)] = s
	return unsafe.Pointer(token)
}

func TakeStringSlot(p unsafe.Pointer) *string {
	strMu.Lock()
	defer strMu.Unlock()
	s, ok := strSlots[p]
	if ok {
		delete(strSlots, p)
	}
	return s
}

func FromCErr(e int) Err { return Err(e) }

func CErr(e Err) int { return int(e) }

var _ = fmt.Sprintf
