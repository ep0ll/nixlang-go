// Package value provides typed access to Nix language values.
package value

import (
	"fmt"
	"sort"

	cexpr "github.com/ep0ll/nixlang-go/internal/c/expr"
	cvalue "github.com/ep0ll/nixlang-go/internal/c/value"
	"github.com/ep0ll/nixlang-go/util"
)

// Type is the runtime type of a Nix value after forcing.
type Type int

const (
	TypeThunk    Type = Type(cvalue.TypeThunk)
	TypeInt      Type = Type(cvalue.TypeInt)
	TypeFloat    Type = Type(cvalue.TypeFloat)
	TypeBool     Type = Type(cvalue.TypeBool)
	TypeString   Type = Type(cvalue.TypeString)
	TypePath     Type = Type(cvalue.TypePath)
	TypeNull     Type = Type(cvalue.TypeNull)
	TypeAttrs    Type = Type(cvalue.TypeAttrs)
	TypeList     Type = Type(cvalue.TypeList)
	TypeFunction Type = Type(cvalue.TypeFunction)
	TypeExternal Type = Type(cvalue.TypeExternal)
	TypeFailed   Type = Type(cvalue.TypeFailed)
)

func (t Type) String() string {
	names := map[Type]string{
		TypeThunk: "thunk", TypeInt: "int", TypeFloat: "float", TypeBool: "bool",
		TypeString: "string", TypePath: "path", TypeNull: "null", TypeAttrs: "attrs",
		TypeList: "list", TypeFunction: "function", TypeExternal: "external", TypeFailed: "failed",
	}
	if s, ok := names[t]; ok {
		return s
	}
	return fmt.Sprintf("unknown(%d)", int(t))
}

// Value is a GC-managed Nix value. Call Close when done.
type Value struct {
	c           *cvalue.Value
	state       *cexpr.EvalState
	ctx         *util.Context
	ownsContext bool
}

// FromInternal wraps a low-level value (used by expr and flake packages).
func FromInternal(ctx *util.Context, state *cexpr.EvalState, cv *cvalue.Value) *Value {
	if cv == nil {
		return nil
	}
	return &Value{c: cv, state: state, ctx: ctx}
}

// Close decrements the GC reference. Idempotent.
func (v *Value) Close() {
	if v == nil {
		return
	}
	if v.c != nil {
		v.c.Decref()
		v.c = nil
	}
	if v.ownsContext && v.ctx != nil {
		v.ctx.Close()
	}
	v.ctx = nil
	v.state = nil
	v.ownsContext = false
}

// Internal returns the low-level value.
func (v *Value) Internal() *cvalue.Value {
	if v == nil {
		return nil
	}
	return v.c
}

// StateInternal returns the evaluator state associated with this value.
// It is intended for integration between the public subpackages.
func (v *Value) StateInternal() *cexpr.EvalState {
	if v == nil {
		return nil
	}
	return v.state
}

// Type returns the current type (call Force first for thunks).
func (v *Value) Type() Type {
	typ, err := v.TypeChecked()
	if err != nil {
		return TypeFailed
	}
	return typ
}

// TypeChecked returns the Nix runtime type and propagates C API errors.
func (v *Value) TypeChecked() (Type, error) {
	if err := v.ensureOpen(); err != nil {
		return TypeFailed, err
	}
	typ := Type(v.c.Type(v.ctx.Internal()))
	if err := v.ctx.Err(); err != nil {
		return TypeFailed, err
	}
	return typ, nil
}

// Force evaluates a thunk in place.
func (v *Value) Force() error {
	if err := v.ensureOpen(); err != nil {
		return err
	}
	code := cexpr.ValueForce(v.ctx.Internal(), v.state, v.c.Ptr())
	return v.ctx.Check(util.ErrorCode(code))
}

// ForceDeep deeply forces the value. Avoid on recursive data.
func (v *Value) ForceDeep() error {
	if err := v.ensureOpen(); err != nil {
		return err
	}
	code := cexpr.ValueForceDeep(v.ctx.Internal(), v.state, v.c.Ptr())
	return v.ctx.Check(util.ErrorCode(code))
}

