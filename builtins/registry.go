// Code generated from the Nix 2.35.6 primop registrations; DO NOT EDIT by hand.
package builtins

// Spec describes one callable Nix evaluator builtin.
type Spec struct {
	// Name is the attribute name in the builtins set.
	Name string
	// Arity is the number of arguments accepted by the fully-applied builtin.
	Arity int
	// Args are the Nix source-level parameter names in call order.
	Args []string
	// GoTypes are the Go parameter types used by the generated wrapper.
	GoTypes []string
	// ExperimentalFeature is the Nix experimental feature gate, if any.
	ExperimentalFeature string
	// Optional is true for platform/configuration-dependent builtins.
	Optional bool
	// Condition documents the platform or evaluator setting controlling availability.
	Condition string
}

// Registry returns the known Nix 2.35.6 callable builtins. The returned slice
// and its nested metadata slices are independent from the package registry.
func Registry() []Spec {
	specs := make([]Spec, len(registry))
	copy(specs, registry)
	for i := range specs {
		specs[i].Args = append([]string(nil), registry[i].Args...)
		specs[i].GoTypes = append([]string(nil), registry[i].GoTypes...)
	}
	return specs
}

var registry = []Spec{
	{Name: "add", Arity: 2, Args: []string{"e1", "e2"}, GoTypes: []string{"any", "any"}},
	{Name: "addDrvOutputDependencies", Arity: 1, Args: []string{"s"}, GoTypes: []string{"string"}},
	{Name: "addErrorContext", Arity: 2, Args: []string{"context", "value"}, GoTypes: []string{"string", "any"}},
	{Name: "abort", Arity: 1, Args: []string{"s"}, GoTypes: []string{"string"}},
	{Name: "throw", Arity: 1, Args: []string{"s"}, GoTypes: []string{"string"}},
	{Name: "all", Arity: 2, Args: []string{"pred", "list"}, GoTypes: []string{"Function", "any"}},
	{Name: "any", Arity: 2, Args: []string{"pred", "list"}, GoTypes: []string{"Function", "any"}},
	{Name: "appendContext", Arity: 2, Args: []string{"string", "context"}, GoTypes: []string{"string", "StringContext"}},
	{Name: "attrNames", Arity: 1, Args: []string{"set"}, GoTypes: []string{"AttrSet"}},
	{Name: "attrValues", Arity: 1, Args: []string{"set"}, GoTypes: []string{"AttrSet"}},
	{Name: "baseNameOf", Arity: 1, Args: []string{"x"}, GoTypes: []string{"any"}},
	{Name: "bitAnd", Arity: 2, Args: []string{"e1", "e2"}, GoTypes: []string{"int64", "int64"}},
	{Name: "bitOr", Arity: 2, Args: []string{"e1", "e2"}, GoTypes: []string{"int64", "int64"}},
	{Name: "bitXor", Arity: 2, Args: []string{"e1", "e2"}, GoTypes: []string{"int64", "int64"}},
	{Name: "break", Arity: 1, Args: []string{"v"}, GoTypes: []string{"any"}},
	{Name: "catAttrs", Arity: 2, Args: []string{"attr", "list"}, GoTypes: []string{"string", "[]AttrSet"}},
	{Name: "ceil", Arity: 1, Args: []string{"number"}, GoTypes: []string{"any"}},
	{Name: "compareVersions", Arity: 2, Args: []string{"s1", "s2"}, GoTypes: []string{"string", "string"}},
	{Name: "concatLists", Arity: 1, Args: []string{"lists"}, GoTypes: []string{"any"}},
	{Name: "concatMap", Arity: 2, Args: []string{"f", "list"}, GoTypes: []string{"Function", "any"}},
	{Name: "concatStringsSep", Arity: 2, Args: []string{"separator", "list"}, GoTypes: []string{"string", "[]string"}},
	{Name: "convertHash", Arity: 1, Args: []string{"args"}, GoTypes: []string{"ConvertHashArgs"}},
	{Name: "deepSeq", Arity: 2, Args: []string{"e1", "e2"}, GoTypes: []string{"any", "any"}},
	{Name: "derivationStrict", Arity: 1, Args: []string{"attrs"}, GoTypes: []string{"DerivationAttrs"}},
	{Name: "dirOf", Arity: 1, Args: []string{"s"}, GoTypes: []string{"any"}},
	{Name: "div", Arity: 2, Args: []string{"e1", "e2"}, GoTypes: []string{"any", "any"}},
	{Name: "elem", Arity: 2, Args: []string{"x", "xs"}, GoTypes: []string{"any", "any"}},
	{Name: "elemAt", Arity: 2, Args: []string{"xs", "n"}, GoTypes: []string{"any", "int64"}},
	{Name: "exec", Arity: 1, Args: []string{"command"}, GoTypes: []string{"[]any"}, Optional: true, Condition: "native-code on non-Windows platforms"},
	{Name: "fetchClosure", Arity: 1, Args: []string{"args"}, GoTypes: []string{"FetchClosureArgs"}, ExperimentalFeature: "fetch-closure"},
	{Name: "fetchGit", Arity: 1, Args: []string{"args"}, GoTypes: []string{"any"}},
	{Name: "fetchMercurial", Arity: 1, Args: []string{"args"}, GoTypes: []string{"any"}},
	{Name: "fetchTarball", Arity: 1, Args: []string{"args"}, GoTypes: []string{"any"}},
	{Name: "fetchTree", Arity: 1, Args: []string{"input"}, GoTypes: []string{"any"}, ExperimentalFeature: "fetch-tree"},
	{Name: "filter", Arity: 2, Args: []string{"f", "list"}, GoTypes: []string{"Function", "any"}},
	{Name: "filterSource", Arity: 2, Args: []string{"e1", "e2"}, GoTypes: []string{"Function", "any"}},
	{Name: "findFile", Arity: 2, Args: []string{"search-path", "lookup-path"}, GoTypes: []string{"[]SearchPathEntry", "string"}},
	{Name: "floor", Arity: 1, Args: []string{"number"}, GoTypes: []string{"any"}},
	{Name: "foldl'", Arity: 3, Args: []string{"op", "nul", "list"}, GoTypes: []string{"Function", "any", "any"}},
	{Name: "fromJSON", Arity: 1, Args: []string{"e"}, GoTypes: []string{"string"}},
	{Name: "fromTOML", Arity: 1, Args: []string{"e"}, GoTypes: []string{"string"}},
	{Name: "functionArgs", Arity: 1, Args: []string{"f"}, GoTypes: []string{"Function"}},
	{Name: "genList", Arity: 2, Args: []string{"generator", "length"}, GoTypes: []string{"Function", "int64"}},
	{Name: "genericClosure", Arity: 1, Args: []string{"attrset"}, GoTypes: []string{"GenericClosureArgs"}},
	{Name: "getAttr", Arity: 2, Args: []string{"s", "set"}, GoTypes: []string{"string", "AttrSet"}},
	{Name: "getContext", Arity: 1, Args: []string{"s"}, GoTypes: []string{"string"}},
	{Name: "getEnv", Arity: 1, Args: []string{"s"}, GoTypes: []string{"string"}},
	{Name: "getFlake", Arity: 1, Args: []string{"args"}, GoTypes: []string{"any"}, ExperimentalFeature: "flakes"},
	{Name: "groupBy", Arity: 2, Args: []string{"f", "list"}, GoTypes: []string{"Function", "any"}},
	{Name: "hashFile", Arity: 2, Args: []string{"type", "p"}, GoTypes: []string{"string", "any"}},
	{Name: "hashString", Arity: 2, Args: []string{"type", "s"}, GoTypes: []string{"string", "string"}},
	{Name: "hasAttr", Arity: 2, Args: []string{"s", "set"}, GoTypes: []string{"string", "AttrSet"}},
	{Name: "hasContext", Arity: 1, Args: []string{"s"}, GoTypes: []string{"string"}},
	{Name: "head", Arity: 1, Args: []string{"list"}, GoTypes: []string{"any"}},
	{Name: "import", Arity: 1, Args: []string{"path"}, GoTypes: []string{"any"}},
	{Name: "importNative", Arity: 2, Args: []string{"path", "symbol"}, GoTypes: []string{"any", "string"}, Optional: true, Condition: "native-code on non-Windows platforms"},
	{Name: "intersectAttrs", Arity: 2, Args: []string{"e1", "e2"}, GoTypes: []string{"AttrSet", "AttrSet"}},
	{Name: "isAttrs", Arity: 1, Args: []string{"e"}, GoTypes: []string{"any"}},
	{Name: "isBool", Arity: 1, Args: []string{"e"}, GoTypes: []string{"any"}},
	{Name: "isFloat", Arity: 1, Args: []string{"e"}, GoTypes: []string{"any"}},
	{Name: "isFunction", Arity: 1, Args: []string{"e"}, GoTypes: []string{"any"}},
	{Name: "isInt", Arity: 1, Args: []string{"e"}, GoTypes: []string{"any"}},
	{Name: "isList", Arity: 1, Args: []string{"e"}, GoTypes: []string{"any"}},
	{Name: "isNull", Arity: 1, Args: []string{"e"}, GoTypes: []string{"any"}},
	{Name: "isPath", Arity: 1, Args: []string{"e"}, GoTypes: []string{"any"}},
	{Name: "isString", Arity: 1, Args: []string{"e"}, GoTypes: []string{"any"}},
	{Name: "length", Arity: 1, Args: []string{"e"}, GoTypes: []string{"any"}},
	{Name: "lessThan", Arity: 2, Args: []string{"e1", "e2"}, GoTypes: []string{"any", "any"}},
	{Name: "listToAttrs", Arity: 1, Args: []string{"e"}, GoTypes: []string{"[]NameValue"}},
	{Name: "map", Arity: 2, Args: []string{"f", "list"}, GoTypes: []string{"Function", "any"}},
	{Name: "mapAttrs", Arity: 2, Args: []string{"f", "attrset"}, GoTypes: []string{"Function", "AttrSet"}},
	{Name: "match", Arity: 2, Args: []string{"regex", "str"}, GoTypes: []string{"string", "string"}},
	{Name: "mul", Arity: 2, Args: []string{"e1", "e2"}, GoTypes: []string{"any", "any"}},
	{Name: "outputOf", Arity: 2, Args: []string{"derivation-reference", "output-name"}, GoTypes: []string{"string", "string"}, ExperimentalFeature: "dynamic-derivations"},
	{Name: "parseDrvName", Arity: 1, Args: []string{"s"}, GoTypes: []string{"string"}},
	{Name: "parseFlakeRef", Arity: 1, Args: []string{"flake-ref"}, GoTypes: []string{"string"}, ExperimentalFeature: "flakes"},
	{Name: "partition", Arity: 2, Args: []string{"pred", "list"}, GoTypes: []string{"Function", "any"}},
	{Name: "path", Arity: 1, Args: []string{"args"}, GoTypes: []string{"PathArgs"}},
	{Name: "pathExists", Arity: 1, Args: []string{"path"}, GoTypes: []string{"any"}},
	{Name: "placeholder", Arity: 1, Args: []string{"output"}, GoTypes: []string{"string"}},
	{Name: "readDir", Arity: 1, Args: []string{"path"}, GoTypes: []string{"any"}},
	{Name: "readFile", Arity: 1, Args: []string{"path"}, GoTypes: []string{"any"}},
	{Name: "readFileType", Arity: 1, Args: []string{"p"}, GoTypes: []string{"any"}},
	{Name: "removeAttrs", Arity: 2, Args: []string{"set", "list"}, GoTypes: []string{"AttrSet", "[]string"}},
	{Name: "replaceStrings", Arity: 3, Args: []string{"from", "to", "s"}, GoTypes: []string{"[]string", "[]string", "string"}},
	{Name: "scopedImport", Arity: 2, Args: []string{"scope", "path"}, GoTypes: []string{"AttrSet", "any"}},
	{Name: "seq", Arity: 2, Args: []string{"e1", "e2"}, GoTypes: []string{"any", "any"}},
	{Name: "sort", Arity: 2, Args: []string{"comparator", "list"}, GoTypes: []string{"Function", "any"}},
	{Name: "split", Arity: 2, Args: []string{"regex", "str"}, GoTypes: []string{"string", "string"}},
	{Name: "splitVersion", Arity: 1, Args: []string{"s"}, GoTypes: []string{"string"}},
	{Name: "storePath", Arity: 1, Args: []string{"path"}, GoTypes: []string{"any"}},
	{Name: "stringLength", Arity: 1, Args: []string{"e"}, GoTypes: []string{"string"}},
	{Name: "sub", Arity: 2, Args: []string{"e1", "e2"}, GoTypes: []string{"any", "any"}},
	{Name: "substring", Arity: 3, Args: []string{"start", "len", "s"}, GoTypes: []string{"int64", "int64", "string"}},
	{Name: "tail", Arity: 1, Args: []string{"list"}, GoTypes: []string{"any"}},
	{Name: "toFile", Arity: 2, Args: []string{"name", "s"}, GoTypes: []string{"string", "string"}},
	{Name: "toJSON", Arity: 1, Args: []string{"e"}, GoTypes: []string{"any"}},
	{Name: "toPath", Arity: 1, Args: []string{"s"}, GoTypes: []string{"any"}},
	{Name: "toString", Arity: 1, Args: []string{"e"}, GoTypes: []string{"any"}},
	{Name: "toXML", Arity: 1, Args: []string{"e"}, GoTypes: []string{"any"}},
	{Name: "trace", Arity: 2, Args: []string{"e1", "e2"}, GoTypes: []string{"any", "any"}},
	{Name: "traceVerbose", Arity: 2, Args: []string{"e1", "e2"}, GoTypes: []string{"any", "any"}},
	{Name: "tryEval", Arity: 1, Args: []string{"e"}, GoTypes: []string{"any"}},
	{Name: "typeOf", Arity: 1, Args: []string{"e"}, GoTypes: []string{"any"}},
	{Name: "unsafeDiscardOutputDependency", Arity: 1, Args: []string{"s"}, GoTypes: []string{"string"}},
	{Name: "unsafeDiscardStringContext", Arity: 1, Args: []string{"s"}, GoTypes: []string{"any"}},
	{Name: "unsafeGetAttrPos", Arity: 2, Args: []string{"s", "set"}, GoTypes: []string{"string", "AttrSet"}},
	{Name: "warn", Arity: 2, Args: []string{"e1", "e2"}, GoTypes: []string{"any", "any"}},
	{Name: "zipAttrsWith", Arity: 2, Args: []string{"f", "list"}, GoTypes: []string{"Function", "[]AttrSet"}},
	{Name: "fetchurl", Arity: 1, Args: []string{"arg"}, GoTypes: []string{"any"}},
	{Name: "flakeRefToString", Arity: 1, Args: []string{"attrs"}, GoTypes: []string{"FlakeRefAttrs"}, ExperimentalFeature: "flakes"},
	{Name: "derivation", Arity: 1, Args: []string{"attrs"}, GoTypes: []string{"DerivationAttrs"}},
}
