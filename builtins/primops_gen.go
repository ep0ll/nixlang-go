// Code generated from the Nix 2.35.6 builtin registry; DO NOT EDIT by hand.
package builtins

import "github.com/ep0ll/nixlang-go/value"

// Add calls builtins.add.
func (b *Builtins) Add(left any, right any) (*value.Value, error) {
	return b.Call("add", left, right)
}

// AddDrvOutputDependencies calls builtins.addDrvOutputDependencies.
func (b *Builtins) AddDrvOutputDependencies(input any) (*value.Value, error) {
	return b.Call("addDrvOutputDependencies", input)
}

// AddErrorContext calls builtins.addErrorContext.
func (b *Builtins) AddErrorContext(context any, value any) (*value.Value, error) {
	return b.Call("addErrorContext", context, value)
}

// Abort calls builtins.abort.
func (b *Builtins) Abort(message string) (*value.Value, error) {
	return b.Call("abort", message)
}

// Throw calls builtins.throw.
func (b *Builtins) Throw(message string) (*value.Value, error) {
	return b.Call("throw", message)
}

// All calls builtins.all.
func (b *Builtins) All(predicate Function, list []any) (*value.Value, error) {
	return b.Call("all", predicate, list)
}

// Any calls builtins.any.
func (b *Builtins) Any(predicate Function, list []any) (*value.Value, error) {
	return b.Call("any", predicate, list)
}

// AppendContext calls builtins.appendContext.
func (b *Builtins) AppendContext(original string, context StringContext) (*value.Value, error) {
	return b.Call("appendContext", original, context)
}

// AttrNames calls builtins.attrNames.
func (b *Builtins) AttrNames(set AttrSet) (*value.Value, error) {
	return b.Call("attrNames", set)
}

// AttrValues calls builtins.attrValues.
func (b *Builtins) AttrValues(set AttrSet) (*value.Value, error) {
	return b.Call("attrValues", set)
}

// BaseNameOf calls builtins.baseNameOf.
func (b *Builtins) BaseNameOf(path any) (*value.Value, error) {
	return b.Call("baseNameOf", path)
}

// BitAnd calls builtins.bitAnd.
func (b *Builtins) BitAnd(left int64, right int64) (*value.Value, error) {
	return b.Call("bitAnd", left, right)
}

// BitOr calls builtins.bitOr.
func (b *Builtins) BitOr(left int64, right int64) (*value.Value, error) {
	return b.Call("bitOr", left, right)
}

// BitXor calls builtins.bitXor.
func (b *Builtins) BitXor(left int64, right int64) (*value.Value, error) {
	return b.Call("bitXor", left, right)
}

// Break calls builtins.break.
func (b *Builtins) Break(value any) (*value.Value, error) {
	return b.Call("break", value)
}

// CatAttrs calls builtins.catAttrs.
func (b *Builtins) CatAttrs(name string, list []AttrSet) (*value.Value, error) {
	return b.Call("catAttrs", name, list)
}

// Ceil calls builtins.ceil.
func (b *Builtins) Ceil(number any) (*value.Value, error) {
	return b.Call("ceil", number)
}

// CompareVersions calls builtins.compareVersions.
func (b *Builtins) CompareVersions(left string, right string) (*value.Value, error) {
	return b.Call("compareVersions", left, right)
}

// ConcatLists calls builtins.concatLists.
func (b *Builtins) ConcatLists(lists [][]any) (*value.Value, error) {
	return b.Call("concatLists", lists)
}

// ConcatMap calls builtins.concatMap.
func (b *Builtins) ConcatMap(function Function, list []any) (*value.Value, error) {
	return b.Call("concatMap", function, list)
}

// ConcatStringsSep calls builtins.concatStringsSep.
func (b *Builtins) ConcatStringsSep(separator string, strings []any) (*value.Value, error) {
	return b.Call("concatStringsSep", separator, strings)
}

