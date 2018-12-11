package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Flag struct {
	flag.Value
	Name     string
	Env      string
	Usage    string
	Required bool
	TypeHint string

	set bool
}

func (f *Flag) Set(v string) error {
	err := f.Value.Set(v)
	if err != nil {
		return err // TODO(yazgazan): wrap error
	}

	f.set = true
	return nil
}

func (f Flag) missingError() error {
	if f.Name != "" && f.Env != "" {
		return errors.Errorf("missing command line flag -%s or environment variable %s", f.Name, f.Env)
	}

	if f.Name != "" {
		return errors.Errorf("missing command line flag -%s", f.Name)
	}

	return errors.Errorf("missing environment variable %s", f.Env)
}

func Parse(flags ...*Flag) error {
	config := &Config{
		FlagSet: flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		Flags:   flags,
	}

	return config.Parse(os.Args[1:])
}

type Config struct {
	FlagSet *flag.FlagSet

	Flags []*Flag
}

func (c *Config) Parse(arguments []string) error {
	// TODO(yazgazan): support https://golang.org/pkg/flag/#FlagSet.ErrorHandling
	// TODO(yazgazan): support https://golang.org/pkg/flag/#FlagSet.Arg
	// TODO(yazgazan): support https://golang.org/pkg/flag/#FlagSet.Args
	c.FlagSet.Usage = c.Usage

	for _, f := range c.Flags {
		if f.Name == "" {
			continue
		}
		c.FlagSet.Var(f, f.Name, f.Usage)
	}

	err := c.FlagSet.Parse(arguments)
	if err != nil {
		return err // TODO(yazgazan): wrap error
	}

	for _, f := range c.Flags {
		if f.Env == "" || f.set { // environment variables should not overwrite the command-line arguments
			continue
		}
		v := os.Getenv(f.Env)
		if v == "" {
			continue
		}

		if f.Name != "" { // we want to maintain `"flag".FlagSet.Visit`'s behaviour
			err = c.FlagSet.Set(f.Name, v)
			if err != nil {
				return err
			}
			continue
		}

		err = f.Set(v)
		if err != nil {
			return err
		}
	}

	hasErrors := false
	for _, f := range c.Flags {
		if !f.Required || f.set {
			continue
		}

		_, _ = fmt.Fprintln(c.FlagSet.Output(), f.missingError())
		hasErrors = true
	}
	if hasErrors {
		c.Usage()
		os.Exit(2)
	}

	return nil
}

func (c *Config) Usage() {
	_, _ = fmt.Fprintf(c.FlagSet.Output(), "Usage of %s:\n", os.Args[0])
	for _, f := range c.Flags {
		if f.Name == "" || f.Env == "" {
			continue
		}

		if f.Name != "" && f.Env != "" {
			_, _ = fmt.Fprintf(c.FlagSet.Output(), "  -%s value, %s=value", f.Name, f.Env)
		} else if f.Name != "" {
			_, _ = fmt.Fprintf(c.FlagSet.Output(), "  -%s value", f.Name)
		} else if f.Env != "" {
			_, _ = fmt.Fprintf(c.FlagSet.Output(), "  %s=value", f.Env)
		}
		if f.TypeHint != "" {
			_, _ = fmt.Fprintf(c.FlagSet.Output(), " (%s)", f.TypeHint)
		}

		_, _ = fmt.Fprint(c.FlagSet.Output(), "\n")
		defaultValue := f.Value.String()
		if f.Usage != "" {
			_, _ = fmt.Fprintf(c.FlagSet.Output(), "        %s", f.Usage)
			if defaultValue != "" && !f.Required {
				_, _ = fmt.Fprintf(c.FlagSet.Output(), " (default %q)", defaultValue)
			}
			_, _ = fmt.Fprint(c.FlagSet.Output(), "\n")
		} else if defaultValue != "" && !f.Required {
			_, _ = fmt.Fprintf(c.FlagSet.Output(), "        (default %q)\n", defaultValue)
		}
	}
}