// Bool returns the boolean value.
func (v *Value) Bool() (bool, error) {
	if err := v.ensureType(TypeBool); err != nil {
		return false, err
	}
	got := v.c.GetBool(v.ctx.Internal())
	if err := v.ctx.Err(); err != nil {
		return false, err
	}
	return got, nil
}

// Int returns the integer value.
func (v *Value) Int() (int64, error) {
	if err := v.ensureType(TypeInt); err != nil {
		return 0, err
	}
	got := v.c.GetInt(v.ctx.Internal())
	if err := v.ctx.Err(); err != nil {
		return 0, err
	}
	return got, nil
}

// Float returns the float value.
func (v *Value) Float() (float64, error) {
	if err := v.ensureType(TypeFloat); err != nil {
		return 0, err
	}
	got := v.c.GetFloat(v.ctx.Internal())
	if err := v.ctx.Err(); err != nil {
		return 0, err
	}
	return got, nil
}

// String returns the string contents.
func (v *Value) String() (string, error) {
	if err := v.ensureType(TypeString); err != nil {
		return "", err
	}
	s, code := v.c.GetString(v.ctx.Internal())
	if err := v.ctx.Check(util.ErrorCode(code)); err != nil {
		return "", err
	}
	return s, nil
}

// Path returns a path value as string.
func (v *Value) Path() (string, error) {
	if err := v.ensureType(TypePath); err != nil {
		return "", err
	}
	path := v.c.GetPathString(v.ctx.Internal())
	if err := v.ctx.Err(); err != nil {
		return "", err
	}
	return path, nil
}

// IsNull reports whether the value is null.
func (v *Value) IsNull() (bool, error) {
	if err := v.Force(); err != nil {
		return false, err
	}
	typ, err := v.TypeChecked()
	if err != nil {
		return false, err
	}
	return typ == TypeNull, nil
}

// ListLen returns the length of a list.
func (v *Value) ListLen() (int, error) {
	if err := v.ensureType(TypeList); err != nil {
		return 0, err
	}
	n := int(v.c.ListSize(v.ctx.Internal()))
	if err := v.ctx.Err(); err != nil {
		return 0, err
	}
	return n, nil
}

// ListAt returns the element at index i (forced).
func (v *Value) ListAt(i int) (*Value, error) {
	if i < 0 {
		return nil, fmt.Errorf("nix: list index must not be negative")
	}
	if err := v.ensureType(TypeList); err != nil {
		return nil, err
	}
	n := int(v.c.ListSize(v.ctx.Internal()))
	if err := v.ctx.Err(); err != nil {
		return nil, err
	}
	if i >= n {
		return nil, fmt.Errorf("nix: list index %d out of range", i)
	}
	cv := v.c.ListByIdx(v.ctx.Internal(), v.state, uint(i))
	if cv == nil {
		if err := v.ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: list index %d out of range", i)
	}
	return FromInternal(v.ctx, v.state, cv), nil
}

// ListAtLazy returns the element at index i without forcing it.
func (v *Value) ListAtLazy(i int) (*Value, error) {
	if i < 0 {
		return nil, fmt.Errorf("nix: list index must not be negative")
	}
	if err := v.ensureType(TypeList); err != nil {
		return nil, err
	}
	if i >= int(v.c.ListSize(v.ctx.Internal())) {
		return nil, fmt.Errorf("nix: list index %d out of range", i)
	}
	cv := v.c.ListByIdxLazy(v.ctx.Internal(), v.state, uint(i))
	if cv == nil {
		if err := v.ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: list index %d out of range", i)
	}
	return FromInternal(v.ctx, v.state, cv), nil
}

