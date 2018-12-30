package config

type StringValidator func(string) error

type stringValidators struct {
	*StringValue
	validators []StringValidator
}

func (v stringValidators) Set(s string) error {
	err := v.StringValue.Set(s)
	if err != nil {
		return err
	}

	for _, validator := range v.validators {
		err = validator(string(*v.StringValue))
		if err != nil {
			return err
		}
	}

	return nil
}

type StringValue string

func (s StringValue) String() string {
	return string(s)
}

func (s *StringValue) Set(val string) error {
	*s = StringValue(val)
	return nil
}

func String(v *string, flag, env, usage string, validators ...StringValidator) *Flag {
	return &Flag{
		Value: stringValidators{
			StringValue: (*StringValue)(v),
			validators:  validators,
		},
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: "string",
	}
}