func Required(f *Flag) *Flag {
	typeHint := f.TypeHint
	switch typeHint {
	default:
		typeHint += ", required"
	case "":
		typeHint += "required"
	}

	return &Flag{
		Value:    f.Value,
		Name:     f.Name,
		Env:      f.Env,
		Usage:    f.Usage,
		Required: true,
		TypeHint: typeHint,
	}
}

func HintType(f *Flag, typeHint string) *Flag {
	return &Flag{
		Value:    f.Value,
		Name:     f.Name,
		Env:      f.Env,
		Usage:    f.Usage,
		Required: f.Required,
		TypeHint: typeHint,
	}
}

func Var(v flag.Value, flag, env, usage string) *Flag {
	return &Flag{
		Value: v,
		Name:  flag,
		Env:   env,
		Usage: usage,
	}
}

type stringValue string

func (s stringValue) String() string {
	return string(s)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func String(v *string, flag, env, usage string) *Flag {
	return &Flag{
		Value:    (*stringValue)(v),
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "string",
	}
}

type boolValue bool

func (b boolValue) String() string {
	return strconv.FormatBool(bool(b))
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	*b = boolValue(v)
	return err
}

func Bool(v *bool, flag, env, usage string) *Flag {
	return &Flag{
		Value:    (*boolValue)(v),
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "boolean",
	}
}

type durationValue time.Duration

func (d durationValue) String() string {
	return time.Duration(d).String()
}

func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	*d = durationValue(v)
	return err
}

func Duration(v *time.Duration, flag, env, usage string) *Flag {
	return &Flag{
		Value:    (*durationValue)(v),
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "duration",
	}
}

type float64Value float64

func (f float64Value) String() string {
	return strconv.FormatFloat(float64(f), 'g', -1, 64)
}

func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	*f = float64Value(v)
	return err
}

func Float64(v *float64, flag, env, usage string) *Flag {
	return &Flag{
		Value:    (*float64Value)(v),
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "float",
	}
}

type intValue int

func (i intValue) String() string {
	return strconv.Itoa(int(i))
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	*i = intValue(v)
	return err
}

func Int(v *int, flag, env, usage string) *Flag {
	return &Flag{
		Value:    (*intValue)(v),
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "integer",
	}
}

type int64Value int64

func (i int64Value) String() string {
	return strconv.Itoa(int(i))
}

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = int64Value(v)
	return err
}

func Int64(v *int64, flag, env, usage string) *Flag {
	return &Flag{
		Value:    (*int64Value)(v),
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "64 bit integer",
	}
}

type uintValue uint

func (i uintValue) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, strconv.IntSize)
	*i = uintValue(v)
	return err
}

func Uint(v *uint, flag, env, usage string) *Flag {
	return &Flag{
		Value:    (*uintValue)(v),
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "unsigned integer",
	}
}

type uint64Value uint64

func (i uint64Value) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uint64Value(v)
	return err
}

func Uint64(v *uint64, flag, env, usage string) *Flag {
	return &Flag{
		Value:    (*uint64Value)(v),
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "unsigned 64 bit integer",
	}
}

// V is sugar for the default flag functions. It will panic if the type is not supported.
// Supported types are: *string, *bool, *time.Duration, *float64, *int, *int64, *uint, *uint64
func V(v interface{}, flag, env, usage string) *Flag {
	switch t := v.(type) {
	default:
		panic(errors.Errorf("unsupported type %T. Supported types: *string, *bool, *time.Duration, *float64, *int, *int64, *uint, *uint64", v))
	case *string:
		return String(t, flag, env, usage)
	case *bool:
		return Bool(t, flag, env, usage)
	case *time.Duration:
		return Duration(t, flag, env, usage)
	case *float64:
		return Float64(t, flag, env, usage)
	case *int:
		return Int(t, flag, env, usage)
	case *int64:
		return Int64(t, flag, env, usage)
	case *uint:
		return Uint(t, flag, env, usage)
	case *uint64:
		return Uint64(t, flag, env, usage)
	}
}
