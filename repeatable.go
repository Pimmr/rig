package rig

import (
	"flag"
	"fmt"
	"reflect"
	"strings"

	"github.com/Pimmr/rig/validators"
	"github.com/pkg/errors"
)

// A Generator is a function that returns new values of a type implementing flag.Value.
// Generators are used with Repeatable to create a new value to be appended to the
// target slice.
type Generator func() flag.Value

type sliceValue struct {
	value      reflect.Value
	generator  Generator
	validators []validators.Var
}

func (vs sliceValue) String() string {
	if vs.value.CanInterface() {
		stringer, ok := vs.value.Interface().(fmt.Stringer)
		if ok {
			return stringer.String()
		}
	}

	value := reflect.Indirect(vs.value)

	ss := make([]string, value.Len())
	for i := 0; i < value.Len(); i++ {
		ss[i] = fmt.Sprint(value.Index(i))
	}
	return "[" + strings.Join(ss, ",") + "]"
}

func (vs sliceValue) Set(s string) error {
	ss := splitRepeatable(s)
	for _, sub := range ss {
		err := vs.set(sub)
		if err != nil {
			return err
		}
	}

	return nil
}

func splitRepeatable(in string) []string {
	var out []string
	var current []rune
	var escaping bool

	for _, c := range in {
		if escaping {
			escaping = false
			current = append(current, c)
			continue
		}
		switch c {
		default:
			current = append(current, c)
		case '\\':
			escaping = true
		case ',':
			out = append(out, string(current))
			current = []rune{}
		}
	}
	out = append(out, string(current))

	return out
}

func (vs sliceValue) set(s string) error {
	if vs.value.Kind() != reflect.Ptr {
		return errors.Errorf("expected pointer to slice, got %s instead", vs.value.Kind())
	}
	ind := reflect.Indirect(vs.value)
	if ind.Kind() != reflect.Slice {
		return errors.Errorf("expected pointer to slice, got pointer to %s instead", ind.Kind())
	}
	if !ind.CanSet() {
		return errors.Errorf("expected pointer to slice to be settable")
	}

	v := vs.generator()
	err := v.Set(s)
	if err != nil {
		return err
	}

	for _, validator := range vs.validators {
		err = validator(v)
		if err != nil {
			return err
		}
	}

	vv := reflect.Indirect(reflect.ValueOf(v))
	if !vv.Type().ConvertibleTo(ind.Type().Elem()) {
		return errors.Errorf("type %s cannot be converted to %s", vv.Type(), ind.Type().Elem())
	}
	vv = vv.Convert(ind.Type().Elem())
	ind.Set(reflect.Append(ind, vv))

	return nil
}

// Repeatable creates a flag that is repeatable. The variable `v` provided should be a pointer to a slice.
// The Generator should generates values that are assignable to the slice's emlements type.
func Repeatable(v interface{}, generator Generator, flag, env, usage string, validators ...validators.Var) *Flag {
	value := reflect.ValueOf(v)

	return &Flag{
		Value: sliceValue{
			value:      value,
			generator:  generator,
			validators: validators,
		},
		Name:     flag,
		Env:      env,
		Usage:    usage,
		TypeHint: strings.Replace(reflect.Indirect(value).Type().String(), "[]main.", "[]", 1),
	}
}

// MakeGenerator creates a Generator that will create values of the underlying type of `v`.
func MakeGenerator(v flag.Value) Generator {
	// TODO(yazgazan): This function will necessitate great examples
	val := reflect.ValueOf(v)
	isPtr := val.Kind() == reflect.Ptr
	t := reflect.Indirect(val).Type()
	return func() flag.Value {
		vv := reflect.New(t)
		if !isPtr {
			vv = reflect.Indirect(vv)
		}
		ret, ok := vv.Interface().(flag.Value)
		if !ok {
			panic(errors.Errorf("expected to be able to cast to flag.Value when generating for %s", t))
		}

		return ret
	}
}
