package validators

import "flag"

// A Var validator should return an error if the flag.Value provided is not considered valid, nil otherwise.
type Var func(flag.Value) error
