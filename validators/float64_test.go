package validators

import (
	"testing"
)

func TestFloat64Range(t *testing.T) {
	for _, test := range []struct {
		min, max    float64
		value       float64
		expectError bool
	}{
		{
			min:         2,
			max:         4,
			value:       1,
			expectError: true,
		},
		{
			min:         2,
			max:         4,
			value:       2,
			expectError: false,
		},
		{
			min:         2,
			max:         4,
			value:       3,
			expectError: false,
		},
		{
			min:         2,
			max:         4,
			value:       4,
			expectError: false,
		},
		{
			min:         2,
			max:         4,
			value:       5,
			expectError: true,
		},
	} {
		err := Float64Range(test.min, test.max)(test.value)
		if test.expectError && err == nil {
			t.Errorf("Float64Range(%f, %f)(%f): expected error, got nil", test.min, test.max, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("Float64Range(%f, %f)(%f): unexpected error: %s", test.min, test.max, test.value, err)
		}
	}
}

func TestFloat64Min(t *testing.T) {
	for _, test := range []struct {
		min         float64
		value       float64
		expectError bool
	}{
		{
			min:         2,
			value:       1,
			expectError: true,
		},
		{
			min:         2,
			value:       2,
			expectError: false,
		},
		{
			min:         2,
			value:       3,
			expectError: false,
		},
	} {
		err := Float64Min(test.min)(test.value)
		if test.expectError && err == nil {
			t.Errorf("Float64Min(%f)(%f): expected error, got nil", test.min, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("Float64Min(%f)(%f): unexpected error: %s", test.min, test.value, err)
		}
	}
}

func TestFloat64Max(t *testing.T) {
	for _, test := range []struct {
		max         float64
		value       float64
		expectError bool
	}{
		{
			max:         4,
			value:       3,
			expectError: false,
		},
		{
			max:         4,
			value:       4,
			expectError: false,
		},
		{
			max:         4,
			value:       5,
			expectError: true,
		},
	} {
		err := Float64Max(test.max)(test.value)
		if test.expectError && err == nil {
			t.Errorf("Float64Max(%f)(%f): expected error, got nil", test.max, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("Float64Max(%f)(%f): unexpected error: %s", test.max, test.value, err)
		}
	}
}
