package config

import (
	"testing"

	"github.com/Pimmr/config/validators"
	"github.com/pkg/errors"
)

func TestStringValue(t *testing.T) {
	for _, test := range []struct {
		value          string
		expectedString string
		input          string
		expectedSet    string
	}{
		{
			value:          "foo",
			expectedString: "foo",
			input:          "bar",
			expectedSet:    "bar",
		},
	} {
		s := stringValue(test.value)

		if s.String() != test.expectedString {
			t.Errorf("String(&%q).String() = %q, expected %q", test.value, s, test.expectedString)
		}

		err := s.Set(test.input)
		if err != nil {
			t.Errorf("String().Set(%q): unexpected error: %s", test.input, err)
			continue
		}
		if string(s) != test.expectedSet {
			t.Errorf("String(&s).Set(%q): expected f to be %q, got %q instead", test.input, test.expectedSet, string(s))
		}
	}
}

func TestString(t *testing.T) {
	var v = "foo"
	flag := "flag"
	env := "ENV"
	usage := "usage"
	f := String(&v, flag, env, usage)

	if f.TypeHint == "" {
		t.Error("String().TypeHint = \"\": expected .TypeHint to be set")
	}
	if f.Name != flag {
		t.Errorf("String(...).Name = %q, expected %q", f.Name, flag)
	}
	if f.Env != env {
		t.Errorf("String(...).Env = %q, expected %q", f.Env, env)
	}
	if f.Usage != usage {
		t.Errorf("String(...).Usage = %q, expected %q", f.Usage, usage)
	}

	expectedString := "foo"
	if f.String() != expectedString {
		t.Errorf("String(&2)).String() = %q, expected %q", f.String(), expectedString)
	}

	s := "bar"
	err := f.Set(s)
	if err != nil {
		t.Errorf("String().Set(%q): unexpected error: %s", s, err)
	}
	if v != "bar" {
		t.Errorf("String(&v).Set(%q): expected v to be %q, got %q instead", s, "bar", v)
	}

	if f.IsBoolFlag() {
		t.Error("String().IsBoolFlag() = true, expected false")
	}
}

func TestStringValidators(t *testing.T) {
	testValidator := func(shouldFail bool) (validator validators.String, called *bool) {
		called = new(bool)
		return func(string) error {
			*called = true
			if shouldFail {
				return errors.New("failing validator")
			}
			return nil
		}, called
	}

	t.Run("passing validators", func(t *testing.T) {
		var val string
		v1, v1Called := testValidator(false)
		v2, v2Called := testValidator(false)
		f := String(&val, "flag", "ENV", "testing string validators", v1, v2)
		in := "foo"
		err := f.Set(in)
		if err != nil {
			t.Errorf("String(..., v1, v2).Set(%q): unexpected error: %s", in, err)
		}
		if !*v1Called || !*v2Called {
			t.Errorf("String(..., v1, v2).Set(%q): some validator wasn't called (v1: %v, v2: %v)", in, *v1Called, *v2Called)
		}
	})

	t.Run("failing validators", func(t *testing.T) {
		var val string
		v1, v1Called := testValidator(true)
		f := String(&val, "flag", "ENV", "testing string validators", v1)
		in := "bar"
		err := f.Set(in)
		if err == nil {
			t.Errorf("String(..., failingV1).Set(%q): expected error, got nil", in)
		}
		if !*v1Called {
			t.Errorf("String(..., failingV1).Set(%q): validator should have been called", in)
		}
	})
}

func TestStringGenerator(t *testing.T) {
	g := StringGenerator()
	s := g()
	if _, ok := s.(*stringValue); !ok {
		t.Errorf("StringGenerator(): expected type *stringValue, got %T instead", s)
	}
}
