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
	_ func(*builtins.Builtins, builtins.Function, any) (*value.Value, error) = (*builtins.Builtins).Map
	_ func(*builtins.Builtins, builtins.ConvertHashArgs) (*value.Value, error) = (*builtins.Builtins).ConvertHash
	_ func(*builtins.Builtins, builtins.FetchClosureArgs) (*value.Value, error) = (*builtins.Builtins).FetchClosure
	_ func(*builtins.Builtins, builtins.PathArgs) (*value.Value, error) = (*builtins.Builtins).Path
	_ func(*builtins.Builtins, []builtins.SearchPathEntry, string) (*value.Value, error) = (*builtins.Builtins).FindFile
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
}
