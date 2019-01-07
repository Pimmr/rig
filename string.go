package config

import (
	"flag"

	"github.com/Pimmr/config/validators"
)

type stringValidators struct {
	*stringValue
	validators []validators.String
}

func (v stringValidators) Set(s string) error {
	err := v.stringValue.Set(s)
	if err != nil {
		return err
	}

	for _, validator := range v.validators {
		err = validator(string(*v.stringValue))
		if err != nil {
			return err
		}
	}

	return nil
}

type stringValue string

func (s stringValue) String() string {
	return string(s)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func String(v *string, flag, env, usage string, validators ...validators.String) *Flag {
	return &Flag{
		Value: stringValidators{
			stringValue: (*stringValue)(v),
			validators:  validators,
		},
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "string",
	}
}

func StringGenerator() Generator {
	return func() flag.Value {
		return new(stringValue)
	}
}
