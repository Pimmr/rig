package rig

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Pimmr/rig/validators"
)

type stringeringString string

func (s stringeringString) String() string {
	return "s:" + string(s)
}

type stringeringStringSlice []string

func (ss stringeringStringSlice) String() string {
	return strings.Join([]string(ss), ",")
}

func pointerToSlice(elems ...string) *[]string {
	ss := make([]string, 0, len(elems))
	ss = append(ss, elems...)

	return &ss
}

func TestSliceValueString(t *testing.T) {
	for _, test := range []struct {
		value    interface{}
		expected string
	}{
		{
			value:    []string{"foo", "bar"},
			expected: "[foo,bar]",
		},
		{
			value:    []stringeringString{"foo", "bar"},
			expected: "[s:foo,s:bar]",
		},
		{
			value:    pointerToSlice("foo", "bar"),
			expected: "[foo,bar]",
		},
		{
			value:    stringeringStringSlice{"foo", "bar"},
			expected: "foo,bar",
		},
	} {
		got := sliceValue{
			value: reflect.ValueOf(test.value),
		}.String()

		if got != test.expected {
			t.Errorf("sliceValue{%#v}.String() = %q, expected %q", test.value, got, test.expected)
		}
	}
}

type testUnexported struct {
	unexported []string
}

type testExported struct {
	Exported []string
}

func TestSliceValueSet(t *testing.T) {
	t.Run("using value from exported field", func(t *testing.T) {
		p := &testExported{}
		v := reflect.Indirect(reflect.ValueOf(p))
		pf := v.FieldByName("Exported")

		sv := sliceValue{
			value:     pf.Addr(),
			generator: StringGenerator(),
		}
		err := sv.set("foo")
		if err != nil {
			t.Errorf("sliceValue{valueFromExportedField}.set(\"foo\"): unexpected error: %s", err)
		}
	})

	t.Run("using value from unexported field", func(t *testing.T) {
		p := &testUnexported{}
		v := reflect.Indirect(reflect.ValueOf(p))
		pf := v.FieldByName("unexported")

		sv := sliceValue{
			value:     pf.Addr(),
			generator: StringGenerator(),
		}
		err := sv.set("foo")
		if err == nil {
			t.Logf("sliceValue.value was obtained by accessing an unexported struct field")
			t.Errorf("sliceValue{valueFromUnexportedField}.set(\"foo\"): expected error, got nil")
		}
	})
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

	testValidator := func(shouldFail bool) validators.Repeatable {
		return func(interface{}) error {
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
		var uu []*url.URL

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
		expected := []*url.URL{in1URL, in2URL}
		if !reflect.DeepEqual(uu, expected) {
			t.Errorf("Repeatable(&[]URLValue).Set(...) = %q, expected %q", uu, expected)
		}
	})

	t.Run("Regexp", func(t *testing.T) {
		var rr []*regexp.Regexp

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
		expected := []*regexp.Regexp{in1Regexp, in2Regexp}
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
	if _, ok := v.(nopValue); !ok {
		t.Errorf("MakeGenerator(nopValue{})() = %T, expected nopValue", v)
	}
}

func ExampleRepeatable() {
	var bb []bool
	var ss []string
	var dd []time.Duration

	c := &Config{
		FlagSet: testingFlagset(),
		Flags: []*Flag{
			Repeatable(&bb, BoolGenerator(), "bool", "BOOL", "repeatable boolean flag"),
			Repeatable(&ss, StringGenerator(), "string", "STRING", "repeatable string flag"),
			Repeatable(&dd, DurationGenerator(), "duration", "DURATION", "repeatable duration flag"),
		},
	}

	err := c.Parse([]string{"-bool=t,f,t", "-string=foo", "-string=bar", "-duration=5m2s,3m44s"})
	if err != nil {
		return
	}

	fmt.Printf("booleans: %v\nstrings: %q\ndurations: %v\n", bb, ss, dd)

	// Output:
	// booleans: [true false true]
	// strings: ["foo" "bar"]
	// durations: [5m2s 3m44s]
}

type CustomType string

func (c CustomType) String() string {
	return string(c)
}

func (c *CustomType) Set(s string) error {
	*c = CustomType(s)

	return nil
}

func ExampleMakeGenerator() {
	var cc []CustomType // implements the "flag".Value interface

	c := &Config{
		FlagSet: testingFlagset(),
		Flags: []*Flag{
			Repeatable(&cc, MakeGenerator(new(CustomType)), "custom", "CUSTOM", "Repeatable flag with a custom type"),
		},
	}

	err := c.Parse([]string{"-custom=foo,bar"})
	if err != nil {
		return
	}

	fmt.Printf("%v\n", cc)

	// Output: [foo bar]
}
