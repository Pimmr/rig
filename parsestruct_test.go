package rig

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestApplyTypeHint(t *testing.T) {
	for _, test := range []struct {
		flag     Flag
		typeHint string
		expected *Flag
	}{
		{
			flag: Flag{
				Name: "foo",
				Env:  "bar",
			},
			typeHint: "",
			expected: &Flag{
				Name: "foo",
				Env:  "bar",
			},
		},
		{
			flag: Flag{
				Name: "foo",
				Env:  "bar",
			},
			typeHint: "foobar",
			expected: &Flag{
				Name:     "foo",
				Env:      "bar",
				TypeHint: "foobar",
			},
		},
	} {
		flag := test.flag
		got := applyTypeHint(&flag, test.typeHint)
		if !reflect.DeepEqual(got, test.expected) {
			t.Errorf("applyTypeHint(%+v, %q) = %+v, expected %+v", test.flag, test.typeHint, got, test.expected)
		}
	}
}

func TestApplyRequired(t *testing.T) {
	for _, test := range []struct {
		flag     Flag
		required bool
		expected *Flag
	}{
		{
			flag: Flag{
				Name: "foo",
				Env:  "bar",
			},
			required: false,
			expected: &Flag{
				Name: "foo",
				Env:  "bar",
			},
		},
		{
			flag: Flag{
				Name: "foo",
				Env:  "bar",
			},
			required: true,
			expected: &Flag{
				Name:     "foo",
				Env:      "bar",
				Required: true,
			},
		},
		{
			flag: Flag{
				Name:     "foo",
				Env:      "bar",
				Required: true,
			},
			required: true,
			expected: &Flag{
				Name:     "foo",
				Env:      "bar",
				Required: true,
			},
		},
		{
			flag: Flag{
				Name:     "foo",
				Env:      "bar",
				Required: true,
			},
			required: false,
			expected: &Flag{
				Name:     "foo",
				Env:      "bar",
				Required: true,
			},
		},
	} {
		flag := test.flag
		got := applyRequired(&flag, test.required)
		got.TypeHint = "" // We don't want to test this field
		if !reflect.DeepEqual(got, test.expected) {
			t.Errorf("applyRequired(%#v, %v) = %#v, expected %#v", test.flag, test.required, got, test.expected)
		}
	}
}

func TestPrefix(t *testing.T) {
	setupFlags := func(t *testing.T) []*Flag {
		t.Helper()
		v := struct {
			FlagA string `flag:"flag-a" env:"FLAG_A"`
			FlagB string `flag:"flag-b"`
			FlagC string `env:"FLAG_C"`
			FlagD string
		}{}

		flags, err := StructToFlags(&v)
		if err != nil {
			t.Fatalf("StructToFlags(): unexpected error while setting test up: %v", err)
		}
		return flags
	}

	t.Run("required=false", func(t *testing.T) {
		flags := setupFlags(t)

		flagPrefix := "prefix"
		envPrefix := "PREFIX"
		flags = prefix(flags, flagPrefix, envPrefix, false)
		for _, f := range flags {
			if f.Name != "" && !strings.HasPrefix(f.Name, "prefix-") {
				t.Errorf("prefix(flags, %q, %q, false): expected flag name %q to have '%s-' prefix", flagPrefix, envPrefix, f.Name, flagPrefix)
			}
			if f.Env != "" && !strings.HasPrefix(f.Env, "PREFIX_") {
				t.Errorf("prefix(flags, %q, %q, false): expected flag env %q to have '%s' prefix", flagPrefix, envPrefix, f.Env, envPrefix)
			}
			if f.Required {
				t.Errorf("prefix(flags, %q, %q, false): expected .Required to be false", flagPrefix, envPrefix)
			}
		}
	})

	t.Run("required=true", func(t *testing.T) {
		flags := setupFlags(t)

		flagPrefix := "prefix"
		envPrefix := "PREFIX"
		flags = prefix(flags, flagPrefix, envPrefix, true)
		for _, f := range flags {
			if f.Name != "" && !strings.HasPrefix(f.Name, "prefix-") {
				t.Errorf("prefix(flags, %q, %q, true): expected flag name %q to have '%s-' prefix", flagPrefix, envPrefix, f.Name, flagPrefix)
			}
			if f.Env != "" && !strings.HasPrefix(f.Env, "PREFIX_") {
				t.Errorf("prefix(flags, %q, %q, true): expected flag env %q to have '%s' prefix", flagPrefix, envPrefix, f.Env, envPrefix)
			}
			if !f.Required {
				t.Errorf("prefix(flags, %q, %q, true): expected .Required to be true", flagPrefix, envPrefix)
			}
		}
	})
}

