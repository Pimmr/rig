package rig

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Pimmr/rig/validators"
)

func testingFlagset() *flag.FlagSet {
	f := flag.NewFlagSet("rig-test", flag.ContinueOnError)
	f.SetOutput(os.Stdout)

	return f
}

// commandLineFlags get flags from the command line. Parse and ParseStruct uses os.Args, so we need to copy `go test`'s flags to avoid errors on -test.* args
func commandLineFlags(t *testing.T) []*Flag {
	t.Helper()
	flags := []*Flag{}
	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		var s string

		flags = append(flags, String(&s, f.Name, "", ""))
	})

	return flags
}

func TestParse(t *testing.T) {
	var (
		s string
		i int
	)

	flags := commandLineFlags(t)

	flags = append(flags, String(&s, "string-flag", "STRING_ENV", ""))
	flags = append(flags, Int(&i, "int-flag", "INT_ENV", ""))

	os.Clearenv()
	os.Setenv("STRING_ENV", "foo")
	os.Setenv("INT_ENV", "42")
	err := Parse(flags...)
	if err != nil {
		t.Errorf("Parse(...): unexpected error: %s", err)
	}

	if s != "foo" {
		t.Errorf("Parse(...): STRING_ENV = %q, expected %q", s, "foo")
	}
	if i != 42 {
		t.Errorf("Parse(...): INT_ENV = %d, expected %d", i, 42)
	}
}

func TestConfigSetDefaultValues(t *testing.T) {
	var (
		s1 = "foo"
		s2 = "bar"
	)

	c := &Config{
		FlagSet: flag.NewFlagSet("flagset", flag.ContinueOnError),
		Flags: []*Flag{
			String(&s1, "string-1", "STRING_1", ""),
			String(&s2, "string-2", "STRING_2", ""),
		},
	}

	for _, f := range c.Flags {
		if f.defaultValue != "" {
			t.Errorf("defaultValue should not have been set for flag %q yet", f.Name)
		}
	}
	c.setDefaultValues()
	expected := []string{"foo", "bar"}
	for i, f := range c.Flags {
		if f.defaultValue != expected[i] {
			t.Errorf("flag %q has .defaultValue = %q, expected %q", f.Name, f.defaultValue, expected[i])
		}
	}

	s1 = "baz"
	s2 = "fuzz"
	c.setDefaultValues()
	for i, f := range c.Flags {
		if f.defaultValue != expected[i] {
			t.Errorf("flag %q has .defaultValue = %q, expected %q (defaultValue shouldn't have been re-set)", f.Name, f.defaultValue, expected[i])
		}
	}
}

