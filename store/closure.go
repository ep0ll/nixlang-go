package store

import (
	"fmt"

	cstore "github.com/ep0ll/nixlang-go/internal/c/store"
	"github.com/ep0ll/nixlang-go/util"
)

// ClosureOptions controls the direction and members returned by GetFSClosure.
type ClosureOptions struct {
	// Reverse follows references to the input path instead of its dependencies.
	Reverse bool
	// IncludeOutputs adds outputs of derivations to a forward closure.
	IncludeOutputs bool
	// IncludeDerivers adds derivations that produced paths to a forward closure.
	IncludeDerivers bool
}

// GetFSClosure enumerates the filesystem closure of path. The visitor receives
// independently-owned paths and must close each path when finished. Iteration
// order is unspecified and the callback runs synchronously.
func (s *Store) GetFSClosure(path *Path, options ClosureOptions, visit func(*Path)) error {
	if s == nil || s.c == nil || path == nil || path.c == nil {
		return fmt.Errorf("nix: store and path are required")
	}
	return s.ctx.Check(util.ErrorCode(s.c.GetFSClosure(
		s.ctx.Internal(),
		path.c,
		options.Reverse,
		options.IncludeOutputs,
		options.IncludeDerivers,
		func(path *cstore.StorePath) {
			wrapped := PathFromInternal(path)
			if visit == nil {
				wrapped.Close()
				return
			}
			visit(wrapped)
		},
	)))
}
