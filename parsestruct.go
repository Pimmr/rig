package rig

import (
	"flag"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/pkg/errors"
)

type fieldInfo struct {
	field reflect.Value
	typ   reflect.StructField

	flag     string
	env      string
	usage    string
	typeHint string
	required bool

	isStruct bool
}

func getFieldInfo(field reflect.Value, typ reflect.StructField) (*fieldInfo, error) {
	flagName, required, err := getFlagName(typ.Name, typ.Tag.Get("flag"))
	if err != nil {
		return nil, err
	}
	if !field.CanInterface() || flagName == "-" {
		return nil, nil
	}
	envName, err := getEnvName(typ.Name, typ.Tag.Get("env"))
	if err != nil {
		return nil, err
	}

	info := &fieldInfo{
		field: field,
		typ:   typ,

		flag:     flagName,
		env:      envName,
		usage:    typ.Tag.Get("usage"),
		typeHint: typ.Tag.Get("typehint"),
		required: required,

		isStruct: field.Kind() == reflect.Struct && !isFlagValue(field),
	}

	if info.flag == "" && info.env == "" && !info.isStruct {
		return nil, nil
	}
	if !info.field.CanAddr() {
		return nil, errors.Errorf(".%s: cannot get address", info.typ.Name)
	}
	info.field = info.field.Addr()

	return info, nil
}

func isFlagValue(field reflect.Value) bool {
	return field.Addr().Type().Implements(reflect.TypeOf((*flag.Value)(nil)).Elem())
}

const (
	inlineOpt  = "inline"
	requireOpt = "require"
)

func getFlagName(fieldName, tag string) (flagName string, required bool, err error) {
	inline := false
	tt := strings.Split(tag, ",")
	if len(tt) > 0 {
		flagName = tt[0]
		tt = tt[1:]
	}

	for _, t := range tt {
		if t == inlineOpt {
			inline = true
			flagName = ""
			continue
		}
		if t == requireOpt {
			required = true
			continue
		}

		return flagName, required, errors.Errorf("unknown flag option %q", t)
	}

	if !inline && flagName == "" {
		flagName = toSnakeCase(fieldName, "-")
	}

	return flagName, required, nil
}

func getEnvName(fieldName, tag string) (envName string, err error) {
	tt := strings.Split(tag, ",")
	if len(tt) > 0 {
		envName = tt[0]
		tt = tt[1:]
	}
	if len(tt) > 1 {
		return envName, errors.Errorf("too many env options")
	}
	if len(tt) == 1 {
		if tt[0] == inlineOpt {
			return "", nil
		}
		return envName, errors.Errorf("unknown env option %q", tt[0])
	}

	if envName == "" {
		envName = toUpperSnakeCase(fieldName, "_")
	}

	return envName, nil
}

func toSnakeCase(s, sep string) string {
	ret := ""
	prev := '\000'

	rr := []rune(s)
	for i, r := range rr {
		if i != 0 && unicode.IsUpper(r) && unicode.IsLower(prev) {
			ret += sep
		} else if i != 0 && i != len(rr)-1 && unicode.IsUpper(r) && unicode.IsUpper(prev) && unicode.IsLower(rr[i+1]) {
			ret += sep
		}
		prev = r
		ret += string(unicode.ToLower(r))
	}

	return ret
}

func toUpperSnakeCase(s, sep string) string {
	ret := ""
	prev := '\000'

	rr := []rune(s)
	for i, r := range rr {
		if i != 0 && unicode.IsUpper(r) && unicode.IsLower(prev) {
			ret += sep
		} else if i != 0 && i != len(rr)-1 && unicode.IsUpper(r) && unicode.IsUpper(prev) && unicode.IsLower(rr[i+1]) {
			ret += sep
		}
		prev = r
		ret += string(unicode.ToUpper(r))
	}

	return ret
}

// ParseStruct uses a default Config to parse the flages provided using os.Args.
// StructtoFlags is used to generate the flags. the additionalFlags are applied after
// the flags derived from the provided struct.
func ParseStruct(v interface{}, additionalFlags ...*Flag) error {
	flags, err := StructToFlags(v)
	if err != nil {
		return err
	}

	flags = append(flags, additionalFlags...)

	config := &Config{
		FlagSet: DefaultFlagSet(),
		Flags:   flags,
	}

	return config.Parse(os.Args[1:])
}

