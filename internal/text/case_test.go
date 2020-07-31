package text

import "testing"

func TestToSnakeCase(t *testing.T) {
	sep := "-"

	for _, test := range []struct {
		In       string
		Expected string
	}{
		{In: "", Expected: ""},
		{In: "f", Expected: "f"},
		{In: "foo", Expected: "foo"},
		{In: "Foo", Expected: "foo"},
		{In: "fooBar", Expected: "foo-bar"},
		{In: "FooBar", Expected: "foo-bar"},
		{In: "fooBarBaz", Expected: "foo-bar-baz"},
		{In: "fooBARBaz", Expected: "foo-bar-baz"},
	} {
		got := ToSnakeCase(test.In, sep)
		if got != test.Expected {
			t.Errorf("ToSnakeCase(%q, %q) = %q, expected %q", test.In, sep, got, test.Expected)
		}
	}
}

func TestToUpperSnakeCase(t *testing.T) {
	sep := "_"

	for _, test := range []struct {
		In       string
		Expected string
	}{
		{In: "", Expected: ""},
		{In: "F", Expected: "F"},
		{In: "foo", Expected: "FOO"},
		{In: "Foo", Expected: "FOO"},
		{In: "fooBar", Expected: "FOO_BAR"},
		{In: "FooBar", Expected: "FOO_BAR"},
		{In: "fooBarBaz", Expected: "FOO_BAR_BAZ"},
		{In: "fooBARBaz", Expected: "FOO_BAR_BAZ"},
	} {
		got := ToUpperSnakeCase(test.In, sep)
		if got != test.Expected {
			t.Errorf("ToUpperSnakeCase(%q, %q) = %q, expected %q", test.In, sep, got, test.Expected)
		}
	}
}
