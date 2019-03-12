package validators

import "testing"

func TestUint64Range(t *testing.T) {
	for _, test := range []struct {
		min, max    uint64
		value       uint64
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
		err := Uint64Range(test.min, test.max)(test.value)
		if test.expectError && err == nil {
			t.Errorf("Uint64Range(%d, %d)(%d): expected error, got nil", test.min, test.max, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("Uint64Range(%d, %d)(%d): unexpected error: %s", test.min, test.max, test.value, err)
		}
	}
}

func TestUint64Min(t *testing.T) {
	for _, test := range []struct {
		min         uint64
		value       uint64
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
		err := Uint64Min(test.min)(test.value)
		if test.expectError && err == nil {
			t.Errorf("Uint64Min(%d)(%d): expected error, got nil", test.min, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("Uint64Min(%d)(%d): unexpected error: %s", test.min, test.value, err)
		}
	}
}

func TestUint64Max(t *testing.T) {
	for _, test := range []struct {
		max         uint64
		value       uint64
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
		err := Uint64Max(test.max)(test.value)
		if test.expectError && err == nil {
			t.Errorf("Uint64Max(%d)(%d): expected error, got nil", test.max, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("Uint64Max(%d)(%d): unexpected error: %s", test.max, test.value, err)
		}
	}
}
