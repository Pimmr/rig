package rig

func TypeHint(f *Flag, typeHint string) *Flag {
	return &Flag{
		Value:    f.Value,
		Name:     f.Name,
		Env:      f.Env,
		Usage:    f.Usage,
		Required: f.Required,
		TypeHint: typeHint,
	}
}
