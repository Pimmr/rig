package validators

import "fmt"

// A Uint32 validator should return an error if the uint32 provided is not considered valid, nil otherwise.
type Uint32 func(uint32) error

// Uint32Range creates a Uint32 validator that fails when the uint32 is strictly smaller than `min` or strictly larger than `max`.
func Uint32Range(min, max uint32) Uint32 {
	return func(i uint32) error {
		if i < min {
			return fmt.Errorf("unsigned 32-bit integer should be %d or more", min)
		}
		if i > max {
			return fmt.Errorf("unsigned 32-bit integer should be %d or less", max)
		}

		return nil
	}
}

// Uint32Min creates a Uint32 validator that fails when the uint32 is strictly smaller than `min`.
func Uint32Min(min uint32) Uint32 {
	return func(i uint32) error {
		if i < min {
			return fmt.Errorf("unsigned 32-bit integer should be %d or more", min)
		}

		return nil
	}
}

// Uint32Max creates a Uint32 validator that fails when the uint32 is strictly larger than `max`.
func Uint32Max(max uint32) Uint32 {
	return func(i uint32) error {
		if i > max {
			return fmt.Errorf("unsigned 32-bit integer should be %d or less", max)
		}

		return nil
	}
}
