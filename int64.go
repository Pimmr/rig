package config

import (
	"flag"
	"strconv"

	"github.com/Pimmr/config/validators"
)

type int64Validators struct {
	*int64Value
	validators []validators.Int64
}

func (v int64Validators) Set(s string) error {
	err := v.int64Value.Set(s)
	if err != nil {
		return err
	}

	for _, validator := range v.validators {
		err = validator(int64(*v.int64Value))
		if err != nil {
			return err
		}
	}

	return nil
}

type int64Value int64

func (i int64Value) String() string {
	return strconv.Itoa(int(i))
}

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = int64Value(v)
	return err
}

func Int64(v *int64, flag, env, usage string, validators ...validators.Int64) *Flag {
	return &Flag{
		Value: int64Validators{
			int64Value: (*int64Value)(v),
			validators: validators,
		},
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "64 bit integer",
	}
}

func Int64Generator() Generator {
	return func() flag.Value {
		return new(int64Value)
	}
}
