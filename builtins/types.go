// Package builtins defines statically typed Go argument types for Nix builtins.
package builtins

import (
	"crypto"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/ep0ll/nixlang-go/value"
)

// PathLike is the closed union of a Nix path value and a Nix string containing
// a path. Use Path for the former and StringPath for the latter.
type PathLike interface {
	nixPathLike() any
}

// StringPath represents a path passed to Nix as a string.
type StringPath string

func (p StringPath) nixPathLike() any { return string(p) }

// FilePath normalizes a host filesystem path and marks it as a Nix path value.
// filepath is used only for host-path normalization; the Nix value remains
// distinct from a Nix string.
func FilePath(path string) Path {
	return Path(filepath.ToSlash(filepath.Clean(path)))
}

// StorePath is a canonical Nix store path, not an arbitrary filesystem path.
type StorePath string

// HashAlgorithm reuses Go's standard crypto.Hash identifiers. Nix supports
// only these four algorithms for hashFile, hashString, and convertHash.
type HashAlgorithm = crypto.Hash

const (
	HashMD5    = crypto.MD5
	HashSHA1   = crypto.SHA1
	HashSHA256 = crypto.SHA256
	HashSHA512 = crypto.SHA512
)

func hashAlgorithmName(algorithm crypto.Hash) (string, error) {
	switch algorithm {
	case crypto.MD5:
		return "md5", nil
	case crypto.SHA1:
		return "sha1", nil
	case crypto.SHA256:
		return "sha256", nil
	case crypto.SHA512:
		return "sha512", nil
	default:
		return "", fmt.Errorf("unsupported Nix hash algorithm %v", algorithm)
	}
}

// HashFormat identifies the output representation accepted by
// builtins.convertHash. The zero value is deliberately invalid.
type HashFormat uint8

const (
	HashFormatInvalid HashFormat = iota
	HashFormatBase16
	HashFormatBase32
	HashFormatBase64
	HashFormatSRI
)

func (f HashFormat) String() string {
	switch f {
	case HashFormatBase16:
		return "base16"
	case HashFormatBase32:
		return "base32"
	case HashFormatBase64:
		return "base64"
	case HashFormatSRI:
		return "sri"
	default:
		return ""
	}
}

// Number is the integer-or-float union accepted by Nix arithmetic primops.
type Number struct {
	isFloat bool
	integer int64
	float   float64
}

// Integer constructs an integer Nix number.
func Integer(n int64) Number { return Number{integer: n} }

// Floating constructs a floating-point Nix number.
func Floating(n float64) Number { return Number{isFloat: true, float: n} }

// ComparisonValue is the number/string/path union accepted by builtins.lessThan.
// Nix requires both operands to have the same comparable type.
type ComparisonValue struct {
	kind    comparisonKind
	number  Number
	text    string
	path    Path
}

type comparisonKind uint8

const (
	comparisonInvalid comparisonKind = iota
	comparisonNumber
	comparisonString
	comparisonPath
)

func CompareNumbers(n Number) ComparisonValue {
	return ComparisonValue{kind: comparisonNumber, number: n}
}

func CompareStrings(s string) ComparisonValue {
	return ComparisonValue{kind: comparisonString, text: s}
}

func ComparePaths(p Path) ComparisonValue {
	return ComparisonValue{kind: comparisonPath, path: p}
}

// StringCoercible is the set of Go-side values that Nix can coerce to strings.
// Construct values with StringValue, IntegerString, PathStringValue,
// AttrStringValue, or ValueStringValue.
type StringCoercible struct {
	value any
}

func StringValue(s string) StringCoercible { return StringCoercible{value: s} }
func IntegerString(n int64) StringCoercible { return StringCoercible{value: n} }
func PathStringValue(p Path) StringCoercible { return StringCoercible{value: p} }
func AttrStringValue(attrs AttrSet) StringCoercible { return StringCoercible{value: attrs} }
func ValueStringValue(v *value.Value) StringCoercible { return StringCoercible{value: v} }

// Function is an independently retained, checked Nix function value from the
// same evaluator as Builtins. Keep it alive until lazy results that reference
// it have been consumed.
type Function struct {
	value *value.Value
}

