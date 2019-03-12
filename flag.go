package config

import (
	"flag"

	"github.com/pkg/errors"
)

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
