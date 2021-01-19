package rig

import (
	"flag"
	"strconv"

	"github.com/Pimmr/rig/validators"
)

type float64Validators struct {
	*float64Value
	validators []validators.Float64
}

func (v float64Validators) Set(s string) error {
	err := v.float64Value.Set(s)
	if err != nil {
		return err
	}

	for _, validator := range v.validators {
		err = validator(float64(*v.float64Value))
		if err != nil {
			return err
		}
	}

	return nil
}

func (v float64Validators) New(i interface{}) flag.Value {
	return float64Validators{
		float64Value: (*float64Value)(i.(*float64)),
		validators:   v.validators,
	}
}

func (v float64Validators) IsNil() bool {
	return v.float64Value == nil
}

type float64Value float64

func (f float64Value) String() string {
	return strconv.FormatFloat(float64(f), 'g', -1, 64)
}

func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	*f = float64Value(v)
	return err
}

// Float64 creates a flag for a float64 variable.
func Float64(v *float64, flag, env, usage string, validators ...validators.Float64) *Flag {
	return &Flag{
		Value: float64Validators{
			float64Value: (*float64Value)(v),
			validators:   validators,
		},
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "float",
	}
}

// Float64Generator is the default float64 generator, to be used with Repeatable for float64 slices.
func Float64Generator() Generator {
	return func() flag.Value {
		return new(float64Value)
	}
}
