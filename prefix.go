package config

func Prefix(namePrefix, envPrefix string, flags ...*Flag) []*Flag {
	ret := make([]*Flag, len(flags))

	for i, f := range flags {
		ret[i] = prefixFlag(namePrefix, envPrefix, f)
	}

	return ret
}

func prefixFlag(namePrefix, envPrefix string, f *Flag) *Flag {
	flagName := f.Name
	if namePrefix != "" {
		flagName = namePrefix + flagName
	}
	envName := f.Env
	if envPrefix != "" {
		envName = envPrefix + envName
	}

	return &Flag{
		Value:    f.Value,
		Name:     flagName,
		Env:      envName,
		Usage:    f.Usage,
		Required: f.Required,
		TypeHint: f.TypeHint,

		set:          f.set,
		defaultValue: f.defaultValue,
	}
}
