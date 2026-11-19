package builtins_test

import (
	"net/url"
	"testing"

	"github.com/ep0ll/nixlang-go/builtins"
	"github.com/ep0ll/nixlang-go/value"
)

// These assignments make the exported Go signatures part of the compile-time
// contract. In particular, simple Nix scalar arguments must not regress to any.
var (
	_ func(*builtins.Builtins, builtins.Number, builtins.Number) (*value.Value, error) = (*builtins.Builtins).Add
	_ func(*builtins.Builtins, builtins.StringCoercible) (*value.Value, error) = (*builtins.Builtins).AddDrvOutputDependencies
	_ func(*builtins.Builtins, builtins.StringCoercible, any) (*value.Value, error) = (*builtins.Builtins).AddErrorContext
	_ func(*builtins.Builtins, builtins.StringCoercible) (*value.Value, error) = (*builtins.Builtins).Abort
	_ func(*builtins.Builtins, builtins.StringCoercible) (*value.Value, error) = (*builtins.Builtins).Throw
	_ func(*builtins.Builtins, builtins.Function, []any) (*value.Value, error) = (*builtins.Builtins).All
	_ func(*builtins.Builtins, builtins.Function, []any) (*value.Value, error) = (*builtins.Builtins).Any
	_ func(*builtins.Builtins, string, builtins.StringContext) (*value.Value, error) = (*builtins.Builtins).AppendContext
	_ func(*builtins.Builtins, builtins.AttrSet) (*value.Value, error) = (*builtins.Builtins).AttrNames
	_ func(*builtins.Builtins, builtins.AttrSet) (*value.Value, error) = (*builtins.Builtins).AttrValues
	_ func(*builtins.Builtins, builtins.PathLike) (*value.Value, error) = (*builtins.Builtins).BaseNameOf
	_ func(*builtins.Builtins, int64, int64) (*value.Value, error) = (*builtins.Builtins).BitAnd
	_ func(*builtins.Builtins, int64, int64) (*value.Value, error) = (*builtins.Builtins).BitOr
	_ func(*builtins.Builtins, int64, int64) (*value.Value, error) = (*builtins.Builtins).BitXor
	_ func(*builtins.Builtins, any) (*value.Value, error) = (*builtins.Builtins).Break
	_ func(*builtins.Builtins, string, []builtins.AttrSet) (*value.Value, error) = (*builtins.Builtins).CatAttrs
	_ func(*builtins.Builtins, builtins.Number) (*value.Value, error) = (*builtins.Builtins).Ceil
	_ func(*builtins.Builtins, string, string) (*value.Value, error) = (*builtins.Builtins).CompareVersions
	_ func(*builtins.Builtins, [][]any) (*value.Value, error) = (*builtins.Builtins).ConcatLists
	_ func(*builtins.Builtins, builtins.Function, []any) (*value.Value, error) = (*builtins.Builtins).ConcatMap
	_ func(*builtins.Builtins, string, []builtins.StringCoercible) (*value.Value, error) = (*builtins.Builtins).ConcatStringsSep
	_ func(*builtins.Builtins, builtins.ConvertHashArgs) (*value.Value, error) = (*builtins.Builtins).ConvertHash
	_ func(*builtins.Builtins, any, any) (*value.Value, error) = (*builtins.Builtins).DeepSeq
	_ func(*builtins.Builtins, builtins.DerivationAttrs) (*value.Value, error) = (*builtins.Builtins).DerivationStrict
	_ func(*builtins.Builtins, builtins.PathLike) (*value.Value, error) = (*builtins.Builtins).DirOf
	_ func(*builtins.Builtins, builtins.Number, builtins.Number) (*value.Value, error) = (*builtins.Builtins).Div
	_ func(*builtins.Builtins, any, []any) (*value.Value, error) = (*builtins.Builtins).Elem
	_ func(*builtins.Builtins, []any, int64) (*value.Value, error) = (*builtins.Builtins).ElemAt
	_ func(*builtins.Builtins, []any) (*value.Value, error) = (*builtins.Builtins).Exec
	_ func(*builtins.Builtins, builtins.FetchClosureArgs) (*value.Value, error) = (*builtins.Builtins).FetchClosure
	_ func(*builtins.Builtins, builtins.FetchGitArgs) (*value.Value, error) = (*builtins.Builtins).FetchGit
	_ func(*builtins.Builtins, builtins.FetchMercurialArgs) (*value.Value, error) = (*builtins.Builtins).FetchMercurial
	_ func(*builtins.Builtins, builtins.FetchTarballArgs) (*value.Value, error) = (*builtins.Builtins).FetchTarball
	_ func(*builtins.Builtins, builtins.FetchTreeArgs) (*value.Value, error) = (*builtins.Builtins).FetchTree
	_ func(*builtins.Builtins, builtins.Function, []any) (*value.Value, error) = (*builtins.Builtins).Filter
	_ func(*builtins.Builtins, builtins.Function, builtins.PathLike) (*value.Value, error) = (*builtins.Builtins).FilterSource
	_ func(*builtins.Builtins, []builtins.SearchPathEntry, string) (*value.Value, error) = (*builtins.Builtins).FindFile
	_ func(*builtins.Builtins, builtins.Number) (*value.Value, error) = (*builtins.Builtins).Floor
	_ func(*builtins.Builtins, builtins.Function, any, []any) (*value.Value, error) = (*builtins.Builtins).FoldlPrime
	_ func(*builtins.Builtins, string) (*value.Value, error) = (*builtins.Builtins).FromJSON
	_ func(*builtins.Builtins, string) (*value.Value, error) = (*builtins.Builtins).FromTOML
	_ func(*builtins.Builtins, builtins.Function) (*value.Value, error) = (*builtins.Builtins).FunctionArgs
	_ func(*builtins.Builtins, builtins.Function, int64) (*value.Value, error) = (*builtins.Builtins).GenList
	_ func(*builtins.Builtins, builtins.GenericClosureArgs) (*value.Value, error) = (*builtins.Builtins).GenericClosure
	_ func(*builtins.Builtins, string, builtins.AttrSet) (*value.Value, error) = (*builtins.Builtins).GetAttr
	_ func(*builtins.Builtins, string) (*value.Value, error) = (*builtins.Builtins).GetContext
	_ func(*builtins.Builtins, string) (*value.Value, error) = (*builtins.Builtins).GetEnv
	_ func(*builtins.Builtins, builtins.FlakeReference) (*value.Value, error) = (*builtins.Builtins).GetFlake
	_ func(*builtins.Builtins, builtins.Function, []any) (*value.Value, error) = (*builtins.Builtins).GroupBy
	_ func(*builtins.Builtins, builtins.HashAlgorithm, builtins.PathLike) (*value.Value, error) = (*builtins.Builtins).HashFile
	_ func(*builtins.Builtins, builtins.HashAlgorithm, string) (*value.Value, error) = (*builtins.Builtins).HashString
	_ func(*builtins.Builtins, string, builtins.AttrSet) (*value.Value, error) = (*builtins.Builtins).HasAttr
	_ func(*builtins.Builtins, string) (*value.Value, error) = (*builtins.Builtins).HasContext
	_ func(*builtins.Builtins, []any) (*value.Value, error) = (*builtins.Builtins).Head
	_ func(*builtins.Builtins, builtins.PathLike) (*value.Value, error) = (*builtins.Builtins).Import
	_ func(*builtins.Builtins, builtins.PathLike, string) (*value.Value, error) = (*builtins.Builtins).ImportNative
	_ func(*builtins.Builtins, builtins.AttrSet, builtins.AttrSet) (*value.Value, error) = (*builtins.Builtins).IntersectAttrs
	_ func(*builtins.Builtins, any) (*value.Value, error) = (*builtins.Builtins).IsAttrs
	_ func(*builtins.Builtins, any) (*value.Value, error) = (*builtins.Builtins).IsBool
	_ func(*builtins.Builtins, any) (*value.Value, error) = (*builtins.Builtins).IsFloat
	_ func(*builtins.Builtins, any) (*value.Value, error) = (*builtins.Builtins).IsFunction
	_ func(*builtins.Builtins, any) (*value.Value, error) = (*builtins.Builtins).IsInt
	_ func(*builtins.Builtins, any) (*value.Value, error) = (*builtins.Builtins).IsList
	_ func(*builtins.Builtins, any) (*value.Value, error) = (*builtins.Builtins).IsNull
	_ func(*builtins.Builtins, any) (*value.Value, error) = (*builtins.Builtins).IsPath
	_ func(*builtins.Builtins, any) (*value.Value, error) = (*builtins.Builtins).IsString
	_ func(*builtins.Builtins, []any) (*value.Value, error) = (*builtins.Builtins).Length
	_ func(*builtins.Builtins, builtins.ComparisonValue, builtins.ComparisonValue) (*value.Value, error) = (*builtins.Builtins).LessThan
	_ func(*builtins.Builtins, []builtins.NameValue) (*value.Value, error) = (*builtins.Builtins).ListToAttrs
	_ func(*builtins.Builtins, builtins.Function, []any) (*value.Value, error) = (*builtins.Builtins).Map
	_ func(*builtins.Builtins, builtins.Function, builtins.AttrSet) (*value.Value, error) = (*builtins.Builtins).MapAttrs
	_ func(*builtins.Builtins, string, string) (*value.Value, error) = (*builtins.Builtins).Match
	_ func(*builtins.Builtins, builtins.Number, builtins.Number) (*value.Value, error) = (*builtins.Builtins).Mul
	_ func(*builtins.Builtins, string, string) (*value.Value, error) = (*builtins.Builtins).OutputOf
	_ func(*builtins.Builtins, string) (*value.Value, error) = (*builtins.Builtins).ParseDrvName
	_ func(*builtins.Builtins, string) (*value.Value, error) = (*builtins.Builtins).ParseFlakeRef
	_ func(*builtins.Builtins, builtins.Function, []any) (*value.Value, error) = (*builtins.Builtins).Partition
	_ func(*builtins.Builtins, builtins.PathArgs) (*value.Value, error) = (*builtins.Builtins).Path
	_ func(*builtins.Builtins, builtins.PathLike) (*value.Value, error) = (*builtins.Builtins).PathExists
	_ func(*builtins.Builtins, string) (*value.Value, error) = (*builtins.Builtins).Placeholder
	_ func(*builtins.Builtins, builtins.PathLike) (*value.Value, error) = (*builtins.Builtins).ReadDir
	_ func(*builtins.Builtins, builtins.PathLike) (*value.Value, error) = (*builtins.Builtins).ReadFile
	_ func(*builtins.Builtins, builtins.PathLike) (*value.Value, error) = (*builtins.Builtins).ReadFileType
	_ func(*builtins.Builtins, builtins.AttrSet, []string) (*value.Value, error) = (*builtins.Builtins).RemoveAttrs
	_ func(*builtins.Builtins, []string, []string, string) (*value.Value, error) = (*builtins.Builtins).ReplaceStrings
	_ func(*builtins.Builtins, builtins.AttrSet, builtins.PathLike) (*value.Value, error) = (*builtins.Builtins).ScopedImport
	_ func(*builtins.Builtins, any, any) (*value.Value, error) = (*builtins.Builtins).Seq
	_ func(*builtins.Builtins, builtins.Function, []any) (*value.Value, error) = (*builtins.Builtins).Sort
	_ func(*builtins.Builtins, string, string) (*value.Value, error) = (*builtins.Builtins).Split
	_ func(*builtins.Builtins, string) (*value.Value, error) = (*builtins.Builtins).SplitVersion
	_ func(*builtins.Builtins, builtins.PathLike) (*value.Value, error) = (*builtins.Builtins).StorePath
	_ func(*builtins.Builtins, string) (*value.Value, error) = (*builtins.Builtins).StringLength
	_ func(*builtins.Builtins, builtins.Number, builtins.Number) (*value.Value, error) = (*builtins.Builtins).Sub
	_ func(*builtins.Builtins, int64, int64, string) (*value.Value, error) = (*builtins.Builtins).Substring
	_ func(*builtins.Builtins, []any) (*value.Value, error) = (*builtins.Builtins).Tail
	_ func(*builtins.Builtins, string, string) (*value.Value, error) = (*builtins.Builtins).ToFile
	_ func(*builtins.Builtins, any) (*value.Value, error) = (*builtins.Builtins).ToJSON
	_ func(*builtins.Builtins, builtins.PathLike) (*value.Value, error) = (*builtins.Builtins).ToPath
	_ func(*builtins.Builtins, builtins.StringCoercible) (*value.Value, error) = (*builtins.Builtins).ToString
	_ func(*builtins.Builtins, any) (*value.Value, error) = (*builtins.Builtins).ToXML
	_ func(*builtins.Builtins, any, any) (*value.Value, error) = (*builtins.Builtins).Trace
	_ func(*builtins.Builtins, any, any) (*value.Value, error) = (*builtins.Builtins).TraceVerbose
	_ func(*builtins.Builtins, any) (*value.Value, error) = (*builtins.Builtins).TryEval
	_ func(*builtins.Builtins, any) (*value.Value, error) = (*builtins.Builtins).TypeOf
	_ func(*builtins.Builtins, builtins.StringCoercible) (*value.Value, error) = (*builtins.Builtins).UnsafeDiscardOutputDependency
	_ func(*builtins.Builtins, builtins.StringCoercible) (*value.Value, error) = (*builtins.Builtins).UnsafeDiscardStringContext
	_ func(*builtins.Builtins, string, builtins.AttrSet) (*value.Value, error) = (*builtins.Builtins).UnsafeGetAttrPos
	_ func(*builtins.Builtins, string, any) (*value.Value, error) = (*builtins.Builtins).Warn
	_ func(*builtins.Builtins, builtins.Function, []builtins.AttrSet) (*value.Value, error) = (*builtins.Builtins).ZipAttrsWith
	_ func(*builtins.Builtins, builtins.FetchURLArgs) (*value.Value, error) = (*builtins.Builtins).FetchURL
	_ func(*builtins.Builtins, builtins.FlakeRefAttrs) (*value.Value, error) = (*builtins.Builtins).FlakeRefToString
	_ func(*builtins.Builtins, builtins.DerivationAttrs) (*value.Value, error) = (*builtins.Builtins).Derivation
)

