package rig

import (
	"flag"
	"regexp"

	"github.com/Pimmr/rig/validators"
)

type regexpValidators struct {
	*regexpValue
	validators []validators.Regexp
}

func (v regexpValidators) Set(s string) error {
	err := v.regexpValue.Set(s)
	if err != nil {
		return err
	}

	for _, validator := range v.validators {
		err = validator(*v.regexpValue.Regexp)
		if err != nil {
			return err
		}
	}

	return nil
}

// A regexpValue is a wrapper used to manipulate *regexp.Regexp flags.
// When using Repeatable for *regexp.Regexp, the slice should be of type []regexpValue.
type regexpValue struct {
	Regexp **regexp.Regexp
}

func (r regexpValue) String() string {
	if *r.Regexp == nil {
		return ""
	}
	return (*r.Regexp).String()
}

// Set compiles and sets the regexp represented by `s`.
func (r *regexpValue) Set(s string) error {
	var err error

	*r.Regexp, err = regexp.Compile(s)
	return err
}

func (r regexpValue) Value() interface{} {
	return r.Regexp
}

// Regexp creates a flag for a *regexp.Regexp variable.
func Regexp(v **regexp.Regexp, flag, env, usage string, validators ...validators.Regexp) *Flag {
	return &Flag{
		Value: regexpValidators{
			regexpValue: &regexpValue{
				Regexp: v,
			},
			validators: validators,
		},
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "regex",
	}
}

// RegexpGenerator is the default *regexp.Regexp generator, to be used with Repeatable for regexp slices.
func RegexpGenerator() Generator {
	return func() flag.Value {
		return &regexpValue{
			Regexp: new(*regexp.Regexp),
		}
	}
}