// NewFunction validates and independently retains a Nix function value.
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
	retained, err := v.Retain()
	if err != nil {
		return Function{}, err
	}
	return Function{value: retained}, nil
}

// Close releases the retained function value. It is safe to call more than once.
func (f *Function) Close() {
	if f == nil || f.value == nil {
		return
	}
	f.value.Close()
	f.value = nil
}

// Value returns a borrowed pointer to the retained evaluator value.
func (f Function) Value() *value.Value { return f.value }

// AttrSet is a Nix attribute set whose values are recursively converted from Go.
type AttrSet = map[string]any

// AttrConvertible describes a typed Go struct that converts to a Nix attrset.
type AttrConvertible interface {
	ToNixAttrs() AttrSet
}

// ArgumentValidator lets structured builtin arguments reject invalid
// combinations before entering the evaluator.
type ArgumentValidator interface {
	Validate() error
}

// NameValue is one element accepted by builtins.listToAttrs.
type NameValue struct {
	Name  string
	Value any
}

func (a NameValue) ToNixAttrs() AttrSet {
	return AttrSet{"name": a.Name, "value": a.Value}
}

func (a NameValue) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("listToAttrs entry name must not be empty")
	}
	return nil
}

// SearchPathEntry describes one entry supplied to builtins.findFile.
type SearchPathEntry struct {
	Prefix string
	Path   PathLike
}

func (e SearchPathEntry) ToNixAttrs() AttrSet {
	attrs := AttrSet{"path": e.Path}
	if e.Prefix != "" {
		attrs["prefix"] = e.Prefix
	}
	return attrs
}

func (e SearchPathEntry) Validate() error {
	if e.Path == nil {
		return fmt.Errorf("findFile search path entry requires a path")
	}
	return nil
}

// PathArgs describes the structured argument to builtins.path.
type PathArgs struct {
	Path      PathLike
	Name      string
	Filter    *Function
	Recursive *bool
	SHA256    string
	Extra     AttrSet
}

func (a PathArgs) ToNixAttrs() AttrSet {
	attrs := cloneAttrs(a.Extra)
	if a.Path != nil {
		attrs["path"] = a.Path
	}
	if a.Name != "" {
		attrs["name"] = a.Name
	}
	if a.Filter != nil {
		attrs["filter"] = *a.Filter
	}
	if a.Recursive != nil {
		attrs["recursive"] = *a.Recursive
	}
	if a.SHA256 != "" {
		attrs["sha256"] = a.SHA256
	}
	return attrs
}

