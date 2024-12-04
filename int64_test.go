package rig

import (
	"errors"
	"testing"

	"github.com/Pimmr/rig/validators"
)

func TestInt64Value(t *testing.T) {
	for _, test := range []struct {
		value          int64
		expectedString string
		input          string
		expectedSet    int64
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
			input:          "not-an-int64",
			expectedError:  true,
		},
	} {
		i := int64Value(test.value)

		if i.String() != test.expectedString {
			t.Errorf("Int64(&%d).String() = %q, expected %q", test.value, i, test.expectedString)
		}

		err := i.Set(test.input)
		if test.expectedError && err == nil {
			t.Errorf("Int64().Set(%q): expected error, got nil instead", test.input)
			continue
		}
		if !test.expectedError && err != nil {
			t.Errorf("Int64().Set(%q): unexpected error: %s", test.input, err)
			continue
		}
		if int64(i) != test.expectedSet {
			t.Errorf("Int64(&i).Set(%q): expected f to be %d, got %d instead", test.input, test.expectedSet, int64(i))
		}
	}
}

func TestInt64(t *testing.T) {
	var v int64 = 2
	flag := "flag"
	env := "ENV"
	usage := "usage"
	f := Int64(&v, flag, env, usage)

	if f.TypeHint == "" {
		t.Error("Int64().TypeHint = \"\": expected .TypeHint to be set")
	}
	if f.Name != flag {
		t.Errorf("Int64(...).Name = %q, expected %q", f.Name, flag)
	}
	if f.Env != env {
		t.Errorf("Int64(...).Env = %q, expected %q", f.Env, env)
	}
	if f.Usage != usage {
		t.Errorf("Int64(...).Usage = %q, expected %q", f.Usage, usage)
	}

	expectedString := "2"
	if f.String() != expectedString {
		t.Errorf("Int64(&2)).String() = %q, expected %q", f.String(), expectedString)
	}

	s := "1"
	err := f.Set(s)
	if err != nil {
		t.Errorf("Int64().Set(%q): unexpected error: %s", s, err)
	}
	if v != 1 {
		t.Errorf("Int64(&v).Set(%q): expected v to be %d, got %d instead", s, 1, v)
	}

	s = "notanint64"
	err = f.Set(s)
	if err == nil {
		t.Errorf("Int64().Set(%q): expected error, got nil", s)
	}

	if f.IsBoolFlag() {
		t.Error("Int64().IsBoolFlag() = true, expected false")
	}
}

func TestInt64Validators(t *testing.T) {
	testValidator := func(shouldFail bool) (validator validators.Int64, called *bool) {
		called = new(bool)
		return func(int64) error {
			*called = true
			if shouldFail {
				return errors.New("failing validator")
			}
			return nil
		}, called
	}

	t.Run("valid input passing validators", func(t *testing.T) {
		var val int64
		v1, v1Called := testValidator(false)
		v2, v2Called := testValidator(false)
		f := Int64(&val, "flag", "ENV", "testing int64 validators", v1, v2)
		in := "1"
		err := f.Set(in)
		if err != nil {
			t.Errorf("Int64(..., v1, v2).Set(%q): unexpected error: %s", in, err)
		}
		if !*v1Called || !*v2Called {
			t.Errorf("Int64(..., v1, v2).Set(%q): some validator wasn't called (v1: %v, v2: %v)", in, *v1Called, *v2Called)
		}
	})

	t.Run("invalid input passing validators", func(t *testing.T) {
		var val int64
		v1, v1Called := testValidator(false)
		f := Int64(&val, "flag", "ENV", "testing int64 validators", v1)
		in := ""
		err := f.Set(in)
		if err == nil {
			t.Errorf("Int64(..., v1).Set(%q): expected error, got nil", in)
		}
		if *v1Called {
			t.Errorf("Int64(..., v1).Set(%q): validator shouldn't have been called", in)
		}
	})

	t.Run("valid input failing validators", func(t *testing.T) {
		var val int64
		v1, v1Called := testValidator(true)
		f := Int64(&val, "flag", "ENV", "testing int64 validators", v1)
		in := "2"
		err := f.Set(in)
		if err == nil {
			t.Errorf("Int64(..., failingV1).Set(%q): expected error, got nil", in)
		}
		if !*v1Called {
			t.Errorf("Int64(..., failingV1).Set(%q): validator should have been called", in)
		}
	})
}

func TestInt64Generator(t *testing.T) {
	g := Int64Generator()
	i := g()
	if _, ok := i.(*int64Value); !ok {
		t.Errorf("Int64Generator(): expected type *int64Value, got %T instead", i)
	}
}
