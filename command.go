package rig

import (
	"flag"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

type Command interface {
	Name() string
	Usage() string
	Call(invokedName string, arguments []string) error
}

type CommandsConfig struct {
	CommandName  string
	CommandUsage string
	Config       *Config
	Commands     map[string]Command

	lastError error
}

func Commands(commands ...Command) *CommandsConfig {
	cs := &CommandsConfig{
		Config: &Config{
			FlagSet: DefaultFlagSet(),
		},
		Commands: make(map[string]Command, len(commands)),
	}

	docs := make([][]string, len(commands))
	for i, command := range commands {
		if _, ok := cs.Commands[command.Name()]; ok {
			cs.lastError = errors.Errorf("command with name %q already set", command.Name())
			return cs
		}
		cs.Commands[command.Name()] = command
		docs[i] = append(docs[i], command.Name())
		usage := command.Usage()
		if usage != "" {
			docs[i] = append(docs[i], usage)
		}
	}

	cs.Config.UsageExtra = func() string {
		out := &strings.Builder{}

		printUsageLines(out, docs, 2, 4)

		return "\nCommands:\n" + out.String()
	}

	return cs
}

func SubCommands(name, usage string, commands ...Command) *CommandsConfig {
	cs := Commands(commands...)
	cs.CommandName = name
	cs.CommandUsage = usage

	return cs
}

func (cs *CommandsConfig) AdditionalFlags(flags ...*Flag) *CommandsConfig {
	cs.Config.Flags = append(cs.Config.Flags, flags...)

	return cs
}

func (cs *CommandsConfig) ParseAndCall() error {
	return cs.Call(os.Args[0], os.Args[1:])
}

func (cs *CommandsConfig) Call(name string, arguments []string) error {
	if cs.lastError != nil {
		return cs.Config.handleError(cs.lastError)
	}

	resetFlagSet(cs.Config, name)

	err := cs.Config.Parse(arguments)
	if err != nil {
		return err
	}

	args := cs.Config.Args()
	if len(args) == 0 {
		return cs.Config.handleError(errors.New("missing command"))
	}
	cmd, ok := cs.Commands[args[0]]
	if !ok {
		return cs.Config.handleError(errors.Errorf("%q: unknown command", args[0]))
	}

	return cmd.Call(name+" "+args[0], args[1:])
}

func (cs *CommandsConfig) Name() string {
	return cs.CommandName
}

func (cs *CommandsConfig) Usage() string {
	return cs.CommandUsage
}

type CallbackCommandConfig struct {
	CommandName  string
	CommandUsage string
	Config       *Config
	Callback     Command
}

func CallbackCommand(name string, fn Command, usage string, flags ...*Flag) *CallbackCommandConfig {
	return &CallbackCommandConfig{
		CommandName:  name,
		CommandUsage: usage,
		Config: &Config{
			FlagSet: DefaultFlagSet(),
			Flags:   flags,
		},
		Callback: fn,
	}
}

func (cc *CallbackCommandConfig) Call(name string, arguments []string) error {
	resetFlagSet(cc.Config, name)

	err := cc.Config.Parse(arguments)
	if err != nil {
		return err
	}

	return cc.Callback.Call(name, cc.Config.Args())
}

func (cc CallbackCommandConfig) Name() string {
	return cc.CommandName
}

func (cc CallbackCommandConfig) Usage() string {
	return cc.CommandUsage
}

type StructCommandConfig struct {
	CommandName  string
	CommandUsage string
	Config       *Config

	fn        reflect.Value
	arg       reflect.Value
	lastError error
}

func StructCommand(name string, fn interface{}, usage string, additionalFlags ...*Flag) *StructCommandConfig {
	fnV := reflect.ValueOf(fn)
	if fnV.Kind() != reflect.Func {
		return &StructCommandConfig{
			Config: &Config{
				FlagSet: DefaultFlagSet(),
			},
			lastError: errors.Errorf("expected fn to be a function, got %T instead", fn),
		}
	}
	fnT := fnV.Type()
	if fnT.NumIn() != 1 {
		return &StructCommandConfig{
			Config: &Config{
				FlagSet: DefaultFlagSet(),
			},
			lastError: errors.Errorf("expected fn to take 1 argument, got %d instead", fnT.NumIn()),
		}
	}
	if fnT.NumOut() != 1 {
		return &StructCommandConfig{
			Config: &Config{
				FlagSet: DefaultFlagSet(),
			},
			lastError: errors.Errorf("expected fn to return 1 value, got %d instead", fnT.NumOut()),
		}
	}
	if !fnT.Out(0).AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
		return &StructCommandConfig{
			Config: &Config{
				FlagSet: DefaultFlagSet(),
			},
			lastError: errors.Errorf("expected fn to return an 'error', got %v instead", fnT.Out(0)),
		}
	}

	argV := reflect.New(fnT.In(0))
	if !argV.CanInterface() {
		return &StructCommandConfig{
			Config: &Config{
				FlagSet: DefaultFlagSet(),
			},
			lastError: errors.Errorf("failed to instanciate value of type %v", fnT.In(0)),
		}
	}
	flags, err := StructToFlags(argV.Interface())
	if err != nil {
		return &StructCommandConfig{
			Config: &Config{
				FlagSet: DefaultFlagSet(),
			},
			lastError: err,
		}
	}

	return &StructCommandConfig{
		CommandName:  name,
		CommandUsage: usage,
		Config: &Config{
			FlagSet: DefaultFlagSet(),
			Flags:   append(flags, additionalFlags...),
		},
		fn:  fnV,
		arg: reflect.Indirect(argV),
	}
}

func (sc *StructCommandConfig) Call(name string, arguments []string) error {
	if sc.lastError != nil {
		return sc.Config.handleError(sc.lastError)
	}

	resetFlagSet(sc.Config, name)

	err := sc.Config.Parse(arguments)
	if err != nil {
		return err
	}

	ret := sc.fn.Call([]reflect.Value{sc.arg})
	if ret[0].IsNil() {
		return nil
	}
	if !ret[0].CanInterface() {
		return errors.New("failed to get error back")
	}

	err, ok := ret[0].Interface().(error)
	if !ok {
		return errors.Errorf("expected return value to be of type 'error', got %T instead", ret[0].Interface())
	}

	return err
}

func (sc *StructCommandConfig) Name() string {
	return sc.CommandName
}

func (sc *StructCommandConfig) Usage() string {
	return sc.CommandUsage
}

func resetFlagSet(config *Config, name string) {
	errorHandling := flag.ExitOnError
	output := io.Writer(os.Stderr)
	if config.FlagSet != nil {
		errorHandling = config.FlagSet.ErrorHandling()
		output = config.FlagSet.Output()
	}
	config.FlagSet = flag.NewFlagSet(name, errorHandling)
	config.FlagSet.SetOutput(output)
}