// ConvertHash calls builtins.convertHash.
func (b *Builtins) ConvertHash(args ConvertHashArgs) (*value.Value, error) {
	return b.Call("convertHash", args)
}

// DeepSeq calls builtins.deepSeq.
func (b *Builtins) DeepSeq(first any, second any) (*value.Value, error) {
	return b.Call("deepSeq", first, second)
}

// DerivationStrict calls builtins.derivationStrict.
func (b *Builtins) DerivationStrict(attrs DerivationAttrs) (*value.Value, error) {
	return b.Call("derivationStrict", attrs)
}

// DirOf calls builtins.dirOf.
func (b *Builtins) DirOf(path any) (*value.Value, error) {
	return b.Call("dirOf", path)
}

// Div calls builtins.div.
func (b *Builtins) Div(left any, right any) (*value.Value, error) {
	return b.Call("div", left, right)
}

// Elem calls builtins.elem.
func (b *Builtins) Elem(element any, list []any) (*value.Value, error) {
	return b.Call("elem", element, list)
}

// ElemAt calls builtins.elemAt.
func (b *Builtins) ElemAt(list []any, index int64) (*value.Value, error) {
	return b.Call("elemAt", list, index)
}

// Exec calls builtins.exec.
func (b *Builtins) Exec(command []any) (*value.Value, error) {
	return b.Call("exec", command)
}

// FetchClosure calls builtins.fetchClosure.
func (b *Builtins) FetchClosure(args FetchClosureArgs) (*value.Value, error) {
	return b.Call("fetchClosure", args)
}

// FetchGit calls builtins.fetchGit.
func (b *Builtins) FetchGit(args any) (*value.Value, error) {
	return b.Call("fetchGit", args)
}

// FetchMercurial calls builtins.fetchMercurial.
func (b *Builtins) FetchMercurial(args any) (*value.Value, error) {
	return b.Call("fetchMercurial", args)
}

// FetchTarball calls builtins.fetchTarball.
func (b *Builtins) FetchTarball(args any) (*value.Value, error) {
	return b.Call("fetchTarball", args)
}

// FetchTree calls builtins.fetchTree.
func (b *Builtins) FetchTree(input any) (*value.Value, error) {
	return b.Call("fetchTree", input)
}

// Filter calls builtins.filter.
func (b *Builtins) Filter(predicate Function, list []any) (*value.Value, error) {
	return b.Call("filter", predicate, list)
}

// FilterSource calls builtins.filterSource.
func (b *Builtins) FilterSource(predicate Function, path any) (*value.Value, error) {
	return b.Call("filterSource", predicate, path)
}

// FindFile calls builtins.findFile.
func (b *Builtins) FindFile(searchPath []SearchPathEntry, lookupPath string) (*value.Value, error) {
	return b.Call("findFile", searchPath, lookupPath)
}

// Floor calls builtins.floor.
func (b *Builtins) Floor(number any) (*value.Value, error) {
	return b.Call("floor", number)
}

// FoldlPrime calls builtins.foldl'.
func (b *Builtins) FoldlPrime(function Function, initial any, list []any) (*value.Value, error) {
	return b.Call("foldl'", function, initial, list)
}

// FromJSON calls builtins.fromJSON.
func (b *Builtins) FromJSON(json string) (*value.Value, error) {
	return b.Call("fromJSON", json)
}

// FromTOML calls builtins.fromTOML.
func (b *Builtins) FromTOML(toml string) (*value.Value, error) {
	return b.Call("fromTOML", toml)
}

// FunctionArgs calls builtins.functionArgs.
func (b *Builtins) FunctionArgs(function Function) (*value.Value, error) {
	return b.Call("functionArgs", function)
}

// GenList calls builtins.genList.
func (b *Builtins) GenList(generator Function, length int64) (*value.Value, error) {
	return b.Call("genList", generator, length)
}

// GenericClosure calls builtins.genericClosure.
func (b *Builtins) GenericClosure(args GenericClosureArgs) (*value.Value, error) {
	return b.Call("genericClosure", args)
}

