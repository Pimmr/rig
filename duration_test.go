package config

import (
	"testing"
	"time"

	"github.com/Pimmr/config/validators"
	"github.com/pkg/errors"
)

func TestDurationValue(t *testing.T) {
	for _, test := range []struct {
		value          time.Duration
		expectedString string
		input          string
		expectedSet    time.Duration
		expectedError  bool
	}{
		{
			value:          4 * time.Minute,
			expectedString: "4m0s",
			input:          "2m",
			expectedSet:    2 * time.Minute,
			expectedError:  false,
		},
		{
			value:          1 * time.Minute,
			expectedString: "1m0s",
			input:          "not-a-duration",
			expectedError:  true,
		},
	} {
		d := durationValue(test.value)

		if d.String() != test.expectedString {
			t.Errorf("Duration(&%s).String() = %q, expected %q", test.value, d, test.expectedString)
		}

		err := d.Set(test.input)
		if test.expectedError && err == nil {
			t.Errorf("Duration().Set(%q): expected error, got nil instead", test.input)
			continue
		}
		if !test.expectedError && err != nil {
			t.Errorf("Duration().Set(%q): unexpected error: %s", test.input, err)
			continue
		}
		if time.Duration(d) != test.expectedSet {
			t.Errorf("Duration(&d).Set(%q): expected f to be %s, got %s instead", test.input, test.expectedSet, time.Duration(d))
		}
	}
}

func TestDuration(t *testing.T) {
	v := 2 * time.Minute
	flag := "flag"
	env := "ENV"
	usage := "usage"
	f := Duration(&v, flag, env, usage)

	if f.TypeHint == "" {
		t.Error("Duration().TypeHint = \"\": expected .TypeHint to be set")
	}
	if f.Name != flag {
		t.Errorf("Duration(...).Name = %q, expected %q", f.Name, flag)
	}
	if f.Env != env {
		t.Errorf("Duration(...).Env = %q, expected %q", f.Env, env)
	}
	if f.Usage != usage {
		t.Errorf("Duration(...).Usage = %q, expected %q", f.Usage, usage)
	}

	expectedString := "2m0s"
	if f.String() != expectedString {
		t.Errorf("Duration(&2)).String() = %q, expected %q", f.String(), expectedString)
	}

	s := "1m"
	err := f.Set(s)
	if err != nil {
		t.Errorf("Duration().Set(%q): unexpected error: %s", s, err)
	}
	if v != 1*time.Minute {
		t.Errorf("Duration(&v).Set(%q): expected v to be %d, got %d instead", s, 1*time.Minute, v)
	}

	s = "notaduration"
	err = f.Set(s)
	if err == nil {
		t.Errorf("Duration().Set(%q): expected error, got nil", s)
	}

	if f.IsBoolFlag() {
		t.Error("Bool().IsBoolFlag() = true, expected false")
	}
}

func TestDurationValidators(t *testing.T) {
	testValidator := func(shouldFail bool) (validator validators.Duration, called *bool) {
		called = new(bool)
		return func(time.Duration) error {
			*called = true
			if shouldFail {
				return errors.New("failing validator")
			}
			return nil
		}, called
	}

	t.Run("valid input passing validators", func(t *testing.T) {
		var val time.Duration
		v1, v1Called := testValidator(false)
		v2, v2Called := testValidator(false)
		f := Duration(&val, "flag", "ENV", "testing duration validators", v1, v2)
		in := "1m"
		err := f.Set(in)
		if err != nil {
			t.Errorf("Duration(..., v1, v2).Set(%q): unexpected error: %s", in, err)
		}
		if !*v1Called || !*v2Called {
			t.Errorf("Duration(..., v1, v2).Set(%q): some validator wasn't called (v1: %v, v2: %v)", in, *v1Called, *v2Called)
		}
	})

	t.Run("invalid input passing validators", func(t *testing.T) {
		var val time.Duration
		v1, v1Called := testValidator(false)
		f := Duration(&val, "flag", "ENV", "testing duration validators", v1)
		in := "notaduration"
		err := f.Set(in)
		if err == nil {
			t.Errorf("Duration(..., v1).Set(%q): expected error, got nil", in)
		}
		if *v1Called {
			t.Errorf("Duration(..., v1).Set(%q): validator shouldn't have been called", in)
		}
	})

	t.Run("valid input failing validators", func(t *testing.T) {
		var val time.Duration
		v1, v1Called := testValidator(true)
		f := Duration(&val, "flag", "ENV", "testing duration validators", v1)
		in := "2m"
		err := f.Set(in)
		if err == nil {
			t.Errorf("Duration(..., failingV1).Set(%q): expected error, got nil", in)
		}
		if !*v1Called {
			t.Errorf("Duration(..., failingV1).Set(%q): validator should have been called", in)
		}
	})
}

func TestDurationGenerator(t *testing.T) {
	g := DurationGenerator()
	d := g()
	if _, ok := d.(*durationValue); !ok {
		t.Errorf("DurationGenerator(): expected type *durationValue, got %T instead", d)
	}
}
