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

func ExampleTypeHint() {
	var s string

	c := &Config{
		FlagSet: testingFlagset(),
		Flags: []*Flag{
			// This will overwrite the default "string" typehint in the usage
			TypeHint(String(&s, "contact", "CONTACT", "Administrative contact"), "email address"),
		},
	}

	c.Usage()

	// Output:
	// Usage: rig-test
	//   -contact "email address"    CONTACT="email address"    Administrative contact
}
