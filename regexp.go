package rig

import (
	"flag"
	"regexp"

	"github.com/Pimmr/rig/validators"
)

type regexpValidators struct {
	*RegexpValue
	validators []validators.Regexp
}

func (v regexpValidators) Set(s string) error {
	err := v.RegexpValue.Set(s)
	if err != nil {
		return err
	}

	for _, validator := range v.validators {
		err = validator(*v.RegexpValue.Regexp)
		if err != nil {
			return err
		}
	}

	return nil
}

// A RegexpValue is a wrapper used to manipulate *regexp.Regexp flags.
// When using Repeatable for *regexp.Regexp, the slice should be of type []RegexpValue
type RegexpValue struct {
	Regexp **regexp.Regexp
}

func (r RegexpValue) String() string {
	if *r.Regexp == nil {
		return ""
	}
	return (*r.Regexp).String()
}

// Set compiles and sets the regexp represented by `s`
func (r *RegexpValue) Set(s string) error {
	var err error

	*r.Regexp, err = regexp.Compile(s)
	return err
}

// Regexp creates a flag for a *regexp.Regexp variable.
func Regexp(v **regexp.Regexp, flag, env, usage string, validators ...validators.Regexp) *Flag {
	return &Flag{
		Value: regexpValidators{
			RegexpValue: &RegexpValue{
				Regexp: v,
			},
			validators: validators,
		},
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "Regexp",
	}
}

// RegexpGenerator is the default *regexp.Regexp generator, to be used with Repeatable for regexp slices.
// the slices type must be []RegexpValue for the generator to work
func RegexpGenerator() Generator {
	return func() flag.Value {
		return &RegexpValue{
			Regexp: new(*regexp.Regexp),
		}
	}
}
