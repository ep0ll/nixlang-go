// Package store provides a high-level Go API for Nix stores.
//
// Open stores, inspect paths, realise outputs, and manipulate closures
// and derivations without invoking the Nix CLI.
package store

import (
	"fmt"

	cstore "github.com/ep0ll/nixlang-go/internal/c/store"
	"github.com/ep0ll/nixlang-go/util"
)

// Store is a handle to a Nix store (local, remote, dummy, etc.).
type Store struct {
	c   *cstore.Store
	ctx *util.Context
}

// Path is a Nix store path.
type Path struct {
	c *cstore.StorePath
}

// Derivation is a Nix derivation object.
type Derivation struct {
	c *cstore.Derivation
}

// Init initializes libstore. Call after util.Init.
func Init(ctx *util.Context) error {
	return ctx.Check(util.ErrorCode(cstore.Init(ctx.Internal())))
}

// InitNoConfig initializes libstore without loading configuration files.
func InitNoConfig(ctx *util.Context) error {
	return ctx.Check(util.ErrorCode(cstore.InitNoConfig(ctx.Internal())))
}

// Open opens a Nix store.
// uri examples: "", "auto", "daemon", "local", "dummy://", "ssh://host".
// params are optional store parameters (currently reserved; prefer URI query or util.SetSetting).
func Open(ctx *util.Context, uri string, params map[string]string) (*Store, error) {
	s := cstore.Open(ctx.Internal(), uri, params)
	if s == nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to open store %q", uri)
	}
	return &Store{c: s, ctx: ctx}, nil
}

// Close releases the store. Idempotent.
func (s *Store) Close() {
	if s != nil && s.c != nil {
		s.c.Free()
		s.c = nil
	}
}

// Internal returns the low-level store for sibling packages.
func (s *Store) Internal() *cstore.Store {
	if s == nil {
		return nil
	}
	return s.c
}

// URI returns the store URI.
func (s *Store) URI() (string, error) {
	v, code := s.c.GetURI(s.ctx.Internal())
	if err := s.ctx.Check(util.ErrorCode(code)); err != nil {
		return "", err
	}
	return v, nil
}

// StoreDir returns the logical store directory (e.g. /nix/store).
func (s *Store) StoreDir() (string, error) {
	v, code := s.c.GetStoreDir(s.ctx.Internal())
	if err := s.ctx.Check(util.ErrorCode(code)); err != nil {
		return "", err
	}
	return v, nil
}

// Version returns the store protocol/version string.
func (s *Store) Version() (string, error) {
	v, code := s.c.GetVersion(s.ctx.Internal())
	if err := s.ctx.Check(util.ErrorCode(code)); err != nil {
		return "", err
	}
	return v, nil
}

// ParsePath parses a full store path string into a Path.
func (s *Store) ParsePath(path string) (*Path, error) {
	p := s.c.ParsePath(s.ctx.Internal(), path)
	if p == nil {
		if err := s.ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: invalid store path %q", path)
	}
	return &Path{c: p}, nil
}

// IsValid reports whether the path exists in the store.
func (s *Store) IsValid(path *Path) bool {
	if path == nil || path.c == nil {
		return false
	}
	return s.c.IsValidPath(s.ctx.Internal(), path.c)
}

// RealPath returns the physical filesystem location of a store path.
func (s *Store) RealPath(path *Path) (string, error) {
	v, code := s.c.RealPath(s.ctx.Internal(), path.c)
	if err := s.ctx.Check(util.ErrorCode(code)); err != nil {
		return "", err
	}
	return v, nil
}

// Realise builds/realises a derivation or path.
// The callback receives each output name and path (borrowed for the callback duration).
func (s *Store) Realise(path *Path, cb func(outname string, out *Path)) error {
	code := s.c.Realise(s.ctx.Internal(), path.c, func(name string, sp *cstore.StorePath) {
		var p *Path
		if sp != nil {
			p = &Path{c: sp}
		}
		cb(name, p)
	})
	return s.ctx.Check(util.ErrorCode(code))
}

