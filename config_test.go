package config

import (
	"bytes"
	"flag"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/Pimmr/config/validators"

	"github.com/pkg/errors"
)

func TestFlagMissingError(t *testing.T) {
	for _, test := range []struct {
		flag               Flag
		errorShouldContain []string
	}{
		{
			flag:               Flag{Name: "", Env: ""},
			errorShouldContain: []string{},
		},
		{
			flag:               Flag{Name: "foo", Env: ""},
			errorShouldContain: []string{"-foo"},
		},
		{
			flag:               Flag{Name: "", Env: "BAR"},
			errorShouldContain: []string{"BAR"},
		},
		{
			flag:               Flag{Name: "foo", Env: "BAR"},
			errorShouldContain: []string{"-foo", "BAR"},
		},
	} {
		err := test.flag.missingError()
		if err == nil {
			t.Errorf("Flag(%+v).missingError(): expected error, got nil", test.flag)
			continue
		}

		errStr := err.Error()
		for _, s := range test.errorShouldContain {
			if !strings.Contains(errStr, s) {
				t.Errorf("Flag(%+v).missingError() = %q: expected to find %q in error string.", test.flag, errStr, s)
			}
		}
	}
}

func TestConfigUsage(t *testing.T) {
	const (
		stringFlag    = "string-flag"
		stringDefault = "string default"
		intEnv        = "INT_ENV"
		intUsage      = "int usage"
		intDefault    = "32"
		boolFlag      = "bool-flag"
		boolEnv       = "BOOL_ENV"
		boolUsage     = "bool usage"
		boolRequired  = "required"
	)

	var (
		s = stringDefault
		i = 32
		b bool
		f float64
	)

	c := &Config{
		FlagSet: flag.NewFlagSet("flagset", flag.ContinueOnError),
		Flags: []*Flag{
			String(&s, stringFlag, "", ""),
			Int(&i, "", intEnv, intUsage),
			Required(Bool(&b, boolFlag, boolEnv, boolUsage)),
			Float64(&f, "", "", "no flag or env set for this one"),
		},
	}
	buf := &bytes.Buffer{}
	c.FlagSet.SetOutput(buf)

	c.Usage()

	if buf.Len() == 0 {
		t.Errorf("Config.Usage(): expected usage to be written to the flagset's output")
	}

	expected := []string{
		stringFlag, stringDefault,
		intEnv, intUsage, intDefault,
		boolFlag, boolEnv, boolUsage, boolRequired,
	}

	bufStr := buf.String()
	for _, s := range expected {
		if !strings.Contains(bufStr, s) {
			t.Errorf("c.Usage() output: expected to find %q", s)
		}
	}
}

func TestRequired(t *testing.T) {
	var s string

	f := String(&s, "string-flag", "STRING_ENV", "testing Required on String")
	r := Required(f)

	if f.Required {
		t.Errorf("String(...).Required = true, expected false")
	}
	if !r.Required {
		t.Errorf("Required(String(...)).Required = false, expected true")
	}

	expectedSuffix := ", required"
	if !strings.HasSuffix(r.TypeHint, expectedSuffix) {
		t.Errorf("Required(String(...)).TypeHint = %q, expected %q suffix", r.TypeHint, expectedSuffix)
	}

	f = Var(new(stringValue), "var-flag", "VAR_ENV", "testing Required on Var")
	r = Required(f)

	if f.Required {
		t.Errorf("Var(...).Required = true, expected false")
	}
	if !r.Required {
		t.Errorf("Required(Var(...)).Required = false, expected true")
	}

	expected := "required"
	if r.TypeHint != expected {
		t.Errorf("Required(Var(...)).TypeHint = %q, expected %q", r.TypeHint, expected)
	}
}

