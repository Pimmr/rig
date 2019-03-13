package main

import (
	"fmt"

	"github.com/Pimmr/rig"
	"github.com/Pimmr/rig/structToFlags"
)

type bar struct {
	FlagE int     `rig-flag:"flag-e" rig-env:"FLAG_E" rig-usage:"Flag E"`
	FlagF float64 `rig-flag:"flag-f" rig-env:"FLAG_F" rig-usage:"Flag F"`
}

type foo struct {
	FlagA string         `rig-flag:"flag-a" rig-env:"FLAG_A" rig-usage:"Flag A"`
	FlagB int            `rig-flag:"flag-b" rig-env:"FLAG_B" rig-usage:"Flag B" rig-required:"true"`
	FlagC []int          `rig-flag:"flag-c" rig-env:"FLAG_C" rig-usage:"Flag C" rig-typehint:"many ints"`
	FlagD []rig.URLValue `rig-flag:"flag-d" rig-usage:"Flag D"`
	Bar   bar            `rig-flag:"bar" rig-env:"BAR"`
	Baz   bar
}

func main() {
	var f foo

	flags, err := structToFlags.StructToFlags(&f)
	if err != nil {
		panic(err)
	}

	err = rig.Parse(flags...)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", f)
}
