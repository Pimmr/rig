package validators

import "testing"

func TestIntRange(t *testing.T) {
	for _, test := range []struct {
		min, max    int
		value       int
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
		err := IntRange(test.min, test.max)(test.value)
		if test.expectError && err == nil {
			t.Errorf("IntRange(%d, %d)(%d): expected error, got nil", test.min, test.max, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("IntRange(%d, %d)(%d): unexpected error: %s", test.min, test.max, test.value, err)
		}
	}
}

func TestIntMin(t *testing.T) {
	for _, test := range []struct {
		min         int
		value       int
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
		err := IntMin(test.min)(test.value)
		if test.expectError && err == nil {
			t.Errorf("IntMin(%d)(%d): expected error, got nil", test.min, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("IntMin(%d)(%d): unexpected error: %s", test.min, test.value, err)
		}
	}
}

func TestIntMax(t *testing.T) {
	for _, test := range []struct {
		max         int
		value       int
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
		err := IntMax(test.max)(test.value)
		if test.expectError && err == nil {
			t.Errorf("IntMax(%d)(%d): expected error, got nil", test.max, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("IntMax(%d)(%d): unexpected error: %s", test.max, test.value, err)
		}
	}
}
