package config

import (
	"net/url"
	"testing"

	"github.com/Pimmr/config/validators"
	"github.com/pkg/errors"
)

func urlMustParse(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(errors.Wrapf(err, "parsing url %q", s))
	}

	return u
}

func TestURLValue(t *testing.T) {
	for _, test := range []struct {
		value          *url.URL
		expectedString string
		input          string
		expectedError  bool
	}{
		{
			value:          urlMustParse("http://example.com/foo?bar=42"),
			expectedString: "http://example.com/foo?bar=42",
			input:          "https://example.com/baz?fizz=buzz",
			expectedError:  false,
		},
		{
			value:          urlMustParse("http://example.com/foo?bar=42"),
			expectedString: "http://example.com/foo?bar=42",
			input:          ":not-a-valid-url",
			expectedError:  true,
		},
	} {
		v := &url.URL{}
		*v = *test.value
		d := URLValue{URL: &v}

		if d.String() != test.expectedString {
			t.Errorf("URL(&%s).String() = %q, expected %q", test.value, d, test.expectedString)
		}

		err := d.Set(test.input)
		if test.expectedError && err == nil {
			t.Errorf("URL().Set(%q): expected error, got nil instead", test.input)
			continue
		}
		if !test.expectedError && err != nil {
			t.Errorf("URL().Set(%q): unexpected error: %s", test.input, err)
			continue
		}
		if err != nil {
			continue
		}
		if (*d.URL).String() != test.input {
			t.Errorf("URL(&d).Set(%q): expected f to be %s, got %v instead", test.input, test.input, d.URL)
		}
	}
}

func TestURL(t *testing.T) {
	var v *url.URL
	flag := "flag"
	env := "ENV"
	usage := "usage"
	f := URL(&v, flag, env, usage)

	if f.TypeHint == "" {
		t.Error("URL().TypeHint = \"\": expected .TypeHint to be set")
	}
	if f.Name != flag {
		t.Errorf("URL(...).Name = %q, expected %q", f.Name, flag)
	}
	if f.Env != env {
		t.Errorf("URL(...).Env = %q, expected %q", f.Env, env)
	}
	if f.Usage != usage {
		t.Errorf("URL(...).Usage = %q, expected %q", f.Usage, usage)
	}

	expectedString := ""
	if f.String() != expectedString {
		t.Errorf("URL(&2)).String() = %q, expected %q", f.String(), expectedString)
	}

	s := "http://example.com"
	err := f.Set(s)
	if err != nil {
		t.Errorf("URL().Set(%q): unexpected error: %s", s, err)
	}
	if v.String() != s {
		t.Errorf("URL(&v).Set(%q): expected v to be %s, got %v instead", s, s, v)
	}

	s = ":notavalidurl"
	err = f.Set(s)
	if err == nil {
		t.Errorf("URL().Set(%q): expected error, got nil", s)
	}

	if f.IsBoolFlag() {
		t.Error("Bool().IsBoolFlag() = true, expected false")
	}
}

func TestURLValidators(t *testing.T) {
	testValidator := func(shouldFail bool) (validator validators.URL, called *bool) {
		called = new(bool)
		return func(*url.URL) error {
			*called = true
			if shouldFail {
				return errors.New("failing validator")
			}
			return nil
		}, called
	}

	t.Run("valid input passing validators", func(t *testing.T) {
		var val *url.URL
		v1, v1Called := testValidator(false)
		v2, v2Called := testValidator(false)
		f := URL(&val, "flag", "ENV", "testing url validators", v1, v2)
		in := "http://example.org"
		err := f.Set(in)
		if err != nil {
			t.Errorf("URL(..., v1, v2).Set(%q): unexpected error: %s", in, err)
		}
		if !*v1Called || !*v2Called {
			t.Errorf("URL(..., v1, v2).Set(%q): some validator wasn't called (v1: %v, v2: %v)", in, *v1Called, *v2Called)
		}
	})

	t.Run("invalid input passing validators", func(t *testing.T) {
		var val *url.URL
		v1, v1Called := testValidator(false)
		f := URL(&val, "flag", "ENV", "testing url validators", v1)
		in := ":notavalidurl"
		err := f.Set(in)
		if err == nil {
			t.Errorf("URL(..., v1).Set(%q): expected error, got nil", in)
		}
		if *v1Called {
			t.Errorf("URL(..., v1).Set(%q): validator shouldn't have been called", in)
		}
	})

	t.Run("valid input failing validators", func(t *testing.T) {
		var val *url.URL
		v1, v1Called := testValidator(true)
		f := URL(&val, "flag", "ENV", "testing url validators", v1)
		in := "http://example.org"
		err := f.Set(in)
		if err == nil {
			t.Errorf("URL(..., failingV1).Set(%q): expected error, got nil", in)
		}
		if !*v1Called {
			t.Errorf("URL(..., failingV1).Set(%q): validator should have been called", in)
		}
	})
}

func TestURLGenerator(t *testing.T) {
	g := URLGenerator()
	d := g()
	if _, ok := d.(*URLValue); !ok {
		t.Errorf("URLGenerator(): expected type *URLValue, got %T instead", d)
	}
}
