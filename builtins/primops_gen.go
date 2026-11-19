// Code generated from the Nix 2.35.6 builtin registry; DO NOT EDIT by hand.
package builtins

import "github.com/ep0ll/nixlang-go/value"

// Add calls builtins.add with 2 Go arguments.
func (b *Builtins) Add(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("add", arg1, arg2)
}

// AddDrvOutputDependencies calls builtins.addDrvOutputDependencies with 1 Go argument.
func (b *Builtins) AddDrvOutputDependencies(arg1 any) (*value.Value, error) {
    return b.Call("addDrvOutputDependencies", arg1)
}

// AddErrorContext calls builtins.addErrorContext with 2 Go arguments.
func (b *Builtins) AddErrorContext(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("addErrorContext", arg1, arg2)
}

// Abort calls builtins.abort with 1 Go argument.
func (b *Builtins) Abort(arg1 any) (*value.Value, error) {
    return b.Call("abort", arg1)
}

// Throw calls builtins.throw with 1 Go argument.
func (b *Builtins) Throw(arg1 any) (*value.Value, error) {
    return b.Call("throw", arg1)
}

// All calls builtins.all with 2 Go arguments.
func (b *Builtins) All(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("all", arg1, arg2)
}

// Any calls builtins.any with 2 Go arguments.
func (b *Builtins) Any(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("any", arg1, arg2)
}

// AppendContext calls builtins.appendContext with 2 Go arguments.
func (b *Builtins) AppendContext(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("appendContext", arg1, arg2)
}

// AttrNames calls builtins.attrNames with 1 Go argument.
func (b *Builtins) AttrNames(arg1 any) (*value.Value, error) {
    return b.Call("attrNames", arg1)
}

// AttrValues calls builtins.attrValues with 1 Go argument.
func (b *Builtins) AttrValues(arg1 any) (*value.Value, error) {
    return b.Call("attrValues", arg1)
}

// BaseNameOf calls builtins.baseNameOf with 1 Go argument.
func (b *Builtins) BaseNameOf(arg1 any) (*value.Value, error) {
    return b.Call("baseNameOf", arg1)
}

// BitAnd calls builtins.bitAnd with 2 Go arguments.
func (b *Builtins) BitAnd(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("bitAnd", arg1, arg2)
}

// BitOr calls builtins.bitOr with 2 Go arguments.
func (b *Builtins) BitOr(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("bitOr", arg1, arg2)
}

// BitXor calls builtins.bitXor with 2 Go arguments.
func (b *Builtins) BitXor(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("bitXor", arg1, arg2)
}

// Break calls builtins.break with 1 Go argument.
func (b *Builtins) Break(arg1 any) (*value.Value, error) {
    return b.Call("break", arg1)
}

// CatAttrs calls builtins.catAttrs with 2 Go arguments.
func (b *Builtins) CatAttrs(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("catAttrs", arg1, arg2)
}

// Ceil calls builtins.ceil with 1 Go argument.
func (b *Builtins) Ceil(arg1 any) (*value.Value, error) {
    return b.Call("ceil", arg1)
}

// CompareVersions calls builtins.compareVersions with 2 Go arguments.
func (b *Builtins) CompareVersions(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("compareVersions", arg1, arg2)
}

// ConcatLists calls builtins.concatLists with 1 Go argument.
func (b *Builtins) ConcatLists(arg1 any) (*value.Value, error) {
    return b.Call("concatLists", arg1)
}

// ConcatMap calls builtins.concatMap with 2 Go arguments.
func (b *Builtins) ConcatMap(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("concatMap", arg1, arg2)
}

// ConcatStringsSep calls builtins.concatStringsSep with 2 Go arguments.
func (b *Builtins) ConcatStringsSep(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("concatStringsSep", arg1, arg2)
}

// ConvertHash calls builtins.convertHash with 1 Go argument.
func (b *Builtins) ConvertHash(arg1 any) (*value.Value, error) {
    return b.Call("convertHash", arg1)
}

