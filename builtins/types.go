// Package builtins defines the Go-side argument types used by the generated
// Nix builtin wrappers.
package builtins

import (
	"fmt"

	"github.com/ep0ll/nixlang-go/value"
)

// Function is a checked Nix function value from the same evaluator as Builtins.
// Wrap an evaluator value with NewFunction before passing it to higher-order
// builtins. The underlying *value.Value remains owned by the caller and must
// stay open for as long as the Function is used.
type Function struct {
	value *value.Value
}

// NewFunction validates and wraps a Nix function value.
func NewFunction(v *value.Value) (Function, error) {
	if v == nil || v.Internal() == nil || v.Internal().Ptr() == nil {
		return Function{}, fmt.Errorf("nix: function value is nil or closed")
	}
	if err := v.Force(); err != nil {
		return Function{}, err
	}
	typ, err := v.TypeChecked()
	if err != nil {
		return Function{}, err
	}
	if typ != value.TypeFunction {
		return Function{}, fmt.Errorf("nix: expected function value, got %s", typ)
	}
	return Function{value: v}, nil
}

// Value returns the underlying evaluator value. The caller retains ownership.
func (f Function) Value() *value.Value {
	return f.value
}

// AttrSet is a Nix attribute set whose values are recursively converted from
// Go values. Use it for Nix APIs whose attribute names are open-ended.
type AttrSet = map[string]any

// AttrConvertible is implemented by the typed option structs in this package.
// The returned AttrSet is converted recursively using the same Go-to-Nix rules
// as Builtins.Call.
type AttrConvertible interface {
	ToNixAttrs() AttrSet
}


// NameValue is one element accepted by builtins.listToAttrs.
type NameValue struct {
	Name  string
	Value any
}

func (a NameValue) ToNixAttrs() AttrSet {
	return AttrSet{"name": a.Name, "value": a.Value}
}

// SearchPathEntry describes an entry supplied to builtins.findFile.
type SearchPathEntry struct {
	Prefix string
	Path   any // Path or string; Nix accepts both.
}

func (e SearchPathEntry) ToNixAttrs() AttrSet {
	attrs := AttrSet{"path": e.Path}
	if e.Prefix != "" {
		attrs["prefix"] = e.Prefix
	}
	return attrs
}

// PathArgs describes the structured form of builtins.path.
type PathArgs struct {
	Path      any      // Path or string.
	Name      string   // Optional store name.
	Filter    Function // Optional Nix predicate (path: type: bool).
	Recursive *bool    // Optional; nil uses Nix's default.
	SHA256    string   // Optional expected SHA-256 hash.
	Extra     AttrSet  // Reserved for forward-compatible extension fields.
}

func (a PathArgs) ToNixAttrs() AttrSet {
	attrs := cloneAttrs(a.Extra)
	attrs["path"] = a.Path
	if a.Name != "" {
		attrs["name"] = a.Name
	}
	if a.Filter != nil {
		attrs["filter"] = a.Filter
	}
	if a.Recursive != nil {
		attrs["recursive"] = *a.Recursive
	}
	if a.SHA256 != "" {
		attrs["sha256"] = a.SHA256
	}
	return attrs
}

// StringContextEntry describes the context flags for one store path in
// builtins.appendContext.
type StringContextEntry struct {
	Path       bool
	AllOutputs bool
	Outputs    []string
}

func (e StringContextEntry) ToNixAttrs() AttrSet {
	attrs := AttrSet{}
	if e.Path {
		attrs["path"] = true
	}
	if e.AllOutputs {
		attrs["allOutputs"] = true
	}
	if len(e.Outputs) != 0 {
		attrs["outputs"] = e.Outputs
	}
	return attrs
}

// StringContext maps store paths to their Nix string-context metadata.
type StringContext map[string]StringContextEntry

// ConvertHashArgs describes builtins.convertHash.
type ConvertHashArgs struct {
	Hash         string
	HashAlgo     string // Optional when the hash encodes its algorithm.
	ToHashFormat string
	Extra        AttrSet
}

func (a ConvertHashArgs) ToNixAttrs() AttrSet {
	attrs := cloneAttrs(a.Extra)
	attrs["hash"] = a.Hash
	attrs["toHashFormat"] = a.ToHashFormat
	if a.HashAlgo != "" {
		attrs["hashAlgo"] = a.HashAlgo
	}
	return attrs
}

// GenericClosureArgs describes the input to builtins.genericClosure.
type GenericClosureArgs struct {
	StartSet []AttrSet
	Operator Function
	Extra    AttrSet
}

func (a GenericClosureArgs) ToNixAttrs() AttrSet {
	attrs := cloneAttrs(a.Extra)
	attrs["startSet"] = a.StartSet
	attrs["operator"] = a.Operator
	return attrs
}

// FetchClosureArgs describes the structured form of builtins.fetchClosure.
type FetchClosureArgs struct {
	FromPath       Path
	FromStore      string
	ToPath         any
	InputAddressed *bool
}

func (a FetchClosureArgs) ToNixAttrs() AttrSet {
	attrs := AttrSet{"fromPath": a.FromPath}
	if a.FromStore != "" {
		attrs["fromStore"] = a.FromStore
	}
	if a.ToPath != nil {
		attrs["toPath"] = a.ToPath
	}
	if a.InputAddressed != nil {
		attrs["inputAddressed"] = *a.InputAddressed
	}
	return attrs
}

// FetchURLArgs describes the attribute-set form of builtins.fetchurl.
type FetchURLArgs struct {
	URL   string
	Name  string
	Extra AttrSet
}

