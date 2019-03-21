package rig

import (
	"flag"
	"time"

	"github.com/Pimmr/rig/validators"
)

type durationValidators struct {
	*durationValue
	validators []validators.Duration
}

func (v durationValidators) Set(s string) error {
	err := v.durationValue.Set(s)
	if err != nil {
		return err
	}

	for _, validator := range v.validators {
		err = validator(time.Duration(*v.durationValue))
		if err != nil {
			return err
		}
	}

	return nil
}

type durationValue time.Duration

func (d durationValue) String() string {
	return time.Duration(d).String()
}

func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	*d = durationValue(v)
	return err
}

// Duration creates a flag for a time.Duration variable.
func Duration(v *time.Duration, flag, env, usage string, validators ...validators.Duration) *Flag {
	return &Flag{
		Value: durationValidators{
			durationValue: (*durationValue)(v),
			validators:    validators,
		},
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "duration",
	}
}

// DurationGenerator is the default time.Duration generator, to be used with Repeatable for time.Duration slices.
func DurationGenerator() Generator {
	return func() flag.Value {
		return new(durationValue)
	}
}