// DeepSeq calls builtins.deepSeq with 2 Go arguments.
func (b *Builtins) DeepSeq(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("deepSeq", arg1, arg2)
}

// DerivationStrict calls builtins.derivationStrict with 1 Go argument.
func (b *Builtins) DerivationStrict(arg1 any) (*value.Value, error) {
    return b.Call("derivationStrict", arg1)
}

// DirOf calls builtins.dirOf with 1 Go argument.
func (b *Builtins) DirOf(arg1 any) (*value.Value, error) {
    return b.Call("dirOf", arg1)
}

// Div calls builtins.div with 2 Go arguments.
func (b *Builtins) Div(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("div", arg1, arg2)
}

// Elem calls builtins.elem with 2 Go arguments.
func (b *Builtins) Elem(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("elem", arg1, arg2)
}

// ElemAt calls builtins.elemAt with 2 Go arguments.
func (b *Builtins) ElemAt(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("elemAt", arg1, arg2)
}

// Exec calls builtins.exec with 1 Go argument.
func (b *Builtins) Exec(arg1 any) (*value.Value, error) {
    return b.Call("exec", arg1)
}

// FetchClosure calls builtins.fetchClosure with 1 Go argument.
func (b *Builtins) FetchClosure(arg1 any) (*value.Value, error) {
    return b.Call("fetchClosure", arg1)
}

// FetchGit calls builtins.fetchGit with 1 Go argument.
func (b *Builtins) FetchGit(arg1 any) (*value.Value, error) {
    return b.Call("fetchGit", arg1)
}

// FetchMercurial calls builtins.fetchMercurial with 1 Go argument.
func (b *Builtins) FetchMercurial(arg1 any) (*value.Value, error) {
    return b.Call("fetchMercurial", arg1)
}

// FetchTarball calls builtins.fetchTarball with 1 Go argument.
func (b *Builtins) FetchTarball(arg1 any) (*value.Value, error) {
    return b.Call("fetchTarball", arg1)
}

// FetchTree calls builtins.fetchTree with 1 Go argument.
func (b *Builtins) FetchTree(arg1 any) (*value.Value, error) {
    return b.Call("fetchTree", arg1)
}

// Filter calls builtins.filter with 2 Go arguments.
func (b *Builtins) Filter(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("filter", arg1, arg2)
}

// FilterSource calls builtins.filterSource with 2 Go arguments.
func (b *Builtins) FilterSource(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("filterSource", arg1, arg2)
}

// FindFile calls builtins.findFile with 2 Go arguments.
func (b *Builtins) FindFile(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("findFile", arg1, arg2)
}

// Floor calls builtins.floor with 1 Go argument.
func (b *Builtins) Floor(arg1 any) (*value.Value, error) {
    return b.Call("floor", arg1)
}

// FoldlPrime calls builtins.foldl' with 3 Go arguments.
func (b *Builtins) FoldlPrime(arg1 any, arg2 any, arg3 any) (*value.Value, error) {
    return b.Call("foldl'", arg1, arg2, arg3)
}

// FromJSON calls builtins.fromJSON with 1 Go argument.
func (b *Builtins) FromJSON(arg1 any) (*value.Value, error) {
    return b.Call("fromJSON", arg1)
}

// FromTOML calls builtins.fromTOML with 1 Go argument.
func (b *Builtins) FromTOML(arg1 any) (*value.Value, error) {
    return b.Call("fromTOML", arg1)
}

// FunctionArgs calls builtins.functionArgs with 1 Go argument.
func (b *Builtins) FunctionArgs(arg1 any) (*value.Value, error) {
    return b.Call("functionArgs", arg1)
}

// GenList calls builtins.genList with 2 Go arguments.
func (b *Builtins) GenList(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("genList", arg1, arg2)
}

// GenericClosure calls builtins.genericClosure with 1 Go argument.
func (b *Builtins) GenericClosure(arg1 any) (*value.Value, error) {
    return b.Call("genericClosure", arg1)
}

// GetAttr calls builtins.getAttr with 2 Go arguments.
func (b *Builtins) GetAttr(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("getAttr", arg1, arg2)
}

