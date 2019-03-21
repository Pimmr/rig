package validators

import (
	"fmt"
	"net/url"
)

// A URL validator should return an error if the *url.URL provided is not considered valid, nil otherwise.
type URL func(*url.URL) error

// URLScheme creates a URL validator that fails when the url.URL does not use the scheme `scheme`.
// The validator never fails of `scheme` is empty.
func URLScheme(scheme string) URL {
	if scheme == "" {
		return func(*url.URL) error {
			return nil
		}
	}

	return func(u *url.URL) error {
		if u.Scheme != scheme {
			return fmt.Errorf("url should use %q scheme", scheme)
		}

		return nil
	}
}

// URLExcludeScheme creates a URL validator that fails when the url.URL uses the scheme `scheme`.
// The validator never fails of `scheme` is empty.
func URLExcludeScheme(scheme string) URL {
	if scheme == "" {
		return func(*url.URL) error {
			return nil
		}
	}

	return func(u *url.URL) error {
		if u.Scheme == scheme {
			return fmt.Errorf("url should not use %q scheme", scheme)
		}

		return nil
	}
}