// AttrsLen returns the number of attributes.
func (v *Value) AttrsLen() (int, error) {
	if err := v.ensureType(TypeAttrs); err != nil {
		return 0, err
	}
	n := int(v.c.AttrsSize(v.ctx.Internal()))
	if err := v.ctx.Err(); err != nil {
		return 0, err
	}
	return n, nil
}

// Attr returns the attribute named name.
func (v *Value) Attr(name string) (*Value, error) {
	if err := v.ensureType(TypeAttrs); err != nil {
		return nil, err
	}
	cv := v.c.AttrByName(v.ctx.Internal(), v.state, name)
	if cv == nil {
		if err := v.ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: attribute %q not found", name)
	}
	return FromInternal(v.ctx, v.state, cv), nil
}

// AttrLazy returns the attribute without forcing its value.
func (v *Value) AttrLazy(name string) (*Value, error) {
	if err := v.ensureType(TypeAttrs); err != nil {
		return nil, err
	}
	cv := v.c.AttrByNameLazy(v.ctx.Internal(), v.state, name)
	if cv == nil {
		if err := v.ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: attribute %q not found", name)
	}
	return FromInternal(v.ctx, v.state, cv), nil
}

// AttrAt returns the name and value at the attribute index. Nix does not
// guarantee a stable attribute iteration order; sort names for deterministic output.
func (v *Value) AttrAt(index int) (string, *Value, error) {
	if index < 0 {
		return "", nil, fmt.Errorf("nix: attribute index must not be negative")
	}
	if err := v.ensureType(TypeAttrs); err != nil {
		return "", nil, err
	}
	n := int(v.c.AttrsSize(v.ctx.Internal()))
	if err := v.ctx.Err(); err != nil {
		return "", nil, err
	}
	if index >= n {
		return "", nil, fmt.Errorf("nix: attribute index %d out of range", index)
	}
	name, cv := v.c.AttrByIdx(v.ctx.Internal(), v.state, uint(index))
	if cv == nil {
		if err := v.ctx.Err(); err != nil {
			return "", nil, err
		}
		return "", nil, fmt.Errorf("nix: attribute index %d out of range", index)
	}
	return name, FromInternal(v.ctx, v.state, cv), nil
}

// AttrAtLazy returns the name and an unforced attribute value at the index.
// The attribute set must already be evaluated. Attribute ordering is unspecified.
func (v *Value) AttrAtLazy(index int) (string, *Value, error) {
	if index < 0 {
		return "", nil, fmt.Errorf("nix: attribute index must not be negative")
	}
	if err := v.ensureType(TypeAttrs); err != nil {
		return "", nil, err
	}
	if index >= int(v.c.AttrsSize(v.ctx.Internal())) {
		return "", nil, fmt.Errorf("nix: attribute index %d out of range", index)
	}
	name, cv := v.c.AttrByIdxLazy(v.ctx.Internal(), v.state, uint(index))
	if cv == nil {
		if err := v.ctx.Err(); err != nil {
			return "", nil, err
		}
		return "", nil, fmt.Errorf("nix: attribute index %d out of range", index)
	}
	return name, FromInternal(v.ctx, v.state, cv), nil
}

// HasAttr reports whether the attribute exists.
func (v *Value) HasAttr(name string) (bool, error) {
	if err := v.ensureType(TypeAttrs); err != nil {
		return false, err
	}
	exists := v.c.HasAttr(v.ctx.Internal(), v.state, name)
	if err := v.ctx.Err(); err != nil {
		return false, err
	}
	return exists, nil
}