func (a FetchURLArgs) ToNixAttrs() AttrSet {
	attrs := cloneAttrs(a.Extra)
	attrs["url"] = a.URL
	if a.Name != "" {
		attrs["name"] = a.Name
	}
	return attrs
}

// FetchTarballArgs describes the attribute-set form of builtins.fetchTarball.
type FetchTarballArgs struct {
	URL    string
	SHA256 string
	Name   string
	Extra  AttrSet
}

func (a FetchTarballArgs) ToNixAttrs() AttrSet {
	attrs := cloneAttrs(a.Extra)
	attrs["url"] = a.URL
	if a.SHA256 != "" {
		attrs["sha256"] = a.SHA256
	}
	if a.Name != "" {
		attrs["name"] = a.Name
	}
	return attrs
}

// FetchGitArgs describes the attribute-set form of builtins.fetchGit.
type FetchGitArgs struct {
	URL          string
	Name         string
	Rev          string
	Ref          string
	Submodules   *bool
	Shallow      *bool
	AllRefs      *bool
	ExportIgnore *bool
	Extra        AttrSet
}

func (a FetchGitArgs) ToNixAttrs() AttrSet {
	attrs := cloneAttrs(a.Extra)
	attrs["url"] = a.URL
	if a.Name != "" {
		attrs["name"] = a.Name
	}
	if a.Rev != "" {
		attrs["rev"] = a.Rev
	}
	if a.Ref != "" {
		attrs["ref"] = a.Ref
	}
	if a.Submodules != nil {
		attrs["submodules"] = *a.Submodules
	}
	if a.Shallow != nil {
		attrs["shallow"] = *a.Shallow
	}
	if a.AllRefs != nil {
		attrs["allRefs"] = *a.AllRefs
	}
	if a.ExportIgnore != nil {
		attrs["exportIgnore"] = *a.ExportIgnore
	}
	return attrs
}

// FetchMercurialArgs describes the attribute-set form of builtins.fetchMercurial.
type FetchMercurialArgs struct {
	URL   any // Path or string.
	Name  string
	Rev   string
	Extra AttrSet
}

func (a FetchMercurialArgs) ToNixAttrs() AttrSet {
	attrs := cloneAttrs(a.Extra)
	attrs["url"] = a.URL
	if a.Name != "" {
		attrs["name"] = a.Name
	}
	if a.Rev != "" {
		attrs["rev"] = a.Rev
	}
	return attrs
}

// FetchTreeArgs describes the attribute-set form of builtins.fetchTree.
// URL and Path accept both builtins.Path and string because Nix fetcher input
// attributes support path/string coercion. Extra carries backend-specific
// input attributes not represented by the common fields.
type FetchTreeArgs struct {
	Type         string
	URL          any
	Path         any
	Rev          string
	Ref          string
	NarHash      string
	Name         string
	Submodules   *bool
	Shallow      *bool
	AllRefs      *bool
	ExportIgnore *bool
	Extra        AttrSet
}

func (a FetchTreeArgs) ToNixAttrs() AttrSet {
	attrs := cloneAttrs(a.Extra)
	if a.Type != "" {
		attrs["type"] = a.Type
	}
	if a.URL != nil {
		attrs["url"] = a.URL
	}
	if a.Path != nil {
		attrs["path"] = a.Path
	}
	if a.Rev != "" {
		attrs["rev"] = a.Rev
	}
	if a.Ref != "" {
		attrs["ref"] = a.Ref
	}
	if a.NarHash != "" {
		attrs["narHash"] = a.NarHash
	}
	if a.Name != "" {
		attrs["name"] = a.Name
	}
	if a.Submodules != nil {
		attrs["submodules"] = *a.Submodules
	}
	if a.Shallow != nil {
		attrs["shallow"] = *a.Shallow
	}
	if a.AllRefs != nil {
		attrs["allRefs"] = *a.AllRefs
	}
	if a.ExportIgnore != nil {
		attrs["exportIgnore"] = *a.ExportIgnore
	}
	return attrs
}

// FlakeRefAttrs describes the attribute-set form accepted by
// builtins.flakeRefToString. Extra can hold additional typed flake-ref fields.
type FlakeRefAttrs struct {
	Type    string
	Owner   string
	Repo    string
	Ref     string
	Rev     string
	Dir     string
	NarHash string
	Host    string
	Extra   AttrSet
}

func (a FlakeRefAttrs) ToNixAttrs() AttrSet {
	attrs := cloneAttrs(a.Extra)
	for key, val := range map[string]string{
		"type": a.Type, "owner": a.Owner, "repo": a.Repo, "ref": a.Ref,
		"rev": a.Rev, "dir": a.Dir, "narHash": a.NarHash, "host": a.Host,
	} {
		if val != "" {
			attrs[key] = val
		}
	}
	return attrs
}

// DerivationAttrs describes common derivation inputs. Extra carries additional
// derivation attributes (including arbitrary environment variables) without
// weakening the types of the common fields.
type DerivationAttrs struct {
	Name    string
	System  string
	Builder any // Path or string.
	Args    []string
	Outputs []string
	Env     AttrSet
	Extra   AttrSet
}

func (a DerivationAttrs) ToNixAttrs() AttrSet {
	attrs := cloneAttrs(a.Extra)
	if a.Name != "" {
		attrs["name"] = a.Name
	}
	if a.System != "" {
		attrs["system"] = a.System
	}
	if a.Builder != nil {
		attrs["builder"] = a.Builder
	}
	if a.Args != nil {
		attrs["args"] = a.Args
	}
	if a.Outputs != nil {
		attrs["outputs"] = a.Outputs
	}
	for key, val := range a.Env {
		attrs[key] = val
	}
	return attrs
}

func cloneAttrs(in AttrSet) AttrSet {
	out := make(AttrSet, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}
