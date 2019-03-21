package validators

import (
	"fmt"
	"strings"
)

// A String validator should return an error if the string provided is not considered valid, nil otherwise.
type String func(string) error

// StringNotEmpty creates a String validator that fails when the string is empty (after calling strings.TrimSpace).
func StringNotEmpty() String {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("string should not be empty")
		}

		return nil
	}
}

// StringLengthRange creates a String validator that fails when the string is strictly shorter than `min` or strictly longer than `max`.
func StringLengthRange(min, max int) String {
	return func(s string) error {
		if len(s) < min {
			return fmt.Errorf("string should be at least %d characters long", min)
		}
		if len(s) > max {
			return fmt.Errorf("string should be at most %d characters long", max)
		}

		return nil
	}
}

// StringLengthMin creates a String validator that fails when the string is strictly shorter than `min`.
func StringLengthMin(min int) String {
	return func(s string) error {
		if len(s) < min {
			return fmt.Errorf("string should be at least %d characters long", min)
		}

		return nil
	}
}

// StringLengthMax creates a String validator that fails when the string is strictly longer than `max`.
func StringLengthMax(max int) String {
	return func(s string) error {
		if len(s) > max {
			return fmt.Errorf("string should be at most %d characters long", max)
		}

		return nil
	}
}

// StringExcludeChars creates a String validator that fails when the string contains one ore more of the characters in `chars`.
func StringExcludeChars(chars string) String {
	return func(s string) error {
		if strings.ContainsAny(s, chars) {
			return fmt.Errorf("string should not contain any of %q", chars)
		}

		return nil
	}
}

// StringExcludePrefix creates a String validator that fails when the string starts with `prefix`.
func StringExcludePrefix(prefix string) String {
	return func(s string) error {
		if strings.HasPrefix(s, prefix) {
			return fmt.Errorf("string should not start with %q", prefix)
		}

		return nil
	}
}

// StringExcludeSuffix creates a String validator that fails when the string ends with `suffix`.
func StringExcludeSuffix(suffix string) String {
	return func(s string) error {
		if strings.HasSuffix(s, suffix) {
			return fmt.Errorf("string should not end with %q", suffix)
		}

		return nil
	}
}