// GetContext calls builtins.getContext with 1 Go argument.
func (b *Builtins) GetContext(arg1 any) (*value.Value, error) {
    return b.Call("getContext", arg1)
}

// GetEnv calls builtins.getEnv with 1 Go argument.
func (b *Builtins) GetEnv(arg1 any) (*value.Value, error) {
    return b.Call("getEnv", arg1)
}

// GetFlake calls builtins.getFlake with 1 Go argument.
func (b *Builtins) GetFlake(arg1 any) (*value.Value, error) {
    return b.Call("getFlake", arg1)
}

// GroupBy calls builtins.groupBy with 2 Go arguments.
func (b *Builtins) GroupBy(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("groupBy", arg1, arg2)
}

// HashFile calls builtins.hashFile with 2 Go arguments.
func (b *Builtins) HashFile(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("hashFile", arg1, arg2)
}

// HashString calls builtins.hashString with 2 Go arguments.
func (b *Builtins) HashString(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("hashString", arg1, arg2)
}

// HasAttr calls builtins.hasAttr with 2 Go arguments.
func (b *Builtins) HasAttr(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("hasAttr", arg1, arg2)
}

// HasContext calls builtins.hasContext with 1 Go argument.
func (b *Builtins) HasContext(arg1 any) (*value.Value, error) {
    return b.Call("hasContext", arg1)
}

// Head calls builtins.head with 1 Go argument.
func (b *Builtins) Head(arg1 any) (*value.Value, error) {
    return b.Call("head", arg1)
}

// Import calls builtins.import with 1 Go argument.
func (b *Builtins) Import(arg1 any) (*value.Value, error) {
    return b.Call("import", arg1)
}

// ImportNative calls builtins.importNative with 2 Go arguments.
func (b *Builtins) ImportNative(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("importNative", arg1, arg2)
}

// IntersectAttrs calls builtins.intersectAttrs with 2 Go arguments.
func (b *Builtins) IntersectAttrs(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("intersectAttrs", arg1, arg2)
}

// IsAttrs calls builtins.isAttrs with 1 Go argument.
func (b *Builtins) IsAttrs(arg1 any) (*value.Value, error) {
    return b.Call("isAttrs", arg1)
}

// IsBool calls builtins.isBool with 1 Go argument.
func (b *Builtins) IsBool(arg1 any) (*value.Value, error) {
    return b.Call("isBool", arg1)
}

// IsFloat calls builtins.isFloat with 1 Go argument.
func (b *Builtins) IsFloat(arg1 any) (*value.Value, error) {
    return b.Call("isFloat", arg1)
}

// IsFunction calls builtins.isFunction with 1 Go argument.
func (b *Builtins) IsFunction(arg1 any) (*value.Value, error) {
    return b.Call("isFunction", arg1)
}

// IsInt calls builtins.isInt with 1 Go argument.
func (b *Builtins) IsInt(arg1 any) (*value.Value, error) {
    return b.Call("isInt", arg1)
}

// IsList calls builtins.isList with 1 Go argument.
func (b *Builtins) IsList(arg1 any) (*value.Value, error) {
    return b.Call("isList", arg1)
}

// IsNull calls builtins.isNull with 1 Go argument.
func (b *Builtins) IsNull(arg1 any) (*value.Value, error) {
    return b.Call("isNull", arg1)
}

// IsPath calls builtins.isPath with 1 Go argument.
func (b *Builtins) IsPath(arg1 any) (*value.Value, error) {
    return b.Call("isPath", arg1)
}

// IsString calls builtins.isString with 1 Go argument.
func (b *Builtins) IsString(arg1 any) (*value.Value, error) {
    return b.Call("isString", arg1)
}

// Length calls builtins.length with 1 Go argument.
func (b *Builtins) Length(arg1 any) (*value.Value, error) {
    return b.Call("length", arg1)
}

// LessThan calls builtins.lessThan with 2 Go arguments.
func (b *Builtins) LessThan(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("lessThan", arg1, arg2)
}

