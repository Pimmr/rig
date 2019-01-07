package config

import (
	"flag"
	"regexp"

	"github.com/Pimmr/config/validators"
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

type RegexpValue struct {
	Regexp **regexp.Regexp
}

func (r RegexpValue) String() string {
	if *r.Regexp == nil {
		return ""
	}
	return (*r.Regexp).String()
}

func (r *RegexpValue) Set(s string) error {
	var err error

	*r.Regexp, err = regexp.Compile(s)
	return err
}

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

func RegexpGenerator() Generator {
	// TODO(yazgazan): might not work, needs testing
	return func() flag.Value {
		return &RegexpValue{
			Regexp: new(*regexp.Regexp),
		}
	}
}
