package rig

// Required marks a flag as required, updating the typehint with "required".
// Noop if Flag.Required is true.
func Required(f *Flag) *Flag {
	if f.Required {
		return f
	}

	typeHint := f.TypeHint
	switch typeHint {
	default:
		typeHint += ", required"
	case "":
		typeHint += "required"
	}

	return &Flag{
		Value:    f.Value,
		Name:     f.Name,
		Env:      f.Env,
		Usage:    f.Usage,
		Required: true,
		TypeHint: typeHint,
	}
}
