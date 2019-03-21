package validators

import (
	"fmt"
	"time"
)

// A Duration validator should return an error if the time.Duration provided is not considered valid, nil otherwise.
type Duration func(time.Duration) error

// DurationRange creates a Duration validator that fails when the time.Duration is strictly less than `min` or strictly more than `max`.
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

// DurationMin creates a Duration validator that fails when the time.Duration is strictly less than `min`.
func DurationMin(min time.Duration) Duration {
	return func(d time.Duration) error {
		if d < min {
			return fmt.Errorf("duration should be %s or more", min)
		}

		return nil
	}
}

// DurationMax creates a Duration validator that fails when the time.Duration is strictly more than `max`.
func DurationMax(max time.Duration) Duration {
	return func(d time.Duration) error {
		if d > max {
			return fmt.Errorf("duration should be %s or less", max)
		}

		return nil
	}
}

// DurationRounded creates a Duration validator that fails when the time.Duration is not a multiple of `r`
func DurationRounded(r time.Duration) Duration {
	return func(d time.Duration) error {
		if d.Round(r) != d {
			return fmt.Errorf("duration should be a multiple of %s", r)
		}

		return nil
	}
}
