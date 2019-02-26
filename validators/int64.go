package validators

import "fmt"

type Int64 func(int64) error

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

func Int64Min(min int64) Int64 {
	return func(i int64) error {
		if i < min {
			return fmt.Errorf("64-bit integer should be %d or more", min)
		}

		return nil
	}
}

func Int64Max(max int64) Int64 {
	return func(i int64) error {
		if i > max {
			return fmt.Errorf("64-bit integer should be %d or less", max)
		}

		return nil
	}
}
