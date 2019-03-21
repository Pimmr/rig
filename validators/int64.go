package validators

import "fmt"

// A Int64 validator should return an error if the int64 provided is not considered valid, nil otherwise.
type Int64 func(int64) error

// Int64Range creates a Int64 validator that fails when the int64 is strictly smaller than `min` or strictly larger than `max`.
func Int64Range(min, max int64) Int64 {
	return func(i int64) error {
		if i < min {
			return fmt.Errorf("64-bit integer should be %d or more", min)
		}
		if i > max {
			return fmt.Errorf("64-bit integer should be %d or less", max)
		}

		return nil
	}
}

// Int64Min creates a Int64 validator that fails when the int64 is strictly smaller than `min`.
func Int64Min(min int64) Int64 {
	return func(i int64) error {
		if i < min {
			return fmt.Errorf("64-bit integer should be %d or more", min)
		}

		return nil
	}
}

// Int64Max creates a Int64 validator that fails when the int64 is strictly larger than `max`.
func Int64Max(max int64) Int64 {
	return func(i int64) error {
		if i > max {
			return fmt.Errorf("64-bit integer should be %d or less", max)
		}

		return nil
	}
}
