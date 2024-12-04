package rig

import (
	"errors"
	"flag"
	"fmt"
)

// A Flag represents the state and definition of a flag.
type Flag struct {
	flag.Value
	Name     string
	Env      string
	Usage    string
	TypeHint string
	Required bool

	set          bool
	defaultValue string
}

type isBoolFlagger interface {
	IsBoolFlag() bool
}

// Set proxies the .Set method on the underlying flag.Value. It is used to keep track
// of wether a flag has been set or not.
func (f *Flag) Set(v string) error {
	err := f.Value.Set(v)
	if err != nil {
		return err
	}

	f.set = true
	return nil
}

// IsBoolFlag proxies the .IsBoolFlag method on the underlying flag.Value, if defined.
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
		return fmt.Errorf("missing command line flag -%s or environment variable %s", f.Name, f.Env)
	case f.Name != "":
		return fmt.Errorf("missing command line flag -%s", f.Name)
	case f.Env != "":
		return fmt.Errorf("missing environment variable %s", f.Env)
	}
}

// IsSet returns true if the flag was set via the command line or the environment
func (f Flag) IsSet() bool {
	return f.set
}
