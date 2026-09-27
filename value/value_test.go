package value_test

import (
	"testing"

	"github.com/ep0ll/nixlang-go/value"
)

func TestTypeString(t *testing.T) {
	cases := []struct {
		typ  value.Type
		want string
	}{
		{value.TypeInt, "int"},
		{value.TypeString, "string"},
		{value.TypeAttrs, "attrs"},
		{value.TypeList, "list"},
		{value.TypeNull, "null"},
		{value.TypeFunction, "function"},
	}
	for _, tc := range cases {
		if got := tc.typ.String(); got != tc.want {
			t.Errorf("Type(%d).String() = %q, want %q", int(tc.typ), got, tc.want)
		}
	}
}

func TestUnknownTypeString(t *testing.T) {
	s := value.Type(999).String()
	if s == "" {
		t.Fatal("empty string for unknown type")
	}
}
