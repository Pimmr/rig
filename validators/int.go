package validators

import (
	"fmt"
)

// A Int validator should return an error if the int provided is not considered valid, nil otherwise.
type Int func(int) error

// IntRange creates a Int validator that fails when the int is strictly smaller than `min` or strictly larger than `max`.
func IntRange(min, max int) Int {
	return func(i int) error {
		if i < min {
			return fmt.Errorf("integer should be %d or more", min)
		}
		if i > max {
			return fmt.Errorf("integer should be %d or less", max)
		}

		return nil
	}
}

// IntMin creates a Int validator that fails when the int is strictly smaller than `min`.
func IntMin(min int) Int {
	return func(i int) error {
		if i < min {
			return fmt.Errorf("integer should be %d or more", min)
		}

		return nil
	}
}

// IntMax creates a Int validator that fails when the int is strictly larger than `max`.
func IntMax(max int) Int {
	return func(i int) error {
		if i > max {
			return fmt.Errorf("integer should be %d or less", max)
		}

		return nil
	}
}
