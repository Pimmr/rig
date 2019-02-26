package validators

import "fmt"

type Uint64 func(uint64) error

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

func Uint64Min(min uint64) Uint64 {
	return func(i uint64) error {
		if i < min {
			return fmt.Errorf("unsigned 64-bit integer should be %d or more", min)
		}

		return nil
	}
}

func Uint64Max(max uint64) Uint64 {
	return func(i uint64) error {
		if i > max {
			return fmt.Errorf("unsigned 64-bit integer should be %d or less", max)
		}

		return nil
	}
}