// ListToAttrs calls builtins.listToAttrs with 1 Go argument.
func (b *Builtins) ListToAttrs(arg1 any) (*value.Value, error) {
    return b.Call("listToAttrs", arg1)
}

// Map calls builtins.map with 2 Go arguments.
func (b *Builtins) Map(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("map", arg1, arg2)
}

// MapAttrs calls builtins.mapAttrs with 2 Go arguments.
func (b *Builtins) MapAttrs(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("mapAttrs", arg1, arg2)
}

// Match calls builtins.match with 2 Go arguments.
func (b *Builtins) Match(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("match", arg1, arg2)
}

// Mul calls builtins.mul with 2 Go arguments.
func (b *Builtins) Mul(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("mul", arg1, arg2)
}

// OutputOf calls builtins.outputOf with 2 Go arguments.
func (b *Builtins) OutputOf(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("outputOf", arg1, arg2)
}

// ParseDrvName calls builtins.parseDrvName with 1 Go argument.
func (b *Builtins) ParseDrvName(arg1 any) (*value.Value, error) {
    return b.Call("parseDrvName", arg1)
}

// ParseFlakeRef calls builtins.parseFlakeRef with 1 Go argument.
func (b *Builtins) ParseFlakeRef(arg1 any) (*value.Value, error) {
    return b.Call("parseFlakeRef", arg1)
}

// Partition calls builtins.partition with 2 Go arguments.
func (b *Builtins) Partition(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("partition", arg1, arg2)
}

// Path calls builtins.path with 1 Go argument.
func (b *Builtins) Path(arg1 any) (*value.Value, error) {
    return b.Call("path", arg1)
}

// PathExists calls builtins.pathExists with 1 Go argument.
func (b *Builtins) PathExists(arg1 any) (*value.Value, error) {
    return b.Call("pathExists", arg1)
}

// Placeholder calls builtins.placeholder with 1 Go argument.
func (b *Builtins) Placeholder(arg1 any) (*value.Value, error) {
    return b.Call("placeholder", arg1)
}

// ReadDir calls builtins.readDir with 1 Go argument.
func (b *Builtins) ReadDir(arg1 any) (*value.Value, error) {
    return b.Call("readDir", arg1)
}

// ReadFile calls builtins.readFile with 1 Go argument.
func (b *Builtins) ReadFile(arg1 any) (*value.Value, error) {
    return b.Call("readFile", arg1)
}

// ReadFileType calls builtins.readFileType with 1 Go argument.
func (b *Builtins) ReadFileType(arg1 any) (*value.Value, error) {
    return b.Call("readFileType", arg1)
}

// RemoveAttrs calls builtins.removeAttrs with 2 Go arguments.
func (b *Builtins) RemoveAttrs(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("removeAttrs", arg1, arg2)
}

// ReplaceStrings calls builtins.replaceStrings with 3 Go arguments.
func (b *Builtins) ReplaceStrings(arg1 any, arg2 any, arg3 any) (*value.Value, error) {
    return b.Call("replaceStrings", arg1, arg2, arg3)
}

// ScopedImport calls builtins.scopedImport with 2 Go arguments.
func (b *Builtins) ScopedImport(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("scopedImport", arg1, arg2)
}

// Seq calls builtins.seq with 2 Go arguments.
func (b *Builtins) Seq(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("seq", arg1, arg2)
}

// Sort calls builtins.sort with 2 Go arguments.
func (b *Builtins) Sort(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("sort", arg1, arg2)
}

// Split calls builtins.split with 2 Go arguments.
func (b *Builtins) Split(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("split", arg1, arg2)
}

// SplitVersion calls builtins.splitVersion with 1 Go argument.
func (b *Builtins) SplitVersion(arg1 any) (*value.Value, error) {
    return b.Call("splitVersion", arg1)
}

// StorePath calls builtins.storePath with 1 Go argument.
func (b *Builtins) StorePath(arg1 any) (*value.Value, error) {
    return b.Call("storePath", arg1)
}

