package validators

import "regexp"

// A Regexp validator should return an error if the *regexp.Regexp provided is not considered valid, nil otherwise.
type Regexp func(*regexp.Regexp) error
