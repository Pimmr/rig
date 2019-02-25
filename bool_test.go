package config

import "testing"

func TestBoolValue(t *testing.T) {
	for _, test := range []struct {
		value          bool
		expectedString string
	}{
		{value: true, expectedString: "true"},
		{value: false, expectedString: "false"},
	} {
		b := boolValue(test.value)

		if b.String() != test.expectedString {
			t.Errorf("Bool(&%v).String() = %q, expected %q", test.value, b, test.expectedString)
		}

		s := "true"
		err := b.Set(s)
		if err != nil {
			t.Errorf("Bool(&%v).Set(%q): unexpected error: %s", test.value, s, err)
		}
		if !b {
			t.Errorf("Bool(&%v).Set(%q): expected value to be true, got false instead", test.value, s)
		}

		if !b.IsBoolFlag() {
			t.Error("Bool().IsBoolFlag() = false, expected true")
		}
	}
}

func TestBool(t *testing.T) {
	v := true
	flag := "flag"
	env := "ENV"
	usage := "usage"
	b := Bool(&v, flag, env, usage)

	if b.TypeHint == "" {
		t.Error("Bool().TypeHint = \"\": expected .TypeHint to be set")
	}
	if b.Name != flag {
		t.Errorf("Bool(...).Name = %q, expected %q", b.Name, flag)
	}
	if b.Env != env {
		t.Errorf("Bool(...).Env = %q, expected %q", b.Env, env)
	}
	if b.Usage != usage {
		t.Errorf("Bool(...).Usage = %q, expected %q", b.Usage, usage)
	}

	expectedString := "true"
	if b.String() != expectedString {
		t.Errorf("Bool(&true)).String() = %q, expected %q", b.String(), expectedString)
	}

	s := "false"
	err := b.Set(s)
	if err != nil {
		t.Errorf("Bool().Set(%q): unexpected error: %s", s, err)
	}
	if v {
		t.Errorf("Bool(&v).Set(%q): expected v to be false, got %v instead", s, v)
	}

	s = "notabool"
	err = b.Set(s)
	if err == nil {
		t.Errorf("Bool().Set(%q): expected error, got nil", s)
	}

	if !b.IsBoolFlag() {
		t.Error("Bool().IsBoolFlag() = false, expected true")
	}
}

func TestBoolGenerator(t *testing.T) {
	g := BoolGenerator()
	b := g()
	if _, ok := b.(*boolValue); !ok {
		t.Errorf("BoolGenerator(): expected type *boolValue, got %T instead", b)
	}
}
