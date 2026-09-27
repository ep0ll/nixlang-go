// Package store provides low-level cgo bindings for nix_api_store.h and related headers.
package store

/*
#cgo pkg-config: nix-store-c nix-util-c
#include <nix_api_store.h>
#include <nix_api_util.h>
#include <stdlib.h>
#include <string.h>

typedef void (*nixgo_realise_cb)(void *userdata, const char *outname, const StorePath *out);
extern void cgoRealiseCB(void *userdata, const char *outname, const StorePath *out);
*/
import "C"
import (
	"runtime"
	"unsafe"

	cutil "github.com/ep0ll/nixlang-go/internal/c/util"
)

// Store wraps a Nix store handle.
type Store struct {
	ptr *C.Store
}

// StorePath wraps a store path.
type StorePath struct {
	ptr *C.StorePath
}

// Derivation wraps nix_derivation.
type Derivation struct {
	ptr *C.nix_derivation
}

// Init initializes libstore.
func Init(ctx *cutil.Context) cutil.Err {
	return cutil.FromCErr(C.nix_libstore_init(ctx.Ptr()))
}

// InitNoConfig initializes libstore without loading Nix config.
func InitNoConfig(ctx *cutil.Context) cutil.Err {
	return cutil.FromCErr(C.nix_libstore_init_no_load_config(ctx.Ptr()))
}

// Open opens a store.
// uri may be "" for the default store.
// params is currently reserved; pass nil. Prefer util.SetSetting / store URI query params.
func Open(ctx *cutil.Context, uri string, params map[string]string) *Store {
	var curi *C.char
	if uri != "" {
		curi = C.CString(uri)
		defer C.free(unsafe.Pointer(curi))
	}
	// Store open params (const char ***) are version-sensitive; use URI/settings instead.
	_ = params
	p := C.nix_store_open(ctx.Ptr(), curi, nil)
	if p == nil {
		return nil
	}
	s := &Store{ptr: p}
	runtime.SetFinalizer(s, func(s *Store) { s.Free() })
	return s
}

// Free releases the store.
func (s *Store) Free() {
	if s != nil && s.ptr != nil {
		C.nix_store_free(s.ptr)
		s.ptr = nil
		runtime.SetFinalizer(s, nil)
	}
}

// Ptr returns the C pointer.
func (s *Store) Ptr() *C.Store {
	if s == nil {
		return nil
	}
	return s.ptr
}

func stringViaCB(ctx *cutil.Context, call func(cb C.nix_get_string_callback, ud unsafe.Pointer) C.nix_err) (string, cutil.Err) {
	var out string
	ud := cutil.RegisterStringSlot(&out)
	err := cutil.FromCErr(call(C.nix_get_string_callback(C.cgoStoreStringCB), ud))
	return out, err
}

// GetURI returns the store URI.
func (s *Store) GetURI(ctx *cutil.Context) (string, cutil.Err) {
	return stringViaCB(ctx, func(cb C.nix_get_string_callback, ud unsafe.Pointer) C.nix_err {
		return C.nix_store_get_uri(ctx.Ptr(), s.ptr, cb, ud)
	})
}

// GetStoreDir returns the store directory.
func (s *Store) GetStoreDir(ctx *cutil.Context) (string, cutil.Err) {
	return stringViaCB(ctx, func(cb C.nix_get_string_callback, ud unsafe.Pointer) C.nix_err {
		return C.nix_store_get_storedir(ctx.Ptr(), s.ptr, cb, ud)
	})
}

// GetVersion returns the store version string.
func (s *Store) GetVersion(ctx *cutil.Context) (string, cutil.Err) {
	return stringViaCB(ctx, func(cb C.nix_get_string_callback, ud unsafe.Pointer) C.nix_err {
		return C.nix_store_get_version(ctx.Ptr(), s.ptr, cb, ud)
	})
}

// ParsePath parses a full store path string.
func (s *Store) ParsePath(ctx *cutil.Context, path string) *StorePath {
	cp := C.CString(path)
	defer C.free(unsafe.Pointer(cp))
	p := C.nix_store_parse_path(ctx.Ptr(), s.ptr, cp)
	if p == nil {
		return nil
	}
	sp := &StorePath{ptr: p}
	runtime.SetFinalizer(sp, func(sp *StorePath) { sp.Free() })
	return sp
}

// IsValidPath reports whether the path exists in the store.
func (s *Store) IsValidPath(ctx *cutil.Context, path *StorePath) bool {
	return bool(C.nix_store_is_valid_path(ctx.Ptr(), s.ptr, path.ptr))
}

// RealPath returns the physical location of a store path.
func (s *Store) RealPath(ctx *cutil.Context, path *StorePath) (string, cutil.Err) {
	return stringViaCB(ctx, func(cb C.nix_get_string_callback, ud unsafe.Pointer) C.nix_err {
		return C.nix_store_real_path(ctx.Ptr(), s.ptr, path.ptr, cb, ud)
	})
}

// Realise builds/realises a store path. callback is invoked per output.
func (s *Store) Realise(ctx *cutil.Context, path *StorePath, cb func(outname string, out *StorePath)) cutil.Err {
	reg := registerRealiseCB(cb)
	defer unregisterRealiseCB(reg)
	return cutil.FromCErr(C.nix_store_realise(ctx.Ptr(), s.ptr, path.ptr, unsafe.Pointer(reg),
		C.nixgo_realise_cb(C.cgoRealiseCB)))
}

// CopyClosure copies the closure of path from src to dst.
func CopyClosure(ctx *cutil.Context, src, dst *Store, path *StorePath) cutil.Err {
	return cutil.FromCErr(C.nix_store_copy_closure(ctx.Ptr(), src.ptr, dst.ptr, path.ptr))
}

