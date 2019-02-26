package validators

import (
	"fmt"
	"time"
)

type Duration func(time.Duration) error

func DurationRange(min, max time.Duration) Duration {
	return func(d time.Duration) error {
		if d < min {
			return fmt.Errorf("duration should be %s or more", min)
		}
		if d > max {
			return fmt.Errorf("duration should be %s or less", max)
		}

		return nil
	}
}

func DurationMin(min time.Duration) Duration {
	return func(d time.Duration) error {
		if d < min {
			return fmt.Errorf("duration should be %s or more", min)
		}

		return nil
	}
}

func DurationMax(max time.Duration) Duration {
	return func(d time.Duration) error {
		if d > max {
			return fmt.Errorf("duration should be %s or less", max)
		}

		return nil
	}
}

func DurationRounded(r time.Duration) Duration {
	return func(d time.Duration) error {
		if d.Round(r) != d {
			return fmt.Errorf("duration should be a multiple of %s", r)
		}

		return nil
	}
}