func TestConfigParse(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		c := &Config{
			FlagSet: flag.NewFlagSet("flagset", flag.ContinueOnError),
			Flags:   []*Flag{},
		}
		buf := &bytes.Buffer{}
		c.FlagSet.SetOutput(buf)

		os.Clearenv()
		err := c.Parse([]string{})
		if err != nil {
			t.Errorf("Config.Parse([]): unexpected error: %s", err)
		}
	})

	t.Run("valid inputs from args", func(t *testing.T) {
		var (
			s string
			i int
		)
		c := &Config{
			FlagSet: flag.NewFlagSet("flagset", flag.ContinueOnError),
			Flags: []*Flag{
				String(&s, "string-flag", "STRING_ENV", ""),
				Int(&i, "int-flag", "INT_ENV", ""),
			},
		}
		buf := &bytes.Buffer{}
		c.FlagSet.SetOutput(buf)

		os.Clearenv()
		err := c.Parse([]string{"-string-flag=foo", "-int-flag=42"})
		if err != nil {
			t.Errorf("Config.Parse(...): unexpected error: %s", err)
		}

		expectedString := "foo"
		if s != expectedString {
			t.Errorf("-string-flag: got %q, expected %q", s, expectedString)
		}

		expectedInt := 42
		if i != expectedInt {
			t.Errorf("-int-flag: got %d, expected %d", i, expectedInt)
		}
	})

	t.Run("valid inputs from env", func(t *testing.T) {
		var (
			s string
			i int
		)
		c := &Config{
			FlagSet: flag.NewFlagSet("flagset", flag.ContinueOnError),
			Flags: []*Flag{
				String(&s, "string-flag", "STRING_ENV", ""),
				Int(&i, "int-flag", "INT_ENV", ""),
			},
		}
		buf := &bytes.Buffer{}
		c.FlagSet.SetOutput(buf)

		os.Clearenv()
		os.Setenv("STRING_ENV", "foo")
		os.Setenv("INT_ENV", "42")
		err := c.Parse([]string{})
		if err != nil {
			t.Errorf("Config.Parse([]): unexpected error: %s", err)
		}

		expectedString := "foo"
		if s != expectedString {
			t.Errorf("STRING_ENV: got %q, expected %q", s, expectedString)
		}

		expectedInt := 42
		if i != expectedInt {
			t.Errorf("INT_ENV: got %d, expected %d", i, expectedInt)
		}
	})

	t.Run("invalid inputs from args", func(t *testing.T) {
		var (
			s string
			i int
		)
		c := &Config{
			FlagSet: flag.NewFlagSet("flagset", flag.ContinueOnError),
			Flags: []*Flag{
				String(&s, "string-flag", "STRING_ENV", ""),
				Int(&i, "int-flag", "INT_ENV", ""),
			},
		}
		buf := &bytes.Buffer{}
		c.FlagSet.SetOutput(buf)

		os.Clearenv()
		err := c.Parse([]string{"-string-flag=foo", "-int-flag=bar"})
		if err == nil {
			t.Errorf("Config.Parse(invalidInput): expected error, got nil")
		}
	})

	t.Run("invalid inputs from env", func(t *testing.T) {
		var (
			s string
			i int
		)
		c := &Config{
			FlagSet: flag.NewFlagSet("flagset", flag.ContinueOnError),
			Flags: []*Flag{
				String(&s, "string-flag", "STRING_ENV", ""),
				Int(&i, "int-flag", "INT_ENV", ""),
			},
		}
		buf := &bytes.Buffer{}
		c.FlagSet.SetOutput(buf)

		os.Clearenv()
		os.Setenv("STRING_ENV", "foo")
		os.Setenv("INT_ENV", "bar")
		err := c.Parse([]string{})
		if err == nil {
			t.Errorf("Config.Parse([]): expected error, got nil")
		}
	})

	t.Run("no input, no required", func(t *testing.T) {
		var (
			s = "foo"
			i = 42
		)
		c := &Config{
			FlagSet: flag.NewFlagSet("flagset", flag.ContinueOnError),
			Flags: []*Flag{
				String(&s, "string-flag", "STRING_ENV", ""),
				Int(&i, "int-flag", "INT_ENV", ""),
			},
		}
		buf := &bytes.Buffer{}
		c.FlagSet.SetOutput(buf)

		os.Clearenv()
		err := c.Parse([]string{})
		if err != nil {
			t.Errorf("Config.Parse([]): unexpected error: %s", err)
		}

		expectedString := "foo"
		if s != expectedString {
			t.Errorf("STRING_ENV: got %q, expected %q", s, expectedString)
		}

		expectedInt := 42
		if i != expectedInt {
			t.Errorf("INT_ENV: got %d, expected %d", i, expectedInt)
		}
	})

	t.Run("no input, required", func(t *testing.T) {
		var (
			s string
			i int
		)
		c := &Config{
			FlagSet: flag.NewFlagSet("flagset", flag.ContinueOnError),
			Flags: []*Flag{
				Required(String(&s, "string-flag", "STRING_ENV", "")),
				Required(Int(&i, "int-flag", "INT_ENV", "")),
			},
		}
		buf := &bytes.Buffer{}
		c.FlagSet.SetOutput(buf)

		os.Clearenv()
		err := c.Parse([]string{})
		if err == nil {
			t.Errorf("Config.Parse([]): expected error, got nil")
		}
	})

	t.Run("missing required positional", func(t *testing.T) {
		var (
			s string
			i int
		)
		c := &Config{
			FlagSet: flag.NewFlagSet("flagset", flag.ContinueOnError),
			Flags: []*Flag{
				String(&s, "string-flag", "STRING_ENV", ""),
				Positional(Required(Int(&i, "int-flag", "INT_ENV", ""))),
			},
		}
		buf := &bytes.Buffer{}
		c.FlagSet.SetOutput(buf)

		os.Clearenv()
		err := c.Parse([]string{})
		if err == nil {
			t.Errorf("Config.Parse([]): expected error, got nil")
		}
	})

	t.Run("malformed positional", func(t *testing.T) {
		var (
			s string
			i int
		)
		c := &Config{
			FlagSet: flag.NewFlagSet("flagset", flag.ContinueOnError),
			Flags: []*Flag{
				String(&s, "string-flag", "STRING_ENV", ""),
				Positional(Int(&i, "int-flag", "INT_ENV", "")),
			},
		}
		buf := &bytes.Buffer{}
		c.FlagSet.SetOutput(buf)

		os.Clearenv()
		err := c.Parse([]string{"foo"})
		if err == nil {
			t.Errorf("Config.Parse([]): expected error, got nil")
		}
	})

	t.Run("too many positionals", func(t *testing.T) {
		var (
			s string
			i int
		)
		c := &Config{
			FlagSet: flag.NewFlagSet("flagset", flag.ContinueOnError),
			Flags: []*Flag{
				String(&s, "string-flag", "STRING_ENV", ""),
				Positional(Int(&i, "int-flag", "INT_ENV", "")),
			},
		}
		buf := &bytes.Buffer{}
		c.FlagSet.SetOutput(buf)

		os.Clearenv()
		err := c.Parse([]string{"42", "14"})
		if err == nil {
			t.Errorf("Config.Parse([]): expected error, got nil")
		}
	})

	t.Run("valid inputs from env, no flags", func(t *testing.T) {
		var (
			s string
			i int
		)
		c := &Config{
			FlagSet: flag.NewFlagSet("flagset", flag.ContinueOnError),
			Flags: []*Flag{
				String(&s, "", "STRING_ENV", ""),
				Int(&i, "", "INT_ENV", ""),
			},
		}
		buf := &bytes.Buffer{}
		c.FlagSet.SetOutput(buf)

		os.Clearenv()
		os.Setenv("STRING_ENV", "foo")
		os.Setenv("INT_ENV", "42")
		err := c.Parse([]string{})
		if err != nil {
			t.Errorf("Config.Parse([]): unexpected error: %s", err)
		}

		expectedString := "foo"
		if s != expectedString {
			t.Errorf("STRING_ENV: got %q, expected %q", s, expectedString)
		}

		expectedInt := 42
		if i != expectedInt {
			t.Errorf("INT_ENV: got %d, expected %d", i, expectedInt)
		}
	})

	t.Run("invalid inputs from env, no flags", func(t *testing.T) {
		var (
			s string
			i int
		)
		c := &Config{
			FlagSet: flag.NewFlagSet("flagset", flag.ContinueOnError),
			Flags: []*Flag{
				String(&s, "", "STRING_ENV", ""),
				Int(&i, "", "INT_ENV", ""),
			},
		}
		buf := &bytes.Buffer{}
		c.FlagSet.SetOutput(buf)

		os.Clearenv()
		os.Setenv("STRING_ENV", "foo")
		os.Setenv("INT_ENV", "bar")
		err := c.Parse([]string{})
		if err == nil {
			t.Errorf("Config.Parse([]): expected error, got nil")
		}
	})
}

