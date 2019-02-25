package config

import (
	"flag"
	"net/url"

	"github.com/Pimmr/config/validators"
)

type urlValidators struct {
	*URLValue
	validators []validators.URL
}

func (v urlValidators) Set(s string) error {
	err := v.URLValue.Set(s)
	if err != nil {
		return err
	}

	for _, validator := range v.validators {
		err = validator(*v.URLValue.URL)
		if err != nil {
			return err
		}
	}

	return nil
}

type URLValue struct {
	URL **url.URL
}

func (u URLValue) String() string {
	if *u.URL == nil {
		return ""
	}

	return (*u.URL).String()
}

func (u *URLValue) Set(s string) error {
	v, err := url.Parse(s)
	*u.URL = v
	return err
}

func URL(v **url.URL, flag, env, usage string, validators ...validators.URL) *Flag {
	return &Flag{
		Value: urlValidators{
			URLValue: &URLValue{
				URL: v,
			},
			validators: validators,
		},
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "URL",
	}
}

func URLGenerator() Generator {
	return func() flag.Value {
		return &URLValue{
			URL: new(*url.URL),
		}
	}
}