type testFlagValue struct{}

func (v testFlagValue) String() string {
	return "testFlagValue"
}

func (v testFlagValue) Set(string) error {
	return nil
}

func TestFlagFromInterface(t *testing.T) {
	flagName := "test-flag-name"
	envName := "TEST_ENV_NAME"
	usage := "test usage"

	t.Run("non-slice", func(t *testing.T) {
		for _, test := range []struct {
			test        string
			in          interface{}
			expected    *Flag
			expectError bool
		}{
			{
				test:     "int",
				in:       new(int),
				expected: Int(new(int), flagName, envName, usage),
			},
			{
				test:     "int64",
				in:       new(int64),
				expected: Int64(new(int64), flagName, envName, usage),
			},
			{
				test:     "uint",
				in:       new(uint),
				expected: Uint(new(uint), flagName, envName, usage),
			},
			{
				test:     "uint64",
				in:       new(uint64),
				expected: Uint64(new(uint64), flagName, envName, usage),
			},
			{
				test:     "string",
				in:       new(string),
				expected: String(new(string), flagName, envName, usage),
			},
			{
				test:     "bool",
				in:       new(bool),
				expected: Bool(new(bool), flagName, envName, usage),
			},
			{
				test:     "duration",
				in:       new(time.Duration),
				expected: Duration(new(time.Duration), flagName, envName, usage),
			},
			{
				test:     "float64",
				in:       new(float64),
				expected: Float64(new(float64), flagName, envName, usage),
			},
			{
				test:     "regexp",
				in:       new(*regexp.Regexp),
				expected: Regexp(new(*regexp.Regexp), flagName, envName, usage),
			},
			{
				test:     "url",
				in:       new(*url.URL),
				expected: URL(new(*url.URL), flagName, envName, usage),
			},
			{
				test:     "flag.Value",
				in:       testFlagValue{},
				expected: Var(testFlagValue{}, flagName, envName, usage),
			},
			{
				test:        "unsupported type",
				in:          &struct{}{},
				expectError: true,
			},
		} {
			f, err := flagFromInterface(test.in, flagName, envName, usage)
			if test.expectError && err == nil {
				t.Errorf("flagFromInterface(%s): expected error, got nil", test.test)
				continue
			}
			if !test.expectError && err != nil {
				t.Errorf("flagFromInterface(%s): unexpected error: %v", test.test, err)
				continue
			}

			if !reflect.DeepEqual(f, test.expected) {
				t.Errorf("flagFromInterface(%s) = %#v, expected %#v", test.test, f, test.expected)
			}
		}
	})

	t.Run("slice", func(t *testing.T) {
		for _, test := range []struct {
			test     string
			in       interface{}
			expected *Flag
		}{
			{
				test:     "[]int",
				in:       new([]int),
				expected: Repeatable(new([]int), IntGenerator(), flagName, envName, usage),
			},
			{
				test:     "[]int64",
				in:       new([]int64),
				expected: Repeatable(new([]int64), Int64Generator(), flagName, envName, usage),
			},
			{
				test:     "[]uint",
				in:       new([]uint),
				expected: Repeatable(new([]uint), UintGenerator(), flagName, envName, usage),
			},
			{
				test:     "[]uint64",
				in:       new([]uint64),
				expected: Repeatable(new([]uint64), Uint64Generator(), flagName, envName, usage),
			},
			{
				test:     "[]string",
				in:       new([]string),
				expected: Repeatable(new([]string), StringGenerator(), flagName, envName, usage),
			},
			{
				test:     "[]bool",
				in:       new([]bool),
				expected: Repeatable(new([]bool), BoolGenerator(), flagName, envName, usage),
			},
			{
				test:     "[]duration",
				in:       new([]time.Duration),
				expected: Repeatable(new([]time.Duration), DurationGenerator(), flagName, envName, usage),
			},
			{
				test:     "[]float64",
				in:       new([]float64),
				expected: Repeatable(new([]float64), Float64Generator(), flagName, envName, usage),
			},
			{
				test:     "[]regexp",
				in:       new([]*regexp.Regexp),
				expected: Repeatable(new([]*regexp.Regexp), RegexpGenerator(), flagName, envName, usage),
			},
			{
				test:     "[]url",
				in:       new([]*url.URL),
				expected: Repeatable(new([]*url.URL), URLGenerator(), flagName, envName, usage),
			},
		} {
			f, err := flagFromInterface(test.in, flagName, envName, usage)
			if err != nil {
				t.Errorf("flagFromInterface(%s): unexpected error: %v", test.test, err)
				continue
			}

			testSv, ok := test.expected.Value.(sliceValue)
			if !ok {
				t.Errorf("expected test.expected.Value to be of type sliceValue, got %T instead", test.expected.Value)
				continue
			}

			if sv, ok := f.Value.(sliceValue); !ok {
				t.Errorf("flag.Value: expected type 'sliceValue', got %T instead", f.Value)
			} else if reflect.TypeOf(sv.value) != reflect.TypeOf(testSv.value) {
				t.Errorf("flagFromInterface(%s).Value.value = %T, expected %T", test.test, sv.value, testSv.value)
				continue
			}

			f.Value = nil
			test.expected.Value = nil
			if !reflect.DeepEqual(f, test.expected) {
				t.Errorf("flagFromInterface(%s) = %#v, expected %#v", test.test, f, test.expected)
			}
		}
	})
}

