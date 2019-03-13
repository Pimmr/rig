package structToFlags

import (
	"flag"
	"net/url"
	"reflect"
	"regexp"
	"time"

	"github.com/Pimmr/rig"
	"github.com/pkg/errors"
)

func StructToFlags(v interface{}) ([]*rig.Flag, error) {
	val := reflect.Indirect(reflect.ValueOf(v))
	if val.Kind() != reflect.Struct {
		return nil, errors.Errorf("%T is not a struct", v)
	}
	valType := val.Type()

	flags := make([]*rig.Flag, 0, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := valType.Field(i)
		flagName := fieldType.Tag.Get("rig-flag")
		env := fieldType.Tag.Get("rig-env")
		usage := fieldType.Tag.Get("rig-usage")
		typeHint := fieldType.Tag.Get("rig-typehint")
		required := fieldType.Tag.Get("rig-required")

		if flagName == "" && env == "" && field.Kind() != reflect.Struct {
			continue
		}
		if !field.CanAddr() {
			return nil, errors.Errorf("%s.%s: cannot get address", valType, fieldType.Name)
		}
		isStruct := field.Kind() == reflect.Struct
		field = field.Addr()

		if isStruct {
			ff, err := StructToFlags(field.Interface())
			if err != nil {
				return nil, err
			}
			flags = append(flags, prefix(ff, flagName, env)...)
			continue
		}

		f, err := flagFromInterface(field.Interface(), flagName, env, usage)
		if err != nil {
			return nil, err
		}
		f = applyTypeHint(f, typeHint)
		f = applyRequired(f, required == "true")
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
	if !required {
		return f
	}

	return rig.Required(f)
}

func prefix(ff []*rig.Flag, flagName, env string) []*rig.Flag {
	for _, f := range ff {
		if flagName != "" && f.Name != "" {
			f.Name = flagName + "-" + f.Name
		}
		if env != "" && f.Env != "" {
			f.Env = env + "_" + f.Env
		}
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
