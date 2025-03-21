package rig

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// Parse uses a default Config to parse the flags provided using os.Args.
// This default Config uses a flag.FlagSet with its ErrorHandling set to flag.ExitOnError.
func Parse(flags ...*Flag) error {
	config := &Config{
		FlagSet: DefaultFlagSet(),
		Flags:   flags,
	}

	return config.Parse(os.Args[1:])
}

// DefaultFlagSet returns the default FlagSet used by the Parse and ParseStruct functions.
func DefaultFlagSet() *flag.FlagSet {
	return flag.NewFlagSet(os.Args[0], flag.ExitOnError)
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
		v, ok := os.LookupEnv(f.Env)
		if !ok {
			continue
		}

		if f.Name != "" { // we want to maintain `"flag".FlagSet.Visit`'s behavior
			err = c.FlagSet.Set(f.Name, v)
			if err != nil {
				return c.handleError(fmt.Errorf("invalid value %q for env variable %q: %w", v, f.Env, err))
			}
			continue
		}

		err = f.Set(v)
		if err != nil {
			return c.handleError(fmt.Errorf("invalid value %q for env variable %q: %w", v, f.Env, err))
		}
	}

	err = parsePositionals(c.Flags, c.FlagSet.Args())
	if err != nil {
		return c.handleError(err)
	}

	err = c.handleMissingFlags()
	if err != nil {
		return c.handleError(err)
	}

	return nil
}

func parsePositionals(flags []*Flag, args []string) error {
	positionalFlags := []*Flag{}
	for _, flag := range flags {
		if !flag.Positional {
			continue
		}
		positionalFlags = append(positionalFlags, flag)
	}

	if len(positionalFlags) == 0 {
		return nil
	}

	for _, arg := range args {
		if len(positionalFlags) == 0 {
			return errors.New("unexpected positional argument")
		}

		err := positionalFlags[0].Set(arg)
		if err != nil {
			return err
		}

		if _, ok := positionalFlags[0].Value.(sliceValue); len(positionalFlags) > 1 || !ok {
			positionalFlags = positionalFlags[1:]
		}
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

func sortFlags(flags []*Flag) []*Flag {
	positional := []*Flag{}
	required := []*Flag{}
	others := make([]*Flag, 0, len(flags))

	for _, f := range flags {
		if f.Positional {
			positional = append(positional, f)
			continue
		}
		if f.Required {
			required = append(required, f)
			continue
		}

		others = append(others, f)
	}

	return append(required, append(others, positional...)...)
}

// Usage prints the usage for the flags to the output defined on the underlying flag.FlagSet.
func (c *Config) Usage() {
	c.setDefaultValues()

	c.Flags = sortFlags(c.Flags)

	lines := make([][]string, 0, len(c.Flags))
	hasNonPos := false
	positionals := []*Flag{}
	for _, f := range c.Flags {
		if f.Name == "" && f.Env == "" {
			continue
		}
		if !f.Positional {
			hasNonPos = true
			continue
		}

		line := c.flagUsage(f)
		lines = append(lines, line)

		positionals = append(positionals, f)
	}

	if len(lines) > 0 {
		lines = append(lines, []string{""})
	}

	for _, f := range c.Flags {
		if (f.Name == "" && f.Env == "") || f.Positional {
			continue
		}

		line := c.flagUsage(f)
		lines = append(lines, line)
	}

	b := &strings.Builder{}
	fmt.Fprintf(b, "Usage of %s", c.FlagSet.Name())
	if hasNonPos {
		fmt.Fprint(b, " [options]")
	}
	for _, pos := range positionals {
		name := pos.Env
		if name == "" {
			name = pos.Name
		}
		if _, ok := pos.Value.(sliceValue); ok {
			name += "..."
		}
		if pos.Required {
			fmt.Fprintf(b, " %s", name)
		} else {
			fmt.Fprintf(b, " [%s]", name)
		}
	}

	fmt.Fprint(b, ":\n")
	fmt.Fprint(c.FlagSet.Output(), b.String())
	// fmt.Fprintf(c.FlagSet.Output(), "Usage of %s:\n", c.FlagSet.Name())

	printUsageLines(c.FlagSet.Output(), lines, 2, 4)
}

func printUsageLines(output io.Writer, lines [][]string, margin, sep int) {
	offsets := offsetsForLines(lines, margin, sep)

	for _, line := range lines {
		delta := 0
		totalOffset := 0
		if len(line) == 1 && line[0] == "" {
			fmt.Fprint(output, "\n")
			continue
		}
		for i, col := range line {
			for j := 0; j < offsets[i]-delta; j++ {
				fmt.Fprint(output, " ")
			}
			if i < len(line)-1 {
				fmt.Fprintf(output, "%s", col)
				delta = len(col)
				totalOffset += offsets[i]
				continue
			}

			colLines := strings.Split(col, "\n")
			for k, colLine := range colLines {
				if k == 0 {
					fmt.Fprintf(output, "%s", colLine)
					continue
				}
				fmt.Fprint(output, "\n")
				for j := 0; j < offsets[i]+totalOffset; j++ {
					fmt.Fprint(output, " ")
				}
				fmt.Fprintf(output, "%s", strings.TrimLeft(colLine, " \t\r"))
			}
		}
		fmt.Fprintln(output)
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
	typ := f.TypeHint
	if typ == "" {
		typ = "value"
	}

	line := []string{}
	switch {
	case f.Name != "" && f.Env != "":
		line = append(line, flagUsageExample(f, typ), f.Env+"="+formatTypeHint(typ))
	case f.Name != "":
		line = append(line, flagUsageExample(f, typ), "")
	case f.Env != "":
		line = append(line, "", f.Env+"="+formatTypeHint(typ))
	}

	usage := c.flagUsageDoc(f)
	if usage != "" {
		line = append(line, usage)
	}

	return line
}

func formatTypeHint(typ string) string {
	if strings.IndexFunc(typ, unicode.IsSpace) == -1 {
		return typ
	}

	return strconv.Quote(typ)
}

func flagUsageExample(f *Flag, typ string) string {
	if f.IsBoolFlag() {
		return fmt.Sprintf("-%s", f.Name)
	}

	return fmt.Sprintf("-%s %s", f.Name, formatTypeHint(typ))
}

func (c *Config) flagUsageDoc(f *Flag) string {
	s := ""

	switch {
	case f.Usage != "":
		s += f.Usage
		if f.defaultValue != "" && !f.Required {
			s += fmt.Sprintf(" (default %q)", f.defaultValue)
		} else if f.Required {
			s += " (required)"
		}
	case f.defaultValue != "" && !f.Required:
		s += fmt.Sprintf("(default %q)", f.defaultValue)
	case f.Required:
		s += "(required)"
	}

	if f.Positional {
		switch s {
		default:
			s += " (positional)"
		case "":
			s = "(positional)"
		}
	}

	return s
}
