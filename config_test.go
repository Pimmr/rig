package config

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

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
			t.Errorf("expected proccess to exit with error, got Success instead")
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
