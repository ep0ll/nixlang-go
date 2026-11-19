// Package builtins provides idiomatic Go access to the Nix language builtins.
//
// It invokes the actual builtins in the Nix evaluator. The C API is used only
// underneath the existing expr/value packages for generic value construction
// and function application; callers pass ordinary Go values, never C pointers
// or C ABI types.
package builtins

import (
	"fmt"
	"reflect"

	"github.com/ep0ll/nixlang-go/expr"
	"github.com/ep0ll/nixlang-go/value"
)

// Path marks a Go string as a Nix path value. Plain Go strings become Nix
// strings; use Path when a builtin must receive the Nix path type.
type Path string

// Builtins is a handle to the builtins attribute set for one evaluator state.
// The evaluator state must outlive this handle and all values returned by it.
// Like expr.State, Builtins must not be used concurrently without external
// synchronization.
type Builtins struct {
	state  *expr.State
	set    *value.Value
	closed bool
}

// New obtains the evaluator's builtins attribute set. It does not create a
// separate evaluator or reimplement builtin semantics.
func New(state *expr.State) (*Builtins, error) {
	if state == nil {
		return nil, fmt.Errorf("nix: evaluator state is required")
	}
	set, err := state.EvalString("builtins", ".")
	if err != nil {
		return nil, fmt.Errorf("nix: obtain builtins attribute set: %w", err)
	}
	if set == nil {
		return nil, fmt.Errorf("nix: evaluator returned a nil builtins value")
	}
	return &Builtins{state: state, set: set}, nil
}

// Close releases this handle's reference to the builtins attribute set.
// Values returned by Get or Call are independently owned and remain valid
// until closed, provided the originating expr.State remains open.
func (b *Builtins) Close() {
	if b == nil || b.closed {
		return
	}
	b.closed = true
	if b.set != nil {
		b.set.Close()
		b.set = nil
	}
	b.state = nil
}

func (b *Builtins) ensureOpen() error {
	if b == nil || b.closed || b.state == nil || b.set == nil {
		return fmt.Errorf("nix: builtins handle is closed")
	}
	return nil
}

// Get returns a named value from the builtins set. The returned value is owned
// by the caller and must be closed. This also exposes constants such as
// currentTime, currentSystem, nixVersion, storeDir, langVersion and nixPath.
func (b *Builtins) Get(name string) (*value.Value, error) {
	if err := b.ensureOpen(); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, fmt.Errorf("nix: builtin name must not be empty")
	}
	v, err := b.set.Attr(name)
	if err != nil {
		return nil, fmt.Errorf("nix: builtin %q: %w", name, err)
	}
	return v, nil
}

// Has reports whether the running evaluator exposes name in builtins. This is
// useful for feature and Nix-version detection.
func (b *Builtins) Has(name string) (bool, error) {
	if err := b.ensureOpen(); err != nil {
		return false, err
	}
	return b.set.HasAttr(name)
}

// Names returns every attribute currently exposed by the running evaluator's
// builtins set, including constants and version-dependent entries.
func (b *Builtins) Names() ([]string, error) {
	if err := b.ensureOpen(); err != nil {
		return nil, err
	}
	return b.set.AttrNames()
}

// Call invokes a builtin with Go arguments and returns its Nix value.
//
// Supported Go arguments are nil (Nix null), bool, signed and unsigned
// integers that fit in int64, float32/float64, string (Nix string), Path (Nix
// path), slices/arrays (Nix lists), string-keyed maps (Nix attribute sets),
// and *value.Value from the same evaluator state. []byte is treated as a Nix
// string. Function-valued arguments should be supplied as *value.Value,
// typically produced by EvalString.
//
// Call intentionally permits partial application: passing fewer arguments
// than the builtin's arity returns the resulting Nix function value.
func (b *Builtins) Call(name string, args ...any) (*value.Value, error) {
	if err := b.ensureOpen(); err != nil {
		return nil, err
	}
	fn, err := b.Get(name)
	if err != nil {
		return nil, err
	}
	defer fn.Close()

	converted := make([]*value.Value, len(args))
	owned := make([]*value.Value, 0, len(args))
	defer func() {
		for _, v := range owned {
			v.Close()
		}
	}()

	for i, arg := range args {
		v, closeValue, err := b.fromGo(arg)
		if err != nil {
			return nil, fmt.Errorf("nix: builtin %q argument %d: %w", name, i+1, err)
		}
		converted[i] = v
		if closeValue {
			owned = append(owned, v)
		}
	}

	result, err := b.state.CallMulti(fn, converted...)
	if err != nil {
		return nil, fmt.Errorf("nix: call builtin %q: %w", name, err)
	}
	return result, nil
}

