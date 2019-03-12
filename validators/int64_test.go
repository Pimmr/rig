package validators

import "testing"

func TestInt64Range(t *testing.T) {
	for _, test := range []struct {
		min, max    int64
		value       int64
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
		err := Int64Range(test.min, test.max)(test.value)
		if test.expectError && err == nil {
			t.Errorf("Int64Range(%d, %d)(%d): expected error, got nil", test.min, test.max, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("Int64Range(%d, %d)(%d): unexpected error: %s", test.min, test.max, test.value, err)
		}
	}
}

func TestInt64Min(t *testing.T) {
	for _, test := range []struct {
		min         int64
		value       int64
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
		err := Int64Min(test.min)(test.value)
		if test.expectError && err == nil {
			t.Errorf("Int64Min(%d)(%d): expected error, got nil", test.min, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("Int64Min(%d)(%d): unexpected error: %s", test.min, test.value, err)
		}
	}
}

func TestInt64Max(t *testing.T) {
	for _, test := range []struct {
		max         int64
		value       int64
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
		err := Int64Max(test.max)(test.value)
		if test.expectError && err == nil {
			t.Errorf("Int64Max(%d)(%d): expected error, got nil", test.max, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("Int64Max(%d)(%d): unexpected error: %s", test.max, test.value, err)
		}
	}
}
