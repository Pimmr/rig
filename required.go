package rig

// Required marks a flag as required, updating the typehint with "required".
// Noop if Flag.Required is true.
func Required(f *Flag) *Flag {
	if f.Required {
		return f
	}

	ret := *f
	ret.Required = true
	return &ret
}

func Positional(f *Flag) *Flag {
	if f.Positional {
		return f
	}

	ret := *f
	ret.Positional = true
	return &ret
}
