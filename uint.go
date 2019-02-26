package config

import (
	"flag"
	"strconv"

	"github.com/Pimmr/config/validators"
)

type uintValidators struct {
	*uintValue
	validators []validators.Uint
}

func (v uintValidators) Set(s string) error {
	err := v.uintValue.Set(s)
	if err != nil {
		return err
	}

	for _, validator := range v.validators {
		err = validator(uint(*v.uintValue))
		if err != nil {
			return err
		}
	}

	return nil
}

type uintValue uint

func (i uintValue) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, strconv.IntSize)
	*i = uintValue(v)
	return err
}

func Uint(v *uint, flag, env, usage string, validators ...validators.Uint) *Flag {
	return &Flag{
		Value: uintValidators{
			uintValue:  (*uintValue)(v),
			validators: validators,
		},
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "unsigned integer",
	}
}

func UintGenerator() Generator {
	return func() flag.Value {
		return new(uintValue)
	}
}
