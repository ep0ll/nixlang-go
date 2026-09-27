//go:build !cgo

// Package store provides stubs when CGO is disabled.
package store

import (
	"sync"
	"unsafe"

	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

type Store struct{ once sync.Once }
type StorePath struct{ once sync.Once }
type Derivation struct{ once sync.Once }

func Init(ctx *cutil.Context) cutil.Err         { return cutil.OK }
func InitNoConfig(ctx *cutil.Context) cutil.Err { return cutil.OK }
func Open(ctx *cutil.Context, uri string, params map[string]string) *Store {
	return &Store{}
}
func (s *Store) Free() {
	if s != nil {
		s.once.Do(func() {})
	}
}
func (s *Store) Ptr() unsafe.Pointer { return nil }
func (s *Store) GetURI(ctx *cutil.Context) (string, cutil.Err) {
	return "stub://", cutil.OK
}
func (s *Store) GetStoreDir(ctx *cutil.Context) (string, cutil.Err) {
	return "/nix/store", cutil.OK
}
func (s *Store) GetVersion(ctx *cutil.Context) (string, cutil.Err) {
	return "", cutil.OK
}
func (s *Store) ParsePath(ctx *cutil.Context, path string) *StorePath { return nil }
func (s *Store) IsValidPath(ctx *cutil.Context, path *StorePath) bool { return false }
func (s *Store) RealPath(ctx *cutil.Context, path *StorePath) (string, cutil.Err) {
	return "", cutil.OK
}
func (s *Store) Realise(ctx *cutil.Context, path *StorePath, cb func(outname string, out *StorePath)) cutil.Err {
	return cutil.OK
}
func CopyClosure(ctx *cutil.Context, src, dst *Store, path *StorePath) cutil.Err {
	return cutil.OK
}
func CopyPath(ctx *cutil.Context, src, dst *Store, path *StorePath, repair, checkSigs bool) cutil.Err {
	return cutil.OK
}
func (s *Store) QueryPathFromHashPart(ctx *cutil.Context, hash string) *StorePath {
	return nil
}
func (s *Store) DrvFromStorePath(ctx *cutil.Context, path *StorePath) *Derivation {
	return nil
}
func DerivationFromJSON(ctx *cutil.Context, s *Store, json string) *Derivation {
	return nil
}
func (s *Store) AddDerivation(ctx *cutil.Context, d *Derivation) *StorePath {
	return nil
}
func (p *StorePath) Free() {
	if p != nil {
		p.once.Do(func() {})
	}
}
func (p *StorePath) Ptr() unsafe.Pointer { return nil }
func (p *StorePath) Clone() *StorePath   { return nil }
func (p *StorePath) Name() string        { return "" }
func (p *StorePath) Hash(ctx *cutil.Context) ([20]byte, cutil.Err) {
	return [20]byte{}, cutil.OK
}
func CreateFromParts(ctx *cutil.Context, hash [20]byte, name string) *StorePath {
	return nil
}
func (d *Derivation) Free() {
	if d != nil {
		d.once.Do(func() {})
	}
}
func (d *Derivation) Ptr() unsafe.Pointer { return nil }
func (d *Derivation) Clone() *Derivation  { return nil }
func (d *Derivation) ToJSON(ctx *cutil.Context) (string, cutil.Err) {
	return "", cutil.OK
}