// CopyPath copies a single path between stores.
func CopyPath(ctx *cutil.Context, src, dst *Store, path *StorePath, repair, checkSigs bool) cutil.Err {
	return cutil.FromCErr(C.nix_store_copy_path(ctx.Ptr(), src.ptr, dst.ptr, path.ptr, C.bool(repair), C.bool(checkSigs)))
}

// QueryPathFromHashPart looks up a path by hash part.
func (s *Store) QueryPathFromHashPart(ctx *cutil.Context, hash string) *StorePath {
	ch := C.CString(hash)
	defer C.free(unsafe.Pointer(ch))
	p := C.nix_store_query_path_from_hash_part(ctx.Ptr(), s.ptr, ch)
	if p == nil {
		return nil
	}
	sp := &StorePath{ptr: p}
	runtime.SetFinalizer(sp, func(sp *StorePath) { sp.Free() })
	return sp
}

// DrvFromStorePath returns the derivation for a .drv path.
func (s *Store) DrvFromStorePath(ctx *cutil.Context, path *StorePath) *Derivation {
	p := C.nix_store_drv_from_store_path(ctx.Ptr(), s.ptr, path.ptr)
	if p == nil {
		return nil
	}
	d := &Derivation{ptr: p}
	runtime.SetFinalizer(d, func(d *Derivation) { d.Free() })
	return d
}

// DerivationFromJSON creates a derivation from JSON.
func DerivationFromJSON(ctx *cutil.Context, s *Store, json string) *Derivation {
	cj := C.CString(json)
	defer C.free(unsafe.Pointer(cj))
	p := C.nix_derivation_from_json(ctx.Ptr(), s.ptr, cj)
	if p == nil {
		return nil
	}
	d := &Derivation{ptr: p}
	runtime.SetFinalizer(d, func(d *Derivation) { d.Free() })
	return d
}

// AddDerivation inserts a derivation into the store, returning its path.
func (s *Store) AddDerivation(ctx *cutil.Context, d *Derivation) *StorePath {
	p := C.nix_add_derivation(ctx.Ptr(), s.ptr, d.ptr)
	if p == nil {
		return nil
	}
	sp := &StorePath{ptr: p}
	runtime.SetFinalizer(sp, func(sp *StorePath) { sp.Free() })
	return sp
}

// --- StorePath ---

// Free releases a StorePath.
func (p *StorePath) Free() {
	if p != nil && p.ptr != nil {
		C.nix_store_path_free(p.ptr)
		p.ptr = nil
		runtime.SetFinalizer(p, nil)
	}
}

// Ptr returns the C pointer.
func (p *StorePath) Ptr() *C.StorePath {
	if p == nil {
		return nil
	}
	return p.ptr
}

// Clone copies a StorePath.
func (p *StorePath) Clone() *StorePath {
	if p == nil || p.ptr == nil {
		return nil
	}
	cp := C.nix_store_path_clone(p.ptr)
	if cp == nil {
		return nil
	}
	sp := &StorePath{ptr: cp}
	runtime.SetFinalizer(sp, func(sp *StorePath) { sp.Free() })
	return sp
}

// Name returns the path name component.
func (p *StorePath) Name() string {
	var out string
	ud := cutil.RegisterStringSlot(&out)
	C.nix_store_path_name(p.ptr, C.nix_get_string_callback(C.cgoStoreStringCB), ud)
	return out
}

// Hash returns the 20-byte hash part.
func (p *StorePath) Hash(ctx *cutil.Context) ([20]byte, cutil.Err) {
	var hp C.nix_store_path_hash_part
	err := cutil.FromCErr(C.nix_store_path_hash(ctx.Ptr(), p.ptr, &hp))
	var out [20]byte
	for i := 0; i < 20; i++ {
		out[i] = byte(hp.bytes[i])
	}
	return out, err
}

// CreateFromParts builds a StorePath from hash and name.
func CreateFromParts(ctx *cutil.Context, hash [20]byte, name string) *StorePath {
	var hp C.nix_store_path_hash_part
	for i := 0; i < 20; i++ {
		hp.bytes[i] = C.uint8_t(hash[i])
	}
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	p := C.nix_store_create_from_parts(ctx.Ptr(), &hp, cn, C.size_t(len(name)))
	if p == nil {
		return nil
	}
	sp := &StorePath{ptr: p}
	runtime.SetFinalizer(sp, func(sp *StorePath) { sp.Free() })
	return sp
}

// --- Derivation ---

func (d *Derivation) Free() {
	if d != nil && d.ptr != nil {
		C.nix_derivation_free(d.ptr)
		d.ptr = nil
		runtime.SetFinalizer(d, nil)
	}
}

func (d *Derivation) Ptr() *C.nix_derivation {
	if d == nil {
		return nil
	}
	return d.ptr
}

func (d *Derivation) Clone() *Derivation {
	if d == nil || d.ptr == nil {
		return nil
	}
	p := C.nix_derivation_clone(d.ptr)
	if p == nil {
		return nil
	}
	nd := &Derivation{ptr: p}
	runtime.SetFinalizer(nd, func(nd *Derivation) { nd.Free() })
	return nd
}

func (d *Derivation) ToJSON(ctx *cutil.Context) (string, cutil.Err) {
	return stringViaCB(ctx, func(cb C.nix_get_string_callback, ud unsafe.Pointer) C.nix_err {
		return C.nix_derivation_to_json(ctx.Ptr(), d.ptr, cb, ud)
	})
}
