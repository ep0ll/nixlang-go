package builtins_test

import (
	"testing"

	"github.com/ep0ll/nixlang-go/builtins"
	"github.com/ep0ll/nixlang-go/value"
)

// These assignments make the exported Go signatures part of the compile-time
// contract. In particular, simple Nix scalar arguments must not regress to any.
var (
	_ func(*builtins.Builtins, string, any) (*value.Value, error) = (*builtins.Builtins).HashFile
	_ func(*builtins.Builtins, string, string) (*value.Value, error) = (*builtins.Builtins).HashString
	_ func(*builtins.Builtins, string, string) (*value.Value, error) = (*builtins.Builtins).CompareVersions
	_ func(*builtins.Builtins, string, string) (*value.Value, error) = (*builtins.Builtins).ToFile
	_ func(*builtins.Builtins, int64, int64) (*value.Value, error) = (*builtins.Builtins).Substring
	_ func(*builtins.Builtins, int64, int64) (*value.Value, error) = (*builtins.Builtins).BitAnd
	_ func(*builtins.Builtins, builtins.Function, []any) (*value.Value, error) = (*builtins.Builtins).Map
	_ func(*builtins.Builtins, string, []any) (*value.Value, error) = (*builtins.Builtins).ConcatStringsSep
	_ func(*builtins.Builtins, any, any) (*value.Value, error) = (*builtins.Builtins).AddErrorContext
	_ func(*builtins.Builtins, any) (*value.Value, error) = (*builtins.Builtins).Abort
	_ func(*builtins.Builtins, any) (*value.Value, error) = (*builtins.Builtins).Throw
	_ func(*builtins.Builtins, string, any) (*value.Value, error) = (*builtins.Builtins).Warn
	_ func(*builtins.Builtins, builtins.ConvertHashArgs) (*value.Value, error) = (*builtins.Builtins).ConvertHash
	_ func(*builtins.Builtins, builtins.FetchClosureArgs) (*value.Value, error) = (*builtins.Builtins).FetchClosure
	_ func(*builtins.Builtins, builtins.PathArgs) (*value.Value, error) = (*builtins.Builtins).Path
	_ func(*builtins.Builtins, []builtins.SearchPathEntry, string) (*value.Value, error) = (*builtins.Builtins).FindFile
	_ func(*builtins.Builtins, string, []builtins.AttrSet) (*value.Value, error) = (*builtins.Builtins).CatAttrs
	_ func(*builtins.Builtins, []builtins.NameValue) (*value.Value, error) = (*builtins.Builtins).ListToAttrs
	_ func(*builtins.Builtins, builtins.Function, []builtins.AttrSet) (*value.Value, error) = (*builtins.Builtins).ZipAttrsWith
	_ func(*builtins.Builtins, [][]any) (*value.Value, error) = (*builtins.Builtins).ConcatLists
	_ func(*builtins.Builtins, []any, int64) (*value.Value, error) = (*builtins.Builtins).ElemAt
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

	fetchClosure := builtins.FetchClosureArgs{
		FromPath:  builtins.Path("/nix/store/example-source"),
		FromStore: "https://cache.nixos.org",
		ToPath:    "",
	}
	closureAttrs := fetchClosure.ToNixAttrs()
	if closureAttrs["toPath"] != "" || closureAttrs["fromStore"] != "https://cache.nixos.org" {
		t.Fatalf("FetchClosureArgs conversion = %#v", closureAttrs)
	}

	shallow := true
	fetchGit := builtins.FetchGitArgs{
		URL:        "https://example.com/source.git",
		Rev:        "0123456789abcdef",
		Submodules: &shallow,
	}
	gitAttrs := fetchGit.ToNixAttrs()
	if gitAttrs["url"] != "https://example.com/source.git" || gitAttrs["submodules"] != true {
		t.Fatalf("FetchGitArgs conversion = %#v", gitAttrs)
	}

	fetchTree := builtins.FetchTreeArgs{
		Type:    "git",
		URL:     "https://example.com/source.git",
		NarHash: "sha256-example",
	}
	treeAttrs := fetchTree.ToNixAttrs()
	if treeAttrs["type"] != "git" || treeAttrs["narHash"] != "sha256-example" {
		t.Fatalf("FetchTreeArgs conversion = %#v", treeAttrs)
	}

	flakeRef := builtins.FlakeRefAttrs{
		Type:  "github",
		Owner: "NixOS",
		Repo:  "nixpkgs",
		Ref:   "nixos-unstable",
	}
	flakeAttrs := flakeRef.ToNixAttrs()
	if flakeAttrs["type"] != "github" || flakeAttrs["owner"] != "NixOS" {
		t.Fatalf("FlakeRefAttrs conversion = %#v", flakeAttrs)
	}
}
