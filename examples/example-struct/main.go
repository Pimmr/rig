package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Pimmr/rig"
)

type bar struct {
	FlagE int     `usage:"Flag E" env:"-"`
	FlagF float64 `usage:"Flag F"`
}

type customType struct {
	A, B string
}

func (t customType) String() string {
	return t.A + "-" + t.B
}

func (t *customType) Set(s string) error {
	ss := strings.Split(s, "-")
	if len(ss) != 2 {
		return fmt.Errorf("malformed customTime %q: expected lhs-rhs", s)
	}

	t.A = ss[0]
	t.B = ss[1]

	return nil
}

type foo struct {
	FlagA string     `usage:"Flag A"`
	FlagB int        `flag:",require" usage:"Flag B"`
	FlagC []int      `usage:"Flag C" typehint:"many ints"`
	FlagD []*url.URL `usage:"Flag D"`
	FlagG customType
	Bar   bar `flag:",require"`
	Baz   bar `flag:",inline" env:",inline"`
	Blah  bar `flag:"-"`
	FlagH *int
}

func main() {
	var f foo

	err := rig.ParseStruct(&f)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", f)
	if f.FlagH != nil {
		fmt.Printf("FlagH = %d\n", *f.FlagH)
	}
}
