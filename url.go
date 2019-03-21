package rig

import (
	"flag"
	"net/url"

	"github.com/Pimmr/rig/validators"
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

// A URLValue is a wrapper used to manipulate *url.URL flags.
// When using Repeatable for *url.URL, the slice should be of type []URLValue
type URLValue struct {
	URL **url.URL
}

func (u URLValue) String() string {
	if *u.URL == nil {
		return ""
	}

	return (*u.URL).String()
}

// Set parses and sets the url represented by `s`
func (u *URLValue) Set(s string) error {
	v, err := url.Parse(s)
	*u.URL = v
	return err
}

// URL creates a flag for a *url.URL variable.
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

// URLGenerator is the default *url.URL generator, to be used with Repeatable for url slices.
// the slices type must be []URLValue for the generator to work
func URLGenerator() Generator {
	return func() flag.Value {
		return &URLValue{
			URL: new(*url.URL),
		}
	}
}
