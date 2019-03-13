package rig

import (
	"testing"

	"github.com/Pimmr/rig/validators"
	"github.com/pkg/errors"
)

func TestUintValue(t *testing.T) {
	for _, test := range []struct {
		value          uint
		expectedString string
		input          string
		expectedSet    uint
		expectedError  bool
	}{
		{
			value:          4,
			expectedString: "4",
			input:          "2",
			expectedSet:    2,
			expectedError:  false,
		},
		{
			value:          1,
			expectedString: "1",
			input:          "not-a-uint",
			expectedError:  true,
		},
	} {
		i := uintValue(test.value)

		if i.String() != test.expectedString {
			t.Errorf("Uint(&%d).String() = %q, expected %q", test.value, i, test.expectedString)
		}

		err := i.Set(test.input)
		if test.expectedError && err == nil {
			t.Errorf("Uint().Set(%q): expected error, got nil instead", test.input)
			continue
		}
		if !test.expectedError && err != nil {
			t.Errorf("Uint().Set(%q): unexpected error: %s", test.input, err)
			continue
		}
		if uint(i) != test.expectedSet {
			t.Errorf("Uint(&i).Set(%q): expected f to be %d, got %d instead", test.input, test.expectedSet, uint(i))
		}
	}
}

func TestUint(t *testing.T) {
	var v uint = 2
	flag := "flag"
	env := "ENV"
	usage := "usage"
	f := Uint(&v, flag, env, usage)

	if f.TypeHint == "" {
		t.Error("Uint().TypeHint = \"\": expected .TypeHint to be set")
	}
	if f.Name != flag {
		t.Errorf("Uint(...).Name = %q, expected %q", f.Name, flag)
	}
	if f.Env != env {
		t.Errorf("Uint(...).Env = %q, expected %q", f.Env, env)
	}
	if f.Usage != usage {
		t.Errorf("Uint(...).Usage = %q, expected %q", f.Usage, usage)
	}

	expectedString := "2"
	if f.String() != expectedString {
		t.Errorf("Uint(&2)).String() = %q, expected %q", f.String(), expectedString)
	}

	s := "1"
	err := f.Set(s)
	if err != nil {
		t.Errorf("Uint().Set(%q): unexpected error: %s", s, err)
	}
	if v != 1 {
		t.Errorf("Uint(&v).Set(%q): expected v to be %d, got %d instead", s, 1, v)
	}

	s = "notauint"
	err = f.Set(s)
	if err == nil {
		t.Errorf("Uint().Set(%q): expected error, got nil", s)
	}

	if f.IsBoolFlag() {
		t.Error("Uint().IsBoolFlag() = true, expected false")
	}
}

func TestUintValidators(t *testing.T) {
	testValidator := func(shouldFail bool) (validator validators.Uint, called *bool) {
		called = new(bool)
		return func(uint) error {
			*called = true
			if shouldFail {
				return errors.New("failing validator")
			}
			return nil
		}, called
	}

	t.Run("valid input passing validators", func(t *testing.T) {
		var val uint
		v1, v1Called := testValidator(false)
		v2, v2Called := testValidator(false)
		f := Uint(&val, "flag", "ENV", "testing uint validators", v1, v2)
		in := "1"
		err := f.Set(in)
		if err != nil {
			t.Errorf("Uint(..., v1, v2).Set(%q): unexpected error: %s", in, err)
		}
		if !*v1Called || !*v2Called {
			t.Errorf("Uint(..., v1, v2).Set(%q): some validator wasn't called (v1: %v, v2: %v)", in, *v1Called, *v2Called)
		}
	})

	t.Run("invalid input passing validators", func(t *testing.T) {
		var val uint
		v1, v1Called := testValidator(false)
		f := Uint(&val, "flag", "ENV", "testing uint validators", v1)
		in := ""
		err := f.Set(in)
		if err == nil {
			t.Errorf("Uint(..., v1).Set(%q): expected error, got nil", in)
		}
		if *v1Called {
			t.Errorf("Uint(..., v1).Set(%q): validator shouldn't have been called", in)
		}
	})

	t.Run("valid input failing validators", func(t *testing.T) {
		var val uint
		v1, v1Called := testValidator(true)
		f := Uint(&val, "flag", "ENV", "testing uint validators", v1)
		in := "2"
		err := f.Set(in)
		if err == nil {
			t.Errorf("Uint(..., failingV1).Set(%q): expected error, got nil", in)
		}
		if !*v1Called {
			t.Errorf("Uint(..., failingV1).Set(%q): validator should have been called", in)
		}
	})
}

func TestUintGenerator(t *testing.T) {
	g := UintGenerator()
	i := g()
	if _, ok := i.(*uintValue); !ok {
		t.Errorf("UintGenerator(): expected type *uintValue, got %T instead", i)
	}
}
