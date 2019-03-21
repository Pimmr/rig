package validators

import "fmt"

// A Uint validator should return an error if the uint provided is not considered valid, nil otherwise.
type Uint func(uint) error

// UintRange creates a Uint validator that fails when the uint is strictly smaller than `min` or strictly larger than `max`.
func UintRange(min, max uint) Uint {
	return func(i uint) error {
		if i < min {
			return fmt.Errorf("unsigned integer should be %d or more", min)
		}
		if i > max {
			return fmt.Errorf("unsigned integer should be %d or less", max)
		}

		return nil
	}
}

// UintMin creates a Uint validator that fails when the uint is strictly smaller than `min`.
func UintMin(min uint) Uint {
	return func(i uint) error {
		if i < min {
			return fmt.Errorf("unsigned integer should be %d or more", min)
		}

		return nil
	}
}

// UintMax creates a Uint validator that fails when the uint is strictly larger than `max`.
func UintMax(max uint) Uint {
	return func(i uint) error {
		if i > max {
			return fmt.Errorf("unsigned integer should be %d or less", max)
		}

		return nil
	}
}
