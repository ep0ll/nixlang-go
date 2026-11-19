// Package value — extended type/copy/string-realise helpers.
package value

import (
	"fmt"

	cstore "github.com/ep0ll/nixlang-go/internal/c/store"
	cvalue "github.com/ep0ll/nixlang-go/internal/c/value"
	nixstore "github.com/ep0ll/nixlang-go/store"
	"github.com/ep0ll/nixlang-go/util"
)

// TypeName returns the Nix type name (e.g. "int", "attrs", "string").
func (v *Value) TypeName() string {
	name, err := v.TypeNameChecked()
	if err != nil {
		return ""
	}
	return name
}

// TypeNameChecked returns the Nix runtime type name and propagates C API errors.
func (v *Value) TypeNameChecked() (string, error) {
	if err := v.ensureOpen(); err != nil {
		return "", err
	}
	name := v.c.GetTypename(v.ctx.Internal())
	if err := v.ctx.Err(); err != nil {
		return "", err
	}
	if name == "" {
		return "", fmt.Errorf("nix: failed to obtain value type name")
	}
	return name, nil
}

// CopyFrom copies source into v (both must be allocated on the same state).
func (v *Value) CopyFrom(source *Value) error {
	if err := v.ensureOpen(); err != nil {
		return err
	}
	if err := source.ensureOpen(); err != nil {
		return err
	}
	if v.state == nil || source.state == nil || !v.state.Same(source.state) {
		return fmt.Errorf("nix: values belong to different evaluator states")
	}
	code := v.c.CopyValue(v.ctx.Internal(), source.c)
	return v.ctx.Check(util.ErrorCode(code))
}

// RealiseString forces a string and returns its content plus store-path context count.
// The returned buffer is the realised string; pathCount is the number of context paths.
func (v *Value) RealiseString(allowIFD bool) (buf string, pathCount int, err error) {
	if err := v.ensureOpen(); err != nil {
		return "", 0, err
	}
	rs := cvalue.StringRealise(v.ctx.Internal(), v.state, v.c, allowIFD)
	if rs == nil {
		if e := v.ctx.Err(); e != nil {
			return "", 0, e
		}
		return "", 0, fmt.Errorf("nix: string realise failed")
	}
	defer rs.Free()
	buf = rs.Buffer()
	count := rs.StorePathCount()
	if err := v.ctx.Err(); err != nil {
		return "", 0, err
	}
	return buf, count, nil
}

// RealisedString owns a realised Nix string and its store-path context.
// Call Close when finished; the finalizer is only a fallback.
type RealisedString struct {
	c *cvalue.RealisedString
}

// RealiseStringWithContext realises a string while preserving its referenced
// store paths. Use StorePathAt or StorePaths to inspect those references.
func (v *Value) RealiseStringWithContext(allowIFD bool) (*RealisedString, error) {
	if err := v.ensureOpen(); err != nil {
		return nil, err
	}
	realised := cvalue.StringRealise(v.ctx.Internal(), v.state, v.c, allowIFD)
	if realised == nil {
		if err := v.ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: string realise failed")
	}
	return &RealisedString{c: realised}, nil
}

// Close releases the realised string. It is idempotent.
func (r *RealisedString) Close() {
	if r != nil && r.c != nil {
		r.c.Free()
		r.c = nil
	}
}

// String returns the realised string buffer.
func (r *RealisedString) String() string {
	if r == nil || r.c == nil {
		return ""
	}
	return r.c.Buffer()
}

// StorePathCount returns the number of store paths in the string context.
func (r *RealisedString) StorePathCount() int {
	if r == nil || r.c == nil {
		return 0
	}
	return r.c.StorePathCount()
}

// StorePathAt returns an owned clone of a referenced store path. The caller
// must close the returned path.
func (r *RealisedString) StorePathAt(index int) (*nixstore.Path, error) {
	if r == nil || r.c == nil {
		return nil, fmt.Errorf("nix: realised string is closed")
	}
	if index < 0 || index >= r.c.StorePathCount() {
		return nil, fmt.Errorf("nix: realised string store-path index %d out of range", index)
	}
	cloned := cstore.ClonePath(r.c.StorePathAt(index))
	if cloned == nil {
		return nil, fmt.Errorf("nix: failed to clone realised string store path %d", index)
	}
	return nixstore.PathFromInternal(cloned), nil
}

// StorePaths returns owned clones of every store path referenced by the string.
// The caller must close every returned path.
func (r *RealisedString) StorePaths() ([]*nixstore.Path, error) {
	if r == nil || r.c == nil {
		return nil, fmt.Errorf("nix: realised string is closed")
	}
	count := r.c.StorePathCount()
	paths := make([]*nixstore.Path, 0, count)
	for i := 0; i < count; i++ {
		path, err := r.StorePathAt(i)
		if err != nil {
			for _, opened := range paths {
				opened.Close()
			}
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}

// Retain creates an independently closeable reference to the same Nix value.
// The returned value belongs to the same evaluator state and must be closed
// separately. Retain is useful when a value must outlive the wrapper that
// originally returned it.
func (v *Value) Retain() (*Value, error) {
	if err := v.ensureOpen(); err != nil {
		return nil, err
	}
	retained, code := v.c.Retain(v.ctx.Internal())
	if err := v.ctx.Check(util.ErrorCode(code)); err != nil {
		return nil, err
	}
	if retained == nil {
		return nil, fmt.Errorf("nix: failed to retain value")
	}
	return FromInternal(v.ctx, v.state, retained), nil
}
