package rig

import (
	"flag"
	"reflect"
	"testing"

	"github.com/Pimmr/rig/validators"
	"github.com/pkg/errors"
)

func newStringValue(s string) *stringValue {
	v := new(stringValue)
	*v = stringValue(s)

	return v
}

func newIntValue(i int) *intValue {
	v := new(intValue)
	*v = intValue(i)

	return v
}

func TestVar(t *testing.T) {
	testValidator := func(shouldFail bool) validators.Var {
		return func(flag.Value) error {
			if shouldFail {
				return errors.New("failing validator")
			}
			return nil
		}
	}

	for _, test := range []struct {
		val         flag.Value
		validators  []validators.Var
		input       string
		expected    flag.Value
		expectError bool
	}{
		{
			val:         newStringValue(""),
			input:       "foo",
			expected:    newStringValue("foo"),
			expectError: false,
		},
		{
			val:         newIntValue(0),
			input:       "42",
			expected:    newIntValue(42),
			expectError: false,
		},
		{
			val:         newIntValue(0),
			input:       "notanint",
			expectError: true,
		},
		{
			val:         newIntValue(0),
			validators:  []validators.Var{testValidator(false)},
			input:       "42",
			expected:    newIntValue(42),
			expectError: false,
		},
		{
			val:         newIntValue(0),
			validators:  []validators.Var{testValidator(true)},
			input:       "42",
			expectError: true,
		},
	} {
		v := Var(test.val, "flag", "ENV", "testing Var", test.validators...)
		err := v.Set(test.input)
		if test.expectError && err == nil {
			t.Errorf("Var(%T).Set(%q): expected error, got nil instead", test.val, test.input)
			continue
		}
		if !test.expectError && err != nil {
			t.Errorf("Var(%T).Set(%q): unexpected error: %s", test.val, test.input, err)
			continue
		}
		if err != nil {
			continue
		}
		if !reflect.DeepEqual(test.val, test.expected) {
			t.Errorf("Var(%T).Set(%q) = %+v, expected %+v", test.val, test.input, test.val, test.expected)
		}
	}
}