// AttrNames returns all attribute names sorted lexicographically.
func (v *Value) AttrNames() ([]string, error) {
	n, err := v.AttrsLen()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, n)
	for i := 0; i < n; i++ {
		name := v.c.AttrNameByIdx(v.ctx.Internal(), v.state, uint(i))
		if err := v.ctx.Err(); err != nil {
			return nil, err
		}
		if name != "" {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names, nil
}

// Attrs returns a map of all attributes (forced). Caller must Close each value.
func (v *Value) Attrs() (map[string]*Value, error) {
	n, err := v.AttrsLen()
	if err != nil {
		return nil, err
	}
	m := make(map[string]*Value, n)
	for i := 0; i < n; i++ {
		name, cv := v.c.AttrByIdx(v.ctx.Internal(), v.state, uint(i))
		if cv == nil {
			if err := v.ctx.Err(); err != nil {
				for _, opened := range m {
					opened.Close()
				}
				return nil, err
			}
			continue
		}
		m[name] = FromInternal(v.ctx, v.state, cv)
	}
	return m, nil
}

func (v *Value) ensureOpen() error {
	if v == nil || v.c == nil || v.c.Ptr() == nil {
		return fmt.Errorf("nix: value is closed")
	}
	if v.ctx == nil || v.state == nil || v.state.Ptr() == nil {
		return fmt.Errorf("nix: value has no evaluation context")
	}
	return nil
}

func (v *Value) ensureType(want Type) error {
	if err := v.ensureOpen(); err != nil {
		return err
	}
	if err := v.Force(); err != nil {
		return err
	}
	got, err := v.TypeChecked()
	if err != nil {
		return err
	}
	if got != want {
		return fmt.Errorf("nix: expected %s, got %s", want, got)
	}
	return nil
}

// SetNull initializes this value to null.
func (v *Value) SetNull() error {
	if err := v.ensureOpen(); err != nil {
		return err
	}
	return v.ctx.Check(util.ErrorCode(v.c.InitNull(v.ctx.Internal())))
}

// SetBool initializes this value to a boolean.
func (v *Value) SetBool(b bool) error {
	if err := v.ensureOpen(); err != nil {
		return err
	}
	return v.ctx.Check(util.ErrorCode(v.c.InitBool(v.ctx.Internal(), b)))
}

// SetInt initializes this value to an integer.
func (v *Value) SetInt(i int64) error {
	if err := v.ensureOpen(); err != nil {
		return err
	}
	return v.ctx.Check(util.ErrorCode(v.c.InitInt(v.ctx.Internal(), i)))
}

// SetFloat initializes this value to a float.
func (v *Value) SetFloat(f float64) error {
	if err := v.ensureOpen(); err != nil {
		return err
	}
	return v.ctx.Check(util.ErrorCode(v.c.InitFloat(v.ctx.Internal(), f)))
}

// SetString initializes this value to a string.
func (v *Value) SetString(s string) error {
	if err := v.ensureOpen(); err != nil {
		return err
	}
	return v.ctx.Check(util.ErrorCode(v.c.InitString(v.ctx.Internal(), s)))
}

// SetPath initializes this value to a path.
func (v *Value) SetPath(path string) error {
	if err := v.ensureOpen(); err != nil {
		return err
	}
	return v.ctx.Check(util.ErrorCode(v.c.InitPathString(v.ctx.Internal(), v.state, path)))
}

// SetList initializes this value as a list containing elems. All elements
// must belong to the same evaluator state as the receiver.
func (v *Value) SetList(elems []*Value) error {
	if err := v.ensureOpen(); err != nil {
		return err
	}
	for i, elem := range elems {
		if elem == nil || elem.c == nil || elem.c.Ptr() == nil {
			return fmt.Errorf("nix: list element %d is nil or closed", i)
		}
		if elem.state == nil || v.state == nil || !elem.state.Same(v.state) {
			return fmt.Errorf("nix: list element %d belongs to a different evaluator state", i)
		}
	}
	builder := cvalue.MakeListBuilder(v.ctx.Internal(), v.state, uint(len(elems)))
	if builder == nil {
		if err := v.ctx.Err(); err != nil {
			return err
		}
		return fmt.Errorf("nix: failed to allocate list builder")
	}
	defer builder.Free()
	for i, elem := range elems {
		if err := v.ctx.Check(util.ErrorCode(builder.Insert(v.ctx.Internal(), uint(i), elem.c))); err != nil {
			return err
		}
	}
	return v.ctx.Check(util.ErrorCode(v.c.MakeList(v.ctx.Internal(), builder)))
}

// SetAttrs initializes this value as an attribute set. All values must belong
// to the same evaluator state as the receiver.
func (v *Value) SetAttrs(attrs map[string]*Value) error {
	if err := v.ensureOpen(); err != nil {
		return err
	}
	for name, elem := range attrs {
		if elem == nil || elem.c == nil || elem.c.Ptr() == nil {
			return fmt.Errorf("nix: attribute %q has a nil or closed value", name)
		}
		if elem.state != v.state {
			return fmt.Errorf("nix: attribute %q belongs to a different evaluator state", name)
		}
	}
	builder := cvalue.MakeBindingsBuilder(v.ctx.Internal(), v.state, uint(len(attrs)))
	if builder == nil {
		if err := v.ctx.Err(); err != nil {
			return err
		}
		return fmt.Errorf("nix: failed to allocate attribute bindings builder")
	}
	defer builder.Free()
	for name, elem := range attrs {
		if err := v.ctx.Check(util.ErrorCode(builder.Insert(v.ctx.Internal(), name, elem.c))); err != nil {
			return err
		}
	}
	return v.ctx.Check(util.ErrorCode(v.c.MakeAttrs(v.ctx.Internal(), builder)))
}

// NewList builds a list value from elements.
func NewList(ctx *util.Context, state *cexpr.EvalState, elems []*Value) (*Value, error) {
	if ctx == nil || state == nil {
		return nil, fmt.Errorf("nix: context and evaluator state are required")
	}
	for i, elem := range elems {
		if elem == nil || elem.c == nil || elem.c.Ptr() == nil {
			return nil, fmt.Errorf("nix: list element %d is nil or closed", i)
		}
		if elem.state == nil || !elem.state.Same(state) {
			return nil, fmt.Errorf("nix: list element %d belongs to a different evaluator state", i)
		}
	}
	v := cvalue.Alloc(ctx.Internal(), state)
	if v == nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to allocate list value")
	}
	b := cvalue.MakeListBuilder(ctx.Internal(), state, uint(len(elems)))
	if b == nil {
		v.Decref()
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to allocate list builder")
	}
	defer b.Free()
	for i, e := range elems {
		if err := ctx.Check(util.ErrorCode(b.Insert(ctx.Internal(), uint(i), e.c))); err != nil {
			v.Decref()
			return nil, err
		}
	}
	if err := ctx.Check(util.ErrorCode(v.MakeList(ctx.Internal(), b))); err != nil {
		v.Decref()
		return nil, err
	}
	return FromInternal(ctx, state, v), nil
}

