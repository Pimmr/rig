package config

import (
	"flag"
	"strconv"
)

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

func BoolGenerator() Generator {
	return func() flag.Value {
		return new(boolValue)
	}
}
