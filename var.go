package config

import "flag"

type VarValidator func(flag.Value) error

type varValidators struct {
	flag.Value
	validators []VarValidator
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

func Var(v flag.Value, flag, env, usage string, validators ...VarValidator) *Flag {
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
