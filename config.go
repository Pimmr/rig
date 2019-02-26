package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/Pimmr/config/validators"

	"github.com/pkg/errors"
)

type Flag struct {
	flag.Value
	Name     string
	Env      string
	Usage    string
	Required bool
	TypeHint string

	set          bool
	defaultValue string
}

type isBoolFlagger interface {
	IsBoolFlag() bool
}

func (f *Flag) Set(v string) error {
	err := f.Value.Set(v)
	if err != nil {
		return err
	}

	f.set = true
	return nil
}

func (f *Flag) IsBoolFlag() bool {
	if boolFlagger, ok := f.Value.(isBoolFlagger); ok {
		return boolFlagger.IsBoolFlag()
	}

	return false
}

func (f Flag) missingError() error {
	switch {
	default:
		return errors.New("configuration variable doesn't have a flag or environment variable specified")
	case f.Name != "" && f.Env != "":
		return errors.Errorf("missing command line flag -%s or environment variable %s", f.Name, f.Env)
	case f.Name != "":
		return errors.Errorf("missing command line flag -%s", f.Name)
	case f.Env != "":
		return errors.Errorf("missing environment variable %s", f.Env)
	}
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
	c.FlagSet.Usage = c.Usage

	for _, f := range c.Flags {
		f.defaultValue = f.Value.String()
		if f.Name == "" {
			continue
		}
		c.FlagSet.Var(f, f.Name, f.Usage)
	}

	err := c.FlagSet.Parse(arguments)
	if err != nil {
		return c.handleError(err)
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
				return c.handleError(errors.Wrapf(err, "invalid value %q for env variable %q", v, f.Env))
			}
			continue
		}

		err = f.Set(v)
		if err != nil {
			return c.handleError(errors.Wrapf(err, "invalid value %q for env variable %q", v, f.Env))
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

func (c *Config) Arg(i int) string {
	return c.FlagSet.Arg(i)
}

func (c *Config) Args() []string {
	return c.FlagSet.Args()
}

func (c *Config) handleError(err error) error {
	_, _ = fmt.Fprintf(c.FlagSet.Output(), "%s\n", err)
	c.Usage()
	switch c.FlagSet.ErrorHandling() {
	case flag.ExitOnError:
		os.Exit(2)
	case flag.PanicOnError:
		panic(err)
	}
	return err
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
		if f.Usage != "" {
			_, _ = fmt.Fprintf(c.FlagSet.Output(), "        %s", f.Usage)
			if f.defaultValue != "" && !f.Required {
				_, _ = fmt.Fprintf(c.FlagSet.Output(), " (default %q)", f.defaultValue)
			}
			_, _ = fmt.Fprint(c.FlagSet.Output(), "\n")
		} else if f.defaultValue != "" && !f.Required {
			_, _ = fmt.Fprintf(c.FlagSet.Output(), "        (default %q)\n", f.defaultValue)
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

func TypeHint(f *Flag, typeHint string) *Flag {
	return &Flag{
		Value:    f.Value,
		Name:     f.Name,
		Env:      f.Env,
		Usage:    f.Usage,
		Required: f.Required,
		TypeHint: typeHint,
	}
}

type Generator func() flag.Value

type sliceValue struct {
	value      reflect.Value
	generator  Generator
	validators []validators.Var
}

func (vs sliceValue) String() string {
	return "[...]"
}

func (vs sliceValue) Set(s string) error {
	ss := splitRepeatable(s)
	for _, sub := range ss {
		err := vs.set(sub)
		if err != nil {
			return err
		}
	}

	return nil
}

func splitRepeatable(in string) []string {
	var out []string
	var current []rune
	var escaping bool

	for _, c := range in {
		if escaping {
			escaping = false
			current = append(current, c)
			continue
		}
		switch c {
		default:
			current = append(current, c)
		case '\\':
			escaping = true
		case ',':
			out = append(out, string(current))
			current = []rune{}
		}
	}
	out = append(out, string(current))

	return out
}

func (vs sliceValue) set(s string) error {
	if vs.value.Kind() != reflect.Ptr {
		return errors.Errorf("expected pointer to slice, got %s instead", vs.value.Kind())
	}
	ind := reflect.Indirect(vs.value)
	if ind.Kind() != reflect.Slice {
		return errors.Errorf("expected pointer to slice, got pointer to %s instead", ind.Kind())
	}
	if !ind.CanSet() {
		return errors.Errorf("expected pointer to slice to be settable")
	}

	v := vs.generator()
	err := v.Set(s)
	if err != nil {
		return err
	}

	for _, validator := range vs.validators {
		err = validator(v)
		if err != nil {
			return err
		}
	}

	vv := reflect.Indirect(reflect.ValueOf(v))
	if !vv.Type().ConvertibleTo(ind.Type().Elem()) {
		return errors.Errorf("type %s cannot be converted to %s", vv.Type(), ind.Type().Elem())
	}
	vv = vv.Convert(ind.Type().Elem())
	ind.Set(reflect.Append(ind, vv))

	return nil
}

func Repeatable(v interface{}, generator Generator, flag, env, usage string, validators ...validators.Var) *Flag {
	return &Flag{
		Value: sliceValue{
			value:      reflect.ValueOf(v),
			generator:  generator,
			validators: validators,
		},
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "repeatable",
	}
}

func MakeGenerator(v flag.Value) Generator {
	// TODO(yazgazan): This function will necessitate great examples
	t := reflect.Indirect(reflect.ValueOf(v)).Type()
	return func() flag.Value {
		vv := reflect.New(t)
		ret, ok := vv.Interface().(flag.Value)
		if !ok {
			panic(errors.Errorf("expected to be able to cast to flag.Value when generating for %s", t))
		}

		return ret
	}
}
