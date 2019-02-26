package validators

import "testing"

func TestUintRange(t *testing.T) {
	for _, test := range []struct {
		min, max    uint
		value       uint
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
		err := UintRange(test.min, test.max)(test.value)
		if test.expectError && err == nil {
			t.Errorf("UintRange(%d, %d)(%d): expected error, got nil", test.min, test.max, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("UintRange(%d, %d)(%d): unexpected error: %s", test.min, test.max, test.value, err)
		}
	}
}

func TestUintMin(t *testing.T) {
	for _, test := range []struct {
		min         uint
		value       uint
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
		err := UintMin(test.min)(test.value)
		if test.expectError && err == nil {
			t.Errorf("UintMin(%d)(%d): expected error, got nil", test.min, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("UintMin(%d)(%d): unexpected error: %s", test.min, test.value, err)
		}
	}
}

func TestUintMax(t *testing.T) {
	for _, test := range []struct {
		max         uint
		value       uint
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
		err := UintMax(test.max)(test.value)
		if test.expectError && err == nil {
			t.Errorf("UintMax(%d)(%d): expected error, got nil", test.max, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("UintMax(%d)(%d): unexpected error: %s", test.max, test.value, err)
		}
	}
}
