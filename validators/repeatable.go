package validators

import (
	"reflect"

	"github.com/pkg/errors"
)

// A Repeatable validator should return an error if the value provided is not considered valid, nil otherwise.
// This validator is used on individual values of a rig.Repeatable.
type Repeatable func(interface{}) error

// ToRepeatable turns some validator (i.e a func(int) error) into a validators.Repeatable (func(interface{}) error),
// removing the need to implement separate validators when dealing with repeatables.
func ToRepeatable(validator interface{}) Repeatable {
	val := reflect.ValueOf(validator)
	if val.Kind() != reflect.Func {
		panic(errors.Errorf("ToRepeatable: expected a function, got %T", validator))
	}
	valT := val.Type()
	if valT.NumIn() != 1 {
		panic(errors.Errorf("ToRepeatable: expected validator to accept 1 argument, got %d", valT.NumIn()))
	}
	argT := valT.In(0)

	if valT.NumOut() != 1 {
		panic(errors.Errorf("ToRepeatable: expected validator to return 1 value, got %d", valT.NumOut()))
	}
	retT := valT.Out(0)
	if !retT.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		panic(errors.Errorf("ToRepeatable: expected validator to return value of type error, got %v", retT))
	}

	return func(value interface{}) error {
		v := reflect.Indirect(reflect.ValueOf(value))
		vT := v.Type()

		if !vT.AssignableTo(argT) && vT.ConvertibleTo(argT) {
			v = v.Convert(argT)
			vT = v.Type()
		}

		if !vT.AssignableTo(argT) {
			return errors.Errorf("cannot use validator on type %v, expected %v", vT, argT)
		}

		out := val.Call([]reflect.Value{v})
		if out[0].IsNil() {
			return nil
		}

		return out[0].Interface().(error)
	}
}
