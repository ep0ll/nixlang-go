package store

import cstore "github.com/ep0ll/nixlang-go/internal/c/store"

// PathFromInternal wraps an owned low-level StorePath returned by another
// package in this module. Ownership transfers to the returned Path.
func PathFromInternal(path *cstore.StorePath) *Path {
	if path == nil {
		return nil
	}
	return &Path{c: path}
}