func TestGetFieldInfo(t *testing.T) {
	t.Run("valid ignore", func(t *testing.T) {
		type validIgnore struct {
			FlagA int `flag:"-"`
		}
		v := &validIgnore{}
		val := reflect.Indirect(reflect.ValueOf(v))
		typ := val.Type()
		fieldTyp, _ := typ.FieldByName("FlagA")
		field := val.FieldByName("FlagA")

		fi, err := getFieldInfo(field, fieldTyp)
		if err != nil {
			t.Errorf("getFieldInfo(%T): unexpected error: %v", v, err)
			return
		}
		if fi != nil {
			t.Errorf("getFieldInfo(%T) = %+v, expected nil", v, fi)
		}
	})

	t.Run("invalid flag option", func(t *testing.T) {
		type invalidIgnore struct {
			FlagA int `flag:"flag-a,foobar"`
		}
		v := &invalidIgnore{}
		val := reflect.Indirect(reflect.ValueOf(v))
		typ := val.Type()
		fieldTyp, _ := typ.FieldByName("FlagA")
		field := val.FieldByName("FlagA")

		_, err := getFieldInfo(field, fieldTyp)
		if err == nil {
			t.Errorf("getFieldInfo(%T): expected error, got nil", v)
		}
	})

	t.Run("valid required", func(t *testing.T) {
		type validRequired struct {
			FlagA int `flag:"flag-a,require"`
		}
		v := &validRequired{}
		val := reflect.Indirect(reflect.ValueOf(v))
		typ := val.Type()
		fieldTyp, _ := typ.FieldByName("FlagA")
		field := val.FieldByName("FlagA")

		fi, err := getFieldInfo(field, fieldTyp)
		if err != nil {
			t.Errorf("getFieldInfo(%T): unexpected error: %v", v, err)
			return
		}
		if fi == nil {
			t.Errorf("getFieldInfo(%T) = nil, expected value", v)
		}

		if !fi.required {
			t.Errorf("getFieldInfo(%T).required = false, expected true", v)
		}
	})

	t.Run("valid default", func(t *testing.T) {
		type validDefault struct {
			FlagA int
		}
		v := &validDefault{}
		val := reflect.Indirect(reflect.ValueOf(v))
		typ := val.Type()
		fieldTyp, _ := typ.FieldByName("FlagA")
		field := val.FieldByName("FlagA")

		fi, err := getFieldInfo(field, fieldTyp)
		if err != nil {
			t.Errorf("getFieldInfo(%T): unexpected error: %v", v, err)
			return
		}
		if fi == nil {
			t.Errorf("getFieldInfo(%T) = nil, expected value", v)
			return
		}

		if fi.flag != "flag-a" {
			t.Errorf("getFieldInfo(%T).flag = %q, expected %q", v, fi.flag, "flag-a")
		}
		if fi.required {
			t.Errorf("getFieldInfo(%T).required = true, expected false", v)
		}
	})

	t.Run("invalid inline", func(t *testing.T) {
		type invalidInline struct {
			FlagA int `flag:",inline" env:",inline"`
		}
		v := &invalidInline{}
		val := reflect.Indirect(reflect.ValueOf(v))
		typ := val.Type()
		fieldTyp, _ := typ.FieldByName("FlagA")
		field := val.FieldByName("FlagA")

		fi, err := getFieldInfo(field, fieldTyp)
		if err != nil {
			t.Errorf("getFieldInfo(%T): unexpected error: %v", v, err)
			return
		}
		if fi != nil {
			t.Errorf("getFieldInfo(%T) = %v, expected nil", v, fi)
		}
	})

	t.Run("invalid env option", func(t *testing.T) {
		type invalidRequired struct {
			FlagA int `env:"FLAG_A,foobar"`
		}
		v := &invalidRequired{}
		val := reflect.Indirect(reflect.ValueOf(v))
		typ := val.Type()
		fieldTyp, _ := typ.FieldByName("FlagA")
		field := val.FieldByName("FlagA")

		_, err := getFieldInfo(field, fieldTyp)
		if err == nil {
			t.Errorf("getFieldInfo(%T): expected error, got nil", v)
		}
	})

	t.Run("non-addressable field", func(t *testing.T) {
		type nonAddressableField struct {
			FlagA int `flag:"flag-a,require"`
		}
		v := nonAddressableField{}
		val := reflect.Indirect(reflect.ValueOf(v))
		typ := val.Type()
		fieldTyp, _ := typ.FieldByName("FlagA")
		field := val.FieldByName("FlagA")

		_, err := getFieldInfo(field, fieldTyp)
		if err == nil {
			t.Errorf("getFieldInfo(%T): expected error, got nil", v)
		}
	})
}

