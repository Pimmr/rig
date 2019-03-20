package main

import (
	"fmt"

	"github.com/Pimmr/rig"
)

type bar struct {
	FlagE int     `flag:"flag-e" env:"FLAG_E" usage:"Flag E"`
	FlagF float64 `flag:"flag-f" env:"FLAG_F" usage:"Flag F"`
}

type foo struct {
	FlagA string         `flag:"flag-a" env:"FLAG_A" usage:"Flag A"`
	FlagB int            `flag:"flag-b" env:"FLAG_B" usage:"Flag B" required:"true"`
	FlagC []int          `flag:"flag-c" env:"FLAG_C" usage:"Flag C" typehint:"many ints"`
	FlagD []rig.URLValue `flag:"flag-d" usage:"Flag D"`
	Bar   bar            `flag:"bar" env:"BAR" required:"t"`
	Baz   bar
	Blah  bar `ignore:"true"`
}

func main() {
	var f foo

	err := rig.ParseStruct(&f)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", f)
}
