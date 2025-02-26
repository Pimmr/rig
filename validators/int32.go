package validators

import "fmt"

// A Int32 validator should return an error if the int32 provided is not considered valid, nil otherwise.
type Int32 func(int32) error

// Int32Range creates a Int32 validator that fails when the int32 is strictly smaller than `min` or strictly larger than `max`.
func Int32Range(min, max int32) Int32 {
	return func(i int32) error {
		if i < min {
			return fmt.Errorf("32-bit integer should be %d or more", min)
		}
		if i > max {
			return fmt.Errorf("32-bit integer should be %d or less", max)
		}

		return nil
	}
}

// Int32Min creates a Int32 validator that fails when the int32 is strictly smaller than `min`.
func Int32Min(min int32) Int32 {
	return func(i int32) error {
		if i < min {
			return fmt.Errorf("32-bit integer should be %d or more", min)
		}

		return nil
	}
}

// Int32Max creates a Int32 validator that fails when the int32 is strictly larger than `max`.
func Int32Max(max int32) Int32 {
	return func(i int32) error {
		if i > max {
			return fmt.Errorf("32-bit integer should be %d or less", max)
		}

		return nil
	}
}