func TestGetFlagName(t *testing.T) {
	for _, test := range []struct {
		Field string
		Tag   string

		// Expected:
		FlagName string
		Required bool
		Error    bool
	}{
		{Field: "", Tag: "", FlagName: "", Required: false, Error: false},
		{Field: "FooBar", Tag: "", FlagName: "foo-bar", Required: false, Error: false},
		{Field: "FooBar", Tag: "bar-baz", FlagName: "bar-baz", Required: false, Error: false},
		{Field: "FooBar", Tag: "bar-baz,require", FlagName: "bar-baz", Required: true, Error: false},
		{Field: "FooBar", Tag: "bar-baz,inline", FlagName: "", Required: false, Error: false},
		{Field: "FooBar", Tag: "bar-baz,inline,require", FlagName: "", Required: true, Error: false},
		{Field: "FooBar", Tag: ",inline,require", FlagName: "", Required: true, Error: false},
		{Field: "FooBar", Tag: ",require", FlagName: "foo-bar", Required: true, Error: false},
		{Field: "FooBar", Tag: ",invalidoption", FlagName: "", Required: false, Error: true},
		{Field: "FooBar", Tag: ",", FlagName: "", Required: false, Error: true},
	} {
		got, required, err := getFlagName(test.Field, test.Tag)
		if test.Error && err == nil {
			t.Errorf("getFlagName(%q, %q): expected error, got nil", test.Field, test.Tag)
			continue
		}
		if !test.Error && err != nil {
			t.Errorf("getFlagName(%q, %q): unexpected error: %v", test.Field, test.Tag, err)
		}
		if err != nil {
			continue
		}

		if got != test.FlagName {
			t.Errorf("getFlagName(%q, %q) = %q, expected %q", test.Field, test.Tag, got, test.FlagName)
		}
		if got != test.FlagName {
			t.Errorf("getFlagName(%q, %q) required = %v, expected %v", test.Field, test.Tag, required, test.Required)
		}
	}
}

