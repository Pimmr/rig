package rig

import "testing"

func TestPointer(t *testing.T) {
	t.Run("Int", func(t *testing.T) {
		var i *int

		p := Pointer(Int(i, "i", "I", "integer"), &i)

		s := p.String()
		expectedNil := "<nil>"
		if s != expectedNil {
			t.Errorf("Pointer(Int(nil)).String() = %q, expected %q", s, expectedNil)
			return
		}

		s = "42"
		err := p.Set(s)
		if err != nil {
			t.Errorf("Pointer.Set(%q): unexpected error: %v", s, err)
			return
		}

		expected := 42
		if *i != expected {
			t.Errorf("Pointer.Set(%q) = %d, expected %d", s, *i, expected)
		}
	})

	t.Run("Var", func(t *testing.T) {
		var v *stringValue

		p := Pointer(Var(v, "v", "V", "var"), &v)

		s := p.String()
		expectedNil := "<nil>"
		if s != expectedNil {
			t.Errorf("Pointer(Var(nil)).String() = %q, expected %q", s, expectedNil)
			return
		}

		s = "foo"
		err := p.Set(s)
		if err != nil {
			t.Errorf("Pointer.Set(%q): unexpected error: %v", s, err)
			return
		}

		if string(*v) != s {
			t.Errorf("Pointer.Set(%q) = %q, expected %q", s, *v, s)
		}
	})
}
