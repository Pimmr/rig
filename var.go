package rig

import (
	"flag"
	"reflect"

	"github.com/Pimmr/rig/validators"
)

type varValidators struct {
	flag.Value
	validators []validators.Var
}

func (v varValidators) Set(s string) error {
	err := v.Value.Set(s)
	if err != nil {
		return err
	}

	for _, validator := range v.validators {
		err = validator(v.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v varValidators) New(i interface{}) flag.Value {
	return varValidators{
		Value:      i.(flag.Value),
		validators: v.validators,
	}
}

func (v varValidators) IsNil() bool {
	return v.Value == nil || reflect.ValueOf(v.Value).IsNil()
}

// Var creates a flag for a flag.Value variable.
func Var(v flag.Value, flag, env, usage string, validators ...validators.Var) *Flag {
	return &Flag{
		Value: varValidators{
			Value:      v,
			validators: validators,
		},
		Name:  flag,
		Env:   env,
		Usage: usage,
	}
}