func TestGetEnvName(t *testing.T) {
	for _, test := range []struct {
		Field string
		Tag   string

		// Expected:
		EnvName string
		Error   bool
	}{
		{Field: "", Tag: "", EnvName: "", Error: false},
		{Field: "FooBar", Tag: "", EnvName: "FOO_BAR", Error: false},
		{Field: "FooBar", Tag: "BAR_BAZ", EnvName: "BAR_BAZ", Error: false},
		{Field: "FooBar", Tag: "BAR_BAZ,inline", EnvName: "", Error: false},
		{Field: "FooBar", Tag: ",invalidoption", EnvName: "", Error: true},
		{Field: "FooBar", Tag: ",inline,inline", EnvName: "", Error: true},
		{Field: "FooBar", Tag: ",", EnvName: "", Error: true},
	} {
		got, err := getEnvName(test.Field, test.Tag)
		if test.Error && err == nil {
			t.Errorf("getEnvName(%q, %q): expected error, got nil", test.Field, test.Tag)
			continue
		}
		if !test.Error && err != nil {
			t.Errorf("getEnvName(%q, %q): unexpected error: %v", test.Field, test.Tag, err)
		}
		if err != nil {
			continue
		}

		if got != test.EnvName {
			t.Errorf("getEnvName(%q, %q) = %q, expected %q", test.Field, test.Tag, got, test.EnvName)
		}
	}
}

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
		got := toSnakeCase(test.In, sep)
		if got != test.Expected {
			t.Errorf("toSnakeCase(%q, %q) = %q, expected %q", test.In, sep, got, test.Expected)
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
		got := toUpperSnakeCase(test.In, sep)
		if got != test.Expected {
			t.Errorf("toUpperSnakeCase(%q, %q) = %q, expected %q", test.In, sep, got, test.Expected)
		}
	}
}

