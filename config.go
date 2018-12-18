package config

import (
	"flag"
	"fmt"
	"os"

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
