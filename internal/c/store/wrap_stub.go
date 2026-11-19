//go:build !cgo

package store

import "unsafe"

// WrapPath is unavailable when CGO is disabled.
func WrapPath(unsafe.Pointer) *StorePath { return nil }


// ClonePath is unavailable when CGO is disabled.
func ClonePath(unsafe.Pointer) *StorePath { return nil }
