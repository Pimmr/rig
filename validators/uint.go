package validators

import "fmt"

type Uint func(uint) error

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

func UintMin(min uint) Uint {
	return func(i uint) error {
		if i < min {
			return fmt.Errorf("unsigned integer should be %d or more", min)
		}

		return nil
	}
}

func UintMax(max uint) Uint {
	return func(i uint) error {
		if i > max {
			return fmt.Errorf("unsigned integer should be %d or less", max)
		}

		return nil
	}
}
