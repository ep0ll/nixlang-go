package util_test

import (
	"testing"

	"github.com/ep0ll/nixlang-go/util"
)

func TestNixErrorFormat(t *testing.T) {
	e := &util.NixError{Code: util.ErrUnknown, Msg: "test"}
	if e.Error() == "" {
		t.Fatal("empty error string")
	}
}

func TestErrorCodeConstants(t *testing.T) {
	if util.OK != 0 {
		t.Fatalf("OK should be 0, got %d", util.OK)
	}
}
