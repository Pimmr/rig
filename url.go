package rig

import (
	"flag"
	"net/url"

	"github.com/Pimmr/rig/validators"
)

type urlValidators struct {
	*urlValue
	validators []validators.URL
}

func (v urlValidators) Set(s string) error {
	err := v.urlValue.Set(s)
	if err != nil {
		return err
	}

	for _, validator := range v.validators {
		err = validator(*v.urlValue.URL)
		if err != nil {
			return err
		}
	}

	return nil
}

// A urlValue is a wrapper used to manipulate *url.URL flags.
// When using Repeatable for *url.URL, the slice should be of type []urlValue
type urlValue struct {
	URL **url.URL
}

func (u urlValue) String() string {
	if *u.URL == nil {
		return ""
	}

	return (*u.URL).String()
}

// Set parses and sets the url represented by `s`
func (u *urlValue) Set(s string) error {
	v, err := url.Parse(s)
	*u.URL = v
	return err
}

func (u *urlValue) Value() interface{} {
	return u.URL
}

// URL creates a flag for a *url.URL variable.
func URL(v **url.URL, flag, env, usage string, validators ...validators.URL) *Flag {
	return &Flag{
		Value: urlValidators{
			urlValue: &urlValue{
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
// the slices type must be []urlValue for the generator to work
func URLGenerator() Generator {
	return func() flag.Value {
		return &urlValue{
			URL: new(*url.URL),
		}
	}
}