// StringLength calls builtins.stringLength with 1 Go argument.
func (b *Builtins) StringLength(arg1 any) (*value.Value, error) {
    return b.Call("stringLength", arg1)
}

// Sub calls builtins.sub with 2 Go arguments.
func (b *Builtins) Sub(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("sub", arg1, arg2)
}

// Substring calls builtins.substring with 3 Go arguments.
func (b *Builtins) Substring(arg1 any, arg2 any, arg3 any) (*value.Value, error) {
    return b.Call("substring", arg1, arg2, arg3)
}

// Tail calls builtins.tail with 1 Go argument.
func (b *Builtins) Tail(arg1 any) (*value.Value, error) {
    return b.Call("tail", arg1)
}

// ToFile calls builtins.toFile with 2 Go arguments.
func (b *Builtins) ToFile(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("toFile", arg1, arg2)
}

// ToJSON calls builtins.toJSON with 1 Go argument.
func (b *Builtins) ToJSON(arg1 any) (*value.Value, error) {
    return b.Call("toJSON", arg1)
}

// ToPath calls builtins.toPath with 1 Go argument.
func (b *Builtins) ToPath(arg1 any) (*value.Value, error) {
    return b.Call("toPath", arg1)
}

// function toString() { [native code] } calls builtins.toString with 1 Go argument.
func (b *Builtins) function toString() { [native code] }(arg1 any) (*value.Value, error) {
    return b.Call("toString", arg1)
}

// ToXML calls builtins.toXML with 1 Go argument.
func (b *Builtins) ToXML(arg1 any) (*value.Value, error) {
    return b.Call("toXML", arg1)
}

// Trace calls builtins.trace with 2 Go arguments.
func (b *Builtins) Trace(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("trace", arg1, arg2)
}

// TraceVerbose calls builtins.traceVerbose with 2 Go arguments.
func (b *Builtins) TraceVerbose(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("traceVerbose", arg1, arg2)
}

// TryEval calls builtins.tryEval with 1 Go argument.
func (b *Builtins) TryEval(arg1 any) (*value.Value, error) {
    return b.Call("tryEval", arg1)
}

// TypeOf calls builtins.typeOf with 1 Go argument.
func (b *Builtins) TypeOf(arg1 any) (*value.Value, error) {
    return b.Call("typeOf", arg1)
}

// UnsafeDiscardOutputDependency calls builtins.unsafeDiscardOutputDependency with 1 Go argument.
func (b *Builtins) UnsafeDiscardOutputDependency(arg1 any) (*value.Value, error) {
    return b.Call("unsafeDiscardOutputDependency", arg1)
}

// UnsafeDiscardStringContext calls builtins.unsafeDiscardStringContext with 1 Go argument.
func (b *Builtins) UnsafeDiscardStringContext(arg1 any) (*value.Value, error) {
    return b.Call("unsafeDiscardStringContext", arg1)
}

// UnsafeGetAttrPos calls builtins.unsafeGetAttrPos with 2 Go arguments.
func (b *Builtins) UnsafeGetAttrPos(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("unsafeGetAttrPos", arg1, arg2)
}

// Warn calls builtins.warn with 2 Go arguments.
func (b *Builtins) Warn(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("warn", arg1, arg2)
}

// ZipAttrsWith calls builtins.zipAttrsWith with 2 Go arguments.
func (b *Builtins) ZipAttrsWith(arg1 any, arg2 any) (*value.Value, error) {
    return b.Call("zipAttrsWith", arg1, arg2)
}

// FetchURL calls builtins.fetchurl with 1 Go argument.
func (b *Builtins) FetchURL(arg1 any) (*value.Value, error) {
    return b.Call("fetchurl", arg1)
}

// FlakeRefToString calls builtins.flakeRefToString with 1 Go argument.
func (b *Builtins) FlakeRefToString(arg1 any) (*value.Value, error) {
    return b.Call("flakeRefToString", arg1)
}

// Derivation calls builtins.derivation with 1 Go argument.
func (b *Builtins) Derivation(arg1 any) (*value.Value, error) {
    return b.Call("derivation", arg1)
}