func TestTypedArgumentStructs(t *testing.T) {
	recursive := false
	attrs := builtins.PathArgs{
		Path:      builtins.Path("./source"),
		Name:      "source",
		Recursive: &recursive,
		SHA256:    "sha256-placeholder",
	}
	if got := attrs.ToNixAttrs(); got["path"] != builtins.Path("./source") {
		t.Fatalf("PathArgs path = %#v", got["path"])
	}

	nameValue := builtins.NameValue{Name: "answer", Value: int64(42)}
	if got := nameValue.ToNixAttrs(); got["name"] != "answer" || got["value"] != int64(42) {
		t.Fatalf("NameValue conversion = %#v", got)
	}

	context := builtins.StringContext{
		"/nix/store/example": {
			Path:    true,
			Outputs: []string{"out"},
		},
	}
	entry := context["/nix/store/example"].ToNixAttrs()
	if entry["path"] != true {
		t.Fatalf("StringContext path = %#v, want true", entry["path"])
	}
	if got := entry["outputs"].([]string); len(got) != 1 || got[0] != "out" {
		t.Fatalf("StringContext outputs = %#v", entry["outputs"])
	}

	contextAttrs := context.ToNixAttrs()
	convertedEntry, ok := contextAttrs["/nix/store/example"].(builtins.AttrSet)
	if !ok || convertedEntry["path"] != true {
		t.Fatalf("StringContext conversion = %#v", contextAttrs)
	}

	derivation := builtins.DerivationAttrs{
		Name:   "example",
		System: "x86_64-linux",
		Builder: builtins.Path("/bin/sh"),
		Args:   []string{"-c", "true"},
		Env:    builtins.AttrSet{"MESSAGE": "hello"},
		Extra:  builtins.AttrSet{"preferLocalBuild": true},
	}
	drvAttrs := derivation.ToNixAttrs()
	if drvAttrs["name"] != "example" || drvAttrs["MESSAGE"] != "hello" {
		t.Fatalf("DerivationAttrs conversion = %#v", drvAttrs)
	}
	searchPath := builtins.SearchPathEntry{
		Prefix: "nixpkgs",
		Path:   builtins.Path("./nixpkgs"),
	}
	if got := searchPath.ToNixAttrs(); got["prefix"] != "nixpkgs" {
		t.Fatalf("SearchPathEntry conversion = %#v", got)
	}

	emptyToPath := builtins.StorePath("")
	fetchClosure := builtins.FetchClosureArgs{
		FromPath:  builtins.StorePath("/nix/store/example-source"),
		FromStore: url.URL{Scheme: "https", Host: "cache.nixos.org"},
		ToPath:    &emptyToPath,
	}
	closureAttrs := fetchClosure.ToNixAttrs()
	if closureAttrs["toPath"] != "" || closureAttrs["fromStore"] != "https://cache.nixos.org" {
		t.Fatalf("FetchClosureArgs conversion = %#v", closureAttrs)
	}

	shallow := true
	fetchGit := builtins.FetchGitArgs{
		URL:        &url.URL{Scheme: "https", Host: "example.com", Path: "/source.git"},
		Rev:        "0123456789abcdef",
		Submodules: &shallow,
	}
	gitAttrs := fetchGit.ToNixAttrs()
	if gitAttrs["url"] != "https://example.com/source.git" || gitAttrs["submodules"] != true {
		t.Fatalf("FetchGitArgs conversion = %#v", gitAttrs)
	}

	fetchTree := builtins.FetchTreeArgs{
		Type:    "git",
		URL:     &url.URL{Scheme: "https", Host: "example.com", Path: "/source.git"},
		NarHash: "sha256-example",
	}
	treeAttrs := fetchTree.ToNixAttrs()
	if treeAttrs["type"] != "git" || treeAttrs["narHash"] != "sha256-example" {
		t.Fatalf("FetchTreeArgs conversion = %#v", treeAttrs)
	}

	flakeRef := builtins.FlakeRefAttrs{
		Type:  builtins.FlakeRefGitHub,
		Owner: "NixOS",
		Repo:  "nixpkgs",
		Ref:   "nixos-unstable",
	}
	flakeAttrs := flakeRef.ToNixAttrs()
	if flakeAttrs["type"] != "github" || flakeAttrs["owner"] != "NixOS" {
		t.Fatalf("FlakeRefAttrs conversion = %#v", flakeAttrs)
	}
}
