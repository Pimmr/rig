package validators

import "fmt"

// A Uint64 validator should return an error if the uint64 provided is not considered valid, nil otherwise.
type Uint64 func(uint64) error

// Uint64Range creates a Uint64 validator that fails when the uint64 is strictly smaller than `min` or strictly larger than `max`.
func Uint64Range(min, max uint64) Uint64 {
	return func(i uint64) error {
		if i < min {
			return fmt.Errorf("unsigned 64-bit integer should be %d or more", min)
		}
		if i > max {
			return fmt.Errorf("unsigned 64-bit integer should be %d or less", max)
		}

		return nil
	}
}

// Uint64Min creates a Uint64 validator that fails when the uint64 is strictly smaller than `min`.
func Uint64Min(min uint64) Uint64 {
	return func(i uint64) error {
		if i < min {
			return fmt.Errorf("unsigned 64-bit integer should be %d or more", min)
		}

		return nil
	}
}

// Uint64Max creates a Uint64 validator that fails when the uint64 is strictly larger than `max`.
func Uint64Max(max uint64) Uint64 {
	return func(i uint64) error {
		if i > max {
			return fmt.Errorf("unsigned 64-bit integer should be %d or less", max)
		}

		return nil
	}
}
