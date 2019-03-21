package validators

import (
	"fmt"
)

// A Float64 validator should return an error if the float64 provided is not considered valid, nil otherwise.
type Float64 func(float64) error

// Float64Range creates a Float64 validator that fails when the float64 is strictly smaller than `min` or strictly larger than `max`.
func Float64Range(min, max float64) Float64 {
	return func(f float64) error {
		if f < min {
			return fmt.Errorf("float64 should be %f or more", min)
		}
		if f > max {
			return fmt.Errorf("float64 should be %f or less", max)
		}

		return nil
	}
}

// Float64Min creates a Float64 validator that fails when the float64 is strictly smaller than `min`.
func Float64Min(min float64) Float64 {
	return func(f float64) error {
		if f < min {
			return fmt.Errorf("float64 should be %f or more", min)
		}

		return nil
	}
}

// Float64Max creates a Float64 validator that fails when the float64 is strictly larger than `max`.
func Float64Max(max float64) Float64 {
	return func(f float64) error {
		if f > max {
			return fmt.Errorf("float64 should be %f or less", max)
		}

		return nil
	}
}
