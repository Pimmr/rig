package config

import (
	"testing"

	"github.com/Pimmr/config/validators"
	"github.com/pkg/errors"
)

func TestFloat64Value(t *testing.T) {
	for _, test := range []struct {
		value          float64
		expectedString string
		input          string
		expectedSet    float64
		expectedError  bool
	}{
		{
			value:          4.2,
			expectedString: "4.2",
			input:          "2.3",
			expectedSet:    2.3,
			expectedError:  false,
		},
		{
			value:          1.1,
			expectedString: "1.1",
			input:          "not-a-float",
			expectedError:  true,
		},
	} {
		f := float64Value(test.value)

		if f.String() != test.expectedString {
			t.Errorf("Float64(&%f).String() = %q, expected %q", test.value, f, test.expectedString)
		}

		err := f.Set(test.input)
		if test.expectedError && err == nil {
			t.Errorf("Float64().Set(%q): expected error, got nil instead", test.input)
			continue
		}
		if !test.expectedError && err != nil {
			t.Errorf("Float64().Set(%q): unexpected error: %s", test.input, err)
			continue
		}
		if float64(f) != test.expectedSet {
			t.Errorf("Float64(&f).Set(%q): expected f to be %f, got %f instead", test.input, test.expectedSet, float64(f))
		}
	}
}

func TestFloat64(t *testing.T) {
	v := 2.4
	flag := "flag"
	env := "ENV"
	usage := "usage"
	f := Float64(&v, flag, env, usage)

	if f.TypeHint == "" {
		t.Error("Float64().TypeHint = \"\": expected .TypeHint to be set")
	}
	if f.Name != flag {
		t.Errorf("Float64(...).Name = %q, expected %q", f.Name, flag)
	}
	if f.Env != env {
		t.Errorf("Float64(...).Env = %q, expected %q", f.Env, env)
	}
	if f.Usage != usage {
		t.Errorf("Float64(...).Usage = %q, expected %q", f.Usage, usage)
	}

	expectedString := "2.4"
	if f.String() != expectedString {
		t.Errorf("Float64(&2.4)).String() = %q, expected %q", f.String(), expectedString)
	}

	s := "1.2"
	err := f.Set(s)
	if err != nil {
		t.Errorf("Float64().Set(%q): unexpected error: %s", s, err)
	}
	if v != 1.2 {
		t.Errorf("Float64(&v).Set(%q): expected v to be %f, got %f instead", s, 1.2, v)
	}

	s = "notafloat"
	err = f.Set(s)
	if err == nil {
		t.Errorf("Float64().Set(%q): expected error, got nil", s)
	}

	if f.IsBoolFlag() {
		t.Error("Float64().IsBoolFlag() = true, expected false")
	}
}

func TestFloat64Validators(t *testing.T) {
	testValidator := func(shouldFail bool) (validator validators.Float64, called *bool) {
		called = new(bool)
		return func(float64) error {
			*called = true
			if shouldFail {
				return errors.New("failing validator")
			}
			return nil
		}, called
	}

	t.Run("valid input passing validators", func(t *testing.T) {
		var val float64
		v1, v1Called := testValidator(false)
		v2, v2Called := testValidator(false)
		f := Float64(&val, "flag", "ENV", "testing float64 validators", v1, v2)
		in := "1.2"
		err := f.Set(in)
		if err != nil {
			t.Errorf("Float64(..., v1, v2).Set(%q): unexpected error: %s", in, err)
		}
		if !*v1Called || !*v2Called {
			t.Errorf("Float64(..., v1, v2).Set(%q): some validator wasn't called (v1: %v, v2: %v)", in, *v1Called, *v2Called)
		}
	})

	t.Run("invalid input passing validators", func(t *testing.T) {
		var val float64
		v1, v1Called := testValidator(false)
		f := Float64(&val, "flag", "ENV", "testing float64 validators", v1)
		in := ""
		err := f.Set(in)
		if err == nil {
			t.Errorf("Float64(..., v1).Set(%q): expected error, got nil", in)
		}
		if *v1Called {
			t.Errorf("Float64(..., v1).Set(%q): validator shouldn't have been called", in)
		}
	})

	t.Run("valid input failing validators", func(t *testing.T) {
		var val float64
		v1, v1Called := testValidator(true)
		f := Float64(&val, "flag", "ENV", "testing float64 validators", v1)
		in := "2.1"
		err := f.Set(in)
		if err == nil {
			t.Errorf("Float64(..., failingV1).Set(%q): expected error, got nil", in)
		}
		if !*v1Called {
			t.Errorf("Float64(..., failingV1).Set(%q): validator should have been called", in)
		}
	})
}

func TestFloat64Generator(t *testing.T) {
	g := Float64Generator()
	f := g()
	if _, ok := f.(*float64Value); !ok {
		t.Errorf("Float64Generator(): expected type *float64Value, got %T instead", f)
	}
}
