package main

import (
	"fmt"

	"github.com/Pimmr/rig"
)

type bar struct {
	FlagE int     `usage:"Flag E"`
	FlagF float64 `usage:"Flag F"`
}

type foo struct {
	FlagA string         `usage:"Flag A"`
	FlagB int            `flag:",require" usage:"Flag B"`
	FlagC []int          `usage:"Flag C" typehint:"many ints"`
	FlagD []rig.URLValue `usage:"Flag D"`
	Bar   bar            `flag:",require"`
	Baz   bar            `flag:",inline" env:",inline"`
	Blah  bar            `flag:"-"`
}

func main() {
	var f foo

	err := rig.ParseStruct(&f)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", f)
}
