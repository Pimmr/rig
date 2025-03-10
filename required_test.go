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

	f = String(&s, "string-flag", "STRING_ENV", "testing Required on String")
	r = Required(Required(f))

	if f.Required {
		t.Errorf("String(...).Required = true, expected false")
	}
	if !r.Required {
		t.Errorf("Required(Required(String(...))).Required = false, expected true")
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

func TestPositional(t *testing.T) {
	var s string

	f := String(&s, "string-flag", "STRING_ENV", "testing Positional on String")
	r := Positional(f)

	if f.Positional {
		t.Errorf("String(...).Positional = true, expected false")
	}
	if !r.Positional {
		t.Errorf("Positional(String(...)).Positional = false, expected true")
	}

	f = String(&s, "string-flag", "STRING_ENV", "testing Positional on String")
	r = Positional(Positional(f))

	if f.Positional {
		t.Errorf("String(...).Positional = true, expected false")
	}
	if !r.Positional {
		t.Errorf("Positional(Positional(String(...))).Positional = false, expected true")
	}

	f = Var(new(stringValue), "var-flag", "VAR_ENV", "testing Positional on Var")
	r = Positional(f)

	if f.Positional {
		t.Errorf("Var(...).Positional = true, expected false")
	}
	if !r.Positional {
		t.Errorf("Positional(Var(...)).Positional = false, expected true")
	}
}
