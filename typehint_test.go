package rig

import "testing"

func TestTypeHint(t *testing.T) {
	var s string

	f := String(&s, "string-flag", "STRING_ENV", "testing TypeHint on String")
	typeHint := "type hint"
	h := TypeHint(f, typeHint)

	if h.TypeHint != typeHint {
		t.Errorf("TypeHint(String(...)).TypeHint = %q, expected %q", h.TypeHint, typeHint)
	}
}
