package config

import (
	"testing"

	"github.com/Pimmr/config/validators"
	"github.com/pkg/errors"
)

func TestIntValue(t *testing.T) {
	for _, test := range []struct {
		value          int
		expectedString string
		input          string
		expectedSet    int
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
			input:          "not-an-int",
			expectedError:  true,
		},
	} {
		i := intValue(test.value)

		if i.String() != test.expectedString {
			t.Errorf("Int(&%d).String() = %q, expected %q", test.value, i, test.expectedString)
		}

		err := i.Set(test.input)
		if test.expectedError && err == nil {
			t.Errorf("Int().Set(%q): expected error, got nil instead", test.input)
			continue
		}
		if !test.expectedError && err != nil {
			t.Errorf("Int().Set(%q): unexpected error: %s", test.input, err)
			continue
		}
		if int(i) != test.expectedSet {
			t.Errorf("Int(&i).Set(%q): expected f to be %d, got %d instead", test.input, test.expectedSet, int(i))
		}
	}
}

func TestInt(t *testing.T) {
	v := 2
	flag := "flag"
	env := "ENV"
	usage := "usage"
	f := Int(&v, flag, env, usage)

	if f.TypeHint == "" {
		t.Error("Int().TypeHint = \"\": expected .TypeHint to be set")
	}
	if f.Name != flag {
		t.Errorf("Int(...).Name = %q, expected %q", f.Name, flag)
	}
	if f.Env != env {
		t.Errorf("Int(...).Env = %q, expected %q", f.Env, env)
	}
	if f.Usage != usage {
		t.Errorf("Int(...).Usage = %q, expected %q", f.Usage, usage)
	}

	expectedString := "2"
	if f.String() != expectedString {
		t.Errorf("Int(&2)).String() = %q, expected %q", f.String(), expectedString)
	}

	s := "1"
	err := f.Set(s)
	if err != nil {
		t.Errorf("Int().Set(%q): unexpected error: %s", s, err)
	}
	if v != 1 {
		t.Errorf("Int(&v).Set(%q): expected v to be %d, got %d instead", s, 1, v)
	}

	s = "notanint"
	err = f.Set(s)
	if err == nil {
		t.Errorf("Int().Set(%q): expected error, got nil", s)
	}

	if f.IsBoolFlag() {
		t.Error("Bool().IsBoolFlag() = true, expected false")
	}
}

func TestIntValidators(t *testing.T) {
	testValidator := func(shouldFail bool) (validator validators.Int, called *bool) {
		called = new(bool)
		return func(int) error {
			*called = true
			if shouldFail {
				return errors.New("failing validator")
			}
			return nil
		}, called
	}

	t.Run("valid input passing validators", func(t *testing.T) {
		var val int
		v1, v1Called := testValidator(false)
		v2, v2Called := testValidator(false)
		f := Int(&val, "flag", "ENV", "testing int validators", v1, v2)
		in := "1"
		err := f.Set(in)
		if err != nil {
			t.Errorf("Int(..., v1, v2).Set(%q): unexpected error: %s", in, err)
		}
		if !*v1Called || !*v2Called {
			t.Errorf("Int(..., v1, v2).Set(%q): some validator wasn't called (v1: %v, v2: %v)", in, *v1Called, *v2Called)
		}
	})

	t.Run("invalid input passing validators", func(t *testing.T) {
		var val int
		v1, v1Called := testValidator(false)
		f := Int(&val, "flag", "ENV", "testing int validators", v1)
		in := ""
		err := f.Set(in)
		if err == nil {
			t.Errorf("Int(..., v1).Set(%q): expected error, got nil", in)
		}
		if *v1Called {
			t.Errorf("Int(..., v1).Set(%q): validator shouldn't have been called", in)
		}
	})

	t.Run("valid input failing validators", func(t *testing.T) {
		var val int
		v1, v1Called := testValidator(true)
		f := Int(&val, "flag", "ENV", "testing int validators", v1)
		in := "2"
		err := f.Set(in)
		if err == nil {
			t.Errorf("Int(..., failingV1).Set(%q): expected error, got nil", in)
		}
		if !*v1Called {
			t.Errorf("Int(..., failingV1).Set(%q): validator should have been called", in)
		}
	})
}

func TestIntGenerator(t *testing.T) {
	g := IntGenerator()
	i := g()
	if _, ok := i.(*intValue); !ok {
		t.Errorf("IntGenerator(): expected type *intValue, got %T instead", i)
	}
}