// GetAttr calls builtins.getAttr.
func (b *Builtins) GetAttr(name string, set AttrSet) (*value.Value, error) {
	return b.Call("getAttr", name, set)
}

// GetContext calls builtins.getContext.
func (b *Builtins) GetContext(input string) (*value.Value, error) {
	return b.Call("getContext", input)
}

// GetEnv calls builtins.getEnv.
func (b *Builtins) GetEnv(name string) (*value.Value, error) {
	return b.Call("getEnv", name)
}

// GetFlake calls builtins.getFlake.
func (b *Builtins) GetFlake(ref any) (*value.Value, error) {
	return b.Call("getFlake", ref)
}

// GroupBy calls builtins.groupBy.
func (b *Builtins) GroupBy(function Function, list []any) (*value.Value, error) {
	return b.Call("groupBy", function, list)
}

// HashFile calls builtins.hashFile.
func (b *Builtins) HashFile(algorithm string, path any) (*value.Value, error) {
	return b.Call("hashFile", algorithm, path)
}

// HashString calls builtins.hashString.
func (b *Builtins) HashString(algorithm string, input string) (*value.Value, error) {
	return b.Call("hashString", algorithm, input)
}

// HasAttr calls builtins.hasAttr.
func (b *Builtins) HasAttr(name string, set AttrSet) (*value.Value, error) {
	return b.Call("hasAttr", name, set)
}

// HasContext calls builtins.hasContext.
func (b *Builtins) HasContext(input string) (*value.Value, error) {
	return b.Call("hasContext", input)
}

// Head calls builtins.head.
func (b *Builtins) Head(list []any) (*value.Value, error) {
	return b.Call("head", list)
}

// Import calls builtins.import.
func (b *Builtins) Import(path any) (*value.Value, error) {
	return b.Call("import", path)
}

// ImportNative calls builtins.importNative.
func (b *Builtins) ImportNative(path any, symbol string) (*value.Value, error) {
	return b.Call("importNative", path, symbol)
}

// IntersectAttrs calls builtins.intersectAttrs.
func (b *Builtins) IntersectAttrs(first AttrSet, second AttrSet) (*value.Value, error) {
	return b.Call("intersectAttrs", first, second)
}

// IsAttrs calls builtins.isAttrs.
func (b *Builtins) IsAttrs(value any) (*value.Value, error) {
	return b.Call("isAttrs", value)
}

// IsBool calls builtins.isBool.
func (b *Builtins) IsBool(value any) (*value.Value, error) {
	return b.Call("isBool", value)
}

// IsFloat calls builtins.isFloat.
func (b *Builtins) IsFloat(value any) (*value.Value, error) {
	return b.Call("isFloat", value)
}

// IsFunction calls builtins.isFunction.
func (b *Builtins) IsFunction(value any) (*value.Value, error) {
	return b.Call("isFunction", value)
}

// IsInt calls builtins.isInt.
func (b *Builtins) IsInt(value any) (*value.Value, error) {
	return b.Call("isInt", value)
}

// IsList calls builtins.isList.
func (b *Builtins) IsList(value any) (*value.Value, error) {
	return b.Call("isList", value)
}

// IsNull calls builtins.isNull.
func (b *Builtins) IsNull(value any) (*value.Value, error) {
	return b.Call("isNull", value)
}

// IsPath calls builtins.isPath.
func (b *Builtins) IsPath(value any) (*value.Value, error) {
	return b.Call("isPath", value)
}

// IsString calls builtins.isString.
func (b *Builtins) IsString(value any) (*value.Value, error) {
	return b.Call("isString", value)
}

// Length calls builtins.length.
func (b *Builtins) Length(list []any) (*value.Value, error) {
	return b.Call("length", list)
}

// LessThan calls builtins.lessThan.
func (b *Builtins) LessThan(left any, right any) (*value.Value, error) {
	return b.Call("lessThan", left, right)
}

