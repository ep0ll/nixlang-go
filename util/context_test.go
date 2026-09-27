package util_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/ep0ll/nixlang-go/util"
)

func TestNixErrorFormat(t *testing.T) {
	e := &util.NixError{Code: util.ErrUnknown, Msg: "test"}
	s := e.Error()
	if s == "" {
		t.Fatal("empty error string")
	}
	if !strings.Contains(s, "test") {
		t.Fatalf("error string %q missing message", s)
	}
	if !strings.Contains(s, "nix:") {
		t.Fatalf("error string %q missing nix: prefix", s)
	}
}

func TestNixErrorAs(t *testing.T) {
	var err error = &util.NixError{Code: util.ErrKey, Msg: "missing"}
	var ne *util.NixError
	if !errors.As(err, &ne) {
		t.Fatal("errors.As failed for *NixError")
	}
	if ne.Code != util.ErrKey {
		t.Fatalf("Code = %d, want ErrKey", ne.Code)
	}
}

func TestErrorCodeConstants(t *testing.T) {
	if util.OK != 0 {
		t.Fatalf("OK should be 0, got %d", util.OK)
	}
}

func TestNixErrorIs(t *testing.T) {
	e1 := &util.NixError{Code: util.ErrKey, Msg: "a"}
	e2 := &util.NixError{Code: util.ErrKey, Msg: "b"}
	e3 := &util.NixError{Code: util.ErrOverflow, Msg: "c"}
	if !errors.Is(e1, e2) {
		t.Fatal("same Code should match via errors.Is")
	}
	if errors.Is(e1, e3) {
		t.Fatal("different Code should not match")
	}
	if !util.IsKeyError(e1) {
		t.Fatal("IsKeyError")
	}
	if util.IsRecoverable(e1) {
		t.Fatal("ErrKey is not recoverable")
	}
	rec := &util.NixError{Code: util.ErrRecoverable, Msg: "retry"}
	if !util.IsRecoverable(rec) {
		t.Fatal("IsRecoverable")
	}
	ne, ok := util.AsNixError(e1)
	if !ok || ne.Code != util.ErrKey {
		t.Fatal("AsNixError")
	}
}
