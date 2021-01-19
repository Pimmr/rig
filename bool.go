package rig

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

func (b *boolValue) IsBoolFlag() bool {
	return true
}

func (b boolValue) New(i interface{}) flag.Value {
	return (*boolValue)(i.(*bool))
}

func (b *boolValue) IsNil() bool {
	return b == nil
}

// Bool creates a flag for a boolean variable.
func Bool(v *bool, flag, env, usage string) *Flag {
	return &Flag{
		Value:    (*boolValue)(v),
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "bool",
	}
}

// BoolGenerator is the default boolean generator, to be used with Repeatable for boolean slices.
func BoolGenerator() Generator {
	return func() flag.Value {
		return new(boolValue)
	}
}
