package rig

import (
	"flag"

	"github.com/Pimmr/rig/validators"
)

type stringValidators struct {
	*stringValue
	validators []validators.String
}

func (v stringValidators) Set(s string) error {
	_ = v.stringValue.Set(s) // stringValue.Set cannot return an error

	for _, validator := range v.validators {
		err := validator(string(*v.stringValue))
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

// String creates a flag for a string variable.
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

// StringGenerator is the default string generator, to be used with Repeatable for string slices.
func StringGenerator() Generator {
	return func() flag.Value {
		return new(stringValue)
	}
}
