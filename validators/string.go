package validators

import (
	"fmt"
	"strings"
)

type String func(string) error

func StringNotEmpty() String {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("string should not be empty")
		}

		return nil
	}
}

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

func StringLengthMin(min int) String {
	return func(s string) error {
		if len(s) < min {
			return fmt.Errorf("string should be at least %d characters long", min)
		}

		return nil
	}
}

func StringLengthMax(max int) String {
	return func(s string) error {
		if len(s) > max {
			return fmt.Errorf("string should be at most %d characters long", max)
		}

		return nil
	}
}

func StringExcludeChars(chars string) String {
	return func(s string) error {
		if strings.ContainsAny(s, chars) {
			return fmt.Errorf("string should not contain any of %q", chars)
		}

		return nil
	}
}

func StringExcludePrefix(prefix string) String {
	return func(s string) error {
		if strings.HasPrefix(s, prefix) {
			return fmt.Errorf("string should not start with %q", prefix)
		}

		return nil
	}
}

func StringExcludeSuffix(suffix string) String {
	return func(s string) error {
		if strings.HasSuffix(s, suffix) {
			return fmt.Errorf("string should not end with %q", suffix)
		}

		return nil
	}
}
