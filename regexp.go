package config

import "regexp"

type RegexpValidator func(*regexp.Regexp) error

type regexpValidators struct {
	*regexpValue
	validators []RegexpValidator
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

type regexpValue struct {
	Regexp **regexp.Regexp
}

func (r regexpValue) String() string {
	if *r.Regexp == nil {
		return ""
	}
	return (*r.Regexp).String()
}

func (r *regexpValue) Set(s string) error {
	var err error

	*r.Regexp, err = regexp.Compile(s)
	return err
}

func Regexp(v **regexp.Regexp, flag, env, usage string, validators ...RegexpValidator) *Flag {
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
		TypeHint: "Regexp",
	}
}