func TestStructToFlags(t *testing.T) {
	t.Run("not a struct", func(t *testing.T) {
		var v int

		_, err := StructToFlags(&v)
		if err == nil {
			t.Errorf("StructToFlags(%T): expected error, got nil", &v)
		}
	})

	t.Run("invalid flag option", func(t *testing.T) {
		type invalidField struct {
			FlagA int `flag:"flag-a,foobar"`
		}
		v := &invalidField{}

		_, err := StructToFlags(v)
		if err == nil {
			t.Errorf("StructToFlags(%T): expected error, got nil", v)
		}
	})

	t.Run("unsupported type", func(t *testing.T) {
		type notAValidType interface{}
		type unsupportedType struct {
			FlagA notAValidType `flag:"flag-a"`
		}
		v := &unsupportedType{}
		_, err := StructToFlags(v)
		if err == nil {
			t.Errorf("StructToFlags(%T): expected error, got nil", v)
		}
	})

	t.Run("nested struct failing", func(t *testing.T) {
		type validSubStruct struct {
			FlagA int `flag:"flag-a"`
		}
		type nestedStruct struct {
			SubA validSubStruct
		}
		v := &nestedStruct{}

		flags, err := StructToFlags(v)
		if err != nil {
			t.Errorf("StructToFlags(%T): unexpected error: %v", v, err)
		}
		if len(flags) != 1 {
			t.Errorf("len(StructToFlags(%T)) = %d, expected 1", v, len(flags))
		}
	})

	t.Run("nested struct failing", func(t *testing.T) {
		type invalidField struct {
			FlagA int `flag:"flag-a,foobar"`
		}
		type nestedStructInvalidField struct {
			SubA invalidField
		}
		v := &nestedStructInvalidField{}

		_, err := StructToFlags(v)
		if err == nil {
			t.Errorf("StructToFlags(%T): expected error, got nil", v)
		}
	})

	t.Run("ignored field", func(t *testing.T) {
		type ignoredField struct {
			FlagA int `flag:"-"`
		}
		v := &ignoredField{}

		_, err := StructToFlags(v)
		if err != nil {
			t.Errorf("StructToFlags(%T): unexpected error: %v", v, err)
		}
	})
}

func TestParseStruct(t *testing.T) {
	t.Run("invalid flag option", func(t *testing.T) {
		type invalidField struct {
			FlagA int `flag:"flag-a,foobar"`
		}
		v := &invalidField{}

		commandLineFlags := commandLineFlags(t)
		err := ParseStruct(v, commandLineFlags...)
		if err == nil {
			t.Errorf("ParseStruct(%T): expected error, got nil", v)
		}
	})

	t.Run("valid flags", func(t *testing.T) {

		type validField struct {
			FlagA int `env:"FLAG_A"`
		}
		v := validField{}

		os.Clearenv()
		os.Setenv("FLAG_A", "42")

		commandLineFlags := commandLineFlags(t)
		err := ParseStruct(&v, commandLineFlags...)
		if err != nil {
			t.Errorf("ParseStruct(%T): unexpected error: %v", v, err)
			return
		}

		expected := 42
		if v.FlagA != expected {
			t.Errorf("ParseStruct(%T).FlagA = %d, expected %d", v, v.FlagA, expected)
		}
	})
}

func ExampleParseStruct() {
	type Configuration struct {
		URL      *url.URL `flag:",required" typehint:"website_url"`
		Strings  []string
		Bool     bool `flag:"boolean" env:"BOOLEAN" usage:"a boolean flag"`
		Timeouts struct {
			ReadTimeout  time.Duration
			WriteTimeout time.Duration
		} `flag:",inline" env:",inline"`

		IgnoreMe float64 `flag:"-"`
	}

	conf := Configuration{}

	err := ParseStruct(&conf)
	if err != nil {
		return
	}
}

func ExampleStructToFlags() {
	type Configuration struct {
		URL      *url.URL `flag:",require" typehint:"website_url"`
		Strings  []string
		Bool     bool `flag:"boolean" env:"BOOLEAN" usage:"a boolean flag"`
		Timeouts struct {
			ReadTimeout  time.Duration
			WriteTimeout time.Duration
		} `flag:",inline" env:",inline"`

		IgnoreMe float64 `flag:"-"`
	}

	conf := Configuration{}

	flags, err := StructToFlags(&conf)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	c := &Config{
		FlagSet: flag.NewFlagSet("test-rig", flag.ContinueOnError),
		Flags:   flags,
	}
	c.FlagSet.SetOutput(os.Stdout)

	c.Usage()

	// Output:
	// Usage of test-rig:
	//   -url website_url           URL=website_url           (required)
	//   -strings []string          STRINGS=[]string          (default "[]")
	//   -boolean                   BOOLEAN=bool              a boolean flag (default "false")
	//   -read-timeout duration     READ_TIMEOUT=duration     (default "0s")
	//   -write-timeout duration    WRITE_TIMEOUT=duration    (default "0s")

}