func TestConfigArgArgs(t *testing.T) {
	var s string

	c := &Config{
		FlagSet: flag.NewFlagSet("flagset", flag.ContinueOnError),
		Flags: []*Flag{
			String(&s, "string-flag", "STRING_ENV", ""),
		},
	}

	const (
		flagArg = "-string-flag=foo"
		arg1    = "arg1"
		arg2    = "arg2"
		arg3    = "arg3"
	)
	args := []string{flagArg, arg1, arg2, arg3}
	err := c.Parse(args)
	if err != nil {
		t.Errorf("Config.Parse(%q): unexpected error: %s", args, err)
		t.FailNow()
	}

	expected := "foo"
	if s != expected {
		t.Errorf("Config.Parse(%q): string-flag should have been set to %q, got %q instead", args, expected, s)
	}

	if c.Arg(0) != arg1 {
		t.Errorf("Config.Arg(0) = %q, expected %q", c.Arg(0), arg1)
	}
	if c.Arg(1) != arg2 {
		t.Errorf("Config.Arg(1) = %q, expected %q", c.Arg(1), arg2)
	}
	if c.Arg(2) != arg3 {
		t.Errorf("Config.Arg(2) = %q, expected %q", c.Arg(2), arg3)
	}

	expectedArgs := []string{arg1, arg2, arg3}
	if !reflect.DeepEqual(c.Args(), expectedArgs) {
		t.Errorf("Config.Args() = %q, expected %q", c.Args(), expectedArgs)
	}
}

const testHandleErrorExitOnErrorEnv = "TEST_HANDLE_ERROR_CRASHER"

var errTest = errors.New("test error")

