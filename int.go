package config

import (
	"flag"
	"strconv"

	"github.com/Pimmr/config/validators"
)

type intValidators struct {
	*intValue
	validators []validators.Int
}

func (v intValidators) Set(s string) error {
	err := v.intValue.Set(s)
	if err != nil {
		return err
	}

	for _, validator := range v.validators {
		err = validator(int(*v.intValue))
		if err != nil {
			return err
		}
	}

	return nil
}

type intValue int

func (i intValue) String() string {
	return strconv.Itoa(int(i))
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	*i = intValue(v)
	return err
}

func Int(v *int, flag, env, usage string, validators ...validators.Int) *Flag {
	return &Flag{
		Value: intValidators{
			intValue:   (*intValue)(v),
			validators: validators,
		},
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "integer",
	}
}

func IntGenerator() Generator {
	return func() flag.Value {
		return new(intValue)
	}
}
