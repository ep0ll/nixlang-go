//go:build !cgo

package store

import cutil "github.com/ep0ll/nixlang-go/internal/c/util"

// GetFSClosure is unavailable when CGO is disabled.
func (s *Store) GetFSClosure(
	ctx *cutil.Context,
	path *StorePath,
	flipDirection, includeOutputs, includeDerivers bool,
	cb func(*StorePath),
) cutil.Err {
	return cutil.ErrUnknown
}