func (a PathArgs) Validate() error {
	if a.Path == nil {
		return fmt.Errorf("path builtin requires Path")
	}
	return nil
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

func (c StringContext) ToNixAttrs() AttrSet {
	attrs := make(AttrSet, len(c))
	for storePath, entry := range c {
		attrs[storePath] = entry.ToNixAttrs()
	}
	return attrs
}

// ConvertHashArgs describes builtins.convertHash.
type ConvertHashArgs struct {
	Hash         string
	HashAlgo     HashAlgorithm // Optional when the hash encodes its algorithm.
	ToHashFormat HashFormat
	Extra        AttrSet
}

func (a ConvertHashArgs) ToNixAttrs() AttrSet {
	attrs := cloneAttrs(a.Extra)
	attrs["hash"] = a.Hash
	if a.HashAlgo != 0 {
		if algorithm, err := hashAlgorithmName(a.HashAlgo); err == nil {
			attrs["hashAlgo"] = algorithm
		}
	}
	attrs["toHashFormat"] = a.ToHashFormat.String()
	return attrs
}

func (a ConvertHashArgs) Validate() error {
	if a.Hash == "" {
		return fmt.Errorf("convertHash requires Hash")
	}
	if a.ToHashFormat.String() == "" {
		return fmt.Errorf("convertHash requires a valid ToHashFormat")
	}
	if a.HashAlgo != 0 {
		if _, err := hashAlgorithmName(a.HashAlgo); err != nil {
			return err
		}
	}
	return nil
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

func (a GenericClosureArgs) Validate() error {
	if a.Operator.value == nil {
		return fmt.Errorf("genericClosure requires Operator")
	}
	return nil
}

// FetchClosureArgs describes the structured form of builtins.fetchClosure.
type FetchClosureArgs struct {
	FromPath       StorePath
	FromStore      url.URL
	ToPath         *StorePath
	InputAddressed *bool
}

func (a FetchClosureArgs) ToNixAttrs() AttrSet {
	attrs := AttrSet{"fromPath": string(a.FromPath)}
	if a.FromStore.String() != "" {
		attrs["fromStore"] = a.FromStore.String()
	}
	if a.ToPath != nil {
		attrs["toPath"] = string(*a.ToPath)
	}
	if a.InputAddressed != nil {
		attrs["inputAddressed"] = *a.InputAddressed
	}
	return attrs
}

func (a FetchClosureArgs) Validate() error {
	if a.FromPath == "" {
		return fmt.Errorf("fetchClosure requires FromPath")
	}
	if a.FromStore.Scheme == "" || a.FromStore.Host == "" {
		return fmt.Errorf("fetchClosure requires an absolute FromStore URL")
	}
	if a.InputAddressed != nil && *a.InputAddressed && a.ToPath != nil {
		return fmt.Errorf("fetchClosure cannot combine InputAddressed=true with ToPath")
	}
	return nil
}

// FetchURLArgs describes the attribute-set form of builtins.fetchurl.
type FetchURLArgs struct {
	URL    url.URL
	Name   string
	SHA256 string
	Extra  AttrSet
}

func (a FetchURLArgs) ToNixAttrs() AttrSet {
	attrs := cloneAttrs(a.Extra)
	attrs["url"] = a.URL.String()
	if a.Name != "" {
		attrs["name"] = a.Name
	}
	if a.SHA256 != "" {
		attrs["sha256"] = a.SHA256
	}
	return attrs
}

func (a FetchURLArgs) Validate() error {
	if a.URL.Scheme == "" {
		return fmt.Errorf("fetchurl requires an absolute URL")
	}
	return nil
}

// FetchTarballArgs describes the attribute-set form of builtins.fetchTarball.
type FetchTarballArgs struct {
	URL    url.URL
	SHA256 string
	Name   string
	Extra  AttrSet
}

func (a FetchTarballArgs) ToNixAttrs() AttrSet {
	attrs := cloneAttrs(a.Extra)
	attrs["url"] = a.URL.String()
	if a.SHA256 != "" {
		attrs["sha256"] = a.SHA256
	}
	if a.Name != "" {
		attrs["name"] = a.Name
	}
	return attrs
}

func (a FetchTarballArgs) Validate() error {
	if a.URL.Scheme == "" {
		return fmt.Errorf("fetchTarball requires an absolute URL")
	}
	return nil
}

// FetchGitArgs describes the attribute-set form of builtins.fetchGit.
type FetchGitArgs struct {
	URL          *url.URL
	Path         *Path
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
	if a.URL != nil {
		attrs["url"] = a.URL.String()
	}
	if a.Path != nil {
		attrs["url"] = *a.Path
	}
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

func (a FetchGitArgs) Validate() error {
	if (a.URL == nil) == (a.Path == nil) {
		return fmt.Errorf("fetchGit requires exactly one of URL or Path")
	}
	if a.URL != nil && a.URL.Scheme == "" {
		return fmt.Errorf("fetchGit URL must be absolute")
	}
	return nil
}

// FetchMercurialArgs describes the attribute-set form of builtins.fetchMercurial.
type FetchMercurialArgs struct {
	URL   *url.URL
	Path  *Path
	Name  string
	Rev   string
	Extra AttrSet
}

func (a FetchMercurialArgs) ToNixAttrs() AttrSet {
	attrs := cloneAttrs(a.Extra)
	if a.URL != nil {
		attrs["url"] = a.URL.String()
	}
	if a.Path != nil {
		attrs["url"] = *a.Path
	}
	if a.Name != "" {
		attrs["name"] = a.Name
	}
	if a.Rev != "" {
		attrs["rev"] = a.Rev
	}
	return attrs
}

func (a FetchMercurialArgs) Validate() error {
	if (a.URL == nil) == (a.Path == nil) {
		return fmt.Errorf("fetchMercurial requires exactly one of URL or Path")
	}
	if a.URL != nil && a.URL.Scheme == "" {
		return fmt.Errorf("fetchMercurial URL must be absolute")
	}
	return nil
}

// FetchTreeArgs describes the structured argument to builtins.fetchTree.
// URL and Path are represented separately to preserve their Nix value kinds.
// Extra carries backend-specific fetcher attributes.
type FetchTreeArgs struct {
	Type         string
	URL          *url.URL
	Path         *Path
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
		attrs["url"] = a.URL.String()
	}
	if a.Path != nil {
		attrs["path"] = *a.Path
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

func (a FetchTreeArgs) Validate() error {
	if a.Type == "" {
		return fmt.Errorf("fetchTree requires Type")
	}
	if a.URL != nil && a.Path != nil {
		return fmt.Errorf("fetchTree cannot combine URL and Path")
	}
	if a.URL != nil && a.URL.Scheme == "" {
		return fmt.Errorf("fetchTree URL must be absolute")
	}
	return nil
}

// FlakeReference is the URL-or-path union accepted by builtins.getFlake.
type FlakeReference struct {
	URL  *url.URL
	Path *Path
}

func (r FlakeReference) Validate() error {
	if (r.URL == nil) == (r.Path == nil) {
		return fmt.Errorf("getFlake requires exactly one of URL or Path")
	}
	if r.URL != nil && r.URL.Scheme == "" {
		return fmt.Errorf("getFlake URL must be absolute")
	}
	return nil
}

func (r FlakeReference) ToNixArgument() any {
	if r.Path != nil {
		return *r.Path
	}
	if r.URL != nil {
		return r.URL.String()
	}
	return nil
}

// FlakeRefType is the backend kind accepted by a flake reference.
type FlakeRefType string

const (
	FlakeRefGitHub   FlakeRefType = "github"
	FlakeRefGitLab   FlakeRefType = "gitlab"
	FlakeRefSourceHut FlakeRefType = "sourcehut"
	FlakeRefGit      FlakeRefType = "git"
	FlakeRefPath     FlakeRefType = "path"
	FlakeRefIndirect FlakeRefType = "indirect"
)

// FlakeRefAttrs describes the structured argument to builtins.flakeRefToString.
type FlakeRefAttrs struct {
	Type       FlakeRefType
	Owner      string
	Repo       string
	Ref        string
	Rev        string
	Dir        string
	NarHash    string
	Host       string
	Shallow    *bool
	Submodules *bool
	AllRefs    *bool
	Extra      AttrSet
}

func (a FlakeRefAttrs) ToNixAttrs() AttrSet {
	attrs := cloneAttrs(a.Extra)
	for key, val := range map[string]string{
		"type": string(a.Type), "owner": a.Owner, "repo": a.Repo, "ref": a.Ref,
		"rev": a.Rev, "dir": a.Dir, "narHash": a.NarHash, "host": a.Host,
	} {
		if val != "" {
			attrs[key] = val
		}
	}
	if a.Shallow != nil {
		attrs["shallow"] = *a.Shallow
	}
	if a.Submodules != nil {
		attrs["submodules"] = *a.Submodules
	}
	if a.AllRefs != nil {
		attrs["allRefs"] = *a.AllRefs
	}
	return attrs
}

// DerivationAttrs describes common derivation inputs. Extra carries attributes
// such as environment variables and Nix derivation options.
type DerivationAttrs struct {
	Name    string
	System  string
	Builder PathLike
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

func (a DerivationAttrs) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("derivation requires Name")
	}
	if a.System == "" {
		return fmt.Errorf("derivation requires System")
	}
	if a.Builder == nil {
		return fmt.Errorf("derivation requires Builder")
	}
	return nil
}

func cloneAttrs(in AttrSet) AttrSet {
	out := make(AttrSet, len(in))
	for key, val := range in {
		out[key] = val
	}
	return out
}
