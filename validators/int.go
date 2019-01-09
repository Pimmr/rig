package validators

import (
	"fmt"
)

type Int func(int) error

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

func IntMin(min int) Int {
	return func(i int) error {
		if i < min {
			return fmt.Errorf("integer should be %d or more", min)
		}

		return nil
	}
}

func IntMax(max int) Int {
	return func(i int) error {
		if i > max {
			return fmt.Errorf("integer should be %d or less", max)
		}

		return nil
	}
}
