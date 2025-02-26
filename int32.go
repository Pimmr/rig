package rig

import (
	"flag"
	"strconv"

	"github.com/Pimmr/rig/validators"
)

type int32Validators struct {
	*int32Value
	validators []validators.Int32
}

func (v int32Validators) Set(s string) error {
	err := v.int32Value.Set(s)
	if err != nil {
		return err
	}

	for _, validator := range v.validators {
		err = validator(int32(*v.int32Value))
		if err != nil {
			return err
		}
	}

	return nil
}

func (v int32Validators) New(i interface{}) flag.Value {
	return int32Validators{
		int32Value: (*int32Value)(i.(*int32)),
		validators: v.validators,
	}
}

func (v int32Validators) IsNil() bool {
	return v.int32Value == nil
}

type int32Value int32

func (i int32Value) String() string {
	return strconv.Itoa(int(i))
}

func (i *int32Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 32)
	*i = int32Value(v)
	return err
}

// Int32 creates a flag for a int32 variable.
func Int32(v *int32, flag, env, usage string, validators ...validators.Int32) *Flag {
	return &Flag{
		Value: int32Validators{
			int32Value: (*int32Value)(v),
			validators: validators,
		},
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "int32",
	}
}

// Int32Generator is the default int32 generator, to be used with Repeatable for int32 slices.
func Int32Generator() Generator {
	return func() flag.Value {
		return new(int32Value)
	}
}
