package rig

import (
	"flag"
	"strconv"

	"github.com/Pimmr/rig/validators"
)

type uint64Validators struct {
	*uint64Value
	validators []validators.Uint64
}

func (v uint64Validators) Set(s string) error {
	err := v.uint64Value.Set(s)
	if err != nil {
		return err
	}

	for _, validator := range v.validators {
		err = validator(uint64(*v.uint64Value))
		if err != nil {
			return err
		}
	}

	return nil
}

type uint64Value uint64

func (i uint64Value) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uint64Value(v)
	return err
}

func Uint64(v *uint64, flag, env, usage string, validators ...validators.Uint64) *Flag {
	return &Flag{
		Value: uint64Validators{
			uint64Value: (*uint64Value)(v),
			validators:  validators,
		},
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "unsigned 64 bit integer",
	}
}

func Uint64Generator() Generator {
	return func() flag.Value {
		return new(uint64Value)
	}
}