func TestHandleErrorExitOnError(t *testing.T) {
	if os.Getenv(testHandleErrorExitOnErrorEnv) != "1" {
		t.SkipNow()
	}
	c := &Config{
		FlagSet: flag.NewFlagSet("flagset", flag.ExitOnError),
	}

	_ = c.handleError(errTest)
	t.Logf("should have os.Exit, this code shouldn't have been reached")
}

func TestHandleError(t *testing.T) {
	t.Run("flag.ExitOnError", func(t *testing.T) {
		buf := &bytes.Buffer{}
		cmd := exec.Command(os.Args[0], "-test.run=TestHandleErrorExitOnError")
		cmd.Env = append(os.Environ(), testHandleErrorExitOnErrorEnv+"=1")
		cmd.Stdout = ioutil.Discard
		cmd.Stderr = buf

		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); !ok || e.Success() {
			t.Errorf("expected process to exit with error, got Success instead")
		}

		firstLine, err := buf.ReadString('\n')
		if err != nil {
			t.Errorf("unexpected error reading output's first line: %s", err)
			t.FailNow()
		}
		if firstLine != errTest.Error()+"\n" {
			t.Errorf("expected output buffer to start with %q, got %q instead", errTest, firstLine)
		}
	})

	t.Run("flag.PanicOnError", func(t *testing.T) {
		c := &Config{
			FlagSet: flag.NewFlagSet("flagset", flag.PanicOnError),
		}
		buf := &bytes.Buffer{}
		c.FlagSet.SetOutput(buf)
		defer func() {
			err := recover()
			if err != errTest {
				t.Errorf("expected panic to be %v, got %v instead", errTest, err)
			}
			firstLine, err := buf.ReadString('\n')
			if err != nil {
				t.Errorf("unexpected error reading output's first line: %s", err)
				t.FailNow()
			}
			if firstLine != errTest.Error()+"\n" {
				t.Errorf("expected output buffer to start with %q, got %q instead", errTest, firstLine)
			}
		}()

		_ = c.handleError(errTest)
		t.Errorf("should have panicked, this code shouldn't have been reached")
	})

	t.Run("flag.ContinueOnError", func(t *testing.T) {
		c := &Config{
			FlagSet: flag.NewFlagSet("flagset", flag.ContinueOnError),
		}
		buf := &bytes.Buffer{}
		c.FlagSet.SetOutput(buf)

		err := c.handleError(errTest)
		if err != errTest {
			t.Errorf("c.handleError(errTest) = %v, expected %v", err, errTest)
		}

		firstLine, err := buf.ReadString('\n')
		if err != nil {
			t.Errorf("unexpected error reading output's first line: %s", err)
			t.FailNow()
		}
		if firstLine != errTest.Error()+"\n" {
			t.Errorf("expected output buffer to start with %q, got %q instead", errTest, firstLine)
		}
	})
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
			TypeHint(String(&s, stringFlag, "", ""), ""),
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

func ExampleParse() {
	var ss []string
	var timeout time.Duration
	var f float64

	err := Parse(
		Repeatable(&ss, StringGenerator(), "strings", "STRINGS", "repeatable strings flag", validators.ToRepeatable(validators.StringLengthMin(2))),
		Duration(&timeout, "timeout", "TIMEOUT", "duration flag"),
		Float64(&f, "float64", "FLOAT64", "float64 flag", validators.Float64Max(1.4), validators.Float64Max(3.2)),
	)
	if err != nil {
		os.Exit(2)
	}
}

func ExampleConfig_Parse() {
	var ss []string
	var timeout time.Duration
	var f float64

	c := &Config{
		FlagSet: flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
		Flags: []*Flag{
			Repeatable(&ss, StringGenerator(), "strings", "STRINGS", "repeatable strings flag", validators.ToRepeatable(validators.StringLengthMin(2))),
			Duration(&timeout, "timeout", "TIMEOUT", "duration flag"),
			Float64(&f, "float64", "FLOAT64", "float64 flag", validators.Float64Min(1.4), validators.Float64Max(3.2)),
		},
	}
	err := c.Parse([]string{"-strings=foo,bar", "-timeout=1h20s", "-float64=2.1"})
	if err != nil {
		os.Exit(2)
	}

	fmt.Printf("ss: %q\n", ss)
	fmt.Printf("timeout: %v\n", timeout)
	fmt.Printf("f: %.2f\n", f)

	// Output:
	// ss: ["foo" "bar"]
	// timeout: 1h0m20s
	// f: 2.10
}
