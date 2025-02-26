package validators

import "testing"

func TestUint32Range(t *testing.T) {
	for _, test := range []struct {
		min, max    uint32
		value       uint32
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
		err := Uint32Range(test.min, test.max)(test.value)
		if test.expectError && err == nil {
			t.Errorf("Uint32Range(%d, %d)(%d): expected error, got nil", test.min, test.max, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("Uint32Range(%d, %d)(%d): unexpected error: %s", test.min, test.max, test.value, err)
		}
	}
}

func TestUint32Min(t *testing.T) {
	for _, test := range []struct {
		min         uint32
		value       uint32
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
		err := Uint32Min(test.min)(test.value)
		if test.expectError && err == nil {
			t.Errorf("Uint32Min(%d)(%d): expected error, got nil", test.min, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("Uint32Min(%d)(%d): unexpected error: %s", test.min, test.value, err)
		}
	}
}

func TestUint32Max(t *testing.T) {
	for _, test := range []struct {
		max         uint32
		value       uint32
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
		err := Uint32Max(test.max)(test.value)
		if test.expectError && err == nil {
			t.Errorf("Uint32Max(%d)(%d): expected error, got nil", test.max, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("Uint32Max(%d)(%d): unexpected error: %s", test.max, test.value, err)
		}
	}
}
