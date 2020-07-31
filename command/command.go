package command

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Pimmr/rig"
	"github.com/Pimmr/rig/internal/text"
	"github.com/pkg/errors"
)

type Args struct {
	Arg0             string
	Args             []string
	EnvPrefix        string
	UsageDescription string
}

type Func func(remainder Args) error

type Config struct {
	Name string
	Fn   Func
}

func Parse(flags []*rig.Flag, cmds ...Config) error {
	c := rig.Config{
		FlagSet: flag.NewFlagSet(os.Args[0]+" [global flags] <command>", flag.ExitOnError),
		Flags:   flags,
	}
	c.UsageDescription = configsToUsageDescription(os.Args[0], cmds)

	err := c.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.HandleError(errors.New("missing command"))
	}

	cmd, ok := findCommand(cmds, args[0])
	if !ok {
		return c.HandleError(errors.Errorf("unknown command %q", args[0]))
	}

	cmdArgs := Args{
		Arg0:             os.Args[0] + " [global flags] " + args[0],
		Args:             args[1:],
		EnvPrefix:        text.ToUpperSnakeCase(cmd.Name, "_"),
		UsageDescription: fmt.Sprintf("Global flags:\n%s", c.FlagsUsage()),
	}

	return cmd.Fn(cmdArgs)
}

func configsToUsageDescription(arg0 string, configs []Config) string {
	if len(configs) == 0 {
		return ""
	}

	ss := make([]string, len(configs))
	for i, config := range configs {
		ss[i] = config.Name
	}

	return fmt.Sprintf("Commands: %s\nSee %s <command> -h for command's flags\n", strings.Join(ss, ", "), arg0)
}

func findCommand(cmds []Config, name string) (Config, bool) {
	for _, cmd := range cmds {
		if cmd.Name == name {
			return cmd, true
		}
	}

	return Config{}, false
}

func New(name string, fn Func) Config {
	return Config{
		Name: name,
		Fn:   fn,
	}
}

func (args Args) Parse(flags ...*rig.Flag) error {
	flags = rig.Prefix(flags, "", args.EnvPrefix, false)

	c := rig.Config{
		FlagSet:          flag.NewFlagSet(args.Arg0+" [flags]", flag.ExitOnError),
		Flags:            flags,
		UsageDescription: args.UsageDescription,
	}

	return c.Parse(args.Args)
}

func (args Args) ParseStruct(v interface{}, additionalFlags ...*rig.Flag) error {
	flags, err := rig.StructToFlags(v)
	if err != nil {
		return err
	}

	flags = append(flags, additionalFlags...)

	return args.Parse(flags...)
}

func (args Args) ParseCommands(flags []*rig.Flag, cmds ...Config) error {
	flags = rig.Prefix(flags, "", args.EnvPrefix, false)

	c := rig.Config{
		FlagSet:          flag.NewFlagSet(args.Arg0+" [flags] <command>", flag.ExitOnError),
		Flags:            flags,
		UsageDescription: args.UsageDescription + "\n" + configsToUsageDescription(args.Arg0, cmds),
	}

	err := c.Parse(args.Args)
	if err != nil {
		return err
	}

	cargs := c.Args()
	if len(cargs) == 0 {
		return c.HandleError(errors.New("missing command"))
	}

	cmd, ok := findCommand(cmds, cargs[0])
	if !ok {
		return c.HandleError(errors.Errorf("unknown command %q", cargs[0]))
	}

	envPrefix := text.ToUpperSnakeCase(cmd.Name, "_")
	if args.EnvPrefix != "" {
		envPrefix = args.EnvPrefix + "_" + envPrefix
	}

	cmdArgs := Args{
		Arg0:             args.Arg0 + " [flags] " + cargs[0],
		Args:             cargs[1:],
		EnvPrefix:        envPrefix,
		UsageDescription: fmt.Sprintf("Flags for %s:\n%s\n%s", args.Arg0, c.FlagsUsage(), args.UsageDescription),
	}

	return cmd.Fn(cmdArgs)
}
