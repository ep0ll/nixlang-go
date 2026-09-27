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
	c     *cvalue.Value
	state *cexpr.EvalState
	ctx   *util.Context
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
	if v != nil && v.c != nil {
		v.c.Decref()
		v.c = nil
	}
}

// Internal returns the low-level value.
func (v *Value) Internal() *cvalue.Value {
	if v == nil {
		return nil
	}
	return v.c
}

// Type returns the current type (call Force first for thunks).
func (v *Value) Type() Type {
	return Type(v.c.Type(v.ctx.Internal()))
}

// Force evaluates a thunk in place.
func (v *Value) Force() error {
	code := cexpr.ValueForce(v.ctx.Internal(), v.state, v.c.Ptr())
	return v.ctx.Check(util.ErrorCode(code))
}

// ForceDeep deeply forces the value. Avoid on recursive data.
func (v *Value) ForceDeep() error {
	code := cexpr.ValueForceDeep(v.ctx.Internal(), v.state, v.c.Ptr())
	return v.ctx.Check(util.ErrorCode(code))
}

// Bool returns the boolean value.
func (v *Value) Bool() (bool, error) {
	if err := v.ensureType(TypeBool); err != nil {
		return false, err
	}
	return v.c.GetBool(v.ctx.Internal()), nil
}

// Int returns the integer value.
func (v *Value) Int() (int64, error) {
	if err := v.ensureType(TypeInt); err != nil {
		return 0, err
	}
	return v.c.GetInt(v.ctx.Internal()), nil
}

// Float returns the float value.
func (v *Value) Float() (float64, error) {
	if err := v.ensureType(TypeFloat); err != nil {
		return 0, err
	}
	return v.c.GetFloat(v.ctx.Internal()), nil
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
	return v.c.GetPathString(v.ctx.Internal()), nil
}

// IsNull reports whether the value is null.
func (v *Value) IsNull() (bool, error) {
	if err := v.Force(); err != nil {
		return false, err
	}
	return v.Type() == TypeNull, nil
}

// ListLen returns the length of a list.
func (v *Value) ListLen() (int, error) {
	if err := v.ensureType(TypeList); err != nil {
		return 0, err
	}
	return int(v.c.ListSize(v.ctx.Internal())), nil
}

// ListAt returns the element at index i (forced).
func (v *Value) ListAt(i int) (*Value, error) {
	if err := v.ensureType(TypeList); err != nil {
		return nil, err
	}
	cv := v.c.ListByIdx(v.ctx.Internal(), v.state, uint(i))
	if cv == nil {
		return nil, v.ctx.Err()
	}
	return FromInternal(v.ctx, v.state, cv), nil
}

// ListAtLazy returns the element at index i without forcing it.
func (v *Value) ListAtLazy(i int) (*Value, error) {
	if err := v.ensureType(TypeList); err != nil {
		return nil, err
	}
	cv := v.c.ListByIdxLazy(v.ctx.Internal(), v.state, uint(i))
	if cv == nil {
		return nil, v.ctx.Err()
	}
	return FromInternal(v.ctx, v.state, cv), nil
}

// AttrsLen returns the number of attributes.
func (v *Value) AttrsLen() (int, error) {
	if err := v.ensureType(TypeAttrs); err != nil {
		return 0, err
	}
	return int(v.c.AttrsSize(v.ctx.Internal())), nil
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

// HasAttr reports whether the attribute exists.
func (v *Value) HasAttr(name string) (bool, error) {
	if err := v.ensureType(TypeAttrs); err != nil {
		return false, err
	}
	return v.c.HasAttr(v.ctx.Internal(), v.state, name), nil
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
			continue
		}
		m[name] = FromInternal(v.ctx, v.state, cv)
	}
	return m, nil
}

func (v *Value) ensureType(want Type) error {
	if err := v.Force(); err != nil {
		return err
	}
	got := v.Type()
	if got != want {
		return fmt.Errorf("nix: expected %s, got %s", want, got)
	}
	return nil
}

// SetNull initializes this value to null.
func (v *Value) SetNull() error {
	return v.ctx.Check(util.ErrorCode(v.c.InitNull(v.ctx.Internal())))
}

// SetBool initializes this value to a boolean.
func (v *Value) SetBool(b bool) error {
	return v.ctx.Check(util.ErrorCode(v.c.InitBool(v.ctx.Internal(), b)))
}

// SetInt initializes this value to an integer.
func (v *Value) SetInt(i int64) error {
	return v.ctx.Check(util.ErrorCode(v.c.InitInt(v.ctx.Internal(), i)))
}

// SetFloat initializes this value to a float.
func (v *Value) SetFloat(f float64) error {
	return v.ctx.Check(util.ErrorCode(v.c.InitFloat(v.ctx.Internal(), f)))
}

// SetString initializes this value to a string.
func (v *Value) SetString(s string) error {
	return v.ctx.Check(util.ErrorCode(v.c.InitString(v.ctx.Internal(), s)))
}

// SetPath initializes this value to a path.
func (v *Value) SetPath(path string) error {
	return v.ctx.Check(util.ErrorCode(v.c.InitPathString(v.ctx.Internal(), v.state, path)))
}

// NewList builds a list value from elements.
func NewList(ctx *util.Context, state *cexpr.EvalState, elems []*Value) (*Value, error) {
	v := cvalue.Alloc(ctx.Internal(), state)
	if v == nil {
		return nil, ctx.Err()
	}
	b := cvalue.MakeListBuilder(ctx.Internal(), state, uint(len(elems)))
	if b == nil {
		v.Decref()
		return nil, ctx.Err()
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
	v := cvalue.Alloc(ctx.Internal(), state)
	if v == nil {
		return nil, ctx.Err()
	}
	b := cvalue.MakeBindingsBuilder(ctx.Internal(), state, uint(len(attrs)))
	if b == nil {
		v.Decref()
		return nil, ctx.Err()
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
