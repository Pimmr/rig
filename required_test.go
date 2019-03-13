package rig

import (
	"strings"
	"testing"
)

func TestRequired(t *testing.T) {
	var s string

	f := String(&s, "string-flag", "STRING_ENV", "testing Required on String")
	r := Required(f)

	if f.Required {
		t.Errorf("String(...).Required = true, expected false")
	}
	if !r.Required {
		t.Errorf("Required(String(...)).Required = false, expected true")
	}

	expectedSuffix := ", required"
	if !strings.HasSuffix(r.TypeHint, expectedSuffix) {
		t.Errorf("Required(String(...)).TypeHint = %q, expected %q suffix", r.TypeHint, expectedSuffix)
	}

	f = Var(new(stringValue), "var-flag", "VAR_ENV", "testing Required on Var")
	r = Required(f)

	if f.Required {
		t.Errorf("Var(...).Required = true, expected false")
	}
	if !r.Required {
		t.Errorf("Required(Var(...)).Required = false, expected true")
	}

	expected := "required"
	if r.TypeHint != expected {
		t.Errorf("Required(Var(...)).TypeHint = %q, expected %q", r.TypeHint, expected)
	}
}
