package structToFlags

import (
	"flag"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/Pimmr/rig"
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

func parseBool(s string) (bool, error) {
	if s == "" {
		return false, nil
	}

	return strconv.ParseBool(s)
}

func getFieldInfo(field reflect.Value, typ reflect.StructField) (*fieldInfo, error) {
	ignore, err := parseBool(typ.Tag.Get("ignore"))
	if err != nil || ignore {
		return nil, err
	}

	required, err := parseBool(typ.Tag.Get("required"))
	if err != nil {
		return nil, err
	}

	info := &fieldInfo{
		field: field,
		typ:   typ,

		flag:     typ.Tag.Get("flag"),
		env:      typ.Tag.Get("env"),
		usage:    typ.Tag.Get("usage"),
		typeHint: typ.Tag.Get("typehint"),
		required: required,

		isStruct: field.Kind() == reflect.Struct,
	}

	if info.flag == "" && info.env == "" && !info.isStruct {
		return nil, nil
	}
	if !info.field.CanAddr() {
		return nil, errors.Errorf(".%s: cannot get address", info.typ.Name)
	}
	info.field = info.field.Addr()
	if !info.field.CanInterface() {
		return nil, errors.Errorf(".%s: cannot get interface", info.typ.Name)
	}

	return info, nil
}

func StructToFlags(v interface{}) ([]*rig.Flag, error) {
	val := reflect.Indirect(reflect.ValueOf(v))
	if val.Kind() != reflect.Struct {
		return nil, errors.Errorf("%T is not a struct", v)
	}
	valType := val.Type()

	flags := make([]*rig.Flag, 0, val.NumField())
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

func applyTypeHint(f *rig.Flag, typeHint string) *rig.Flag {
	if typeHint == "" {
		return f
	}

	return rig.TypeHint(f, typeHint)
}

func applyRequired(f *rig.Flag, required bool) *rig.Flag {
	if !required || f.Required {
		return f
	}

	return rig.Required(f)
}

func prefix(ff []*rig.Flag, flagName, env string, required bool) []*rig.Flag {
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
func flagFromInterface(i interface{}, flagName, env, usage string) (*rig.Flag, error) {
	switch t := i.(type) {
	default:
		v, ok := i.(flag.Value)
		if ok {
			return rig.Var(v, flagName, env, usage), nil
		}

		return nil, errors.Errorf("unsupported type %T", i)
	case *int:
		return rig.Int(t, flagName, env, usage), nil
	case *int64:
		return rig.Int64(t, flagName, env, usage), nil
	case *uint:
		return rig.Uint(t, flagName, env, usage), nil
	case *uint64:
		return rig.Uint64(t, flagName, env, usage), nil
	case *string:
		return rig.String(t, flagName, env, usage), nil
	case *bool:
		return rig.Bool(t, flagName, env, usage), nil
	case *time.Duration:
		return rig.Duration(t, flagName, env, usage), nil
	case *float64:
		return rig.Float64(t, flagName, env, usage), nil
	case **regexp.Regexp:
		return rig.Regexp(t, flagName, env, usage), nil
	case **url.URL:
		return rig.URL(t, flagName, env, usage), nil

	case *[]int:
		return rig.Repeatable(t, rig.IntGenerator(), flagName, env, usage), nil
	case *[]int64:
		return rig.Repeatable(t, rig.Int64Generator(), flagName, env, usage), nil
	case *[]uint:
		return rig.Repeatable(t, rig.UintGenerator(), flagName, env, usage), nil
	case *[]uint64:
		return rig.Repeatable(t, rig.Uint64Generator(), flagName, env, usage), nil
	case *[]string:
		return rig.Repeatable(t, rig.StringGenerator(), flagName, env, usage), nil
	case *[]bool:
		return rig.Repeatable(t, rig.BoolGenerator(), flagName, env, usage), nil
	case *[]time.Duration:
		return rig.Repeatable(t, rig.DurationGenerator(), flagName, env, usage), nil
	case *[]float64:
		return rig.Repeatable(t, rig.Float64Generator(), flagName, env, usage), nil
	case *[]rig.RegexpValue:
		return rig.Repeatable(t, rig.RegexpGenerator(), flagName, env, usage), nil
	case *[]rig.URLValue:
		return rig.Repeatable(t, rig.URLGenerator(), flagName, env, usage), nil
	}
}