// StructToFlags generates a set of Flag based on the provided struct.
//
// StructToFlags recognizes four struct flags: "flag", "env", "typehint" and "usage".
// The flag and env names are inferred based on the field name unless values are provided in
// the struct tags.
// The field names are transformed from CamelCase to snake_case (using "-" as a separator for the flag).
//
// Additional options "inline" and "require" can be specified in the struct tags ("require" should be specified on the "flag" tag).
//
// A flag or env can be marked as ignored by using `flag:"-"` and `env:"-"` respectively
func StructToFlags(v interface{}) ([]*Flag, error) {
	val := reflect.Indirect(reflect.ValueOf(v))
	if val.Kind() != reflect.Struct {
		return nil, errors.Errorf("%T is not a struct", v)
	}
	valType := val.Type()

	flags := make([]*Flag, 0, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		info, err := getFieldInfo(val.Field(i), valType.Field(i))
		if err != nil {
			return nil, err
		}
		if info == nil {
			continue
		}

		if info.isStruct {
			ff, err := StructToFlags(info.field.Interface())
			if err != nil {
				return nil, err
			}
			flags = append(flags, prefix(ff, info.flag, info.env, info.required)...)
			continue
		}

		f, err := flagFromInterface(info.field.Interface(), info.flag, info.env, info.usage)
		if err != nil {
			return nil, err
		}
		f = applyTypeHint(f, info.typeHint)
		f = applyRequired(f, info.required)
		flags = append(flags, f)
	}

	return flags, nil
}

func applyTypeHint(f *Flag, typeHint string) *Flag {
	if typeHint == "" {
		return f
	}

	return TypeHint(f, typeHint)
}

func applyRequired(f *Flag, required bool) *Flag {
	if !required || f.Required {
		return f
	}

	return Required(f)
}

func prefix(ff []*Flag, flagName, env string, required bool) []*Flag {
	for i, f := range ff {
		if flagName != "" && f.Name != "" {
			f.Name = flagName + "-" + f.Name
		}
		if env != "" && f.Env != "" {
			f.Env = env + "_" + f.Env
		}
		ff[i] = applyRequired(f, required)
	}

	return ff
}

//nolint:gocyclo
func flagFromInterface(i interface{}, flagName, env, usage string) (*Flag, error) {
	switch t := i.(type) {
	default:
		v, ok := i.(flag.Value)
		if ok {
			return Var(v, flagName, env, usage), nil
		}

		return nil, errors.Errorf("unsupported type %T", i)
	case *int:
		return Int(t, flagName, env, usage), nil
	case *int64:
		return Int64(t, flagName, env, usage), nil
	case *uint:
		return Uint(t, flagName, env, usage), nil
	case *uint64:
		return Uint64(t, flagName, env, usage), nil
	case *string:
		return String(t, flagName, env, usage), nil
	case *bool:
		return Bool(t, flagName, env, usage), nil
	case *time.Duration:
		return Duration(t, flagName, env, usage), nil
	case *float64:
		return Float64(t, flagName, env, usage), nil
	case **regexp.Regexp:
		return Regexp(t, flagName, env, usage), nil
	case **url.URL:
		return URL(t, flagName, env, usage), nil

	case *[]int:
		return Repeatable(t, IntGenerator(), flagName, env, usage), nil
	case *[]int64:
		return Repeatable(t, Int64Generator(), flagName, env, usage), nil
	case *[]uint:
		return Repeatable(t, UintGenerator(), flagName, env, usage), nil
	case *[]uint64:
		return Repeatable(t, Uint64Generator(), flagName, env, usage), nil
	case *[]string:
		return Repeatable(t, StringGenerator(), flagName, env, usage), nil
	case *[]bool:
		return Repeatable(t, BoolGenerator(), flagName, env, usage), nil
	case *[]time.Duration:
		return Repeatable(t, DurationGenerator(), flagName, env, usage), nil
	case *[]float64:
		return Repeatable(t, Float64Generator(), flagName, env, usage), nil
	case *[]*regexp.Regexp:
		return Repeatable(t, RegexpGenerator(), flagName, env, usage), nil
	case *[]*url.URL:
		return Repeatable(t, URLGenerator(), flagName, env, usage), nil
	}
}
