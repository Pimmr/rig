package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func Parse(flags ...*Flag) error {
	config := &Config{
		FlagSet: flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		Flags:   flags,
	}

	return config.Parse(os.Args[1:])
}

type Config struct {
	FlagSet *flag.FlagSet

	Flags            []*Flag
	defaultValuesSet bool
}

func (c *Config) setDefaultValues() {
	if c.defaultValuesSet {
		return
	}

	for _, f := range c.Flags {
		f.defaultValue = f.Value.String()
		if f.Name == "" {
			continue
		}
	}
	c.defaultValuesSet = true
}

func (c *Config) Parse(arguments []string) error {
	c.FlagSet.Usage = c.Usage

	c.setDefaultValues()
	for _, f := range c.Flags {
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

	hasMissing := false
	for _, f := range c.Flags {
		if !f.Required || f.set {
			continue
		}

		_, _ = fmt.Fprintln(c.FlagSet.Output(), f.missingError())
		hasMissing = true
	}
	if hasMissing {
		return c.handleError(errors.New("missing required values"))
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
	c.setDefaultValues()

	_, _ = fmt.Fprintf(c.FlagSet.Output(), "Usage of %s:\n", os.Args[0])
	for _, f := range c.Flags {
		if f.Name == "" && f.Env == "" {
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
