package validators

import "regexp"

type Regexp func(*regexp.Regexp) error
