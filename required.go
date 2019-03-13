package rig

func Required(f *Flag) *Flag {
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
