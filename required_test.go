package rig

import (
	"fmt"
	"testing"
	"time"
)

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

	f = String(&s, "string-flag", "STRING_ENV", "testing Required on String")
	r = Required(Required(f))

	if f.Required {
		t.Errorf("String(...).Required = true, expected false")
	}
	if !r.Required {
		t.Errorf("Required(Required(String(...))).Required = false, expected true")
	}

	f = Var(new(stringValue), "var-flag", "VAR_ENV", "testing Required on Var")
	r = Required(f)

	if f.Required {
		t.Errorf("Var(...).Required = true, expected false")
	}
	if !r.Required {
		t.Errorf("Required(Var(...)).Required = false, expected true")
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
