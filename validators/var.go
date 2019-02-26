package validators

import "flag"

type Var func(flag.Value) error
