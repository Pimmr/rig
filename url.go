package config

import "net/url"

type URLValidator func(*url.URL) error

type urlValidators struct {
	*urlValue
	validators []URLValidator
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

type urlValue struct {
	URL **url.URL
}

func (u urlValue) String() string {
	if *u.URL == nil {
		return ""
	}

	return (*u.URL).String()
}

func (u *urlValue) Set(s string) error {
	v, err := url.Parse(s)
	*u.URL = v
	return err
}

func URL(v **url.URL, flag, env, usage string, validators ...URLValidator) *Flag {
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