func TestRepeatable(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		var ii []int

		f := Repeatable(&ii, IntGenerator(), "flag", "ENV", "testing Repeatable integers")
		in1 := "42"
		err := f.Set(in1)
		if err != nil {
			t.Errorf("Repeatable(&[]int).Set(%q): unexpected error: %s", in1, err)
		}
		in2 := "23"
		err = f.Set(in2)
		if err != nil {
			t.Errorf("Repeatable(&[]int).Set(%q): unexpected error: %s", in2, err)
		}

		expected := []int{42, 23}
		if !reflect.DeepEqual(ii, expected) {
			t.Errorf("Repeatable(&[]int).Set(...) = %q, expected %q", ii, expected)
		}
	})

	testValidator := func(shouldFail bool) validators.Var {
		return func(flag.Value) error {
			if shouldFail {
				return errors.New("failing validator")
			}

			return nil
		}
	}

	t.Run("int passing validator", func(t *testing.T) {
		var ii []int

		f := Repeatable(&ii, IntGenerator(), "flag", "ENV", "testing Repeatable integers", testValidator(false))
		in1 := "42"
		err := f.Set(in1)
		if err != nil {
			t.Errorf("Repeatable(&[]int).Set(%q, passingValidator): unexpected error: %s", in1, err)
		}

		expected := []int{42}
		if !reflect.DeepEqual(ii, expected) {
			t.Errorf("Repeatable(&[]int).Set(..., passingValidator) = %q, expected %q", ii, expected)
		}
	})

	t.Run("int failing validator", func(t *testing.T) {
		var ii []int

		f := Repeatable(&ii, IntGenerator(), "flag", "ENV", "testing Repeatable integers", testValidator(true))
		in1 := "42"
		err := f.Set(in1)
		if err == nil {
			t.Errorf("Repeatable(&[]int).Set(%q, failingValidator): expected error, got nil", in1)
		}
	})

	t.Run("int multiple values", func(t *testing.T) {
		var ii []int

		f := Repeatable(&ii, IntGenerator(), "flag", "ENV", "testing Repeatable integers")
		in := "42,23"
		err := f.Set(in)
		if err != nil {
			t.Errorf("Repeatable(&[]int).Set(%q): unexpected error: %s", in, err)
		}

		expected := []int{42, 23}
		if !reflect.DeepEqual(ii, expected) {
			t.Errorf("Repeatable(&[]int).Set(...) = %q, expected %q", ii, expected)
		}
	})

	t.Run("string multiple values with escaping", func(t *testing.T) {
		var ss []string

		f := Repeatable(&ss, StringGenerator(), "flag", "ENV", "testing Repeatable strings")
		in := "foo\\,bar,baz"
		err := f.Set(in)
		if err != nil {
			t.Errorf("Repeatable(&[]string).Set(%q): unexpected error: %s", in, err)
		}

		expected := []string{"foo,bar", "baz"}
		if !reflect.DeepEqual(ss, expected) {
			t.Errorf("Repeatable(&[]string).Set(...) = %q, expected %q", ss, expected)
		}
	})

	t.Run("URL", func(t *testing.T) {
		var uu []URLValue

		f := Repeatable(&uu, URLGenerator(), "flag", "ENV", "testing Repeatable URLs")
		in1 := "http://example.com/foo"
		err := f.Set(in1)
		if err != nil {
			t.Errorf("Repeatable(&[]URLValue).Set(%q): unexpected error: %s", in1, err)
		}
		in2 := "http://example.com/bar"
		err = f.Set(in2)
		if err != nil {
			t.Errorf("Repeatable(&[]URLValue).Set(%q): unexpected error: %s", in2, err)
		}

		in1URL := urlMustParse(in1)
		in2URL := urlMustParse(in2)
		expected := []URLValue{
			URLValue{
				URL: &in1URL,
			},
			URLValue{
				URL: &in2URL,
			},
		}
		if !reflect.DeepEqual(uu, expected) {
			t.Errorf("Repeatable(&[]URLValue).Set(...) = %q, expected %q", uu, expected)
		}
	})

	t.Run("Regexp", func(t *testing.T) {
		var rr []RegexpValue

		f := Repeatable(&rr, RegexpGenerator(), "flag", "ENV", "testing Repeatable Regexps")
		in1 := "[a-d][0-7]+"
		err := f.Set(in1)
		if err != nil {
			t.Errorf("Repeatable(&[]RegexpValue).Set(%q): unexpected error: %s", in1, err)
		}
		in2 := "[e-g][8-9]{2}"
		err = f.Set(in2)
		if err != nil {
			t.Errorf("Repeatable(&[]RegexpValue).Set(%q): unexpected error: %s", in2, err)
		}

		in1Regexp := regexp.MustCompile(in1)
		in2Regexp := regexp.MustCompile(in2)
		expected := []RegexpValue{
			RegexpValue{
				Regexp: &in1Regexp,
			},
			RegexpValue{
				Regexp: &in2Regexp,
			},
		}
		if !reflect.DeepEqual(rr, expected) {
			t.Errorf("Repeatable(&[]RegexpValue).Set(...) = %q, expected %q", rr, expected)
		}
	})

	t.Run("not a pointer", func(t *testing.T) {
		var ss []string

		f := Repeatable(ss, RegexpGenerator(), "flag", "ENV", "testing Repeatable Regexps")
		in := "foo"
		err := f.Set(in)
		if err == nil {
			t.Errorf("Repeatable([]string).Set(%q): expected error, got nil", in)
		}
	})

	t.Run("not a slice", func(t *testing.T) {
		var s string

		f := Repeatable(&s, RegexpGenerator(), "flag", "ENV", "testing Repeatable Regexps")
		in := "foo"
		err := f.Set(in)
		if err == nil {
			t.Errorf("Repeatable(&string).Set(%q): expected error, got nil", in)
		}
	})

	t.Run("invalid value", func(t *testing.T) {
		var ii []int

		f := Repeatable(&ii, IntGenerator(), "flag", "ENV", "testing Repeatable integers")
		in1 := "foo"
		err := f.Set(in1)
		if err == nil {
			t.Errorf("Repeatable(&[]int).Set(%q): expected error, got nil", in1)
		}
	})

	t.Run("generator incompatible with values", func(t *testing.T) {
		var ii []int

		f := Repeatable(&ii, StringGenerator(), "flag", "ENV", "testing Repeatable integers")
		in1 := "42"
		err := f.Set(in1)
		if err == nil {
			t.Errorf("Repeatable(&[]int, StringGenerator()).Set(%q): expected error, got nil", in1)
		}
	})
}

type nopValue struct{}

func (v nopValue) String() string {
	return ""
}

func (v nopValue) Set(string) error {
	return nil
}

func TestMakeGenerator(t *testing.T) {
	g := MakeGenerator(new(stringValue))
	v := g()

	if _, ok := v.(*stringValue); !ok {
		t.Errorf("MakeGenerator(new(stringValue))() = %T, expected *stringValue", v)
	}

	g = MakeGenerator(nopValue{})
	v = g()
	if _, ok := v.(*nopValue); !ok {
		t.Errorf("MakeGenerator(nopValue{})() = %T, expected *nopValue", v)
	}
}
