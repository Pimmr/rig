package rig

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

// Parse uses a default Config to parse the flags provided using os.Args.
// This default Config uses a flag.FlagSet with its ErrorHandling set to flag.ExitOnError.
func Parse(flags ...*Flag) error {
	config := &Config{
		FlagSet: flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		Flags:   flags,
	}

	return config.Parse(os.Args[1:])
}

// A Config represents a set of flags to be parsed. The flags are only set on the underlying
// flag.FlagSet when Config.Parse is called.
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

// Parse parses the arguments provided, along with the environment variables (using os.Getenv).
// Flags parsed from the `arguments` take precedence over the environment variables.
// The argument list provided should not include the command name.
func (c *Config) Parse(arguments []string) error {
	c.FlagSet.Usage = c.Usage

	c.setDefaultValues()

	err := c.parseFlagset(arguments)
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

		if f.Name != "" { // we want to maintain `"flag".FlagSet.Visit`'s behavior
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

	err = c.handleMissingFlags()
	if err != nil {
		return c.handleError(err)
	}

	return nil
}

func (c *Config) parseFlagset(arguments []string) error {
	for _, f := range c.Flags {
		if f.Name == "" {
			continue
		}
		c.FlagSet.Var(f, f.Name, f.Usage)
	}

	err := c.FlagSet.Parse(arguments)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) handleMissingFlags() error {
	hasMissing := false
	for _, f := range c.Flags {
		if !f.Required || f.set {
			continue
		}

		fmt.Fprintln(c.FlagSet.Output(), f.missingError())
		hasMissing = true
	}
	if hasMissing {
		return errors.New("missing required values")
	}

	return nil
}

// Arg proxies the .Arg method on the underlying flag.Flagset
func (c *Config) Arg(i int) string {
	return c.FlagSet.Arg(i)
}

// Args proxies the .Args method on the underlying flag.Flagset
func (c *Config) Args() []string {
	return c.FlagSet.Args()
}

func (c *Config) handleError(err error) error {
	fmt.Fprintf(c.FlagSet.Output(), "%s\n", err)
	c.Usage()
	switch c.FlagSet.ErrorHandling() {
	case flag.ExitOnError:
		os.Exit(2)
	case flag.PanicOnError:
		panic(err)
	}
	return err
}

// Usage prints the usage for the flags to the output defined on the underlying flag.FlagSet.
func (c *Config) Usage() {
	c.setDefaultValues()

	fmt.Fprintf(c.FlagSet.Output(), "Usage of %s:\n", c.FlagSet.Name())
	lines := make([][]string, 0, len(c.Flags))
	for _, f := range c.Flags {
		if f.Name == "" && f.Env == "" {
			continue
		}

		line := c.flagUsage(f)
		lines = append(lines, line)
	}

	offsets := offsetsForLines(lines, 2, 4)

	for _, line := range lines {
		delta := 0
		for i, col := range line {
			for j := 0; j < offsets[i]-delta; j++ {
				fmt.Fprint(c.FlagSet.Output(), " ")
			}
			fmt.Fprintf(c.FlagSet.Output(), "%s", col)
			delta = len(col)
		}
		fmt.Fprintln(c.FlagSet.Output())
	}
}

func offsetsForLines(lines [][]string, margin, sep int) []int {
	offsets := []int{}
	for _, line := range lines {
		for i, col := range line {
			if i >= len(offsets) {
				offsets = append(offsets, len(col)+sep)
			} else if len(col)+sep > offsets[i] {
				offsets[i] = len(col) + sep
			}
		}
	}
	offsets = append([]int{margin}, offsets...)

	return offsets
}

func (c *Config) flagUsage(f *Flag) []string {
	line := []string{}
	switch {
	case f.Name != "" && f.Env != "":
		line = append(line, flagUsageExample(f), f.Env+"=value")
	case f.Name != "":
		line = append(line, flagUsageExample(f), "")
	case f.Env != "":
		line = append(line, "", f.Env+"=value")
	}
	if f.TypeHint != "" {
		line = append(line, fmt.Sprintf("(%s)", f.TypeHint))
	} else {
		line = append(line, "")
	}

	usage := c.flagUsageDoc(f)
	if usage != "" {
		line = append(line, usage)
	}

	return line
}

func flagUsageExample(f *Flag) string {
	if f.IsBoolFlag() {
		return fmt.Sprintf("-%s", f.Name)
	}

	return fmt.Sprintf("-%s value", f.Name)
}

func (c *Config) flagUsageDoc(f *Flag) string {
	s := ""

	if f.Usage != "" {
		s += f.Usage
		if f.defaultValue != "" && !f.Required {
			s += fmt.Sprintf(" (default %q)", f.defaultValue)
		}
	} else if f.defaultValue != "" && !f.Required {
		s += fmt.Sprintf("(default %q)", f.defaultValue)
	}

	return s
}
