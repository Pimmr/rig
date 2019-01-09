package validators

import (
	"fmt"
)

type Float64 func(float64) error

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

func Float64Min(min float64) Float64 {
	return func(f float64) error {
		if f < min {
			return fmt.Errorf("float64 should be %f or more", min)
		}

		return nil
	}
}

func Float64Max(max float64) Float64 {
	return func(f float64) error {
		if f > max {
			return fmt.Errorf("float64 should be %f or less", max)
		}

		return nil
	}
}