// CallInt invokes name and requires an integer result.
func (b *Builtins) CallInt(name string, args ...any) (int64, error) {
	v, err := b.Call(name, args...)
	if err != nil {
		return 0, err
	}
	defer v.Close()
	return v.Int()
}

// CallFloat invokes name and requires a floating-point result.
func (b *Builtins) CallFloat(name string, args ...any) (float64, error) {
	v, err := b.Call(name, args...)
	if err != nil {
		return 0, err
	}
	defer v.Close()
	return v.Float()
}

// CallBool invokes name and requires a Boolean result.
func (b *Builtins) CallBool(name string, args ...any) (bool, error) {
	v, err := b.Call(name, args...)
	if err != nil {
		return false, err
	}
	defer v.Close()
	return v.Bool()
}

// CallString invokes name and requires a string result.
func (b *Builtins) CallString(name string, args ...any) (string, error) {
	v, err := b.Call(name, args...)
	if err != nil {
		return "", err
	}
	defer v.Close()
	return v.String()
}

func (b *Builtins) fromGo(input any) (*value.Value, bool, error) {
	if existing, ok := input.(*value.Value); ok {
		if existing == nil || existing.Internal() == nil || existing.Internal().Ptr() == nil {
			return nil, false, fmt.Errorf("value argument is nil or closed")
		}
		if existing.StateInternal() == nil || b.state.Internal() == nil ||
			!existing.StateInternal().Same(b.state.Internal()) {
			return nil, false, fmt.Errorf("value belongs to a different evaluator state")
		}
		return existing, false, nil
	}

	out, err := b.state.AllocValue()
	if err != nil {
		return nil, false, err
	}
	fail := func(err error) (*value.Value, bool, error) {
		out.Close()
		return nil, false, err
	}

	switch x := input.(type) {
	case nil:
		err = out.SetNull()
	case Path:
		err = out.SetPath(string(x))
	case string:
		err = out.SetString(x)
	case []byte:
		err = out.SetString(string(x))
	default:
		rv := reflect.ValueOf(input)
		for rv.IsValid() && rv.Kind() == reflect.Interface {
			if rv.IsNil() {
				return fail(fmt.Errorf("nil interface is not a supported typed value; use nil for Nix null"))
			}
			rv = rv.Elem()
		}
		if !rv.IsValid() {
			err = out.SetNull()
			break
		}

		switch rv.Kind() {
		case reflect.Bool:
			err = out.SetBool(rv.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err = out.SetInt(rv.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n := rv.Uint()
			if n > uint64(^uint64(0)>>1) {
				return fail(fmt.Errorf("unsigned integer %d exceeds Nix's int64 range", n))
			}
			err = out.SetInt(int64(n))
		case reflect.Float32, reflect.Float64:
			err = out.SetFloat(rv.Float())
		case reflect.Slice, reflect.Array:
			if rv.Kind() == reflect.Slice && rv.Type().Elem().Kind() == reflect.Uint8 {
				err = out.SetString(string(rv.Bytes()))
				break
			}
			elems := make([]*value.Value, rv.Len())
			owned := make([]*value.Value, 0, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				if !rv.Index(i).CanInterface() {
					for _, child := range owned {
						child.Close()
					}
					return fail(fmt.Errorf("list element %d cannot be accessed", i))
				}
				child, closeChild, childErr := b.fromGo(rv.Index(i).Interface())
				if childErr != nil {
					for _, opened := range owned {
						opened.Close()
					}
					return fail(fmt.Errorf("list element %d: %w", i, childErr))
				}
				elems[i] = child
				if closeChild {
					owned = append(owned, child)
				}
			}
			err = out.SetList(elems)
			for _, child := range owned { child.Close() }
		case reflect.Map:
			if rv.Type().Key().Kind() != reflect.String {
				return fail(fmt.Errorf("attribute-set map keys must be strings, got %s", rv.Type().Key()))
			}
			attrs := make(map[string]*value.Value, rv.Len())
			owned := make([]*value.Value, 0, rv.Len())
			iter := rv.MapRange()
			for iter.Next() {
				key := iter.Key().String()
				if !iter.Value().CanInterface() {
					for _, child := range owned { child.Close() }
					return fail(fmt.Errorf("attribute %q cannot be accessed", key))
				}
				child, closeChild, childErr := b.fromGo(iter.Value().Interface())
				if childErr != nil {
					for _, opened := range owned { opened.Close() }
					return fail(fmt.Errorf("attribute %q: %w", key, childErr))
				}
				attrs[key] = child
				if closeChild { owned = append(owned, child) }
			}
			err = out.SetAttrs(attrs)
			for _, child := range owned { child.Close() }
		default:
			return fail(fmt.Errorf("unsupported Go argument type %T; use primitives, Path, slices, string-keyed maps, or *value.Value", input))
		}
	}

	if err != nil {
		return fail(err)
	}
	return out, true, nil
}