// NewAttrs builds an attribute set from a map.
func NewAttrs(ctx *util.Context, state *cexpr.EvalState, attrs map[string]*Value) (*Value, error) {
	if ctx == nil || state == nil {
		return nil, fmt.Errorf("nix: context and evaluator state are required")
	}
	for name, elem := range attrs {
		if elem == nil || elem.c == nil {
			return nil, fmt.Errorf("nix: attribute %q has a nil or closed value", name)
		}
		if elem.state != state {
			return nil, fmt.Errorf("nix: attribute %q belongs to a different evaluator state", name)
		}
	}
	v := cvalue.Alloc(ctx.Internal(), state)
	if v == nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to allocate attribute-set value")
	}
	b := cvalue.MakeBindingsBuilder(ctx.Internal(), state, uint(len(attrs)))
	if b == nil {
		v.Decref()
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("nix: failed to allocate attribute bindings builder")
	}
	defer b.Free()
	for name, e := range attrs {
		if err := ctx.Check(util.ErrorCode(b.Insert(ctx.Internal(), name, e.c))); err != nil {
			v.Decref()
			return nil, err
		}
	}
	if err := ctx.Check(util.ErrorCode(v.MakeAttrs(ctx.Internal(), b))); err != nil {
		v.Decref()
		return nil, err
	}
	return FromInternal(ctx, state, v), nil
}
