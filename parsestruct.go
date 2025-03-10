package rig

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

type fieldInfo struct {
	field reflect.Value
	typ   reflect.StructField

	flag       string
	env        string
	usage      string
	typeHint   string
	required   bool
	positional bool

	isStruct bool
}

func getFieldInfo(field reflect.Value, typ reflect.StructField) (*fieldInfo, error) {
	flagName, required, positional, err := getFlagName(typ.Name, typ.Tag.Get("flag"))
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

		flag:       flagName,
		env:        envName,
		usage:      typ.Tag.Get("usage"),
		typeHint:   typ.Tag.Get("typehint"),
		required:   required,
		positional: positional,

		isStruct: field.Kind() == reflect.Struct && !isFlagValue(field),
	}

	if info.flag == "" && info.env == "" && !info.isStruct {
		return nil, nil
	}
	if !info.field.CanAddr() {
		return nil, fmt.Errorf(".%s: cannot get address", info.typ.Name)
	}
	info.field = info.field.Addr()

	return info, nil
}

func isFlagValue(field reflect.Value) bool {
	return field.Addr().Type().Implements(reflect.TypeOf((*flag.Value)(nil)).Elem())
}

const (
	inlineOpt     = "inline"
	requireOpt    = "require"
	positionalOpt = "positional"
)

func getFlagName(fieldName, tag string) (flagName string, required, positional bool, err error) {
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
		if t == positionalOpt {
			positional = true
			continue
		}

		return flagName, required, positional, fmt.Errorf("unknown flag option %q", t)
	}

	if !inline && flagName == "" {
		flagName = toSnakeCase(fieldName, "-")
	}

	return flagName, required, positional, nil
}

func getEnvName(fieldName, tag string) (envName string, err error) {
	tt := strings.Split(tag, ",")
	if len(tt) > 0 {
		envName = tt[0]
		tt = tt[1:]
	}
	if len(tt) > 1 {
		return envName, fmt.Errorf("too many env options")
	}
	if len(tt) == 1 {
		if tt[0] == inlineOpt {
			return "", nil
		}
		return envName, fmt.Errorf("unknown env option %q", tt[0])
	}

	if envName == "" {
		envName = toUpperSnakeCase(fieldName, "_")
	}

	if envName == "-" {
		return "", nil
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
		return nil, fmt.Errorf("%T is not a struct", v)
	}

	fields, err := flagInfo(val)
	if err != nil {
		return nil, err
	}
	flags := make([]*Flag, 0, len(fields))
	for _, info := range fields {
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
		f.Positional = info.positional
		flags = append(flags, f)
	}

	return flags, nil
}

func flagInfo(val reflect.Value) ([]*fieldInfo, error) {
	valType := val.Type()

	fields := make([]*fieldInfo, 0, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		info, err := getFieldInfo(val.Field(i), valType.Field(i))
		if err != nil {
			return fields, err
		}
		if info == nil {
			continue
		}

		fields = append(fields, info)
	}

	sort.Slice(fields, func(i, j int) bool {
		if fields[j].positional {
			return false
		}

		return strings.Compare(fields[i].flag, fields[j].flag) < 0
	})

	return fields, nil
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

func getCompatiblePointerToPointerElem(i interface{}) (reflect.Value, bool) {
	switch i.(type) {
	case **url.URL, **regexp.Regexp:
		return reflect.Value{}, false
	}

	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		return reflect.Value{}, false
	}
	if v.Type().Elem().Kind() != reflect.Ptr {
		return reflect.Value{}, false
	}

	return v.Elem(), true
}

func flagFromInterface(i interface{}, flagName, env, usage string) (f *Flag, err error) {
	elem, ok := getCompatiblePointerToPointerElem(i)
	if !ok {
		return flagFromInterfaceConcrete(i, flagName, env, usage)
	}

	f, err = flagFromInterfaceConcrete(elem.Interface(), flagName, env, usage)
	if err != nil {
		return nil, err
	}

	defer func() {
		r := recover()
		if r == nil {
			return
		}

		if errv, ok := r.(error); ok {
			err = errv
		} else {
			panic(r)
		}
	}()

	return Pointer(f, i), nil
}

//nolint:gocyclo
func flagFromInterfaceConcrete(i interface{}, flagName, env, usage string) (*Flag, error) {
	switch t := i.(type) {
	default:
		v, ok := i.(flag.Value)
		if ok {
			return Var(v, flagName, env, usage), nil
		}

		return nil, fmt.Errorf("unsupported type %T", i)
	case *int:
		return Int(t, flagName, env, usage), nil
	case *int64:
		return Int64(t, flagName, env, usage), nil
	case *int32:
		return Int32(t, flagName, env, usage), nil
	case *uint:
		return Uint(t, flagName, env, usage), nil
	case *uint64:
		return Uint64(t, flagName, env, usage), nil
	case *uint32:
		return Uint32(t, flagName, env, usage), nil
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
	case *[]int32:
		return Repeatable(t, Int32Generator(), flagName, env, usage), nil
	case *[]uint:
		return Repeatable(t, UintGenerator(), flagName, env, usage), nil
	case *[]uint64:
		return Repeatable(t, Uint64Generator(), flagName, env, usage), nil
	case *[]uint32:
		return Repeatable(t, Uint32Generator(), flagName, env, usage), nil
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
