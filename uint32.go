package rig

import (
	"flag"
	"strconv"

	"github.com/Pimmr/rig/validators"
)

type uint32Validators struct {
	*uint32Value
	validators []validators.Uint32
}

func (v uint32Validators) Set(s string) error {
	err := v.uint32Value.Set(s)
	if err != nil {
		return err
	}

	for _, validator := range v.validators {
		err = validator(uint32(*v.uint32Value))
		if err != nil {
			return err
		}
	}

	return nil
}

func (v uint32Validators) New(i interface{}) flag.Value {
	return uint32Validators{
		uint32Value: (*uint32Value)(i.(*uint32)),
		validators:  v.validators,
	}
}

func (v uint32Validators) IsNil() bool {
	return v.uint32Value == nil
}

type uint32Value uint32

func (i uint32Value) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

func (i *uint32Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 32)
	*i = uint32Value(v)
	return err
}

// Uint32 creates a flag for a uint32 variable.
func Uint32(v *uint32, flag, env, usage string, validators ...validators.Uint32) *Flag {
	return &Flag{
		Value: uint32Validators{
			uint32Value: (*uint32Value)(v),
			validators:  validators,
		},
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "uint32",
	}
}

// Uint32Generator is the default uint32 generator, to be used with Repeatable for uint32 slices.
func Uint32Generator() Generator {
	return func() flag.Value {
		return new(uint32Value)
	}
}