// ListToAttrs calls builtins.listToAttrs.
func (b *Builtins) ListToAttrs(entries []NameValue) (*value.Value, error) {
	return b.Call("listToAttrs", entries)
}

// Map calls builtins.map.
func (b *Builtins) Map(function Function, list []any) (*value.Value, error) {
	return b.Call("map", function, list)
}

// MapAttrs calls builtins.mapAttrs.
func (b *Builtins) MapAttrs(function Function, attrs AttrSet) (*value.Value, error) {
	return b.Call("mapAttrs", function, attrs)
}

// Match calls builtins.match.
func (b *Builtins) Match(regex string, input string) (*value.Value, error) {
	return b.Call("match", regex, input)
}

// Mul calls builtins.mul.
func (b *Builtins) Mul(left any, right any) (*value.Value, error) {
	return b.Call("mul", left, right)
}

// OutputOf calls builtins.outputOf.
func (b *Builtins) OutputOf(derivationReference string, outputName string) (*value.Value, error) {
	return b.Call("outputOf", derivationReference, outputName)
}

// ParseDrvName calls builtins.parseDrvName.
func (b *Builtins) ParseDrvName(name string) (*value.Value, error) {
	return b.Call("parseDrvName", name)
}

// ParseFlakeRef calls builtins.parseFlakeRef.
func (b *Builtins) ParseFlakeRef(ref string) (*value.Value, error) {
	return b.Call("parseFlakeRef", ref)
}

// Partition calls builtins.partition.
func (b *Builtins) Partition(predicate Function, list []any) (*value.Value, error) {
	return b.Call("partition", predicate, list)
}

// Path calls builtins.path.
func (b *Builtins) Path(args PathArgs) (*value.Value, error) {
	return b.Call("path", args)
}

// PathExists calls builtins.pathExists.
func (b *Builtins) PathExists(path any) (*value.Value, error) {
	return b.Call("pathExists", path)
}

// Placeholder calls builtins.placeholder.
func (b *Builtins) Placeholder(output string) (*value.Value, error) {
	return b.Call("placeholder", output)
}

// ReadDir calls builtins.readDir.
func (b *Builtins) ReadDir(path any) (*value.Value, error) {
	return b.Call("readDir", path)
}

// ReadFile calls builtins.readFile.
func (b *Builtins) ReadFile(path any) (*value.Value, error) {
	return b.Call("readFile", path)
}

// ReadFileType calls builtins.readFileType.
func (b *Builtins) ReadFileType(path any) (*value.Value, error) {
	return b.Call("readFileType", path)
}

// RemoveAttrs calls builtins.removeAttrs.
func (b *Builtins) RemoveAttrs(set AttrSet, names []string) (*value.Value, error) {
	return b.Call("removeAttrs", set, names)
}

// ReplaceStrings calls builtins.replaceStrings.
func (b *Builtins) ReplaceStrings(from []string, to []string, input string) (*value.Value, error) {
	return b.Call("replaceStrings", from, to, input)
}

// ScopedImport calls builtins.scopedImport.
func (b *Builtins) ScopedImport(scope AttrSet, path any) (*value.Value, error) {
	return b.Call("scopedImport", scope, path)
}

// Seq calls builtins.seq.
func (b *Builtins) Seq(first any, second any) (*value.Value, error) {
	return b.Call("seq", first, second)
}

// Sort calls builtins.sort.
func (b *Builtins) Sort(comparator Function, list []any) (*value.Value, error) {
	return b.Call("sort", comparator, list)
}

// Split calls builtins.split.
func (b *Builtins) Split(regex string, input string) (*value.Value, error) {
	return b.Call("split", regex, input)
}

// SplitVersion calls builtins.splitVersion.
func (b *Builtins) SplitVersion(version string) (*value.Value, error) {
	return b.Call("splitVersion", version)
}

