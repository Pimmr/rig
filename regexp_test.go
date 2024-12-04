package rig

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/Pimmr/rig/validators"
)

func TestRegexpValue(t *testing.T) {
	for _, test := range []struct {
		value          *regexp.Regexp
		expectedString string
		input          string
		expectedError  bool
	}{
		{
			value:          regexp.MustCompile("[a-d][0-9]+"),
			expectedString: "[a-d][0-9]+",
			input:          "[x-z][2-5]{2}",
			expectedError:  false,
		},
		{
			value:          regexp.MustCompile("[a-d][0-9]+"),
			expectedString: "[a-d][0-9]+",
			input:          "[not-a-valid-regexp",
			expectedError:  true,
		},
	} {
		v := test.value.Copy()
		d := regexpValue{Regexp: &v}

		if d.String() != test.expectedString {
			t.Errorf("Regexp(&%s).String() = %q, expected %q", test.value, d, test.expectedString)
		}

		err := d.Set(test.input)
		if test.expectedError && err == nil {
			t.Errorf("Regexp().Set(%q): expected error, got nil instead", test.input)
			continue
		}
		if !test.expectedError && err != nil {
			t.Errorf("Regexp().Set(%q): unexpected error: %s", test.input, err)
			continue
		}
		if err != nil {
			continue
		}
		if (*d.Regexp).String() != test.input {
			t.Errorf("Regexp(&d).Set(%q): expected f to be %s, got %v instead", test.input, test.input, d.Regexp)
		}
	}
}

func TestRegexp(t *testing.T) {
	var v *regexp.Regexp
	flag := "flag"
	env := "ENV"
	usage := "usage"
	f := Regexp(&v, flag, env, usage)

	if f.TypeHint == "" {
		t.Error("Regexp().TypeHint = \"\": expected .TypeHint to be set")
	}
	if f.Name != flag {
		t.Errorf("Regexp(...).Name = %q, expected %q", f.Name, flag)
	}
	if f.Env != env {
		t.Errorf("Regexp(...).Env = %q, expected %q", f.Env, env)
	}
	if f.Usage != usage {
		t.Errorf("Regexp(...).Usage = %q, expected %q", f.Usage, usage)
	}

	expectedString := ""
	if f.String() != expectedString {
		t.Errorf("Regexp(&2)).String() = %q, expected %q", f.String(), expectedString)
	}

	s := "[x-z][2-5]{2}"
	err := f.Set(s)
	if err != nil {
		t.Errorf("Regexp().Set(%q): unexpected error: %s", s, err)
	}
	if v.String() != s {
		t.Errorf("Regexp(&v).Set(%q): expected v to be %s, got %v instead", s, s, v)
	}

	s = "[notavalidregexp"
	err = f.Set(s)
	if err == nil {
		t.Errorf("Regexp().Set(%q): expected error, got nil", s)
	}

	if f.IsBoolFlag() {
		t.Error("Regexp().IsBoolFlag() = true, expected false")
	}
}

func TestRegexpValidators(t *testing.T) {
	testValidator := func(shouldFail bool) (validator validators.Regexp, called *bool) {
		called = new(bool)
		return func(*regexp.Regexp) error {
			*called = true
			if shouldFail {
				return errors.New("failing validator")
			}
			return nil
		}, called
	}

	t.Run("valid input passing validators", func(t *testing.T) {
		var val *regexp.Regexp
		v1, v1Called := testValidator(false)
		v2, v2Called := testValidator(false)
		f := Regexp(&val, "flag", "ENV", "testing regexp validators", v1, v2)
		in := "[x-z][2-5]{2}"
		err := f.Set(in)
		if err != nil {
			t.Errorf("Regexp(..., v1, v2).Set(%q): unexpected error: %s", in, err)
		}
		if !*v1Called || !*v2Called {
			t.Errorf("Regexp(..., v1, v2).Set(%q): some validator wasn't called (v1: %v, v2: %v)", in, *v1Called, *v2Called)
		}
	})

	t.Run("invalid input passing validators", func(t *testing.T) {
		var val *regexp.Regexp
		v1, v1Called := testValidator(false)
		f := Regexp(&val, "flag", "ENV", "testing regexp validators", v1)
		in := "[notavalidregexp"
		err := f.Set(in)
		if err == nil {
			t.Errorf("Regexp(..., v1).Set(%q): expected error, got nil", in)
		}
		if *v1Called {
			t.Errorf("Regexp(..., v1).Set(%q): validator shouldn't have been called", in)
		}
	})

	t.Run("valid input failing validators", func(t *testing.T) {
		var val *regexp.Regexp
		v1, v1Called := testValidator(true)
		f := Regexp(&val, "flag", "ENV", "testing regexp validators", v1)
		in := "[a-d][0-9]+"
		err := f.Set(in)
		if err == nil {
			t.Errorf("Regexp(..., failingV1).Set(%q): expected error, got nil", in)
		}
		if !*v1Called {
			t.Errorf("Regexp(..., failingV1).Set(%q): validator should have been called", in)
		}
	})
}

func TestRegexpGenerator(t *testing.T) {
	g := RegexpGenerator()
	d := g()
	r, ok := d.(*regexpValue)
	if !ok {
		t.Errorf("RegexpGenerator(): expected type *RegexpValue, got %T instead", d)
	}

	in := "[a-d]+"
	err := d.Set(in)
	if err != nil {
		t.Errorf("RegexpGenerator()().Set(%q): unexpected error: %s", in, err)
		t.FailNow()
	}
	if (*r.Regexp).String() != in {
		t.Errorf("RegexpGenerator()().Set(%q) = %q, expected %q", in, (*r.Regexp).String(), in)
	}
}

func ExampleRegexpGenerator() {
	var rr []*regexp.Regexp

	c := &Config{
		FlagSet: testingFlagset(),
		Flags: []*Flag{
			Repeatable(&rr, RegexpGenerator(), "regexp", "REGEXP", "Repeatable regexp flag"),
		},
	}

	err := c.Parse([]string{"-regexp=^foo.*", "-regexp=[0-9]bar{1,5}"})
	if err != nil {
		return
	}

	fmt.Printf("%v\n", rr)

	// Output: [^foo.* [0-9]bar{1 5}]
}
