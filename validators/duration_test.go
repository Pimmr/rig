package validators

import (
	"testing"
	"time"
)

func TestDurationRange(t *testing.T) {
	for _, test := range []struct {
		min, max    time.Duration
		value       time.Duration
		expectError bool
	}{
		{
			min:         2 * time.Minute,
			max:         4 * time.Minute,
			value:       1 * time.Minute,
			expectError: true,
		},
		{
			min:         2 * time.Minute,
			max:         4 * time.Minute,
			value:       2 * time.Minute,
			expectError: false,
		},
		{
			min:         2 * time.Minute,
			max:         4 * time.Minute,
			value:       3 * time.Minute,
			expectError: false,
		},
		{
			min:         2 * time.Minute,
			max:         4 * time.Minute,
			value:       4 * time.Minute,
			expectError: false,
		},
		{
			min:         2 * time.Minute,
			max:         4 * time.Minute,
			value:       5 * time.Minute,
			expectError: true,
		},
	} {
		err := DurationRange(test.min, test.max)(test.value)
		if test.expectError && err == nil {
			t.Errorf("DurationRange(%s, %s)(%s): expected error, got nil", test.min, test.max, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("DurationRange(%s, %s)(%s): unexpected error: %s", test.min, test.max, test.value, err)
		}
	}
}

func TestDurationMin(t *testing.T) {
	for _, test := range []struct {
		min         time.Duration
		value       time.Duration
		expectError bool
	}{
		{
			min:         2 * time.Minute,
			value:       1 * time.Minute,
			expectError: true,
		},
		{
			min:         2 * time.Minute,
			value:       2 * time.Minute,
			expectError: false,
		},
		{
			min:         2 * time.Minute,
			value:       3 * time.Minute,
			expectError: false,
		},
	} {
		err := DurationMin(test.min)(test.value)
		if test.expectError && err == nil {
			t.Errorf("DurationMin(%s)(%s): expected error, got nil", test.min, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("DurationMin(%s)(%s): unexpected error: %s", test.min, test.value, err)
		}
	}
}

func TestDurationMax(t *testing.T) {
	for _, test := range []struct {
		max         time.Duration
		value       time.Duration
		expectError bool
	}{
		{
			max:         4 * time.Minute,
			value:       3 * time.Minute,
			expectError: false,
		},
		{
			max:         4 * time.Minute,
			value:       4 * time.Minute,
			expectError: false,
		},
		{
			max:         4 * time.Minute,
			value:       5 * time.Minute,
			expectError: true,
		},
	} {
		err := DurationMax(test.max)(test.value)
		if test.expectError && err == nil {
			t.Errorf("DurationMax(%s)(%s): expected error, got nil", test.max, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("DurationMax(%s)(%s): unexpected error: %s", test.max, test.value, err)
		}
	}
}

func TestDurationRounded(t *testing.T) {
	for _, test := range []struct {
		rounding    time.Duration
		value       time.Duration
		expectError bool
	}{
		{
			rounding:    2 * time.Minute,
			value:       8 * time.Minute,
			expectError: false,
		},
		{
			rounding:    2 * time.Minute,
			value:       7 * time.Minute,
			expectError: true,
		},
		{
			rounding:    2*time.Minute + 30*time.Second,
			value:       5 * time.Minute,
			expectError: false,
		},
		{
			rounding:    2*time.Minute + 30*time.Second,
			value:       6 * time.Minute,
			expectError: true,
		},
	} {
		err := DurationRounded(test.rounding)(test.value)
		if test.expectError && err == nil {
			t.Errorf("DurationRounded(%s)(%s): expected error, got nil", test.rounding, test.value)
		}
		if !test.expectError && err != nil {
			t.Errorf("DurationRounded(%s)(%s): unexpected error: %s", test.rounding, test.value, err)
		}
	}
}
