package validators

import (
	"flag"
	"net/url"
	"regexp"
	"time"
)

type Duration func(time.Duration) error
type Float64 func(float64) error
type Int func(int) error
type Int64 func(int64) error
type Regexp func(*regexp.Regexp) error
type String func(string) error
type Uint func(uint) error
type Uint64 func(uint64) error
type URL func(*url.URL) error
type Var func(flag.Value) error