// QueryPathFromHashPart looks up a store path by its hash component.
func (s *Store) QueryPathFromHashPart(hash string) (*Path, error) {
	p := s.c.QueryPathFromHashPart(s.ctx.Internal(), hash)
	if p == nil {
		return nil, s.ctx.Err()
	}
	return &Path{c: p}, nil
}

// DerivationFromPath returns the derivation object for a .drv store path.
func (s *Store) DerivationFromPath(path *Path) (*Derivation, error) {
	d := s.c.DrvFromStorePath(s.ctx.Internal(), path.c)
	if d == nil {
		if err := s.ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: no derivation at path")
	}
	return &Derivation{c: d}, nil
}

// DerivationFromJSON creates a derivation from a JSON string.
func (s *Store) DerivationFromJSON(json string) (*Derivation, error) {
	d := cstore.DerivationFromJSON(s.ctx.Internal(), s.c, json)
	if d == nil {
		if err := s.ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to parse derivation JSON")
	}
	return &Derivation{c: d}, nil
}

// AddDerivation inserts a derivation into the store.
func (s *Store) AddDerivation(d *Derivation) (*Path, error) {
	p := s.c.AddDerivation(s.ctx.Internal(), d.c)
	if p == nil {
		if err := s.ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to add derivation")
	}
	return &Path{c: p}, nil
}

// Close frees the path.
func (p *Path) Close() {
	if p != nil && p.c != nil {
		p.c.Free()
		p.c = nil
	}
}

// Internal returns the low-level path.
func (p *Path) Internal() *cstore.StorePath {
	if p == nil {
		return nil
	}
	return p.c
}

// Name returns the name component of the store path.
func (p *Path) Name() string {
	return p.c.Name()
}

// Hash returns the 20-byte hash part of the store path.
func (p *Path) Hash(ctx *util.Context) ([20]byte, error) {
	h, code := p.c.Hash(ctx.Internal())
	if err := ctx.Check(util.ErrorCode(code)); err != nil {
		return [20]byte{}, err
	}
	return h, nil
}

// Clone returns an independent copy of the path.
func (p *Path) Clone() *Path {
	c := p.c.Clone()
	if c == nil {
		return nil
	}
	return &Path{c: c}
}

// CreatePathFromParts builds a Path from hash and name (no store required).
func CreatePathFromParts(ctx *util.Context, hash [20]byte, name string) (*Path, error) {
	p := cstore.CreateFromParts(ctx.Internal(), hash, name)
	if p == nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to create store path from parts")
	}
	return &Path{c: p}, nil
}

// Close frees the derivation.
func (d *Derivation) Close() {
	if d != nil && d.c != nil {
		d.c.Free()
		d.c = nil
	}
}

// ToJSON serializes the derivation to JSON.
func (d *Derivation) ToJSON(ctx *util.Context) (string, error) {
	s, code := d.c.ToJSON(ctx.Internal())
	if err := ctx.Check(util.ErrorCode(code)); err != nil {
		return "", err
	}
	return s, nil
}

// Clone returns an independent copy.
func (d *Derivation) Clone() *Derivation {
	c := d.c.Clone()
	if c == nil {
		return nil
	}
	return &Derivation{c: c}
}

// CopyClosure copies the closure of path from src to dst.
func CopyClosure(ctx *util.Context, src, dst *Store, path *Path) error {
	return ctx.Check(util.ErrorCode(cstore.CopyClosure(ctx.Internal(), src.c, dst.c, path.c)))
}

// CopyPath copies a single path between stores.
func CopyPath(ctx *util.Context, src, dst *Store, path *Path, repair, checkSigs bool) error {
	return ctx.Check(util.ErrorCode(cstore.CopyPath(ctx.Internal(), src.c, dst.c, path.c, repair, checkSigs)))
}