// StorePath calls builtins.storePath.
func (b *Builtins) StorePath(path any) (*value.Value, error) {
	return b.Call("storePath", path)
}

// StringLength calls builtins.stringLength.
func (b *Builtins) StringLength(input string) (*value.Value, error) {
	return b.Call("stringLength", input)
}

// Sub calls builtins.sub.
func (b *Builtins) Sub(left any, right any) (*value.Value, error) {
	return b.Call("sub", left, right)
}

// Substring calls builtins.substring.
func (b *Builtins) Substring(start int64, length int64, input string) (*value.Value, error) {
	return b.Call("substring", start, length, input)
}

// Tail calls builtins.tail.
func (b *Builtins) Tail(list []any) (*value.Value, error) {
	return b.Call("tail", list)
}

// ToFile calls builtins.toFile.
func (b *Builtins) ToFile(name string, contents string) (*value.Value, error) {
	return b.Call("toFile", name, contents)
}

// ToJSON calls builtins.toJSON.
func (b *Builtins) ToJSON(expr any) (*value.Value, error) {
	return b.Call("toJSON", expr)
}

// ToPath calls builtins.toPath.
func (b *Builtins) ToPath(path any) (*value.Value, error) {
	return b.Call("toPath", path)
}

// ToString calls builtins.toString.
func (b *Builtins) ToString(value any) (*value.Value, error) {
	return b.Call("toString", value)
}

// ToXML calls builtins.toXML.
func (b *Builtins) ToXML(expr any) (*value.Value, error) {
	return b.Call("toXML", expr)
}

// Trace calls builtins.trace.
func (b *Builtins) Trace(first any, second any) (*value.Value, error) {
	return b.Call("trace", first, second)
}

// TraceVerbose calls builtins.traceVerbose.
func (b *Builtins) TraceVerbose(expr any, value any) (*value.Value, error) {
	return b.Call("traceVerbose", expr, value)
}

// TryEval calls builtins.tryEval.
func (b *Builtins) TryEval(expr any) (*value.Value, error) {
	return b.Call("tryEval", expr)
}

// TypeOf calls builtins.typeOf.
func (b *Builtins) TypeOf(expr any) (*value.Value, error) {
	return b.Call("typeOf", expr)
}

// UnsafeDiscardOutputDependency calls builtins.unsafeDiscardOutputDependency.
func (b *Builtins) UnsafeDiscardOutputDependency(input any) (*value.Value, error) {
	return b.Call("unsafeDiscardOutputDependency", input)
}

// UnsafeDiscardStringContext calls builtins.unsafeDiscardStringContext.
func (b *Builtins) UnsafeDiscardStringContext(input any) (*value.Value, error) {
	return b.Call("unsafeDiscardStringContext", input)
}

// UnsafeGetAttrPos calls builtins.unsafeGetAttrPos.
func (b *Builtins) UnsafeGetAttrPos(name string, set AttrSet) (*value.Value, error) {
	return b.Call("unsafeGetAttrPos", name, set)
}

// Warn calls builtins.warn.
func (b *Builtins) Warn(message string, value any) (*value.Value, error) {
	return b.Call("warn", message, value)
}

// ZipAttrsWith calls builtins.zipAttrsWith.
func (b *Builtins) ZipAttrsWith(function Function, list []AttrSet) (*value.Value, error) {
	return b.Call("zipAttrsWith", function, list)
}

// FetchURL calls builtins.fetchurl.
func (b *Builtins) FetchURL(args any) (*value.Value, error) {
	return b.Call("fetchurl", args)
}

// FlakeRefToString calls builtins.flakeRefToString.
func (b *Builtins) FlakeRefToString(attrs FlakeRefAttrs) (*value.Value, error) {
	return b.Call("flakeRefToString", attrs)
}

// Derivation calls builtins.derivation.
func (b *Builtins) Derivation(attrs DerivationAttrs) (*value.Value, error) {
	return b.Call("derivation", attrs)
}
