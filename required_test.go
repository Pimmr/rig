package rig

import (
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

	f = Var(new(stringValue), "var-flag", "VAR_ENV", "testing Required on Var")
	r = Required(f)

	if f.Required {
		t.Errorf("Var(...).Required = true, expected false")
	}
	if !r.Required {
		t.Errorf("Required(Var(...)).Required = false, expected true")
	}
}
