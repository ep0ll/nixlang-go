// Package fetchers provides fetcher settings for Nix (used by flakes and fetch builtins).
package fetchers

import (
	"fmt"

	cfetch "github.com/ep0ll/nixlang-go/internal/c/fetchers"
	"github.com/ep0ll/nixlang-go/util"
)

// Settings holds shared fetcher configuration.
type Settings struct {
	c *cfetch.Settings
}

// New creates default fetcher settings.
func New(ctx *util.Context) (*Settings, error) {
	s := cfetch.New(ctx.Internal())
	if s == nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to create fetchers settings")
	}
	return &Settings{c: s}, nil
}

// Close frees the settings.
func (s *Settings) Close() {
	if s != nil && s.c != nil {
		s.c.Free()
		s.c = nil
	}
}

// Internal returns the low-level handle.
func (s *Settings) Internal() *cfetch.Settings {
	if s == nil {
		return nil
	}
	return s.c
}
